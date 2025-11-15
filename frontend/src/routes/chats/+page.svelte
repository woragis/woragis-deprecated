<script lang="ts">
	import ChatFilters from './_components/ChatFilters.svelte';
	import ChatHeader from './_components/ChatHeader.svelte';
	import ConversationPanel from './_components/ConversationPanel.svelte';
	import ConversationsList from './_components/ConversationsList.svelte';
	import CreateConversationModal from './_components/CreateConversationModal.svelte';
	import { createChatsLogic } from './chats.logic';

	const {
		conversationsQuery,
		messagesQuery,
		transcriptsQuery,
		assignmentsQuery,
		selectedConversation,
		selectedConversationIds,
		searchInput,
		includeArchived,
		bulkStatus,
		composeContent,
		composeRole,
		generateReply,
		transcriptStatus,
		showCreateModal,
		createError,
		createForm,
		ideaQuery,
		projectQuery,
		ideaDropdownOpen,
		projectDropdownOpen,
		refreshConversations,
		applyFilters,
		setSearchInput,
		setIncludeArchived,
		selectConversation,
		toggleConversationSelection,
		clearSelections,
		handleSendMessage,
		handleShareTranscript,
		applyBulkAction,
		openCreateModal,
		closeCreateModal,
		handleCreateConversation,
		updateCreateFormField,
		selectProject,
		clearProjectSelection,
		selectIdea,
		clearIdeaSelection,
		getFilteredProjects,
		getFilteredIdeas,
		setIdeaQuery,
		setProjectQuery,
		setIdeaDropdownOpen,
		setProjectDropdownOpen,
		setComposeContent,
		setComposeRole,
		setGenerateReply,
		createConversationMutation,
		appendMessageMutation,
		shareTranscriptMutation
	} = createChatsLogic();

	$: conversationList = $conversationsQuery.data ?? [];
	$: messages = $messagesQuery.data ?? [];
	$: transcripts = $transcriptsQuery.data ?? [];
	$: assignments = $assignmentsQuery.data ?? [];
	$: selectionCount = $selectedConversationIds.size;
	$: projectSuggestions = getFilteredProjects();
	$: ideaSuggestions = getFilteredIdeas();
</script>

<svelte:head>
	<title>Chats · Woragis</title>
</svelte:head>

<div class="flex flex-col gap-6">
	<ChatHeader
		isRefreshing={$conversationsQuery.isFetching}
		onRefresh={refreshConversations}
		onCreate={openCreateModal}
	/>

	<ChatFilters
		searchValue={$searchInput}
		includeArchived={$includeArchived}
		bulkStatus={$bulkStatus}
		selectionCount={selectionCount}
		onSearchInput={setSearchInput}
		onIncludeArchivedChange={setIncludeArchived}
		onSearch={applyFilters}
		onBulkArchive={() => applyBulkAction('archive')}
		onBulkRestore={() => applyBulkAction('restore')}
		onBulkDelete={() => applyBulkAction('delete')}
		onClearSelection={clearSelections}
	/>

	<div class="grid gap-6 lg:grid-cols-[320px_1fr]">
		<ConversationsList
			conversations={conversationList}
			isLoading={$conversationsQuery.isPending}
			isFetching={$conversationsQuery.isFetching}
			selectedConversationId={$selectedConversation?.id ?? null}
			selectedIds={$selectedConversationIds}
			onSelect={selectConversation}
			onToggleSelection={toggleConversationSelection}
		/>

		<ConversationPanel
			conversation={$selectedConversation}
			messages={messages}
			messagesLoading={$messagesQuery.isPending || $messagesQuery.isFetching}
			transcripts={transcripts}
			transcriptsLoading={$transcriptsQuery.isPending || $transcriptsQuery.isFetching}
			assignments={assignments}
			assignmentsLoading={$assignmentsQuery.isPending || $assignmentsQuery.isFetching}
			transcriptStatus={$transcriptStatus}
			composeContent={$composeContent}
			composeRole={$composeRole}
			generateReply={$generateReply}
			onComposeContentChange={setComposeContent}
			onComposeRoleChange={setComposeRole}
			onGenerateReplyChange={setGenerateReply}
			onSendMessage={handleSendMessage}
			onShareTranscript={handleShareTranscript}
			isSending={$appendMessageMutation.isPending}
			isSharing={$shareTranscriptMutation.isPending}
		/>
	</div>
</div>

<CreateConversationModal
	open={$showCreateModal}
	form={$createForm}
	createError={$createError}
	ideaQueryValue={$ideaQuery}
	projectQueryValue={$projectQuery}
	ideaDropdownOpen={$ideaDropdownOpen}
	projectDropdownOpen={$projectDropdownOpen}
	filteredProjects={projectSuggestions}
	filteredIdeas={ideaSuggestions}
	isSubmitting={$createConversationMutation.isPending}
	onClose={closeCreateModal}
	onSubmit={handleCreateConversation}
	onFieldChange={updateCreateFormField}
	onIdeaQueryChange={setIdeaQuery}
	onProjectQueryChange={setProjectQuery}
	onIdeaDropdownChange={setIdeaDropdownOpen}
	onProjectDropdownChange={setProjectDropdownOpen}
	onSelectIdea={selectIdea}
	onClearIdea={clearIdeaSelection}
	onSelectProject={selectProject}
	onClearProject={clearProjectSelection}
/>

