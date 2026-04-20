package errors

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
)

const (
	mysqlErrCodeDuplicateEntry = 1062
)

func ParseDuplicateField(err error) error {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return err
	}

	if mysqlErr.Number != mysqlErrCodeDuplicateEntry {
		return err
	}

	msg := mysqlErr.Message
	switch {
	case strings.Contains(msg, "idx_users_username") || strings.Contains(msg, "username"):
		return entity.ErrUsernameAlreadyExists
	case strings.Contains(msg, "idx_users_email") || strings.Contains(msg, "email"):
		return entity.ErrEmailAlreadyExists
	}

	return entity.ErrUserAlreadyExists
}
