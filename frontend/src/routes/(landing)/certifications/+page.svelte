<script lang="ts">
	import { onMount } from 'svelte';
	import {
		listCertifications,
		createCertification,
		updateCertification,
		deleteCertification,
		type Certification,
		type CreateCertificationInput,
		type UpdateCertificationInput,
		type CertificationStatus,
		type CertificationCategory
	} from '$lib/api/certifications';

	let certifications: Certification[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let showCreateModal = $state(false);
	let showEditModal = $state(false);
	let editingCertification: Certification | null = $state(null);
	let searchQuery = $state('');

	// Form state
	let formName = $state('');
	let formIssuer = $state('');
	let formIssueDate = $state('');
	let formExpiryDate = $state('');
	let formCredentialId = $state('');
	let formVerificationUrl = $state('');
	let formCertificateUrl = $state('');
	let formDescription = $state('');
	let formStatus = $state<CertificationStatus>('active');
	let formCategory = $state<CertificationCategory>('other');
	let formFeatured = $state(false);
	let formDisplayOrder = $state<number | ''>(0);

	const statuses: CertificationStatus[] = ['active', 'expired', 'revoked', 'renewed'];
	const categories: CertificationCategory[] = [
		'cloud',
		'security',
		'programming',
		'database',
		'devops',
		'architecture',
		'other'
	];

	onMount(async () => {
		await fetchCertifications();
	});

	async function fetchCertifications() {
		loading = true;
		error = null;
		try {
			certifications = await listCertifications();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load certifications';
			console.error('Error fetching certifications:', err);
		} finally {
			loading = false;
		}
	}

	function openCreateModal() {
		resetForm();
		showCreateModal = true;
	}

	function openEditModal(cert: Certification) {
		editingCertification = cert;
		formName = cert.name;
		formIssuer = cert.issuer;
		formIssueDate = cert.issueDate.split('T')[0];
		formExpiryDate = cert.expiryDate ? cert.expiryDate.split('T')[0] : '';
		formCredentialId = cert.credentialId || '';
		formVerificationUrl = cert.verificationUrl || '';
		formCertificateUrl = cert.certificateUrl || '';
		formDescription = cert.description || '';
		formStatus = cert.status;
		formCategory = cert.category;
		formFeatured = cert.featured;
		formDisplayOrder = cert.displayOrder;
		showEditModal = true;
	}

	function resetForm() {
		formName = '';
		formIssuer = '';
		formIssueDate = '';
		formExpiryDate = '';
		formCredentialId = '';
		formVerificationUrl = '';
		formCertificateUrl = '';
		formDescription = '';
		formStatus = 'active';
		formCategory = 'other';
		formFeatured = false;
		formDisplayOrder = 0;
		editingCertification = null;
	}

	async function handleCreate() {
		if (!formName.trim() || !formIssuer.trim() || !formIssueDate) {
			alert('Name, issuer, and issue date are required');
			return;
		}

		try {
			const input: CreateCertificationInput = {
				name: formName.trim(),
				issuer: formIssuer.trim(),
				issueDate: formIssueDate,
				expiryDate: formExpiryDate || undefined,
				credentialId: formCredentialId.trim() || undefined,
				verificationUrl: formVerificationUrl.trim() || undefined,
				certificateUrl: formCertificateUrl.trim() || undefined,
				description: formDescription.trim() || undefined,
				status: formStatus,
				category: formCategory,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0
			};

			await createCertification(input);
			showCreateModal = false;
			resetForm();
			await fetchCertifications();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to create certification');
			console.error('Error creating certification:', err);
		}
	}

	async function handleUpdate() {
		if (!editingCertification || !formName.trim() || !formIssuer.trim() || !formIssueDate) {
			alert('Name, issuer, and issue date are required');
			return;
		}

		try {
			const input: UpdateCertificationInput = {
				name: formName.trim(),
				issuer: formIssuer.trim(),
				issueDate: formIssueDate,
				expiryDate: formExpiryDate || undefined,
				credentialId: formCredentialId.trim() || undefined,
				verificationUrl: formVerificationUrl.trim() || undefined,
				certificateUrl: formCertificateUrl.trim() || undefined,
				description: formDescription.trim() || undefined,
				status: formStatus,
				category: formCategory,
				featured: formFeatured,
				displayOrder: formDisplayOrder ? Number(formDisplayOrder) : 0
			};

			await updateCertification(editingCertification.id, input);
			showEditModal = false;
			resetForm();
			await fetchCertifications();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update certification');
			console.error('Error updating certification:', err);
		}
	}

	async function handleDelete(id: string) {
		if (!confirm('Are you sure you want to delete this certification?')) return;

		try {
			await deleteCertification(id);
			await fetchCertifications();
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete certification');
			console.error('Error deleting certification:', err);
		}
	}

	function filteredCertifications() {
		if (!searchQuery.trim()) return certifications;
		const query = searchQuery.toLowerCase();
		return certifications.filter(
			(c) =>
				c.name.toLowerCase().includes(query) ||
				c.issuer.toLowerCase().includes(query) ||
				c.description?.toLowerCase().includes(query)
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
			<h1>Certifications Management</h1>
			<p>Manage certifications and credentials</p>
		</div>
		<button onclick={openCreateModal}>Create Certification</button>
	</div>

	<div class="search-bar">
		<input
			type="text"
			placeholder="Search certifications..."
			bind:value={searchQuery}
			class="search-input"
		/>
	</div>

	{#if error}
		<div class="error">{error}</div>
	{/if}

	{#if loading}
		<div class="loading">Loading...</div>
	{:else if filteredCertifications().length === 0}
		<div class="empty">No certifications found</div>
	{:else}
		<table class="table">
			<thead>
				<tr>
					<th>Name</th>
					<th>Issuer</th>
					<th>Category</th>
					<th>Status</th>
					<th>Issue Date</th>
					<th>Expiry Date</th>
					<th>Featured</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each filteredCertifications() as cert}
					<tr>
						<td>
							<strong>{cert.name}</strong>
							{#if cert.description}
								<br />
								<small>{cert.description}</small>
							{/if}
						</td>
						<td>{cert.issuer}</td>
						<td>{cert.category}</td>
						<td>
							<span class="status status-{cert.status}">{cert.status}</span>
						</td>
						<td>{formatDate(cert.issueDate)}</td>
						<td>{formatDate(cert.expiryDate)}</td>
						<td>{cert.featured ? 'Yes' : 'No'}</td>
						<td>
							<button onclick={() => openEditModal(cert)}>Edit</button>
							<button onclick={() => handleDelete(cert.id)} class="delete-btn">Delete</button>
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
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Create Certification</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Issuer *</label>
					<input type="text" bind:value={formIssuer} />
				</div>
				<div class="form-group">
					<label>Issue Date *</label>
					<input type="date" bind:value={formIssueDate} />
				</div>
				<div class="form-group">
					<label>Expiry Date</label>
					<input type="date" bind:value={formExpiryDate} />
				</div>
				<div class="form-group">
					<label>Category</label>
					<select bind:value={formCategory}>
						{#each categories as cat}
							<option value={cat}>{cat}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Status</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{status}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Credential ID</label>
					<input type="text" bind:value={formCredentialId} />
				</div>
				<div class="form-group">
					<label>Verification URL</label>
					<input type="url" bind:value={formVerificationUrl} />
				</div>
				<div class="form-group">
					<label>Certificate URL</label>
					<input type="url" bind:value={formCertificateUrl} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
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
{#if showEditModal && editingCertification}
	<div class="modal-overlay" onclick={() => (showEditModal = false)}>
		<div class="modal" onclick={(e) => e.stopPropagation()}>
			<h2>Edit Certification</h2>
			<div class="form">
				<div class="form-group">
					<label>Name *</label>
					<input type="text" bind:value={formName} />
				</div>
				<div class="form-group">
					<label>Issuer *</label>
					<input type="text" bind:value={formIssuer} />
				</div>
				<div class="form-group">
					<label>Issue Date *</label>
					<input type="date" bind:value={formIssueDate} />
				</div>
				<div class="form-group">
					<label>Expiry Date</label>
					<input type="date" bind:value={formExpiryDate} />
				</div>
				<div class="form-group">
					<label>Category</label>
					<select bind:value={formCategory}>
						{#each categories as cat}
							<option value={cat}>{cat}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Status</label>
					<select bind:value={formStatus}>
						{#each statuses as status}
							<option value={status}>{status}</option>
						{/each}
					</select>
				</div>
				<div class="form-group">
					<label>Credential ID</label>
					<input type="text" bind:value={formCredentialId} />
				</div>
				<div class="form-group">
					<label>Verification URL</label>
					<input type="url" bind:value={formVerificationUrl} />
				</div>
				<div class="form-group">
					<label>Certificate URL</label>
					<input type="url" bind:value={formCertificateUrl} />
				</div>
				<div class="form-group">
					<label>Description</label>
					<textarea bind:value={formDescription} rows="3"></textarea>
				</div>
				<div class="form-group">
					<label>
						<input type="checkbox" bind:checked={formFeatured} />
						Featured
					</label>
				</div>
				<div class="form-group">
					<label>Display Order</label>
					<input type="number" bind:value={formDisplayOrder} />
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

	.table small {
		color: #666;
		font-size: 0.875rem;
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

	.status-expired {
		background: #f8d7da;
		color: #721c24;
	}

	.status-revoked {
		background: #f8d7da;
		color: #721c24;
	}

	.status-renewed {
		background: #d1ecf1;
		color: #0c5460;
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

