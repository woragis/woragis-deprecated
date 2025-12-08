import { apiClient } from '@clients/apiClient';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

// Import types from types.ts to maintain consistency
import type { Project, ProjectStatus } from './types';

// Re-export for convenience
export type { Project, ProjectStatus };

export interface CreateProjectInput {
	name: string;
	description?: string;
	status?: ProjectStatus;
	healthScore?: number;
	mrr?: number;
	cac?: number;
	ltv?: number;
	churnRate?: number;
}

export interface UpdateProjectStatusInput {
	status: ProjectStatus;
}

export interface UpdateProjectMetricsInput {
	healthScore?: number;
	mrr?: number;
	cac?: number;
	ltv?: number;
	churnRate?: number;
}

// Helper function to map API response to Project type
const mapProject = (apiProject: any): Project => ({
	id: apiProject.id,
	user_id: apiProject.userId || apiProject.user_id,
	name: apiProject.name,
	description: apiProject.description || '',
	slug: apiProject.slug,
	status: apiProject.status,
	health_score: apiProject.healthScore ?? apiProject.health_score ?? 0,
	mrr: apiProject.mrr ?? 0,
	cac: apiProject.cac ?? 0,
	ltv: apiProject.ltv ?? 0,
	churn_rate: apiProject.churnRate ?? apiProject.churn_rate ?? 0,
	created_at: apiProject.createdAt || apiProject.created_at,
	updated_at: apiProject.updatedAt || apiProject.updated_at
});

export async function listProjects(): Promise<Project[]> {
	const response = await apiClient.get<ApiResponse<any[]>>('/projects');
	return (response.data.data ?? []).map(mapProject);
}

export async function getProject(id: string): Promise<Project> {
	const response = await apiClient.get<ApiResponse<any>>(`/projects/${id}`);
	return mapProject(response.data.data);
}

export async function getProjectBySlug(slug: string): Promise<Project> {
	const response = await apiClient.get<ApiResponse<any>>(`/projects/slug/${slug}`);
	return mapProject(response.data.data);
}

export async function createProject(input: CreateProjectInput): Promise<Project> {
	const response = await apiClient.post<ApiResponse<any>>('/projects', input);
	return mapProject(response.data.data);
}

export async function updateProjectStatus(id: string, status: ProjectStatus): Promise<Project> {
	const response = await apiClient.patch<ApiResponse<any>>(`/projects/${id}/status`, { status });
	return mapProject(response.data.data);
}

export async function updateProjectMetrics(id: string, input: UpdateProjectMetricsInput): Promise<Project> {
	const response = await apiClient.patch<ApiResponse<any>>(`/projects/${id}/metrics`, input);
	return mapProject(response.data.data);
}

export async function deleteProject(id: string): Promise<void> {
	await apiClient.delete(`/projects/${id}`);
}

// Import types for API object
import type {
	KanbanBoard,
	Milestone,
	ProjectDependency,
	ProjectDocumentation,
	DocumentationSection,
	ProjectTechnology,
	ProjectFileStructure,
	ProjectArchitectureDiagram,
	DocumentationVisibility,
	TechnologyCategory,
	ArchitectureDiagramType,
	UUID
} from './types';

// Kanban Board
export async function getKanbanBoard(projectId: string): Promise<KanbanBoard> {
	const response = await apiClient.get<ApiResponse<KanbanBoard>>(`/projects/${projectId}/board`);
	return response.data.data;
}

export interface CreateKanbanColumnInput {
	name: string;
	wipLimit?: number;
	position?: number;
}

export async function createKanbanColumn(
	projectId: string,
	input: CreateKanbanColumnInput
): Promise<any> {
	const response = await apiClient.post<ApiResponse<any>>(`/projects/${projectId}/columns`, input);
	return response.data.data;
}

export interface UpdateKanbanColumnInput {
	name?: string;
	wipLimit?: number;
	position?: number;
}

export async function updateKanbanColumn(
	projectId: string,
	columnId: string,
	input: UpdateKanbanColumnInput
): Promise<any> {
	const response = await apiClient.patch<ApiResponse<any>>(
		`/projects/${projectId}/columns/${columnId}`,
		input
	);
	return response.data.data;
}

export async function reorderKanbanColumns(projectId: string, columnOrder: string[]): Promise<void> {
	await apiClient.post(`/projects/${projectId}/columns/reorder`, { columnOrder });
}

