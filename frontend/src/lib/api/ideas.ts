import { apiClient } from '@clients/apiClient';
import type {
	Idea,
	IdeaCollaborator,
	IdeaDocument,
	IdeaLink,
	IdeaNode,
	IdeaNodeConnection,
	IdeaVersion,
	ConnectionDirection,
	UUID
} from './types';

interface ApiResponse<T> {
	success: boolean;
	data: T;
}

type IdeaDTO = {
	ID: UUID;
	UserID: UUID;
	Title: string;
	Description?: string;
	PosX: number;
	PosY: number;
	Color?: string;
	ProjectID?: UUID | null;
	Version: number;
	CreatedAt: string;
	UpdatedAt: string;
};

type IdeaLinkDTO = {
	ID: UUID;
	UserID: UUID;
	SourceIdeaID: UUID;
	TargetIdeaID: UUID;
	Relation: string;
	Weight: number;
	Bidirectional: boolean;
	CreatedAt: string;
};

type IdeaVersionDTO = {
	ID: UUID;
	IdeaID: UUID;
	UserID: UUID;
	EditorID: UUID;
	Version: number;
	Title: string;
	Description?: string;
	PosX: number;
	PosY: number;
	Color: string;
	ChangeType: string;
	CreatedAt: string;
};

type IdeaCollaboratorDTO = {
	ID: UUID;
	OwnerID: UUID;
	CollaboratorID: UUID;
	Role: string;
	CreatedAt: string;
	UpdatedAt: string;
};

const DEFAULT_COLOR = '#2563eb';

const pick = <T, K extends keyof T, F extends string>(
	dto: T & Record<string, any>,
	keys: (K | F)[]
) => {
	for (const key of keys) {
		const value = dto[key as keyof typeof dto];
		if (value !== undefined && value !== null) {
			return value as T[K];
		}
	}
	return undefined;
};

const sanitizeSlug = (value: string) =>
	value
		.toLowerCase()
		.trim()
		.replace(/[^a-z0-9]+/g, '-')
		.replace(/^-+|-+$/g, '') || 'idea';

const buildIdeaSlug = (title: string) => sanitizeSlug(title) || 'idea';

const mapIdea = (dto: IdeaDTO & Record<string, any>): Idea => {
	const id = pick(dto, ['ID', 'id', 'Id']);
	const userId = pick(dto, ['UserID', 'userId', 'UserId']);
	const posX = pick(dto, ['PosX', 'posX']);
	const posY = pick(dto, ['PosY', 'posY']);
	const createdAt = pick(dto, ['CreatedAt', 'createdAt']);
	const updatedAt = pick(dto, ['UpdatedAt', 'updatedAt']) ?? createdAt;
	const title = dto.Title ?? dto.title ?? '';
	const ideaId = id ?? crypto.randomUUID();

	return {
		id: ideaId,
		user_id: userId ?? null,
		title,
		description: dto.Description ?? dto.description,
		slug: dto.Slug ?? dto.slug ?? buildIdeaSlug(title),
		pos_x: Number(posX ?? 0),
		pos_y: Number(posY ?? 0),
		color: (dto.Color ?? dto.color ?? DEFAULT_COLOR) || DEFAULT_COLOR,
		project_id: dto.ProjectID ?? dto.projectId ?? null,
		version: dto.Version ?? dto.version ?? 1,
		created_at: (createdAt as string) ?? new Date().toISOString(),
		updated_at: (updatedAt as string) ?? new Date().toISOString()
	};
};

const mapLink = (dto: IdeaLinkDTO & Record<string, any>): IdeaLink => ({
	id: pick(dto, ['ID', 'id']) ?? crypto.randomUUID(),
	user_id: pick(dto, ['UserID', 'userId']) ?? null,
	source_idea_id: pick(dto, ['SourceIdeaID', 'sourceIdeaId']) ?? '',
	target_idea_id: pick(dto, ['TargetIdeaID', 'targetIdeaId']) ?? '',
	relation: dto.Relation ?? dto.relation ?? 'relates',
	weight: dto.Weight ?? dto.weight ?? 1,
	bidirectional: Boolean(dto.Bidirectional ?? dto.bidirectional),
	created_at: (pick(dto, ['CreatedAt', 'createdAt']) as string) ?? new Date().toISOString()
});

