package authdomain

import "time"

type AuthClaims struct {
	Sub int		// subject	
	Exp time.Time	// expiration time
	Iat time.Time	// issued at
	Admin      bool
	ResourceRoles map[string]any
}

