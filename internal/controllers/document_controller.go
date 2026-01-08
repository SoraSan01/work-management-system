package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"workms/internal/models"
	"workms/internal/repositories"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentController struct {
	Repo        *repositories.DocumentRepository
	ProjectRepo *repositories.ProjectRepository
	TaskRepo    *repositories.TaskRepository
	UserRepo    *repositories.UserRepository
}

func NewDocumentController(
	repo *repositories.DocumentRepository,
	projectRepo *repositories.ProjectRepository,
	taskRepo *repositories.TaskRepository,
	userRepo *repositories.UserRepository,
) *DocumentController {
	return &DocumentController{
		Repo:        repo,
		ProjectRepo: projectRepo,
		TaskRepo:    taskRepo,
		UserRepo:    userRepo,
	}
}

// For rendering documents in templates
type DocumentView struct {
	ID            uint64
	Name          string
	OriginalName  string
	FilePath      string
	FileType      string
	FileSize      string // formatted string for display (e.g., "1.2 MB")
	FileSizeBytes int64  // raw bytes for JS sorting
	MimeType      string
	Description   string
	Version       int
	ProjectID     *uint64
	TaskID        *uint64
	Project       *models.Project
	Task          *models.Task
	UserFullName  string
	IsApproved    bool       // ADDED: Missing field
	ApprovedBy    *uint64    // ADDED: Missing field
	ApprovedAt    *time.Time // ADDED: Missing field
	CreatedAt     time.Time
}

// ListDocuments renders the document management page
func (dc *DocumentController) ListDocuments(c *gin.Context) {
	// Fetch all documents
	documents, err := dc.Repo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch documents: " + err.Error(),
		})
		return
	}

	// Debug: Log the number of documents fetched
	fmt.Printf("Found %d documents from database\n", len(documents))

	// Get statistics
	total, totalSize, recentUploads, fileTypes, err := dc.Repo.GetDocumentStatistics()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to fetch statistics: " + err.Error(),
		})
		return
	}

	// Fetch projects and tasks for filters
	projects, _ := dc.ProjectRepo.FindAll()
	tasks, _ := dc.TaskRepo.FindAll()

	// Get current user
	session := sessions.Default(c)
	var currentUser *models.User
	if uid := session.Get("user_id"); uid != nil {
		id := uid.(uint64)
		currentUser, _ = dc.UserRepo.FindByID(id)
	}

	// Enhance documents data
	enhancedDocuments := make([]DocumentView, len(documents))
	for i, doc := range documents {
		userFullName := ""
		if doc.User != nil {
			userFullName = doc.User.FirstName + " " + doc.User.LastName
		}

		var project *models.Project
		if doc.Project != nil {
			project = doc.Project
		}

		var task *models.Task
		if doc.Task != nil {
			task = doc.Task
		}

		enhancedDocuments[i] = DocumentView{
			ID:            doc.ID,
			Name:          doc.Name,
			OriginalName:  doc.OriginalName,
			FilePath:      doc.FilePath,
			FileType:      strings.ToLower(doc.FileType),
			FileSize:      formatFileSize(doc.FileSize),
			FileSizeBytes: doc.FileSize,
			MimeType:      doc.MimeType,
			Description:   doc.Description,
			Version:       doc.Version,
			ProjectID:     doc.ProjectID,
			TaskID:        doc.TaskID,
			Project:       project,
			Task:          task,
			UserFullName:  userFullName,
			IsApproved:    doc.IsApproved, // ADDED
			ApprovedBy:    doc.ApprovedBy, // ADDED
			ApprovedAt:    doc.ApprovedAt, // ADDED
			CreatedAt:     doc.CreatedAt,
		}
	}

	// Debug output
	fmt.Printf("Enhanced documents count: %d\n", len(enhancedDocuments))
	if len(enhancedDocuments) > 0 {
		fmt.Printf("First document: %+v\n", enhancedDocuments[0])
	}

	// Pass to template
	c.HTML(http.StatusOK, "tasks/document.html", gin.H{
		"title": "Documents",
		"stats": gin.H{
			"totalDocuments": total,
			"totalSize":      formatFileSize(totalSize),
			"recentUploads":  recentUploads,
			"fileTypes":      fileTypes,
		},
		"CurrentUser": currentUser,
		"documents":   enhancedDocuments,
		"projects":    projects,
		"tasks":       tasks,
	})
}