const mapVersion = (dto: IdeaVersionDTO & Record<string, any>): IdeaVersion => ({
	id: pick(dto, ['ID', 'id']) ?? crypto.randomUUID(),
	idea_id: pick(dto, ['IdeaID', 'ideaId']) ?? '',
	user_id: pick(dto, ['UserID', 'userId']) ?? '',
	editor_id: pick(dto, ['EditorID', 'editorId']) ?? '',
	version: dto.Version ?? dto.version ?? 1,
	title: dto.Title ?? dto.title ?? '',
	description: dto.Description ?? dto.description,
	pos_x: Number(pick(dto, ['PosX', 'posX']) ?? 0),
	pos_y: Number(pick(dto, ['PosY', 'posY']) ?? 0),
	color: dto.Color ?? dto.color ?? DEFAULT_COLOR,
	change_type: dto.ChangeType ?? dto.changeType ?? 'update',
	created_at: (pick(dto, ['CreatedAt', 'createdAt']) as string) ?? new Date().toISOString()
});

const mapCollaborator = (dto: IdeaCollaboratorDTO): IdeaCollaborator => ({
	id: dto.ID,
	owner_id: dto.OwnerID,
	collaborator_id: dto.CollaboratorID,
	role: dto.Role,
	created_at: dto.CreatedAt,
	updated_at: dto.UpdatedAt
});

export interface CreateIdeaInput {
	title: string;
	description?: string;
	pos_x?: number;
	pos_y?: number;
	color?: string;
	project_id?: UUID | null;
}

export interface UpdateIdeaInput {
	title?: string;
	description?: string;
	color?: string;
	project_id?: UUID | null;
}

export interface MoveIdeaInput {
	pos_x: number;
	pos_y: number;
}

export interface CreateLinkInput {
	source_idea_id: UUID;
	target_idea_id: UUID;
	relation: string;
	weight?: number;
	bidirectional?: boolean;
}

export interface ListIdeasParams {
	owner_id?: UUID;
}

export interface ListLinksParams {
	owner_id?: UUID;
	idea_id?: UUID;
	relation?: string;
	search?: string;
	min_weight?: number;
	max_weight?: number;
	bidirectional?: boolean;
}

export async function fetchIdeas(params: ListIdeasParams = {}): Promise<Idea[]> {
	const response = await apiClient.get<ApiResponse<IdeaDTO[]>>('/ideas', {
		params: {
			owner_id: params.owner_id
		}
	});
	return (response.data.data ?? []).map(mapIdea);
}

export async function fetchIdeaBySlug(slug: string): Promise<Idea> {
	const response = await apiClient.get<ApiResponse<IdeaDTO>>(`/ideas/slug/${slug}`);
	return mapIdea(response.data.data);
}

export async function createIdea(input: CreateIdeaInput): Promise<Idea> {
	const response = await apiClient.post<ApiResponse<IdeaDTO>>('/ideas', {
		title: input.title,
		description: input.description,
		pos_x: input.pos_x ?? 0,
		pos_y: input.pos_y ?? 0,
		color: input.color ?? DEFAULT_COLOR,
		project_id: input.project_id ?? null
	});
	return mapIdea(response.data.data);
}

export async function updateIdea(ideaId: UUID, input: UpdateIdeaInput): Promise<Idea> {
	const response = await apiClient.patch<ApiResponse<IdeaDTO>>(`/ideas/${ideaId}`, {
		title: input.title,
		description: input.description,
		color: input.color,
		project_id: input.project_id ?? null
	});
	return mapIdea(response.data.data);
}

export async function moveIdea(ideaId: UUID, input: MoveIdeaInput): Promise<void> {
	await apiClient.patch(`/ideas/${ideaId}/position`, {
		pos_x: input.pos_x,
		pos_y: input.pos_y
	});
}

export async function bulkMoveIdeas(moves: { idea_id: UUID; pos_x: number; pos_y: number }[]): Promise<void> {
	if (moves.length === 0) {
		return;
	}
	await apiClient.post('/ideas/bulk/move', {
		items: moves.map((move) => ({
			idea_id: move.idea_id,
			pos_x: move.pos_x,
			pos_y: move.pos_y
		}))
	});
}

export interface BulkUpdateItem extends UpdateIdeaInput {
	idea_id: UUID;
}

