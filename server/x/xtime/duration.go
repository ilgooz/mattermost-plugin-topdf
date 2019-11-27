// package xtime extends features of package "time".
package xtime

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration is a time.Duration with JSON decoding support.
type Duration time.Duration

// UnmarshalJSON tries to unmarshal a JSON value as Duration.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if value, ok := v.(string); ok {
		dr, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(dr)
		return nil
	}
	return fmt.Errorf("invalid duration %q", b)
}
