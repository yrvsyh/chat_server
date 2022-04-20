package config

const (
	PasswordSalt = "some salt"
	CookieKey    = "auth"
	SessionKey   = "session key"

	AvatarFileSizeLimit = 1 * 1024 * 1024
	AvatarPathPrefix    = "./static/avatar/"

	ImagePathPrefix = "./static/upload/image/"
	FilePathPrefix  = "./static/upload/file/"

	RedisChannelGroupMessageKeyPrefix = "group_message_"
)
