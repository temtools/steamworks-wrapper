package steamworks_wrapper

import "strconv"

type User struct {
}

func GetUser() *User {
	return &User{}
}

func (user *User) GetSteamID() string {
	ret, _, _ := steamApiISteamUserGetSteamID.Call(steamApiContext.m_pSteamUser)
	return strconv.FormatInt(int64(ret), 10)
}
