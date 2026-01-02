// Role Modal Elements
const roleModal = document.getElementById('role-modal');
const openModalBtn = document.getElementById('open-modal');
const closeModalBtn = document.getElementById('close-modal');
const roleForm = document.getElementById('role-form');
const modalTitle = document.getElementById('modal-title');
const submitBtn = document.getElementById('submit-btn');
const roleIdInput = document.getElementById('role-id');
const roleNameInput = document.getElementById('role-name');

// Delete Modal Elements
const deleteModal = document.getElementById('delete-modal');
const cancelDeleteBtn = document.getElementById('cancel-delete');
const confirmDeleteBtn = document.getElementById('confirm-delete');
const deleteRoleName = document.getElementById('delete-role-name');

let currentDeleteId = null;

// Open modal for creating new role
openModalBtn?.addEventListener('click', () => {
    resetForm();
    modalTitle.textContent = 'Create New Role';
    submitBtn.textContent = 'Create Role';
    roleForm.action = '/roles';
    roleForm.method = 'POST';
    roleModal.classList.remove('hidden');
});

// Close modal
closeModalBtn?.addEventListener('click', () => {
    roleModal.classList.add('hidden');
    resetForm();
});

// Close modal when clicking outside
roleModal?.addEventListener('click', (e) => {
    if (e.target === roleModal) {
        roleModal.classList.add('hidden');
        resetForm();
    }
});

// Handle edit button clicks
document.querySelectorAll('.edit-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        const row = e.target.closest('tr');
        const id = row.dataset.id;
        const name = row.dataset.name;
        
        // Populate form with existing data
        roleIdInput.value = id;
        roleNameInput.value = name;
        
        // Update modal title and button
        modalTitle.textContent = 'Edit Role';
        submitBtn.textContent = 'Update Role';
        
        // Update form action for update
        roleForm.action = `/roles/${id}`;
        roleForm.method = 'POST';
        
        // Show modal
        roleModal.classList.remove('hidden');
    });
});

// Handle delete button clicks
document.querySelectorAll('.delete-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        currentDeleteId = e.target.dataset.id;
        const roleName = e.target.dataset.name;
        
        deleteRoleName.textContent = roleName;
        deleteModal.classList.remove('hidden');
    });
});

// Cancel delete
cancelDeleteBtn?.addEventListener('click', () => {
    deleteModal.classList.add('hidden');
    currentDeleteId = null;
});

// Close delete modal when clicking outside
deleteModal?.addEventListener('click', (e) => {
    if (e.target === deleteModal) {
        deleteModal.classList.add('hidden');
        currentDeleteId = null;
    }
});

// Confirm delete
confirmDeleteBtn?.addEventListener('click', () => {
    if (currentDeleteId) {
        // Create and submit a form for deletion
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = `/roles/${currentDeleteId}/delete`;
        document.body.appendChild(form);
        form.submit();
    }
});

// Reset form to initial state
function resetForm() {
    roleForm.reset();
    roleIdInput.value = '';
    roleNameInput.value = '';
}

// Close modals on Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        if (!roleModal.classList.contains('hidden')) {
            roleModal.classList.add('hidden');
            resetForm();
        }
        if (!deleteModal.classList.contains('hidden')) {
            deleteModal.classList.add('hidden');
            currentDeleteId = null;
        }
    }
});