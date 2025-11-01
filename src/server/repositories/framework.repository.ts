import { eq, desc, asc, and, like, sql } from "drizzle-orm";
import { db } from "../db";
import { frameworks, projectFrameworks } from "../db/schemas";
import type {
  Framework,
  NewFramework,
  FrameworkFilters,
  FrameworkType,
  ProficiencyLevel,
} from "@/types";

export class FrameworkRepository {
  // Basic CRUD operations
  async findAll(userId?: string): Promise<Framework[]> {
    let query = db.select().from(frameworks) as any;
    if (userId) {
      query = query.where(eq(frameworks.userId, userId)) as any;
    }
    return await query.orderBy(asc(frameworks.order), asc(frameworks.name));
  }

  async findVisible(userId?: string): Promise<Framework[]> {
    const conditions = [eq(frameworks.visible, true)];
    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }
    return await db
      .select()
      .from(frameworks)
      .where(and(...conditions))
      .orderBy(asc(frameworks.order), asc(frameworks.name));
  }

  async findById(id: string, userId?: string): Promise<Framework | null> {
    const conditions = [eq(frameworks.id, id)];
    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }
    const result = await db
      .select()
      .from(frameworks)
      .where(and(...conditions));
    return result[0] || null;
  }

  async findBySlug(slug: string, userId?: string): Promise<Framework | null> {
    const conditions = [eq(frameworks.slug, slug)];
    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }
    const result = await db
      .select()
      .from(frameworks)
      .where(and(...conditions));
    return result[0] || null;
  }

  async create(framework: NewFramework): Promise<Framework> {
    const result = await db.insert(frameworks).values(framework).returning();
    return result[0];
  }

  async update(
    id: string,
    framework: Partial<NewFramework>,
    userId?: string
  ): Promise<Framework | null> {
    const conditions = [eq(frameworks.id, id)];
    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }
    const result = await db
      .update(frameworks)
      .set({ ...framework, updatedAt: new Date() })
      .where(and(...conditions))
      .returning();
    return result[0] || null;
  }

  async delete(id: string, userId?: string): Promise<void> {
    const conditions = [eq(frameworks.id, id)];
    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }
    await db.delete(frameworks).where(and(...conditions));
  }

  // Search and filtering
  async search(
    filters: FrameworkFilters,
    userId?: string
  ): Promise<Framework[]> {
    const conditions = [];

    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }

    if (filters.visible !== undefined) {
      conditions.push(eq(frameworks.visible, filters.visible));
    }

    if (filters.public !== undefined) {
      conditions.push(eq(frameworks.public, filters.public));
    }

    if (filters.type) {
      conditions.push(eq(frameworks.type, filters.type));
    }

    if (filters.search) {
      conditions.push(like(frameworks.name, `%${filters.search}%`));
    }

    let query = db.select().from(frameworks) as any;

    if (conditions.length > 0) {
      query = query.where(and(...conditions)) as any;
    }

    if (filters.limit) {
      query = query.limit(filters.limit) as any;
    }

    if (filters.offset) {
      query = query.offset(filters.offset) as any;
    }

    // Dynamic sorting
    const sortField = filters.sortBy || "order";
    const sortDirection = filters.sortOrder || "asc";
    
    // Map sort field to column
    const sortColumn = sortField === "name" ? frameworks.name :
                       sortField === "type" ? frameworks.type :
                       sortField === "createdAt" ? frameworks.createdAt :
                       sortField === "updatedAt" ? frameworks.updatedAt :
                       frameworks.order;
    
    if (sortDirection === "desc") {
      query = query.orderBy(desc(sortColumn)) as any;
    } else {
      query = query.orderBy(asc(sortColumn)) as any;
    }
    
    // Add secondary sort by name for consistency
    if (sortField !== "name") {
      query = query.orderBy(asc(frameworks.name)) as any;
    }

    return await query;
  }

  // Type-specific methods
  async findByType(type: FrameworkType, userId?: string): Promise<Framework[]> {
    const conditions = [eq(frameworks.type, type)];
    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }
    return await db
      .select()
      .from(frameworks)
      .where(and(...conditions))
      .orderBy(asc(frameworks.order), asc(frameworks.name));
  }

  async findVisibleByType(
    type: FrameworkType,
    userId?: string
  ): Promise<Framework[]> {
    const conditions = [
      eq(frameworks.type, type),
      eq(frameworks.visible, true),
    ];
    if (userId) {
      conditions.push(eq(frameworks.userId, userId));
    }
    return await db
      .select()
      .from(frameworks)
      .where(and(...conditions))
      .orderBy(asc(frameworks.order), asc(frameworks.name));
  }

  // Project relations
  async getProjectCount(frameworkId: string): Promise<number> {
    const result = await db
      .select({ count: sql<number>`count(*)` })
      .from(projectFrameworks)
      .where(eq(projectFrameworks.frameworkId, frameworkId));

    return result[0]?.count || 0;
  }

  async findWithProjectCount(
    frameworkId: string
  ): Promise<{ framework: Framework; projectCount: number } | null> {
    const framework = await this.findById(frameworkId);
    if (!framework) return null;

    const projectCount = await this.getProjectCount(frameworkId);
    return { framework, projectCount };
  }

  async findPopular(
    limit: number = 10
  ): Promise<Array<{ framework: Framework; projectCount: number }>> {
    const result = await db
      .select({
        framework: frameworks,
        projectCount: sql<number>`count(${projectFrameworks.id})`,
      })
      .from(frameworks)
      .leftJoin(
        projectFrameworks,
        eq(frameworks.id, projectFrameworks.frameworkId)
      )
      .groupBy(frameworks.id)
      .orderBy(desc(sql`count(${projectFrameworks.id})`))
      .limit(limit);

    return result.map(({ framework, projectCount }) => ({
      framework,
      projectCount,
    }));
  }

  // Version analysis
  async getVersionDistribution(
    frameworkId: string
  ): Promise<Array<{ version: string; count: number }>> {
    const result = await db
      .select({
        version: projectFrameworks.version,
        count: sql<number>`count(*)`,
      })
      .from(projectFrameworks)
      .where(eq(projectFrameworks.frameworkId, frameworkId))
      .groupBy(projectFrameworks.version)
      .orderBy(desc(sql`count(*)`));

    return result.map(({ version, count }) => ({
      version: version || "unknown",
      count,
    }));
  }

  // Proficiency analysis
  async getProficiencyDistribution(
    frameworkId: string
  ): Promise<Array<{ proficiency: string; count: number }>> {
    const result = await db
      .select({
        proficiency: projectFrameworks.proficiency,
        count: sql<number>`count(*)`,
      })
      .from(projectFrameworks)
      .where(eq(projectFrameworks.frameworkId, frameworkId))
      .groupBy(projectFrameworks.proficiency)
      .orderBy(desc(sql`count(*)`));

    return result.map(({ proficiency, count }) => ({
      proficiency: proficiency || "unknown",
      count,
    }));
  }

  // Ordering operations
  async updateOrder(
    frameworkOrders: { id: string; order: number }[]
  ): Promise<void> {
    const promises = frameworkOrders.map(({ id, order }) =>
      db
        .update(frameworks)
        .set({ order, updatedAt: new Date() })
        .where(eq(frameworks.id, id))
    );
    await Promise.all(promises);
  }
}
