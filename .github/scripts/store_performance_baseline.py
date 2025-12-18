#!/usr/bin/env python3
"""
Store Performance Baseline Script

This script stores performance test results as baselines for regression detection.
Run this after confirming good performance test results.
"""

import json
import os
import re
import sys
from pathlib import Path
from typing import Dict, Optional


def parse_throughput(log_content: str) -> Optional[float]:
    """Extract throughput from test logs."""
    pattern = r"Throughput:\s+([\d.]+)\s+(?:msg|jobs)/s"
    match = re.search(pattern, log_content)
    if match:
        return float(match.group(1))
    return None


def parse_latency(log_content: str) -> Optional[Dict[str, float]]:
    """Extract latency metrics from test logs."""
    latencies = {}
    
    avg_pattern = r"Avg latency:\s+([\d.]+)\s*(?:ms|µs|ns)"
    avg_match = re.search(avg_pattern, log_content)
    if avg_match:
        value = float(avg_match.group(1))
        context_start = max(0, avg_match.start() - 20)
        context_end = min(len(log_content), avg_match.end() + 20)
        context = log_content[context_start:context_end]
        if "µs" in context:
            value = value / 1000
        elif "ns" in context:
            value = value / 1000000
        latencies["avg"] = value
    
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
    pattern = r"Benchmark\w+Throughput[^\s]+\s+\d+\s+([\d.]+)\s+ns/op"
    match = re.search(pattern, log_content)
    if match:
        ns_per_op = float(match.group(1))
        ops_per_sec = 1_000_000_000 / ns_per_op
        return {"throughput": ops_per_sec}
    return None


def main():
    """Main entry point."""
    workspace = Path(os.getenv("GITHUB_WORKSPACE", "."))
    baseline_file = workspace / ".github" / "performance-baselines.json"
    
    result_files = {
        "email-worker": workspace / "backend" / "email-worker" / "email-worker-performance.txt",
        "translation-worker": workspace / "backend" / "translation-worker" / "translation-worker-performance.txt",
        "whatsapp-worker": workspace / "backend" / "whatsapp-worker" / "whatsapp-worker-performance.txt",
    }
    
    baselines = {}
    
    for worker_name, log_file in result_files.items():
        if not log_file.exists():
            print(f"⚠️ Performance log not found: {log_file}")
            continue
        
        with open(log_file, 'r') as f:
            log_content = f.read()
        
        baseline = {}
        
        # Parse throughput
        throughput = parse_throughput(log_content)
        if not throughput:
            benchmark_results = parse_benchmark_results(log_content)
            if benchmark_results:
                throughput = benchmark_results.get("throughput")
        
        if throughput:
            baseline["throughput"] = throughput
        
        # Parse latency
        latency = parse_latency(log_content)
        if latency:
            if "avg" in latency:
                baseline["latency_avg"] = latency["avg"]
            if "p95" in latency:
                baseline["latency_p95"] = latency["p95"]
        
        if baseline:
            baselines[worker_name] = baseline
            print(f"✅ Stored baseline for {worker_name}: {baseline}")
    
    if baselines:
        baseline_file.parent.mkdir(parents=True, exist_ok=True)
        with open(baseline_file, 'w') as f:
            json.dump(baselines, f, indent=2)
        print(f"\n✅ Baselines saved to {baseline_file}")
    else:
        print("⚠️ No baselines to save")
        sys.exit(1)


if __name__ == "__main__":
    main()
