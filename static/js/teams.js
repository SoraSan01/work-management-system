document.addEventListener("DOMContentLoaded", () => {

	/* =====================
	   Elements
	===================== */
	const teamModal = document.getElementById("team-modal");
	const deleteModal = document.getElementById("delete-team-modal");
	const memberModal = document.getElementById("member-modal");

	const openModalBtn = document.getElementById("open-team-modal");
	const closeModalBtn = document.getElementById("close-team-modal");
	const closeMemberBtn = document.getElementById("close-member-modal");

	const teamForm = document.getElementById("team-form");
	const memberForm = document.getElementById("member-form");
	const modalTitle = document.getElementById("team-modal-title");

	const teamIdInput = document.getElementById("team-id");
	const teamNameInput = document.getElementById("team-name");
	const teamDescInput = document.getElementById("team-description");
	const memberTeamIdInput = document.getElementById("member-team-id");

	const deleteName = document.getElementById("delete-team-name");
	const confirmDeleteBtn = document.getElementById("confirm-team-delete");
	const cancelDeleteBtn = document.getElementById("cancel-team-delete");

	const teamCards = document.querySelectorAll('.team-card');
	const membersList = document.getElementById('members-list');
	const memberCountEl = document.getElementById('member-count');
	const addMemberBtn = document.getElementById('add-member-btn');
	
	// Store team data for client-side operations
	const teamsData = new Map();
	
	let deleteTeamId = null;
	let currentTeamId = null;

	/* =====================
	   Initialize Teams Data
	===================== */
	teamCards.forEach(card => {
		const teamId = card.dataset.id;
		teamsData.set(teamId, {
			id: teamId,
			name: card.dataset.name,
			description: card.dataset.description,
			memberCount: card.dataset.memberCount
		});
	});

	// Set first team as active on load
	if (teamCards.length > 0) {
		currentTeamId = teamCards[0].dataset.id;
		if (memberTeamIdInput) {
			memberTeamIdInput.value = currentTeamId;
		}
	}

	/* =====================
	   Modal Helpers
	===================== */
	const openModal = (modal) => {
		modal?.classList.remove("hidden");
		// Focus first input when modal opens
		setTimeout(() => {
			const firstInput = modal?.querySelector('input:not([type="hidden"]), select, textarea');
			firstInput?.focus();
		}, 100);
	};
	
	const closeModal = (modal) => modal?.classList.add("hidden");

	const resetTeamForm = () => {
		teamForm.reset();
		teamIdInput.value = "";
		teamForm.action = "/teams";
		modalTitle.textContent = "Create Team";
	};

	/* =====================
	   Update Members Display
	===================== */
	const updateMembersDisplay = (teamId) => {
		const teamData = teamsData.get(teamId);
		if (!teamData || !membersList) return;

		// Update count
		if (memberCountEl) {
			const count = teamData.memberCount;
			memberCountEl.textContent = `${count} member${count !== '1' ? 's' : ''}`;
		}

		// Show loading state
		membersList.innerHTML = '<div class="text-center py-12 text-xs text-gray-500">Loading members...</div>';

		// Fetch fresh member data from server
		fetch(`/teams/${teamId}/members`)
			.then(res => res.json())
			.then(data => {
				if (!data.members || data.members.length === 0) {
					membersList.innerHTML = '<div class="text-center py-12 text-xs text-gray-500">No members in this team</div>';
					return;
				}

				membersList.innerHTML = data.members.map(member => `
					<div class="flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:border-gray-300 transition-colors">
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 bg-gray-500 rounded-full flex items-center justify-center text-white font-semibold text-sm">
								${member.User.FirstName.charAt(0)}${member.User.LastName.charAt(0)}
							</div>
							<div>
								<div class="text-sm font-medium text-gray-900 flex items-center gap-1">
									${member.User.FirstName} ${member.User.LastName}
									${member.Role === 'leader' ? `
										<svg class="w-4 h-4 text-yellow-500" fill="currentColor" viewBox="0 0 20 20">
											<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
										</svg>
									` : ''}
								</div>
								<span class="text-xs text-gray-500">${member.User.Email}</span>
							</div>
						</div>
						<div class="flex items-center gap-3">
							<span class="px-2 py-0.5 text-xs rounded-md ${member.Role === 'leader' ? 'bg-yellow-50 text-yellow-700' : 'bg-gray-100 text-gray-700'}">
								${member.Role === 'leader' ? 'Leader' : 'Member'}
							</span>
							<button class="remove-member-btn text-red-600 hover:text-red-700 p-1" 
								data-member-id="${member.ID}"
								data-member-name="${member.User.FirstName} ${member.User.LastName}"
								data-team-id="${teamId}">
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
								</svg>
							</button>
						</div>
					</div>
				`).join('');

				// Re-attach event listeners to new remove buttons
				attachRemoveMemberListeners();
			})
			.catch(err => {
				console.error('Error fetching members:', err);
				membersList.innerHTML = '<div class="text-center py-12 text-xs text-red-500">Failed to load members</div>';
			});
	};

	/* =====================
	   Team Selection
	===================== */
	teamCards.forEach(card => {
		card.addEventListener('click', function(e) {
			// Don't trigger if clicking on action buttons
			if (e.target.closest('.edit-team-btn') || 
			    e.target.closest('.delete-team-btn') || 
			    e.target.closest('button')) {
				return;
			}

			// Remove active state from all cards
			teamCards.forEach(c => c.classList.remove('bg-gray-50', 'border-l-4', 'border-l-gray-900'));
			
			// Add active state to clicked card
			this.classList.add('bg-gray-50', 'border-l-4', 'border-l-gray-900');

			// Update current team ID
			currentTeamId = this.dataset.id;

			// Update hidden input for member modal
			if (memberTeamIdInput) {
				memberTeamIdInput.value = currentTeamId;
			}

			// Update members display
			updateMembersDisplay(currentTeamId);
		});
	});

	/* =====================
	   Create Team Modal
	===================== */
	openModalBtn?.addEventListener("click", () => {
		resetTeamForm();
		openModal(teamModal);
	});

	closeModalBtn?.addEventListener("click", () => {
		closeModal(teamModal);
	});

	// Close modal on backdrop click
	teamModal?.addEventListener('click', (e) => {
		if (e.target === teamModal) {
			closeModal(teamModal);
		}
	});

	/* =====================
	   Edit Team
	===================== */
	const attachEditListeners = () => {
		document.querySelectorAll(".edit-team-btn").forEach(btn => {
			btn.addEventListener("click", (e) => {
				e.stopPropagation();

				teamIdInput.value = btn.dataset.id;
				teamNameInput.value = btn.dataset.name;
				teamDescInput.value = btn.dataset.description;

				teamForm.action = `/teams/${btn.dataset.id}`;
				modalTitle.textContent = "Edit Team";

				openModal(teamModal);
			});
		});
	};
	attachEditListeners();

	/* =====================
	   Delete Team
	===================== */
	const attachDeleteListeners = () => {
		document.querySelectorAll(".delete-team-btn").forEach(btn => {
			btn.addEventListener("click", (e) => {
				e.stopPropagation();
				
				deleteTeamId = btn.dataset.id;
				const teamName = btn.dataset.name || teamsData.get(deleteTeamId)?.name || 'this team';
				
				if (deleteName) {
					deleteName.textContent = teamName;
				}
				
				openModal(deleteModal);
			});
		});
	};
	attachDeleteListeners();

	cancelDeleteBtn?.addEventListener("click", () => {
		closeModal(deleteModal);
		deleteTeamId = null;
	});

	// Close delete modal on backdrop click
	deleteModal?.addEventListener('click', (e) => {
		if (e.target === deleteModal) {
			closeModal(deleteModal);
			deleteTeamId = null;
		}
	});

	confirmDeleteBtn?.addEventListener("click", () => {
		if (!deleteTeamId) return;

		// Disable button to prevent double-clicks
		confirmDeleteBtn.disabled = true;
		confirmDeleteBtn.textContent = "Deleting...";

		fetch(`/teams/${deleteTeamId}/delete`, {
			method: "POST",
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
			alert("Failed to delete team. Please try again.");
			confirmDeleteBtn.disabled = false;
			confirmDeleteBtn.textContent = "Delete Team";
		});
	});

	/* =====================
	   Add Member Modal
	===================== */
	addMemberBtn?.addEventListener('click', (e) => {
		e.stopPropagation();
		
		if (!currentTeamId) {
			alert('Please select a team first');
			return;
		}

		memberTeamIdInput.value = currentTeamId;
		openModal(memberModal);
	});

	closeMemberBtn?.addEventListener('click', () => {
		closeModal(memberModal);
	});

	// Close member modal on backdrop click
	memberModal?.addEventListener('click', (e) => {
		if (e.target === memberModal) {
			closeModal(memberModal);
		}
	});

	/* =====================
	   Remove Member
	===================== */
	const attachRemoveMemberListeners = () => {
		document.querySelectorAll('.remove-member-btn').forEach(btn => {
			btn.addEventListener('click', (e) => {
				e.stopPropagation();
				
				const memberId = btn.dataset.memberId;
				const teamId = btn.dataset.teamId || currentTeamId;
				const memberName = btn.dataset.memberName;

				if (!confirm(`Remove ${memberName} from this team?`)) {
					return;
				}

				// Disable button
				btn.disabled = true;
				btn.innerHTML = '<svg class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>';
				
				fetch(`/teams/${teamId}/members/${memberId}`, {
					method: 'DELETE',
					headers: {
						'Content-Type': 'application/json',
					}
				})
				.then(res => {
					if (!res.ok) throw new Error('Remove member failed');
					location.reload();
				})
				.catch(err => {
					console.error(err);
					alert('Failed to remove member. Please try again.');
					btn.disabled = false;
					btn.innerHTML = '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>';
				});
			});
		});
	};
	attachRemoveMemberListeners();

	/* =====================
	   Form Validation
	===================== */
	teamForm?.addEventListener('submit', (e) => {
		const nameValue = teamNameInput.value.trim();
		
		if (!nameValue) {
			e.preventDefault();
			alert('Team name is required');
			teamNameInput.focus();
			return false;
		}

		if (nameValue.length < 3) {
			e.preventDefault();
			alert('Team name must be at least 3 characters long');
			teamNameInput.focus();
			return false;
		}
	});

	memberForm?.addEventListener('submit', (e) => {
		const userId = document.getElementById('member-user-id')?.value;
		const teamId = memberTeamIdInput?.value;
		
		if (!userId) {
			e.preventDefault();
			alert('Please select a user');
			return false;
		}

		if (!teamId) {
			e.preventDefault();
			alert('No team selected');
			return false;
		}
	});

	/* =====================
	   Keyboard Shortcuts
	===================== */
	document.addEventListener('keydown', (e) => {
		// ESC to close modals
		if (e.key === 'Escape') {
			if (!teamModal?.classList.contains('hidden')) {
				closeModal(teamModal);
			}
			if (!memberModal?.classList.contains('hidden')) {
				closeModal(memberModal);
			}
			if (!deleteModal?.classList.contains('hidden')) {
				closeModal(deleteModal);
				deleteTeamId = null;
			}
		}
	});

	/* =====================
	   Accessibility Improvements
	===================== */
	const trapFocus = (modal) => {
		const focusableElements = modal.querySelectorAll(
			'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])'
		);
		
		if (focusableElements.length === 0) return;
		
		const firstElement = focusableElements[0];
		const lastElement = focusableElements[focusableElements.length - 1];

		const handleTabKey = (e) => {
			if (e.key !== 'Tab') return;

			if (e.shiftKey) {
				if (document.activeElement === firstElement) {
					lastElement.focus();
					e.preventDefault();
				}
			} else {
				if (document.activeElement === lastElement) {
					firstElement.focus();
					e.preventDefault();
				}
			}
		};

		modal.addEventListener('keydown', handleTabKey);
	};

	if (teamModal) trapFocus(teamModal);
	if (memberModal) trapFocus(memberModal);
	if (deleteModal) trapFocus(deleteModal);

});