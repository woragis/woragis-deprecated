import { apiClient } from '@clients/apiClient';
import type {
	ArchitectureDiagramType,
	DocumentationSection,
	DocumentationSectionType,
	DocumentationVisibility,
	KanbanBoard,
	Milestone,
	Project,
	ProjectArchitectureDiagram,
	ProjectDependency,
	ProjectFileStructure,
	ProjectStatus,
	ProjectTechnology,
	TechnologyCategory,
	UUID
} from './types';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

export interface CreateProjectInput {
	name: string;
	description: string;
	status?: ProjectStatus;
	healthScore?: number;
	mrr?: number;
	cac?: number;
	ltv?: number;
	churnRate?: number;
}

export async function listProjects(): Promise<Project[]> {
	const response = await apiClient.get<ApiResponse<Project[]>>('/projects');
	return response.data.data ?? [];
}

export async function getProjectBySlug(slug: string): Promise<Project> {
	const response = await apiClient.get<ApiResponse<Project>>(`/projects/slug/${slug}`);
	return response.data.data;
}

export async function createProject(input: CreateProjectInput): Promise<Project> {
	const response = await apiClient.post<ApiResponse<Project>>('/projects', {
		name: input.name,
		description: input.description,
		status: input.status ?? 'planning',
		health_score: input.healthScore ?? 50,
		mrr: input.mrr ?? 0,
		cac: input.cac ?? 0,
		ltv: input.ltv ?? 0,
		churn_rate: input.churnRate ?? 0
	});
	return response.data.data;
}

export async function updateProjectStatus(
	projectId: UUID,
	status: ProjectStatus
): Promise<Project> {
	const response = await apiClient.patch<ApiResponse<Project>>(`/projects/${projectId}/status`, {
		status
	});
	return response.data.data;
}

export async function updateProjectMetrics(
	projectId: UUID,
	payload: Partial<{
		healthScore: number;
		mrr: number;
		cac: number;
		ltv: number;
		churnRate: number;
	}>
): Promise<Project> {
	const response = await apiClient.patch<ApiResponse<Project>>(`/projects/${projectId}/metrics`, {
		health_score: payload.healthScore,
		mrr: payload.mrr,
		cac: payload.cac,
		ltv: payload.ltv,
		churn_rate: payload.churnRate
	});
	return response.data.data;
}

export async function addMilestone(
	projectId: UUID,
	payload: { title: string; description: string; dueDate: string }
): Promise<Milestone> {
	const response = await apiClient.post<ApiResponse<Milestone>>(
		`/projects/${projectId}/milestones`,
		{
			title: payload.title,
			description: payload.description,
			due_date: payload.dueDate
		}
	);
	return response.data.data;
}

export async function listMilestones(projectId: UUID): Promise<Milestone[]> {
	const response = await apiClient.get<ApiResponse<Milestone[]>>(
		`/projects/${projectId}/milestones`
	);
	return response.data.data ?? [];
}

export async function bulkUpdateMilestones(
	projectId: UUID,
	updates: Array<{
		milestoneId: UUID;
		title?: string;
		description?: string;
		dueDate?: string;
		completed?: boolean;
	}>
): Promise<Milestone[]> {
	const payload = updates.map((item) => ({
		milestone_id: item.milestoneId,
		title: item.title,
		description: item.description,
		due_date: item.dueDate,
		completed: item.completed
	}));

	const response = await apiClient.post<ApiResponse<Milestone[]>>(
		`/projects/${projectId}/milestones/bulk`,
		{
			updates: payload
		}
	);
	return response.data.data ?? [];
}

export async function toggleMilestone(milestoneId: UUID, completed: boolean): Promise<Milestone> {
	const response = await apiClient.patch<ApiResponse<Milestone>>(
		`/projects/milestones/${milestoneId}`,
		{
			completed
		}
	);
	return response.data.data;
}

export async function getKanbanBoard(projectId: UUID): Promise<KanbanBoard> {
	const response = await apiClient.get<ApiResponse<KanbanBoard>>(`/projects/${projectId}/kanban`);
	return response.data.data;
}

