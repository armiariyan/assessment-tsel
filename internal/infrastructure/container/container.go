package container

import (
	"os"
	"time"

	"github.com/armiariyan/assessment-tsel/internal/config"
	"github.com/armiariyan/assessment-tsel/internal/domain/repositories"
	"github.com/armiariyan/assessment-tsel/internal/infrastructure/mysql"
	"github.com/armiariyan/assessment-tsel/internal/pkg/log"
	"github.com/armiariyan/assessment-tsel/internal/usecase/customers"
	"github.com/armiariyan/assessment-tsel/internal/usecase/healthcheck"
	"github.com/armiariyan/assessment-tsel/internal/usecase/invoices"
	"github.com/armiariyan/assessment-tsel/internal/usecase/items"

	"github.com/armiariyan/bepkg/logger"
)

type Container struct {
	Config             *config.DefaultConfig
	DB                 *config.DB
	Logger             logger.Logger
	HealthCheckService healthcheck.Service
	CustomersService   customers.CustomersService
	ItemsService       items.ItemsService
	InvoicesService    invoices.InvoicesService
}

func (c *Container) Validate() *Container {
	if c.Config == nil {
		panic("Config is nil")
	}
	if c.Logger == nil {
		panic("Logger is nil")
	}
	if c.HealthCheckService == nil {
		panic("HealthCheckService is nil")
	}
	if c.CustomersService == nil {
		panic("CustomersService is nil")
	}
	if c.ItemsService == nil {
		panic("ItemsService is nil")
	}
	if c.InvoicesService == nil {
		panic("InvoicesService is nil")
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

	dbConf := config.DB{
		URI:         config.GetString("mysql.uri"),
		Timeout:     config.GetInt("mysql.timeout"),
		MaxPoolSize: config.GetInt("mysql.maxPool"),
		MinPoolSize: config.GetInt("mysql.minPoolSize"),
		DebugMode:   config.GetBool("debugMode"),
	}

	log.New()

	db := mysql.NewMySQL(dbConf)

	// * Repositories
	customerRepository := repositories.NewCustomersRepository(db)
	itemsRepository := repositories.NewItemsRepository(db)
	invoicesRepository := repositories.NewInvoicesRepository(db)
	invoiceItemsRepository := repositories.NewInvoiceItemsRepository(db)
	invoiceSummaryRepository := repositories.NewInvoiceSummaryRepository(db)

	// * Wrapper

	// * Services
	healthCheckService := healthcheck.NewService().Validate()
	customersService := customers.NewService().
		SetDB(db).
		SetCustomersRepository(customerRepository).
		Validate()
	itemsService := items.NewService().
		SetDB(db).
		SetItemsRepository(itemsRepository).
		Validate()
	invoicesService := invoices.NewService().
		SetDB(db).
		SetCustomersRepository(customerRepository).
		SetInvoicesRepository(invoicesRepository).
		SetItemsRepository(itemsRepository).
		SetInvoiceItemsRepository(invoiceItemsRepository).
		SetInvoiceSummaryRepository(invoiceSummaryRepository).
		Validate()

	// * Brokers

	// * Workers

	container := &Container{
		Config:             defConfig,
		Logger:             defLogger,
		HealthCheckService: healthCheckService,
		CustomersService:   customersService,
		ItemsService:       itemsService,
		InvoicesService:    invoicesService,
	}
	container.Validate()
	return container

}
