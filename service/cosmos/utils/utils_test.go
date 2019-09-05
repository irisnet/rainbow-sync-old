package cosmos

import (
	"testing"
	"encoding/json"
)

func TestParseRewards(t *testing.T) {
	val := "1786stake"
	data := ParseRewards(val)
	bytemsg, _ := json.Marshal(data)
	t.Log(string(bytemsg))
}