export async function deleteKanbanColumn(projectId: string, columnId: string): Promise<void> {
	await apiClient.delete(`/projects/${projectId}/columns/${columnId}`);
}

export interface CreateKanbanCardInput {
	columnId: string;
	milestoneId?: string;
	title: string;
	description?: string;
	dueDate?: string;
	position?: number;
}

export async function createKanbanCard(
	projectId: string,
	input: CreateKanbanCardInput
): Promise<any> {
	const response = await apiClient.post<ApiResponse<any>>(`/projects/${projectId}/cards`, input);
	return response.data.data;
}

export interface UpdateKanbanCardInput {
	columnId?: string;
	milestoneId?: string;
	title?: string;
	description?: string;
	dueDate?: string;
	position?: number;
}

export async function updateKanbanCard(
	projectId: string,
	cardId: string,
	input: UpdateKanbanCardInput
): Promise<any> {
	const response = await apiClient.patch<ApiResponse<any>>(
		`/projects/${projectId}/cards/${cardId}`,
		input
	);
	return response.data.data;
}

export async function moveKanbanCard(
	projectId: string,
	cardId: string,
	targetColumnId: string,
	targetPosition: number
): Promise<void> {
	await apiClient.post(`/projects/${projectId}/cards/${cardId}/move`, {
		targetColumnId,
		targetPosition
	});
}

export async function deleteKanbanCard(projectId: string, cardId: string): Promise<void> {
	await apiClient.delete(`/projects/${projectId}/cards/${cardId}`);
}

// Milestones
export async function listMilestones(projectId: string): Promise<Milestone[]> {
	const response = await apiClient.get<ApiResponse<Milestone[]>>(`/projects/${projectId}/milestones`);
	return response.data.data ?? [];
}

export interface AddMilestoneInput {
	title: string;
	description: string;
	dueDate: string;
}

export async function addMilestone(projectId: string, input: AddMilestoneInput): Promise<Milestone> {
	const response = await apiClient.post<ApiResponse<Milestone>>(
		`/projects/${projectId}/milestones`,
		input
	);
	return response.data.data;
}

export interface BulkUpdateMilestonesInput {
	updates: Array<{
		id: string;
		completed?: boolean;
		dueDate?: string;
		title?: string;
		description?: string;
	}>;
}

export async function bulkUpdateMilestones(
	projectId: string,
	input: BulkUpdateMilestonesInput
): Promise<void> {
	await apiClient.patch(`/projects/${projectId}/milestones/bulk`, input);
}

// Dependencies
export async function listDependencies(projectId: string): Promise<ProjectDependency[]> {
	const response = await apiClient.get<ApiResponse<ProjectDependency[]>>(
		`/projects/${projectId}/dependencies`
	);
	return response.data.data ?? [];
}

export interface CreateDependencyInput {
	dependsOnProjectId: string;
	type: 'blocks' | 'relates' | 'supports';
}

export async function createDependency(
	projectId: string,
	input: CreateDependencyInput
): Promise<ProjectDependency> {
	const response = await apiClient.post<ApiResponse<ProjectDependency>>(
		`/projects/${projectId}/dependencies`,
		input
	);
	return response.data.data;
}

export async function deleteDependency(projectId: string, dependencyId: string): Promise<void> {
	await apiClient.delete(`/projects/${projectId}/dependencies/${dependencyId}`);
}

// Duplicate Project
export interface DuplicateProjectInput {
	name: string;
	description?: string;
	status?: ProjectStatus;
	copyMilestones?: boolean;
	copyKanban?: boolean;
	copyDependencies?: boolean;
}

export async function duplicateProject(
	templateProjectId: string,
	input: DuplicateProjectInput
): Promise<Project> {
	const response = await apiClient.post<ApiResponse<any>>(
		`/projects/${templateProjectId}/duplicate`,
		input
	);
	return mapProject(response.data.data);
}

// Documentation
export async function getDocumentation(projectId: string): Promise<ProjectDocumentation> {
	const response = await apiClient.get<ApiResponse<ProjectDocumentation>>(
		`/projects/${projectId}/documentation`
	);
	return response.data.data;
}

export interface CreateDocumentationInput {
	visibility: DocumentationVisibility;
}

export async function createDocumentation(
	projectId: string,
	input: CreateDocumentationInput
): Promise<ProjectDocumentation> {
	const response = await apiClient.post<ApiResponse<ProjectDocumentation>>(
		`/projects/${projectId}/documentation`,
		input
	);
	return response.data.data;
}

