package config

import "os"

const (
	PasswordSalt        = "some salt"
	CookieKey           = "auth"
	SessionKey          = "session key"
	SessionUserKey      = "name"
	AvatarPathPrefix    = "resource" + string(os.PathSeparator) + "avatar"
	AvatarFileSizeLimit = 1 * 1024 * 1024
)
