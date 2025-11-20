package ds

import (
	"time"
)

type PRStatus string

const (
	MERGED PRStatus = "MERGED"
	OPEN   PRStatus = "OPEN"
)

type PullRequest struct {
	ID       string     `gorm:"type:varchar(100);primaryKey" json:"pull_request_id"`
	Name     string     `gorm:"type:varchar(50);not null" json:"pull_request_name"`
	AuthorID string     `gorm:"not null" json:"author_id"`
	Status   string     `gorm:"type:varchar(10);default:OPEN;not null" json:"status"`
	MergedAt *time.Time `gorm:"default:null" json:"merged_at"`

	Author            User       `gorm:"foreignKey:AuthorID" json:"-"`
	AssignedReviewers []Reviewer `gorm:"foreignKey:PullRequestID" json:"assigned_reviewers"`
}
