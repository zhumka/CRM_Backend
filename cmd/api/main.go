package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"crm/internal/config"
	"crm/internal/handler"
	"crm/internal/pkg/jwtutil"
	"crm/internal/repository"
	"crm/internal/seed"
	"crm/internal/service"

	_ "crm/docs" // сгенерированная Swagger-документация
)

// @title           CRM API
// @version         1.0
// @description     Бэкенд CRM-системы для предприятия оптово-розничной торговли:
// @description     управление продуктами, категориями, поставщиками, заявками на закупку,
// @description     счетами-фактурами, аналитикой и отчётами.
// @termsOfService  http://swagger.io/terms/

// @contact.name   CRM Backend
// @license.name   MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Введите токен в формате: Bearer <token>

// @tag.name  auth
// @tag.description Аутентификация и регистрация
// @tag.name  users
// @tag.description Управление пользователями (только администратор)
// @tag.name  categories
// @tag.description Категории продукции
// @tag.name  suppliers
// @tag.description Поставщики
// @tag.name  products
// @tag.description Продукты
// @tag.name  purchase-requests
// @tag.description Заявки на закупку
// @tag.name  invoices
// @tag.description Счета-фактуры
// @tag.name  analytics
// @tag.description Аналитика и KPI
// @tag.name  reports
// @tag.description Отчёты (формирование и экспорт)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	gin.SetMode(cfg.GinMode)

	// Слой данных.
	db, err := repository.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("migrations applied")

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	supplierRepo := repository.NewSupplierRepository(db)
	productRepo := repository.NewProductRepository(db)
	purchaseRepo := repository.NewPurchaseRequestRepository(db)
	invoiceRepo := repository.NewInvoiceRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	reportRepo := repository.NewReportRepository(db)
	taxRepo := repository.NewTaxRepository(db)
	unitRepo := repository.NewUnitRepository(db)
	saleRepo := repository.NewSaleRepository(db)

	// Слой бизнес-логики.
	jwtManager := jwtutil.NewManager(cfg.JWTSecret, cfg.JWTTTL)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)

	services := handler.Services{
		Auth:       service.NewAuthService(userRepo, jwtManager),
		Users:      userService,
		Categories: service.NewCategoryService(categoryRepo),
		Suppliers:  service.NewSupplierService(supplierRepo),
		Products:   productService,
		Purchases:  service.NewPurchaseRequestService(purchaseRepo, productRepo),
		Invoices:   service.NewInvoiceService(invoiceRepo),
		Analytics:  analyticsService,
		Reports:    service.NewReportService(reportRepo, analyticsService),
		Taxes:      service.NewTaxService(taxRepo),
		Units:      service.NewUnitService(unitRepo),
		Sales:      service.NewSaleService(saleRepo),
	}

	// Создаём администратора по умолчанию при первом запуске.
	adminUser := getEnv("ADMIN_USERNAME", "admin")
	adminPass := getEnv("ADMIN_PASSWORD", "admin123")
	if err := userService.EnsureAdmin(context.Background(), adminUser, adminPass); err != nil {
		log.Fatalf("seed admin: %v", err)
	}

	// Наполнение демо-данными (идемпотентно). Отключается SEED_DEMO=false.
	if getEnv("SEED_DEMO", "true") == "true" {
		seeder := seed.NewSeeder(userRepo, categoryRepo, supplierRepo, productRepo, purchaseRepo, invoiceRepo, taxRepo, unitRepo, saleRepo)
		if err := seeder.Run(context.Background()); err != nil {
			log.Fatalf("seed demo: %v", err)
		}
	}

	// Транспортный слой.
	h := handler.NewHandler(services, jwtManager)
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      h.Router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
