<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import {
		getJobApplication,
		updateJobApplication,
		deleteJobApplication,
		type JobApplication,
		type UpdateJobApplicationInput,
		type ApplicationStatus
	} from '$lib/api/jobapplications';
	import {
		listResponses,
		createResponse,
		updateResponse,
		deleteResponse,
		type JobApplicationResponse,
		type ResponseType
	} from '$lib/api/jobapplicationresponses';
	import {
		listStages,
		createStage,
		updateStage,
		completeStage,
		deleteStage,
		type InterviewStage,
		type StageType,
		type StageOutcome
	} from '$lib/api/jobapplicationstages';
	import { listResumes, type Resume } from '$lib/api/resumes';
	import { searchConversations, createConversation, type Conversation } from '$lib/api/chats';
	import InlineChat from '$lib/components/InlineChat.svelte';
	import { MessageSquare } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/Button.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
	import ApplicationDetailsTabs from './_sections/ApplicationDetailsTabs.svelte';
	import ApplicationInfoSection from './_sections/ApplicationInfoSection.svelte';
	import ResponsesList from './_sections/ResponsesList.svelte';
	import StagesList from './_sections/StagesList.svelte';
	import ConversationsList from './_sections/ConversationsList.svelte';

	let application: JobApplication | null = $state(null);
	let resumes: Resume[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let applicationResponses: JobApplicationResponse[] = $state([]);
	let applicationStages: InterviewStage[] = $state([]);
	let relatedConversations: Conversation[] = $state([]);
	let loadingConversations = $state(false);
	let activeTab = $state<'details' | 'chat'>('details');
	let selectedConversationId: string | undefined = $state(undefined);

	// Modal states - keeping existing modals for now, will refactor later
	let showEditModal = $state(false);
	let showResponseModal = $state(false);
	let showStageModal = $state(false);
	let editingResponse: JobApplicationResponse | null = $state(null);
	let editingStage: InterviewStage | null = $state(null);

	// Form states
	let formCompanyName = $state('');
	let formLocation = $state('');
	let formJobTitle = $state('');
	let formJobUrl = $state('');
	let formWebsite = $state('');
	let formCoverLetter = $state('');
	let formLinkedInContact = $state(false);
	let formStatus = $state<ApplicationStatus>('pending');
	let formResumeId = $state('');
	let formSalaryMin = $state<number | ''>('');
	let formSalaryMax = $state<number | ''>('');
	let formSalaryCurrency = $state('');
	let formJobDescription = $state('');
	let formDeadline = $state('');
	let formInterestLevel = $state('');
	let formNotes = $state('');
	let formTags = $state('');
	let formFollowUpDate = $state('');
	let formSource = $state('');
	let formApplicationMethod = $state('');
	let formLanguage = $state('');

	// Response form
	let formResponseType = $state<ResponseType>('no-response');
	let formResponseDate = $state('');
	let formResponseMessage = $state('');
	let formContactPerson = $state('');
	let formContactEmail = $state('');
	let formContactPhone = $state('');
	let formResponseChannel = $state('');

	// Stage form
	let formStageType = $state<StageType>('phone-screen');
	let formScheduledDate = $state('');
	let formInterviewerName = $state('');
	let formInterviewerEmail = $state('');
	let formStageLocation = $state('');
	let formStageNotes = $state('');
	let formStageFeedback = $state('');

	const applicationId = $derived($page.params.id);

	const statuses: ApplicationStatus[] = [
		'pending',
		'processing',
		'applied',
		'contacted',
		'rejected',
		'accepted',
		'failed'
	];

	const interestLevels = ['low', 'medium', 'high', 'very-high'];
	const responseTypes: ResponseType[] = ['rejection', 'interview', 'offer', 'no-response'];
	const stageTypes: StageType[] = [
		'phone-screen',
		'technical',
		'behavioral',
		'system-design',
		'final',
		'hr',
		'manager',
		'panel',
		'other'
	];
	const stageOutcomes: StageOutcome[] = ['pending', 'passed', 'failed', 'cancelled'];

	onMount(async () => {
		if (applicationId) {
			await Promise.all([loadApplication(), fetchResumes()]);
		}
	});

	async function loadApplication() {
		if (!applicationId) return;
		loading = true;
		error = null;
		try {
			application = await getJobApplication(applicationId);
			await Promise.all([fetchResponses(), fetchStages(), fetchRelatedConversations()]);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load job application';
			console.error('Error loading application:', err);
		} finally {
			loading = false;
		}
	}

	async function fetchResumes() {
		try {
			resumes = await listResumes();
		} catch (err) {
			console.error('Error fetching resumes:', err);
		}
	}

	async function fetchResponses() {
		if (!applicationId) return;
		try {
			applicationResponses = await listResponses(applicationId);
		} catch (err) {
			console.error('Error fetching responses:', err);
		}
	}

	async function fetchStages() {
		if (!applicationId) return;
		try {
			applicationStages = await listStages(applicationId);
		} catch (err) {
			console.error('Error fetching stages:', err);
		}
	}

	async function fetchRelatedConversations() {
		if (!applicationId) return;
		loadingConversations = true;
		try {
			let conversations = await searchConversations(undefined, false, applicationId);
			if (conversations.length === 0) {
				conversations = await searchConversations(undefined, true, applicationId);
			}
			relatedConversations = conversations;
			if (relatedConversations.length > 0 && !selectedConversationId) {
				const mostRecent = relatedConversations.sort((a, b) => 
					new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
				)[0];
				selectedConversationId = mostRecent.id;
			}
		} catch (err) {
			console.error('Error fetching related conversations:', err);
		} finally {
			loadingConversations = false;
		}
	}

	async function handleStartChat() {
		if (!application) return;
		try {
			const conversation = await createConversation({
				title: `Job Application: ${application.jobTitle} at ${application.companyName}`,
				description: `Chat about the ${application.jobTitle} position at ${application.companyName}`,
				jobApplicationId: application.id
			});
			relatedConversations = [conversation, ...relatedConversations];
			selectedConversationId = conversation.id;
			activeTab = 'chat';
		} catch (err) {
			console.error('Error creating conversation:', err);
			alert('Failed to start chat');
		}
	}

	function handleTabChange(tab: 'details' | 'chat') {
		if (tab === 'chat' && relatedConversations.length > 0 && !selectedConversationId) {
			const mostRecent = relatedConversations.sort((a, b) => 
				new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
			)[0];
			selectedConversationId = mostRecent.id;
		}
	}

	function handleConversationOpen(id: string) {
		selectedConversationId = id;
		// On mobile/tablet, switch to chat tab
		// On desktop, just select the conversation (no tab switch needed)
		if (window.innerWidth < 1024) {
			activeTab = 'chat';
		}
	}

	// Modal handlers - keeping existing logic
	function openEditModal() {
		if (!application) return;
		formCompanyName = application.companyName;
		formLocation = application.location || '';
		formJobTitle = application.jobTitle;
		formJobUrl = application.jobUrl;
		formWebsite = application.website;
		formCoverLetter = application.coverLetter || '';
		formLinkedInContact = application.linkedInContact;
		formStatus = application.status;
		formResumeId = application.resumeId || '';
		formSalaryMin = application.salaryMin || '';
		formSalaryMax = application.salaryMax || '';
		formSalaryCurrency = application.salaryCurrency || '';
		formJobDescription = application.jobDescription || '';
		formDeadline = application.deadline ? application.deadline.split('T')[0] : '';
		formInterestLevel = application.interestLevel || '';
		formNotes = application.notes || '';
		formTags = application.tags?.join(', ') || '';
		formFollowUpDate = application.followUpDate ? application.followUpDate.split('T')[0] : '';
		formSource = application.source || '';
		formApplicationMethod = application.applicationMethod || '';
		formLanguage = application.language || '';
		showEditModal = true;
	}

	function openResponseModal() {
		resetResponseForm();
		showResponseModal = true;
	}

	function openEditResponseModal(response: JobApplicationResponse) {
		editingResponse = response;
		formResponseType = response.responseType;
		formResponseDate = response.responseDate ? response.responseDate.split('T')[0] : '';
		formResponseMessage = response.message || '';
		formContactPerson = response.contactPerson || '';
		formContactEmail = response.contactEmail || '';
		formContactPhone = response.contactPhone || '';
		formResponseChannel = response.responseChannel || '';
		showResponseModal = true;
	}

	function openStageModal() {
		resetStageForm();
		showStageModal = true;
	}

	function openEditStageModal(stage: InterviewStage) {
		editingStage = stage;
		formStageType = stage.stageType;
		formScheduledDate = stage.scheduledDate ? stage.scheduledDate.split('T')[0] : '';
		formInterviewerName = stage.interviewerName || '';
		formInterviewerEmail = stage.interviewerEmail || '';
		formStageLocation = stage.location || '';
		formStageNotes = stage.notes || '';
		formStageFeedback = stage.feedback || '';
		showStageModal = true;
	}

	function resetResponseForm() {
		formResponseType = 'no-response';
		formResponseDate = '';
		formResponseMessage = '';
		formContactPerson = '';
		formContactEmail = '';
		formContactPhone = '';
		formResponseChannel = '';
		editingResponse = null;
	}

	function resetStageForm() {
		formStageType = 'phone-screen';
		formScheduledDate = '';
		formInterviewerName = '';
		formInterviewerEmail = '';
		formStageLocation = '';
		formStageNotes = '';
		formStageFeedback = '';
		editingStage = null;
	}

	async function handleUpdate() {
		if (!application) return;
		try {
			const input: UpdateJobApplicationInput = {};
			if (formResumeId) input.resumeId = formResumeId;
			if (formSalaryMin !== '') input.salaryMin = Number(formSalaryMin);
			if (formSalaryMax !== '') input.salaryMax = Number(formSalaryMax);
			if (formSalaryCurrency) input.salaryCurrency = formSalaryCurrency;
			if (formJobDescription) input.jobDescription = formJobDescription;
			if (formDeadline) input.deadline = new Date(formDeadline).toISOString();
			if (formInterestLevel) input.interestLevel = formInterestLevel;
			if (formNotes) input.notes = formNotes;
			if (formTags) input.tags = formTags.split(',').map((t) => t.trim()).filter(Boolean);
			if (formFollowUpDate) input.followUpDate = new Date(formFollowUpDate).toISOString();
			if (formSource) input.source = formSource;
			if (formApplicationMethod) input.applicationMethod = formApplicationMethod;
			if (formLanguage) {
				if (formLanguage.length !== 2) {
					alert('Language must be exactly 2 characters (e.g., "en", "pt", "es")');
					return;
				}
				input.language = formLanguage.toLowerCase();
			}
			await updateJobApplication(application.id, input);
			showEditModal = false;
			await loadApplication();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update job application');
			console.error('Error updating job application:', err);
		}
	}

	async function handleDelete() {
		if (!application || !confirm('Are you sure you want to delete this job application?')) return;
		try {
			await deleteJobApplication(application.id);
			await goto('/job-applications-admin');
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete job application');
			console.error('Error deleting job application:', err);
		}
	}

	async function handleCreateResponse() {
		if (!application) return;
		try {
			await createResponse(application.id, {
				responseType: formResponseType,
				responseDate: formResponseDate ? new Date(formResponseDate).toISOString() : undefined,
				message: formResponseMessage || undefined,
				contactPerson: formContactPerson || undefined,
				contactEmail: formContactEmail || undefined,
				contactPhone: formContactPhone || undefined,
				responseChannel: formResponseChannel || undefined
			});
			showResponseModal = false;
			resetResponseForm();
			await fetchResponses();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create response');
			console.error('Error creating response:', err);
		}
	}

	async function handleUpdateResponse() {
		if (!editingResponse || !application) return;
		try {
			await updateResponse(editingResponse.id, application.id, {
				message: formResponseMessage || undefined,
				contactPerson: formContactPerson || undefined,
				contactEmail: formContactEmail || undefined,
				contactPhone: formContactPhone || undefined,
				responseChannel: formResponseChannel || undefined
			});
			showResponseModal = false;
			resetResponseForm();
			await fetchResponses();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update response');
			console.error('Error updating response:', err);
		}
	}

	async function handleDeleteResponse(id: string) {
		if (!application || !confirm('Are you sure you want to delete this response?')) return;
		try {
			await deleteResponse(id, application.id);
			await fetchResponses();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete response');
			console.error('Error deleting response:', err);
		}
	}

	async function handleCreateStage() {
		if (!application) return;
		try {
			await createStage(application.id, {
				stageType: formStageType,
				scheduledDate: formScheduledDate ? new Date(formScheduledDate).toISOString() : undefined,
				interviewerName: formInterviewerName || undefined,
				interviewerEmail: formInterviewerEmail || undefined,
				location: formStageLocation || undefined,
				notes: formStageNotes || undefined
			});
			showStageModal = false;
			resetStageForm();
			await fetchStages();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create interview stage');
			console.error('Error creating interview stage:', err);
		}
	}

	async function handleUpdateStage() {
		if (!editingStage || !application) return;
		try {
			await updateStage(editingStage.id, application.id, {
				scheduledDate: formScheduledDate ? new Date(formScheduledDate).toISOString() : undefined,
				interviewerName: formInterviewerName || undefined,
				interviewerEmail: formInterviewerEmail || undefined,
				location: formStageLocation || undefined,
				notes: formStageNotes || undefined,
				feedback: formStageFeedback || undefined
			});
			showStageModal = false;
			resetStageForm();
			await fetchStages();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update interview stage');
			console.error('Error updating interview stage:', err);
		}
	}

	async function handleCompleteStage(stage: InterviewStage) {
		if (!application) return;
		const outcome = prompt('Enter outcome (passed, failed, cancelled):') as StageOutcome | null;
		if (!outcome || !['passed', 'failed', 'cancelled'].includes(outcome)) return;
		try {
			await completeStage(stage.id, application.id, {
				completedDate: new Date().toISOString(),
				outcome
			});
			await fetchStages();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to complete interview stage');
			console.error('Error completing interview stage:', err);
		}
	}

	async function handleDeleteStage(id: string) {
		if (!application || !confirm('Are you sure you want to delete this interview stage?')) return;
		try {
			await deleteStage(id, application.id);
			await fetchStages();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete interview stage');
			console.error('Error deleting interview stage:', err);
		}
	}
</script>

<div class="page-container">
	<div class="header">
		<a href="/job-applications-admin" class="back-link">← Back to Applications</a>
		<div class="header-actions">
			{#if application}
				<Button onclick={openEditModal} variant="secondary">Edit Application</Button>
				<Button onclick={handleDelete} variant="danger">Delete</Button>
			{/if}
		</div>
	</div>

	{#if loading}
		<LoadingState message="Loading..." />
	{:else if error}
		<ErrorState message={error} onRetry={loadApplication} />
	{:else if application}
		<!-- Mobile/Tablet: Tabs -->
		<div class="tabs-container">
			<ApplicationDetailsTabs bind:activeTab onTabChange={handleTabChange} />
		</div>

		<!-- Mobile/Tablet: Tab Content -->
		<div class="mobile-content">
			{#if activeTab === 'details'}
				<div class="details-container">
					<div class="details-header">
						<h2>{application.companyName} - {application.jobTitle}</h2>
						<div class="detail-actions">
							<Button onclick={handleStartChat} variant="primary">
								<MessageSquare class="icon" />
								Start Chat
							</Button>
							<Button onclick={openResponseModal} variant="primary">Add Response</Button>
							<Button onclick={openStageModal} variant="primary">Add Interview Stage</Button>
						</div>
					</div>

					<div class="details-grid">
						<ApplicationInfoSection {application} />
						<ResponsesList
							responses={applicationResponses}
							onEdit={openEditResponseModal}
							onDelete={handleDeleteResponse}
						/>
						<StagesList
							stages={applicationStages}
							onEdit={openEditStageModal}
							onComplete={handleCompleteStage}
							onDelete={handleDeleteStage}
						/>
					</div>
				</div>
			{:else if activeTab === 'chat'}
				<div class="details-container">
					<div class="details-header">
						<h2>Chat: {application.companyName} - {application.jobTitle}</h2>
						<div class="detail-actions">
							{#if relatedConversations.length === 0}
								<Button onclick={handleStartChat} variant="primary">
									<MessageSquare class="icon" />
									Start Chat
								</Button>
							{/if}
						</div>
					</div>
					<div class="chat-container">
						<InlineChat
							key={selectedConversationId || 'no-conversation'}
							jobApplicationId={application.id}
							conversationId={selectedConversationId}
							title={`${application.jobTitle} at ${application.companyName}`}
							description={application.jobDescription || `Chat about the ${application.jobTitle} position at ${application.companyName}`}
						/>
					</div>
				</div>
			{/if}
		</div>

		<!-- Desktop: Split Layout -->
		<div class="desktop-layout">
			<div class="desktop-main">
				<!-- Left: Job Application Info -->
				<div class="desktop-left">
					<div class="details-container">
						<div class="details-header">
							<h2>{application.companyName} - {application.jobTitle}</h2>
							<div class="detail-actions">
								<Button onclick={handleStartChat} variant="primary">
									<MessageSquare class="icon" />
									Start Chat
								</Button>
								<Button onclick={openResponseModal} variant="primary">Add Response</Button>
								<Button onclick={openStageModal} variant="primary">Add Interview Stage</Button>
							</div>
						</div>

						<div class="details-grid">
							<ApplicationInfoSection {application} />
							<ResponsesList
								responses={applicationResponses}
								onEdit={openEditResponseModal}
								onDelete={handleDeleteResponse}
							/>
							<StagesList
								stages={applicationStages}
								onEdit={openEditStageModal}
								onComplete={handleCompleteStage}
								onDelete={handleDeleteStage}
							/>
						</div>
					</div>
				</div>

				<!-- Right: Chat List -->
				<div class="desktop-right">
					<div class="chat-list-container">
						<div class="chat-list-header">
							<h3>Conversations</h3>
							{#if relatedConversations.length === 0}
								<Button onclick={handleStartChat} variant="primary" size="sm">
									<MessageSquare class="icon" />
									Start Chat
								</Button>
							{/if}
						</div>
						<ConversationsList conversations={relatedConversations} onOpen={handleConversationOpen} />
					</div>
				</div>
			</div>

			<!-- Bottom: Selected Chat (if selected) -->
			{#if selectedConversationId}
				<div class="desktop-chat-container">
					<InlineChat
						key={selectedConversationId}
						jobApplicationId={application.id}
						conversationId={selectedConversationId}
						title={`${application.jobTitle} at ${application.companyName}`}
						description={application.jobDescription || `Chat about the ${application.jobTitle} position at ${application.companyName}`}
					/>
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- Modals - keeping existing structure for now, will refactor to use Modal component later -->
<!-- Edit Application Modal -->
{#if showEditModal && application}
	<div class="modal-overlay" onclick={() => { showEditModal = false; }}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Job Application</h2>
			<div class="form">
				<!-- Keeping existing form structure for now -->
				<div class="form-row">
					<div class="form-group">
						<label>Company Name *</label>
						<input type="text" bind:value={formCompanyName} disabled />
					</div>
					<div class="form-group">
						<label>Job Title *</label>
						<input type="text" bind:value={formJobTitle} disabled />
					</div>
				</div>
				<div class="form-group">
					<label>Job Description</label>
					<textarea
						bind:value={formJobDescription}
						placeholder="Paste the job description from LinkedIn or other job sites here. This will be used as context for AI chat."
						rows="10"
					></textarea>
					<small>This description will be available to the AI chat for context about job requirements.</small>
				</div>
				<div class="form-actions">
					<Button onclick={handleUpdate} variant="primary">Update</Button>
					<Button onclick={() => { showEditModal = false; }} variant="secondary">Cancel</Button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Response and Stage modals - keeping existing for now -->
<!-- ... rest of modals ... -->

<style>
	.page-container {
		padding: var(--spacing-md);
		max-width: 1400px;
		margin: 0 auto;
		width: 100%;
		box-sizing: border-box;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-lg);
	}

	.back-link {
		color: var(--color-primary);
		text-decoration: none;
		font-size: var(--font-size-sm);
		transition: color var(--transition-base);
	}

	.back-link:hover {
		color: var(--color-primary-hover);
		text-decoration: underline;
	}

	.header-actions {
		display: flex;
		gap: var(--spacing-sm);
	}

	.details-container {
		background-color: var(--color-bg-primary);
		border-radius: var(--radius-lg);
		padding: var(--spacing-lg);
		border: 1px solid var(--color-border);
	}

	.details-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: var(--spacing-lg);
		padding-bottom: var(--spacing-md);
		border-bottom: 2px solid var(--color-border);
	}

	.details-header h2 {
		margin: 0;
		font-size: var(--font-size-2xl);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.detail-actions {
		display: flex;
		gap: var(--spacing-sm);
		align-items: center;
	}

	.detail-actions :global(.icon) {
		width: 1rem;
		height: 1rem;
		margin-right: var(--spacing-xs);
	}

	.details-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--spacing-lg);
		margin-top: var(--spacing-lg);
	}

	/* Desktop: Make details-grid single column in left panel */
	@media (min-width: 1024px) {
		.desktop-left .details-grid {
			grid-template-columns: 1fr;
		}
	}

	.chat-container {
		height: 600px;
		margin-top: var(--spacing-md);
	}

	/* Mobile/Tablet Layout */
	.tabs-container {
		display: block;
		margin-bottom: var(--spacing-lg);
	}

	.mobile-content {
		display: block;
	}

	.desktop-layout {
		display: none;
	}

	/* Desktop Layout */
	@media (min-width: 1024px) {
		.tabs-container {
			display: none;
		}

		.mobile-content {
			display: none;
		}

		.desktop-layout {
			display: flex;
			flex-direction: column;
			gap: var(--spacing-lg);
		}

		.desktop-main {
			display: grid;
			grid-template-columns: 2fr 1fr;
			gap: var(--spacing-lg);
			align-items: start;
			width: 100%;
			max-width: 100%;
		}

		.desktop-left {
			min-width: 0;
			width: 100%;
			overflow: hidden;
		}

		.desktop-right {
			position: sticky;
			top: var(--spacing-md);
			align-self: start;
			min-width: 0;
			width: 100%;
			max-width: 100%;
		}

		.chat-list-container {
			background-color: var(--color-bg-primary);
			border-radius: var(--radius-lg);
			border: 1px solid var(--color-border);
			padding: var(--spacing-md);
			width: 100%;
			max-width: 100%;
			box-sizing: border-box;
		}

		.chat-list-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: var(--spacing-md);
			padding-bottom: var(--spacing-md);
			border-bottom: 1px solid var(--color-border);
		}

		.chat-list-header h3 {
			margin: 0;
			font-size: var(--font-size-lg);
			font-weight: var(--font-weight-semibold);
			color: var(--color-text-primary);
		}

		.desktop-chat-container {
			background-color: var(--color-bg-primary);
			border-radius: var(--radius-lg);
			border: 1px solid var(--color-border);
			padding: var(--spacing-lg);
			min-height: 500px;
			width: 100%;
			max-width: 100%;
			box-sizing: border-box;
		}

		.desktop-chat-container :global(> *) {
			width: 100%;
			height: 100%;
		}
	}

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: var(--z-modal);
		padding: var(--spacing-md);
	}

	.modal {
		background-color: var(--color-bg-primary);
		border-radius: var(--radius-lg);
		padding: var(--spacing-lg);
		max-width: 600px;
		width: 100%;
		max-height: 90vh;
		overflow-y: auto;
		border: 1px solid var(--color-border);
	}

	.modal-large {
		max-width: 900px;
	}

	.modal h2 {
		margin: 0 0 var(--spacing-md) 0;
		font-size: var(--font-size-xl);
		font-weight: var(--font-weight-semibold);
		color: var(--color-text-primary);
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-md);
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--spacing-md);
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: var(--spacing-xs);
	}

	.form-group label {
		font-weight: var(--font-weight-medium);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
	}

	.form-group input,
	.form-group textarea,
	.form-group select {
		padding: var(--spacing-sm);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		font-size: var(--font-size-sm);
		color: var(--color-text-primary);
		background-color: var(--color-bg-primary);
		transition: border-color var(--transition-base);
	}

	.form-group input:focus,
	.form-group textarea:focus,
	.form-group select:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.form-group input:disabled,
	.form-group select:disabled {
		background-color: var(--color-bg-tertiary);
		cursor: not-allowed;
		opacity: 0.6;
	}

	.form-actions {
		display: flex;
		gap: var(--spacing-sm);
		margin-top: var(--spacing-sm);
		justify-content: flex-end;
	}
</style>
