package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListEmployees(c *gin.Context) {
	c.HTML(http.StatusOK, "employees/list.html", gin.H{
		"title": "Employee List",
	})
}

func CreateEmployeeForm(c *gin.Context) {
	c.HTML(http.StatusOK, "employees/create.html", gin.H{
		"title": "Create Employee",
	})
}

func EditEmployeeForm(c *gin.Context) {
	c.HTML(http.StatusOK, "employees/edit.html", gin.H{
		"title": "Edit Employee",
	})
}

func DetailEmployee(c *gin.Context) {
	c.HTML(http.StatusOK, "employees/details.html", gin.H{
		"title": "Employee Details",
	})
}
