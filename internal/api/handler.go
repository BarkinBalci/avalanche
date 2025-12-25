package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/BarkinBalci/event-analytics-service/docs"
	"github.com/BarkinBalci/event-analytics-service/internal/models"
	"github.com/BarkinBalci/event-analytics-service/internal/service"
)

type Handler struct {
	eventService *service.EventService
	router       *gin.Engine
}

func NewHandler(eventService *service.EventService) *Handler {
	h := &Handler{
		eventService: eventService,
		router:       gin.Default(),
	}

	h.registerRoutes()

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) registerRoutes() {
	h.router.GET("/health", h.healthCheck)
	h.router.POST("/events", h.publishEvent)
	h.router.POST("/events/bulk", h.publishEventsBulk)
	h.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// healthCheck handles health check requests
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) healthCheck(c *gin.Context) {
	// TODO: add a more sophisticated health check
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// publishEvent handles POST /events
// @Summary Publish a single event
// @Description Publish a single analytics event to the queue
// @Tags events
// @Accept json
// @Produce json
// @Param event body models.PublishEventRequest true "Event data"
// @Success 202 {object} models.PublishEventResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events [post]
func (h *Handler) publishEvent(c *gin.Context) {
	var req models.PublishEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	eventID, err := h.eventService.ProcessEvent(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, models.PublishEventResponse{
		EventID: eventID,
		Status:  "accepted",
	})
}

// publishEventsBulk handles POST /events/bulk
// @Summary Publish multiple events
// @Description Publish multiple analytics events in bulk to the queue
// @Tags events
// @Accept json
// @Produce json
// @Param events body models.PublishEventsBulkRequest true "Bulk events data"
// @Success 202 {object} models.PublishBulkEventsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/bulk [post]
func (h *Handler) publishEventsBulk(c *gin.Context) {
	var bulkRequest models.PublishEventsBulkRequest

	if err := c.ShouldBindJSON(&bulkRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	eventIDs, errors, err := h.eventService.ProcessBulkEvents(bulkRequest.Events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	accepted := len(eventIDs)
	rejected := len(errors)

	c.JSON(http.StatusAccepted, models.PublishBulkEventsResponse{
		Accepted: accepted,
		Rejected: rejected,
		EventIDs: eventIDs,
		Errors:   errors,
	})
}
