package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"validation_error"`
	Message string `json:"message,omitempty" example:"event_name is required"`
}

// PublishEventResponse represents a successful event ingestion response
type PublishEventResponse struct {
	EventID string `json:"event_id" example:"evt_1a2b3c4d5e6f"`
	Status  string `json:"status" example:"accepted"`
}

// PublishBulkEventsResponse represents a successful bulk event ingestion response
type PublishBulkEventsResponse struct {
	Accepted int      `json:"accepted" example:"5"`
	Rejected int      `json:"rejected" example:"0"`
	EventIDs []string `json:"event_ids,omitempty" example:"evt_1,evt_2,evt_3"`
	Errors   []string `json:"errors,omitempty" example:"validation error on event 3"`
}
