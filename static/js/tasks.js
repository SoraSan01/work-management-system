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
    const taskPrioritySelect = document.getElementById("task-priority");
    const taskDueDateInput = document.getElementById("task-due-date");

    const deleteTaskModal = document.getElementById("delete-task-modal");
    const deleteTaskTitleSpan = document.getElementById("delete-task-title");
    const cancelTaskDeleteBtn = document.getElementById("cancel-task-delete");
    const confirmTaskDeleteBtn = document.getElementById("confirm-task-delete");

    // Open Create Modal
    openTaskModalBtn.addEventListener("click", () => {
        taskModalTitle.textContent = "Create Task";
        taskForm.reset();
        taskIdInput.value = "";
        taskModal.classList.remove("hidden");
    });

    // Close Modal
    closeTaskModalBtn.addEventListener("click", () => {
        taskModal.classList.add("hidden");
    });

    // Edit Task Buttons
    document.querySelectorAll(".edit-task-btn").forEach((btn) => {
        btn.addEventListener("click", () => {
            const row = btn.closest("tr");
            taskModalTitle.textContent = "Edit Task";
            taskIdInput.value = row.dataset.id;
            taskTitleInput.value = row.dataset.title;
            taskDescriptionInput.value = row.dataset.description;
            taskStatusSelect.value = row.dataset.status;
            taskPrioritySelect.value = row.dataset.priority;
            taskDueDateInput.value = row.dataset.dueDate;

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
            const title = row.dataset.title.toLowerCase();
            const description = row.dataset.description.toLowerCase();
            const status = row.dataset.status;

            const matchesSearch = title.includes(searchText) || description.includes(searchText);
            const matchesStatus = !statusValue || status === statusValue;

            row.style.display = matchesSearch && matchesStatus ? "" : "none";
        });
    }

    searchInput.addEventListener("input", filterTasks);
    statusFilter.addEventListener("change", filterTasks);
});
