package steamworks_wrapper

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestUser_RequestEncryptedAppTicket(t *testing.T) {
	user := GetUser()
	ticket, err := user.RequestEncryptedAppTicket()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("ticket length: %d\n", len(ticket))
	fmt.Printf("%s", hex.Dump(ticket))
}
