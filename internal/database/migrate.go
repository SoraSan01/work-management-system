package database

import "workms/internal/models"

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.TeamMember{},
		&models.Project{},
		&models.Task{},
		&models.Schedule{},
		&models.Notification{},
		&models.Document{},
		&models.ActivityLog{},
		&models.Department{},
		&models.Role{},
	)
	if err != nil {
		panic(err)
	}
}
