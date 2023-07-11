package repository

import (
	"context"

	"github.com/ak-karimzai/bank-api/internel/db"
	"github.com/google/uuid"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (db.Session, error)
}
