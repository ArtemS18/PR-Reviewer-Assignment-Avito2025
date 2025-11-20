package dto

import "reviewer-api/internal/app/ds"

type UserDTO struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func ToUserDTO(userORM ds.User) UserDTO {
	return UserDTO{
		ID:       userORM.ID,
		Name:     userORM.Name,
		IsActive: userORM.IsActive,
	}
}

func ToUserORM(teamID string, userDTO UserDTO) ds.User {
	return ds.User{
		ID:       userDTO.ID,
		Name:     userDTO.Name,
		IsActive: userDTO.IsActive,
		TeamID:   teamID,
	}
}

func ToUserORMList(teamID string, usersDTO []UserDTO) []ds.User {
	usersORM := make([]ds.User, 0, len(usersDTO))
	for _, user := range usersDTO {
		userORM := ToUserORM(teamID, user)
		usersORM = append(usersORM, userORM)
	}
	return usersORM
}
