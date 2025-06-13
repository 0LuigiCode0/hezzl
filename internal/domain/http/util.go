package dhttp

import (
	"fmt"
	"strings"

	"github.com/0LuigiCode0/hezzl/internal/domain/consts"
)

func validateString(v *string, fieldName string) error {
	*v = strings.TrimSpace(*v)
	if *v == "" {
		return fmt.Errorf(consts.ErrFieldEmpty, fieldName)
	}
	return nil
}
