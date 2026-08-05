package authdomain

import "time"

type AuthClaims struct {
	Sub int		// subject	
	Exp time.Time	// expiration time
	Iat time.Time	// issued at
	Username      string
	Admin         bool
	ResourceRoles map[string]any
}

