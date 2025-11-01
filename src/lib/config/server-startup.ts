/**
 * Server-side Startup Validation
 * 
 * This module provides server-side validation that can be used in:
 * - Server components
 * - API routes
 * - Build-time validation
 * 
 * Note: This should NOT be used in middleware (Edge Runtime)
 */

import { validateConnectionsOnStartup } from "./env";

let startupValidated = false;
let startupValidationPromise: Promise<void> | null = null;

/**
 * Validate startup once per server instance
 * This prevents multiple validation calls during the same server session
 * This function will exit the process if validation fails (server-side only)
 */
export async function validateServerStartup(): Promise<void> {
  if (startupValidated) {
    return;
  }

  if (startupValidationPromise) {
    return startupValidationPromise;
  }

  startupValidationPromise = (async () => {
    const result = await validateConnectionsOnStartup();
    
    if (!result.success) {
      console.error("❌ Server startup validation failed:");
      result.errors.forEach(error => console.error(`  - ${error}`));
      
      console.error("\n💡 Make sure your services are running:");
      console.error("  - Database: docker-compose up db");
      console.error("  - Redis: docker-compose up redis");
      console.error("  - Or start all: docker-compose up");
      
      // Only exit in server environment, not in Edge Runtime
      if (typeof process !== 'undefined' && process.exit) {
        process.exit(1);
      } else {
        throw new Error(`Startup validation failed: ${result.errors.join(', ')}`);
      }
    }
    
    startupValidated = true;
  })();
  
  try {
    await startupValidationPromise;
  } catch (error) {
    startupValidationPromise = null; // Reset on failure
    throw error;
  }
}

/**
 * Check if startup has been validated
 */
export function isStartupValidated(): boolean {
  return startupValidated;
}

/**
 * Reset startup validation (useful for testing)
 */
export function resetStartupValidation(): void {
  startupValidated = false;
  startupValidationPromise = null;
}
