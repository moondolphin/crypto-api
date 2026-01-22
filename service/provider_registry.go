package service

import (
	"strings"

	"github.com/moondolphin/crypto-api/domain"
)

type ProviderRegistry struct {
	providers map[string]domain.PriceProvider
}

func NewProviderRegistry(ps ...domain.PriceProvider) *ProviderRegistry {
	m := make(map[string]domain.PriceProvider, len(ps))
	for _, p := range ps {
		m[strings.ToLower(strings.TrimSpace(p.Name()))] = p
	}
	return &ProviderRegistry{providers: m}
}

func (r *ProviderRegistry) Get(name string) (domain.PriceProvider, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	p, ok := r.providers[key]
	return p, ok
}
