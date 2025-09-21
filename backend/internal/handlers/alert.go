package handlers

import (
	"net/http"
	"strconv"

	"news-to-text/internal/middleware"
	"news-to-text/internal/models"
	"news-to-text/internal/services"

	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	alertService services.AlertService
	authService  services.AuthService
}

func NewAlertHandler(alertService services.AlertService, authService services.AuthService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
		authService:  authService,
	}
}

// GetAlerts godoc
// @Summary Get user alerts
// @Description Get all alerts for the authenticated user
// @Tags alerts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.AlertResponse
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alerts [get]
func (h *AlertHandler) GetAlerts(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	alerts, err := h.alertService.GetAlerts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alerts"})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// CreateAlert godoc
// @Summary Create a new alert
// @Description Create a new news alert for the authenticated user
// @Tags alerts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param alert body models.AlertCreateRequest true "Alert data"
// @Success 201 {object} models.AlertResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alerts [post]
func (h *AlertHandler) CreateAlert(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AlertCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.alertService.CreateAlert(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create alert"})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// UpdateAlert godoc
// @Summary Update an alert
// @Description Update an existing alert for the authenticated user
// @Tags alerts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Alert ID"
// @Param alert body models.AlertUpdateRequest true "Alert update data"
// @Success 200 {object} models.AlertResponse
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Alert not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alerts/{id} [put]
func (h *AlertHandler) UpdateAlert(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	alertIDStr := c.Param("id")
	alertID, err := strconv.ParseUint(alertIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	var req models.AlertUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.alertService.UpdateAlert(userID, uint(alertID), &req)
	if err != nil {
		if err.Error() == "alert not found" || err.Error() == "unauthorized access to alert" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update alert"})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// DeleteAlert godoc
// @Summary Delete an alert
// @Description Delete an existing alert for the authenticated user
// @Tags alerts
// @Security BearerAuth
// @Param id path int true "Alert ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Alert not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alerts/{id} [delete]
func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	alertIDStr := c.Param("id")
	alertID, err := strconv.ParseUint(alertIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	err = h.alertService.DeleteAlert(userID, uint(alertID))
	if err != nil {
		if err.Error() == "alert not found" || err.Error() == "unauthorized access to alert" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete alert"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAlertHistory godoc
// @Summary Get alert history
// @Description Get alert history for the authenticated user
// @Tags alerts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.AlertHistory
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alerts/history [get]
func (h *AlertHandler) GetAlertHistory(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	history, err := h.alertService.GetAlertHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get alert history"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// TestAlert godoc
// @Summary Test an alert
// @Description Send a test notification for an alert
// @Tags alerts
// @Security BearerAuth
// @Param id path int true "Alert ID"
// @Success 200 {object} map[string]interface{} "success message"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Alert not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alerts/test [post]
func (h *AlertHandler) TestAlert(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		AlertID uint `json:"alert_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.alertService.TestAlert(userID, req.AlertID)
	if err != nil {
		if err.Error() == "alert not found" || err.Error() == "unauthorized access to alert" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to test alert"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test alert sent successfully"})
}