export async function updateDocumentationVisibility(
	projectId: string,
	visibility: DocumentationVisibility
): Promise<ProjectDocumentation> {
	const response = await apiClient.patch<ApiResponse<ProjectDocumentation>>(
		`/projects/${projectId}/documentation/visibility`,
		{ visibility }
	);
	return response.data.data;
}

export async function listDocumentationSections(projectId: string): Promise<DocumentationSection[]> {
	const response = await apiClient.get<ApiResponse<DocumentationSection[]>>(
		`/projects/${projectId}/documentation/sections`
	);
	return response.data.data ?? [];
}

export interface CreateDocumentationSectionInput {
	type: string;
	title: string;
	content: string;
	position?: number;
}

export async function createDocumentationSection(
	projectId: string,
	input: CreateDocumentationSectionInput
): Promise<DocumentationSection> {
	const response = await apiClient.post<ApiResponse<DocumentationSection>>(
		`/projects/${projectId}/documentation/sections`,
		input
	);
	return response.data.data;
}

export interface UpdateDocumentationSectionInput {
	type?: string;
	title?: string;
	content?: string;
	position?: number;
}

export async function updateDocumentationSection(
	sectionId: string,
	input: UpdateDocumentationSectionInput
): Promise<DocumentationSection> {
	const response = await apiClient.patch<ApiResponse<DocumentationSection>>(
		`/documentation/sections/${sectionId}`,
		input
	);
	return response.data.data;
}

export async function deleteDocumentationSection(sectionId: string): Promise<void> {
	await apiClient.delete(`/documentation/sections/${sectionId}`);
}

export async function reorderDocumentationSections(
	projectId: string,
	sectionOrder: string[]
): Promise<void> {
	await apiClient.post(`/projects/${projectId}/documentation/sections/reorder`, { sectionOrder });
}

// Technologies
export async function listTechnologies(projectId: string): Promise<ProjectTechnology[]> {
	const response = await apiClient.get<ApiResponse<ProjectTechnology[]>>(
		`/projects/${projectId}/technologies`
	);
	return response.data.data ?? [];
}

export interface CreateTechnologyInput {
	name: string;
	version: string;
	category: TechnologyCategory;
	purpose: string;
	link?: string;
}

export async function createTechnology(
	projectId: string,
	input: CreateTechnologyInput
): Promise<ProjectTechnology> {
	const response = await apiClient.post<ApiResponse<ProjectTechnology>>(
		`/projects/${projectId}/technologies`,
		input
	);
	return response.data.data;
}

export interface UpdateTechnologyInput {
	name?: string;
	version?: string;
	category?: TechnologyCategory;
	purpose?: string;
	link?: string;
}

export async function updateTechnology(
	techId: string,
	input: UpdateTechnologyInput
): Promise<ProjectTechnology> {
	const response = await apiClient.patch<ApiResponse<ProjectTechnology>>(
		`/technologies/${techId}`,
		input
	);
	return response.data.data;
}

export async function deleteTechnology(techId: string): Promise<void> {
	await apiClient.delete(`/technologies/${techId}`);
}

export async function bulkCreateTechnologies(
	projectId: string,
	technologies: CreateTechnologyInput[]
): Promise<ProjectTechnology[]> {
	const response = await apiClient.post<ApiResponse<ProjectTechnology[]>>(
		`/projects/${projectId}/technologies/bulk`,
		{ technologies }
	);
	return response.data.data ?? [];
}

export async function bulkUpdateTechnologies(
	projectId: string,
	technologies: Array<{ id: string } & UpdateTechnologyInput>
): Promise<void> {
	await apiClient.patch(`/projects/${projectId}/technologies/bulk`, { technologies });
}

// File Structures
export async function listFileStructures(projectId: string): Promise<ProjectFileStructure[]> {
	const response = await apiClient.get<ApiResponse<ProjectFileStructure[]>>(
		`/projects/${projectId}/file-structures`
	);
	return response.data.data ?? [];
}

export interface CreateFileStructureInput {
	parentId?: string;
	path: string;
	name: string;
	isDirectory: boolean;
	language?: string;
	lineCount?: number;
	purpose?: string;
	position?: number;
}

