ALTER TABLE "testimonials" ALTER COLUMN "featured" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "testimonials" ALTER COLUMN "visible" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "testimonials" ALTER COLUMN "public" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "testimonials" ALTER COLUMN "created_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "testimonials" ALTER COLUMN "updated_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_posts" ALTER COLUMN "visible" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_posts" ALTER COLUMN "public" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_posts" ALTER COLUMN "view_count" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_posts" ALTER COLUMN "like_count" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_posts" ALTER COLUMN "created_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_posts" ALTER COLUMN "updated_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "experiences" ALTER COLUMN "visible" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "experiences" ALTER COLUMN "created_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "experiences" ALTER COLUMN "updated_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_tags" ALTER COLUMN "visible" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_tags" ALTER COLUMN "created_at" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "blog_tags" ALTER COLUMN "updated_at" SET NOT NULL;