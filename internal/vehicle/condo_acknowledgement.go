package vehicle

import "time"

type CondoAcknowledgement struct {
	UserID    int64      `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
