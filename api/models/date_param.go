package models

import (
	"encoding/json"
	"time"
)

// DateParam represents a date parameter with validation
type DateParam struct {
	time.Time
	Valid bool
}

func (dp *DateParam) UnmarshalJSON(b []byte) error {
	// Si viene como string vacío
	if string(b) == `""` {
		dp.Valid = false
		return nil
	}

	// Intentar parsear como RFC3339
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	dp.Time = t
	dp.Valid = true
	return nil
}
