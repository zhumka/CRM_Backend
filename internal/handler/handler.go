package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"crm/internal/pkg/jwtutil"
	"crm/internal/service"
)

// Handler агрегирует сервисы и формирует HTTP-маршруты.
type Handler struct {
	auth       *service.AuthService
	users      *service.UserService
	categories *service.CategoryService
	suppliers  *service.SupplierService
	products   *service.ProductService
	purchases  *service.PurchaseRequestService
	invoices   *service.InvoiceService
	analytics  *service.AnalyticsService
	reports    *service.ReportService
	jwt        *jwtutil.Manager
}

// Services собирает все сервисы для конструктора обработчика.
type Services struct {
	Auth       *service.AuthService
	Users      *service.UserService
	Categories *service.CategoryService
	Suppliers  *service.SupplierService
	Products   *service.ProductService
	Purchases  *service.PurchaseRequestService
	Invoices   *service.InvoiceService
	Analytics  *service.AnalyticsService
	Reports    *service.ReportService
}

func NewHandler(s Services, jwt *jwtutil.Manager) *Handler {
	return &Handler{
		auth:       s.Auth,
		users:      s.Users,
		categories: s.Categories,
		suppliers:  s.Suppliers,
		products:   s.Products,
		purchases:  s.Purchases,
		invoices:   s.Invoices,
		analytics:  s.Analytics,
		reports:    s.Reports,
		jwt:        jwt,
	}
}

// Router конфигурирует все маршруты API.
func (h *Handler) Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), corsMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Swagger UI: http://<host>/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	// Публичные маршруты аутентификации.
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
	}

	// Защищённые маршруты (требуют валидный JWT).
	authorized := api.Group("")
	authorized.Use(authMiddleware(h.jwt))
	{
		authorized.GET("/me", h.me)

		// Управление пользователями — только администратор.
		users := authorized.Group("/users")
		users.Use(adminOnly())
		{
			users.GET("", h.listUsers)
			users.POST("", h.createUser)
			users.GET("/:id", h.getUser)
			users.PUT("/:id", h.updateUser)
			users.DELETE("/:id", h.deleteUser)
		}

		// Категории: чтение — всем, изменение — администратору.
		h.registerCRUD(authorized, "/categories",
			h.listCategories, h.getCategory,
			h.createCategory, h.updateCategory, h.deleteCategory)

		// Поставщики.
		h.registerCRUD(authorized, "/suppliers",
			h.listSuppliers, h.getSupplier,
			h.createSupplier, h.updateSupplier, h.deleteSupplier)

		// Продукты.
		h.registerCRUD(authorized, "/products",
			h.listProducts, h.getProduct,
			h.createProduct, h.updateProduct, h.deleteProduct)

		// Счета-фактуры (создание/изменение — администратор).
		h.registerCRUD(authorized, "/invoices",
			h.listInvoices, h.getInvoice,
			h.createInvoice, h.updateInvoice, h.deleteInvoice)

		// Заявки на закупку — доступны и пользователю, и администратору.
		pr := authorized.Group("/purchase-requests")
		{
			pr.GET("", h.listPurchaseRequests)
			pr.POST("", h.createPurchaseRequest)
			pr.GET("/:id", h.getPurchaseRequest)
			pr.DELETE("/:id", h.deletePurchaseRequest)
			pr.PATCH("/:id/status", adminOnly(), h.updatePurchaseRequestStatus)
		}

		// Аналитика — доступна любому авторизованному пользователю.
		an := authorized.Group("/analytics")
		{
			an.GET("/summary", h.analyticsSummary)
			an.GET("/sales", h.analyticsSales)
			an.GET("/purchase-requests", h.analyticsRequestsByStatus)
			an.GET("/top-products", h.analyticsTopProducts)
		}

		// Отчёты: чтение и экспорт — всем; создание/генерация/удаление — администратору.
		rep := authorized.Group("/reports")
		{
			rep.GET("", h.listReports)
			rep.GET("/:id", h.getReport)
			rep.GET("/:id/export", h.exportReport)
			rep.POST("", adminOnly(), h.createReport)
			rep.POST("/generate-sales", adminOnly(), h.generateSalesReport)
			rep.DELETE("/:id", adminOnly(), h.deleteReport)
		}
	}

	return r
}

// registerCRUD регистрирует стандартный набор: чтение всем, запись администратору.
func (h *Handler) registerCRUD(
	g *gin.RouterGroup, path string,
	list, get, create, update gin.HandlerFunc, del gin.HandlerFunc,
) {
	grp := g.Group(path)
	grp.GET("", list)
	grp.GET("/:id", get)
	grp.POST("", adminOnly(), create)
	grp.PUT("/:id", adminOnly(), update)
	grp.DELETE("/:id", adminOnly(), del)
}
