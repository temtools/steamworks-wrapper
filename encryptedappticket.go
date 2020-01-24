package steamworks_wrapper

import (
	"fmt"
	"time"
)

func (user *User) RequestEncryptedAppTicket() ([]byte, error) {
	ensureLoaded()

	ret, _, _ := steamApiISteamUserRequestEncryptedAppTicket.Call(steamApiContext.m_pSteamUser, 0, 0)
	for {
		if abort, err := IsAPICallCompleted(ret); abort {
			if err != nil {
				return nil, err
			}

			fmt.Printf("got request encrypted app ticket")
		}

		time.Sleep(1 * time.Millisecond)
	}
}
