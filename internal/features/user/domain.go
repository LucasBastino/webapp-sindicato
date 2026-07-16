package user

import "time"

type User struct {
	ID              	int   				`db:"id_user"`

	Username        	string 				`db:"username"`
	PasswordHash    	string 				`db:"password_hash"`
	
	Admin				bool  				`db:"admin"`
	ResourceRoles 		map[string]any
	ResourceRolesJSON 	[]byte 				`db:"resource_roles"`

	CreatedAt       	time.Time			`db:"created_at"`
}
