package dto

import "reviewer-api/internal/app/ds"

type TeamDTO struct {
	Name    string    `json:"team_name"`
	Members []UserDTO `json:"members"`
}

func ToTeamDTO(teamORM ds.Team) TeamDTO {
	usersDTO := make([]UserDTO, 0, len(teamORM.Members))
	for _, user := range teamORM.Members {
		usersDTO = append(usersDTO, ToUserDTO(user))
	}
	return TeamDTO{
		Name:    teamORM.ID,
		Members: usersDTO,
	}
}
