<script lang="ts">
	import { getResumeDownloadUrlById } from '$lib/api/resumes';
	import type { Resume } from '$lib/api/resumes';
	
	let {
		resume
	}: {
		resume: Resume;
	} = $props();

	let previewUrl = $derived(getResumeDownloadUrlById(resume.id));
</script>

<div class="section">
	<h2 class="section-title">Preview</h2>
	<div class="preview-container">
		<iframe
			src={previewUrl}
			class="preview-iframe"
			title="Resume Preview"
		/>
	</div>
</div>

<style>
	.preview-container {
		width: 100%;
		height: 800px;
		border: 1px solid #e5e7eb;
		border-radius: 0.375rem;
		overflow: hidden;
		background-color: #f9fafb;
	}

	.preview-iframe {
		width: 100%;
		height: 100%;
		border: none;
	}
</style>
