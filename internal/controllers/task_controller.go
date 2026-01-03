package controllers

import (
	"math"
	"net/http"
	"strconv"
	"workms/internal/models"
	"workms/internal/repositories"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Repo *repositories.TaskRepository
}

func NewTaskController(repo *repositories.TaskRepository) *TaskController {
	return &TaskController{Repo: repo}
}

// READ ALL DATA
func (tc *TaskController) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize := 10

	task, total, err := tc.Repo.FindPaginated(page, pageSize)
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

	c.HTML(http.StatusOK, "tasks/list.html", gin.H{
		"title": "Tasks",
		"tasks": task,

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

func (tc *TaskController) Store(c *gin.Context) {

	task := models.Task{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Status:      c.PostForm("status"),
		Priority:    c.PostForm("priority"),
	}

	if err := tc.Repo.Create(&task); err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Failed to create task: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/tasks")
}

func (tc *TaskController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Invalid task ID",
		})
		return
	}

	// Fetch project
	task, err := tc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "Task not found",
		})
		return
	}

	// Update fields
	task.Title = c.PostForm("title")

	// Save
	if err := tc.Repo.Update(task); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to update task: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/tasks")
}

// DELETE DATA
func (tc *TaskController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tc.Repo.Delete(id)
	c.Redirect(http.StatusFound, "/tasks")
}
