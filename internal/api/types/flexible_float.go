package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleFloat allows decoding from either a JSON number or a JSON string
type FlexibleFloat float64

func (f *FlexibleFloat) UnmarshalJSON(b []byte) error {
	// Try number first
	var num float64
	if err := json.Unmarshal(b, &num); err == nil {
		*f = FlexibleFloat(num)
		return nil
	}

	// Try string
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		parsed, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return fmt.Errorf("invalid float string: %w", err)
		}
		*f = FlexibleFloat(parsed)
		return nil
	}

	return fmt.Errorf("invalid flexible float value: %s", string(b))
}
