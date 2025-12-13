// Simple translation function - returns the key for now
// Can be extended later with proper i18n support
export function useTranslation() {
	return (key: string, fallback?: string): string => {
		// Simple translation map for common keys
		const translations: Record<string, string> = {
			'jobApplications.title': 'Job Applications',
			'jobApplications.subtitle': 'Manage your job applications',
			'jobApplications.createButton': 'Create Application',
			'jobApplications.searchPlaceholder': 'Search applications...',
			'jobApplications.empty': 'No applications found',
			'jobApplications.table.company': 'Company',
			'jobApplications.table.jobTitle': 'Job Title',
			'jobApplications.table.location': 'Location',
			'jobApplications.table.website': 'Website',
			'jobApplications.table.language': 'Language',
			'jobApplications.table.status': 'Status',
			'jobApplications.table.interest': 'Interest',
			'jobApplications.table.appliedAt': 'Applied At',
			'jobApplications.table.actions': 'Actions',
			'jobApplications.error': 'An error occurred',
			'jobApplications.loading': 'Loading applications...',
			'jobApplications.createError': 'Failed to create application',
			'jobApplications.deleteConfirm': 'Are you sure you want to delete this application?',
			'jobApplications.deleteError': 'Failed to delete application'
		};
		
		return translations[key] || fallback || key;
	};
}
