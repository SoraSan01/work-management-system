package controllers

import (
	"math"
	"net/http"
	"strconv"
	"workms/internal/models"
	"workms/internal/repositories"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	Repo *repositories.TeamRepository
}

func NewTeamController(repo *repositories.TeamRepository) *TeamController {
	return &TeamController{Repo: repo}
}

// ListTeams handles GET /teams - displays all teams with pagination
func (tc *TeamController) ListTeams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize := 10
	teams, total, err := tc.Repo.FindPaginated(page, pageSize)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch teams: " + err.Error(),
		})
		return
	}

	// Fetch all users for the Add Member dropdown
	users, err := tc.Repo.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch users: " + err.Error(),
		})
		return
	}

	// Calculate pagination data
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	start := (page-1)*pageSize + 1
	end := page * pageSize
	if end > int(total) {
		end = int(total)
	}
	if start > int(total) {
		start = int(total)
	}

	pages := make([]int, totalPages)
	for i := 0; i < totalPages; i++ {
		pages[i] = i + 1
	}

	c.HTML(http.StatusOK, "teams/list.html", gin.H{
		"title":       "Teams",
		"teams":       teams,
		"users":       users,
		"currentPage": page,
		"totalPages":  totalPages,
		"total":       total,
		"start":       start,
		"end":         end,
		"pages":       pages,
		"hasPrev":     page > 1,
		"hasNext":     page < totalPages,
		"prevPage":    page - 1,
		"nextPage":    page + 1,
	})
}

// Store handles POST /teams - creates a new team
func (tc *TeamController) Store(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")

	// Validation
	if name == "" {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Team name is required",
		})
		return
	}

	team := models.Team{
		Name:        name,
		Description: description,
	}

	if err := tc.Repo.Create(&team); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to create team: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/teams")
}

// Update handles POST /teams/:id - updates an existing team
func (tc *TeamController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid team ID",
		})
		return
	}

	// Fetch existing team
	team, err := tc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "Team not found: " + err.Error(),
		})
		return
	}

	// Get form data
	name := c.PostForm("name")
	description := c.PostForm("description")

	// Validation
	if name == "" {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Team name is required",
		})
		return
	}

	// Update fields
	team.Name = name
	team.Description = description

	// Save
	if err := tc.Repo.Update(team); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to update team: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/teams")
}

// Delete handles DELETE /teams/:id - deletes a team
func (tc *TeamController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid team ID",
		})
		return
	}

	if err := tc.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete team: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Team deleted successfully",
	})
}

// AddMember handles POST /teams/add-member - adds a member to a team
func (tc *TeamController) AddMember(c *gin.Context) {
	teamIDStr := c.PostForm("team_id")
	userIDStr := c.PostForm("user_id")
	role := c.PostForm("role")

	// Parse and validate team ID
	teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
	if err != nil || teamID == 0 {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid team ID",
		})
		return
	}

	// Parse and validate user ID
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid user ID",
		})
		return
	}

	// Validate role
	if role == "" {
		role = "member"
	}
	if role != "member" && role != "leader" {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid role. Must be 'member' or 'leader'",
		})
		return
	}

	// Create the team member
	member := models.TeamMember{
		TeamID: teamID,
		UserID: userID,
		Role:   role,
	}

	if err := tc.Repo.AddMember(&member); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to add member: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/teams")
}

// RemoveMember handles DELETE /teams/:teamId/members/:memberId - removes a member from a team
func (tc *TeamController) RemoveMember(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid team ID",
		})
		return
	}

	memberID, err := strconv.ParseUint(c.Param("memberId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid member ID",
		})
		return
	}

	if err := tc.Repo.RemoveMember(teamID, memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove member: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Member removed successfully",
	})
}

// ViewTeam handles GET /teams/:id - displays a single team's details
func (tc *TeamController) ViewTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid team ID",
		})
		return
	}

	team, err := tc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "Team not found: " + err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "teams/view.html", gin.H{
		"title": "Team Details",
		"team":  team,
	})
}

// GetTeamMembers handles GET /teams/:id/members - returns team members as JSON
func (tc *TeamController) GetTeamMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid team ID",
		})
		return
	}

	members, err := tc.Repo.GetMembersByTeamID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch members: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"members": members,
	})
}
