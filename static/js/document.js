// Modal Controls
const uploadModal = document.getElementById('uploadModal');
const uploadBtn = document.getElementById('uploadBtn');
const closeModal = document.getElementById('closeModal');
const cancelBtn = document.getElementById('cancelBtn');

uploadBtn.addEventListener('click', () => {
    uploadModal.classList.remove('hidden');
});

closeModal.addEventListener('click', () => {
    uploadModal.classList.add('hidden');
    resetForm();
});

cancelBtn.addEventListener('click', () => {
    uploadModal.classList.add('hidden');
    resetForm();
});

// Close modal on outside click
uploadModal.addEventListener('click', (e) => {
    if (e.target === uploadModal) {
        uploadModal.classList.add('hidden');
        resetForm();
    }
});

// File Upload Handling
const fileInput = document.getElementById('fileInput');
const dropZone = document.getElementById('dropZone');
const filePreview = document.getElementById('filePreview');
const fileName = document.getElementById('fileName');
const fileSize = document.getElementById('fileSize');
const removeFile = document.getElementById('removeFile');

let selectedFile = null;

// File input change
fileInput.addEventListener('change', (e) => {
    handleFiles(e.target.files);
});

// Drag and drop
dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropZone.classList.add('border-blue-500', 'bg-blue-50');
});

dropZone.addEventListener('dragleave', (e) => {
    e.preventDefault();
    dropZone.classList.remove('border-blue-500', 'bg-blue-50');
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('border-blue-500', 'bg-blue-50');
    handleFiles(e.dataTransfer.files);
});

// Handle selected files
function handleFiles(files) {
    if (files.length > 0) {
        const file = files[0];
        
        // Validate file size (10MB)
        if (file.size > 10 * 1024 * 1024) {
            alert('File size must be less than 10MB');
            return;
        }
        
        selectedFile = file;
        fileName.textContent = file.name;
        fileSize.textContent = formatFileSize(file.size);
        filePreview.classList.remove('hidden');
        dropZone.classList.add('hidden');
    }
}

// Remove selected file
removeFile.addEventListener('click', () => {
    selectedFile = null;
    fileInput.value = '';
    filePreview.classList.add('hidden');
    dropZone.classList.remove('hidden');
});

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

// Form submission
const uploadForm = document.getElementById('uploadForm');
uploadForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    if (!selectedFile) {
        alert('Please select a file');
        return;
    }
    
    const formData = new FormData(uploadForm);
    formData.set('file', selectedFile);
    
    try {
        const response = await fetch('/documents/upload', {
            method: 'POST',
            body: formData,
            credentials: 'same-origin'
        });
        
        if (response.ok) {
            const result = await response.json();
            alert('Document uploaded successfully!');
            uploadModal.classList.add('hidden');
            resetForm();
            // Reload page or update document list
            window.location.reload();
        } else {
            const error = await response.json();
            alert('Upload failed: ' + (error.message || 'Unknown error'));
        }
    } catch (error) {
        console.error('Upload error:', error);
        alert('Upload failed: ' + error.message);
    }
});

// Reset form
function resetForm() {
    uploadForm.reset();
    selectedFile = null;
    fileInput.value = '';
    filePreview.classList.add('hidden');
    dropZone.classList.remove('hidden');
}

// Search and Filter
const searchInput = document.getElementById('searchInput');
const projectFilter = document.getElementById('projectFilter');
const fileTypeFilter = document.getElementById('fileTypeFilter');
const sortBy = document.getElementById('sortBy');

function filterDocuments() {
    const searchTerm = searchInput.value.toLowerCase();
    const projectId = projectFilter.value;
    const fileType = fileTypeFilter.value;
    const cards = document.querySelectorAll('.document-card');
    
    cards.forEach(card => {
        const name = card.dataset.name.toLowerCase();
        const cardProject = card.dataset.project;
        const cardType = card.dataset.type;
        
        let show = true;
        
        if (searchTerm && !name.includes(searchTerm)) {
            show = false;
        }
        
        if (projectId && cardProject !== projectId) {
            show = false;
        }
        
        if (fileType) {
            if (fileType === 'doc' && !['doc', 'docx'].includes(cardType)) {
                show = false;
            } else if (fileType === 'image' && !['jpg', 'jpeg', 'png', 'gif'].includes(cardType)) {
                show = false;
            } else if (fileType === 'spreadsheet' && !['xls', 'xlsx', 'csv'].includes(cardType)) {
                show = false;
            } else if (fileType !== 'doc' && fileType !== 'image' && fileType !== 'spreadsheet' && cardType !== fileType) {
                show = false;
            }
        }
        
        card.style.display = show ? 'block' : 'none';
    });
}

searchInput.addEventListener('input', filterDocuments);
projectFilter.addEventListener('change', filterDocuments);
fileTypeFilter.addEventListener('change', filterDocuments);

// Sort functionality
sortBy.addEventListener('change', (e) => {
    const grid = document.getElementById('documentsGrid');
    const cards = Array.from(document.querySelectorAll('.document-card'));
    
    cards.sort((a, b) => {
        switch (e.target.value) {
            case 'newest':
                return b.dataset.created - a.dataset.created;
            case 'oldest':
                return a.dataset.created - b.dataset.created;
            case 'name':
                return a.dataset.name.localeCompare(b.dataset.name);
            case 'size':
                return (b.dataset.size || 0) - (a.dataset.size || 0);
            default:
                return 0;
        }
    });
    
    cards.forEach(card => grid.appendChild(card));
});

// Download document
async function downloadDocument(id) {
    try {
        const response = await fetch(`/documents/${id}/download`);
        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = response.headers.get('Content-Disposition').split('filename=')[1];
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        } else {
            alert('Download failed');
        }
    } catch (error) {
        console.error('Download error:', error);
        alert('Download failed: ' + error.message);
    }
}

// Delete document
async function deleteDocument(id) {
    if (!confirm('Are you sure you want to delete this document?')) {
        return;
    }
    
    try {
        const response = await fetch(`/documents/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            alert('Document deleted successfully');
            window.location.reload();
        } else {
            alert('Delete failed');
        }
    } catch (error) {
        console.error('Delete error:', error);
        alert('Delete failed: ' + error.message);
    }
}

async function approveDocument(id) {
    if (!confirm('Are you sure you want to approve this document?')) return;

    try {
        const res = await fetch(`/documents/${id}/approve`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });

        const data = await res.json();
        if (res.ok) {
            alert(data.message);
            location.reload();
        } else {
            alert('Error: ' + data.error);
        }
    } catch (err) {
        console.error(err);
        alert('Error approving document');
    }
}

