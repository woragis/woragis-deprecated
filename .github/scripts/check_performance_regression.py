#!/usr/bin/env python3
"""
Performance Regression Detection Script

This script analyzes performance test results and detects regressions
by comparing current results with baseline metrics.
"""

import json
import os
import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple

# Performance baselines (in messages/jobs per second)
BASELINES = {
    "email-worker": {
        "throughput": 50.0,  # messages/second
        "latency_avg": 100.0,  # milliseconds
        "latency_p95": 200.0,  # milliseconds
    },
    "translation-worker": {
        "throughput": 10.0,  # jobs/second
        "latency_avg": 200.0,  # milliseconds
        "latency_p95": 500.0,  # milliseconds
    },
    "whatsapp-worker": {
        "throughput": 50.0,  # messages/second
        "latency_avg": 100.0,  # milliseconds
        "latency_p95": 200.0,  # milliseconds
    },
}

# Regression thresholds (percentage degradation)
THROUGHPUT_THRESHOLD = 0.20  # 20% degradation
LATENCY_THRESHOLD = 0.30  # 30% increase


def load_baselines() -> Dict[str, Dict[str, float]]:
    """Load baselines from file or use defaults."""
    workspace = Path(os.getenv("GITHUB_WORKSPACE", "."))
    baseline_file = workspace / ".github" / "performance-baselines.json"
    
    if baseline_file.exists():
        try:
            with open(baseline_file, 'r') as f:
                file_baselines = json.load(f)
                # Merge with defaults (file takes precedence)
                merged = BASELINES.copy()
                merged.update(file_baselines)
                return merged
        except Exception as e:
            print(f"⚠️ Failed to load baseline file: {e}")
            print("Using default baselines")
    
    return BASELINES


def parse_throughput(log_content: str) -> Optional[float]:
    """Extract throughput from test logs."""
    # Look for "Throughput: X.XX msg/s" or "Throughput: X.XX jobs/s"
    pattern = r"Throughput:\s+([\d.]+)\s+(?:msg|jobs)/s"
    match = re.search(pattern, log_content)
    if match:
        return float(match.group(1))
    return None


def parse_latency(log_content: str) -> Optional[Dict[str, float]]:
    """Extract latency metrics from test logs."""
    latencies = {}
    
    # Look for "Avg latency: X.XXms" or "Avg latency: X.XXms"
    avg_pattern = r"Avg latency:\s+([\d.]+)\s*(?:ms|µs|ns)"
    avg_match = re.search(avg_pattern, log_content)
    if avg_match:
        value = float(avg_match.group(1))
        # Convert to milliseconds if needed
        context_start = max(0, avg_match.start() - 20)
        context_end = min(len(log_content), avg_match.end() + 20)
        context = log_content[context_start:context_end]
        if "µs" in context:
            value = value / 1000
        elif "ns" in context:
            value = value / 1000000
        latencies["avg"] = value
    
    # Look for "P95 latency: X.XXms"
    p95_pattern = r"P95 latency:\s+([\d.]+)\s*(?:ms|µs|ns)"
    p95_match = re.search(p95_pattern, log_content)
    if p95_match:
        value = float(p95_match.group(1))
        context_start = max(0, p95_match.start() - 20)
        context_end = min(len(log_content), p95_match.end() + 20)
        context = log_content[context_start:context_end]
        if "µs" in context:
            value = value / 1000
        elif "ns" in context:
            value = value / 1000000
        latencies["p95"] = value
    
    return latencies if latencies else None


def parse_benchmark_results(log_content: str) -> Optional[Dict[str, float]]:
    """Extract benchmark results from Go benchmark output."""
    # Look for benchmark lines like:
    # BenchmarkEmailWorkerThroughput-8    1000    1234 ns/op    512 B/op    2 allocs/op
    pattern = r"Benchmark\w+Throughput[^\s]+\s+\d+\s+([\d.]+)\s+ns/op"
    match = re.search(pattern, log_content)
    if match:
        ns_per_op = float(match.group(1))
        # Convert to operations per second
        ops_per_sec = 1_000_000_000 / ns_per_op
        return {"throughput": ops_per_sec}
    return None


