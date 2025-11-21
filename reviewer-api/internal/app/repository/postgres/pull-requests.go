package postgres

import (
	"fmt"
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/dto"
	"reviewer-api/internal/app/repository"
	"reviewer-api/internal/pkg/utils"
	"time"
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
		Distinct("user_id").
		Where("pull_request_id = ?", pk_id)

	var members_ids []string
	err := p.db.Model(&ds.User{}).
		Select("users.id").
		Where("users.team_id = (?) AND users.is_active = true AND users.id != ?", sub, excluded_id).
		Where("users.id NOT IN (?)", reviewerSub).
		Scan(&members_ids).Error

	if err != nil {
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
	fmt.Println(new_reviewer_id)

	if new_reviewer_id == "" {
		return ds.PullRequest{}, repository.ErrNotFound
	}
	var pk ds.PullRequest
	err = p.db.Where("id = ?", pk_id).First(&pk).Error
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pk")
	}
	if pk.Status == string(ds.MERGED) {
		return ds.PullRequest{}, repository.ErrReassign
	}

	resp := p.db.Model(&ds.Reviewer{}).
		Where("pull_request_id = ? AND user_id =?", pk_id, old_reviewer_id).
		Update("user_id", new_reviewer_id)
	if resp.RowsAffected == 0 {
		return ds.PullRequest{}, repository.ErrNotFound
	}

	if resp.Error != nil {
		return ds.PullRequest{}, repository.ErrUnexpect
	}
	err = p.db.Preload("AssignedReviewers").Where("id = ?", pk_id).First(&pk).Error
	if err != nil {
		return ds.PullRequest{}, repository.ErrUnexpect
	}

	return pk, nil
}

func (p *Postgres) Merged(pk_id string) (ds.PullRequest, error) {
	var pk ds.PullRequest
	err := p.db.Preload("AssignedReviewers").Where("id = ?", pk_id).First(&pk).Error
	if pk.Status == string(ds.MERGED) {
		return pk, nil
	}
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pk")
	}
	t := time.Now().UTC()
	err = p.db.Model(&pk).Updates(ds.PullRequest{Status: string(ds.MERGED), MergedAt: &t}).Error
	if err != nil {
		return ds.PullRequest{}, repository.HandelPgError(err, "pk")
	}

	return pk, nil
}
