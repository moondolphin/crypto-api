package mocks

import (
	"context"
	"reflect"

	"go.uber.org/mock/gomock"

	"github.com/moondolphin/crypto-api/domain"
)

// --- CoinRepository ---

type MockCoinRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCoinRepositoryMockRecorder
}

type MockCoinRepositoryMockRecorder struct {
	mock *MockCoinRepository
}

func NewMockCoinRepository(ctrl *gomock.Controller) *MockCoinRepository {
	m := &MockCoinRepository{ctrl: ctrl}
	m.recorder = &MockCoinRepositoryMockRecorder{mock: m}
	return m
}

func (m *MockCoinRepository) EXPECT() *MockCoinRepositoryMockRecorder {
	return m.recorder
}

func (m *MockCoinRepository) GetEnabledBySymbol(ctx context.Context, symbol string) (*domain.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnabledBySymbol", ctx, symbol)
	c, _ := ret[0].(*domain.Coin)
	err, _ := ret[1].(error)
	return c, err
}

func (r *MockCoinRepositoryMockRecorder) GetEnabledBySymbol(ctx, symbol any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "GetEnabledBySymbol", reflect.TypeOf((*MockCoinRepository)(nil).GetEnabledBySymbol), ctx, symbol)
}

// --- PriceProvider ---

type MockPriceProvider struct {
	ctrl     *gomock.Controller
	recorder *MockPriceProviderMockRecorder
}

type MockPriceProviderMockRecorder struct {
	mock *MockPriceProvider
}

func NewMockPriceProvider(ctrl *gomock.Controller) *MockPriceProvider {
	m := &MockPriceProvider{ctrl: ctrl}
	m.recorder = &MockPriceProviderMockRecorder{mock: m}
	return m
}

func (m *MockPriceProvider) EXPECT() *MockPriceProviderMockRecorder {
	return m.recorder
}

func (m *MockPriceProvider) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	s, _ := ret[0].(string)
	return s
}

func (r *MockPriceProviderMockRecorder) Name() *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Name", reflect.TypeOf((*MockPriceProvider)(nil).Name))
}

func (m *MockPriceProvider) GetCurrentPrice(ctx context.Context, coin domain.Coin, currency string) (domain.PriceQuote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentPrice", ctx, coin, currency)
	q, _ := ret[0].(domain.PriceQuote)
	err, _ := ret[1].(error)
	return q, err
}

func (r *MockPriceProviderMockRecorder) GetCurrentPrice(ctx, coin, currency any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "GetCurrentPrice", reflect.TypeOf((*MockPriceProvider)(nil).GetCurrentPrice), ctx, coin, currency)
}

// --- PriceProviderRegistry ---

type MockPriceProviderRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockPriceProviderRegistryMockRecorder
}

type MockPriceProviderRegistryMockRecorder struct {
	mock *MockPriceProviderRegistry
}

func NewMockPriceProviderRegistry(ctrl *gomock.Controller) *MockPriceProviderRegistry {
	m := &MockPriceProviderRegistry{ctrl: ctrl}
	m.recorder = &MockPriceProviderRegistryMockRecorder{mock: m}
	return m
}

func (m *MockPriceProviderRegistry) EXPECT() *MockPriceProviderRegistryMockRecorder {
	return m.recorder
}

func (m *MockPriceProviderRegistry) Get(name string) (domain.PriceProvider, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", name)
	p, _ := ret[0].(domain.PriceProvider)
	ok, _ := ret[1].(bool)
	return p, ok
}

func (r *MockPriceProviderRegistryMockRecorder) Get(name any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Get", reflect.TypeOf((*MockPriceProviderRegistry)(nil).Get), name)
}

// ListEnabled mocks base method.
func (m *MockCoinRepository) ListEnabled(ctx context.Context) ([]domain.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEnabled", ctx)
	ret0, _ := ret[0].([]domain.Coin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEnabled indicates an expected call of ListEnabled.
func (mr *MockCoinRepositoryMockRecorder) ListEnabled(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEnabled", reflect.TypeOf((*MockCoinRepository)(nil).ListEnabled), ctx)
}

func (m *MockCoinRepository) GetBySymbol(ctx context.Context, symbol string) (*domain.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySymbol", ctx, symbol)

	c, _ := ret[0].(*domain.Coin)
	err, _ := ret[1].(error)
	return c, err
}

func (m *MockCoinRepository) Upsert(ctx context.Context, c domain.Coin) (*domain.Coin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", ctx, c)

	out, _ := ret[0].(*domain.Coin)
	err, _ := ret[1].(error)
	return out, err
}

func (mr *MockCoinRepositoryMockRecorder) Upsert(ctx, c any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(
		mr.mock,
		"Upsert",
		reflect.TypeOf((*MockCoinRepository)(nil).Upsert),
		ctx, c,
	)
}
