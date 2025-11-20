package handlers

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository/postgres"
)

type TeamRepository interface {
	GetTeam(team_name string) (ds.Team, error)
	AddTeam(teamData dto.TeamDTO) (ds.Team, error)
}

type TeamHandler struct {
	repo TeamRepository
}

func NewTeamHandler(pg *postgres.Postgres) (*TeamHandler, error) {
	return &TeamHandler{pg}, nil
}
