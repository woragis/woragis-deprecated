CREATE TYPE "public"."proficiency_level" AS ENUM('beginner', 'intermediate', 'advanced', 'expert');--> statement-breakpoint
ALTER TYPE "public"."framework_type" ADD VALUE 'library' BEFORE 'tool';--> statement-breakpoint
ALTER TYPE "public"."framework_type" ADD VALUE 'database';--> statement-breakpoint
ALTER TYPE "public"."framework_type" ADD VALUE 'other';--> statement-breakpoint
ALTER TABLE "frameworks" ALTER COLUMN "order" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "frameworks" ALTER COLUMN "visible" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "frameworks" ALTER COLUMN "created_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "frameworks" ALTER COLUMN "updated_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "frameworks" ADD COLUMN "proficiency_level" "proficiency_level";--> statement-breakpoint
ALTER TABLE "frameworks" ADD COLUMN "public" boolean DEFAULT true NOT NULL;