package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/pragmaticreviews/golang-gin-poc/entity"
	"gitlab.com/pragmaticreviews/golang-gin-poc/service"
)

type ExperienceController interface {
	FindAll() []entity.Experience
	Save(ctx *gin.Context) entity.Experience
	Delete(id int)
	Update(id int, updateData entity.Experience) entity.Experience
	FindById(id int) (entity.Experience, bool)
}

type experienceController struct {
	service service.ExperienceService
}

func NewExperienceController(service service.ExperienceService) ExperienceController {
	return &experienceController {
		service: service,
	}
}

func (c *experienceController) FindAll() []entity.Experience {
	return c.service.FindAll() 
}

func (c *experienceController) FindById(id int) (entity.Experience, bool) {
	return c.service.FindById(id)
}

func (c *experienceController) Save(ctx *gin.Context) entity.Experience {
	var experience entity.Experience
	ctx.BindJSON(&experience)
	savedExperience := c.service.Save(experience)
	return savedExperience
}

func (c *experienceController) Delete(id int) {
	c.service.Delete(id)
}

func (c *experienceController) Update(id int, updateData entity.Experience) entity.Experience {
	updatedExperience := c.service.Update(id, updateData)
	return updatedExperience
}
