package domain

//go:generate echo Generating mocks for refresh_control_port.go
//go:generate go run go.uber.org/mock/mockgen@v0.5.0 -source=refresh_control_port.go -destination=../test/mocks/refresh_control_port_mock.go -package=mocks

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
