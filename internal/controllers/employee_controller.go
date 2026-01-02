package controllers

import (
	"math"
	"net/http"
	"strconv"
	"workms/internal/models"
	"workms/internal/repositories"
	"workms/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Repo *repositories.UserRepository
}

func NewUserController(repo *repositories.UserRepository) *UserController {
	return &UserController{Repo: repo}
}

// READ ALL DATA
func (uc *UserController) ListEmployees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize := 10

	users, total, err := uc.Repo.FindPaginated(page, pageSize)
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

	// Load departments and roles
	departments, err := uc.Repo.GetActiveDepartments()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	roles, err := uc.Repo.FindAllRole()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "employees/list.html", gin.H{
		"title":       "Employees",
		"users":       users,
		"departments": departments,
		"roles":       roles,

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

func (uc *UserController) Store(c *gin.Context) {
	password, err := bcrypt.GenerateFromPassword(
		[]byte(c.PostForm("password")),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to encrypt password: " + err.Error(),
		})
		return
	}

	deptID := utils.Uint64OrZero(c.PostForm("department_id"))
	roleID := utils.Uint64OrZero(c.PostForm("role_id"))

	user := models.User{
		FirstName:    c.PostForm("first_name"),
		LastName:     c.PostForm("last_name"),
		Email:        c.PostForm("email"),
		Password:     string(password),
		RoleID:       roleID,
		IsActive:     c.PostForm("is_active") == "true",
		DepartmentID: deptID,
	}

	if err := uc.Repo.Create(&user); err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Failed to create user: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/employees")
}

func (uc *UserController) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// Fetch user
	user, err := uc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "User not found",
		})
		return
	}

	// Update basic fields
	user.FirstName = c.PostForm("first_name")
	user.LastName = c.PostForm("last_name")
	user.Email = c.PostForm("email")
	user.IsActive = c.PostForm("is_active") == "true"

	// Update RoleID and DepartmentID
	user.RoleID = utils.Uint64OrZero(c.PostForm("role_id"))
	user.DepartmentID = utils.Uint64OrZero(c.PostForm("department_id"))

	// Update password if provided
	password := c.PostForm("password")
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
				"Message": "Failed to hash password: " + err.Error(),
			})
			return
		}
		user.Password = string(hashed)
	}

	// Save
	if err := uc.Repo.Update(user); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to update user: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/employees")
}

// DELETE DATA
func (uc *UserController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	uc.Repo.Delete(id)
	c.Redirect(http.StatusFound, "/employees")
}
