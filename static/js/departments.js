// Department Modal Elements
const departmentModal = document.getElementById('department-modal');
const openModalBtn = document.getElementById('open-modal');
const closeModalBtn = document.getElementById('close-modal');
const departmentForm = document.getElementById('department-form');
const modalTitle = document.getElementById('modal-title');
const submitBtn = document.getElementById('submit-btn');
const departmentIdInput = document.getElementById('department-id');
const departmentNameInput = document.getElementById('department-name');
const statusActiveInput = document.getElementById('status-active');
const statusInactiveInput = document.getElementById('status-inactive');

// Delete Modal Elements
const deleteModal = document.getElementById('delete-modal');
const cancelDeleteBtn = document.getElementById('cancel-delete');
const confirmDeleteBtn = document.getElementById('confirm-delete');
const deleteDepartmentName = document.getElementById('delete-department-name');

let currentDeleteId = null;

// Open modal for creating new department
openModalBtn?.addEventListener('click', () => {
    resetForm();
    modalTitle.textContent = 'Create New Department';
    submitBtn.textContent = 'Create Department';
    departmentForm.action = '/departments';
    departmentForm.method = 'POST';
    departmentModal.classList.remove('hidden');
});

// Close modal
closeModalBtn?.addEventListener('click', () => {
    departmentModal.classList.add('hidden');
    resetForm();
});

// Close modal when clicking outside
departmentModal?.addEventListener('click', (e) => {
    if (e.target === departmentModal) {
        departmentModal.classList.add('hidden');
        resetForm();
    }
});

// Handle edit button clicks
document.querySelectorAll('.edit-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        const row = e.target.closest('tr');
        const id = row.dataset.id;
        const name = row.dataset.name;
        const isActive = row.dataset.isActive === 'true';
        
        // Populate form with existing data
        departmentIdInput.value = id;
        departmentNameInput.value = name;
        
        // Set status radio buttons
        if (isActive) {
            statusActiveInput.checked = true;
        } else {
            statusInactiveInput.checked = true;
        }
        
        // Update modal title and button
        modalTitle.textContent = 'Edit Department';
        submitBtn.textContent = 'Update Department';
        
        // Update form action for update
        departmentForm.action = `/departments/${id}`;
        departmentForm.method = 'POST';
        
        // Show modal
        departmentModal.classList.remove('hidden');
    });
});

// Handle delete button clicks
document.querySelectorAll('.delete-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        currentDeleteId = e.target.dataset.id;
        const departmentName = e.target.dataset.name;
        
        deleteDepartmentName.textContent = departmentName;
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
        form.action = `/departments/${currentDeleteId}/delete`;
        document.body.appendChild(form);
        form.submit();
    }
});

// Reset form to initial state
function resetForm() {
    departmentForm.reset();
    departmentIdInput.value = '';
    departmentNameInput.value = '';
    statusActiveInput.checked = true;
}

// Close modals on Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        if (!departmentModal.classList.contains('hidden')) {
            departmentModal.classList.add('hidden');
            resetForm();
        }
        if (!deleteModal.classList.contains('hidden')) {
            deleteModal.classList.add('hidden');
            currentDeleteId = null;
        }
    }
});