export async function createKanbanColumn(
	projectId: UUID,
	payload: { name: string; position?: number; wipLimit?: number }
): Promise<KanbanBoard> {
	const response = await apiClient.post<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns`,
		{
			name: payload.name,
			position: payload.position,
			wip_limit: payload.wipLimit
		}
	);
	return response.data.data;
}

export async function updateKanbanColumn(
	projectId: UUID,
	columnId: UUID,
	payload: { name?: string; wipLimit?: number }
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns/${columnId}`,
		{
			name: payload.name,
			wip_limit: payload.wipLimit
		}
	);
	return response.data.data;
}

export async function reorderKanbanColumns(
	projectId: UUID,
	columnOrder: UUID[]
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns/reorder`,
		{
			column_order: columnOrder
		}
	);
	return response.data.data;
}

export async function deleteKanbanColumn(projectId: UUID, columnId: UUID): Promise<KanbanBoard> {
	const response = await apiClient.delete<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/columns/${columnId}`,
		{
			data: {}
		}
	);
	return response.data.data;
}

export async function createKanbanCard(
	projectId: UUID,
	payload: {
		columnId: UUID;
		title: string;
		description?: string;
		dueDate?: string;
		milestoneId?: UUID;
		position?: number;
	}
): Promise<KanbanBoard> {
	const response = await apiClient.post<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards`,
		{
			column_id: payload.columnId,
			title: payload.title,
			description: payload.description,
			due_date: payload.dueDate,
			milestone_id: payload.milestoneId,
			position: payload.position
		}
	);
	return response.data.data;
}

export async function updateKanbanCard(
	projectId: UUID,
	cardId: UUID,
	payload: {
		title?: string;
		description?: string;
		dueDate?: string;
		milestoneId?: UUID;
	}
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards/${cardId}`,
		{
			title: payload.title,
			description: payload.description,
			due_date: payload.dueDate,
			milestone_id: payload.milestoneId
		}
	);
	return response.data.data;
}

