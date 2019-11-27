package xtime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationUnmarshal(t *testing.T) {
	var data struct {
		T Duration `json:t"`
	}
	require.NoError(t, json.Unmarshal([]byte(`{"t":"10s"}`), &data))
	require.Equal(t, time.Second*10, time.Duration(data.T))
}