export async function createFileStructure(
	projectId: string,
	input: CreateFileStructureInput
): Promise<ProjectFileStructure> {
	const response = await apiClient.post<ApiResponse<ProjectFileStructure>>(
		`/projects/${projectId}/file-structures`,
		input
	);
	return response.data.data;
}

export interface UpdateFileStructureInput {
	parentId?: string;
	path?: string;
	name?: string;
	isDirectory?: boolean;
	language?: string;
	lineCount?: number;
	purpose?: string;
	position?: number;
}

export async function updateFileStructure(
	fileStructureId: string,
	input: UpdateFileStructureInput
): Promise<ProjectFileStructure> {
	const response = await apiClient.patch<ApiResponse<ProjectFileStructure>>(
		`/file-structures/${fileStructureId}`,
		input
	);
	return response.data.data;
}

export async function deleteFileStructure(fileStructureId: string): Promise<void> {
	await apiClient.delete(`/file-structures/${fileStructureId}`);
}

export async function bulkCreateFileStructures(
	projectId: string,
	structures: CreateFileStructureInput[]
): Promise<ProjectFileStructure[]> {
	const response = await apiClient.post<ApiResponse<ProjectFileStructure[]>>(
		`/projects/${projectId}/file-structures/bulk`,
		{ structures }
	);
	return response.data.data ?? [];
}

export async function bulkUpdateFileStructures(
	projectId: string,
	structures: Array<{ id: string } & UpdateFileStructureInput>
): Promise<void> {
	await apiClient.patch(`/projects/${projectId}/file-structures/bulk`, { structures });
}

// Architecture Diagrams
export async function listArchitectureDiagrams(
	projectId: string
): Promise<ProjectArchitectureDiagram[]> {
	const response = await apiClient.get<ApiResponse<ProjectArchitectureDiagram[]>>(
		`/projects/${projectId}/architecture-diagrams`
	);
	return response.data.data ?? [];
}

export interface CreateArchitectureDiagramInput {
	type: ArchitectureDiagramType;
	title: string;
	description: string;
	content: string;
	format: string;
	imageUrl?: string;
}

export async function createArchitectureDiagram(
	projectId: string,
	input: CreateArchitectureDiagramInput
): Promise<ProjectArchitectureDiagram> {
	const response = await apiClient.post<ApiResponse<ProjectArchitectureDiagram>>(
		`/projects/${projectId}/architecture-diagrams`,
		input
	);
	return response.data.data;
}

export interface UpdateArchitectureDiagramInput {
	type?: ArchitectureDiagramType;
	title?: string;
	description?: string;
	content?: string;
	format?: string;
	imageUrl?: string;
}

export async function updateArchitectureDiagram(
	diagramId: string,
	input: UpdateArchitectureDiagramInput
): Promise<ProjectArchitectureDiagram> {
	const response = await apiClient.patch<ApiResponse<ProjectArchitectureDiagram>>(
		`/architecture-diagrams/${diagramId}`,
		input
	);
	return response.data.data;
}

export async function deleteArchitectureDiagram(diagramId: string): Promise<void> {
	await apiClient.delete(`/architecture-diagrams/${diagramId}`);
}

// API Object
export const projectsApi = {
	listProjects,
	getProject,
	getProjectBySlug,
	createProject,
	updateProjectStatus,
	updateProjectMetrics,
	deleteProject,
	getKanbanBoard,
	createKanbanColumn,
	updateKanbanColumn,
	reorderKanbanColumns,
	deleteKanbanColumn,
	createKanbanCard,
	updateKanbanCard,
	moveKanbanCard,
	deleteKanbanCard,
	listMilestones,
	addMilestone,
	bulkUpdateMilestones,
	listDependencies,
	createDependency,
	deleteDependency,
	duplicateProject,
	getDocumentation,
	createDocumentation,
	updateDocumentationVisibility,
	listDocumentationSections,
	createDocumentationSection,
	updateDocumentationSection,
	deleteDocumentationSection,
	reorderDocumentationSections,
	listTechnologies,
	createTechnology,
	updateTechnology,
	deleteTechnology,
	bulkCreateTechnologies,
	bulkUpdateTechnologies,
	listFileStructures,
	createFileStructure,
	updateFileStructure,
	deleteFileStructure,
	bulkCreateFileStructures,
	bulkUpdateFileStructures,
	listArchitectureDiagrams,
	createArchitectureDiagram,
	updateArchitectureDiagram,
	deleteArchitectureDiagram
};
