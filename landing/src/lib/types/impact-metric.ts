export type MetricType =
	| 'projects_delivered'
	| 'users_impacted'
	| 'performance_improvement'
	| 'cost_savings'
	| 'time_saved';

export type MetricUnit =
	| 'count'
	| 'percentage'
	| 'currency'
	| 'hours'
	| 'days'
	| 'months'
	| 'years'
	| 'milliseconds'
	| 'seconds'
	| 'minutes';

export type EntityType = 'project' | 'problem_solution' | 'case_study' | 'system_design';

export interface ImpactMetric {
	id: string;
	userId: string;
	type: MetricType;
	value: number;
	unit: MetricUnit;
	description?: string;
	entityType?: EntityType;
	entityId?: string;
	periodStart?: string;
	periodEnd?: string;
	featured: boolean;
	displayOrder: number;
	createdAt: string;
	updatedAt: string;
}

export interface ListImpactMetricsParams {
	type?: MetricType;
	entityType?: EntityType;
	entityId?: string;
	featured?: boolean;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

