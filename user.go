package steamworks_wrapper

type User struct {
	userHandle int32
}

func GetUser() *User {
	return &User{}
}
