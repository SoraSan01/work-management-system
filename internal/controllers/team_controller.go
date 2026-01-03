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

// READ ALL DATA
func (tc *TeamController) ListTeams(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize := 10

	team, total, err := tc.Repo.FindPaginated(page, pageSize)
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

	c.HTML(http.StatusOK, "teams/list.html", gin.H{
		"title": "Teams",
		"teams": team,

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

func (tc *TeamController) Store(c *gin.Context) {

	team := models.Team{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
	}

	if err := tc.Repo.Create(&team); err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Failed to create team: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/teams")
}

func (tc *TeamController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid team ID",
		})
		return
	}

	// Fetch project
	team, err := tc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "Team not found",
		})
		return
	}

	// Update fields
	team.Name = c.PostForm("name")
	team.Description = c.PostForm("description")

	// Save
	if err := tc.Repo.Update(team); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to update team: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/teams")
}

// DELETE DATA
func (tc *TeamController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tc.Repo.Delete(id)
	c.Redirect(http.StatusFound, "/teams")
}
