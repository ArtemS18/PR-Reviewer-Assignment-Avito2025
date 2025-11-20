package repository

import (
	"fmt"
)

var ErrNotFound = fmt.Errorf("not found")
var ErrAlreadyExists = fmt.Errorf("already exists")
var ErrUnexpect = fmt.Errorf("unexpect err")
