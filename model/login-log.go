package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginLog struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	LoginAt time.Time
}
