package app

import (
	"context"
	"errors"
	"log"
	"math"
	"time"

	"github.com/moondolphin/crypto-api/domain"
)

var ErrCooldownActive = errors.New("cooldown_active")

type ManualRefreshWithCooldownUseCase struct {
	RefreshUC   RefreshQuotesUseCase
	ControlRepo domain.RefreshControlRepository
	Now         func() time.Time
	Cooldown    time.Duration // ej: 20 * time.Minute
}

type ManualRefreshWithCooldownOutput struct {
	RefreshQuotesOutput
	RetryAfterSeconds int `json:"retry_after_seconds,omitempty"`
}

func (uc ManualRefreshWithCooldownUseCase) Execute(ctx context.Context) (ManualRefreshWithCooldownOutput, error) {
	nowFn := uc.Now
	if nowFn == nil {
		nowFn = time.Now
	}
	now := nowFn().UTC()

	cd := uc.Cooldown
	if cd <= 0 {
		cd = 20 * time.Minute
	}

	last, ok, err := uc.ControlRepo.GetLastManualRefresh(ctx)
	if err != nil {
		return ManualRefreshWithCooldownOutput{}, err
	}

	if ok {
		next := last.Add(cd)
		if now.Before(next) {
			//remain := time.Until(next)
			remain := next.Sub(now)
			//sec := int(remain.Seconds())
			sec := int(math.Ceil(remain.Seconds()))

			if sec < 1 {
				sec = 1
			}
			return ManualRefreshWithCooldownOutput{
				RetryAfterSeconds: sec,
			}, ErrCooldownActive
		}
	}

	// Ejecuta refresh real
	out, err := uc.RefreshUC.Execute(ctx)
	if err != nil {
		return ManualRefreshWithCooldownOutput{}, err
	}

	// Marca “último manual refresh”
	//_ = uc.ControlRepo.SetLastManualRefresh(ctx, now)
	if err := uc.ControlRepo.SetLastManualRefresh(ctx, now); err != nil {
		// opcional: loguear, pero NO fallar el refresh
		log.Printf("Warning: failed to set last manual refresh time: %v", err)
	}

	return ManualRefreshWithCooldownOutput{RefreshQuotesOutput: out}, nil
}
