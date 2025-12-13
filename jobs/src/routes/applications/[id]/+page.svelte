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
	import { listResumes, requestCVGeneration, type Resume } from '$lib/api/resumes';
	import { searchConversations, createConversation, type Conversation } from '$lib/api/chats';
	import InlineChat from '$lib/components/InlineChat.svelte';
	import { MessageSquare } from 'lucide-svelte';
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/Button.svelte';
	import LoadingState from '$lib/components/ui/LoadingState.svelte';
	import ErrorState from '$lib/components/ui/ErrorState.svelte';
import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
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
	let selectedConversationId: string | undefined = $state(undefined);
	let loadingConversations = $state(false);

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
	let requestingCV = $state(false);

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

	function formatStatus(status: string): string {
		return status.charAt(0).toUpperCase() + status.slice(1);
	}

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
			;
		} catch (err) {
			console.error('Error creating conversation:', err);
			alert('Failed to start chat');
		}
	}


	function handleConversationOpen(id: string) {
		selectedConversationId = id;
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

	async function handleRequestCV() {
		if (!application) return;
		if (!confirm('This will generate a new CV tailored to this job application. Continue?')) return;
		
		requestingCV = true;
		try {
			const generatedResume = await requestCVGeneration({
				jobApplicationId: application.id,
				language: application.language || 'en'
			});
			
			// Update the job application with the generated resume
			await updateJobApplication(application.id, {
				resumeId: generatedResume.id
			});
			
			// Refresh resumes list and application
			await Promise.all([fetchResumes(), loadApplication()]);
			alert('CV generated successfully and associated with this job application!');
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to generate CV');
			console.error('Error generating CV:', err);
		} finally {
			requestingCV = false;
		}
	}

	async function handleSelectCV() {
		if (!application || !formResumeId) return;
		try {
			await updateJobApplication(application.id, {
				resumeId: formResumeId
			});
			await loadApplication();
			showEditModal = false;
			alert('CV associated successfully!');
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to associate CV');
			console.error('Error associating CV:', err);
		}
	}
</script>

<div class="container mx-auto px-6 py-8 max-w-7xl">
	<div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
		<a href="/applications" class="text-blue-600 hover:text-blue-800 font-medium">← Back to Applications</a>
		<div class="flex gap-2">
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
		<!-- Application Header -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-6">
			<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-4">
				<div class="flex-1">
					<h1 class="text-3xl font-bold text-gray-900 mb-2">
						{application.companyName} - {application.jobTitle}
					</h1>
					<div class="flex flex-wrap items-center gap-4 text-sm text-gray-600">
						<div class="flex items-center gap-2">
							<span class="font-medium">Status:</span>
							<StatusBadge status={application.status} type="status">
								{formatStatus(application.status)}
							</StatusBadge>
						</div>
						{#if application.location}
							<div class="flex items-center gap-2">
								<span class="font-medium">Location:</span>
								<span>{application.location}</span>
							</div>
						{/if}
						{#if application.website}
							<div class="flex items-center gap-2">
								<span class="font-medium">Website:</span>
								<span class="capitalize">{application.website}</span>
							</div>
						{/if}
					</div>
				</div>
				<div class="flex flex-wrap gap-2">
					<Button onclick={handleStartChat} variant="primary" size="sm">
						<MessageSquare class="w-4 h-4 mr-2" />
						Start Chat
					</Button>
					<Button onclick={handleRequestCV} variant="primary" size="sm" disabled={requestingCV}>
						{requestingCV ? 'Generating CV...' : 'Request CV'}
					</Button>
					<Button onclick={openResponseModal} variant="secondary" size="sm">Add Response</Button>
					<Button onclick={openStageModal} variant="secondary" size="sm">Add Interview Stage</Button>
				</div>
			</div>
		</div>

		<!-- Application Details Grid -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
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

		<!-- Chat Section at Bottom -->
		<div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
			<div class="mb-4">
				<h2 class="text-xl font-semibold text-gray-900 mb-2">Chat</h2>
				{#if relatedConversations.length > 0}
					<div class="mb-4">
						<ConversationsList conversations={relatedConversations} onOpen={handleConversationOpen} />
					</div>
				{/if}
			</div>
			<div class="min-h-[400px]">
				<InlineChat
					
					jobApplicationId={application.id}
					conversationId={selectedConversationId}
					title={`${application.jobTitle} at ${application.companyName}`}
					description={application.jobDescription || `Chat about the ${application.jobTitle} position at ${application.companyName}`}
				/>
			</div>
		</div>
	{/if}
</div>
