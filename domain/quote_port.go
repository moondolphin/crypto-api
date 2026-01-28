package domain

import "context"

type QuoteRepository interface {
	Insert(ctx context.Context, q Quote) error
}
