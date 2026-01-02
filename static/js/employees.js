// employees.js
document.addEventListener("DOMContentLoaded", () => {
	const modal = document.getElementById("employee-modal");
	const deleteModal = document.getElementById("delete-modal");
	const openBtn = document.getElementById("open-modal");
	const closeBtn = document.getElementById("close-modal");
	const form = document.getElementById("employee-form");
	const modalTitle = document.getElementById("modal-title");
	const submitBtn = document.getElementById("submit-btn");
	const passwordInput = document.getElementById("password");
	const passwordHint = document.getElementById("password-hint");

	let deleteEmployeeId = null;

	// --- Open modal for creating new employee ---
	openBtn.addEventListener("click", () => {
		resetForm();
		modalTitle.textContent = "Create New Employee";
		submitBtn.textContent = "Create";
		passwordInput.required = true;
		passwordInput.placeholder = "Password";
		passwordHint.textContent = "Required for new employees";
		form.action = "/employees";
		modal.classList.remove("hidden");
	});

	// --- Close modal ---
	closeBtn.addEventListener("click", () => {
		modal.classList.add("hidden");
		resetForm();
	});

	modal.addEventListener("click", (e) => {
		if (e.target === modal) {
			modal.classList.add("hidden");
			resetForm();
		}
	});

	function resetForm() {
		form.reset();
		document.getElementById("employee-id").value = "";
	}

	// --- Edit buttons ---
	const editButtons = document.querySelectorAll(".edit-btn");
	editButtons.forEach(btn => {
		btn.addEventListener("click", (e) => {
			const row = e.target.closest("tr");
			const id = row.dataset.id;
			const firstName = row.dataset.firstName;
			const lastName = row.dataset.lastName;
			const email = row.dataset.email;
			const departmentId = row.dataset.departmentId;
			const roleId = row.dataset.roleId;

			document.getElementById("employee-id").value = id;
			document.getElementById("first-name").value = firstName;
			document.getElementById("last-name").value = lastName;
			document.getElementById("email").value = email;
			document.getElementById("department-id").value = departmentId;
			document.getElementById("role-id").value = roleId;

			modalTitle.textContent = "Update Employee";
			submitBtn.textContent = "Update";
			passwordInput.required = false;
			passwordInput.placeholder = "Password (leave blank to keep current)";
			passwordHint.textContent = "Leave blank to keep current password";
			form.action = `/employees/${id}`;

			modal.classList.remove("hidden");
		});
	});

	// --- Delete buttons ---
	const deleteButtons = document.querySelectorAll(".delete-btn");
	deleteButtons.forEach(btn => {
		btn.addEventListener("click", (e) => {
			deleteEmployeeId = e.target.dataset.id;
			const employeeName = e.target.dataset.name;
			document.getElementById("delete-employee-name").textContent = employeeName;
			deleteModal.classList.remove("hidden");
		});
	});

	document.getElementById("cancel-delete").addEventListener("click", () => {
		deleteModal.classList.add("hidden");
		deleteEmployeeId = null;
	});

	document.getElementById("confirm-delete").addEventListener("click", () => {
		if (deleteEmployeeId) {
			const deleteForm = document.createElement("form");
			deleteForm.method = "POST";
			deleteForm.action = `/employees/${deleteEmployeeId}/delete`;

			const methodInput = document.createElement("input");
			methodInput.type = "hidden";
			methodInput.name = "_method";
			methodInput.value = "DELETE";
			deleteForm.appendChild(methodInput);

			document.body.appendChild(deleteForm);
			deleteForm.submit();
		}
	});

	// --- Search functionality ---
	const searchInput = document.getElementById("search-input");
	const departmentFilter = document.getElementById("department-filter");

	function filterTable() {
		const searchTerm = searchInput.value.toLowerCase();
		const selectedDept = departmentFilter.value;
		const rows = document.querySelectorAll("tbody tr[data-id]");

		rows.forEach(row => {
			const firstName = row.dataset.firstName.toLowerCase();
			const lastName = row.dataset.lastName.toLowerCase();
			const email = row.dataset.email.toLowerCase();
			const deptId = row.dataset.departmentId;

			const matchesSearch = firstName.includes(searchTerm) || lastName.includes(searchTerm) || email.includes(searchTerm);
			const matchesDept = !selectedDept || deptId === selectedDept;

			row.style.display = (matchesSearch && matchesDept) ? "" : "none";
		});
	}

	searchInput.addEventListener("input", filterTable);
	departmentFilter.addEventListener("change", filterTable);
});