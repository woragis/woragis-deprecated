// Load environment variables from .env file (server-side only)
// This is safe to import in Edge Runtime as it only reads process.env
try {
  // Only load dotenv if we're in a Node.js environment (not Edge Runtime)
  if (typeof process !== 'undefined' && !process.env.NEXT_RUNTIME) {
    const { config } = require("dotenv");
    config();
  }
} catch (error) {
  // Silently fail in Edge Runtime - environment variables should be set externally
  console.warn("Could not load .env file (this is expected in Edge Runtime)");
}

export const env = {
  // Database
  DATABASE_URL:
    process.env.DATABASE_URL ||
    "postgresql://postgres:password@localhost:5432/woragis",
  POSTGRES_DB: process.env.POSTGRES_DB || "woragis",
  POSTGRES_USER: process.env.POSTGRES_USER || "postgres",
  POSTGRES_PASSWORD: process.env.POSTGRES_PASSWORD || "password",
  POSTGRES_PORT: process.env.POSTGRES_PORT || "5432",

  // Redis
  REDIS_URL: process.env.REDIS_URL || "redis://localhost:6379/1",

  // Next.js
  NODE_ENV: process.env.NODE_ENV || "development",
  NEXT_PUBLIC_APP_URL:
    process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000",

  // App Configuration
  PORT: process.env.PORT || "3000",

  // Storage Configuration
  STORAGE_TYPE: (process.env.STORAGE_TYPE as "local" | "s3") || "local",
  
  // S3 Configuration (required in production)
  AWS_S3_BUCKET: process.env.AWS_S3_BUCKET,
  AWS_REGION: process.env.AWS_REGION || "us-east-1",
  AWS_ACCESS_KEY_ID: process.env.AWS_ACCESS_KEY_ID,
  AWS_SECRET_ACCESS_KEY: process.env.AWS_SECRET_ACCESS_KEY,
  AWS_S3_ENDPOINT: process.env.AWS_S3_ENDPOINT,
  
  // Upload Configuration
  UPLOAD_DIR: process.env.UPLOAD_DIR || "uploads",
  UPLOAD_PUBLIC_URL: process.env.UPLOAD_PUBLIC_URL || "http://localhost:3000/uploads",

  // AI Configuration
  OPENAI_API_KEY: process.env.OPENAI_API_KEY,

  // CORS Configuration
  ALLOWED_ORIGINS: process.env.ALLOWED_ORIGINS || "http://localhost:3000,http://localhost:3001",
  // Note: Native mobile apps don't need CORS origins - only web views and PWAs do
  MOBILE_APP_ORIGINS: process.env.MOBILE_APP_ORIGINS || "capacitor://localhost,ionic://localhost",

  // API Key Configuration
  API_KEY_SECRET: process.env.API_KEY_SECRET || "your-api-key-secret",
  MOBILE_API_KEYS: process.env.MOBILE_API_KEYS || "mobile-key-1,mobile-key-2",
} as const;

// Validate required environment variables
const requiredEnvVars = ["DATABASE_URL", "REDIS_URL", "OPENAI_API_KEY"] as const;

for (const envVar of requiredEnvVars) {
  if (!env[envVar]) {
    throw new Error(`Missing required environment variable: ${envVar}`);
  }
}

// Validate S3 configuration in production
if (env.NODE_ENV === "production") {
  if (env.STORAGE_TYPE === "s3") {
    const requiredS3Vars = [
      "AWS_S3_BUCKET",
      "AWS_ACCESS_KEY_ID", 
      "AWS_SECRET_ACCESS_KEY"
    ] as const;
    
    for (const envVar of requiredS3Vars) {
      if (!env[envVar]) {
        throw new Error(`Missing required S3 environment variable for production: ${envVar}`);
      }
    }
    
    console.log("✅ S3 configuration validated for production");
  } else {
    console.warn("⚠️  WARNING: Using local storage in production. Consider using S3 for better scalability.");
  }
}

// Validate URL formats
try {
  new URL(env.DATABASE_URL);
} catch {
  throw new Error(`Invalid DATABASE_URL format: ${env.DATABASE_URL}`);
}

try {
  new URL(env.REDIS_URL);
} catch {
  throw new Error(`Invalid REDIS_URL format: ${env.REDIS_URL}`);
}

// Production environment warnings
if (env.NODE_ENV === "production") {
  if (env.DATABASE_URL.includes("localhost")) {
    console.warn("⚠️  WARNING: Using localhost database in production environment");
  }
  if (env.REDIS_URL.includes("localhost")) {
    console.warn("⚠️  WARNING: Using localhost Redis in production environment");
  }
}

/**
 * Test database connection on startup
 */
async function testDatabaseConnection(): Promise<boolean> {
  try {
    const { testDatabaseConnection } = await import("./connection-tests");
    const result = await testDatabaseConnection();
    
    if (!result.success) {
      console.error("❌ Database connection failed:", result.error);
      return false;
    }
    
    console.log("✅ Database connection successful");
    return true;
  } catch (error) {
    console.error("❌ Database connection test failed:", error);
    return false;
  }
}

/**
 * Test Redis connection on startup
 */
async function testRedisConnection(): Promise<boolean> {
  try {
    const { testRedisConnection } = await import("./connection-tests");
    const result = await testRedisConnection();
    
    if (!result.success) {
      console.error("❌ Redis connection failed:", result.error);
      return false;
    }
    
    console.log("✅ Redis connection successful");
    return true;
  } catch (error) {
    console.error("❌ Redis connection test failed:", error);
    return false;
  }
}

/**
 * Test all connections and return results
 * This should be called on app startup (server-side only)
 * Note: This function should NOT be called in Edge Runtime (middleware)
 */
export async function validateConnectionsOnStartup(): Promise<{ success: boolean; errors: string[] }> {
  console.log("🔍 Testing connections on startup...");
  
  try {
    const { testAllConnections } = await import("./connection-tests");
    const results = await testAllConnections();
    
    if (!results.overall.success) {
      console.error("❌ Connection validation failed:");
      results.overall.errors.forEach(error => console.error(`  - ${error}`));
      
      console.error("\n💡 Make sure your services are running:");
      console.error("  - Database: docker-compose up db");
      console.error("  - Redis: docker-compose up redis");
      console.error("  - Or start all: docker-compose up");
      
      return { success: false, errors: results.overall.errors };
    }
    
    console.log("✅ All connections validated successfully");
    return { success: true, errors: [] };
  } catch (error) {
    console.error("❌ Connection validation error:", error);
    return { success: false, errors: [error instanceof Error ? error.message : 'Unknown error'] };
  }
}

export default env;
