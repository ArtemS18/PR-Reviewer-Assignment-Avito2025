package postgres

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (p *Postgres) GetTeam(team_name string) (ds.Team, error) {
	var team ds.Team
	err := p.db.Preload("Members").Where("name = ?", team_name).First(&team).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ds.Team{}, repository.ErrNotFound
		}
		return ds.Team{}, err
	}

	return team, nil
}

func (p *Postgres) createMembers(teamID string, usersDTO []dto.UserDTO) ([]ds.User, error) {
	usersORM := dto.ToUserORMList(teamID, usersDTO)
	err := p.db.Model(&ds.User{}).Create(&usersORM).Error
	if err != nil {
		return nil, err
	}
	return usersORM, nil
}
func (p *Postgres) createTeam(teamData dto.TeamDTO) (ds.Team, error) {
	team := ds.Team{Name: teamData.Name}
	err := p.db.Create(&team).Error
	if err != nil {
		return ds.Team{}, err
	}

	members, err := p.createMembers(team.ID, teamData.Members)
	if err != nil {
		return ds.Team{}, err
	}
	team.Members = members
	return team, nil
}

func (p *Postgres) AddTeam(teamData dto.TeamDTO) (ds.Team, error) {
	var team ds.Team
	err := p.db.Preload("Members").Where("name = ?", teamData.Name).First(&team).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			team, err = p.createTeam(teamData)
			if err != nil {
				return ds.Team{}, err
			}
			return team, nil
		default:
			return ds.Team{}, err
		}
	}

	newMembers := dto.ToUserORMList(team.ID, teamData.Members)

	for _, user := range newMembers {
		err := p.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "is_active"}),
		}).Create(&user).Error
		if err != nil {
			return ds.Team{}, err
		}
	}
	err = p.db.Model(&team).Association("Members").Replace(newMembers)
	if err != nil {
		return ds.Team{}, err
	}
	return team, nil
}
