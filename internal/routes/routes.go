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

	// INITIALIZE CONTROLLER
	employeeController := controllers.NewUserController(userRepo)
	departmentController := controllers.NewDepartmentRepository(depRepo)
	roleController := controllers.NewRoleRepository(roleRepo)

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
		proj.GET("/", controllers.ListProjects)
		proj.GET("/create", controllers.CreateProjectForm)
		proj.GET("/edit", controllers.EditProjectForm)
		proj.GET("/detail", controllers.DetailProject)
		proj.GET("/board", controllers.BoardProject)
	}

	task := r.Group("/tasks")
	{
		task.GET("/", controllers.ListTask)
		task.GET("/create", controllers.CreateTaskForm)
		task.GET("/edit", controllers.EditTaskForm)
		task.GET("/detail", controllers.DetailTask)
		task.GET("/board", controllers.BoardTask)
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
}
