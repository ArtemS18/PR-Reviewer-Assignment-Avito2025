package ds

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID   string `gorm:"type:varchar(100);primaryKey" json:"team_id"`
	Name string `gorm:"type:varchar(50);unique;not null" json:"team_name"`

	Members []User `gorm:"foreignKey:TeamID" json:"members"`
}

func (u *Team) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = fmt.Sprintf("team-%s", uuid.NewString())
	return nil
}
