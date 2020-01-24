// +build windows

package steamworks_wrapper

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	uintptrSize           uintptr = 0
	sizeOfSteamApiContext         = 512
	steamApiContextSize           = 2 + (sizeOfSteamApiContext / unsafe.Sizeof(uintptrSize))
)

var (
	initCompleted = false

	steamAPIIsSteamRunning *syscall.Proc

	// internal
	steamInternalContextInit     *syscall.Proc
	steamInternalCreateInterface *syscall.Proc

	// user
	steamApiGetSteamUser *syscall.Proc
	steamApiGetSteamPipe *syscall.Proc

	// steam client
	steamApiISteamClientGetISteamUser  *syscall.Proc
	steamApiISteamClientGetISteamUtils *syscall.Proc

	// utils
	steamApiISteamUtilsIsCallbackCompleted     *syscall.Proc
	steamApiISteamUtilsGetAPICallFailureReason *syscall.Proc

	// encrypted app tickets
	steamApiISteamUserRequestEncryptedAppTicket *syscall.Proc

	// user/pipe state
	user int32
	pipe int32
)

var (
	steamApiContextData [steamApiContextSize]uintptr
	steamApiContext     *steamApiContextType
)

type steamApiContextType struct {
	m_pSteamClient             uintptr
	m_pSteamUser               uintptr
	m_pSteamFriends            uintptr
	m_pSteamUtils              uintptr
	m_pSteamMatchmaking        uintptr
	m_pSteamGameSearch         uintptr
	m_pSteamUserStats          uintptr
	m_pSteamApps               uintptr
	m_pSteamMatchmakingServers uintptr
	m_pSteamNetworking         uintptr
	m_pSteamRemoteStorage      uintptr
	m_pSteamScreenshots        uintptr
	m_pSteamHTTP               uintptr
	m_pController              uintptr
	m_pSteamUGC                uintptr
	m_pSteamAppList            uintptr
	m_pSteamMusic              uintptr
	m_pSteamMusicRemote        uintptr
	m_pSteamHTMLSurface        uintptr
	m_pSteamInventory          uintptr
	m_pSteamVideo              uintptr
	m_pSteamTV                 uintptr
	m_pSteamParentalSettings   uintptr
	m_pSteamInput              uintptr
}

func init() {
	dll := syscall.MustLoadDLL("steam_api64.dll")
	steamAPIInit := dll.MustFindProc("SteamAPI_Init")
	ret, _, err := steamAPIInit.Call()
	if err != nil && err.(syscall.Errno) != 0 {
		panic(err)
	}

	if ret == 1 {
		initCompleted = true
	}

	steamInternalContextInit = dll.MustFindProc("SteamInternal_ContextInit")
	steamInternalCreateInterface = dll.MustFindProc("SteamInternal_CreateInterface")

	steamAPIIsSteamRunning = dll.MustFindProc("SteamAPI_IsSteamRunning")

	steamApiGetSteamUser = dll.MustFindProc("SteamAPI_GetHSteamUser")
	steamApiGetSteamPipe = dll.MustFindProc("SteamAPI_GetHSteamPipe")
	steamApiISteamClientGetISteamUser = dll.MustFindProc("SteamAPI_ISteamClient_GetISteamUser")
	steamApiISteamClientGetISteamUtils = dll.MustFindProc("SteamAPI_ISteamClient_GetISteamUtils")
	steamApiISteamUserRequestEncryptedAppTicket = dll.MustFindProc("SteamAPI_ISteamUser_RequestEncryptedAppTicket")
	steamApiISteamUtilsIsCallbackCompleted = dll.MustFindProc("SteamAPI_ISteamUtils_IsAPICallCompleted")
	steamApiISteamUtilsGetAPICallFailureReason = dll.MustFindProc("SteamAPI_ISteamUtils_GetAPICallFailureReason")

	initContext()
}

func ensureLoaded() {
	for {
		if steamApiContext != nil {
			return
		}
	}
}

func initContext() {
	ret, _, _ := steamApiGetSteamUser.Call()
	if ret == 0 {
		panic(fmt.Errorf("could not get user from steam"))
	}

	user = int32(ret)

	ret, _, _ = steamApiGetSteamPipe.Call()
	if ret == 0 {
		panic(fmt.Errorf("could not get pipe from steam"))
	}

	pipe = int32(ret)

	newCallbackCounterAndContext := &steamApiContextData
	newCallbackCounterAndContext[0] = syscall.NewCallback(onContextInit)
	_, _, _ = steamInternalContextInit.Call(uintptr(unsafe.Pointer(&newCallbackCounterAndContext[0])))
}

func onContextInit(steamContext uintptr) uintptr {
	steamApiContext = (*steamApiContextType)(unsafe.Pointer(steamContext))
	steamApiContext.m_pSteamClient = createInterface("SteamClient020")
	steamApiContext.m_pSteamUser = createSteamClientInterface(steamApiISteamClientGetISteamUser, "SteamUser020")
	steamApiContext.m_pSteamUtils = createSteamPipeInterface(steamApiISteamClientGetISteamUtils, "SteamUtils009")

	return 0
}

func createSteamPipeInterface(ptr *syscall.Proc, name string) uintptr {
	data := []byte(name)
	ret, _, _ := ptr.Call(steamApiContext.m_pSteamClient, uintptr(pipe), uintptr(unsafe.Pointer(&data[0])))
	return ret
}

func createSteamClientInterface(ptr *syscall.Proc, name string) uintptr {
	data := []byte(name)
	ret, _, _ := ptr.Call(steamApiContext.m_pSteamClient, uintptr(user), uintptr(pipe), uintptr(unsafe.Pointer(&data[0])))
	return ret
}

func createInterface(name string) uintptr {
	data := []byte(name)
	ret, _, _ := steamInternalCreateInterface.Call(uintptr(unsafe.Pointer(&data[0])))
	return ret
}
