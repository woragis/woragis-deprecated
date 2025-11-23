export interface CaseStudy {
	id: string;
	userId: string;
	projectId: string;
	projectSlug: string; // Links to the project
	title: string;
	problem: string;
	context: string;
	solution: string;
	approach: string[];
	architecture?: {
		diagram?: string; // Mermaid diagram syntax
		diagramType?: 'mermaid' | 'plantuml';
		description?: string;
		components?: Array<{
			name: string;
			description: string;
			technologies: string[];
		}>;
	};
	metrics?: {
		before?: Array<{ label: string; value: string }>;
		after?: Array<{ label: string; value: string }>;
		impact?: string;
	};
	lessonsLearned?: string[];
	technologies: string[];
	featured: boolean;
	createdAt: string;
	updatedAt: string;
}

