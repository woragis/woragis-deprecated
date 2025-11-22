import type { IconType } from 'svelte-icons-pack';
import { SiGo, SiDocker, SiKubernetes, SiPython } from 'svelte-icons-pack/si';
import { Settings } from 'lucide-svelte';

export interface Skill {
	name: string;
	description: string;
	icon?: IconType | any;
	iconName?: string;
	color: string;
	bgGradient: string;
	borderColor: string;
	hoverBorderColor: string;
	shadowColor: string;
}

export const skills: Skill[] = [
	{
		name: 'Golang',
		description:
			'My primary language for building high-performance backend services. I develop main server applications with Go, leveraging its concurrency model and efficiency for scalable distributed systems.',
		icon: SiGo,
		color: 'cyan',
		bgGradient: 'from-cyan-900/30 to-cyan-800/20',
		borderColor: 'border-cyan-700/30',
		hoverBorderColor: 'hover:border-cyan-500/50',
		shadowColor: 'hover:shadow-cyan-500/20'
	},
	{
		name: 'Python',
		description:
			'Building AI servers and intelligent backend services with Python. I leverage modern ML frameworks and libraries like LangChain, OpenAI, and other AI tools to create intelligent applications that integrate seamlessly with my Golang services.',
		icon: SiPython,
		color: 'yellow',
		bgGradient: 'from-yellow-900/30 to-yellow-800/20',
		borderColor: 'border-yellow-700/30',
		hoverBorderColor: 'hover:border-yellow-500/50',
		shadowColor: 'hover:shadow-yellow-500/20'
	},
	{
		name: 'Docker',
		description:
			'Expertise in containerization for consistent deployments. I containerize applications for development, testing, and production environments, ensuring reproducibility and portability across different platforms.',
		icon: SiDocker,
		color: 'blue',
		bgGradient: 'from-blue-900/30 to-blue-800/20',
		borderColor: 'border-blue-700/30',
		hoverBorderColor: 'hover:border-blue-500/50',
		shadowColor: 'hover:shadow-blue-500/20'
	},
	{
		name: 'Kubernetes',
		description:
			'Orchestrating containerized applications at scale. I design and manage K8s clusters for production-grade infrastructure and deployments, implementing auto-scaling, service discovery, and high availability patterns.',
		icon: SiKubernetes,
		color: 'indigo',
		bgGradient: 'from-indigo-900/30 to-indigo-800/20',
		borderColor: 'border-indigo-700/30',
		hoverBorderColor: 'hover:border-indigo-500/50',
		shadowColor: 'hover:shadow-indigo-500/20'
	},
	{
		name: 'DevOps',
		description:
			'Bridging development and operations. I focus on automation, CI/CD pipelines, infrastructure as code, and cloud-native practices to streamline deployments and ensure reliable, scalable systems.',
		icon: Settings,
		color: 'purple',
		bgGradient: 'from-purple-900/30 to-purple-800/20',
		borderColor: 'border-purple-700/30',
		hoverBorderColor: 'hover:border-purple-500/50',
		shadowColor: 'hover:shadow-purple-500/20'
	}
];

