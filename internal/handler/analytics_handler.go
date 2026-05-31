package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// analyticsSummary godoc
// @Summary  Сводные KPI
// @Tags     analytics
// @Produce  json
// @Security BearerAuth
// @Success  200  {object}  model.Summary
// @Router   /analytics/summary [get]
func (h *Handler) analyticsSummary(c *gin.Context) {
	res, err := h.analytics.Summary(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// analyticsSales godoc
// @Summary  Финансовая аналитика за период
// @Tags     analytics
// @Produce  json
// @Security BearerAuth
// @Param    from  query     string  false  "Дата с (YYYY-MM-DD)"
// @Param    to    query     string  false  "Дата по (YYYY-MM-DD)"
// @Success  200   {object}  model.SalesAnalytics
// @Failure  400   {object}  errorResponse
// @Router   /analytics/sales [get]
func (h *Handler) analyticsSales(c *gin.Context) {
	res, err := h.analytics.Sales(c.Request.Context(), c.Query("from"), c.Query("to"))
	if err != nil {
		badRequest(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// analyticsRequestsByStatus godoc
// @Summary  Количество заявок по статусам
// @Tags     analytics
// @Produce  json
// @Security BearerAuth
// @Success  200  {array}  model.StatusCount
// @Router   /analytics/purchase-requests [get]
func (h *Handler) analyticsRequestsByStatus(c *gin.Context) {
	res, err := h.analytics.RequestsByStatus(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// analyticsTopProducts godoc
// @Summary  Самые востребованные продукты
// @Tags     analytics
// @Produce  json
// @Security BearerAuth
// @Param    limit  query     int  false  "Сколько вернуть (по умолчанию 10)"
// @Success  200    {array}   model.TopProduct
// @Router   /analytics/top-products [get]
func (h *Handler) analyticsTopProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	res, err := h.analytics.TopProducts(c.Request.Context(), limit)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}
