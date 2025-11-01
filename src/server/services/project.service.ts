import { projectRepository } from "@/server/repositories";
import type {
  Project,
  NewProject,
  ProjectFilters,
  ProjectOrderUpdate,
  ProjectWithRelations,
  ApiResponse,
} from "@/types";
import { BaseService } from "./base.service";

export class ProjectService extends BaseService {
  async getAllProjects(userId?: string): Promise<ApiResponse<Project[]>> {
    try {
      const projects = await projectRepository.findAll(userId);
      return this.success(projects);
    } catch (error) {
      return this.handleError(error, "getAllProjects");
    }
  }

  async getVisibleProjects(userId?: string): Promise<ApiResponse<Project[]>> {
    try {
      const projects = await projectRepository.findVisible(userId);
      return this.success(projects);
    } catch (error) {
      return this.handleError(error, "getVisibleProjects");
    }
  }

  async getPublicProjects(): Promise<ApiResponse<Project[]>> {
    try {
      const projects = await projectRepository.findPublic();
      return this.success(projects);
    } catch (error) {
      return this.handleError(error, "getPublicProjects");
    }
  }

  async getFeaturedProjects(
    limit: number = 3,
    userId?: string
  ): Promise<ApiResponse<Project[]>> {
    try {
      const projects = await projectRepository.findFeatured(limit, userId);
      return this.success(projects);
    } catch (error) {
      return this.handleError(error, "getFeaturedProjects");
    }
  }

  async getPublicFeaturedProjects(
    limit: number = 3
  ): Promise<ApiResponse<Project[]>> {
    try {
      const projects = await projectRepository.findPublicFeatured(limit);
      return this.success(projects);
    } catch (error) {
      return this.handleError(error, "getPublicFeaturedProjects");
    }
  }

  async getProjectById(
    id: string,
    userId?: string
  ): Promise<ApiResponse<Project | null>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      const project = await projectRepository.findById(id, userId);
      return this.success(project);
    } catch (error) {
      return this.handleError(error, "getProjectById");
    }
  }

  async getProjectBySlug(
    slug: string,
    userId?: string
  ): Promise<ApiResponse<Project | null>> {
    try {
      if (!slug || slug.trim().length === 0) {
        return {
          success: false,
          error: "Invalid project slug",
        };
      }

      const project = await projectRepository.findBySlug(slug, userId);
      return this.success(project);
    } catch (error) {
      return this.handleError(error, "getProjectBySlug");
    }
  }

  async getProjectWithRelations(
    id: string,
    userId?: string
  ): Promise<ApiResponse<ProjectWithRelations | null>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      const project = await projectRepository.findWithRelations(id, userId);
      return this.success(project);
    } catch (error) {
      return this.handleError(error, "getProjectWithRelations");
    }
  }

  async getPublicProjectById(id: string): Promise<ApiResponse<Project | null>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      const project = await projectRepository.findPublicById(id);
      return this.success(project);
    } catch (error) {
      return this.handleError(error, "getPublicProjectById");
    }
  }

  async getPublicProjectWithRelations(
    id: string
  ): Promise<ApiResponse<ProjectWithRelations | null>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      const project = await projectRepository.findPublicWithRelations(id);
      return this.success(project);
    } catch (error) {
      return this.handleError(error, "getPublicProjectWithRelations");
    }
  }

  async createProject(
    projectData: NewProject,
    userId: string
  ): Promise<ApiResponse<Project>> {
    try {
      const requiredFields: (keyof NewProject)[] = [
        "title",
        "description",
        "slug",
      ];
      const validationErrors = this.validateRequired(
        projectData,
        requiredFields
      );

      if (validationErrors.length > 0) {
        return {
          success: false,
          error: `Validation failed: ${validationErrors.join(", ")}`,
        };
      }

      // Add userId to project data
      const projectWithUser = { ...projectData, userId };
      const project = await projectRepository.create(projectWithUser);
      return this.success(project, "Project created successfully");
    } catch (error) {
      return this.handleError(error, "createProject");
    }
  }

  async updateProject(
    id: string,
    projectData: Partial<NewProject>,
    userId?: string
  ): Promise<ApiResponse<Project | null>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      const project = await projectRepository.update(id, projectData, userId);
      if (!project) {
        return {
          success: false,
          error: "Project not found",
        };
      }

      return this.success(project, "Project updated successfully");
    } catch (error) {
      return this.handleError(error, "updateProject");
    }
  }

  async deleteProject(id: string, userId?: string): Promise<ApiResponse<void>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      await projectRepository.delete(id, userId);
      return this.success(undefined, "Project deleted successfully");
    } catch (error) {
      return this.handleError(error, "deleteProject");
    }
  }

  async searchProjects(
    filters: ProjectFilters,
    userId?: string
  ): Promise<ApiResponse<Project[]>> {
    try {
      const projects = await projectRepository.search(filters, userId);
      return this.success(projects);
    } catch (error) {
      return this.handleError(error, "searchProjects");
    }
  }

  async updateProjectOrder(
    projectOrders: ProjectOrderUpdate[]
  ): Promise<ApiResponse<void>> {
    try {
      if (!projectOrders || projectOrders.length === 0) {
        return {
          success: false,
          error: "No project orders provided",
        };
      }

      await projectRepository.updateOrder(projectOrders);
      return this.success(undefined, "Project order updated successfully");
    } catch (error) {
      return this.handleError(error, "updateProjectOrder");
    }
  }

  async toggleProjectVisibility(
    id: string
  ): Promise<ApiResponse<Project | null>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      const project = await projectRepository.toggleVisibility(id);
      if (!project) {
        return {
          success: false,
          error: "Project not found",
        };
      }

      return this.success(project, "Project visibility toggled successfully");
    } catch (error) {
      return this.handleError(error, "toggleProjectVisibility");
    }
  }

  async toggleProjectFeatured(
    id: string
  ): Promise<ApiResponse<Project | null>> {
    try {
      if (!this.validateId(id)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      const project = await projectRepository.toggleFeatured(id);
      if (!project) {
        return {
          success: false,
          error: "Project not found",
        };
      }

      return this.success(
        project,
        "Project featured status toggled successfully"
      );
    } catch (error) {
      return this.handleError(error, "toggleProjectFeatured");
    }
  }



  // Framework relations
  async assignFrameworksToProject(
    projectId: string,
    frameworkAssignments: {
      frameworkId: string;
      version?: string;
      proficiency?: string;
    }[]
  ): Promise<ApiResponse<void>> {
    try {
      if (!this.validateId(projectId)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      await projectRepository.assignFrameworks(projectId, frameworkAssignments);
      return this.success(undefined, "Frameworks assigned successfully");
    } catch (error) {
      return this.handleError(error, "assignFrameworksToProject");
    }
  }

  async updateProjectFrameworks(
    projectId: string,
    frameworkAssignments: {
      frameworkId: string;
      version?: string;
      proficiency?: string;
    }[]
  ): Promise<ApiResponse<void>> {
    try {
      if (!this.validateId(projectId)) {
        return {
          success: false,
          error: "Invalid project ID",
        };
      }

      await projectRepository.updateFrameworks(projectId, frameworkAssignments);
      return this.success(undefined, "Project frameworks updated successfully");
    } catch (error) {
      return this.handleError(error, "updateProjectFrameworks");
    }
  }
}