export async function bulkUpdateIdeas(updates: BulkUpdateItem[]): Promise<void> {
	if (updates.length === 0) {
		return;
	}
	await apiClient.post('/ideas/bulk/update', {
		items: updates.map((item) => ({
			idea_id: item.idea_id,
			title: item.title,
			description: item.description,
			color: item.color,
			project_id: item.project_id ?? null
		}))
	});
}

export async function deleteIdeas(ids: UUID[]): Promise<void> {
	if (ids.length === 0) {
		return;
	}
	await apiClient.post('/ideas/bulk/delete', {
		idea_ids: ids
	});
}

export async function restoreIdeas(ids: UUID[]): Promise<void> {
	if (ids.length === 0) {
		return;
	}
	await apiClient.post('/ideas/bulk/restore', {
		idea_ids: ids
	});
}

export async function fetchLinks(params: ListLinksParams = {}): Promise<IdeaLink[]> {
	const response = await apiClient.get<ApiResponse<IdeaLinkDTO[]>>('/ideas/links', {
		params: {
			owner_id: params.owner_id,
			idea_id: params.idea_id,
			relation: params.relation,
			search: params.search,
			min_weight: params.min_weight,
			max_weight: params.max_weight,
			bidirectional: params.bidirectional
		}
	});
	return (response.data.data ?? []).map(mapLink);
}

export async function createLink(input: CreateLinkInput): Promise<IdeaLink> {
	const response = await apiClient.post<ApiResponse<IdeaLinkDTO>>('/ideas/links', {
		source_idea_id: input.source_idea_id,
		target_idea_id: input.target_idea_id,
		relation: input.relation,
		weight: input.weight ?? 1,
		bidirectional: Boolean(input.bidirectional)
	});
	return mapLink(response.data.data);
}

export async function fetchVersions(ideaId: UUID, limit = 20): Promise<IdeaVersion[]> {
	const response = await apiClient.get<ApiResponse<IdeaVersionDTO[]>>(`/ideas/${ideaId}/versions`, {
		params: {
			limit
		}
	});
	return (response.data.data ?? []).map(mapVersion);
}

export async function addCollaborator(ownerId: UUID | undefined, collaboratorId: UUID, role?: string) {
	const response = await apiClient.post<ApiResponse<IdeaCollaboratorDTO>>('/ideas/collaborators', {
		owner_id: ownerId,
		collaborator_id: collaboratorId,
		role
	});
	return mapCollaborator(response.data.data);
}

export async function listCollaborators(ownerId?: UUID): Promise<IdeaCollaborator[]> {
	const response = await apiClient.get<ApiResponse<IdeaCollaboratorDTO[]>>('/ideas/collaborators', {
		params: {
			owner_id: ownerId
		}
	});
	return (response.data.data ?? []).map(mapCollaborator);
}

export async function removeCollaborator(ownerId: UUID | undefined, collaboratorId: UUID): Promise<void> {
	await apiClient.delete(`/ideas/collaborators/${collaboratorId}`, {
		params: {
			owner_id: ownerId
		}
	});
}

// IdeaNode API methods

type IdeaNodeDTO = {
	ID: UUID;
	IdeaID: UUID;
	Title: string;
	Description?: string;
	PosX: number;
	PosY: number;
	Width: number;
	Height: number;
	Color: string;
	Type: string;
	Version: number;
	CreatedAt: string;
	UpdatedAt: string;
};

type IdeaNodeConnectionDTO = {
	ID: UUID;
	IdeaID: UUID;
	SourceNodeID: UUID;
	TargetNodeID: UUID;
	Direction: ConnectionDirection;
	Label?: string;
	CreatedAt: string;
};

const mapIdeaNode = (dto: IdeaNodeDTO & Record<string, any>): IdeaNode => ({
	id: pick(dto, ['ID', 'id']) ?? crypto.randomUUID(),
	idea_id: pick(dto, ['IdeaID', 'ideaId', 'idea_id']) ?? '',
	title: dto.Title ?? dto.title ?? '',
	description: dto.Description ?? dto.description,
	pos_x: Number(pick(dto, ['PosX', 'posX', 'pos_x']) ?? 0),
	pos_y: Number(pick(dto, ['PosY', 'posY', 'pos_y']) ?? 0),
	width: Number(pick(dto, ['Width', 'width']) ?? 200),
	height: Number(pick(dto, ['Height', 'height']) ?? 100),
	color: dto.Color ?? dto.color ?? DEFAULT_COLOR,
	type: dto.Type ?? dto.type ?? 'default',
	version: dto.Version ?? dto.version ?? 1,
	created_at: (pick(dto, ['CreatedAt', 'createdAt']) as string) ?? new Date().toISOString(),
	updated_at: (pick(dto, ['UpdatedAt', 'updatedAt']) as string) ?? new Date().toISOString()
});

