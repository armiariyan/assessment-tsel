package container

import (
	"os"
	"time"

	"github.com/armiariyan/assessment-tsel/internal/config"
	"github.com/armiariyan/assessment-tsel/internal/domain/repositories"
	"github.com/armiariyan/assessment-tsel/internal/infrastructure/postgresql"
	"github.com/armiariyan/assessment-tsel/internal/pkg/log"
	"github.com/armiariyan/assessment-tsel/internal/usecase/healthcheck"
	"github.com/armiariyan/assessment-tsel/internal/usecase/products"
	"gorm.io/gorm"

	"github.com/armiariyan/bepkg/logger"
)

type Container struct {
	Config             *config.DefaultConfig
	PostgresqlDB       *config.PostgresqlDB
	ProductsDB         *gorm.DB
	Logger             logger.Logger
	HealthCheckService healthcheck.Service
	ProductService     products.Service
}

func (c *Container) Validate() *Container {
	if c.Config == nil {
		panic("Config is nil")
	}
	if c.ProductsDB == nil {
		panic("ProductsDB is nil")
	}
	if c.Logger == nil {
		panic("Logger is nil")
	}
	if c.HealthCheckService == nil {
		panic("HealthCheckService is nil")
	}
	if c.ProductService == nil {
		panic("ProductService is nil")
	}
	return c
}

func New() *Container {
	config.Load(os.Getenv("env"), ".env")

	fileLoc := config.GetString("logger.fileLocation")
	tdrFileLoc := config.GetString("logger.fileTdrLocation")
	maxAge := time.Duration(config.GetInt("logger.fileMaxAge"))
	stdOut := config.GetBool("logger.stdout")

	defLogger := logger.New(logger.Options{
		FileLocation:    fileLoc,
		FileTdrLocation: tdrFileLoc,
		FileMaxAge:      maxAge,
		Stdout:          stdOut,
	})

	defConfig := &config.DefaultConfig{
		Apps: config.Apps{
			Name:     config.GetString("app.name"),
			Address:  config.GetString("address"),
			HttpPort: config.GetString("port"),
		},
	}

	psqlConfig := &config.PostgresqlDB{
		Host:     config.GetString("postgresql.products.host"),
		User:     config.GetString("postgresql.products.user"),
		Password: config.GetString("postgresql.products.password"),
		Name:     config.GetString("postgresql.products.db"),
		Port:     config.GetInt("postgresql.products.port"),
		SSLMode:  config.GetString("postgresql.products.ssl"),
		Schema:   config.GetString("postgresql.products.schema"),
		Debug:    config.GetBool("postgresql.products.debug"),
	}

	log.New()

	productsDB := postgresql.NewDB(*psqlConfig)

	// * Repositories
	productRepository := repositories.NewProductsRepository(productsDB)

	// * Wrapper

	// * Services
	healthCheckService := healthcheck.NewService().
		Validate()

	productService := products.NewService().
		SetProductsRepository(productRepository).
		Validate()

	// * Brokers

	// * Workers

	container := &Container{
		Config:             defConfig,
		Logger:             defLogger,
		PostgresqlDB:       psqlConfig,
		ProductsDB:         productsDB,
		HealthCheckService: healthCheckService,
		ProductService:     productService,
	}
	container.Validate()
	return container

}
