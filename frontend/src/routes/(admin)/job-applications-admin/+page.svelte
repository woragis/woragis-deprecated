<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listJobApplications,
		createJobApplication,
		updateJobApplication,
		updateJobApplicationStatus,
		deleteJobApplication,
		type JobApplication,
		type CreateJobApplicationInput,
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
		scheduleStage,
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

	let applications: JobApplication[] = $state([]);
	let resumes: Resume[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state('');
	let activeTab = $state<'applications' | 'details' | 'chat'>('applications');
	let selectedApplication: JobApplication | null = $state(null);
	let applicationResponses: JobApplicationResponse[] = $state([]);
	let applicationStages: InterviewStage[] = $state([]);
	let relatedConversations: Conversation[] = $state([]);
	let showChat = $state(false);
	let loadingConversations = $state(false);

	// Modal states
	let showCreateModal = $state(false);
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
		await Promise.all([fetchApplications(), fetchResumes()]);
	});

	async function fetchApplications() {
		loading = true;
		error = null;
		try {
			applications = await listJobApplications();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load job applications';
			console.error('Error fetching job applications:', err);
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

	async function fetchApplicationDetails(app: JobApplication) {
		selectedApplication = app;
		activeTab = 'details';
		try {
			[applicationResponses, applicationStages] = await Promise.all([
				listResponses(app.id),
				listStages(app.id)
			]);
			await fetchRelatedConversations(app.id);
		} catch (err) {
			console.error('Error fetching application details:', err);
		}
	}

	async function fetchRelatedConversations(applicationId: string) {
		loadingConversations = true;
		try {
			relatedConversations = await searchConversations(undefined, false, applicationId);
		} catch (err) {
			console.error('Error fetching related conversations:', err);
		} finally {
			loadingConversations = false;
		}
	}

	async function handleStartChat() {
		if (!selectedApplication) return;

		try {
			const conversation = await createConversation({
				title: `Job Application: ${selectedApplication.jobTitle} at ${selectedApplication.companyName}`,
				description: `Chat about the ${selectedApplication.jobTitle} position at ${selectedApplication.companyName}`,
				jobApplicationId: selectedApplication.id
			});
			relatedConversations = [conversation, ...relatedConversations];
			showChat = true;
			activeTab = 'chat';
		} catch (err) {
			console.error('Error creating conversation:', err);
			alert('Failed to start chat');
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(app: JobApplication) {
		selectedApplication = app;
		formCompanyName = app.companyName;
		formLocation = app.location || '';
		formJobTitle = app.jobTitle;
		formJobUrl = app.jobUrl;
		formWebsite = app.website;
		formCoverLetter = app.coverLetter || '';
		formLinkedInContact = app.linkedInContact;
		formStatus = app.status;
		formResumeId = app.resumeId || '';
		formSalaryMin = app.salaryMin || '';
		formSalaryMax = app.salaryMax || '';
		formSalaryCurrency = app.salaryCurrency || '';
		formJobDescription = app.jobDescription || '';
		formDeadline = app.deadline ? app.deadline.split('T')[0] : '';
		formInterestLevel = app.interestLevel || '';
		formNotes = app.notes || '';
		formTags = app.tags?.join(', ') || '';
		formFollowUpDate = app.followUpDate ? app.followUpDate.split('T')[0] : '';
		formSource = app.source || '';
		formApplicationMethod = app.applicationMethod || '';
		formLanguage = app.language || '';
		showEditModal = true;
	}

	function openResponseModal(app?: JobApplication) {
		resetResponseForm();
		if (app) {
			selectedApplication = app;
		}
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

	function openStageModal(app?: JobApplication) {
		resetStageForm();
		if (app) {
			selectedApplication = app;
		}
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

	function resetForm() {
		formCompanyName = '';
		formLocation = '';
		formJobTitle = '';
		formJobUrl = '';
		formWebsite = '';
		formCoverLetter = '';
		formLinkedInContact = false;
		formStatus = 'pending';
		formResumeId = '';
		formSalaryMin = '';
		formSalaryMax = '';
		formSalaryCurrency = '';
		formJobDescription = '';
		formDeadline = '';
		formInterestLevel = '';
		formNotes = '';
		formTags = '';
		formFollowUpDate = '';
		formSource = '';
		formApplicationMethod = '';
		formLanguage = '';
		selectedApplication = null;
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

	async function handleCreate() {
		if (!formCompanyName.trim() || !formJobTitle.trim() || !formJobUrl.trim() || !formWebsite.trim()) {
			alert('Company name, job title, job URL, and website are required');
			return;
		}

		try {
			const input: CreateJobApplicationInput = {
				companyName: formCompanyName.trim(),
				location: formLocation.trim() || undefined,
				jobTitle: formJobTitle.trim(),
				jobUrl: formJobUrl.trim(),
				website: formWebsite.trim(),
				coverLetter: formCoverLetter.trim() || undefined,
				linkedInContact: formLinkedInContact,
				status: formStatus
			};

			await createJobApplication(input);
			showCreateModal = false;
			resetForm();
			await fetchApplications();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create job application');
			console.error('Error creating job application:', err);
		}
	}

	async function handleUpdate() {
		if (!selectedApplication) return;

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

			await updateJobApplication(selectedApplication.id, input);
			showEditModal = false;
			resetForm();
			await fetchApplications();
			if (selectedApplication) {
				await fetchApplicationDetails(selectedApplication);
			}
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update job application');
			console.error('Error updating job application:', err);
		}
	}

	async function handleUpdateStatus() {
		if (!selectedApplication) return;

		try {
			await updateJobApplicationStatus(selectedApplication.id, { status: formStatus });
			showEditModal = false;
			resetForm();
			await fetchApplications();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update job application status');
			console.error('Error updating job application status:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this job application?')) return;

		try {
			await deleteJobApplication(id);
			await fetchApplications();
			if (selectedApplication?.id === id) {
				selectedApplication = null;
				activeTab = 'applications';
			}
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete job application');
			console.error('Error deleting job application:', err);
		}
	}

	async function handleCreateResponse() {
		if (!selectedApplication) {
			alert('Please select an application first');
			return;
		}

		try {
			await createResponse(selectedApplication.id, {
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
			await fetchApplicationDetails(selectedApplication);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create response');
			console.error('Error creating response:', err);
		}
	}

	async function handleUpdateResponse() {
		if (!editingResponse || !selectedApplication) return;

		try {
			await updateResponse(editingResponse.id, selectedApplication.id, {
				message: formResponseMessage || undefined,
				contactPerson: formContactPerson || undefined,
				contactEmail: formContactEmail || undefined,
				contactPhone: formContactPhone || undefined,
				responseChannel: formResponseChannel || undefined
			});
			showResponseModal = false;
			resetResponseForm();
			await fetchApplicationDetails(selectedApplication);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update response');
			console.error('Error updating response:', err);
		}
	}

	async function handleDeleteResponse(id: string) {
		if (!selectedApplication || !confirm('Are you sure you want to delete this response?')) return;

		try {
			await deleteResponse(id, selectedApplication.id);
			await fetchApplicationDetails(selectedApplication);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete response');
			console.error('Error deleting response:', err);
		}
	}

	async function handleCreateStage() {
		if (!selectedApplication) {
			alert('Please select an application first');
			return;
		}

		try {
			await createStage(selectedApplication.id, {
				stageType: formStageType,
				scheduledDate: formScheduledDate ? new Date(formScheduledDate).toISOString() : undefined,
				interviewerName: formInterviewerName || undefined,
				interviewerEmail: formInterviewerEmail || undefined,
				location: formStageLocation || undefined,
				notes: formStageNotes || undefined
			});
			showStageModal = false;
			resetStageForm();
			await fetchApplicationDetails(selectedApplication);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create interview stage');
			console.error('Error creating interview stage:', err);
		}
	}

	async function handleUpdateStage() {
		if (!editingStage || !selectedApplication) return;

		try {
			await updateStage(editingStage.id, selectedApplication.id, {
				scheduledDate: formScheduledDate ? new Date(formScheduledDate).toISOString() : undefined,
				interviewerName: formInterviewerName || undefined,
				interviewerEmail: formInterviewerEmail || undefined,
				location: formStageLocation || undefined,
				notes: formStageNotes || undefined,
				feedback: formStageFeedback || undefined
			});
			showStageModal = false;
			resetStageForm();
			await fetchApplicationDetails(selectedApplication);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update interview stage');
			console.error('Error updating interview stage:', err);
		}
	}

	async function handleCompleteStage(stage: InterviewStage) {
		if (!selectedApplication) return;
		const outcome = prompt('Enter outcome (passed, failed, cancelled):') as StageOutcome | null;
		if (!outcome || !['passed', 'failed', 'cancelled'].includes(outcome)) return;

		try {
			await completeStage(stage.id, selectedApplication.id, {
				completedDate: new Date().toISOString(),
				outcome
			});
			await fetchApplicationDetails(selectedApplication);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to complete interview stage');
			console.error('Error completing interview stage:', err);
		}
	}

	async function handleDeleteStage(id: string) {
		if (!selectedApplication || !confirm('Are you sure you want to delete this interview stage?'))
			return;

		try {
			await deleteStage(id, selectedApplication.id);
			await fetchApplicationDetails(selectedApplication);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete interview stage');
			console.error('Error deleting interview stage:', err);
		}
	}

	function filteredApplications() {
		if (!searchQuery.trim()) return applications;
		const query = searchQuery.toLowerCase();
		return applications.filter(
			(a) =>
				a.companyName.toLowerCase().includes(query) ||
				a.jobTitle.toLowerCase().includes(query) ||
				a.location?.toLowerCase().includes(query)
		);
	}

	function formatDate(dateString?: string): string {
		if (!dateString) return '—';
		return new Date(dateString).toLocaleDateString();
	}
</script>

<div class="page-container">
	<div class="header">
		<div>
			<h1>Job Applications Management</h1>
			<p>Manage job applications, responses, and interview stages</p>
		</div>
		<button onclick={openCreateModal}>Create Application</button>
	</div>

	<div class="tabs">
		<button
			class="tab {activeTab === 'applications' ? 'active' : ''}"
			onclick={() => (activeTab = 'applications')}
		>
			Applications
		</button>
		{#if selectedApplication}
			<button
				class="tab {activeTab === 'details' ? 'active' : ''}"
				onclick={() => (activeTab = 'details')}
			>
				Details: {selectedApplication.companyName}
			</button>
			<button
				class="tab {activeTab === 'chat' ? 'active' : ''}"
				onclick={() => {
					activeTab = 'chat';
					if (selectedApplication) {
						fetchRelatedConversations(selectedApplication.id);
					}
				}}
			>
				<MessageSquare class="w-4 h-4 inline mr-1" />
				Chat
			</button>
		{/if}
	</div>

	{#if activeTab === 'applications'}
		<div class="search-bar">
			<input
				type="text"
				placeholder="Search applications..."
				bind:value={searchQuery}
				class="search-input"
			/>
		</div>

		{#if error}
			<div class="error">{error}</div>
		{/if}

		{#if loading}
			<div class="loading">Loading...</div>
		{:else if filteredApplications().length === 0}
			<div class="empty">No job applications found</div>
		{:else}
			<table class="table">
				<thead>
					<tr>
						<th>Company</th>
						<th>Job Title</th>
						<th>Location</th>
						<th>Website</th>
						<th>Language</th>
						<th>Status</th>
						<th>Interest</th>
						<th>Applied At</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each filteredApplications() as app}
						<tr>
							<td><strong>{app.companyName}</strong></td>
							<td>{app.jobTitle}</td>
							<td>{app.location || '—'}</td>
							<td>{app.website}</td>
							<td>{app.language || '—'}</td>
							<td>
								<span class="status status-{app.status}">{app.status}</span>
							</td>
							<td>{app.interestLevel || '—'}</td>
							<td>{formatDate(app.appliedAt)}</td>
							<td>
								<button onclick={() => fetchApplicationDetails(app)}>View</button>
								<button onclick={() => openEditModal(app)}>Edit</button>
								<button onclick={() => handleDelete(app.id)} class="delete-btn">Delete</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	{:else if activeTab === 'details' && selectedApplication}
		<div class="details-container">
			<div class="details-header">
				<h2>{selectedApplication.companyName} - {selectedApplication.jobTitle}</h2>
				<div class="detail-actions">
					<button onclick={handleStartChat} class="flex items-center gap-2">
						<MessageSquare class="w-4 h-4" />
						Start Chat
					</button>
					<button onclick={() => openResponseModal()}>Add Response</button>
					<button onclick={() => openStageModal()}>Add Interview Stage</button>
					<button onclick={() => selectedApplication && openEditModal(selectedApplication)}>Edit Application</button>
				</div>
			</div>

			{#if relatedConversations.length > 0}
				<div class="detail-section mt-4">
					<h3>Related Conversations ({relatedConversations.length})</h3>
					<div class="list">
						{#each relatedConversations as conv}
							<div class="list-item">
								<div class="list-item-header">
									<span class="font-semibold">{conv.title}</span>
									<span class="date">{formatDate(conv.updatedAt)}</span>
								</div>
								{#if conv.description}
									<p class="list-item-text">{conv.description}</p>
								{/if}
								<div class="list-item-actions">
									<button onclick={() => {
										activeTab = 'chat';
										showChat = true;
									}}>Open Chat</button>
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/if}

			<div class="details-grid">
				<div class="detail-section">
					<h3>Application Info</h3>
					<div class="detail-item">
						<strong>Status:</strong> <span class="status status-{selectedApplication.status}">{selectedApplication.status}</span>
					</div>
					<div class="detail-item"><strong>Location:</strong> {selectedApplication.location || '—'}</div>
					<div class="detail-item"><strong>Website:</strong> {selectedApplication.website}</div>
					<div class="detail-item"><strong>Interest Level:</strong> {selectedApplication.interestLevel || '—'}</div>
					<div class="detail-item"><strong>Source:</strong> {selectedApplication.source || '—'}</div>
					<div class="detail-item"><strong>Language:</strong> {selectedApplication.language || '—'}</div>
					<div class="detail-item"><strong>Applied At:</strong> {formatDate(selectedApplication.appliedAt)}</div>
					{#if selectedApplication.salaryMin || selectedApplication.salaryMax}
						<div class="detail-item">
							<strong>Salary:</strong> {selectedApplication.salaryMin || '—'} - {selectedApplication.salaryMax || '—'} {selectedApplication.salaryCurrency || ''}
						</div>
					{/if}
					{#if selectedApplication.deadline}
						<div class="detail-item"><strong>Deadline:</strong> {formatDate(selectedApplication.deadline)}</div>
					{/if}
					{#if selectedApplication.followUpDate}
						<div class="detail-item"><strong>Follow-up Date:</strong> {formatDate(selectedApplication.followUpDate)}</div>
					{/if}
					{#if selectedApplication.notes}
						<div class="detail-item"><strong>Notes:</strong> {selectedApplication.notes}</div>
					{/if}
					{#if selectedApplication.tags && selectedApplication.tags.length > 0}
						<div class="detail-item">
							<strong>Tags:</strong> {selectedApplication.tags.join(', ')}
						</div>
					{/if}
				</div>

				<div class="detail-section">
					<h3>Responses ({applicationResponses.length})</h3>
					{#if applicationResponses.length === 0}
						<p class="empty-text">No responses yet</p>
					{:else}
						<div class="list">
							{#each applicationResponses as response}
								<div class="list-item">
									<div class="list-item-header">
										<span class="response-type response-type-{response.responseType}">{response.responseType}</span>
										<span class="date">{formatDate(response.responseDate)}</span>
									</div>
									{#if response.message}
										<p class="list-item-text">{response.message}</p>
									{/if}
									{#if response.contactPerson}
										<p class="list-item-text"><strong>Contact:</strong> {response.contactPerson} {response.contactEmail ? `(${response.contactEmail})` : ''}</p>
									{/if}
									<div class="list-item-actions">
										<button onclick={() => openEditResponseModal(response)}>Edit</button>
										<button onclick={() => handleDeleteResponse(response.id)} class="delete-btn">Delete</button>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>

				<div class="detail-section">
					<h3>Interview Stages ({applicationStages.length})</h3>
					{#if applicationStages.length === 0}
						<p class="empty-text">No interview stages yet</p>
					{:else}
						<div class="list">
							{#each applicationStages as stage}
								<div class="list-item">
									<div class="list-item-header">
										<span class="stage-type">{stage.stageType}</span>
										<span class="outcome outcome-{stage.outcome}">{stage.outcome}</span>
									</div>
									{#if stage.scheduledDate}
										<p class="list-item-text"><strong>Scheduled:</strong> {formatDate(stage.scheduledDate)}</p>
									{/if}
									{#if stage.interviewerName}
										<p class="list-item-text"><strong>Interviewer:</strong> {stage.interviewerName} {stage.interviewerEmail ? `(${stage.interviewerEmail})` : ''}</p>
									{/if}
									{#if stage.notes}
										<p class="list-item-text">{stage.notes}</p>
									{/if}
									<div class="list-item-actions">
										<button onclick={() => openEditStageModal(stage)}>Edit</button>
										{#if stage.outcome === 'pending'}
											<button onclick={() => handleCompleteStage(stage)}>Complete</button>
										{/if}
										<button onclick={() => handleDeleteStage(stage.id)} class="delete-btn">Delete</button>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>
	{:else if activeTab === 'chat' && selectedApplication}
		<div class="details-container">
			<div class="details-header">
				<h2>Chat: {selectedApplication.companyName} - {selectedApplication.jobTitle}</h2>
				<div class="detail-actions">
					{#if relatedConversations.length === 0}
						<button onclick={handleStartChat} class="flex items-center gap-2">
							<MessageSquare class="w-4 h-4" />
							Start Chat
						</button>
					{/if}
				</div>
			</div>

			<div class="mt-4" style="height: 600px;">
				<InlineChat
					jobApplicationId={selectedApplication.id}
					title={`${selectedApplication.jobTitle} at ${selectedApplication.companyName}`}
					description={selectedApplication.jobDescription || `Chat about the ${selectedApplication.jobTitle} position at ${selectedApplication.companyName}`}
				/>
			</div>
		</div>
	{/if}
</div>

<!-- Create/Edit Application Modal -->
{#if showCreateModal || showEditModal}
	<div class="modal-overlay" onclick={() => {
		showCreateModal = false;
		showEditModal = false;
		resetForm();
	}}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>{showCreateModal ? 'Create' : 'Edit'} Job Application</h2>
			<div class="form">
				<div class="form-row">
					<div class="form-group">
						<label>Company Name *</label>
						<input type="text" bind:value={formCompanyName} disabled={showEditModal} />
					</div>
					<div class="form-group">
						<label>Job Title *</label>
						<input type="text" bind:value={formJobTitle} disabled={showEditModal} />
					</div>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Job URL *</label>
						<input type="url" bind:value={formJobUrl} disabled={showEditModal} />
					</div>
					<div class="form-group">
						<label>Website *</label>
						<input type="text" bind:value={formWebsite} placeholder="linkedin, glassdoor, etc." disabled={showEditModal} />
					</div>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Location</label>
						<input type="text" bind:value={formLocation} disabled={showEditModal} />
					</div>
					<div class="form-group">
						<label>Status</label>
						<select bind:value={formStatus}>
							{#each statuses as status}
								<option value={status}>{status}</option>
							{/each}
						</select>
					</div>
				</div>
				{#if showEditModal}
					<div class="form-row">
						<div class="form-group">
							<label>Resume</label>
							<select bind:value={formResumeId}>
								<option value="">None</option>
								{#each resumes as resume}
									<option value={resume.id}>{resume.title}</option>
								{/each}
							</select>
						</div>
						<div class="form-group">
							<label>Interest Level</label>
							<select bind:value={formInterestLevel}>
								<option value="">None</option>
								{#each interestLevels as level}
									<option value={level}>{level}</option>
								{/each}
							</select>
						</div>
					</div>
					<div class="form-row">
						<div class="form-group">
							<label>Salary Min</label>
							<input type="number" bind:value={formSalaryMin} />
						</div>
						<div class="form-group">
							<label>Salary Max</label>
							<input type="number" bind:value={formSalaryMax} />
						</div>
						<div class="form-group">
							<label>Currency</label>
							<input type="text" bind:value={formSalaryCurrency} placeholder="USD, EUR, etc." />
						</div>
					</div>
					<div class="form-row">
						<div class="form-group">
							<label>Deadline</label>
							<input type="date" bind:value={formDeadline} />
						</div>
						<div class="form-group">
							<label>Follow-up Date</label>
							<input type="date" bind:value={formFollowUpDate} />
						</div>
					</div>
					<div class="form-row">
						<div class="form-group">
							<label>Source</label>
							<input type="text" bind:value={formSource} placeholder="referral, job-board, etc." />
						</div>
						<div class="form-group">
							<label>Application Method</label>
							<input type="text" bind:value={formApplicationMethod} placeholder="auto, manual, assisted" />
						</div>
					</div>
					<div class="form-group">
						<label>Language (ISO 639-1)</label>
						<input type="text" bind:value={formLanguage} placeholder="en, pt, es" maxlength="2" pattern="[a-z]{2}" />
						<small style="color: #555; font-size: 0.75rem;">2-character language code (e.g., "en", "pt", "es")</small>
					</div>
					<div class="form-group">
						<label>Tags (comma-separated)</label>
						<input type="text" bind:value={formTags} placeholder="remote, startup, dream-job" />
					</div>
					<div class="form-group">
						<label>Job Description</label>
						<textarea bind:value={formJobDescription} rows="4"></textarea>
					</div>
					<div class="form-group">
						<label>Notes</label>
						<textarea bind:value={formNotes} rows="3"></textarea>
					</div>
				{/if}
				<div class="form-group">
					<label>Cover Letter</label>
					<textarea bind:value={formCoverLetter} rows="6"></textarea>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formLinkedInContact} />
						LinkedIn Contact
					</label>
				</div>
				<div class="form-actions">
					{#if showCreateModal}
						<button onclick={handleCreate}>Create</button>
					{:else}
						<button onclick={handleUpdate}>Update</button>
						<button onclick={handleUpdateStatus}>Update Status Only</button>
					{/if}
					<button onclick={() => {
						showCreateModal = false;
						showEditModal = false;
						resetForm();
					}}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Response Modal -->
{#if showResponseModal && selectedApplication}
	<div class="modal-overlay" onclick={() => {
		showResponseModal = false;
		resetResponseForm();
	}}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>{editingResponse ? 'Edit' : 'Add'} Response</h2>
			<div class="form">
				<div class="form-group">
					<label>Response Type *</label>
					<select bind:value={formResponseType} disabled={!!editingResponse}>
						{#each responseTypes as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Response Date</label>
					<input type="date" bind:value={formResponseDate} />
				</div>
				<div class="form-group">
					<label>Message</label>
					<textarea bind:value={formResponseMessage} rows="4"></textarea>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Contact Person</label>
						<input type="text" bind:value={formContactPerson} />
					</div>
					<div class="form-group">
						<label>Contact Email</label>
						<input type="email" bind:value={formContactEmail} />
					</div>
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Contact Phone</label>
						<input type="tel" bind:value={formContactPhone} />
					</div>
					<div class="form-group">
						<label>Response Channel</label>
						<input type="text" bind:value={formResponseChannel} placeholder="email, phone, linkedin" />
					</div>
				</div>
				<div class="form-actions">
					<button onclick={editingResponse ? handleUpdateResponse : handleCreateResponse}>
						{editingResponse ? 'Update' : 'Create'}
					</button>
					<button onclick={() => {
						showResponseModal = false;
						resetResponseForm();
					}}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Stage Modal -->
{#if showStageModal && selectedApplication}
	<div class="modal-overlay" onclick={() => {
		showStageModal = false;
		resetStageForm();
	}}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>{editingStage ? 'Edit' : 'Add'} Interview Stage</h2>
			<div class="form">
				<div class="form-group">
					<label>Stage Type *</label>
					<select bind:value={formStageType} disabled={!!editingStage}>
						{#each stageTypes as type}
							<option value={type}>{type}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Scheduled Date</label>
					<input type="date" bind:value={formScheduledDate} />
				</div>
				<div class="form-row">
					<div class="form-group">
						<label>Interviewer Name</label>
						<input type="text" bind:value={formInterviewerName} />
					</div>
					<div class="form-group">
						<label>Interviewer Email</label>
						<input type="email" bind:value={formInterviewerEmail} />
					</div>
				</div>
				<div class="form-group">
					<label>Location</label>
					<input type="text" bind:value={formStageLocation} placeholder="in-person, video, phone, or address" />
				</div>
				<div class="form-group">
					<label>Notes</label>
					<textarea bind:value={formStageNotes} rows="3"></textarea>
				</div>
				{#if editingStage}
					<div class="form-group">
						<label>Feedback</label>
						<textarea bind:value={formStageFeedback} rows="3"></textarea>
					</div>
				{/if}
				<div class="form-actions">
					<button onclick={editingStage ? handleUpdateStage : handleCreateStage}>
						{editingStage ? 'Update' : 'Create'}
					</button>
					<button onclick={() => {
						showStageModal = false;
						resetStageForm();
					}}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-container {
		padding: 1rem;
		max-width: 1400px;
		margin: 0 auto;
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
	}

	.header h1 {
		margin: 0 0 0.25rem 0;
		font-size: 1.5rem;
		color: #333;
	}

	.header p {
		margin: 0;
		color: #666;
		font-size: 0.9rem;
	}

	.header button {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
	}

	.header button:hover {
		background: #0056b3;
	}

	.tabs {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 1rem;
		border-bottom: 2px solid #e0e0e0;
	}

	.tab {
		padding: 0.75rem 1.5rem;
		background: none;
		border: none;
		border-bottom: 2px solid transparent;
		cursor: pointer;
		font-size: 0.9rem;
		color: #666;
		margin-bottom: -2px;
	}

	.tab:hover {
		color: #007bff;
	}

	.tab.active {
		color: #007bff;
		border-bottom-color: #007bff;
		font-weight: 600;
	}

	.search-bar {
		margin-bottom: 1rem;
	}

	.search-input {
		width: 100%;
		max-width: 400px;
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
	}

	.error {
		padding: 0.75rem;
		background: #fee;
		color: #c33;
		border: 1px solid #fcc;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.loading,
	.empty {
		padding: 2rem;
		text-align: center;
		color: #666;
	}

	.table {
		width: 100%;
		border-collapse: collapse;
		background: white;
	}

	.table th,
	.table td {
		padding: 0.75rem;
		text-align: left;
		border-bottom: 1px solid #ddd;
		color: #333;
	}

	.table th {
		background: #f5f5f5;
		font-weight: 600;
	}

	.table tbody tr:hover {
		background: #f9f9f9;
	}

	.table button {
		padding: 0.25rem 0.75rem;
		background: #28a745;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
		margin-right: 0.5rem;
	}

	.table button:hover {
		background: #218838;
	}

	.delete-btn {
		background: #dc3545 !important;
	}

	.delete-btn:hover {
		background: #c82333 !important;
	}

	.status {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.status-pending {
		background: #fff3cd;
		color: #856404;
	}

	.status-processing {
		background: #d1ecf1;
		color: #0c5460;
	}

	.status-applied {
		background: #d4edda;
		color: #155724;
	}

	.status-contacted {
		background: #d1ecf1;
		color: #0c5460;
	}

	.status-rejected {
		background: #f8d7da;
		color: #721c24;
	}

	.status-accepted {
		background: #d4edda;
		color: #155724;
	}

	.status-failed {
		background: #f8d7da;
		color: #721c24;
	}

	.details-container {
		background: white;
		border-radius: 8px;
		padding: 1.5rem;
		color: #333;
	}

	.details-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 2px solid #e0e0e0;
	}

	.details-header h2 {
		margin: 0;
		font-size: 1.5rem;
		color: #333;
	}

	.detail-actions {
		display: flex;
		gap: 0.5rem;
	}

	.detail-actions button {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.detail-actions button:hover {
		background: #0056b3;
	}

	.details-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1.5rem;
	}

	.detail-section {
		background: #f9f9f9;
		padding: 1rem;
		border-radius: 4px;
	}

	.detail-section h3 {
		margin: 0 0 1rem 0;
		font-size: 1.1rem;
		color: #333;
	}

	.detail-item {
		margin-bottom: 0.5rem;
		font-size: 0.9rem;
		color: #333;
	}

	.detail-item strong {
		color: #333;
	}

	.empty-text {
		color: #666;
		font-style: italic;
	}

	.list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.list-item {
		background: white;
		padding: 0.75rem;
		border-radius: 4px;
		border: 1px solid #ddd;
		color: #333;
	}

	.list-item-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.5rem;
		color: #333;
	}

	.list-item-header span {
		color: #333;
	}

	.list-item-text {
		margin: 0.25rem 0;
		font-size: 0.875rem;
		color: #555;
	}

	.list-item-text strong {
		color: #333;
	}

	.list-item * {
		color: inherit;
	}

	.list-item-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}

	.list-item-actions button {
		padding: 0.25rem 0.5rem;
		background: #28a745;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.75rem;
	}

	.response-type,
	.stage-type {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
		color: #333;
		background: #e9ecef;
	}

	.response-type-rejection {
		background: #f8d7da;
		color: #721c24;
	}

	.response-type-interview {
		background: #d1ecf1;
		color: #0c5460;
	}

	.response-type-offer {
		background: #d4edda;
		color: #155724;
	}

	.response-type-no-response {
		background: #e2e3e5;
		color: #383d41;
	}

	.outcome {
		padding: 0.25rem 0.5rem;
		border-radius: 4px;
		font-size: 0.75rem;
		font-weight: 500;
	}

	.outcome-pending {
		background: #fff3cd;
		color: #856404;
	}

	.outcome-passed {
		background: #d4edda;
		color: #155724;
	}

	.outcome-failed {
		background: #f8d7da;
		color: #721c24;
	}

	.outcome-cancelled {
		background: #e2e3e5;
		color: #383d41;
	}

	.date {
		font-size: 0.75rem;
		color: #666;
	}

	.modal-overlay {
		position: fixed;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal {
		background: white;
		border-radius: 8px;
		padding: 1.5rem;
		max-width: 600px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
		color: #333;
	}

	.modal-large {
		max-width: 900px;
	}

	.modal h2 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
		color: #333;
	}

	.form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.form-row {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 1rem;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.form-group label {
		font-weight: 500;
		font-size: 0.875rem;
		color: #333;
	}

	.form-group input,
	.form-group textarea,
	.form-group select {
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
		color: #333;
		background: white;
	}

	.form-group input:disabled,
	.form-group select:disabled {
		background: #f5f5f5;
		cursor: not-allowed;
	}

	.form-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 0.5rem;
	}

	.form-actions button {
		padding: 0.5rem 1rem;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.form-actions button:first-child {
		background: #007bff;
		color: white;
	}

	.form-actions button:first-child:hover {
		background: #0056b3;
	}

	.form-actions button:last-child {
		background: #6c757d;
		color: white;
	}

	.form-actions button:last-child:hover {
		background: #5a6268;
	}
</style>
