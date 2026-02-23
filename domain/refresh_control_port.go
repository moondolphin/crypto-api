package domain

import (
	"context"
	"time"
)

type RefreshControlRepository interface {
	// devuelve el último momento en que se ejecutó el refresh manual.
	// ok=false si nunca se ejecutó.
	GetLastManualRefresh(ctx context.Context) (t time.Time, ok bool, err error)

	// guarda el momento del refresh manual
	SetLastManualRefresh(ctx context.Context, t time.Time) error
}
