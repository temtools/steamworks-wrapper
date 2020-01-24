package steamworks_wrapper

import "testing"

func TestGetUser(t *testing.T) {
	u := GetUser()
	if u.userHandle == 0 {
		t.Error("could not get steam user")
	}
}
