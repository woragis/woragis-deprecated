export interface Component {
	name: string;
	description: string;
	technology: string;
}

export interface SystemDesign {
	id: string;
	userId: string;
	title: string;
	description: string;
	components?: Component[];
	dataFlow?: string;
	scalability?: string;
	reliability?: string;
	diagram?: string;
	featured: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface ListSystemDesignsParams {
	featured?: boolean;
	limit?: number;
	offset?: number;
	orderBy?: string;
	order?: 'asc' | 'desc';
}

