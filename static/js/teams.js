document.addEventListener("DOMContentLoaded", () => {

	/* =====================
	   Elements
	===================== */
	const teamModal = document.getElementById("team-modal");
	const deleteModal = document.getElementById("delete-team-modal");

	const openModalBtn = document.getElementById("open-team-modal");
	const closeModalBtn = document.getElementById("close-team-modal");

	const teamForm = document.getElementById("team-form");
	const modalTitle = document.getElementById("team-modal-title");

	const teamIdInput = document.getElementById("team-id");
	const teamNameInput = document.getElementById("team-name");
	const teamDescInput = document.getElementById("team-description");

	const searchInput = document.getElementById("search-input");

	const deleteName = document.getElementById("delete-team-name");
	const confirmDeleteBtn = document.getElementById("confirm-team-delete");
	const cancelDeleteBtn = document.getElementById("cancel-team-delete");

	let deleteTeamId = null;

	/* =====================
	   Helpers
	===================== */
	const openModal = () => teamModal.classList.remove("hidden");
	const closeModal = () => teamModal.classList.add("hidden");

	const openDeleteModal = () => deleteModal.classList.remove("hidden");
	const closeDeleteModal = () => deleteModal.classList.add("hidden");

	const resetForm = () => {
		teamForm.reset();
		teamIdInput.value = "";
		teamForm.action = "/teams";
		modalTitle.textContent = "Create Team";
	};

	/* =====================
	   Create Team
	===================== */
	openModalBtn?.addEventListener("click", () => {
		resetForm();
		openModal();
	});

	closeModalBtn?.addEventListener("click", () => {
		closeModal();
	});

	/* =====================
	   Edit Team
	===================== */
	document.querySelectorAll(".edit-team-btn").forEach(btn => {
		btn.addEventListener("click", (e) => {
			const row = e.target.closest("tr");

			teamIdInput.value = row.dataset.id;
			teamNameInput.value = row.dataset.name;
			teamDescInput.value = row.dataset.description;

			teamForm.action = `/teams/${row.dataset.id}`;
			modalTitle.textContent = "Edit Team";

			openModal();
		});
	});

	/* =====================
	   Delete Team
	===================== */
	document.querySelectorAll(".delete-team-btn").forEach(btn => {
		btn.addEventListener("click", () => {
			deleteTeamId = btn.dataset.id;
			deleteName.textContent = btn.dataset.name;
			openDeleteModal();
		});
	});

	cancelDeleteBtn?.addEventListener("click", closeDeleteModal);

	confirmDeleteBtn?.addEventListener("click", () => {
		if (!deleteTeamId) return;

		fetch(`/teams/${deleteTeamId}`, {
			method: "DELETE",
			headers: {
				"Content-Type": "application/json",
			}
		})
		.then(res => {
			if (!res.ok) throw new Error("Delete failed");
			location.reload();
		})
		.catch(err => {
			console.error(err);
			alert("Failed to delete team");
		});
	});

	/* =====================
	   Search Filter
	===================== */
	searchInput?.addEventListener("input", () => {
		const term = searchInput.value.toLowerCase();

		document.querySelectorAll("tbody tr").forEach(row => {
			const name = row.dataset.name?.toLowerCase() || "";
			const desc = row.dataset.description?.toLowerCase() || "";

			row.style.display =
				name.includes(term) || desc.includes(term)
					? ""
					: "none";
		});
	});

});
