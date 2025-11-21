package postgres

import (
	"fmt"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/pkg/utils"
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
		return nil, repository.HandelPgError(err, "pk")
	}
	return assignedReviewers, nil
}

func (p *Postgres) getMemberIds(excluded_id string, pk_id string) ([]string, error) {
	sub := p.db.Model(&ds.Team{}).
		Select("teams.id").
		Joins("JOIN users ON teams.id = users.team_id").
		Where("users.id = ?", excluded_id)

	reviewerSub := p.db.Model(&ds.Reviewer{}).
		Select("user_id").
		Where("pull_request_id = ?", pk_id)

	var members_ids []string
	err := p.db.Model(&ds.User{}).
		Select("users.id").
		Where("users.team_id = (?) AND users.is_active = true AND users.id != ?", sub, excluded_id).
		Where("users.id NOT IN (?)", reviewerSub).
		Scan(&members_ids).Error

	if err != nil || len(members_ids) == 0 {
		return nil, repository.ErrNotFound
	}
	return members_ids, nil
}

func (p *Postgres) CreatePullRequest(pkDTO dto.PullRequestCreateDTO) (ds.PullRequest, error) {

	members_ids, err := p.getMemberIds(pkDTO.AuthorID, pkDTO.ID)

	if err != nil {
		return ds.PullRequest{}, repository.ErrNotFound
	}
	pk := ds.PullRequest{
		ID:       pkDTO.ID,
		Name:     pkDTO.Name,
		AuthorID: pkDTO.AuthorID,
	}
	err = p.db.Create(&pk).Error
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pr")
	}
	assignReviewers, err := p.assignReviewers(members_ids, pk.ID)
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pr")
	}
	pk.AssignedReviewers = assignReviewers
	return pk, nil
}

func (p *Postgres) ReassignReviewer(pk_id string, old_reviewer_id string) (ds.PullRequest, error) {
	members_ids, err := p.getMemberIds(old_reviewer_id, pk_id)
	fmt.Println(members_ids)

	if err != nil {
		return ds.PullRequest{}, repository.ErrNotFound
	}

	new_reviewer_id := utils.GetRandomNumber(members_ids)

	if new_reviewer_id == "" {
		return ds.PullRequest{}, repository.ErrNotFound
	}

	resp := p.db.Model(&ds.Reviewer{}).
		Where("pull_request_id = ? AND user_id =?", pk_id, old_reviewer_id).
		Update("user_id", new_reviewer_id)

	if resp.Error != nil {
		return ds.PullRequest{}, repository.ErrUnexpect
	}
	var pk ds.PullRequest
	err = p.db.Preload("AssignedReviewers").Where("id = ?", pk_id).First(&pk).Error
	if err != nil {
		fmt.Println(err.Error())
		return ds.PullRequest{}, repository.HandelPgError(err, "pk")
	}
	return pk, nil
}
