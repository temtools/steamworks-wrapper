package steamworks_wrapper

import (
	"fmt"
	"unsafe"
)

func IsAPICallCompleted(cb uintptr) (bool, error) {
	failure := false
	ret, _, _ := steamApiISteamUtilsIsCallbackCompleted.Call(steamApiContext.m_pSteamUtils, cb, uintptr(unsafe.Pointer(&failure)))
	if ret == 1 {
		if failure {
			ret, _, _ := steamApiISteamUtilsGetAPICallFailureReason.Call(steamApiContext.m_pSteamUtils, cb)
			fmt.Printf("steam failure code: %d\n", ret)
			return true, fmt.Errorf("steam failure code: %d\n", ret)
		}

		return true, nil
	}

	return false, nil
}
