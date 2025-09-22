package utils

import (
	"errors"
	"fmt"
	"strings"
)

// GetGroupsByRole - return global role
func GetGroupsByRole(role string) string {
	switch {
	case role == "SA":
		return "/global-SA"
	case strings.HasPrefix(role, "CA-"):
		context := strings.TrimPrefix(role, "CA-")
		return fmt.Sprintf("/%s/CA", context)
	case strings.HasPrefix(role, "CO-"):
		context := strings.TrimPrefix(role, "CO-")
		return fmt.Sprintf("/%s/CO", context)
	default:
		return ""
	}
}

// ParseRole - Return ShortRole, Context, err
func ParseRole(role string) (string, string, error) {
	switch {
	case role == "SA":
		return "SA", "", nil
	case strings.HasPrefix(role, "CA-"):
		context := strings.TrimPrefix(role, "CA-")
		return "CA", context, nil
	case strings.HasPrefix(role, "CO-"):
		context := strings.TrimPrefix(role, "CO-")
		return "CO", context, nil
	default:
		return "", "", errors.New("parse role")
	}
}
