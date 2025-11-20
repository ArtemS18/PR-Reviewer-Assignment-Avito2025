package ds

type User struct {
	ID       string `gorm:"type:varchar(100);primaryKey" json:"user_id"`
	Name     string `gorm:"type:varchar(25);not null" json:"username"`
	IsActive bool   `gorm:"type:boolean" json:"is_active"`
	TeamID   string `gorm:"type:varchar(100);default:null" json:"-"`

	Team     Team          `gorm:"foreignKey:TeamID" json:"-"`
	Assigned []PullRequest `gorm:"many2many:reviewers" json:"pull_requests"`
}
