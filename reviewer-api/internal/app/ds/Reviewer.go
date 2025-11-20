package ds

type Reviewer struct {
	UserID        string `gorm:"type:varchar(100);primaryKey"`
	PullRequestID string `gorm:"type:varchar(100);primaryKey"`

	User        User        `gorm:"foreignKey:UserID"`
	PullRequest PullRequest `gorm:"foreignKey:PullRequestID"`
}
