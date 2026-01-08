// tasks.js

document.addEventListener("DOMContentLoaded", () => {
    // Elements
    const taskModal = document.getElementById("task-modal");
    const openTaskModalBtn = document.getElementById("open-task-modal");
    const closeTaskModalBtn = document.getElementById("close-task-modal");
    const taskForm = document.getElementById("task-form");
    const taskModalTitle = document.getElementById("task-modal-title");
    const taskIdInput = document.getElementById("task-id");
    const taskTitleInput = document.getElementById("task-title");
    const taskDescriptionInput = document.getElementById("task-description");
    const taskStatusSelect = document.getElementById("task-status");
    const taskProjectSelect = document.getElementById("task-project");
    const taskAssignedToSelect = document.getElementById("task-assigned-to");
    const taskPrioritySelect = document.getElementById("task-priority");
    const taskDueDateInput = document.getElementById("task-due-date");

    const deleteTaskModal = document.getElementById("delete-task-modal");
    const deleteTaskTitleSpan = document.getElementById("delete-task-title");
    const cancelTaskDeleteBtn = document.getElementById("cancel-task-delete");
    const confirmTaskDeleteBtn = document.getElementById("confirm-task-delete");

    // Store original users list for create mode
    const originalUsersHTML = taskAssignedToSelect.innerHTML;

    // Open Create Modal
    openTaskModalBtn.addEventListener("click", () => {
        taskModalTitle.textContent = "Create Task";
        taskForm.reset();
        taskIdInput.value = "";
        taskForm.action = "/tasks";
        taskForm.method = "POST";
        
        // Reset to original users list
        taskAssignedToSelect.innerHTML = originalUsersHTML;
        
        taskModal.classList.remove("hidden");
        taskProjectSelect.value = "";
        taskAssignedToSelect.value = "";
    });

    // Close Modal
    closeTaskModalBtn.addEventListener("click", () => {
        taskModal.classList.add("hidden");
    });

    // Close modal when clicking outside
    taskModal.addEventListener("click", (e) => {
        if (e.target === taskModal) {
            taskModal.classList.add("hidden");
        }
    });

    // Edit Task Buttons
    document.querySelectorAll(".edit-task-btn").forEach((btn) => {
        btn.addEventListener("click", () => {
            const row = btn.closest("tr");
            taskModalTitle.textContent = "Edit Task";

            const taskId = row.dataset.id;
            taskIdInput.value = taskId;
            taskTitleInput.value = row.dataset.title || "";
            taskDescriptionInput.value = row.dataset.description || "";
            taskStatusSelect.value = row.dataset.status || "todo";
            taskPrioritySelect.value = row.dataset.priority || "medium";
            
            // Set due date - the data attribute should be in YYYY-MM-DD format
            const dueDate = row.dataset.dueDate;
            taskDueDateInput.value = dueDate || "";

            // Set project select
            const projectId = row.dataset.projectId || "";
            taskProjectSelect.value = projectId;

            // Get assigned user ID
            const assignedTo = row.dataset.assignedTo || "";

            // If there's a project, load its team members
            if (projectId) {
                updateAssignedUsers(projectId, assignedTo);
            } else {
                // No project, use all users
                taskAssignedToSelect.innerHTML = originalUsersHTML;
                taskAssignedToSelect.value = assignedTo;
            }

            // Update form action and method for editing
            taskForm.action = `/tasks/${taskId}`;
            taskForm.method = "POST";

            taskModal.classList.remove("hidden");
        });
    });

    // Delete Task Buttons
    document.querySelectorAll(".delete-task-btn").forEach((btn) => {
        btn.addEventListener("click", () => {
            const row = btn.closest("tr");
            deleteTaskTitleSpan.textContent = row.dataset.title;
            confirmTaskDeleteBtn.dataset.id = row.dataset.id;
            deleteTaskModal.classList.remove("hidden");
        });
    });

    // Cancel Delete
    cancelTaskDeleteBtn.addEventListener("click", () => {
        deleteTaskModal.classList.add("hidden");
    });

    // Close delete modal when clicking outside
    deleteTaskModal.addEventListener("click", (e) => {
        if (e.target === deleteTaskModal) {
            deleteTaskModal.classList.add("hidden");
        }
    });

    // Confirm Delete
    confirmTaskDeleteBtn.addEventListener("click", () => {
        const taskId = confirmTaskDeleteBtn.dataset.id;
        fetch(`/tasks/${taskId}`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
        })
        .then((res) => {
            if (res.ok) {
                window.location.reload();
            } else {
                alert("Failed to delete task.");
            }
        })
        .catch((err) => {
            console.error("Error deleting task:", err);
            alert("Failed to delete task.");
        });
    });

    // Search & Filter (basic client-side filtering)
    const searchInput = document.getElementById("search-input");
    const statusFilter = document.getElementById("status-filter");

    const tableRows = Array.from(document.querySelectorAll("tbody tr"));

    function filterTasks() {
        const searchText = searchInput.value.toLowerCase();
        const statusValue = statusFilter.value;

        tableRows.forEach((row) => {
            const title = (row.dataset.title || "").toLowerCase();
            const description = (row.dataset.description || "").toLowerCase();
            const status = row.dataset.status;

            const matchesSearch = title.includes(searchText) || description.includes(searchText);
            const matchesStatus = !statusValue || status === statusValue;

            row.style.display = matchesSearch && matchesStatus ? "" : "none";
        });
    }

    searchInput.addEventListener("input", filterTasks);
    statusFilter.addEventListener("change", filterTasks);

    // Update assigned users when project select changes
    taskProjectSelect.addEventListener("change", () => {
        const projectId = taskProjectSelect.value;
        if (projectId) {
            updateAssignedUsers(projectId);
        } else {
            // No project selected, restore all users
            taskAssignedToSelect.innerHTML = originalUsersHTML;
        }
    });

    // Function to update assigned users based on project's team
    function updateAssignedUsers(projectId, selectedUserId = "") {
        if (!projectId) {
            taskAssignedToSelect.innerHTML = originalUsersHTML;
            return;
        }

        // Show loading state
        taskAssignedToSelect.innerHTML = '<option value="">Loading...</option>';
        taskAssignedToSelect.disabled = true;

        // Fetch users for the project's team via API
        fetch(`/projects/${projectId}/users`)
            .then(res => {
                if (!res.ok) {
                    throw new Error('Failed to fetch users');
                }
                return res.json();
            })
            .then(users => {
                // Clear and rebuild options
                taskAssignedToSelect.innerHTML = '<option value="">Unassigned</option>';
                
                if (users && users.length > 0) {
                    users.forEach(user => {
                        const option = document.createElement("option");
                        option.value = user.ID;
                        option.textContent = `${user.FirstName} ${user.LastName}`;
                        if (selectedUserId && selectedUserId == user.ID) {
                            option.selected = true;
                        }
                        taskAssignedToSelect.appendChild(option);
                    });
                } else {
                    // No team members found
                    const option = document.createElement("option");
                    option.value = "";
                    option.textContent = "No team members in this project";
                    option.disabled = true;
                    taskAssignedToSelect.appendChild(option);
                }
                
                taskAssignedToSelect.disabled = false;
            })
            .catch(err => {
                console.error("Error fetching project users:", err);
                // Restore all users on error
                taskAssignedToSelect.innerHTML = originalUsersHTML;
                taskAssignedToSelect.disabled = false;
                if (selectedUserId) {
                    taskAssignedToSelect.value = selectedUserId;
                }
            });
    }
});