package steamworks_wrapper

import "testing"

func TestUser_RequestEncryptedAppTicket(t *testing.T) {
	user := GetUser()
	user.RequestEncryptedAppTicket()
}
