<script lang="ts">
	export let show = false;
	export let onClose: () => void;
</script>

{#if show}
	<div class="rounded-2xl border border-slate-800/80 bg-slate-950/60 p-6">
		<div class="mb-4 flex items-center justify-between">
			<h2 class="text-xl font-semibold text-slate-100">📚 Reports Tutorial</h2>
			<button
				class="rounded border border-slate-700 px-2 py-1 text-sm text-slate-300 hover:border-slate-500"
				on:click={onClose}
			>
				Close
			</button>
		</div>

		<div class="space-y-6 text-sm text-slate-300">
			<!-- Overview -->
			<section class="space-y-2">
				<h3 class="text-base font-semibold text-slate-100">What are Reports?</h3>
				<p>
					Reports allow you to generate and schedule automated summaries of your data. Each report consists of
					three main components:
				</p>
				<ul class="ml-4 list-disc space-y-1 text-slate-400">
					<li>
						<strong class="text-slate-200">Definition:</strong> What data to include and how to filter it
					</li>
					<li>
						<strong class="text-slate-200">Schedule:</strong> When to automatically generate the report
					</li>
					<li>
						<strong class="text-slate-200">Delivery:</strong> How to send the report (email, etc.)
					</li>
				</ul>
			</section>

			<!-- Step 1: Create Definition -->
			<section class="space-y-2">
				<h3 class="text-base font-semibold text-slate-100">Step 1: Create a Report Definition</h3>
				<ol class="ml-4 list-decimal space-y-2 text-slate-400">
					<li>Click the <strong class="text-slate-200">"Create Definition"</strong> button</li>
					<li>
						Fill in the report details:
						<ul class="ml-4 mt-1 list-disc space-y-1">
							<li><strong>Name:</strong> A descriptive name for your report</li>
							<li><strong>Description:</strong> Optional description of what the report contains</li>
							<li>
								<strong>Sections:</strong> JSON object defining what sections to include (e.g.,
								<code class="rounded bg-slate-900 px-1 text-xs">{'{"overview": []}'}</code>)
							</li>
							<li>
								<strong>Filters:</strong> JSON object defining date ranges and filters (e.g.,
								<code class="rounded bg-slate-900 px-1 text-xs">{'{"date_range": "last_30_days"}'}</code>)
							</li>
						</ul>
					</li>
					<li>Click <strong class="text-slate-200">"Create definition"</strong></li>
				</ol>
			</section>

			<!-- Step 2: Add Schedule -->
			<section class="space-y-2">
				<h3 class="text-base font-semibold text-slate-100">Step 2: Add a Schedule</h3>
				<p class="text-slate-400">
					After creating a definition, select it and click <strong class="text-slate-200">"Add schedule"</strong>.
				</p>
				<div class="rounded-lg border border-slate-800/60 bg-slate-900/40 p-4">
					<h4 class="mb-2 text-sm font-semibold text-slate-200">Schedule Options:</h4>
					<ul class="space-y-2 text-slate-400">
						<li>
							<strong class="text-slate-200">Quick Presets:</strong> Choose from common schedules like
							"Daily", "Every 7 Days", "Every 14 Days", or "Every 30 Days". Adjust the time (hour and
							minute) as needed.
						</li>
						<li>
							<strong class="text-slate-200">Weekly:</strong> Select a specific weekday (Monday, Tuesday,
							etc.) and time. Perfect for weekly reports on a specific day.
						</li>
						<li>
							<strong class="text-slate-200">Monthly:</strong> Choose a day of the month (1-31) and time.
							Great for monthly reports on the 1st, 15th, or last day of the month.
						</li>
						<li>
							<strong class="text-slate-200">Custom Cron:</strong> For advanced users, enter a custom cron
							expression. Format: <code class="rounded bg-slate-950 px-1 text-xs">minute hour day-of-month month day-of-week</code>
						</li>
					</ul>
				</div>
				<div class="mt-3 rounded-lg border border-blue-500/30 bg-blue-500/10 p-3">
					<p class="text-xs text-blue-200">
						<strong>💡 Tip:</strong> The cron expression is automatically generated when you use presets or
						select weekly/monthly options. You can see the generated cron expression below the schedule
						options.
					</p>
				</div>
			</section>

			<!-- Step 3: Add Delivery -->
			<section class="space-y-2">
				<h3 class="text-base font-semibold text-slate-100">Step 3: Configure Delivery</h3>
				<p class="text-slate-400">
					Click <strong class="text-slate-200">"Add delivery"</strong> to specify how the report should be
					sent.
				</p>
				<ul class="ml-4 list-disc space-y-1 text-slate-400">
					<li>
						<strong class="text-slate-200">Channel:</strong> Choose how to deliver (e.g., "email")
					</li>
					<li>
						<strong class="text-slate-200">Target:</strong> The recipient (e.g., email address)
					</li>
					<li>
						<strong class="text-slate-200">Template:</strong> Optional JSON template for customizing the
						delivery format
					</li>
				</ul>
			</section>

			<!-- Examples -->
			<section class="space-y-2">
				<h3 class="text-base font-semibold text-slate-100">Common Use Cases</h3>
				<div class="space-y-3">
					<div class="rounded-lg border border-slate-800/60 bg-slate-900/40 p-3">
						<h4 class="mb-1 text-sm font-semibold text-slate-200">Daily Report at 8 AM</h4>
						<ol class="ml-4 list-decimal space-y-1 text-xs text-slate-400">
							<li>Create a definition with your desired filters</li>
							<li>Add schedule → Select "Quick Presets" → Click "Daily"</li>
							<li>Set hour to 8, minute to 0</li>
							<li>Add delivery channel (e.g., email)</li>
						</ol>
					</div>
					<div class="rounded-lg border border-slate-800/60 bg-slate-900/40 p-3">
						<h4 class="mb-1 text-sm font-semibold text-slate-200">Weekly Report Every Saturday at 9 AM</h4>
						<ol class="ml-4 list-decimal space-y-1 text-xs text-slate-400">
							<li>Create a definition</li>
							<li>Add schedule → Select "Weekly"</li>
							<li>Choose "Saturday" from weekday dropdown</li>
							<li>Set hour to 9, minute to 0</li>
							<li>Add delivery channel</li>
						</ol>
					</div>
					<div class="rounded-lg border border-slate-800/60 bg-slate-900/40 p-3">
						<h4 class="mb-1 text-sm font-semibold text-slate-200">Monthly Report on the 1st at 6 AM</h4>
						<ol class="ml-4 list-decimal space-y-1 text-xs text-slate-400">
							<li>Create a definition</li>
							<li>Add schedule → Select "Monthly"</li>
							<li>Set day of month to 1</li>
							<li>Set hour to 6, minute to 0</li>
							<li>Add delivery channel</li>
						</ol>
					</div>
					<div class="rounded-lg border border-slate-800/60 bg-slate-900/40 p-3">
						<h4 class="mb-1 text-sm font-semibold text-slate-200">Bi-Weekly Report (Every 14 Days)</h4>
						<ol class="ml-4 list-decimal space-y-1 text-xs text-slate-400">
							<li>Create a definition</li>
							<li>Add schedule → Select "Quick Presets" → Click "Every 14 Days"</li>
							<li>Adjust time as needed</li>
							<li>Add delivery channel</li>
						</ol>
					</div>
				</div>
			</section>

			<!-- Managing Reports -->
			<section class="space-y-2">
				<h3 class="text-base font-semibold text-slate-100">Managing Reports</h3>
				<ul class="ml-4 list-disc space-y-1 text-slate-400">
					<li>
						<strong class="text-slate-200">Queue Runs:</strong> Manually trigger a report generation using
						the bulk actions panel
					</li>
					<li>
						<strong class="text-slate-200">Enable/Disable:</strong> Toggle schedules and deliveries on/off
						without deleting them
					</li>
					<li>
						<strong class="text-slate-200">Favorites:</strong> Mark frequently used reports as favorites
						for quick access
					</li>
					<li>
						<strong class="text-slate-200">Archive:</strong> Archive reports you no longer need without
						deleting them
					</li>
				</ul>
			</section>

			<!-- Cron Reference -->
			<section class="space-y-2">
				<h3 class="text-base font-semibold text-slate-100">Cron Expression Reference</h3>
				<div class="rounded-lg border border-slate-800/60 bg-slate-900/40 p-3">
					<p class="mb-2 text-xs text-slate-400">
						Cron format: <code class="rounded bg-slate-950 px-1">minute hour day-of-month month day-of-week</code>
					</p>
					<div class="space-y-1 text-xs text-slate-400">
						<div><code class="rounded bg-slate-950 px-1">0 8 * * *</code> = Daily at 8:00 AM</div>
						<div><code class="rounded bg-slate-950 px-1">0 9 * * 1</code> = Every Monday at 9:00 AM</div>
						<div><code class="rounded bg-slate-950 px-1">0 6 1 * *</code> = 1st of every month at 6:00 AM</div>
						<div><code class="rounded bg-slate-950 px-1">0 */7 * * *</code> = Every 7 days at midnight</div>
						<div><code class="rounded bg-slate-950 px-1">30 14 * * 5</code> = Every Friday at 2:30 PM</div>
					</div>
					<p class="mt-2 text-xs text-slate-500">
						Day of week: 0 = Sunday, 1 = Monday, ..., 6 = Saturday
					</p>
				</div>
			</section>
		</div>
	</div>
{/if}

