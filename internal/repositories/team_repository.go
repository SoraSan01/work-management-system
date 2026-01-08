package repositories

import (
	"errors"
	"workms/internal/models"

	"gorm.io/gorm"
)

type TeamRepository struct {
	DB *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{DB: db}
}

// FindAll retrieves all teams with their members and users
func (r *TeamRepository) FindAll() ([]models.Team, error) {
	var teams []models.Team
	err := r.DB.
		Preload("Members.User").
		Order("id DESC").
		Find(&teams).Error

	if err != nil {
		return nil, err
	}
	return teams, nil
}

// Create creates a new team
func (r *TeamRepository) Create(team *models.Team) error {
	if team.Name == "" {
		return errors.New("team name is required")
	}
	return r.DB.Create(team).Error
}

// FindByID retrieves a team by ID with all relations
func (r *TeamRepository) FindByID(id uint64) (*models.Team, error) {
	var team models.Team
	err := r.DB.
		Preload("Members.User").
		First(&team, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, err
	}
	return &team, nil
}

// Update updates an existing team
func (r *TeamRepository) Update(team *models.Team) error {
	if team.ID == 0 {
		return errors.New("team ID is required")
	}
	if team.Name == "" {
		return errors.New("team name is required")
	}

	// Check if team exists
	var existing models.Team
	if err := r.DB.First(&existing, team.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("team not found")
		}
		return err
	}

	return r.DB.Save(team).Error
}

// Delete removes a team by ID
func (r *TeamRepository) Delete(id uint64) error {
	// Check if team exists
	var team models.Team
	if err := r.DB.First(&team, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("team not found")
		}
		return err
	}

	// Delete team (this should cascade to members if properly configured)
	return r.DB.Delete(&models.Team{}, id).Error
}

// FindPaginated retrieves teams with pagination
func (r *TeamRepository) FindPaginated(page, pageSize int) (teams []models.Team, total int64, err error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Get total count
	if err = r.DB.Model(&models.Team{}).Count(&total).Error; err != nil {
		return
	}

	// Get paginated results
	err = r.DB.
		Preload("Members.User").
		Limit(pageSize).
		Offset(offset).
		Order("id DESC").
		Find(&teams).Error

	return
}

// GetAllUsers retrieves all users for the member dropdown
func (r *TeamRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Order("first_name ASC, last_name ASC").
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

// AddMember adds a member to a team
func (r *TeamRepository) AddMember(member *models.TeamMember) error {
	if member.TeamID == 0 {
		return errors.New("team ID is required")
	}
	if member.UserID == 0 {
		return errors.New("user ID is required")
	}

	// Check if team exists
	var team models.Team
	if err := r.DB.First(&team, member.TeamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("team not found")
		}
		return err
	}

	// Check if user exists
	var user models.User
	if err := r.DB.First(&user, member.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if user is already a member
	var existing models.TeamMember
	err := r.DB.Where("team_id = ? AND user_id = ?", member.TeamID, member.UserID).First(&existing).Error
	if err == nil {
		return errors.New("user is already a member of this team")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Set default role if not provided
	if member.Role == "" {
		member.Role = "member"
	}

	return r.DB.Create(member).Error
}

// RemoveMember removes a member from a team
func (r *TeamRepository) RemoveMember(teamID, memberID uint64) error {
	if teamID == 0 {
		return errors.New("team ID is required")
	}
	if memberID == 0 {
		return errors.New("member ID is required")
	}

	// Check if member exists
	var member models.TeamMember
	err := r.DB.Where("team_id = ? AND id = ?", teamID, memberID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("member not found")
		}
		return err
	}

	return r.DB.Delete(&member).Error
}

// GetMembersByTeamID retrieves all members of a specific team
func (r *TeamRepository) GetMembersByTeamID(teamID uint64) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.DB.
		Preload("User").
		Where("team_id = ?", teamID).
		Find(&members).Error

	if err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateMemberRole updates a member's role in a team
func (r *TeamRepository) UpdateMemberRole(teamID, memberID uint64, newRole string) error {
	if teamID == 0 {
		return errors.New("team ID is required")
	}
	if memberID == 0 {
		return errors.New("member ID is required")
	}
	if newRole != "member" && newRole != "leader" {
		return errors.New("invalid role: must be 'member' or 'leader'")
	}

	result := r.DB.Model(&models.TeamMember{}).
		Where("team_id = ? AND id = ?", teamID, memberID).
		Update("role", newRole)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("member not found")
	}

	return nil
}

// FindAllWithRelations fetches all teams with their creator
func (r *TeamRepository) FindAllWithRelations() ([]models.Team, error) {
	var teams []models.Team
	err := r.DB.
		Preload("User").
		Preload("Members").
		Find(&teams).Error
	return teams, err
}

// GetProjectCountForTeam returns project count for a team
func (r *TeamRepository) GetProjectCountForTeam(teamID uint64) (total, active int64, err error) {
	err = r.DB.Model(&models.Project{}).Where("team_id = ?", teamID).Count(&total).Error
	if err != nil {
		return
	}

	err = r.DB.Model(&models.Project{}).
		Where("team_id = ? AND status = ?", teamID, "active").
		Count(&active).Error
	return
}

// GetTaskCountForTeam returns total tasks for a team's projects
func (r *TeamRepository) GetTaskCountForTeam(teamID uint64) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Task{}).
		Joins("JOIN projects ON projects.id = tasks.project_id").
		Where("projects.team_id = ?", teamID).
		Count(&count).Error
	return count, err
}

func (r *TeamRepository) FindByIDWithMembers(id uint64) (*models.Team, error) {
	var team models.Team
	err := r.DB.Preload("Members").First(&team, id).Error
	return &team, err
}
