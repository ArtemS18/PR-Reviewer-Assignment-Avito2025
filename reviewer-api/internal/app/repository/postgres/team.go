package postgres

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"

	"gorm.io/gorm/clause"
)

func (p *Postgres) GetTeam(team_name string) (ds.Team, error) {
	var team ds.Team
	err := p.db.Preload("Members").Where("name = ?", team_name).First(&team).Error
	if err != nil {
		return ds.Team{}, repository.HandelPgError(err, "team")
	}

	return team, nil
}

func (p *Postgres) createOrUpdateMembers(teamID string, usersDTO []dto.UserDTO) ([]ds.User, error) {
	usersORM := dto.ToUserORMList(teamID, usersDTO)
	err := p.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "is_active", "team_id"}),
	}).Create(&usersORM).Error
	if err != nil {
		return nil, repository.HandelPgError(err, "team")
	}
	return usersORM, nil
}
func (p *Postgres) AddTeam(teamData dto.TeamDTO) (ds.Team, error) {
	team := ds.Team{Name: teamData.Name}
	err := p.db.Create(&team).Error
	if err != nil {
		return ds.Team{}, repository.HandelPgError(err, "team")
	}
	members, err := p.createOrUpdateMembers(team.ID, teamData.Members)
	if err != nil {
		return ds.Team{}, repository.HandelPgError(err, "team")
	}
	team.Members = members
	return team, nil
}
