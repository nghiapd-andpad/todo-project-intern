package helper

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ParseDueDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"invalid due_date format %q, expected YYYY-MM-DD", dateStr)
	}
	return &t, nil
}
