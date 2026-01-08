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
	// Global middleware
	r.Use(middleware.Logger())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	depRepo := repositories.NewDepartmentRepository(database.DB)
	roleRepo := repositories.NewRoleRepository(database.DB)
	projRepo := repositories.NewProjectRepository(database.DB)
	taskRepo := repositories.NewTaskRepository(database.DB)
	teamRepo := repositories.NewTeamRepository(database.DB)
	docuRepo := repositories.NewDocumentRepository(database.DB)
	authRepo := repositories.NewAuthenticationRepository(database.DB)

	// Initialize controllers
	userController := controllers.NewUserController(userRepo)
	departmentController := controllers.NewDepartmentController(depRepo)
	roleController := controllers.NewRoleController(roleRepo)
	projController := controllers.NewProjectController(projRepo, teamRepo, userRepo)
	taskController := controllers.NewTaskController(taskRepo, projRepo, teamRepo, userRepo)
	teamController := controllers.NewTeamController(teamRepo)
	docuController := controllers.NewDocumentController(docuRepo, projRepo, taskRepo, userRepo)
	authController := controllers.NewAuthenticationController(authRepo)
	dashboardController := controllers.NewDashboardController(userRepo)

	// --------------------
	// Public Routes
	// --------------------
	r.GET("/", authController.Login)
	r.GET("/login", authController.Login) // ✅ ADD THIS
	r.POST("/login", authController.Authenticate)
	r.GET("/logout", middleware.AuthRequired(), authController.Logout)

	// Dashboard
	r.GET("/dashboard", middleware.AuthRequired(), dashboardController.Index)

	// --------------------
	// Employee Routes
	// --------------------
	emp := r.Group("/employees")
	{
		emp.GET("/", userController.ListEmployees)
		emp.POST("/", userController.Store)
		emp.POST("/:id", userController.Update)
		emp.POST("/:id/delete", userController.Delete)
	}

	// --------------------
	// Department Routes
	// --------------------
	dep := r.Group("/departments")
	{
		dep.GET("/", departmentController.ListDepartments)
		dep.POST("/", departmentController.Store)
		dep.POST("/:id", departmentController.Update)
		dep.POST("/:id/delete", departmentController.Delete)
	}

	// --------------------
	// Role Routes
	// --------------------
	roles := r.Group("/roles")
	{
		roles.GET("/", roleController.ListRoles)
		roles.POST("/", roleController.Store)
		roles.POST("/:id", roleController.Update)
		roles.POST("/:id/delete", roleController.Delete)
	}

	// --------------------
	// Project Routes
	// --------------------
	proj := r.Group("/projects")
	{
		proj.GET("/", projController.ListProjects)
		proj.POST("/", projController.Store)
		proj.POST("/:id", projController.Update)
		proj.POST("/:id/delete", projController.Delete)
		proj.GET("/:id/users", projController.GetProjectUsers)
	}

	// --------------------
	// Task Routes
	// --------------------
	task := r.Group("/tasks")
	{
		task.GET("/", taskController.ListTasks)
		task.POST("/", taskController.Store)
		task.POST("/:id", taskController.Update)
		task.DELETE("/:id", taskController.Delete)
		task.GET("/board", taskController.ListBoard)
	}

	// --------------------
	// Document Routes
	// --------------------
	docs := r.Group("/documents")
	{
		docs.GET("/", docuController.ListDocuments)
		docs.POST("/upload", docuController.UploadDocument)
		docs.GET("/:id", docuController.GetDocument)
		docs.GET("/:id/download", docuController.DownloadDocument)
		docs.PUT("/:id", docuController.UpdateDocument)
		docs.DELETE("/:id", docuController.DeleteDocument)
		docs.POST("/documents/:id/approve", middleware.AuthAdminOrManager(), docuController.ApproveDocument)
	}

	// --------------------
	// Other Routes (Teams, Calendar, Reports)
	// --------------------
	teams := r.Group("/teams")
	{
		teams.GET("/", teamController.ListTeams)
		teams.POST("/", teamController.Store)
		teams.POST("/:id", teamController.Update)
		teams.POST("/:id/delete", teamController.Delete)
		teams.POST("/add-member", teamController.AddMember)
		teams.DELETE("/:teamId/members/:memberId", teamController.RemoveMember)
		teams.GET("/:id/members", teamController.GetTeamMembers)
	}

	event := r.Group("/calendar")
	{
		event.GET("/", controllers.CalendarIndexHandler)
		event.GET("/event-create", controllers.CreateEventForm)
		event.GET("/event-details", controllers.DetailEvent)
	}

	reports := r.Group("/reports")
	{
		reports.GET("/tasks", controllers.ReportTask)
		reports.GET("/employees", controllers.ReportEmployee)
		reports.GET("/projects", controllers.ReportProject)
	}

	// 404 fallback
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{})
	})
}
