package model

import "time"

type Bastion struct {
	ID        int        `json:"id,omitempty"`
	UserID    int        `json:"user_id,omitempty"`
	VmID      int        `json:"vm_id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
