package postgres

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) assignReviewers(members_ids []string, pk_id string) ([]ds.Reviewer, error) {

	assignedIds := utils.GetRandomSlice(members_ids)
	if len(assignedIds) == 0 {
		return nil, nil
	}
	assignedReviewers := make([]ds.Reviewer, 0, len(assignedIds))
	for _, rewId := range assignedIds {
		assignedReviewers = append(assignedReviewers, ds.Reviewer{UserID: rewId, PullRequestID: pk_id})
	}

	err := p.db.Create(&assignedReviewers).Error
	if err != nil {
		return nil, repository.ErrUnexpect
	}
	return assignedReviewers, nil
}

func (p *Postgres) CreatePullRequest(pkDTO dto.PullRequestCreateDTO) (ds.PullRequest, error) {
	var members_ids []string

	sub := p.db.Model(&ds.Team{}).
		Select("teams.id").
		Joins("JOIN users ON teams.id = users.team_id").
		Where("users.id = ?", pkDTO.AuthorID)

	err := p.db.Model(&ds.User{}).
		Select("users.id").
		Where("users.team_id = (?) AND users.is_active = true AND users.id != ?", sub, pkDTO.AuthorID).
		Scan(&members_ids).Error

	if err != nil || len(members_ids) == 0 {
		return ds.PullRequest{}, repository.ErrNotFound
	}
	pk := ds.PullRequest{
		ID:       pkDTO.ID,
		Name:     pkDTO.Name,
		AuthorID: pkDTO.AuthorID,
	}
	err = p.db.Create(&pk).Error
	if err != nil {
		pqErr, ok := err.(*pgconn.PgError)
		if ok && pqErr.Code == "23505" {
			return ds.PullRequest{}, repository.ErrAlreadyExists
		}
		if ok && pqErr.Code == "23503" {
			return ds.PullRequest{}, repository.ErrNotFound
		}
		return ds.PullRequest{}, repository.ErrUnexpect
	}
	assignReviewers, err := p.assignReviewers(members_ids, pk.ID)
	if err != nil {
		return ds.PullRequest{}, repository.ErrUnexpect
	}
	pk.AssignedReviewers = assignReviewers
	return pk, nil
}
