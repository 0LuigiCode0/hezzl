package dusecase

import (
	"fmt"
	"strings"
)

func validateString(v *string, fieldName string) error {
	*v = strings.TrimSpace(*v)
	if *v == "" {
		return fmt.Errorf("поле: %s", fieldName)
	}
	return nil
}
