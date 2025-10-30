package shared

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		s := strings.Trim(v, "{}")
		if s == "" {
			*a = StringArray{}
			return nil
		}
		*a = strings.Split(s, ",")
		return nil
	case []byte:
		return a.Scan(string(v))
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (a StringArray) Value() (driver.Value, error) {
	return "{" + strings.Join(a, ",") + "}", nil
}
