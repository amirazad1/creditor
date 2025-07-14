package constants

const (
	// User
	AdminRoleName   string = "admin"
	DefaultRoleName string = "default"
	DefaultUserName string = "admin"

	RedisOtpDefaultKey string = "otp"

	// Claims
	AuthorizationHeaderKey string = "Authorization"
	UserIdKey              string = "UserId"
	FirstNameKey           string = "FirstName"
	LastNameKey            string = "LastName"	
	RolesKey               string = "Roles"
	ExpireTimeKey          string = "Exp"

	// JWT
	RefreshTokenCookieName string = "refresh_token"
)
