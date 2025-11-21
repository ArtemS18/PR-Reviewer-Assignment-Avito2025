package repository

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrNotFound = fmt.Errorf("not found")
var ErrTeamAlreadyExists = fmt.Errorf("team already exists")
var ErrPRAlreadyExists = fmt.Errorf("pr already exists")
var ErrUnexpect = fmt.Errorf("unexpect err")
var ErrReassign = fmt.Errorf("cannot reassign on merged PR")

func HandelPgError(err error, table string) error {
	pqErr, ok := err.(*pgconn.PgError)
	if ok && pqErr.Code == "23505" {
		switch table {
		case "team":
			return ErrTeamAlreadyExists
		case "pr":
			return ErrPRAlreadyExists
		}
	}
	if ok && pqErr.Code == "23503" {
		return ErrNotFound
	}
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return ErrUnexpect
}
