package service

import (
	"context"
	"log"
	"strings"
	"github.com/jackc/pgx/v5"
	"gitlab.com/pragmaticreviews/golang-gin-poc/entity"
)

type ExperienceService interface {
	Save(entity.Experience) entity.Experience
	FindAll() []entity.Experience
	Delete(id int)
	Update(id int, updateData entity.Experience) entity.Experience
	FindById(id int) (entity.Experience, bool)
}

type ExpecienceService struct {
	db *pgx.Conn
}

func NewExperienceService(db *pgx.Conn) ExperienceService {
	return &ExpecienceService{db: db}
}

func (service *ExpecienceService) Save(experience entity.Experience) entity.Experience {
	ctx := context.Background()
	// Convert skills slice to comma-separated string
	skillsStr := strings.Join(experience.Skills, ",")
	var id int
	err := service.db.QueryRow(ctx,
		`INSERT INTO experience (title, description, skills, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		experience.Title, experience.Description, skillsStr, experience.StartDate, experience.EndDate,
	).Scan(&id)
	if err != nil {
		log.Printf("Error saving experience: %v", err)
		return entity.Experience{}
	}
	experience.Id = id
	return experience
}

func (service *ExpecienceService) FindAll() []entity.Experience {
	ctx := context.Background()
	rows, err := service.db.Query(ctx, `SELECT id, title, description, skills, start_date, end_date FROM experience`)
	if err != nil {
		log.Printf("Error fetching experiences: %v", err)
		return []entity.Experience{}
	}
	defer rows.Close()

	var experiences []entity.Experience
	for rows.Next() {
		var experience entity.Experience
		var skillsStr string
		err := rows.Scan(&experience.Id, &experience.Title, &experience.Description, &skillsStr, &experience.StartDate, &experience.EndDate)
		if err != nil {
			log.Printf("Error scanning experience: %v", err)
			continue
		}
		// Convert comma-separated skills string back to slice
		experience.Skills = strings.Split(skillsStr, ",")
		experiences = append(experiences, experience)
	}
	return experiences
}

func (service *ExpecienceService) FindById(id int) (entity.Experience, bool) {
	ctx := context.Background()
	var experience entity.Experience
	var skillsStr string
	err := service.db.QueryRow(ctx, `SELECT id, title, description, skills, start_date, end_date FROM experience WHERE id = $1`, id).
		Scan(&experience.Id, &experience.Title, &experience.Description, &skillsStr, &experience.StartDate, &experience.EndDate)
	if err != nil {
		log.Printf("Error fetching experience by id: %v", err)
		return entity.Experience{}, false
	}
	experience.Skills = strings.Split(skillsStr, ",")
	return experience, true
}

func (service *ExpecienceService) Delete(id int) {
	ctx := context.Background()
	_, err := service.db.Exec(ctx, `DELETE FROM experience WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error deleting experience: %v", err)
	}
}

func (service *ExpecienceService) Update(id int, updateData entity.Experience) entity.Experience {
	ctx := context.Background()
	// Convert skills slice to comma-separated string
	skillsStr := strings.Join(updateData.Skills, ",")
	_, err := service.db.Exec(ctx,
		`UPDATE experience SET title = $1, description = $2, skills = $3, start_date = $4, end_date = $5 WHERE id = $6`,
		updateData.Title, updateData.Description, skillsStr, updateData.StartDate, updateData.EndDate, id,
	)
	if err != nil {
		log.Printf("Error updating experience: %v", err)
		return entity.Experience{}
	}
	updateData.Id = id
	return updateData
}