const mapIdeaNodeConnection = (dto: IdeaNodeConnectionDTO & Record<string, any>): IdeaNodeConnection => ({
	id: pick(dto, ['ID', 'id']) ?? crypto.randomUUID(),
	idea_id: pick(dto, ['IdeaID', 'ideaId', 'idea_id']) ?? '',
	source_node_id: pick(dto, ['SourceNodeID', 'sourceNodeId', 'source_node_id']) ?? '',
	target_node_id: pick(dto, ['TargetNodeID', 'targetNodeId', 'target_node_id']) ?? '',
	direction: (dto.Direction ?? dto.direction) as ConnectionDirection,
	label: dto.Label ?? dto.label,
	created_at: (pick(dto, ['CreatedAt', 'createdAt']) as string) ?? new Date().toISOString()
});

export interface CreateIdeaNodeInput {
	idea_id: UUID;
	title: string;
	description?: string;
	pos_x?: number;
	pos_y?: number;
	width?: number;
	height?: number;
	color?: string;
	type?: string;
}

export interface UpdateIdeaNodeInput {
	title?: string;
	description?: string;
	color?: string;
	type?: string;
}

export interface MoveIdeaNodeInput {
	pos_x: number;
	pos_y: number;
}

export interface ResizeIdeaNodeInput {
	width: number;
	height: number;
}

export interface CreateIdeaNodeConnectionInput {
	idea_id: UUID;
	source_node_id: UUID;
	target_node_id: UUID;
	direction: ConnectionDirection;
	label?: string;
}

export async function fetchIdeaNodes(ideaId: UUID): Promise<IdeaNode[]> {
	const response = await apiClient.get<ApiResponse<IdeaNodeDTO[]>>(`/ideas/${ideaId}/nodes`);
	return (response.data.data ?? []).map(mapIdeaNode);
}

export async function createIdeaNode(input: CreateIdeaNodeInput): Promise<IdeaNode> {
	const response = await apiClient.post<ApiResponse<IdeaNodeDTO>>('/ideas/nodes', {
		idea_id: input.idea_id,
		title: input.title,
		description: input.description,
		pos_x: input.pos_x ?? 0,
		pos_y: input.pos_y ?? 0,
		width: input.width ?? 200,
		height: input.height ?? 100,
		color: input.color ?? DEFAULT_COLOR,
		type: input.type ?? 'default'
	});
	return mapIdeaNode(response.data.data);
}

export async function updateIdeaNode(nodeId: UUID, input: UpdateIdeaNodeInput): Promise<IdeaNode> {
	const response = await apiClient.patch<ApiResponse<IdeaNodeDTO>>(`/ideas/nodes/${nodeId}`, {
		title: input.title,
		description: input.description,
		color: input.color,
		type: input.type
	});
	return mapIdeaNode(response.data.data);
}

export async function moveIdeaNode(nodeId: UUID, input: MoveIdeaNodeInput): Promise<IdeaNode> {
	const response = await apiClient.patch<ApiResponse<IdeaNodeDTO>>(`/ideas/nodes/${nodeId}/position`, {
		pos_x: input.pos_x,
		pos_y: input.pos_y
	});
	return mapIdeaNode(response.data.data);
}

export async function resizeIdeaNode(nodeId: UUID, input: ResizeIdeaNodeInput): Promise<IdeaNode> {
	const response = await apiClient.patch<ApiResponse<IdeaNodeDTO>>(`/ideas/nodes/${nodeId}/resize`, {
		width: input.width,
		height: input.height
	});
	return mapIdeaNode(response.data.data);
}

export async function deleteIdeaNode(nodeId: UUID): Promise<void> {
	await apiClient.delete(`/ideas/nodes/${nodeId}`);
}

