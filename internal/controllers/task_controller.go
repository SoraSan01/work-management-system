package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
	"workms/internal/models"
	"workms/internal/repositories"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Repo        *repositories.TaskRepository
	ProjectRepo *repositories.ProjectRepository
	TeamRepo    *repositories.TeamRepository
	UserRepo    *repositories.UserRepository
}

func NewTaskController(
	repo *repositories.TaskRepository,
	projectRepo *repositories.ProjectRepository,
	teamRepo *repositories.TeamRepository,
	userRepo *repositories.UserRepository,
) *TaskController {
	return &TaskController{
		Repo:        repo,
		ProjectRepo: projectRepo,
		TeamRepo:    teamRepo,
		UserRepo:    userRepo,
	}
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

	projects, _ := tc.ProjectRepo.FindAll()
	users, _ := tc.UserRepo.GetAssignableUsers() // or UserRepo.FindAll()

	c.HTML(http.StatusOK, "tasks/list.html", gin.H{
		"title":    "Tasks",
		"tasks":    task,
		"projects": projects,
		"users":    users,

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

func (tc *TaskController) ListBoard(c *gin.Context) {
	// Fetch tasks with relations
	tasks, err := tc.Repo.FindAllWithRelations()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch tasks: " + err.Error(),
		})
		return
	}

	// Get task statistics
	totalTasks, inProgress, completed, err := tc.Repo.GetTaskStatistics()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch task statistics: " + err.Error(),
		})
		return
	}

	var completionRate int
	if totalTasks > 0 {
		completionRate = int((completed * 100) / totalTasks)
	}

	// Fetch projects with relations
	projects, err := tc.ProjectRepo.FindAllWithRelations()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch projects: " + err.Error(),
		})
		return
	}

	// Get project statistics
	totalProjects, activeProjects, err := tc.ProjectRepo.GetProjectStatistics()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch project statistics: " + err.Error(),
		})
		return
	}

	// Fetch teams
	teams, err := tc.TeamRepo.FindAllWithRelations()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch teams: " + err.Error(),
		})
		return
	}

	// Prepare enhanced tasks data
	enhancedTasks := make([]map[string]interface{}, len(tasks))
	for i, task := range tasks {
		projectName := ""
		if task.Project != nil {
			projectName = task.Project.Name
		}

		assignedTo := ""
		if task.User != nil {
			assignedTo = task.User.FirstName + " " + task.User.LastName
		}

		enhancedTasks[i] = map[string]interface{}{
			"ID":          task.ID,
			"Title":       task.Title,
			"Description": task.Description,
			"Status":      task.Status,
			"Priority":    task.Priority,
			"ProjectName": projectName,
			"AssignedTo":  assignedTo,
			"DueDate":     task.DueDate.Format("1/2/2006"),
		}
	}

	// Prepare enhanced projects data
	enhancedProjects := make([]map[string]interface{}, len(projects))
	for i, project := range projects {
		teamName := ""
		if project.Team != nil {
			teamName = project.Team.Name
		}

		createdBy := ""
		if project.User != nil {
			createdBy = project.User.FirstName + " " + project.User.LastName
		}

		// Get task counts for this project
		taskTotal, taskCompleted, _ := tc.ProjectRepo.GetTaskCountForProject(project.ID)

		timeline := fmt.Sprintf("%s - %s",
			project.StartDate.Format("1/2/2006"),
			project.EndDate.Format("1/2/2006"))

		enhancedProjects[i] = map[string]interface{}{
			"ID":           project.ID,
			"Name":         project.Name,
			"Description":  project.Description,
			"Status":       project.Status,
			"TeamName":     teamName,
			"CreatedBy":    createdBy,
			"TaskProgress": fmt.Sprintf("%d/%d completed", taskCompleted, taskTotal),
			"Timeline":     timeline,
		}
	}

	// Prepare enhanced teams data
	enhancedTeams := make([]map[string]interface{}, len(teams))
	for i, team := range teams {
		// Get project counts for this team
		projectTotal, projectActive, _ := tc.TeamRepo.GetProjectCountForTeam(team.ID)

		// Get task count for this team
		taskCount, _ := tc.TeamRepo.GetTaskCountForTeam(team.ID)

		enhancedTeams[i] = map[string]interface{}{
			"ID":             team.ID,
			"Name":           team.Name,
			"Description":    team.Description,
			"ProjectCount":   fmt.Sprintf("%d", projectTotal),
			"ActiveProjects": fmt.Sprintf("%d", projectActive),
			"TotalTasks":     fmt.Sprintf("%d tasks", taskCount),
		}
	}

	// Count unique team members
	teamMembersCount := len(teams) // Simplified; you might want to count unique members across all teams

	c.HTML(http.StatusOK, "task/board.html", gin.H{
		"stats": gin.H{
			"title":          "Task Board",
			"totalTasks":     totalTasks,
			"completionRate": completionRate,
			"inProgress":     inProgress,
			"completed":      completed,
			"activeProjects": activeProjects,
			"totalProjects":  totalProjects,
			"teamMembers":    teamMembersCount,
			"teams":          len(teams),
		},
		"tasks":    enhancedTasks,
		"projects": enhancedProjects,
		"teams":    enhancedTeams,
	})
}

func (tc *TaskController) Store(c *gin.Context) {
	var projectID *uint64
	if v := c.PostForm("project_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		projectID = &id
	}

	var assignedTo *uint64
	if v := c.PostForm("assigned_to"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		assignedTo = &id
	}

	task := models.Task{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Status:      c.PostForm("status"),
		Priority:    c.PostForm("priority"),
		ProjectID:   projectID,
		AssignedTo:  assignedTo,
	}

	// Parse due date if provided
	if dueDateStr := c.PostForm("due_date"); dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			task.DueDate = dueDate
		}
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

	// Fetch task
	task, err := tc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "Task not found",
		})
		return
	}

	task.Title = c.PostForm("title")
	task.Description = c.PostForm("description")
	task.Status = c.PostForm("status")
	task.Priority = c.PostForm("priority")

	if v := c.PostForm("project_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		task.ProjectID = &id
	} else {
		task.ProjectID = nil
	}

	if v := c.PostForm("assigned_to"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		task.AssignedTo = &id
	} else {
		task.AssignedTo = nil
	}

	// Parse due date if provided
	if dueDateStr := c.PostForm("due_date"); dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			task.DueDate = dueDate
		}
	}

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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := tc.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
