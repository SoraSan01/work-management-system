// ===============================
// Project Modal Elements
// ===============================
const projectModal = document.getElementById('project-modal');
const openModalBtn = document.getElementById('open-project-modal');
const closeModalBtn = document.getElementById('close-project-modal');
const projectForm = document.getElementById('project-form');
const modalTitle = document.getElementById('project-modal-title');
const submitBtn = document.getElementById('project-submit-btn');

const projectIdInput = document.getElementById('project-id');
const projectNameInput = document.getElementById('project-name');
const projectDescriptionInput = document.getElementById('project-description');
const projectStatusInput = document.getElementById('project-status');
const projectTeamInput = document.getElementById('project-team');
const projectStartDateInput = document.getElementById('project-start');
const projectEndDateInput = document.getElementById('project-end');

// ===============================
// Delete Modal Elements
// ===============================
const deleteModal = document.getElementById('delete-project-modal');
const cancelDeleteBtn = document.getElementById('cancel-project-delete');
const confirmDeleteBtn = document.getElementById('confirm-project-delete');
const deleteProjectName = document.getElementById('delete-project-name');

let currentDeleteId = null;

// ===============================
// Open modal (Create)
// ===============================
openModalBtn?.addEventListener('click', () => {
	resetForm();

	modalTitle.textContent = 'Create New Project';
	submitBtn.textContent = 'Create Project';

	projectForm.action = '/projects';
	projectForm.method = 'POST';

	projectModal.classList.remove('hidden');
});

// ===============================
// Close modal
// ===============================
closeModalBtn?.addEventListener('click', () => {
	projectModal.classList.add('hidden');
	resetForm();
});

// Click outside modal
projectModal?.addEventListener('click', (e) => {
	if (e.target === projectModal) {
		projectModal.classList.add('hidden');
		resetForm();
	}
});

// ===============================
// Edit Project
// ===============================
document.querySelectorAll('.edit-project-btn').forEach(btn => {
	btn.addEventListener('click', (e) => {
		const row = e.target.closest('tr');

		projectIdInput.value = row.dataset.id || '';
		projectNameInput.value = row.dataset.name || '';
		projectStatusInput.value = row.dataset.status || '';
		projectDescriptionInput.value = row.dataset.description || '';
		projectStartDateInput.value = row.dataset.startDate || '';
		projectEndDateInput.value = row.dataset.endDate || '';
		projectTeamInput.value = row.dataset.teamId || '';

		modalTitle.textContent = 'Edit Project';
		submitBtn.textContent = 'Update Project';

		projectForm.action = `/projects/${row.dataset.id}`;
		projectForm.method = 'POST';

		projectModal.classList.remove('hidden');
	});
});

// ===============================
// Delete Project
// ===============================
document.querySelectorAll('.delete-project-btn').forEach(btn => {
	btn.addEventListener('click', () => {
		currentDeleteId = btn.dataset.id;
		deleteProjectName.textContent = btn.dataset.name;

		deleteModal.classList.remove('hidden');
	});
});

cancelDeleteBtn?.addEventListener('click', () => {
	deleteModal.classList.add('hidden');
	currentDeleteId = null;
});

deleteModal?.addEventListener('click', (e) => {
	if (e.target === deleteModal) {
		deleteModal.classList.add('hidden');
		currentDeleteId = null;
	}
});

confirmDeleteBtn?.addEventListener('click', () => {
	if (!currentDeleteId) return;

	const form = document.createElement('form');
	form.method = 'POST';
	form.action = `/projects/${currentDeleteId}/delete`;
	document.body.appendChild(form);
	form.submit();
});

// ===============================
// Helpers
// ===============================
function resetForm() {
	projectForm.reset();
	projectIdInput.value = '';
	projectTeamInput.value = '';
}

// Escape key closes modals
document.addEventListener('keydown', (e) => {
	if (e.key === 'Escape') {
		projectModal.classList.add('hidden');
		deleteModal.classList.add('hidden');
		currentDeleteId = null;
	}
});
