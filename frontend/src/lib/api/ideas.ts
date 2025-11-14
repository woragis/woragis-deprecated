import { apiClient } from '@clients/apiClient';
import type {
	Idea,
	IdeaCollaborator,
	IdeaLink,
	IdeaVersion,
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

const mapIdea = (dto: IdeaDTO): Idea => ({
	id: dto.ID,
	user_id: dto.UserID,
	title: dto.Title,
	description: dto.Description,
	pos_x: dto.PosX ?? 0,
	pos_y: dto.PosY ?? 0,
	color: dto.Color && dto.Color.trim().length > 0 ? dto.Color : DEFAULT_COLOR,
	project_id: dto.ProjectID ?? null,
	version: dto.Version ?? 1,
	created_at: dto.CreatedAt,
	updated_at: dto.UpdatedAt
});

const mapLink = (dto: IdeaLinkDTO): IdeaLink => ({
	id: dto.ID,
	user_id: dto.UserID,
	source_idea_id: dto.SourceIdeaID,
	target_idea_id: dto.TargetIdeaID,
	relation: dto.Relation,
	weight: dto.Weight ?? 1,
	bidirectional: Boolean(dto.Bidirectional),
	created_at: dto.CreatedAt
});

const mapVersion = (dto: IdeaVersionDTO): IdeaVersion => ({
	id: dto.ID,
	idea_id: dto.IdeaID,
	user_id: dto.UserID,
	editor_id: dto.EditorID,
	version: dto.Version,
	title: dto.Title,
	description: dto.Description,
	pos_x: dto.PosX ?? 0,
	pos_y: dto.PosY ?? 0,
	color: dto.Color ?? DEFAULT_COLOR,
	change_type: dto.ChangeType,
	created_at: dto.CreatedAt
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

