// internal/routes/routes.go
package routes

import (
	"net/http"
	"workms/internal/controllers"
	"workms/internal/database"
	"workms/internal/middleware"
	"workms/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "errors/error.html", gin.H{})
	})

	// INITIALIZE REPO
	userRepo := repositories.NewUserRepository(database.DB)
	depRepo := repositories.NewDepartmentRepository(database.DB)
	roleRepo := repositories.NewRoleRepository(database.DB)
	projRepo := repositories.NewProjectRepository(database.DB)
	taskRepo := repositories.NewTaskRepository(database.DB)
	teamRepo := repositories.NewTeamRepository(database.DB)

	// INITIALIZE CONTROLLER
	employeeController := controllers.NewUserController(userRepo)
	departmentController := controllers.NewDepartmentRepository(depRepo)
	roleController := controllers.NewRoleRepository(roleRepo)
	projController := controllers.NewProjectController(projRepo)
	taskController := controllers.NewTaskController(taskRepo)
	teamController := controllers.NewTeamController(teamRepo)

	// Apply global middleware (example: simple auth)
	r.Use(middleware.Logger())

	// Public routes
	r.GET("/", controllers.IndexHandler)

	emp := r.Group("/employees")
	{
		emp.GET("/", employeeController.ListEmployees)
		emp.POST("/", employeeController.Store)
		emp.POST("/:id", employeeController.Update)
		emp.POST("/:id/delete", employeeController.Delete)
	}

	proj := r.Group("/projects")
	{
		proj.GET("/", projController.ListProjects)
		proj.POST("/", projController.Store)
		proj.POST("/:id", projController.Update)
		proj.POST("/:id/delete", projController.Delete)
	}

	task := r.Group("/tasks")
	{
		task.GET("/", taskController.ListTasks)
		task.POST("/", taskController.Store)
		task.POST("/:id", taskController.Update)
		task.POST("/:id/delete", taskController.Delete)
	}

	event := r.Group("/calendar")
	{
		event.GET("/", controllers.CalendarIndexHandler)
		event.GET("/event-create", controllers.CreateEventForm)
		event.GET("/event-details", controllers.DetailEvent)
	}

	docu := r.Group("/documents")
	{
		docu.GET("/", controllers.ListDocument)
		docu.GET("/upload", controllers.UploadDocument)
		docu.GET("/detail", controllers.DetailDocument)
	}

	reports := r.Group("/reports")
	{
		reports.GET("/tasks", controllers.ReportTask)
		reports.GET("/employees", controllers.ReportEmployee)
		reports.GET("/projects", controllers.ReportProject)
	}

	dep := r.Group("/departments")
	{
		dep.GET("/", departmentController.ListDepartments)
		dep.POST("/", departmentController.Store)
		dep.POST("/:id", departmentController.Update)
		dep.POST("/:id/delete", departmentController.Delete)
	}

	roles := r.Group("/roles")
	{
		roles.GET("/", roleController.ListRoles)
		roles.POST("/", roleController.Store)
		roles.POST("/:id", roleController.Update)
		roles.POST("/:id/delete", roleController.Delete)
	}

	teams := r.Group("/teams")
	{
		teams.GET("/", teamController.ListTeams)
		teams.POST("/", teamController.Store)
		teams.POST("/:id", teamController.Update)
		teams.POST("/:id/delete", teamController.Delete)
	}
}
