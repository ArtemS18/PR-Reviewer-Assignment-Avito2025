package postgres

import (
	"reviewer-api/internal/app/ds"
	"reviewer-api/internal/app/repository"

	"gorm.io/gorm/clause"
)

func (p *Postgres) SetUserFlag(user_id string, is_active bool) (ds.User, error) {
	var user ds.User
	res := p.db.Model(&ds.User{}).Where("id = ?", user_id).Clauses(clause.Returning{}).Update("is_active", is_active).Scan(&user)
	if res.Error != nil {
		return ds.User{}, res.Error
	}
	if res.RowsAffected == 0 {
		return ds.User{}, repository.ErrNotFound
	}
	return user, nil
}

func (p *Postgres) GetReview(user_id string) (ds.User, error) {
	var user ds.User
	err := p.db.Model(&ds.User{}).Preload("Assigned").Where("id = ?", user_id).First(&user).Error
	if err != nil {
		return ds.User{}, repository.HandelPgError(err, "users")
	}
	return user, nil
}
