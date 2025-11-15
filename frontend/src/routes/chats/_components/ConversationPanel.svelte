<script lang="ts">
	import type {
		ChatAssignment,
		ChatConversation,
		ChatMessage,
		ChatTranscript
	} from '$lib/api/types';

	export let conversation: ChatConversation | null = null;
	export let messages: ChatMessage[] = [];
	export let messagesLoading = false;
	export let transcripts: ChatTranscript[] = [];
	export let transcriptsLoading = false;
	export let assignments: ChatAssignment[] = [];
	export let assignmentsLoading = false;
	export let transcriptStatus = '';
	export let composeContent = '';
	export let composeRole: 'user' | 'assistant' = 'user';
	export let generateReply = false;
	export let onComposeContentChange: (value: string) => void;
	export let onComposeRoleChange: (value: 'user' | 'assistant') => void;
	export let onGenerateReplyChange: (value: boolean) => void;
	export let onSendMessage: (event: SubmitEvent) => void;
	export let onShareTranscript: () => void;
	export let isSending = false;
	export let isSharing = false;

	const formatTimestamp = (value: string) => new Date(value).toLocaleString();
</script>

<section class="flex h-[650px] flex-col gap-4 rounded-2xl border border-slate-800/80 bg-slate-950/60 p-5">
	{#if !conversation}
		<div class="flex flex-1 flex-col items-center justify-center gap-3 text-center">
			<h2 class="text-lg font-semibold text-slate-200">Select a conversation</h2>
			<p class="max-w-sm text-sm text-slate-500">
				Choose a conversation from the list to review messages, stream new activity, and manage transcripts.
			</p>
		</div>
	{:else}
		<header class="flex flex-wrap items-center justify-between gap-3">
			<div>
				<h2 class="text-lg font-semibold text-slate-100">{conversation.title}</h2>
				<p class="text-xs text-slate-400">
					Last updated {formatTimestamp(conversation.updated_at)}
				</p>
			</div>
			<button
				class="rounded-lg border border-slate-700 px-3 py-2 text-xs text-slate-300 transition hover:border-slate-500 hover:text-slate-100 disabled:opacity-60"
				type="button"
				on:click={onShareTranscript}
				disabled={isSharing}
			>
				{isSharing ? 'Sharing…' : 'Share transcript'}
			</button>
		</header>

		{#if transcriptStatus}
			<div class="rounded-lg border border-slate-700/60 bg-slate-900/70 px-3 py-2 text-xs text-slate-300">
				{transcriptStatus}
			</div>
		{/if}

		<div class="grid flex-1 grid-cols-1 gap-4 lg:grid-cols-[2fr_1fr]">
			<div class="flex h-full flex-col overflow-hidden rounded-xl border border-slate-800/80 bg-slate-950/70">
				<div class="flex items-center justify-between border-b border-slate-800/80 px-4 py-3 text-xs text-slate-400">
					<span>Messages</span>
					{#if messagesLoading}
						<span>Loading…</span>
					{/if}
				</div>
				<div class="flex-1 space-y-3 overflow-y-auto p-4">
					{#if messages.length === 0}
						<p class="text-xs text-slate-500">No messages yet.</p>
					{:else}
						{#each messages as message (message.id)}
							<div
								class={`rounded-lg border px-3 py-2 text-sm shadow ${
									message.role === 'assistant'
										? 'border-slate-700/70 bg-slate-900/80 text-slate-200'
										: 'border-primary/30 bg-primary/10 text-slate-200'
								}`}
							>
								<div class="mb-1 flex items-center justify-between text-[10px] uppercase tracking-wide text-slate-400">
									<span>{message.role}</span>
									<span>{new Date(message.created_at).toLocaleTimeString()}</span>
								</div>
								<p class="whitespace-pre-line">{message.content}</p>
							</div>
						{/each}
					{/if}
				</div>
				<form class="border-t border-slate-800/80 p-4" on:submit={onSendMessage}>
					<div class="flex flex-col gap-3">
						<textarea
							class="h-24 rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-sm text-slate-100 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/40"
							placeholder="Write a message…"
							value={composeContent}
							on:input={(event) => onComposeContentChange((event.target as HTMLTextAreaElement).value)}
							required
						></textarea>
						<div class="flex flex-wrap items-center gap-3 text-xs text-slate-300">
							<label class="flex items-center gap-2">
								<span>Role</span>
								<select
									class="rounded border border-slate-700 bg-slate-900 px-2 py-1"
									value={composeRole}
									on:change={(event) =>
										onComposeRoleChange((event.target as HTMLSelectElement).value as 'user' | 'assistant')}
								>
									<option value="user">User</option>
									<option value="assistant">Assistant</option>
								</select>
							</label>
							<label class="flex items-center gap-2">
								<input
									type="checkbox"
									checked={generateReply}
									on:change={(event) =>
										onGenerateReplyChange((event.target as HTMLInputElement).checked)}
								/>
								Request AI reply
							</label>
							<button
								class="ml-auto rounded-lg bg-primary px-4 py-2 text-xs font-semibold text-white transition hover:bg-primary/90 disabled:opacity-60"
								type="submit"
								disabled={isSending}
							>
								{isSending ? 'Sending…' : 'Send'}
							</button>
						</div>
					</div>
				</form>
			</div>

			<div class="flex h-full flex-col gap-4">
				<section class="rounded-xl border border-slate-800/80 bg-slate-950/70">
					<header class="border-b border-slate-800/80 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-slate-400">
						Shared transcripts
					</header>
					<div class="max-h-56 space-y-2 overflow-y-auto px-4 py-3 text-xs text-slate-300">
						{#if transcriptsLoading}
							<p class="text-slate-500">Loading transcripts…</p>
						{:else if transcripts.length === 0}
							<p class="text-slate-500">No transcripts generated yet.</p>
						{:else}
							{#each transcripts as transcript (transcript.id)}
								<div class="rounded border border-slate-700/70 bg-slate-900/60 px-3 py-2">
									<p class="font-mono text-[11px] text-primary">{transcript.share_code}</p>
									<p class="text-[10px] text-slate-500">
										Created {formatTimestamp(transcript.created_at)}
									</p>
									{#if transcript.expires_at}
										<p class="text-[10px] text-slate-500">
											Expires {formatTimestamp(transcript.expires_at)}
										</p>
									{/if}
								</div>
							{/each}
						{/if}
					</div>
				</section>

				<section class="flex-1 rounded-xl border border-slate-800/80 bg-slate-950/70">
					<header class="border-b border-slate-800/80 px-4 py-3 text-xs font-semibold uppercase tracking-wide text-slate-400">
						Assignment history
					</header>
					<div class="max-h-64 space-y-2 overflow-y-auto px-4 py-3 text-xs text-slate-300">
						{#if assignmentsLoading}
							<p class="text-slate-500">Loading assignments…</p>
						{:else if assignments.length === 0}
							<p class="text-slate-500">No assignment records yet.</p>
						{:else}
							{#each assignments as assignment (assignment.id)}
								<div class="rounded border border-slate-700/70 bg-slate-900/60 px-3 py-2">
									<p class="font-semibold text-slate-200">
										{assignment.agent_name || assignment.agent_id}
									</p>
									<p class="text-[10px] text-slate-500">
										Assigned {formatTimestamp(assignment.assigned_at)}
									</p>
									{#if assignment.unassigned_at}
										<p class="text-[10px] text-slate-500">
											Unassigned {formatTimestamp(assignment.unassigned_at)}
										</p>
									{/if}
									{#if assignment.notes}
										<p class="mt-1 text-[11px] text-slate-400">{assignment.notes}</p>
									{/if}
								</div>
							{/each}
						{/if}
					</div>
				</section>
			</div>
		</div>
	{/if}
</section>

