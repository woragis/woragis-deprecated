/**
 * Startup Validation Module
 * 
 * This module provides functions to validate the application state on startup.
 * It should be imported and called early in the application lifecycle.
 */

import { validateConnectionsOnStartup } from "./env";

/**
 * Comprehensive startup validation
 * This function should be called at the very beginning of your application
 * Note: This function will exit the process if validation fails (server-side only)
 */
export async function validateApplicationStartup(): Promise<void> {
  console.log("🚀 Starting application validation...");
  
  try {
    // Test all critical connections
    const result = await validateConnectionsOnStartup();
    
    if (!result.success) {
      console.error("❌ Application startup validation failed:");
      result.errors.forEach(error => console.error(`  - ${error}`));
      
      console.error("\n💡 Make sure your services are running:");
      console.error("  - Database: docker-compose up db");
      console.error("  - Redis: docker-compose up redis");
      console.error("  - Or start all: docker-compose up");
      
      // Only exit in server environment, not in Edge Runtime
      if (typeof process !== 'undefined' && process.exit) {
        process.exit(1);
      } else {
        throw new Error(`Application startup validation failed: ${result.errors.join(', ')}`);
      }
    }
    
    console.log("✅ Application startup validation completed successfully");
  } catch (error) {
    console.error("❌ Application startup validation failed:", error);
    
    // Only exit in server environment, not in Edge Runtime
    if (typeof process !== 'undefined' && process.exit) {
      process.exit(1);
    } else {
      throw error;
    }
  }
}

/**
 * Lightweight startup validation (for development)
 * This version doesn't exit the process on failure, just logs warnings
 */
export async function validateApplicationStartupDev(): Promise<boolean> {
  console.log("🚀 Starting development application validation...");
  
  try {
    const { testAllConnections } = await import("./connection-tests");
    const results = await testAllConnections();
    
    if (!results.overall.success) {
      console.warn("⚠️  Some connections failed (development mode):");
      results.overall.errors.forEach(error => console.warn(`  - ${error}`));
      console.warn("💡 Start services with: docker-compose up db redis");
      return false;
    }
    
    console.log("✅ Development application validation completed successfully");
    return true;
  } catch (error) {
    console.warn("⚠️  Development application validation failed:", error);
    return false;
  }
}

/**
 * Health check endpoint data
 * This can be used for health check endpoints
 */
export async function getHealthCheckData(): Promise<{
  status: "healthy" | "unhealthy";
  timestamp: string;
  services: {
    database: { status: "up" | "down"; error?: string };
    redis: { status: "up" | "down"; error?: string };
  };
}> {
  const { testAllConnections } = await import("./connection-tests");
  const results = await testAllConnections();
  
  return {
    status: results.overall.success ? "healthy" : "unhealthy",
    timestamp: new Date().toISOString(),
    services: {
      database: {
        status: results.database.success ? "up" : "down",
        error: results.database.success ? undefined : results.database.error,
      },
      redis: {
        status: results.redis.success ? "up" : "down",
        error: results.redis.success ? undefined : results.redis.error,
      },
    },
  };
}
