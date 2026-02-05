package bootstrap

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	httpapi "github.com/moondolphin/crypto-api/adapters/primary/httpapi"
	mysqlrepo "github.com/moondolphin/crypto-api/adapters/secondary/persistence/mysql"
	"github.com/moondolphin/crypto-api/adapters/secondary/providers"
	"github.com/moondolphin/crypto-api/adapters/secondary/security"
	"github.com/moondolphin/crypto-api/app"
	"github.com/moondolphin/crypto-api/config"
	"github.com/moondolphin/crypto-api/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start() (*gin.Engine, error) {
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

	// deps
	userRepo := mysqlrepo.NewMySQLUserRepository(db)
	hasher := security.NewBcryptHasher(0)

	jwtSecret, err := config.JWTSecret()
	if err != nil {
		return nil, err
	}
	jwtTTL := config.JWTTTL()
	jwtSvc := security.NewJWTService(jwtSecret, jwtTTL, "crypto-api")

	// use cases
	registerUC := app.RegisterUserUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Now:      time.Now,
	}

	loginUC := app.LoginUseCase{
		UserRepo: userRepo,
		Hasher:   hasher,
		Tokens:   jwtSvc,
		Now:      time.Now,
		TTL:      jwtTTL,
	}

	coinRepo := mysqlrepo.NewMySQLCoinRepository(db)
	reg := service.NewProviderRegistry(
		providers.NewBinanceProvider(),
		providers.NewCoinGeckoProvider(),
	)

	// router
	r := gin.Default()

	quoteRepo := mysqlrepo.NewMySQLQuoteRepository(db)

	lastPriceUC := app.GetLastPriceUseCase{
		CoinRepo:  coinRepo,
		QuoteRepo: quoteRepo,
	}

	refreshUC := app.RefreshQuotesUseCase{
		CoinRepo:  coinRepo,
		QuoteRepo: quoteRepo,
		Providers: reg,
		Now:       time.Now,
		ProviderFX: map[string]string{
			"binance":   "USDT",
			"coingecko": "USD",
		},
	}

	searchQuotesUC := app.SearchQuotesUseCase{
		Repo: quoteRepo,
	}

	createCoinUC := app.CreateCoinUseCase{
		CoinRepo: coinRepo,
	}

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// públicos
	r.POST("/api/v1/auth/register", httpapi.RegisterUserHandler{UC: registerUC}.Handle)
	r.POST("/api/v1/auth/login", httpapi.LoginHandler{UC: loginUC}.Handle)
	r.GET("/api/v1/crypto/price", httpapi.GetCurrentPriceHandler{UC: lastPriceUC}.Handle)
	r.POST("/api/v1/job/refresh", httpapi.RefreshHandler{UC: refreshUC}.Handle)
	r.GET("/api/v1/quotes", httpapi.SearchQuotesHandler{UC: searchQuotesUC}.Handle)

	// privados
	auth := r.Group("/api/v1")
	auth.Use(httpapi.AuthRequired(jwtSecret))
	auth.POST("/coins", httpapi.CreateCoinHandler{UC: createCoinUC}.Handle)

	auth.GET("/me", func(c *gin.Context) {
		v, _ := c.Get("auth")
		c.JSON(200, v)
	})

	return r, nil
}