export async function moveKanbanCard(
	projectId: UUID,
	cardId: UUID,
	targetColumnId: UUID,
	targetPosition: number
): Promise<KanbanBoard> {
	const response = await apiClient.patch<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards/${cardId}/move`,
		{
			target_column_id: targetColumnId,
			target_position: targetPosition
		}
	);
	return response.data.data;
}

export async function deleteKanbanCard(projectId: UUID, cardId: UUID): Promise<KanbanBoard> {
	const response = await apiClient.delete<ApiResponse<KanbanBoard>>(
		`/projects/${projectId}/kanban/cards/${cardId}`,
		{
			data: {}
		}
	);
	return response.data.data;
}

export async function createDependency(
	projectId: UUID,
	payload: { dependsOnProjectId: UUID; type: 'blocks' | 'relates' | 'supports' }
): Promise<ProjectDependency> {
	const response = await apiClient.post<ApiResponse<ProjectDependency>>(
		`/projects/${projectId}/dependencies`,
		{
			depends_on_project_id: payload.dependsOnProjectId,
			type: payload.type
		}
	);
	return response.data.data;
}

export async function listDependencies(projectId: UUID): Promise<ProjectDependency[]> {
	const response = await apiClient.get<ApiResponse<ProjectDependency[]>>(
		`/projects/${projectId}/dependencies`
	);
	return response.data.data ?? [];
}

export async function deleteDependency(projectId: UUID, dependencyId: UUID): Promise<void> {
	await apiClient.delete(`/projects/${projectId}/dependencies/${dependencyId}`, {
		data: {}
	});
}

export async function duplicateProject(
	templateProjectId: UUID,
	payload: {
		name: string;
		description?: string;
		status?: ProjectStatus;
		healthScore?: number;
		mrr?: number;
		cac?: number;
		ltv?: number;
		churnRate?: number;
		copyBoard?: boolean;
		copyMilestones?: boolean;
		copyDependencies?: boolean;
	}
): Promise<Project> {
	const response = await apiClient.post<ApiResponse<Project>>(
		`/projects/${templateProjectId}/duplicate`,
		{
			name: payload.name,
			description: payload.description,
			status: payload.status,
			health_score: payload.healthScore,
			mrr: payload.mrr,
			cac: payload.cac,
			ltv: payload.ltv,
			churn_rate: payload.churnRate,
			copy_board: payload.copyBoard ?? true,
			copy_milestones: payload.copyMilestones ?? true,
			copy_dependencies: payload.copyDependencies ?? false
		}
	);
	return response.data.data;
}

// Documentation API

export async function createDocumentation(
	projectId: UUID,
	payload: { visibility?: DocumentationVisibility }
): Promise<ProjectDocumentation> {
	const response = await apiClient.post<ApiResponse<ProjectDocumentation>>(
		`/projects/${projectId}/documentation`,
		{
			visibility: payload.visibility ?? 'collaborators'
		}
	);
	return response.data.data;
}

export async function getDocumentation(projectId: UUID): Promise<ProjectDocumentation> {
	const response = await apiClient.get<ApiResponse<ProjectDocumentation>>(
		`/projects/${projectId}/documentation`
	);
	return response.data.data;
}

export async function getPublicDocumentation(slug: string): Promise<ProjectDocumentation> {
	const response = await apiClient.get<ApiResponse<ProjectDocumentation>>(
		`/projects/slug/${slug}/documentation`
	);
	return response.data.data;
}

export async function updateDocumentationVisibility(
	projectId: UUID,
	visibility: DocumentationVisibility
): Promise<ProjectDocumentation> {
	const response = await apiClient.patch<ApiResponse<ProjectDocumentation>>(
		`/projects/${projectId}/documentation/visibility`,
		{
			visibility
		}
	);
	return response.data.data;
}

export async function deleteDocumentation(projectId: UUID): Promise<void> {
	await apiClient.delete(`/projects/${projectId}/documentation`, { data: {} });
}

export async function createDocumentationSection(
	projectId: UUID,
	payload: {
		type: DocumentationSectionType;
		title: string;
		content: string;
		position?: number;
	}
): Promise<DocumentationSection> {
	const response = await apiClient.post<ApiResponse<DocumentationSection>>(
		`/projects/${projectId}/documentation/sections`,
		{
			type: payload.type,
			title: payload.title,
			content: payload.content,
			position: payload.position
		}
	);
	return response.data.data;
}

export async function listDocumentationSections(projectId: UUID): Promise<DocumentationSection[]> {
	const response = await apiClient.get<ApiResponse<DocumentationSection[]>>(
		`/projects/${projectId}/documentation/sections`
	);
	return response.data.data ?? [];
}

export async function updateDocumentationSection(
	sectionId: UUID,
	payload: {
		title?: string;
		content?: string;
		position?: number;
	}
): Promise<DocumentationSection> {
	const response = await apiClient.patch<ApiResponse<DocumentationSection>>(
		`/projects/documentation/sections/${sectionId}`,
		{
			title: payload.title,
			content: payload.content,
			position: payload.position
		}
	);
	return response.data.data;
}

export async function deleteDocumentationSection(sectionId: UUID): Promise<void> {
	await apiClient.delete(`/projects/documentation/sections/${sectionId}`, { data: {} });
}

export async function reorderDocumentationSections(
	projectId: UUID,
	sectionOrder: UUID[]
): Promise<DocumentationSection[]> {
	const response = await apiClient.patch<ApiResponse<DocumentationSection[]>>(
		`/projects/${projectId}/documentation/sections/reorder`,
		{
			section_order: sectionOrder
		}
	);
	return response.data.data ?? [];
}

// Technology API

export async function createTechnology(
	projectId: UUID,
	payload: {
		name: string;
		version: string;
		category: TechnologyCategory;
		purpose: string;
		link?: string;
	}
): Promise<ProjectTechnology> {
	const response = await apiClient.post<ApiResponse<ProjectTechnology>>(
		`/projects/${projectId}/technologies`,
		{
			name: payload.name,
			version: payload.version,
			category: payload.category,
			purpose: payload.purpose,
			link: payload.link
		}
	);
	return response.data.data;
}

export async function listTechnologies(projectId: UUID): Promise<ProjectTechnology[]> {
	const response = await apiClient.get<ApiResponse<ProjectTechnology[]>>(
		`/projects/${projectId}/technologies`
	);
	return response.data.data ?? [];
}

export async function updateTechnology(
	techId: UUID,
	payload: {
		name?: string;
		version?: string;
		category?: TechnologyCategory;
		purpose?: string;
		link?: string;
	}
): Promise<ProjectTechnology> {
	const response = await apiClient.patch<ApiResponse<ProjectTechnology>>(
		`/projects/technologies/${techId}`,
		{
			name: payload.name,
			version: payload.version,
			category: payload.category,
			purpose: payload.purpose,
			link: payload.link
		}
	);
	return response.data.data;
}

export async function deleteTechnology(techId: UUID): Promise<void> {
	await apiClient.delete(`/projects/technologies/${techId}`, { data: {} });
}

export async function bulkCreateTechnologies(
	projectId: UUID,
	technologies: Array<{
		name: string;
		version: string;
		category: TechnologyCategory;
		purpose: string;
		link?: string;
	}>
): Promise<ProjectTechnology[]> {
	const response = await apiClient.post<ApiResponse<ProjectTechnology[]>>(
		`/projects/${projectId}/technologies/bulk`,
		{
			technologies
		}
	);
	return response.data.data ?? [];
}

export async function bulkUpdateTechnologies(
	projectId: UUID,
	technologies: Array<{
		tech_id: UUID;
		name?: string;
		version?: string;
		category?: TechnologyCategory;
		purpose?: string;
		link?: string;
	}>
): Promise<ProjectTechnology[]> {
	const response = await apiClient.patch<ApiResponse<ProjectTechnology[]>>(
		`/projects/${projectId}/technologies/bulk`,
		{
			technologies
		}
	);
	return response.data.data ?? [];
}

// File Structure API

export async function createFileStructure(
	projectId: UUID,
	payload: {
		path: string;
		name: string;
		is_directory: boolean;
		parent_id?: UUID;
		language?: string;
		line_count?: number;
		purpose?: string;
		position?: number;
	}
): Promise<ProjectFileStructure> {
	const response = await apiClient.post<ApiResponse<ProjectFileStructure>>(
		`/projects/${projectId}/file-structures`,
		{
			path: payload.path,
			name: payload.name,
			is_directory: payload.is_directory,
			parent_id: payload.parent_id,
			language: payload.language,
			line_count: payload.line_count ?? 0,
			purpose: payload.purpose,
			position: payload.position
		}
	);
	return response.data.data;
}

export async function listFileStructures(projectId: UUID): Promise<ProjectFileStructure[]> {
	const response = await apiClient.get<ApiResponse<ProjectFileStructure[]>>(
		`/projects/${projectId}/file-structures`
	);
	return response.data.data ?? [];
}

export async function updateFileStructure(
	fileStructureId: UUID,
	payload: {
		purpose?: string;
		line_count?: number;
		language?: string;
		position?: number;
	}
): Promise<ProjectFileStructure> {
	const response = await apiClient.patch<ApiResponse<ProjectFileStructure>>(
		`/projects/file-structures/${fileStructureId}`,
		{
			purpose: payload.purpose,
			line_count: payload.line_count,
			language: payload.language,
			position: payload.position
		}
	);
	return response.data.data;
}

export async function deleteFileStructure(fileStructureId: UUID): Promise<void> {
	await apiClient.delete(`/projects/file-structures/${fileStructureId}`, { data: {} });
}

export async function bulkCreateFileStructures(
	projectId: UUID,
	structures: Array<{
		path: string;
		name: string;
		is_directory: boolean;
		parent_id?: UUID;
		language?: string;
		line_count?: number;
		purpose?: string;
		position?: number;
	}>
): Promise<ProjectFileStructure[]> {
	const response = await apiClient.post<ApiResponse<ProjectFileStructure[]>>(
		`/projects/${projectId}/file-structures/bulk`,
		{
			structures
		}
	);
	return response.data.data ?? [];
}

export async function bulkUpdateFileStructures(
	projectId: UUID,
	structures: Array<{
		file_structure_id: UUID;
		purpose?: string;
		line_count?: number;
		language?: string;
		position?: number;
	}>
): Promise<ProjectFileStructure[]> {
	const response = await apiClient.patch<ApiResponse<ProjectFileStructure[]>>(
		`/projects/${projectId}/file-structures/bulk`,
		{
			structures
		}
	);
	return response.data.data ?? [];
}

// Architecture Diagram API

export async function createArchitectureDiagram(
	projectId: UUID,
	payload: {
		type: ArchitectureDiagramType;
		title: string;
		description: string;
		content: string;
		format?: string;
		image_url?: string;
	}
): Promise<ProjectArchitectureDiagram> {
	const response = await apiClient.post<ApiResponse<ProjectArchitectureDiagram>>(
		`/projects/${projectId}/architecture-diagrams`,
		{
			type: payload.type,
			title: payload.title,
			description: payload.description,
			content: payload.content,
			format: payload.format ?? 'mermaid',
			image_url: payload.image_url
		}
	);
	return response.data.data;
}

export async function listArchitectureDiagrams(
	projectId: UUID
): Promise<ProjectArchitectureDiagram[]> {
	const response = await apiClient.get<ApiResponse<ProjectArchitectureDiagram[]>>(
		`/projects/${projectId}/architecture-diagrams`
	);
	return response.data.data ?? [];
}

export async function getArchitectureDiagram(
	diagramId: UUID
): Promise<ProjectArchitectureDiagram> {
	const response = await apiClient.get<ApiResponse<ProjectArchitectureDiagram>>(
		`/projects/architecture-diagrams/${diagramId}`
	);
	return response.data.data;
}

export async function updateArchitectureDiagram(
	diagramId: UUID,
	payload: {
		title?: string;
		description?: string;
		content?: string;
		image_url?: string;
	}
): Promise<ProjectArchitectureDiagram> {
	const response = await apiClient.patch<ApiResponse<ProjectArchitectureDiagram>>(
		`/projects/architecture-diagrams/${diagramId}`,
		{
			title: payload.title,
			description: payload.description,
			content: payload.content,
			image_url: payload.image_url
		}
	);
	return response.data.data;
}

export async function deleteArchitectureDiagram(diagramId: UUID): Promise<void> {
	await apiClient.delete(`/projects/architecture-diagrams/${diagramId}`, { data: {} });
}

export const projectsApi = {
	listProjects,
	createProject,
	updateProjectStatus,
	updateProjectMetrics,
	addMilestone,
	listMilestones,
	bulkUpdateMilestones,
	toggleMilestone,
	getKanbanBoard,
	createKanbanColumn,
	updateKanbanColumn,
	reorderKanbanColumns,
	deleteKanbanColumn,
	createKanbanCard,
	updateKanbanCard,
	moveKanbanCard,
	deleteKanbanCard,
	createDependency,
	listDependencies,
	deleteDependency,
	duplicateProject,
	// Documentation
	createDocumentation,
	getDocumentation,
	getPublicDocumentation,
	updateDocumentationVisibility,
	deleteDocumentation,
	createDocumentationSection,
	listDocumentationSections,
	updateDocumentationSection,
	deleteDocumentationSection,
	reorderDocumentationSections,
	// Technology
	createTechnology,
	listTechnologies,
	updateTechnology,
	deleteTechnology,
	bulkCreateTechnologies,
	bulkUpdateTechnologies,
	// File Structure
	createFileStructure,
	listFileStructures,
	updateFileStructure,
	deleteFileStructure,
	bulkCreateFileStructures,
	bulkUpdateFileStructures,
	// Architecture Diagrams
	createArchitectureDiagram,
	listArchitectureDiagrams,
	getArchitectureDiagram,
	updateArchitectureDiagram,
	deleteArchitectureDiagram
};
