package config

import "os"

const (
	PasswordSalt        = "some salt"
	CookieKey           = "auth"
	SessionKey          = "session key"
	AvatarPathPrefix    = "./static" + string(os.PathSeparator) + "avatar"
	AvatarFileSizeLimit = 1 * 1024 * 1024

	RedisChannelGroupMessageKeyPrefix = "group_message_"
)
