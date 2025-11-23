import type { CaseStudy } from '$lib/types/case-study';

export const caseStudies: CaseStudy[] = [
	{
		id: 'woragis-backend-architecture',
		projectSlug: 'woragis-backend',
		title: 'Woragis Backend: Microservices Architecture & API Design',
		problem:
			'Needed to build a scalable, maintainable backend system that could handle multiple domains (projects, skills, finances, ideas, chats) while maintaining clean separation of concerns and enabling future growth.',
		context:
			'The Woragis platform required a backend that could support a complex personal productivity application with multiple interconnected domains. The system needed to be modular, testable, and easy to extend as new features were added.',
		solution:
			'Designed and implemented a domain-driven architecture using Go Fiber, PostgreSQL, and Redis. Each domain (projects, skills, finances, etc.) is self-contained with its own entities, repositories, services, and handlers. The system uses JWT authentication and API key authentication for public access.',
		approach: [
			'Implemented domain-driven design with clear boundaries between domains',
			'Created a unified repository pattern using GORM for database operations',
			'Designed RESTful APIs with consistent response formats',
			'Implemented middleware for authentication, CORS, and request logging',
			'Added API key authentication for public read-only access',
			'Used Redis for caching and session management',
			'Implemented Docker containerization for easy deployment'
		],
		architecture: {
			diagram: `graph TB
    subgraph "Client Layer"
        A[Web Frontend]
        B[Mobile App]
        C[Landing Page]
    end
    
    subgraph "API Gateway"
        D[Fiber Router]
        E[Middleware Chain]
    end
    
    subgraph "Authentication"
        F[JWT Manager]
        G[API Key Validator]
        H[Session Manager]
    end
    
    subgraph "Domain Services"
        I[Projects Service]
        J[Skills Service]
        K[Finances Service]
        L[Ideas Service]
        M[Chats Service]
        N[Posts Service]
        O[Testimonials Service]
    end
    
    subgraph "Data Layer"
        P[(PostgreSQL)]
        Q[(Redis)]
    end
    
    subgraph "External Services"
        R[AI Service]
        S[Email Service]
        T[WhatsApp Service]
    end
    
    A --> D
    B --> D
    C --> D
    D --> E
    E --> F
    E --> G
    F --> H
    E --> I
    E --> J
    E --> K
    E --> L
    E --> M
    E --> N
    E --> O
    I --> P
    J --> P
    K --> P
    L --> P
    M --> P
    N --> P
    O --> P
    I --> Q
    M --> Q
    M --> R
    I --> S
    M --> T`,
			diagramType: 'mermaid',
			description:
				'The architecture follows a layered approach with clear separation between client, API gateway, authentication, domain services, and data layers.',
			components: [
				{
					name: 'API Gateway',
					description: 'Fiber router with middleware chain for request processing, authentication, and routing.',
					technologies: ['Go Fiber', 'Middleware']
				},
				{
					name: 'Domain Services',
					description: 'Self-contained services for each business domain with their own entities, repositories, and handlers.',
					technologies: ['Go', 'GORM', 'Domain-Driven Design']
				},
				{
					name: 'Data Layer',
					description: 'PostgreSQL for persistent storage and Redis for caching and session management.',
					technologies: ['PostgreSQL', 'Redis', 'GORM']
				},
				{
					name: 'External Services',
					description: 'Integration with AI services, email, and WhatsApp for extended functionality.',
					technologies: ['LangChain', 'SMTP', 'Whatsmeow']
				}
			]
		},
		metrics: {
			before: [
				{ label: 'Response Time', value: 'N/A (New System)' },
				{ label: 'Code Organization', value: 'Monolithic' },
				{ label: 'Testability', value: 'Low' }
			],
			after: [
				{ label: 'Average Response Time', value: '< 100ms' },
				{ label: 'Code Organization', value: 'Domain-Driven' },
				{ label: 'Testability', value: 'High (Unit & Integration)' },
				{ label: 'API Endpoints', value: '50+' },
				{ label: 'Domains', value: '10+' }
			],
			impact:
				'The modular architecture enables rapid feature development, easy testing, and clear separation of concerns. The system can scale horizontally and new domains can be added without affecting existing functionality.'
		},
		lessonsLearned: [
			'Domain-driven design significantly improves code maintainability and testability',
			'Middleware chains in Fiber provide powerful request processing capabilities',
			'API key authentication is essential for public-facing endpoints',
			'GORM simplifies database operations but requires careful relationship management',
			'Redis caching improves performance for frequently accessed data',
			'Docker containerization simplifies deployment and development setup'
		],
		technologies: ['Go', 'Fiber', 'PostgreSQL', 'Redis', 'GORM', 'Docker', 'JWT'],
		featured: true
	},
	{
		id: 'api-key-authentication',
		projectSlug: 'woragis-backend',
		title: 'Flexible API Key Authentication for Public Access',
		problem:
			'Needed to provide read-only access to specific backend endpoints (projects, skills, posts) for a public landing page without requiring full JWT authentication, while maintaining JWT for protected routes.',
		context:
			'The landing page needed to fetch projects and skills data from the backend API, but requiring JWT authentication would complicate the frontend implementation and limit public access.',
		solution:
			'Designed a custom Fiber middleware (`RequireAPIKeyOrAuth`) that first attempts API key validation for GET requests. If a valid API key is present, the request proceeds. Otherwise, it falls back to the existing JWT authentication middleware.',
		approach: [
			'Created API key domain with entity, repository, service, and handler layers',
			'Implemented API key validation middleware that extracts keys from headers',
			'Designed middleware to support both API key (GET) and JWT (all methods)',
			'Updated CORS configuration to allow custom API key headers',
			'Modified handlers to accept either API key-derived or JWT-derived user IDs'
		],
		architecture: {
			diagram: `sequenceDiagram
    participant Client
    participant Middleware
    participant APIKeyService
    participant JWTService
    participant Handler
    
    Client->>Middleware: GET /api/projects (X-API-Key header)
    Middleware->>APIKeyService: Validate API Key
    alt API Key Valid
        APIKeyService-->>Middleware: Valid API Key + UserID
        Middleware->>Middleware: Store UserID in Context
        Middleware->>Handler: Continue Request
        Handler-->>Client: 200 OK (Projects Data)
    else API Key Invalid/Missing
        Middleware->>JWTService: Validate JWT Token
        alt JWT Valid
            JWTService-->>Middleware: Valid JWT + UserID
            Middleware->>Handler: Continue Request
            Handler-->>Client: 200 OK (Projects Data)
        else JWT Invalid
            JWTService-->>Middleware: Unauthorized
            Middleware-->>Client: 401 Unauthorized
        end
    end`,
			diagramType: 'mermaid',
			description:
				'The middleware chain allows flexible authentication where API keys work for GET requests and JWT is required for all other operations.',
			components: [
				{
					name: 'API Key Middleware',
					description: 'Extracts and validates API keys from request headers.',
					technologies: ['Go Fiber', 'Middleware']
				},
				{
					name: 'JWT Middleware',
					description: 'Validates JWT tokens for authenticated requests.',
					technologies: ['Go Fiber', 'JWT']
				},
				{
					name: 'Context Storage',
					description: 'Stores authenticated user ID in Fiber context for handler access.',
					technologies: ['Go Fiber', 'Context']
				}
			]
		},
		metrics: {
			before: [
				{ label: 'Public API Access', value: 'Not Available' },
				{ label: 'Authentication Method', value: 'JWT Only' }
			],
			after: [
				{ label: 'Public API Access', value: 'Available via API Key' },
				{ label: 'Authentication Method', value: 'API Key (GET) or JWT (All)' },
				{ label: 'CORS Headers', value: 'Properly Configured' }
			],
			impact:
				'Enabled seamless integration between the landing page and backend API without requiring user authentication. The system maintains security for protected routes while providing convenient public access.'
		},
		lessonsLearned: [
			'Middleware order in Fiber is crucial for correct authentication flow',
			'CORS preflight headers are case-sensitive, requiring explicit allowance for both uppercase and lowercase header names',
			'Context propagation (e.g., UserID from API key) is essential for handlers to function correctly',
			'API keys should be scoped to read-only operations for security'
		],
		technologies: ['Go', 'Fiber', 'JWT', 'CORS', 'Middleware'],
		featured: true
	}
];

// Get case study by project slug
export function getCaseStudyByProjectSlug(slug: string): CaseStudy | undefined {
	return caseStudies.find((cs) => cs.projectSlug === slug);
}

// Get case study by ID
export function getCaseStudyById(id: string): CaseStudy | undefined {
	return caseStudies.find((cs) => cs.id === id);
}

// Get all featured case studies
export function getFeaturedCaseStudies(): CaseStudy[] {
	return caseStudies.filter((cs) => cs.featured);
}

