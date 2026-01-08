package controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"
	"workms/internal/models"
	"workms/internal/repositories"
	"workms/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	Repo     *repositories.ProjectRepository
	TeamRepo *repositories.TeamRepository
	UserRepo *repositories.UserRepository
}

func NewProjectController(
	projectRepo *repositories.ProjectRepository,
	teamRepo *repositories.TeamRepository,
	userRepo *repositories.UserRepository,
) *ProjectController {
	return &ProjectController{
		Repo:     projectRepo,
		TeamRepo: teamRepo,
		UserRepo: userRepo,
	}
}

// READ ALL DATA
func (pc *ProjectController) ListProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize := 10

	proj, total, err := pc.Repo.FindPaginated(page, pageSize)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	teams, err := pc.TeamRepo.FindAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1 // ensure at least 1 page to avoid template issues
	}

	// Calculate start and end for display
	var start, end int
	if total == 0 {
		start = 0
		end = 0
	} else {
		start = (page-1)*pageSize + 1
		end = page * pageSize
		if end > int(total) {
			end = int(total)
		}
	}

	// Build page numbers slice
	pages := make([]int, totalPages)
	for i := 0; i < totalPages; i++ {
		pages[i] = i + 1
	}

	c.HTML(http.StatusOK, "projects/list.html", gin.H{
		"title":    "Projects",
		"projects": proj,
		"teams":    teams,

		// pagination
		"currentPage": page,
		"totalPages":  totalPages,
		"total":       total,
		"start":       start,
		"end":         end,
		"pages":       pages,

		"hasPrev":  page > 1,
		"hasNext":  page < totalPages,
		"prevPage": page - 1,
		"nextPage": page + 1,
	})
}

func (pc *ProjectController) Store(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.PostForm("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.PostForm("end_date"))

	proj := models.Project{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Status:      c.PostForm("status"),
		StartDate:   startDate,
		EndDate:     endDate,
		TeamID:      utils.ParseUintPtr(c.PostForm("team_id")),
		CreatedBy:   utils.ParseUintPtr(c.PostForm("created_by")),
		CreatedAt:   time.Now(),
	}

	if err := pc.Repo.Create(&proj); err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Failed to create project: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/projects")
}

func (pc *ProjectController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid project ID",
		})
		return
	}

	// Fetch project
	proj, err := pc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "Project not found",
		})
		return
	}

	startDate, _ := time.Parse("2006-01-02", c.PostForm("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.PostForm("end_date"))

	// Update fields
	proj.Name = c.PostForm("name")
	proj.Description = c.PostForm("description")
	proj.Status = c.PostForm("status")
	proj.StartDate = startDate
	proj.EndDate = endDate
	proj.TeamID = utils.ParseUintPtr(c.PostForm("team_id"))

	// Save
	if err := pc.Repo.Update(proj); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to update project: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/projects")
}

// DELETE DATA
func (pc *ProjectController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	pc.Repo.Delete(id)
	c.Redirect(http.StatusFound, "/projects")
}

func (pc *ProjectController) GetUsersForProject(c *gin.Context) {
	projectID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	users, err := pc.UserRepo.GetUsersByProject(projectID) // implement this in UserRepo
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (pc *ProjectController) GetProjectUsers(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Find the project with its team
	project, err := pc.Repo.FindByID(projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// If project has no team, return empty array
	if project.TeamID == nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	// Get team members
	team, err := pc.TeamRepo.FindByIDWithMembers(*project.TeamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch team members"})
		return
	}

	// Return the users
	c.JSON(http.StatusOK, team.Members)
}
