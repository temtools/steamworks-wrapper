package steamworks_wrapper

import "testing"

func TestLoader(t *testing.T) {
	if !initCompleted {
		t.Error("steam api not initialized")
	}
}
