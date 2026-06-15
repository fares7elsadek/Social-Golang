package domain

import "time"


type Role string


const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleUser     Role = "user"
	RoleReadOnly Role = "readonly"
)

var RoleMapping = map[string]int{
	string(RoleAdmin):    1,
	string(RoleManager):  2,
	string(RoleUser):     3,
	string(RoleReadOnly): 4,
}


type User struct {
	ID        int       "json:\"id\""
	Username  string    "json:\"username\""
	Email     string    "json:\"email\""
	Roles     []Role    `json:"roles"`
	IsActive   bool     `json:"is_active"`
	CreatedAt time.Time "json:\"created_at\""
	UpdatedAt time.Time "json:\"updated_at\""
	Password  string    "json:\"-\""
}


func (u *User) HasRule(r Role) bool {
	for _,role := range u.Roles {
		if role == r {
			return true
		}
	}
	return false
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` 
}

type Claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Roles  []Role `json:"roles"`
}