def check_regression(
    worker_name: str,
    current_throughput: Optional[float],
    current_latency: Optional[Dict[str, float]],
    baseline: Dict[str, float],
) -> List[str]:
    """Check for performance regressions."""
    regressions = []
    
    # Check throughput regression
    if current_throughput and "throughput" in baseline:
        baseline_throughput = baseline["throughput"]
        degradation = (baseline_throughput - current_throughput) / baseline_throughput
        
        if degradation > THROUGHPUT_THRESHOLD:
            regressions.append(
                f"⚠️ Throughput regression: {current_throughput:.2f} ops/s "
                f"(baseline: {baseline_throughput:.2f} ops/s, "
                f"degradation: {degradation*100:.1f}%)"
            )
        else:
            print(f"✅ Throughput OK: {current_throughput:.2f} ops/s "
                  f"(baseline: {baseline_throughput:.2f} ops/s)")
    
    # Check latency regression
    if current_latency and "latency_avg" in baseline:
        if "avg" in current_latency:
            baseline_latency = baseline["latency_avg"]
            current_avg = current_latency["avg"]
            increase = (current_avg - baseline_latency) / baseline_latency
            
            if increase > LATENCY_THRESHOLD:
                regressions.append(
                    f"⚠️ Average latency regression: {current_avg:.2f}ms "
                    f"(baseline: {baseline_latency:.2f}ms, "
                    f"increase: {increase*100:.1f}%)"
                )
            else:
                print(f"✅ Average latency OK: {current_avg:.2f}ms "
                      f"(baseline: {baseline_latency:.2f}ms)")
        
        if "p95" in current_latency and "latency_p95" in baseline:
            baseline_p95 = baseline["latency_p95"]
            current_p95 = current_latency["p95"]
            increase = (current_p95 - baseline_p95) / baseline_p95
            
            if increase > LATENCY_THRESHOLD:
                regressions.append(
                    f"⚠️ P95 latency regression: {current_p95:.2f}ms "
                    f"(baseline: {baseline_p95:.2f}ms, "
                    f"increase: {increase*100:.1f}%)"
                )
            else:
                print(f"✅ P95 latency OK: {current_p95:.2f}ms "
                      f"(baseline: {baseline_p95:.2f}ms)")
    
    return regressions


def analyze_worker_performance(worker_name: str, log_file: Path, baseline: Dict[str, float]) -> List[str]:
    """Analyze performance test results for a worker."""
    if not log_file.exists():
        print(f"⚠️ Performance log not found: {log_file}")
        return []
    
    with open(log_file, 'r') as f:
        log_content = f.read()
    
    regressions = []
    
    # Parse throughput
    throughput = parse_throughput(log_content)
    if not throughput:
        # Try parsing benchmark results
        benchmark_results = parse_benchmark_results(log_content)
        if benchmark_results:
            throughput = benchmark_results.get("throughput")
    
    # Parse latency
    latency = parse_latency(log_content)
    
    # Check for regressions
    worker_regressions = check_regression(
        worker_name, throughput, latency, baseline
    )
    regressions.extend(worker_regressions)
    
    return regressions


def main():
    """Main entry point."""
    # Load baselines
    baselines = load_baselines()
    
    # Find performance test result files
    workspace = Path(os.getenv("GITHUB_WORKSPACE", "."))
    result_files = {
        "email-worker": workspace / "backend" / "email-worker" / "email-worker-performance.txt",
        "translation-worker": workspace / "backend" / "translation-worker" / "translation-worker-performance.txt",
        "whatsapp-worker": workspace / "backend" / "whatsapp-worker" / "whatsapp-worker-performance.txt",
    }
    
    all_regressions = []
    
    for worker_name, log_file in result_files.items():
        print(f"\n📊 Analyzing {worker_name} performance...")
        baseline = baselines.get(worker_name, {})
        if not baseline:
            print(f"⚠️ No baseline defined for {worker_name}, skipping")
            continue
        
        regressions = analyze_worker_performance(worker_name, log_file, baseline)
        all_regressions.extend(regressions)
    
    # Output results
    if all_regressions:
        print("\n" + "="*60)
        print("⚠️ PERFORMANCE REGRESSIONS DETECTED")
        print("="*60)
        for regression in all_regressions:
            print(regression)
        print("="*60)
        
        # Write to GitHub Actions output
        if os.getenv("GITHUB_ACTIONS"):
            summary_file = os.getenv("GITHUB_STEP_SUMMARY", "/dev/stdout")
            with open(summary_file, "a") as f:
                f.write("\n## ⚠️ Performance Regressions Detected\n\n")
                for regression in all_regressions:
                    f.write(f"- {regression}\n")
                f.write("\n")
        
        sys.exit(1)
    else:
        print("\n✅ No performance regressions detected!")
        sys.exit(0)


if __name__ == "__main__":
    main()
