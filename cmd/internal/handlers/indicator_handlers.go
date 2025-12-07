package handlers

import (
	"encoding/json"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"

	"github.com/google/uuid"
)

type IndicatorHandlers struct {
	indicatorUseCase *usecases.IndicatorUseCase
	causeUseCase     *usecases.CauseUseCase
	actionUseCase    *usecases.ActionUseCase
}

func NewIndicatorHandlers(
	indicatorUseCase *usecases.IndicatorUseCase,
	causeUseCase *usecases.CauseUseCase,
	actionUseCase *usecases.ActionUseCase,
) *IndicatorHandlers {
	return &IndicatorHandlers{
		indicatorUseCase: indicatorUseCase,
		causeUseCase:     causeUseCase,
		actionUseCase:    actionUseCase,
	}
}

type CreateIndicatorRequest struct {
	IterationID uuid.UUID `json:"iteration_id"`
}

type CreateCauseRequest struct {
	IndicatorID       uuid.UUID `json:"indicator_id"`
	Metric            string    `json:"metric"`
	Description       string    `json:"description"`
	ProductivityLevel string    `json:"productivity_level"`
}

type CreateActionRequest struct {
	IndicatorID uuid.UUID `json:"indicator_id"`
	CauseID     uuid.UUID `json:"cause_id"`
	Description string    `json:"description"`
}

// Get handles GET /indicators
// @Summary Get indicator
// @Description Get indicator with causes and actions for a specific iteration
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param iteration_id query string true "Iteration ID" format(uuid)
// @Success 200 {object} models.Indicator "Indicator details"
// @Failure 400 {string} string "Invalid iteration_id"
// @Failure 404 {string} string "Indicator not found"
// @Router /indicators [get]
func (h *IndicatorHandlers) Get(w http.ResponseWriter, r *http.Request) {
	iterationIDStr := r.URL.Query().Get("iteration_id")
	if iterationIDStr == "" {
		http.Error(w, "iteration_id is required", http.StatusBadRequest)
		return
	}

	iterationID, err := uuid.Parse(iterationIDStr)
	if err != nil {
		http.Error(w, "Invalid iteration_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	indicator, err := h.indicatorUseCase.Get(ctx, iterationID)
	if err != nil {
		http.Error(w, "Indicator not found", http.StatusNotFound)
		return
	}

	// Load causes and actions
	causes, _ := h.causeUseCase.Get(ctx, indicator.ID)
	actions, _ := h.actionUseCase.Get(ctx, indicator.ID)

	indicator.Causes = causes
	indicator.Actions = actions

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(indicator)
}

// Create handles POST /indicators
// @Summary Create a new indicator
// @Description Create a new indicator for an iteration
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param indicator body CreateIndicatorRequest true "Indicator data"
// @Success 201 {object} map[string]interface{} "Indicator created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create indicator"
// @Router /indicators [post]
func (h *IndicatorHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateIndicatorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newIndicator := models.Indicator{
		IterationID: req.IterationID,
	}

	ctx := r.Context()
	indicatorID, err := h.indicatorUseCase.Create(ctx, newIndicator)
	if err != nil {
		http.Error(w, "Failed to create indicator", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":           indicatorID,
		"iteration_id": req.IterationID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// CreateCause handles POST /indicators/causes
// @Summary Create a new cause
// @Description Create a new cause for an indicator
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cause body CreateCauseRequest true "Cause data"
// @Success 201 {object} map[string]interface{} "Cause created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to create cause"
// @Router /indicators/causes [post]
func (h *IndicatorHandlers) CreateCause(w http.ResponseWriter, r *http.Request) {
	var req CreateCauseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newCause := models.Cause{
		IndicatorID:       req.IndicatorID,
		Metric:            models.MetricEnum(req.Metric),
		Description:       req.Description,
		ProductivityLevel: models.ProductivityEnum(req.ProductivityLevel),
	}

	ctx := r.Context()
	causeID, err := h.causeUseCase.Create(ctx, newCause)
	if err != nil {
		http.Error(w, "Failed to create cause", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":                 causeID,
		"indicator_id":       req.IndicatorID,
		"metric":             req.Metric,
		"description":        req.Description,
		"productivity_level": req.ProductivityLevel,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// CreateAction handles POST /indicators/actions
// @Summary Create a new action
// @Description Create a new action for an indicator
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param action body CreateActionRequest true "Action data"
// @Success 201 {object} map[string]interface{} "Action created successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 404 {string} string "Cause not found"
// @Failure 500 {string} string "Failed to create action"
// @Router /indicators/actions [post]
func (h *IndicatorHandlers) CreateAction(w http.ResponseWriter, r *http.Request) {
	var req CreateActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get cause to include in action
	ctx := r.Context()
	causes, err := h.causeUseCase.Get(ctx, req.IndicatorID)
	if err != nil {
		http.Error(w, "Failed to get causes", http.StatusInternalServerError)
		return
	}

	var cause models.Cause
	for _, c := range causes {
		if c.ID == req.CauseID {
			cause = c
			break
		}
	}

	if cause.ID == uuid.Nil {
		http.Error(w, "Cause not found", http.StatusNotFound)
		return
	}

	newAction := models.Action{
		IndicatorID: req.IndicatorID,
		Cause:       cause,
		Description: req.Description,
	}

	actionID, err := h.actionUseCase.Create(ctx, newAction)
	if err != nil {
		http.Error(w, "Failed to create action", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":           actionID,
		"indicator_id": req.IndicatorID,
		"cause_id":     req.CauseID,
		"description":  req.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
