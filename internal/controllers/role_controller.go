package controllers

import (
	"net/http"
	"strconv"
	"workms/internal/models"
	"workms/internal/repositories"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	Repo *repositories.RoleRepository
}

func NewRoleController(repo *repositories.RoleRepository) *RoleController {
	return &RoleController{Repo: repo}
}

// READ ALL DATA
func (uc *RoleController) ListRoles(c *gin.Context) {
	role, err := uc.Repo.FindAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "roles/list.html", gin.H{
		"roles": role,
	})
}

// INSERT DATA
func (uc *RoleController) Store(c *gin.Context) {
	role := models.Role{
		Name: c.PostForm("name"),
	}

	if err := uc.Repo.Create(&role); err != nil {
		c.HTML(http.StatusBadRequest, "errors/error.html", gin.H{
			"Message": "Failed to create Role: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/roles")
}

func (uc *RoleController) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	// Fetch Role
	role, err := uc.Repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
			"Message": "Role not found",
		})
		return
	}

	role.Name = c.PostForm("name")

	// Save
	if err := uc.Repo.Update(role); err != nil {
		c.HTML(http.StatusInternalServerError, "errors/error.html", gin.H{
			"Message": "Failed to update Role: " + err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/roles")
}

// DELETE DATA
func (uc *RoleController) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	uc.Repo.Delete(id)
	c.Redirect(http.StatusFound, "/roles")
}
