package steamworks_wrapper

import "testing"

func TestUser_GetSteamID(t *testing.T) {
	u := GetUser()
	steamID := u.GetSteamID()
	if steamID != "76561198041062849" {
		t.Errorf("steamID was %s", steamID)
	}
}
