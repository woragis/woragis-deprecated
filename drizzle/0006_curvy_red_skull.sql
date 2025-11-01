CREATE TABLE "ideas" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"user_id" uuid NOT NULL,
	"title" text NOT NULL,
	"slug" text NOT NULL,
	"document" text NOT NULL,
	"description" text,
	"featured" boolean DEFAULT false,
	"visible" boolean DEFAULT true,
	"public" boolean DEFAULT false,
	"order" integer DEFAULT 0 NOT NULL,
	"created_at" timestamp DEFAULT now(),
	"updated_at" timestamp DEFAULT now(),
	CONSTRAINT "ideas_slug_unique" UNIQUE("slug")
);
--> statement-breakpoint
CREATE TABLE "idea_nodes" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"idea_id" uuid NOT NULL,
	"title" text NOT NULL,
	"content" text,
	"type" text DEFAULT 'default',
	"position_x" real DEFAULT 0 NOT NULL,
	"position_y" real DEFAULT 0 NOT NULL,
	"width" real DEFAULT 200,
	"height" real DEFAULT 100,
	"color" text,
	"connections" jsonb DEFAULT '[]'::jsonb,
	"visible" boolean DEFAULT true,
	"created_at" timestamp DEFAULT now(),
	"updated_at" timestamp DEFAULT now()
);
--> statement-breakpoint
CREATE TABLE "ai_chats" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"user_id" uuid NOT NULL,
	"idea_node_id" uuid NOT NULL,
	"title" text NOT NULL,
	"messages" jsonb DEFAULT '[]'::jsonb NOT NULL,
	"agent" text DEFAULT 'gpt-4' NOT NULL,
	"model" text,
	"system_prompt" text,
	"temperature" text DEFAULT '0.7',
	"max_tokens" text,
	"visible" boolean DEFAULT true,
	"archived" boolean DEFAULT false,
	"created_at" timestamp DEFAULT now(),
	"updated_at" timestamp DEFAULT now()
);
--> statement-breakpoint
CREATE TABLE "blog_post_translations" (
	"blog_post_id" uuid NOT NULL,
	"language_code" text NOT NULL,
	"title" text NOT NULL,
	"excerpt" text NOT NULL,
	"content" text NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "blog_post_translations_blog_post_id_language_code_pk" PRIMARY KEY("blog_post_id","language_code")
);
--> statement-breakpoint
CREATE TABLE "education_translations" (
	"education_id" uuid NOT NULL,
	"language_code" text NOT NULL,
	"degree" text NOT NULL,
	"school" text NOT NULL,
	"year" text NOT NULL,
	"description" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "education_translations_education_id_language_code_pk" PRIMARY KEY("education_id","language_code")
);
--> statement-breakpoint
CREATE TABLE "experience_translations" (
	"experience_id" uuid NOT NULL,
	"language_code" text NOT NULL,
	"title" text NOT NULL,
	"company" text NOT NULL,
	"period" text NOT NULL,
	"location" text NOT NULL,
	"description" text NOT NULL,
	"achievements" text NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "experience_translations_experience_id_language_code_pk" PRIMARY KEY("experience_id","language_code")
);
--> statement-breakpoint
CREATE TABLE "project_translations" (
	"project_id" uuid NOT NULL,
	"language_code" text NOT NULL,
	"title" text NOT NULL,
	"description" text NOT NULL,
	"long_description" text,
	"content" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "project_translations_project_id_language_code_pk" PRIMARY KEY("project_id","language_code")
);
--> statement-breakpoint
CREATE TABLE "testimonial_translations" (
	"testimonial_id" uuid NOT NULL,
	"language_code" text NOT NULL,
	"name" text NOT NULL,
	"role" text NOT NULL,
	"company" text NOT NULL,
	"content" text NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "testimonial_translations_testimonial_id_language_code_pk" PRIMARY KEY("testimonial_id","language_code")
);
--> statement-breakpoint
ALTER TABLE "ideas" ADD CONSTRAINT "ideas_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "idea_nodes" ADD CONSTRAINT "idea_nodes_idea_id_ideas_id_fk" FOREIGN KEY ("idea_id") REFERENCES "public"."ideas"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "ai_chats" ADD CONSTRAINT "ai_chats_user_id_users_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "ai_chats" ADD CONSTRAINT "ai_chats_idea_node_id_idea_nodes_id_fk" FOREIGN KEY ("idea_node_id") REFERENCES "public"."idea_nodes"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "blog_post_translations" ADD CONSTRAINT "blog_post_translations_blog_post_id_blog_posts_id_fk" FOREIGN KEY ("blog_post_id") REFERENCES "public"."blog_posts"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "education_translations" ADD CONSTRAINT "education_translations_education_id_education_id_fk" FOREIGN KEY ("education_id") REFERENCES "public"."education"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "experience_translations" ADD CONSTRAINT "experience_translations_experience_id_experiences_id_fk" FOREIGN KEY ("experience_id") REFERENCES "public"."experiences"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "project_translations" ADD CONSTRAINT "project_translations_project_id_projects_id_fk" FOREIGN KEY ("project_id") REFERENCES "public"."projects"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "testimonial_translations" ADD CONSTRAINT "testimonial_translations_testimonial_id_testimonials_id_fk" FOREIGN KEY ("testimonial_id") REFERENCES "public"."testimonials"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
CREATE INDEX "blog_post_translations_blog_post_id_idx" ON "blog_post_translations" USING btree ("blog_post_id");--> statement-breakpoint
CREATE INDEX "blog_post_translations_language_code_idx" ON "blog_post_translations" USING btree ("language_code");--> statement-breakpoint
CREATE INDEX "education_translations_education_id_idx" ON "education_translations" USING btree ("education_id");--> statement-breakpoint
CREATE INDEX "education_translations_language_code_idx" ON "education_translations" USING btree ("language_code");--> statement-breakpoint
CREATE INDEX "experience_translations_experience_id_idx" ON "experience_translations" USING btree ("experience_id");--> statement-breakpoint
CREATE INDEX "experience_translations_language_code_idx" ON "experience_translations" USING btree ("language_code");--> statement-breakpoint
CREATE INDEX "project_translations_project_id_idx" ON "project_translations" USING btree ("project_id");--> statement-breakpoint
CREATE INDEX "project_translations_language_code_idx" ON "project_translations" USING btree ("language_code");--> statement-breakpoint
CREATE INDEX "testimonial_translations_testimonial_id_idx" ON "testimonial_translations" USING btree ("testimonial_id");--> statement-breakpoint
CREATE INDEX "testimonial_translations_language_code_idx" ON "testimonial_translations" USING btree ("language_code");