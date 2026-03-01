package service

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"gitlab.com/pragmaticreviews/golang-gin-poc/entity"
)

type ProjectService interface {
	Save(entity.Project) entity.Project
	FindAll() []entity.Project
	Delete(id int)
	Update(id int, updateData entity.Project) entity.Project
	FindById(id int) (entity.Project, bool)
}

type projectService struct {
	db *pgx.Conn
}

// constructor function to create a new project service instance with database connection
func NewProjectService(db *pgx.Conn) ProjectService {
	return &projectService{db: db}
}

// Save method to save a project to the database
func (service *projectService) Save(project entity.Project) entity.Project {
	ctx := context.Background()
	// Convert stack slice to comma-separated string
	stackStr := strings.Join(project.Stack, ",")

	var id int
	err := service.db.QueryRow(ctx,
		`INSERT INTO project (title, description, url, stack) VALUES ($1, $2, $3, $4) RETURNING id`,
		project.Title, project.Description, project.Url, stackStr,
	).Scan(&id)

	if err != nil {
		log.Printf("Error saving project: %v", err)
		return entity.Project{}
	}

	project.Id = id
	return project
}

// FindAll method to return all projects from the database
func (service *projectService) FindAll() []entity.Project {
	ctx := context.Background()
	rows, err := service.db.Query(ctx, `SELECT id, title, description, url, stack FROM project`)
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		return []entity.Project{}
	}
	defer rows.Close()

	var projects []entity.Project
	for rows.Next() {
		var project entity.Project
		var stackStr string
		err := rows.Scan(&project.Id, &project.Title, &project.Description, &project.Url, &stackStr)
		if err != nil {
			log.Printf("Error scanning project: %v", err)
			continue
		}
		// Convert comma-separated string to slice
		if stackStr != "" {
			project.Stack = strings.Split(stackStr, ",")
		} else {
			project.Stack = []string{}
		}
		projects = append(projects, project)
	}

	return projects
}

// FindById method to find a project by ID from the database
func (service *projectService) FindById(id int) (entity.Project, bool) {
	ctx := context.Background()
	var project entity.Project
	var stackStr string

	err := service.db.QueryRow(ctx,
		`SELECT id, title, description, url, stack FROM project WHERE id = $1`, id,
	).Scan(&project.Id, &project.Title, &project.Description, &project.Url, &stackStr)

	if err != nil {
		log.Printf("Error finding project: %v", err)
		return entity.Project{}, false
	}

	if stackStr != "" {
		project.Stack = strings.Split(stackStr, ",")
	} else {
		project.Stack = []string{}
	}

	return project, true
}

// Delete method to delete a project from the database
func (service *projectService) Delete(id int) {
	ctx := context.Background()
	_, err := service.db.Exec(ctx, `DELETE FROM project WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting project: %v", err)
	}
}

// Update method to update a project in the database
func (service *projectService) Update(id int, updateData entity.Project) entity.Project {
	ctx := context.Background()
	stackStr := strings.Join(updateData.Stack, ",")

	_, err := service.db.Exec(ctx,
		`UPDATE project SET title = $1, description = $2, url = $3, stack = $4 WHERE id = $5`,
		updateData.Title, updateData.Description, updateData.Url, stackStr, id,
	)

	if err != nil {
		log.Printf("Error updating project: %v", err)
		return entity.Project{}
	}

	updateData.Id = id
	return updateData
}