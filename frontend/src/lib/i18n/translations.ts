import type { Locale } from './store';

export type TranslationKey = 
	// Job Applications
	| 'jobApplications.title'
	| 'jobApplications.subtitle'
	| 'jobApplications.createButton'
	| 'jobApplications.searchPlaceholder'
	| 'jobApplications.loading'
	| 'jobApplications.empty'
	| 'jobApplications.error'
	| 'jobApplications.table.company'
	| 'jobApplications.table.jobTitle'
	| 'jobApplications.table.location'
	| 'jobApplications.table.website'
	| 'jobApplications.table.language'
	| 'jobApplications.table.status'
	| 'jobApplications.table.interest'
	| 'jobApplications.table.appliedAt'
	| 'jobApplications.table.actions'
	| 'jobApplications.table.view'
	| 'jobApplications.table.delete'
	| 'jobApplications.modal.createTitle'
	| 'jobApplications.modal.companyName'
	| 'jobApplications.modal.jobTitle'
	| 'jobApplications.modal.jobUrl'
	| 'jobApplications.modal.website'
	| 'jobApplications.modal.location'
	| 'jobApplications.modal.status'
	| 'jobApplications.modal.coverLetter'
	| 'jobApplications.modal.linkedInContact'
	| 'jobApplications.modal.create'
	| 'jobApplications.modal.cancel'
	| 'jobApplications.modal.required'
	| 'jobApplications.deleteConfirm'
	| 'jobApplications.createError'
	| 'jobApplications.deleteError'
	| 'jobApplications.status.pending'
	| 'jobApplications.status.processing'
	| 'jobApplications.status.applied'
	| 'jobApplications.status.contacted'
	| 'jobApplications.status.rejected'
	| 'jobApplications.status.accepted'
	| 'jobApplications.status.failed'
	// Projects
	| 'projects.title'
	| 'projects.subtitle'
	| 'projects.createButton'
	| 'projects.searchPlaceholder'
	| 'projects.loading'
	| 'projects.empty'
	| 'projects.error'
	| 'projects.table.name'
	| 'projects.table.status'
	| 'projects.table.healthScore'
	| 'projects.table.mrr'
	| 'projects.table.cac'
	| 'projects.table.ltv'
	| 'projects.table.churnRate'
	| 'projects.table.created'
	| 'projects.table.actions'
	| 'projects.table.view'
	| 'projects.table.delete'
	| 'projects.modal.createTitle'
	| 'projects.modal.name'
	| 'projects.modal.description'
	| 'projects.modal.status'
	| 'projects.modal.healthScore'
	| 'projects.modal.mrr'
	| 'projects.modal.cac'
	| 'projects.modal.ltv'
	| 'projects.modal.churnRate'
	| 'projects.modal.create'
	| 'projects.modal.cancel'
	| 'projects.modal.required'
	| 'projects.deleteConfirm'
	| 'projects.createError'
	| 'projects.deleteError'
	| 'projects.status.idea'
	| 'projects.status.planning'
	| 'projects.status.executing'
	| 'projects.status.monitoring'
	| 'projects.status.completed'
	// Resumes
	| 'resumes.title'
	| 'resumes.subtitle'
	| 'resumes.createButton'
	| 'resumes.loading'
	| 'resumes.empty'
	| 'resumes.emptySubtext'
	| 'resumes.error'
	| 'resumes.table.title'
	| 'resumes.table.file'
	| 'resumes.table.status'
	| 'resumes.table.created'
	| 'resumes.table.actions'
	| 'resumes.table.main'
	| 'resumes.table.featured'
	| 'resumes.modal.createTitle'
	| 'resumes.modal.title'
	| 'resumes.modal.file'
	| 'resumes.modal.filePath'
	| 'resumes.modal.fileName'
	| 'resumes.modal.fileSize'
	| 'resumes.modal.chooseFile'
	| 'resumes.modal.manualEntry'
	| 'resumes.modal.create'
	| 'resumes.modal.cancel'
	| 'resumes.modal.creating'
	| 'resumes.modal.required'
	| 'resumes.deleteConfirm'
	| 'resumes.createError'
	| 'resumes.deleteError'
	| 'resumes.uploadSuccess'
	| 'resumes.createSuccess'
	| 'resumes.deleteSuccess'
	| 'resumes.loadError'
	// Common
	| 'common.required'
	| 'common.cancel'
	| 'common.create'
	| 'common.delete'
	| 'common.view'
	| 'common.loading'
	| 'common.error'
	| 'common.empty';

