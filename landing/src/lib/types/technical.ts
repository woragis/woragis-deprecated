export interface TechnicalCaseStudy {
	id: string;
	title: string;
	description: string;
	challenge: string;
	solution: string;
	technologies: string[];
	architecture?: string; // Architecture description or diagram URL
	metrics?: {
		label: string;
		value: string;
		improvement?: string;
	}[];
	tradeoffs?: {
		decision: string;
		pros: string[];
		cons: string[];
	}[];
	lessonsLearned?: string[];
}

export interface SystemDesign {
	id: string;
	title: string;
	description: string;
	components: {
		name: string;
		description: string;
		technology: string;
	}[];
	dataFlow?: string;
	scalability?: string;
	reliability?: string;
	diagram?: string; // URL to diagram
}

export interface ProblemSolution {
	id: string;
	problem: string;
	context: string;
	solution: string;
	technologies: string[];
	impact: string;
	metrics?: {
		before: string;
		after: string;
		improvement: string;
	};
}

