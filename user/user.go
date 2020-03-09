package user

type User struct {
	Id       int64
	Username string
	Nickname string
}

var (
	User1 = User{
		Id:       1024,
		Username: "sez001",
		Nickname: "sez",
	}
	User2 = User{
		Id:       2048,
		Username: "zed001",
		Nickname: "zed",
	}
)