export type Translations = Record<TranslationKey, string>;

const translations: Record<Locale, Translations> = {
	en: {
		// Job Applications
		'jobApplications.title': 'Job Applications Management',
		'jobApplications.subtitle': 'Manage job applications, responses, and interview stages',
		'jobApplications.createButton': 'Create Application',
		'jobApplications.searchPlaceholder': 'Search applications...',
		'jobApplications.loading': 'Loading...',
		'jobApplications.empty': 'No job applications found',
		'jobApplications.error': 'Failed to load job applications',
		'jobApplications.table.company': 'Company',
		'jobApplications.table.jobTitle': 'Job Title',
		'jobApplications.table.location': 'Location',
		'jobApplications.table.website': 'Website',
		'jobApplications.table.language': 'Language',
		'jobApplications.table.status': 'Status',
		'jobApplications.table.interest': 'Interest',
		'jobApplications.table.appliedAt': 'Applied At',
		'jobApplications.table.actions': 'Actions',
		'jobApplications.table.view': 'View',
		'jobApplications.table.delete': 'Delete',
		'jobApplications.modal.createTitle': 'Create Job Application',
		'jobApplications.modal.companyName': 'Company Name',
		'jobApplications.modal.jobTitle': 'Job Title',
		'jobApplications.modal.jobUrl': 'Job URL',
		'jobApplications.modal.website': 'Website',
		'jobApplications.modal.location': 'Location',
		'jobApplications.modal.status': 'Status',
		'jobApplications.modal.coverLetter': 'Cover Letter',
		'jobApplications.modal.linkedInContact': 'LinkedIn Contact',
		'jobApplications.modal.create': 'Create',
		'jobApplications.modal.cancel': 'Cancel',
		'jobApplications.modal.required': '*',
		'jobApplications.deleteConfirm': 'Are you sure you want to delete this job application?',
		'jobApplications.createError': 'Failed to create job application',
		'jobApplications.deleteError': 'Failed to delete job application',
		'jobApplications.status.pending': 'pending',
		'jobApplications.status.processing': 'processing',
		'jobApplications.status.applied': 'applied',
		'jobApplications.status.contacted': 'contacted',
		'jobApplications.status.rejected': 'rejected',
		'jobApplications.status.accepted': 'accepted',
		'jobApplications.status.failed': 'failed',
		// Projects
		'projects.title': 'Projects Management',
		'projects.subtitle': 'Manage projects',
		'projects.createButton': 'Create Project',
		'projects.searchPlaceholder': 'Search projects...',
		'projects.loading': 'Loading...',
		'projects.empty': 'No projects found',
		'projects.error': 'Failed to load projects',
		'projects.table.name': 'Name',
		'projects.table.status': 'Status',
		'projects.table.healthScore': 'Health Score',
		'projects.table.mrr': 'MRR',
		'projects.table.cac': 'CAC',
		'projects.table.ltv': 'LTV',
		'projects.table.churnRate': 'Churn Rate',
		'projects.table.created': 'Created',
		'projects.table.actions': 'Actions',
		'projects.table.view': 'View',
		'projects.table.delete': 'Delete',
		'projects.modal.createTitle': 'Create Project',
		'projects.modal.name': 'Name',
		'projects.modal.description': 'Description',
		'projects.modal.status': 'Status',
		'projects.modal.healthScore': 'Health Score',
		'projects.modal.mrr': 'MRR (Monthly Recurring Revenue)',
		'projects.modal.cac': 'CAC (Customer Acquisition Cost)',
		'projects.modal.ltv': 'LTV (Lifetime Value)',
		'projects.modal.churnRate': 'Churn Rate (%)',
		'projects.modal.create': 'Create',
		'projects.modal.cancel': 'Cancel',
		'projects.modal.required': '*',
		'projects.deleteConfirm': 'Are you sure you want to delete this project? This will also delete all related milestones, kanban boards, dependencies, documentation, technologies, file structures, and architecture diagrams.',
		'projects.createError': 'Failed to create project',
		'projects.deleteError': 'Failed to delete project',
		'projects.status.idea': 'idea',
		'projects.status.planning': 'planning',
		'projects.status.executing': 'executing',
		'projects.status.monitoring': 'monitoring',
		'projects.status.completed': 'completed',
		// Resumes
		'resumes.title': 'Resumes',
		'resumes.subtitle': 'Manage your resume files',
		'resumes.createButton': 'Create Resume',
		'resumes.loading': 'Loading resumes...',
		'resumes.empty': 'No resumes found',
		'resumes.emptySubtext': 'Resumes are created when you generate them using the resume-worker',
		'resumes.error': 'Failed to load resumes',
		'resumes.table.title': 'Title',
		'resumes.table.file': 'File',
		'resumes.table.status': 'Status',
		'resumes.table.created': 'Created',
		'resumes.table.actions': 'Actions',
		'resumes.table.main': 'Main',
		'resumes.table.featured': 'Featured',
		'resumes.modal.createTitle': 'Create Resume',
		'resumes.modal.title': 'Title',
		'resumes.modal.file': 'Resume File (PDF)',
		'resumes.modal.filePath': 'File Path',
		'resumes.modal.fileName': 'File Name',
		'resumes.modal.fileSize': 'File Size (bytes)',
		'resumes.modal.chooseFile': 'Choose File',
		'resumes.modal.manualEntry': 'Or manually enter file details below',
		'resumes.modal.create': 'Create',
		'resumes.modal.cancel': 'Cancel',
		'resumes.modal.creating': 'Creating...',
		'resumes.modal.required': '*',
		'resumes.deleteConfirm': 'Are you sure you want to delete',
		'resumes.createError': 'Failed to create resume',
		'resumes.deleteError': 'Failed to delete resume',
		'resumes.uploadSuccess': 'Resume uploaded successfully',
		'resumes.createSuccess': 'Resume created successfully',
		'resumes.deleteSuccess': 'Resume deleted successfully',
		'resumes.loadError': 'Failed to load resumes',
		// Common
		'common.required': '*',
		'common.cancel': 'Cancel',
		'common.create': 'Create',
		'common.delete': 'Delete',
		'common.view': 'View',
		'common.loading': 'Loading...',
		'common.error': 'Error',
		'common.empty': 'No items found'
	},
	pt: {
		// Job Applications
		'jobApplications.title': 'Gerenciamento de Candidaturas',
		'jobApplications.subtitle': 'Gerencie candidaturas, respostas e etapas de entrevista',
		'jobApplications.createButton': 'Criar Candidatura',
		'jobApplications.searchPlaceholder': 'Buscar candidaturas...',
		'jobApplications.loading': 'Carregando...',
		'jobApplications.empty': 'Nenhuma candidatura encontrada',
		'jobApplications.error': 'Falha ao carregar candidaturas',
		'jobApplications.table.company': 'Empresa',
		'jobApplications.table.jobTitle': 'Cargo',
		'jobApplications.table.location': 'Localização',
		'jobApplications.table.website': 'Website',
		'jobApplications.table.language': 'Idioma',
		'jobApplications.table.status': 'Status',
		'jobApplications.table.interest': 'Interesse',
		'jobApplications.table.appliedAt': 'Candidatado em',
		'jobApplications.table.actions': 'Ações',
		'jobApplications.table.view': 'Ver',
		'jobApplications.table.delete': 'Excluir',
		'jobApplications.modal.createTitle': 'Criar Candidatura',
		'jobApplications.modal.companyName': 'Nome da Empresa',
		'jobApplications.modal.jobTitle': 'Cargo',
		'jobApplications.modal.jobUrl': 'URL da Vaga',
		'jobApplications.modal.website': 'Website',
		'jobApplications.modal.location': 'Localização',
		'jobApplications.modal.status': 'Status',
		'jobApplications.modal.coverLetter': 'Carta de Apresentação',
		'jobApplications.modal.linkedInContact': 'Contato LinkedIn',
		'jobApplications.modal.create': 'Criar',
		'jobApplications.modal.cancel': 'Cancelar',
		'jobApplications.modal.required': '*',
		'jobApplications.deleteConfirm': 'Tem certeza que deseja excluir esta candidatura?',
		'jobApplications.createError': 'Falha ao criar candidatura',
		'jobApplications.deleteError': 'Falha ao excluir candidatura',
		'jobApplications.status.pending': 'pendente',
		'jobApplications.status.processing': 'processando',
		'jobApplications.status.applied': 'candidatado',
		'jobApplications.status.contacted': 'contatado',
		'jobApplications.status.rejected': 'rejeitado',
		'jobApplications.status.accepted': 'aceito',
		'jobApplications.status.failed': 'falhou',
		// Projects
		'projects.title': 'Gerenciamento de Projetos',
		'projects.subtitle': 'Gerencie projetos',
		'projects.createButton': 'Criar Projeto',
		'projects.searchPlaceholder': 'Buscar projetos...',
		'projects.loading': 'Carregando...',
		'projects.empty': 'Nenhum projeto encontrado',
		'projects.error': 'Falha ao carregar projetos',
		'projects.table.name': 'Nome',
		'projects.table.status': 'Status',
		'projects.table.healthScore': 'Pontuação de Saúde',
		'projects.table.mrr': 'MRR',
		'projects.table.cac': 'CAC',
		'projects.table.ltv': 'LTV',
		'projects.table.churnRate': 'Taxa de Cancelamento',
		'projects.table.created': 'Criado',
		'projects.table.actions': 'Ações',
		'projects.table.view': 'Ver',
		'projects.table.delete': 'Excluir',
		'projects.modal.createTitle': 'Criar Projeto',
		'projects.modal.name': 'Nome',
		'projects.modal.description': 'Descrição',
		'projects.modal.status': 'Status',
		'projects.modal.healthScore': 'Pontuação de Saúde',
		'projects.modal.mrr': 'MRR (Receita Recorrente Mensal)',
		'projects.modal.cac': 'CAC (Custo de Aquisição de Cliente)',
		'projects.modal.ltv': 'LTV (Valor do Tempo de Vida)',
		'projects.modal.churnRate': 'Taxa de Cancelamento (%)',
		'projects.modal.create': 'Criar',
		'projects.modal.cancel': 'Cancelar',
		'projects.modal.required': '*',
		'projects.deleteConfirm': 'Tem certeza que deseja excluir este projeto? Isso também excluirá todos os marcos relacionados, quadros kanban, dependências, documentação, tecnologias, estruturas de arquivos e diagramas de arquitetura.',
		'projects.createError': 'Falha ao criar projeto',
		'projects.deleteError': 'Falha ao excluir projeto',
		'projects.status.idea': 'ideia',
		'projects.status.planning': 'planejamento',
		'projects.status.executing': 'executando',
		'projects.status.monitoring': 'monitorando',
		'projects.status.completed': 'concluído',
		// Resumes
		'resumes.title': 'Currículos',
		'resumes.subtitle': 'Gerencie seus arquivos de currículo',
		'resumes.createButton': 'Criar Currículo',
		'resumes.loading': 'Carregando currículos...',
		'resumes.empty': 'Nenhum currículo encontrado',
		'resumes.emptySubtext': 'Currículos são criados quando você os gera usando o resume-worker',
		'resumes.error': 'Falha ao carregar currículos',
		'resumes.table.title': 'Título',
		'resumes.table.file': 'Arquivo',
		'resumes.table.status': 'Status',
		'resumes.table.created': 'Criado',
		'resumes.table.actions': 'Ações',
		'resumes.table.main': 'Principal',
		'resumes.table.featured': 'Destaque',
		'resumes.modal.createTitle': 'Criar Currículo',
		'resumes.modal.title': 'Título',
		'resumes.modal.file': 'Arquivo de Currículo (PDF)',
		'resumes.modal.filePath': 'Caminho do Arquivo',
		'resumes.modal.fileName': 'Nome do Arquivo',
		'resumes.modal.fileSize': 'Tamanho do Arquivo (bytes)',
		'resumes.modal.chooseFile': 'Escolher Arquivo',
		'resumes.modal.manualEntry': 'Ou insira os detalhes do arquivo manualmente abaixo',
		'resumes.modal.create': 'Criar',
		'resumes.modal.cancel': 'Cancelar',
		'resumes.modal.creating': 'Criando...',
		'resumes.modal.required': '*',
		'resumes.deleteConfirm': 'Tem certeza que deseja excluir',
		'resumes.createError': 'Falha ao criar currículo',
		'resumes.deleteError': 'Falha ao excluir currículo',
		'resumes.uploadSuccess': 'Currículo enviado com sucesso',
		'resumes.createSuccess': 'Currículo criado com sucesso',
		'resumes.deleteSuccess': 'Currículo excluído com sucesso',
		'resumes.loadError': 'Falha ao carregar currículos',
		// Common
		'common.required': '*',
		'common.cancel': 'Cancelar',
		'common.create': 'Criar',
		'common.delete': 'Excluir',
		'common.view': 'Ver',
		'common.loading': 'Carregando...',
		'common.error': 'Erro',
		'common.empty': 'Nenhum item encontrado'
	}
};

export function getTranslation(locale: Locale, key: TranslationKey): string {
	return translations[locale][key] || key;
}

