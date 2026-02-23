package mocks

import (
	"context"
	"reflect"
	"time"

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

func (mr *MockCoinRepositoryMockRecorder) GetBySymbol(ctx, symbol any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySymbol", reflect.TypeOf((*MockCoinRepository)(nil).GetBySymbol), ctx, symbol)
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

// --- QuoteRepository ---

type MockQuoteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockQuoteRepositoryMockRecorder
}

type MockQuoteRepositoryMockRecorder struct {
	mock *MockQuoteRepository
}

func NewMockQuoteRepository(ctrl *gomock.Controller) *MockQuoteRepository {
	m := &MockQuoteRepository{ctrl: ctrl}
	m.recorder = &MockQuoteRepositoryMockRecorder{mock: m}
	return m
}

func (m *MockQuoteRepository) EXPECT() *MockQuoteRepositoryMockRecorder {
	return m.recorder
}

func (m *MockQuoteRepository) GetLatest(ctx context.Context, symbol, provider, currency string) (*domain.PriceQuote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatest", ctx, symbol, provider, currency)
	q, _ := ret[0].(*domain.PriceQuote)
	err, _ := ret[1].(error)
	return q, err
}

func (r *MockQuoteRepositoryMockRecorder) GetLatest(ctx, symbol, provider, currency any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "GetLatest", reflect.TypeOf((*MockQuoteRepository)(nil).GetLatest), ctx, symbol, provider, currency)
}

func (m *MockQuoteRepository) ListFilter(ctx context.Context, f domain.QuoteFilter) ([]domain.Quote, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFilter", ctx, f)
	quotes, _ := ret[0].([]domain.Quote)
	total, _ := ret[1].(int)
	err, _ := ret[2].(error)
	return quotes, total, err
}

func (r *MockQuoteRepositoryMockRecorder) ListFilter(ctx, f any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "ListFilter", reflect.TypeOf((*MockQuoteRepository)(nil).ListFilter), ctx, f)
}

func (m *MockQuoteRepository) Insert(ctx context.Context, q domain.Quote) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, q)
	err, _ := ret[0].(error)
	return err
}

func (r *MockQuoteRepositoryMockRecorder) Insert(ctx, q any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Insert", reflect.TypeOf((*MockQuoteRepository)(nil).Insert), ctx, q)
}

// --- UserRepository ---

type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	m := &MockUserRepository{ctrl: ctrl}
	m.recorder = &MockUserRepositoryMockRecorder{mock: m}
	return m
}

func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsByEmail", ctx, email)
	exists, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return exists, err
}

func (r *MockUserRepositoryMockRecorder) ExistsByEmail(ctx, email any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "ExistsByEmail", reflect.TypeOf((*MockUserRepository)(nil).ExistsByEmail), ctx, email)
}

func (m *MockUserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, u)
	user, _ := ret[0].(domain.User)
	err, _ := ret[1].(error)
	return user, err
}

func (r *MockUserRepositoryMockRecorder) Create(ctx, u any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), ctx, u)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	user, _ := ret[0].(*domain.User)
	err, _ := ret[1].(error)
	return user, err
}

func (r *MockUserRepositoryMockRecorder) FindByEmail(ctx, email any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), ctx, email)
}

// --- PasswordHasher ---

type MockPasswordHasher struct {
	ctrl     *gomock.Controller
	recorder *MockPasswordHasherMockRecorder
}

type MockPasswordHasherMockRecorder struct {
	mock *MockPasswordHasher
}

func NewMockPasswordHasher(ctrl *gomock.Controller) *MockPasswordHasher {
	m := &MockPasswordHasher{ctrl: ctrl}
	m.recorder = &MockPasswordHasherMockRecorder{mock: m}
	return m
}

func (m *MockPasswordHasher) EXPECT() *MockPasswordHasherMockRecorder {
	return m.recorder
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", password)
	hash, _ := ret[0].(string)
	err, _ := ret[1].(error)
	return hash, err
}

func (r *MockPasswordHasherMockRecorder) Hash(password any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Hash", reflect.TypeOf((*MockPasswordHasher)(nil).Hash), password)
}

func (m *MockPasswordHasher) Compare(hash, password string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compare", hash, password)
	ok, _ := ret[0].(bool)
	err, _ := ret[1].(error)
	return ok, err
}

func (r *MockPasswordHasherMockRecorder) Compare(hash, password any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Compare", reflect.TypeOf((*MockPasswordHasher)(nil).Compare), hash, password)
}

// --- TokenService ---

type MockTokenService struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceMockRecorder
}

type MockTokenServiceMockRecorder struct {
	mock *MockTokenService
}

func NewMockTokenService(ctrl *gomock.Controller) *MockTokenService {
	m := &MockTokenService{ctrl: ctrl}
	m.recorder = &MockTokenServiceMockRecorder{mock: m}
	return m
}

func (m *MockTokenService) EXPECT() *MockTokenServiceMockRecorder {
	return m.recorder
}

func (m *MockTokenService) Generate(userID int64, email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", userID, email)
	token, _ := ret[0].(string)
	err, _ := ret[1].(error)
	return token, err
}

func (r *MockTokenServiceMockRecorder) Generate(userID, email any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "Generate", reflect.TypeOf((*MockTokenService)(nil).Generate), userID, email)
}

// --- RefreshControlRepository ---

type MockRefreshControlRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRefreshControlRepositoryMockRecorder
}

type MockRefreshControlRepositoryMockRecorder struct {
	mock *MockRefreshControlRepository
}

func NewMockRefreshControlRepository(ctrl *gomock.Controller) *MockRefreshControlRepository {
	m := &MockRefreshControlRepository{ctrl: ctrl}
	m.recorder = &MockRefreshControlRepositoryMockRecorder{mock: m}
	return m
}

func (m *MockRefreshControlRepository) EXPECT() *MockRefreshControlRepositoryMockRecorder {
	return m.recorder
}

func (m *MockRefreshControlRepository) GetLastManualRefresh(ctx context.Context) (time.Time, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastManualRefresh", ctx)
	t, _ := ret[0].(time.Time)
	ok, _ := ret[1].(bool)
	err, _ := ret[2].(error)
	return t, ok, err
}

func (r *MockRefreshControlRepositoryMockRecorder) GetLastManualRefresh(ctx any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "GetLastManualRefresh", reflect.TypeOf((*MockRefreshControlRepository)(nil).GetLastManualRefresh), ctx)
}

func (m *MockRefreshControlRepository) SetLastManualRefresh(ctx context.Context, t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLastManualRefresh", ctx, t)
	err, _ := ret[0].(error)
	return err
}

func (r *MockRefreshControlRepositoryMockRecorder) SetLastManualRefresh(ctx, t any) *gomock.Call {
	r.mock.ctrl.T.Helper()
	return r.mock.ctrl.RecordCallWithMethodType(r.mock, "SetLastManualRefresh", reflect.TypeOf((*MockRefreshControlRepository)(nil).SetLastManualRefresh), ctx, t)
}
