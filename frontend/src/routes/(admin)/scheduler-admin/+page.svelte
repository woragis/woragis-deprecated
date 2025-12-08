<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listSchedules,
		createSchedule,
		updateSchedule,
		deleteSchedule,
		bulkActivate,
		bulkDeactivate,
		bulkPause,
		bulkResume,
		type Schedule,
		type CreateScheduleInput,
		type UpdateScheduleInput
	} from '$lib/api/scheduler';

	let schedules: Schedule[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingSchedule: Schedule | null = $state(null);
	let searchQuery = $state('');
	let selectedIds = $state<Set<string>>(new Set());

	// Form state
	let formReportType = $state('');
	let formAgentAlias = $state('');
	let formFrequency = $state('daily');
	let formWeekday = $state('');
	let formTimeOfDay = $state('09:00');
	let formTimezone = $state('UTC');
	let formRrule = $state('');
	let formPriority = $state<number | ''>(0);
	let formEmail = $state('');
	let formPhoneNumber = $state('');
	let formActive = $state(true);
	let formPaused = $state(false);

	const frequencies = ['daily', 'weekly', 'custom'];

	onMount(async () => {
		await fetchSchedules();
	});

	async function fetchSchedules() {
		loading = true;
		error = null;
		try {
			schedules = await listSchedules();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load schedules';
			console.error('Error fetching schedules:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(schedule: Schedule) {
		editingSchedule = schedule;
		formReportType = schedule.reportType;
		formAgentAlias = schedule.agentAlias;
		formFrequency = schedule.frequency;
		formWeekday = schedule.weekday || '';
		formTimeOfDay = schedule.timeOfDay;
		formTimezone = schedule.timezone;
		formRrule = schedule.rrule || '';
		formPriority = schedule.priority;
		formEmail = schedule.email || '';
		formPhoneNumber = schedule.phoneNumber || '';
		formActive = schedule.active;
		formPaused = schedule.paused;
		showEditModal = true;
	}

	function resetForm() {
		formReportType = '';
		formAgentAlias = '';
		formFrequency = 'daily';
		formWeekday = '';
		formTimeOfDay = '09:00';
		formTimezone = 'UTC';
		formRrule = '';
		formPriority = 0;
		formEmail = '';
		formPhoneNumber = '';
		formActive = true;
		formPaused = false;
		editingSchedule = null;
	}

	async function handleCreate() {
		if (!formReportType.trim() || !formAgentAlias.trim() || !formTimeOfDay.trim()) {
			alert('Report type, agent alias, and time of day are required');
			return;
		}

		try {
			const input: CreateScheduleInput = {
				reportType: formReportType.trim(),
				agentAlias: formAgentAlias.trim(),
				frequency: formFrequency,
				weekday: formWeekday.trim() || undefined,
				timeOfDay: formTimeOfDay.trim(),
				timezone: formTimezone.trim(),
				rrule: formRrule.trim() || undefined,
				priority: formPriority ? Number(formPriority) : 0,
				email: formEmail.trim() || undefined,
				phoneNumber: formPhoneNumber.trim() || undefined,
				active: formActive
			};

			await createSchedule(input);
			showCreateModal = false;
			resetForm();
			await fetchSchedules();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create schedule');
			console.error('Error creating schedule:', err);
		}
	}

	async function handleUpdate() {
		if (!editingSchedule) {
			alert('No schedule selected');
			return;
		}

		try {
			const input: UpdateScheduleInput = {
				reportType: formReportType.trim(),
				agentAlias: formAgentAlias.trim(),
				frequency: formFrequency,
				weekday: formWeekday.trim() || undefined,
				timeOfDay: formTimeOfDay.trim(),
				timezone: formTimezone.trim(),
				rrule: formRrule.trim() || undefined,
				priority: formPriority ? Number(formPriority) : 0,
				email: formEmail.trim() || undefined,
				phoneNumber: formPhoneNumber.trim() || undefined,
				active: formActive,
				paused: formPaused
			};

			await updateSchedule(editingSchedule.id, input);
			showEditModal = false;
			resetForm();
			await fetchSchedules();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update schedule');
			console.error('Error updating schedule:', err);
		}
	}

	function toggleSelect(id: string) {
		if (selectedIds.has(id)) {
			selectedIds.delete(id);
		} else {
			selectedIds.add(id);
		}
		selectedIds = new Set(selectedIds);
	}

	async function handleBulkActivate() {
		if (selectedIds.size === 0) return;
		try {
			await bulkActivate(Array.from(selectedIds));
			selectedIds.clear();
			await fetchSchedules();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to activate schedules');
		}
	}

	async function handleBulkDeactivate() {
		if (selectedIds.size === 0) return;
		try {
			await bulkDeactivate(Array.from(selectedIds));
			selectedIds.clear();
			await fetchSchedules();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to deactivate schedules');
		}
	}

	async function handleBulkPause() {
		if (selectedIds.size === 0) return;
		try {
			await bulkPause(Array.from(selectedIds));
			selectedIds.clear();
			await fetchSchedules();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to pause schedules');
		}
	}

	async function handleBulkResume() {
		if (selectedIds.size === 0) return;
		try {
			await bulkResume(Array.from(selectedIds));
			selectedIds.clear();
			await fetchSchedules();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to resume schedules');
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this schedule? This will also delete all execution runs.')) return;

		try {
			await deleteSchedule(id);
			await fetchSchedules();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete schedule');
			console.error('Error deleting schedule:', err);
		}
	}

	function filteredSchedules() {
		if (!searchQuery.trim()) return schedules;
		const query = searchQuery.toLowerCase();
		return schedules.filter(
			(s) =>
				s.reportType.toLowerCase().includes(query) ||
				s.agentAlias.toLowerCase().includes(query) ||
				s.frequency.toLowerCase().includes(query)
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
			<h1>Scheduler Management</h1>
			<p>Manage scheduled reports</p>
		</div>
		<button onclick={openCreateModal}>Create Schedule</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search schedules..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if selectedIds.size > 0}
		<div class="bulk-actions">
			<span>{selectedIds.size} selected</span>
			<button onclick={handleBulkActivate}>Activate</button>
			<button onclick={handleBulkDeactivate}>Deactivate</button>
			<button onclick={handleBulkPause}>Pause</button>
			<button onclick={handleBulkResume}>Resume</button>
			<button onclick={() => (selectedIds.clear())}>Clear</button>
		</div>
	{/if}

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredSchedules().length === 0}
		<div class="empty">No schedules found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>
						<input
							type="checkbox"
							checked={selectedIds.size === filteredSchedules().length && filteredSchedules().length > 0}
							onchange={() => {
								if (selectedIds.size === filteredSchedules().length) {
									selectedIds.clear();
								} else {
									selectedIds = new Set(filteredSchedules().map((s) => s.id));
								}
							}}
						/>
					</th>
					<th>Report Type</th>
					<th>Agent</th>
					<th>Frequency</th>
					<th>Time</th>
					<th>Timezone</th>
					<th>Status</th>
					<th>Next Run</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredSchedules() as schedule}
					<tr>
						<td>
							<input
								type="checkbox"
								checked={selectedIds.has(schedule.id)}
								onchange={() => toggleSelect(schedule.id)}
							/>
						</td>
						<td>{schedule.reportType}</td>
						<td>{schedule.agentAlias}</td>
						<td>{schedule.frequency}</td>
						<td>{schedule.timeOfDay}</td>
						<td>{schedule.timezone}</td>
						<td>
							{#if schedule.paused}
								<span class="status status-paused">Paused</span>
							{:else if schedule.active}
								<span class="status status-active">Active</span>
							{:else}
								<span class="status status-inactive">Inactive</span>
							{/if}
						</td>
						<td>{formatDate(schedule.nextRun)}</td>
						<td>
							<button onclick={() => openEditModal(schedule)}>Edit</button>
							<button onclick={() => handleDelete(schedule.id)} class="delete-btn">Delete</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>

<!-- Create Modal -->
{#if showCreateModal}
	<div class="modal-overlay" onclick={() => (showCreateModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Create Schedule</h2>
			<div class="form">
				<div class="form-group">
					<label>Report Type *</label>
					<input type="text" bind:value={formReportType} />
				</div>
				<div class="form-group">
					<label>Agent Alias *</label>
					<input type="text" bind:value={formAgentAlias} />
				</div>
				<div class="form-group">
					<label>Frequency</label>
					<select bind:value={formFrequency}>
						{#each frequencies as freq}
							<option value={freq}>{freq}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Weekday (for weekly)</label>
					<input type="text" bind:value={formWeekday} placeholder="monday, tuesday, etc." />
				</div>
				<div class="form-group">
					<label>Time of Day *</label>
					<input type="time" bind:value={formTimeOfDay} />
				</div>
				<div class="form-group">
					<label>Timezone</label>
					<input type="text" bind:value={formTimezone} />
				</div>
				<div class="form-group">
					<label>RRule (for custom)</label>
					<input type="text" bind:value={formRrule} />
				</div>
				<div class="form-group">
					<label>Priority</label>
					<input type="number" bind:value={formPriority} />
				</div>
				<div class="form-group">
					<label>Email</label>
					<input type="email" bind:value={formEmail} />
				</div>
				<div class="form-group">
					<label>Phone Number</label>
					<input type="tel" bind:value={formPhoneNumber} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formActive} />
						Active
					</label>
				</div>
				<div class="form-actions">
					<button onclick={handleCreate}>Create</button>
					<button onclick={() => (showCreateModal = false)}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal && editingSchedule}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal modal-large" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Schedule</h2>
			<div class="form">
				<div class="form-group">
					<label>Report Type *</label>
					<input type="text" bind:value={formReportType} />
				</div>
				<div class="form-group">
					<label>Agent Alias *</label>
					<input type="text" bind:value={formAgentAlias} />
				</div>
				<div class="form-group">
					<label>Frequency</label>
					<select bind:value={formFrequency}>
						{#each frequencies as freq}
							<option value={freq}>{freq}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Weekday (for weekly)</label>
					<input type="text" bind:value={formWeekday} placeholder="monday, tuesday, etc." />
				</div>
				<div class="form-group">
					<label>Time of Day *</label>
					<input type="time" bind:value={formTimeOfDay} />
				</div>
				<div class="form-group">
					<label>Timezone</label>
					<input type="text" bind:value={formTimezone} />
				</div>
				<div class="form-group">
					<label>RRule (for custom)</label>
					<input type="text" bind:value={formRrule} />
				</div>
				<div class="form-group">
					<label>Priority</label>
					<input type="number" bind:value={formPriority} />
				</div>
				<div class="form-group">
					<label>Email</label>
					<input type="email" bind:value={formEmail} />
				</div>
				<div class="form-group">
					<label>Phone Number</label>
					<input type="tel" bind:value={formPhoneNumber} />
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formActive} />
						Active
					</label>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formPaused} />
						Paused
					</label>
				</div>
				<div class="form-actions">
					<button onclick={handleUpdate}>Update</button>
					<button onclick={() => (showEditModal = false)}>Cancel</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.page-container {
		padding: 1rem;
		max-width: 1200px;
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

	.bulk-actions {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.75rem;
		background: #f5f5f5;
		border-radius: 4px;
		margin-bottom: 1rem;
	}

	.bulk-actions button {
		padding: 0.375rem 0.75rem;
		background: #6c757d;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.875rem;
	}

	.bulk-actions button:hover {
		background: #5a6268;
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

	.status-active {
		background: #d4edda;
		color: #155724;
	}

	.status-paused {
		background: #fff3cd;
		color: #856404;
	}

	.status-inactive {
		background: #e2e3e5;
		color: #383d41;
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
		max-width: 800px;
		width: 90%;
		max-height: 90vh;
		overflow-y: auto;
	}

	.modal-large {
		max-width: 900px;
	}

	.modal h2 {
		margin: 0 0 1rem 0;
		font-size: 1.25rem;
	}

	.form {
		display: flex;
		flex-direction: column;
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
	}

	.form-group input,
	.form-group textarea,
	.form-group select {
		padding: 0.5rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 0.875rem;
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