// UploadDocument handles document upload
func (dc *DocumentController) UploadDocument(c *gin.Context) {
	session := sessions.Default(c)
	userIDVal := session.Get("user_id")
	if userIDVal == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uploadedBy := userIDVal.(uint64)

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file size (10MB max)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads/documents"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Get form data
	name := c.PostForm("name")
	if name == "" {
		name = file.Filename
	}
	description := c.PostForm("description")

	// Parse optional foreign keys
	var projectID *uint64
	if v := c.PostForm("project_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		projectID = &id
	}

	var taskID *uint64
	if v := c.PostForm("task_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		taskID = &id
	}

	// Determine file type
	fileType := strings.TrimPrefix(ext, ".")

	// Get MIME type
	mimeType := file.Header.Get("Content-Type")

	// Create document record
	document := models.Document{
		Name:         name,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileType:     fileType,
		FileSize:     file.Size,
		MimeType:     mimeType,
		Description:  description,
		UploadedBy:   &uploadedBy,
		ProjectID:    projectID,
		TaskID:       taskID,
		Version:      1,
		IsDeleted:    false,
	}

	// Save to database
	if err := dc.Repo.Create(&document); err != nil {
		// Clean up uploaded file if database save fails
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Document uploaded successfully",
		"document": document,
	})
}

// DownloadDocument handles document download
func (dc *DocumentController) DownloadDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Fetch document
	document, err := dc.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on server"})
		return
	}

	// Open file
	file, err := os.Open(document.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// Set headers for download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", document.OriginalName))
	c.Header("Content-Type", document.MimeType)

	// Stream file to response
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}
}

// DeleteDocument handles document deletion
func (dc *DocumentController) DeleteDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Fetch document to get file path
	if _, err := dc.Repo.FindByID(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Soft delete in database
	if err := dc.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

// GetDocument returns document details as JSON
func (dc *DocumentController) GetDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	document, err := dc.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, document)
}

// UpdateDocument handles document metadata updates
func (dc *DocumentController) UpdateDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	document, err := dc.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Update fields
	if name := c.PostForm("name"); name != "" {
		document.Name = name
	}

	if description := c.PostForm("description"); description != "" {
		document.Description = description
	}

	if v := c.PostForm("project_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		document.ProjectID = &id
	}

	if v := c.PostForm("task_id"); v != "" {
		id, _ := strconv.ParseUint(v, 10, 64)
		document.TaskID = &id
	}

	document.UpdatedAt = time.Now()

	// Save
	if err := dc.Repo.Update(document); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Document updated successfully",
		"document": document,
	})
}

func (dc *DocumentController) ApproveDocument(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id").(uint64)

	docID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := dc.Repo.FindByID(docID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	if doc.IsApproved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document already approved"})
		return
	}

	now := time.Now()
	doc.IsApproved = true
	doc.ApprovedBy = &userID
	doc.ApprovedAt = &now

	if err := dc.Repo.Update(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve document"})
		return
	}

	// Check if all docs for task are approved
	if doc.TaskID != nil {
		allDocs, _ := dc.Repo.FindByTaskID(*doc.TaskID)
		allApproved := true
		for _, d := range allDocs {
			if !d.IsApproved {
				allApproved = false
				break
			}
		}
		if allApproved {
			task, _ := dc.TaskRepo.FindByID(*doc.TaskID)
			task.Status = "done"
			dc.TaskRepo.Update(task)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document approved successfully"})
}

// Helper function to format file size
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
