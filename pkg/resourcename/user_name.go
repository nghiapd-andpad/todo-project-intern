package resourcename

import (
	"fmt"
	"regexp"
	"strconv"
)

// users/{user_id}
var userRegex = regexp.MustCompile(`^users/(\d+)$`)

func ParseUserResourceName(name string) (int64, error) {
	matches := userRegex.FindStringSubmatch(name)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid user resource name: %q", name)
	}
	uID, _ := strconv.ParseInt(matches[1], 10, 64)
	return uID, nil
}
