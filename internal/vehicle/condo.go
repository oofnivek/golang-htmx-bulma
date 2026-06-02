package vehicle

import (
	"time"

	"golang-htmx-bulma/internal/pkg/status"
)

type Condo struct {
	ID         int64         `json:"id"`
	Name       string        `json:"name"`
	Status     status.Status `json:"status"`
	McstNumber string        `json:"mcst_number"`
	McstEmail  string        `json:"mcst_email"`
	Address    string        `json:"address"`
	CreatedBy  string        `json:"created_by"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedBy  *string       `json:"updated_by"`
	UpdatedAt  *time.Time    `json:"updated_at"`
}
