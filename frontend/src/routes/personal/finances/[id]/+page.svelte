<script lang="ts">
	import { derived } from 'svelte/store';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { createTransactionDetailLogic } from './transaction-detail.logic';
	import FinancesHero from '../_components/FinancesHero.svelte';
	import AuthNotice from '../_components/AuthNotice.svelte';
	import ActionError from '../_components/ActionError.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import StatCard from '$lib/components/StatCard.svelte';

	const transactionIdStore = derived(page, ($page) => $page.params.id ?? null);

	const {
		isAuthenticated,
		error,
		transactionQuery,
		transaction,
		refresh,
		handleToggle,
		numberFormatter
	} = createTransactionDetailLogic(transactionIdStore);
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
	<FinancesHero isAuthenticated={$isAuthenticated} onRefresh={refresh} />

	<!-- Main Content -->
	<div class="mx-auto max-w-4xl px-6 py-8 lg:px-8">
		{#if $error}
			<div class="mb-6">
				<ActionError message={$error} />
			</div>
		{/if}

		{#if !$isAuthenticated}
			<div class="mx-auto max-w-2xl">
				<AuthNotice />
			</div>
		{:else if $transactionQuery.isLoading}
			<div
				class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-12 text-center backdrop-blur-sm"
			>
				<div class="mx-auto mb-4 h-12 w-12 animate-spin rounded-full border-4 border-slate-700 border-t-emerald-500"></div>
				<p class="text-sm font-medium text-slate-400">Loading transaction details...</p>
			</div>
		{:else if !$transaction}
			<EmptyState
				title="Transaction not found"
				description="The transaction you're looking for doesn't exist or has been removed."
			>
				<button
					class="mt-4 inline-block rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white transition-all hover:bg-emerald-700"
					onclick={() => goto('/personal/finances')}
				>
					Return to Finances
				</button>
			</EmptyState>
		{:else}
			<div class="space-y-6">
				<!-- Back Button -->
				<button
					class="flex items-center gap-2 text-sm font-medium text-slate-400 transition-colors hover:text-slate-200"
					onclick={() => goto('/personal/finances')}
				>
					<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M15 19l-7-7 7-7"
						></path>
					</svg>
					Back to Transactions
				</button>

				<!-- Transaction Header Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-gradient-to-br from-slate-900/60 to-slate-800/40 p-8 shadow-xl backdrop-blur-sm"
					class:border-emerald-500={$transaction.type === 'income'}
					class:border-rose-500={$transaction.type === 'expense'}
				>
					<div class="flex items-start justify-between">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-4">
								<span
									class="inline-flex h-10 w-10 items-center justify-center rounded-full {$transaction.type === 'income'
										? 'bg-emerald-500/20 text-emerald-300'
										: 'bg-rose-500/20 text-rose-300'} text-lg font-bold"
								>
									{$transaction.type === 'income' ? '↑' : '↓'}
								</span>
								<span
									class="inline-flex items-center rounded-full border px-3 py-1 text-sm font-medium capitalize {$transaction.type === 'income'
										? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30'
										: 'bg-rose-500/20 text-rose-300 border-rose-500/30'}"
								>
									{$transaction.type}
								</span>
							</div>
							<h1 class="text-3xl font-bold text-white mb-2">{$transaction.category}</h1>
							{#if $transaction.description}
								<p class="text-slate-400 mb-4">{$transaction.description}</p>
							{/if}
						</div>
					</div>

					<!-- Amount Display -->
					<div class="mt-6 pt-6 border-t border-slate-800/50">
						<p class="text-sm text-slate-400 mb-2">Amount</p>
						<p
							class="text-4xl font-bold {$transaction.type === 'income'
								? 'text-emerald-400'
								: 'text-rose-400'}"
						>
							{numberFormatter.format($transaction.amount)} {$transaction.currency}
						</p>
						{#if $transaction.currency !== $transaction.base_currency}
							<p class="mt-2 text-sm text-slate-500">
								≈ {numberFormatter.format($transaction.normalized_amount)} {$transaction.base_currency}
							</p>
						{/if}
					</div>
				</div>

				<!-- Stats Grid -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<StatCard
						label="Date"
						value={new Date($transaction.occurred_at).toLocaleDateString()}
						accentColor="blue"
					/>
					<StatCard
						label="Time"
						value={new Date($transaction.occurred_at).toLocaleTimeString()}
						accentColor="purple"
					/>
					<StatCard
						label="Base Currency"
						value={$transaction.base_currency}
						accentColor="amber"
					/>
				</div>

				<!-- Details Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Details</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Category</p>
							<p class="text-sm font-medium text-slate-200">{$transaction.category}</p>
						</div>
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Type</p>
							<p class="text-sm font-medium text-slate-200 capitalize">{$transaction.type}</p>
						</div>
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Currency</p>
							<p class="text-sm font-medium text-slate-200">{$transaction.currency}</p>
						</div>
						<div>
							<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-1">Base Currency</p>
							<p class="text-sm font-medium text-slate-200">{$transaction.base_currency}</p>
						</div>
						{#if $transaction.tags && $transaction.tags.length > 0}
							<div class="md:col-span-2">
								<p class="text-xs font-medium uppercase tracking-wider text-slate-400 mb-2">Tags</p>
								<div class="flex flex-wrap gap-2">
									{#each $transaction.tags as tag}
										<span
											class="inline-flex items-center rounded-full bg-purple-500/20 px-3 py-1 text-xs font-medium text-purple-300"
										>
											{tag}
										</span>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</div>

				<!-- Flags Card -->
				<div
					class="rounded-xl border border-slate-800/50 bg-slate-900/40 p-6 shadow-xl backdrop-blur-sm"
				>
					<h2 class="text-lg font-semibold text-slate-100 mb-4">Flags & Status</h2>
					<div class="space-y-3">
						<label class="flex items-center justify-between cursor-pointer">
							<div class="flex items-center gap-3">
								<span class="text-sm font-medium text-slate-200">Recurring</span>
								<span
									class="inline-flex items-center rounded-full bg-indigo-500/20 px-2 py-0.5 text-xs font-medium text-indigo-300"
								>
									{$transaction.is_recurring ? 'Yes' : 'No'}
								</span>
							</div>
							<button
								class="rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-indigo-500/50 hover:bg-slate-800/80"
								onclick={() => handleToggle('recurring', !$transaction.is_recurring)}
							>
								{$transaction.is_recurring ? 'Disable' : 'Enable'}
							</button>
						</label>
						<label class="flex items-center justify-between cursor-pointer">
							<div class="flex items-center gap-3">
								<span class="text-sm font-medium text-slate-200">Essential</span>
								<span
									class="inline-flex items-center rounded-full bg-amber-500/20 px-2 py-0.5 text-xs font-medium text-amber-300"
								>
									{$transaction.is_essential ? 'Yes' : 'No'}
								</span>
							</div>
							<button
								class="rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-amber-500/50 hover:bg-slate-800/80"
								onclick={() => handleToggle('essential', !$transaction.is_essential)}
							>
								{$transaction.is_essential ? 'Disable' : 'Enable'}
							</button>
						</label>
						<label class="flex items-center justify-between cursor-pointer">
							<div class="flex items-center gap-3">
								<span class="text-sm font-medium text-slate-200">Archived</span>
								<span
									class="inline-flex items-center rounded-full bg-slate-700/40 px-2 py-0.5 text-xs font-medium text-slate-300"
								>
									{$transaction.is_archived ? 'Yes' : 'No'}
								</span>
							</div>
							<button
								class="rounded-lg border border-slate-700 bg-slate-800/50 px-4 py-2 text-sm font-medium text-slate-200 transition-all hover:border-rose-500/50 hover:bg-slate-800/80"
								onclick={() => handleToggle('archive', !$transaction.is_archived)}
							>
								{$transaction.is_archived ? 'Unarchive' : 'Archive'}
							</button>
						</label>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

