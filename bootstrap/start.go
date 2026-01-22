package bootstrap

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	httpapi "github.com/moondolphin/crypto-api/adapters/primary/httpapi"
	mysqlrepo "github.com/moondolphin/crypto-api/adapters/secondary/persistence/mysql"
	"github.com/moondolphin/crypto-api/adapters/secondary/providers"
	"github.com/moondolphin/crypto-api/app"
	"github.com/moondolphin/crypto-api/config"
	"github.com/moondolphin/crypto-api/service"

	"github.com/moondolphin/crypto-api/domain"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start() (*gin.Engine, error) {

	_ = domain.PriceQuote{}

	dsn, err := config.MySQLDSN()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	coinRepo := mysqlrepo.NewMySQLCoinRepository(db)
	reg := service.NewProviderRegistry(
		providers.NewBinanceProvider(),
		providers.NewCoinGeckoProvider(),
	)

	uc := app.GetCurrentPriceUseCase{
		CoinRepo:  coinRepo,
		Providers: reg,
		Now:       time.Now,
	}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// @Summary Consultar precio actual de criptomoneda
	// @Description Consulta el precio actual de una criptomoneda habilitada
	// @Tags Crypto
	// @Param symbol query string true "Símbolo (BTC, ETH)"
	// @Param currency query string true "Moneda (USD, USDT)"
	// @Param provider query string true "Proveedor (binance)"
	// @Success 200 {object} domain.PriceQuote
	// @Failure 400 {object} map[string]string
	// @Failure 404 {object} map[string]string
	// @Failure 503 {object} map[string]string
	// @Router /api/v1/crypto/price [get]

	r.GET("/api/v1/crypto/price", httpapi.GetCurrentPriceHandler{UC: uc}.Handle)

	return r, nil
}
