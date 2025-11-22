import type { IconType } from 'svelte-icons-pack';
import { SiRedis } from 'svelte-icons-pack/si';
import { Brain, GitBranch } from 'lucide-svelte';

export interface Interest {
	title: string;
	description: string;
	icon?: IconType | any;
	iconName?: string;
	color: string;
	bgGradient: string;
	borderColor: string;
	hoverBorderColor: string;
	shadowColor: string;
	fullWidth?: boolean;
}

export const interests: Interest[] = [
	{
		title: 'AI & RAG',
		description:
			'Fascinated by Artificial Intelligence and Retrieval-Augmented Generation (RAG) systems. I\'m exploring how to build intelligent applications that can retrieve and synthesize information effectively. I work with AI servers built in Python, leveraging modern ML frameworks and libraries to create intelligent backend services.',
		icon: Brain,
		color: 'pink-purple',
		bgGradient: 'from-pink-900/30 to-purple-900/20',
		borderColor: 'border-pink-700/30',
		hoverBorderColor: 'hover:border-pink-500/50',
		shadowColor: 'hover:shadow-pink-500/20'
	},
	{
		title: 'Redis & Pub/Sub',
		description:
			'Deep interest in Redis Pub/Sub patterns for real-time communication between distributed services. I design systems where multiple servers communicate seamlessly through Redis messaging. I implement inter-service communication architectures using Redis as the backbone, enabling scalable and responsive distributed applications.',
		icon: SiRedis,
		color: 'red-orange',
		bgGradient: 'from-red-900/30 to-orange-900/20',
		borderColor: 'border-red-700/30',
		hoverBorderColor: 'hover:border-red-500/50',
		shadowColor: 'hover:shadow-red-500/20'
	},
	{
		title: 'Distributed Architecture',
		description:
			'I specialize in building hybrid architectures where the main server is built with Golang for performance and reliability, while AI services are implemented in Python to leverage the rich ML ecosystem. These services communicate through Redis Pub/Sub, enabling a microservices architecture that\'s both scalable and maintainable. The combination allows each service to use the best tools for its specific domain.',
		icon: GitBranch,
		color: 'green-emerald',
		bgGradient: 'from-green-900/30 to-emerald-900/20',
		borderColor: 'border-green-700/30',
		hoverBorderColor: 'hover:border-green-500/50',
		shadowColor: 'hover:shadow-green-500/20',
		fullWidth: true
	}
];

