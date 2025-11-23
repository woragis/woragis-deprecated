import type { TechnicalCaseStudy, SystemDesign, ProblemSolution } from '$lib/types/technical';

export const caseStudies: TechnicalCaseStudy[] = [
	{
		id: 'woragis-backend',
		title: 'Woragis Backend Architecture',
		description:
			'Built a scalable backend system using Go, Fiber, and PostgreSQL with microservices architecture supporting multiple domains including projects, finances, skills, and AI integrations.',
		challenge:
			'Need to build a robust, scalable backend that can handle multiple business domains while maintaining clean architecture, API key authentication, and real-time capabilities.',
		solution:
			'Implemented a domain-driven design with separate modules for each business domain, JWT and API key authentication, Redis for caching and pub/sub, and a clean separation of concerns with repositories, services, and handlers.',
		technologies: ['Go', 'Fiber', 'PostgreSQL', 'Redis', 'Docker', 'Kubernetes'],
		architecture:
			'Microservices-ready architecture with domain-driven design. Each domain (projects, skills, finances, etc.) has its own repository, service, and handler layers. API key authentication for public read access, JWT for authenticated operations.',
		metrics: [
			{ label: 'Response Time', value: '< 50ms', improvement: 'P95 latency' },
			{ label: 'Uptime', value: '99.9%', improvement: 'Target SLA' },
			{ label: 'API Endpoints', value: '50+', improvement: 'RESTful APIs' }
		],
		tradeoffs: [
			{
				decision: 'API Key vs JWT for public endpoints',
				pros: [
					'Simpler for public consumption',
					'No token expiration management',
					'Better for static site integration'
				],
				cons: ['Less granular permissions', 'Requires separate key management']
			},
			{
				decision: 'Monolithic vs Microservices',
				pros: [
					'Easier development and deployment',
					'Simpler debugging',
					'Lower operational complexity'
				],
				cons: [
					'Potential scaling bottlenecks',
					'Shared database can become a bottleneck'
				]
			}
		],
		lessonsLearned: [
			'Domain-driven design makes code more maintainable and testable',
			'API key authentication is essential for public-facing APIs',
			'Redis pub/sub enables real-time features without WebSocket complexity',
			'Clean architecture pays off in long-term maintainability'
		]
	}
];

export const systemDesigns: SystemDesign[] = [
	{
		id: 'api-architecture',
		title: 'RESTful API Architecture',
		description:
			'Clean, scalable API architecture with middleware-based authentication, CORS support, and domain-driven design.',
		components: [
			{
				name: 'API Gateway',
				description: 'Fiber router with CORS, logging, and recovery middleware',
				technology: 'Go Fiber'
			},
			{
				name: 'Authentication Layer',
				description: 'JWT for authenticated users, API keys for public read access',
				technology: 'JWT, API Keys'
			},
			{
				name: 'Domain Services',
				description: 'Business logic separated by domain (projects, skills, finances)',
				technology: 'Go'
			},
			{
				name: 'Data Layer',
				description: 'Repository pattern with GORM for database access',
				technology: 'GORM, PostgreSQL'
			},
			{
				name: 'Cache Layer',
				description: 'Redis for caching and pub/sub messaging',
				technology: 'Redis'
			}
		],
		dataFlow:
			'Request → CORS Middleware → Auth Middleware → Domain Handler → Service → Repository → Database. Responses cached in Redis when appropriate.',
		scalability:
			'Horizontal scaling ready with stateless services. Database connection pooling and Redis clustering support. Kubernetes-ready deployment.',
		reliability:
			'Error recovery middleware, structured logging, health checks, and graceful shutdown. Database migrations and rollback support.'
	},
	{
		id: 'authentication-flow',
		title: 'Dual Authentication System',
		description:
			'Flexible authentication supporting both JWT tokens for authenticated users and API keys for public read-only access.',
		components: [
			{
				name: 'JWT Authentication',
				description: 'Token-based auth for authenticated operations (POST, PATCH, DELETE)',
				technology: 'JWT, Go'
			},
			{
				name: 'API Key Authentication',
				description: 'Key-based auth for public read operations (GET)',
				technology: 'SHA256 Hashing'
			},
			{
				name: 'Middleware Chain',
				description: 'RequireAPIKeyOrAuth middleware checks API key first, falls back to JWT',
				technology: 'Go Fiber Middleware'
			},
			{
				name: 'Context Storage',
				description: 'User ID and API key stored in Fiber context for downstream handlers',
				technology: 'Fiber Context'
			}
		],
		dataFlow:
			'Request → Extract API Key/JWT → Validate → Store in Context → Handler → Service. API keys validated via SHA256 hash comparison.',
		scalability:
			'Stateless authentication allows horizontal scaling. API keys can be rate-limited per key. JWT tokens validated without database lookup.',
		reliability:
			'Graceful fallback from API key to JWT. Comprehensive error handling and logging. Secure key storage with hashing.'
	}
];

export const problemSolutions: ProblemSolution[] = [
	{
		id: 'cors-preflight',
		problem: 'CORS preflight requests failing for custom headers',
		context:
			'Frontend landing page couldn\'t make API requests with X-API-Key header due to CORS preflight failures. Browser normalizes custom headers to lowercase for preflight checks.',
		solution:
			'Updated CORS configuration to accept both X-API-Key and x-api-key in allowed headers. Browsers send lowercase headers in preflight, but the actual request can use either case.',
		technologies: ['CORS', 'Go Fiber', 'HTTP Headers'],
		impact:
			'Enabled seamless API integration from static landing page without CORS errors. Public API access now works reliably.',
		metrics: {
			before: 'CORS errors blocking all requests',
			after: '100% successful preflight requests',
			improvement: 'Zero CORS-related failures'
		}
	},
	{
		id: 'api-key-middleware',
		problem: 'API key authentication not working for GET requests',
		context:
			'Needed to allow public read access to projects and skills endpoints using API keys, while maintaining JWT authentication for write operations.',
		solution:
			'Created RequireAPIKeyOrAuth middleware that checks for API key first on GET requests, then falls back to JWT. Ensured middleware order so API key routes are registered before protected routes.',
		technologies: ['Go', 'Fiber Middleware', 'Authentication'],
		impact:
			'Landing page can now fetch projects and skills data without requiring user authentication. Clean separation between public read and authenticated write operations.',
		metrics: {
			before: 'All endpoints required JWT authentication',
			after: 'GET endpoints support API keys, write operations require JWT',
			improvement: 'Flexible authentication model'
		}
	},
	{
		id: 'domain-architecture',
		problem: 'Code organization and maintainability as system grows',
		context:
			'Multiple business domains (projects, skills, finances, etc.) needed clean separation to maintain code quality and enable team collaboration.',
		solution:
			'Implemented domain-driven design with each domain having its own package structure: repository (data access), service (business logic), handler (HTTP layer), and routes. Clear boundaries and interfaces.',
		technologies: ['Go', 'Domain-Driven Design', 'Clean Architecture'],
		impact:
			'Code is more maintainable, testable, and scalable. New domains can be added without affecting existing ones. Clear ownership and boundaries.',
		metrics: {
			before: 'Monolithic structure with mixed concerns',
			after: '8+ independent domains with clear boundaries',
			improvement: 'Modular, maintainable architecture'
		}
	}
];