export async function fetchIdeaNodeConnections(ideaId: UUID): Promise<IdeaNodeConnection[]> {
	const response = await apiClient.get<ApiResponse<IdeaNodeConnectionDTO[]>>(`/ideas/${ideaId}/node-connections`);
	return (response.data.data ?? []).map(mapIdeaNodeConnection);
}

export async function createIdeaNodeConnection(input: CreateIdeaNodeConnectionInput): Promise<IdeaNodeConnection> {
	const response = await apiClient.post<ApiResponse<IdeaNodeConnectionDTO>>('/ideas/node-connections', {
		idea_id: input.idea_id,
		source_node_id: input.source_node_id,
		target_node_id: input.target_node_id,
		direction: input.direction,
		label: input.label
	});
	return mapIdeaNodeConnection(response.data.data);
}

export async function deleteIdeaNodeConnection(connId: UUID): Promise<void> {
	await apiClient.delete(`/ideas/node-connections/${connId}`);
}

// Document API methods

type IdeaDocumentDTO = {
	ID: UUID;
	IdeaID: UUID;
	NodeID?: UUID | null;
	Title: string;
	Content: string;
	Version: number;
	CreatedAt: string;
	UpdatedAt: string;
};

const mapIdeaDocument = (dto: IdeaDocumentDTO & Record<string, any>): IdeaDocument => ({
	id: pick(dto, ['ID', 'id']) ?? crypto.randomUUID(),
	idea_id: pick(dto, ['IdeaID', 'ideaId', 'idea_id']) ?? '',
	node_id: pick(dto, ['NodeID', 'nodeId', 'node_id']) ?? null,
	title: dto.Title ?? dto.title ?? '',
	content: dto.Content ?? dto.content ?? '',
	version: dto.Version ?? dto.version ?? 1,
	created_at: (pick(dto, ['CreatedAt', 'createdAt']) as string) ?? new Date().toISOString(),
	updated_at: (pick(dto, ['UpdatedAt', 'updatedAt']) as string) ?? new Date().toISOString()
});

export interface CreateIdeaDocumentInput {
	idea_id: UUID;
	node_id?: UUID | null;
	title: string;
	content: string;
}

export interface UpdateIdeaDocumentInput {
	title?: string;
	content?: string;
}

export interface ListDocumentsParams {
	node_id?: UUID | null;
}

export async function fetchIdeaDocuments(ideaId: UUID, params?: ListDocumentsParams): Promise<IdeaDocument[]> {
	const response = await apiClient.get<ApiResponse<IdeaDocumentDTO[]>>(`/ideas/${ideaId}/documents`, {
		params: {
			node_id: params?.node_id
		}
	});
	return (response.data.data ?? []).map(mapIdeaDocument);
}

export async function createIdeaDocument(input: CreateIdeaDocumentInput): Promise<IdeaDocument> {
	const response = await apiClient.post<ApiResponse<IdeaDocumentDTO>>('/ideas/documents', {
		idea_id: input.idea_id,
		node_id: input.node_id ?? null,
		title: input.title,
		content: input.content
	});
	return mapIdeaDocument(response.data.data);
}

export async function updateIdeaDocument(docId: UUID, input: UpdateIdeaDocumentInput): Promise<IdeaDocument> {
	const response = await apiClient.patch<ApiResponse<IdeaDocumentDTO>>(`/ideas/documents/${docId}`, {
		title: input.title,
		content: input.content
	});
	return mapIdeaDocument(response.data.data);
}

export async function deleteIdeaDocument(docId: UUID): Promise<void> {
	await apiClient.delete(`/ideas/documents/${docId}`);
}

export const ideasApi = {
	fetchIdeaBySlug,
	fetchIdeas,
	createIdea,
	updateIdea,
	moveIdea,
	bulkMoveIdeas,
	bulkUpdateIdeas,
	deleteIdeas,
	restoreIdeas,
	fetchLinks,
	createLink,
	fetchVersions,
	addCollaborator,
	listCollaborators,
	removeCollaborator,
	fetchIdeaNodes,
	createIdeaNode,
	updateIdeaNode,
	moveIdeaNode,
	resizeIdeaNode,
	deleteIdeaNode,
	fetchIdeaNodeConnections,
	createIdeaNodeConnection,
	deleteIdeaNodeConnection,
	fetchIdeaDocuments,
	createIdeaDocument,
	updateIdeaDocument,
	deleteIdeaDocument
};

