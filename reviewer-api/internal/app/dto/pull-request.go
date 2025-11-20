package dto

import (
	"reviewer-api/internal/app/ds"
	"time"
)

type PullRequestCreateDTO struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}
type PullRequestDTO struct {
	ID       string     `json:"pull_request_id"`
	Name     string     `json:"pull_request_name"`
	AuthorID string     `json:"author_id"`
	Status   string     `json:"status"`
	MergedAt *time.Time `json:"merged_at"`

	AssignedReviewers []string `gorm:"foreignKey:PullRequestID" json:"assigned_reviewers"`
}

func ToPullRequestDTO(orm ds.PullRequest) PullRequestDTO {
	rews := make([]string, 0, len(orm.AssignedReviewers))
	for _, user := range orm.AssignedReviewers {
		rews = append(rews, user.UserID)
	}
	return PullRequestDTO{
		ID:                orm.ID,
		Name:              orm.Name,
		AuthorID:          orm.AuthorID,
		Status:            orm.Status,
		MergedAt:          orm.MergedAt,
		AssignedReviewers: rews,
	}
}
