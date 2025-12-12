import { api } from '$lib/constants';

// Resume download function - triggers download
export async function downloadResume(userId: string, language: string = 'en'): Promise<void> {
	try {
		const params = new URLSearchParams({
			userId,
			language
		});
		
		// Create download URL
		const downloadUrl = `${api.baseURL}/public/resume/download?${params.toString()}`;
		
		// Create a temporary anchor element to trigger download
		const link = document.createElement('a');
		link.href = downloadUrl;
		link.download = ''; // Let the server set the filename via Content-Disposition
		link.target = '_blank';
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	} catch (error) {
		console.error('Error downloading resume:', error);
		throw error;
	}
}

// Resume preview function - opens in new window/tab for preview
export async function previewResume(userId: string, language: string = 'en'): Promise<void> {
	try {
		const params = new URLSearchParams({
			userId,
			language
		});
		
		// Create preview URL
		const previewUrl = `${api.baseURL}/public/resume/preview?${params.toString()}`;
		
		// Open in new window for preview
		window.open(previewUrl, '_blank');
	} catch (error) {
		console.error('Error previewing resume:', error);
		throw error;
	}
}

