package steamworks_wrapper

import (
	"fmt"
	"time"
	"unsafe"
)

func (user *User) RequestEncryptedAppTicket() ([]byte, error) {
	ensureLoaded()

	ret, _, _ := steamApiISteamUserRequestEncryptedAppTicket.Call(steamApiContext.m_pSteamUser, 0, 0)
	for {
		if abort, err := IsAPICallCompleted(ret); abort {
			if err != nil {
				return nil, err
			}

			size := 0
			retGetSizeCheck, _, _ := steamApiISteamUserGetEncryptedAppTicket.Call(steamApiContext.m_pSteamUser, 0, 0, uintptr(unsafe.Pointer(&size)))
			if retGetSizeCheck == 0 && size > 0 {
				// We got a ticket, get it
				ticket := make([]byte, size)
				retGetTicket, _, _ := steamApiISteamUserGetEncryptedAppTicket.Call(steamApiContext.m_pSteamUser, uintptr(unsafe.Pointer(&ticket[0])), uintptr(size), uintptr(unsafe.Pointer(&size)))
				if retGetTicket == 1 {
					return ticket, nil
				}
			}

			return nil, fmt.Errorf("request encrypted app ticket failed")
		}

		time.Sleep(1 * time.Millisecond)
	}
}
