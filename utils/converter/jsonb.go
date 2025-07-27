package converter

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Generic helper for JSONB types
type JSONB[T any] struct {
	Data T
}

func (j JSONB[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

func (j *JSONB[T]) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, &j.Data)
}
