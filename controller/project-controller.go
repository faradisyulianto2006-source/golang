	package controller 

	// Importing necessary packages
	import (
		"github.com/gin-gonic/gin"
		"gitlab.com/pragmaticreviews/golang-gin-poc/entity"
		"gitlab.com/pragmaticreviews/golang-gin-poc/service"
	)

	// Importing necessary packages and defining the project controller interface
	type ProjectController interface {
		FindAll() []entity.Project
		Save(ctx *gin.Context) entity.Project
		Delete(id int)
		Update(id int, updateData entity.Project) entity.Project
		FindById(id int) (entity.Project, bool)
	}

	// implementation of the project controller
	type projectController struct {
		service service.ProjectService
	}

	// constructor function to create a new project controller instance
	func NewProjectController(service service.ProjectService) ProjectController {
		return &projectController {
			service: service,
		}
	}

	// FindAll method to return all projects from the controller
	func (c *projectController) FindAll() []entity.Project {
		return c.service.FindAll() 
	}

	func (c *projectController) FindById(id int) (entity.Project, bool) {
		return c.service.FindById(id)
	}
		
	// Save method to save a project from the controller
	func (c *projectController) Save(ctx *gin.Context) entity.Project {
		// bind the JSON request body to a project entity
		var project entity.Project
		ctx.BindJSON(&project)
		savedProject := c.service.Save(project)
		return savedProject
	}

	// Delete method to delete a project from the controller
	func (c *projectController) Delete(id int) {
		c.service.Delete(id)
	}

	// Update method to update a project from the controller
	func (c *projectController) Update(id int, updateData entity.Project) entity.Project {
		updatedProject := c.service.Update(id, updateData)
		return updatedProject
	}