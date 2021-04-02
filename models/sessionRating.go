package models
import(
"github.com/satori/go.uuid"
"time")

// CREATE TABLE session_ratings (user_id UUID, session_id UUID, rating numeric , created_at timestamptz, updated_at timestamptz null, 
// 	PRIMARY KEY (user_id, session_id));

type SessionRating struct {
	UserID uuid.UUID `gorm:"primary_key"`
	SessionID uuid.UUID `gorm:"primary_key"`
	Rating    float64
	CreatedAt time.Time  `json:"created_at"`
    UpdatedAt *time.Time  `json:"update_at"`
}

