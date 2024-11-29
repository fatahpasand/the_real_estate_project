// internal/core/domain/user.go
package domain

import "time"

type User struct {
    ID            uint      `json:"id"`
    Email         string    `json:"email"`
    Password      string    `json:"-"`
    Name          string    `json:"name"`
    Verified      bool      `json:"verified"`
    GoogleID      string    `json:"google_id,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type OAuthProfile struct {
    ID        string
    Email     string
    Name      string
    Provider  string
}