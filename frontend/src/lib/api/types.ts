export type UUID = string;

export type TransactionType = 'income' | 'expense';

export interface Transaction {
	id: UUID;
	user_id: UUID;
	type: TransactionType;
	category: string;
	description?: string;
	amount: number;
	currency: string;
	base_currency: string;
	normalized_amount: number;
	occurred_at: string;
	is_recurring: boolean;
	is_essential: boolean;
	is_archived: boolean;
	template_id?: UUID | null;
	tags: string[];
	created_at: string;
}

export interface TransactionSummary {
	income_total: number;
	expense_total: number;
	savings_allocation: number;
	base_currency: string;
}

export interface KanbanCard {
	id: UUID;
	project_id: UUID;
	column_id: UUID;
	milestone_id?: UUID | null;
	title: string;
	description?: string;
	due_date?: string | null;
	position: number;
	created_at: string;
	updated_at: string;
}

export interface KanbanColumnMeta {
	id: UUID;
	project_id: UUID;
	name: string;
	wip_limit: number;
	position: number;
	created_at: string;
	updated_at: string;
}

export interface KanbanColumnWithCards {
	column: KanbanColumnMeta;
	cards: KanbanCard[];
}

export interface KanbanBoard {
	project_id: UUID;
	columns: KanbanColumnWithCards[];
}

export type ProjectStatus = 'idea' | 'planning' | 'executing' | 'monitoring' | 'completed';

export interface Project {
	id: UUID;
	user_id: UUID;
	name: string;
	description: string;
	slug: string;
	status: ProjectStatus;
	health_score: number;
	mrr: number;
	cac: number;
	ltv: number;
	churn_rate: number;
	created_at: string;
	updated_at: string;
}

export interface Milestone {
	id: UUID;
	project_id: UUID;
	title: string;
	description: string;
	due_date: string;
	completed: boolean;
	created_at: string;
	updated_at: string;
}

export interface ProjectDependency {
	id: UUID;
	project_id: UUID;
	depends_on_project_id: UUID;
	type: 'blocks' | 'relates' | 'supports';
	created_at: string;
	updated_at: string;
}

export interface ReportDefinition {
	id: UUID;
	name: string;
	description: string;
	sections: Record<string, unknown>;
	filters: Record<string, unknown>;
	is_favorite: boolean;
	archived_at?: string | null;
	created_at: string;
	updated_at: string;
}

export interface ReportSchedule {
	id: UUID;
	report_id: UUID;
	cron: string;
	frequency: string;
	timezone: string;
	next_run?: string | null;
	last_run_at?: string | null;
	enabled: boolean;
	meta?: Record<string, unknown>;
	created_at: string;
	updated_at: string;
}

export interface ReportDelivery {
	id: UUID;
	report_id: UUID;
	channel: string;
	target: string;
	template?: Record<string, unknown>;
	enabled: boolean;
	created_at: string;
	updated_at: string;
}

export interface ReportRun {
	id: UUID;
	report_id: UUID;
	status: 'pending' | 'running' | 'completed' | 'failed';
	started_at?: string | null;
	completed_at?: string | null;
	output_location?: string;
	error_message?: string;
	metadata?: Record<string, unknown>;
	created_at: string;
	updated_at: string;
}

export interface ReportDefinitionDetail {
	definition: ReportDefinition;
	schedules: ReportSchedule[];
	deliveries: ReportDelivery[];
}

export interface Idea {
	id: UUID;
	user_id: UUID;
	title: string;
	description?: string;
	slug: string;
	pos_x: number;
	pos_y: number;
	color: string;
	project_id?: UUID | null;
	version: number;
	created_at: string;
	updated_at: string;
}

export interface IdeaLink {
	id: UUID;
	user_id: UUID;
	source_idea_id: UUID;
	target_idea_id: UUID;
	relation: string;
	weight: number;
	bidirectional: boolean;
	created_at: string;
}

export interface IdeaVersion {
	id: UUID;
	idea_id: UUID;
	user_id: UUID;
	editor_id: UUID;
	version: number;
	title: string;
	description?: string;
	pos_x: number;
	pos_y: number;
	color: string;
	change_type: string;
	created_at: string;
}

export interface IdeaCollaborator {
	id: UUID;
	owner_id: UUID;
	collaborator_id: UUID;
	role: string;
	created_at: string;
	updated_at: string;
}

export interface ChatConversation {
	id: UUID;
	user_id: UUID;
	title: string;
	description?: string;
	idea_id?: UUID | null;
	project_id?: UUID | null;
	assigned_agent_id?: UUID | null;
	shared_transcript?: string;
	archived_at?: string | null;
	deleted_at?: string | null;
	created_at: string;
	updated_at: string;
}

export interface ChatMessage {
	id: UUID;
	conversation_id: UUID;
	role: string;
	content: string;
	created_at: string;
}

export interface ChatTranscript {
	id: UUID;
	share_code: string;
	content?: string;
	created_at: string;
	expires_at?: string | null;
}

export interface ChatAssignment {
	id: UUID;
	agent_id: UUID;
	agent_name?: string;
	assigned_at: string;
	unassigned_at?: string | null;
	notes?: string;
}
