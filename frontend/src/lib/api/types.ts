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
