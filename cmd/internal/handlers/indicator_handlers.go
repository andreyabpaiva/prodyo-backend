package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"prodyo-backend/cmd/internal/models"
	"prodyo-backend/cmd/internal/usecases"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type IndicatorHandlers struct {
	indicatorUseCase      *usecases.IndicatorUseCase
	indicatorRangeUseCase *usecases.IndicatorRangeUseCase
	causeUseCase          *usecases.CauseUseCase
	actionUseCase         *usecases.ActionUseCase
}

func NewIndicatorHandlers(
	indicatorUseCase *usecases.IndicatorUseCase,
	indicatorRangeUseCase *usecases.IndicatorRangeUseCase,
	causeUseCase *usecases.CauseUseCase,
	actionUseCase *usecases.ActionUseCase,
) *IndicatorHandlers {
	return &IndicatorHandlers{
		indicatorUseCase:      indicatorUseCase,
		indicatorRangeUseCase: indicatorRangeUseCase,
		causeUseCase:          causeUseCase,
		actionUseCase:         actionUseCase,
	}
}

type CreateIndicatorRequest struct {
	IterationID uuid.UUID `json:"iteration_id"`
}

type CreateCauseRequest struct {
	IndicatorRangeID  uuid.UUID `json:"indicator_range_id"`
	Metric            string    `json:"metric"`
	Description       string    `json:"description"`
	ProductivityLevel string    `json:"productivity_level"`
}

type CreateActionRequest struct {
	IndicatorRangeID uuid.UUID  `json:"indicator_range_id"`
	Metric           string     `json:"metric"`
	CauseDescription string     `json:"cause_description"`
	Description      string     `json:"description"`
	Status           *string    `json:"status,omitempty"`
	StartAt          *time.Time `json:"start_at,omitempty"`
	EndAt            *time.Time `json:"end_at,omitempty"`
	AssigneeID       *uuid.UUID `json:"assignee_id,omitempty"`
}

type PatchActionRequest struct {
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	StartAt     *time.Time `json:"start_at,omitempty"`
	EndAt       *time.Time `json:"end_at,omitempty"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
}

// RangeValuesRequest represents min/max values for a productivity level
type RangeValuesRequest struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// ProductivityRangeRequest represents the full range configuration
type ProductivityRangeRequest struct {
	Ok       RangeValuesRequest `json:"ok"`
	Alert    RangeValuesRequest `json:"alert"`
	Critical RangeValuesRequest `json:"critical"`
}

// SetRangeRequest is used to create or update a productivity range for an indicator type at project level
type SetRangeRequest struct {
	ProjectID     uuid.UUID                `json:"project_id"`
	IndicatorType string                   `json:"indicator_type"` // SpeedPerIteration, ReworkPerIteration, InstabilityIndex
	Range         ProductivityRangeRequest `json:"range"`
}

// UpdateMetricValuesRequest is used to update calculated metric values
type UpdateMetricValuesRequest struct {
	SpeedValue       float64 `json:"speed_value"`
	ReworkValue      float64 `json:"rework_value"`
	InstabilityValue float64 `json:"instability_value"`
}

// Get handles GET /indicators
// @Summary Get indicator
// @Description Get indicator with causes, actions, and calculated productivity levels for a specific iteration
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

	// Get project ID to fetch indicator ranges
	projectID, err := h.indicatorUseCase.GetProjectIDByIterationID(ctx, iterationID)
	if err != nil {
		http.Error(w, "Failed to get project ID", http.StatusInternalServerError)
		return
	}

	// Get all indicator ranges for the project
	ranges, err := h.indicatorRangeUseCase.GetByProjectID(ctx, projectID)
	if err != nil {
		http.Error(w, "Failed to get indicator ranges", http.StatusInternalServerError)
		return
	}

	// Load causes and actions from all indicator ranges for this project
	var allCauses []models.Cause
	var allActions []models.Action
	for _, r := range ranges {
		causes, _ := h.causeUseCase.Get(ctx, r.ID)
		actions, _ := h.actionUseCase.Get(ctx, r.ID)
		allCauses = append(allCauses, causes...)
		allActions = append(allActions, actions...)
	}

	indicator.Causes = allCauses
	indicator.Actions = allActions

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

// SetRange handles POST /indicators/ranges
// @Summary Set productivity range for an indicator type
// @Description Create or update the productivity range (OK, Alert, Critical min/max values) for a specific indicator type at project level
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param range body SetRangeRequest true "Range configuration"
// @Success 201 {object} map[string]interface{} "Range set successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Failed to set range"
// @Router /indicators/ranges [post]
func (h *IndicatorHandlers) SetRange(w http.ResponseWriter, r *http.Request) {
	var req SetRangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate indicator type
	indicatorType := models.IndicatorEnum(req.IndicatorType)
	if indicatorType != models.IndicatorSpeedPerIteration &&
		indicatorType != models.IndicatorReworkPerIteration &&
		indicatorType != models.IndicatorInstabilityIndex {
		http.Error(w, "Invalid indicator_type. Must be SpeedPerIteration, ReworkPerIteration, or InstabilityIndex", http.StatusBadRequest)
		return
	}

	ir := models.IndicatorRange{
		ProjectID:     req.ProjectID,
		IndicatorType: indicatorType,
		Range: models.ProductivityRange{
			Ok:       models.RangeValues{Min: req.Range.Ok.Min, Max: req.Range.Ok.Max},
			Alert:    models.RangeValues{Min: req.Range.Alert.Min, Max: req.Range.Alert.Max},
			Critical: models.RangeValues{Min: req.Range.Critical.Min, Max: req.Range.Critical.Max},
		},
	}

	ctx := r.Context()
	rangeID, err := h.indicatorRangeUseCase.SetRange(ctx, ir)
	if err != nil {
		http.Error(w, "Failed to set range", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":             rangeID,
		"project_id":     req.ProjectID,
		"indicator_type": req.IndicatorType,
		"range":          req.Range,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetRanges handles GET /projects/{project_id}/indicator-ranges
// @Summary Get all indicator ranges for a project
// @Description Get all productivity ranges configured for a project
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID" format(uuid)
// @Success 200 {array} models.IndicatorRange "List of ranges"
// @Failure 400 {string} string "Invalid project_id"
// @Failure 500 {string} string "Failed to get ranges"
// @Router /projects/{project_id}/indicator-ranges [get]
func (h *IndicatorHandlers) GetRanges(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIDStr := vars["project_id"]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ranges, err := h.indicatorRangeUseCase.GetByProjectID(ctx, projectID)
	if err != nil {
		http.Error(w, "Failed to get ranges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ranges)
}

// GetRangeByIndicatorType handles GET /projects/{project_id}/indicator-ranges/{indicator_type}
// @Summary Get range for a specific indicator type
// @Description Get the productivity range for a specific indicator type of a project
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID" format(uuid)
// @Param indicator_type path string true "Indicator type (SpeedPerIteration, ReworkPerIteration, InstabilityIndex)"
// @Success 200 {object} models.IndicatorRange "Range configuration"
// @Failure 400 {string} string "Invalid parameters"
// @Failure 404 {string} string "Range not found"
// @Router /projects/{project_id}/indicator-ranges/{indicator_type} [get]
func (h *IndicatorHandlers) GetRangeByIndicatorType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIDStr := vars["project_id"]
	indicatorTypeStr := vars["indicator_type"]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	indicatorType := models.IndicatorEnum(indicatorTypeStr)
	if indicatorType != models.IndicatorSpeedPerIteration &&
		indicatorType != models.IndicatorReworkPerIteration &&
		indicatorType != models.IndicatorInstabilityIndex {
		http.Error(w, "Invalid indicator_type", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ir, err := h.indicatorRangeUseCase.GetByIndicatorType(ctx, projectID, indicatorType)
	if err != nil {
		http.Error(w, "Range not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ir)
}

// DeleteRange handles DELETE /indicators/ranges/{range_id}
// @Summary Delete a productivity range
// @Description Remove a productivity range configuration
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param range_id path string true "Range ID" format(uuid)
// @Success 204 "Range deleted successfully"
// @Failure 400 {string} string "Invalid range_id"
// @Failure 500 {string} string "Failed to delete range"
// @Router /indicators/ranges/{range_id} [delete]
func (h *IndicatorHandlers) DeleteRange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rangeIDStr := vars["range_id"]

	rangeID, err := uuid.Parse(rangeIDStr)
	if err != nil {
		http.Error(w, "Invalid range_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.indicatorRangeUseCase.Delete(ctx, rangeID); err != nil {
		http.Error(w, "Failed to delete range", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateMetricValues handles PUT /indicators/{indicator_id}/metrics
// @Summary Update calculated metric values
// @Description Update the calculated values for speed, rework, and instability metrics
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param indicator_id path string true "Indicator ID" format(uuid)
// @Param metrics body UpdateMetricValuesRequest true "Metric values"
// @Success 200 {object} models.Indicator "Updated indicator"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to update metrics"
// @Router /indicators/{indicator_id}/metrics [put]
func (h *IndicatorHandlers) UpdateMetricValues(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	indicatorIDStr := vars["indicator_id"]

	indicatorID, err := uuid.Parse(indicatorIDStr)
	if err != nil {
		http.Error(w, "Invalid indicator_id", http.StatusBadRequest)
		return
	}

	var req UpdateMetricValuesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.indicatorUseCase.UpdateMetricValues(ctx, indicatorID, req.SpeedValue, req.ReworkValue, req.InstabilityValue); err != nil {
		http.Error(w, "Failed to update metrics", http.StatusInternalServerError)
		return
	}

	// Return the updated indicator
	ind, err := h.indicatorUseCase.GetByID(ctx, indicatorID)
	if err != nil {
		http.Error(w, "Failed to get updated indicator", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ind)
}

// GetMetricSummary handles GET /indicators/{indicator_id}/summary
// @Summary Get metric summary
// @Description Get a summary of all indicators with their values and productivity levels
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param indicator_id path string true "Indicator ID" format(uuid)
// @Success 200 {array} models.IndicatorMetricValue "Metric summary"
// @Failure 400 {string} string "Invalid indicator_id"
// @Failure 404 {string} string "Indicator not found"
// @Router /indicators/{indicator_id}/summary [get]
func (h *IndicatorHandlers) GetMetricSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	indicatorIDStr := vars["indicator_id"]

	indicatorID, err := uuid.Parse(indicatorIDStr)
	if err != nil {
		http.Error(w, "Invalid indicator_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ind, err := h.indicatorUseCase.GetByID(ctx, indicatorID)
	if err != nil {
		http.Error(w, "Indicator not found", http.StatusNotFound)
		return
	}

	summary := ind.GetMetricSummary()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
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
		IndicatorRangeID:  req.IndicatorRangeID,
		Metric:            models.MetricEnum(req.Metric),
		Description:       req.Description,
		ProductivityLevel: models.ProductivityEnum(req.ProductivityLevel),
	}

	ctx := r.Context()
	causeID, err := h.causeUseCase.Create(ctx, newCause)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create cause", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":                 causeID,
		"indicator_range_id": req.IndicatorRangeID,
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
// @Failure 500 {string} string "Failed to create action"
// @Router /indicators/actions [post]
func (h *IndicatorHandlers) CreateAction(w http.ResponseWriter, r *http.Request) {
	var req CreateActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	_, err := h.indicatorRangeUseCase.GetByID(ctx, req.IndicatorRangeID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Indicator range not found. Please create an indicator range first.", http.StatusNotFound)
		return
	}

	metric := models.MetricEnum(req.Metric)
	if metric != models.MetricWorkVelocity &&
		metric != models.MetricReworkIndex &&
		metric != models.MetricInstabilityIndex {
		http.Error(w, "Invalid metric. Must be WorkVelocity, ReworkIndex, or InstabilityIndex", http.StatusBadRequest)
		return
	}

	newCause := models.Cause{
		IndicatorRangeID:  req.IndicatorRangeID,
		Metric:            metric,
		Description:       req.CauseDescription,
		ProductivityLevel: models.ProductivityCritical,
	}

	causeID, err := h.causeUseCase.Create(ctx, newCause)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create cause", http.StatusInternalServerError)
		return
	}

	newCause.ID = causeID

	newAction := models.Action{
		IndicatorRangeID: req.IndicatorRangeID,
		Cause:            newCause,
		Description:      req.Description,
		Status:           models.StatusNotStarted,
	}

	if req.Status != nil {
		newAction.Status = models.StatusEnum(*req.Status)
	}

	if req.StartAt != nil {
		newAction.StartAt = *req.StartAt
	}
	if req.EndAt != nil {
		newAction.EndAt = *req.EndAt
	}

	if req.AssigneeID != nil {
		newAction.Assignee = models.User{
			ID: *req.AssigneeID,
		}
	}

	actionID, err := h.actionUseCase.Create(ctx, newAction)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create action", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":                 actionID,
		"indicator_range_id": req.IndicatorRangeID,
		"cause_id":           causeID,
		"metric":             req.Metric,
		"cause_description":  req.CauseDescription,
		"description":        req.Description,
		"status":             newAction.Status,
		"start_at":           req.StartAt,
		"end_at":             req.EndAt,
		"assignee_id":        req.AssigneeID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// CreateDefaultRanges handles POST /projects/{project_id}/indicator-ranges/default
// @Summary Create default indicator ranges for a project
// @Description Create default productivity ranges for all indicator types for a project
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID" format(uuid)
// @Success 201 {object} map[string]interface{} "Default ranges created"
// @Failure 400 {string} string "Invalid project_id"
// @Failure 500 {string} string "Failed to create default ranges"
// @Router /projects/{project_id}/indicator-ranges/default [post]
func (h *IndicatorHandlers) CreateDefaultRanges(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIDStr := vars["project_id"]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.indicatorRangeUseCase.CreateDefaultRanges(ctx, projectID); err != nil {
		http.Error(w, "Failed to create default ranges", http.StatusInternalServerError)
		return
	}

	// Return the created ranges
	ranges, _ := h.indicatorRangeUseCase.GetByProjectID(ctx, projectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"project_id": projectID,
		"ranges":     ranges,
	})
}

// GetIndicatorRangeIDByProjectIDAndType handles GET /projects/{project_id}/indicator-range-ids/{indicator_type}
// @Summary Get indicator range ID by project ID and indicator type
// @Description Returns the indicator range ID based on project ID and indicator type
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID" format(uuid)
// @Param indicator_type path string true "Indicator type (SpeedPerIteration, ReworkPerIteration, InstabilityIndex)"
// @Success 200 {object} map[string]interface{} "Indicator range ID"
// @Failure 400 {string} string "Invalid parameters"
// @Failure 404 {string} string "Indicator range not found"
// @Router /projects/{project_id}/indicator-range-ids/{indicator_type} [get]
func (h *IndicatorHandlers) GetIndicatorRangeIDByProjectIDAndType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectIDStr := vars["project_id"]
	indicatorTypeStr := vars["indicator_type"]

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project_id", http.StatusBadRequest)
		return
	}

	indicatorType := models.IndicatorEnum(indicatorTypeStr)
	if indicatorType != models.IndicatorSpeedPerIteration &&
		indicatorType != models.IndicatorReworkPerIteration &&
		indicatorType != models.IndicatorInstabilityIndex {
		http.Error(w, "Invalid indicator_type. Must be SpeedPerIteration, ReworkPerIteration, or InstabilityIndex", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ir, err := h.indicatorRangeUseCase.GetByIndicatorType(ctx, projectID, indicatorType)
	if err != nil {
		http.Error(w, "Indicator range not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"indicator_range_id": ir.ID,
	})
}

// GetCausesAndActionsByIteration godoc
// @Summary Get all causes and actions for an iteration
// @Description Retrieves all causes and actions associated with a specific iteration
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param iteration_id path string true "Iteration ID"
// @Success 200 {object} map[string]interface{} "Causes and actions retrieved successfully"
// @Failure 400 {string} string "Invalid iteration ID"
// @Failure 500 {string} string "Failed to retrieve causes and actions"
// @Router /iterations/{iteration_id}/causes-actions [get]
func (h *IndicatorHandlers) GetCausesAndActionsByIteration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	iterationIDStr := vars["iteration_id"]

	iterationID, err := uuid.Parse(iterationIDStr)
	if err != nil {
		http.Error(w, "Invalid iteration ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	causes, err := h.causeUseCase.GetByIterationID(ctx, iterationID)
	if err != nil {
		log.Printf("Failed to get causes: %v", err)
		http.Error(w, "Failed to retrieve causes", http.StatusInternalServerError)
		return
	}

	actions, err := h.actionUseCase.GetByIterationID(ctx, iterationID)
	if err != nil {
		log.Printf("Failed to get actions: %v", err)
		http.Error(w, "Failed to retrieve actions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"iteration_id": iterationID,
		"causes":       causes,
		"actions":      actions,
	})
}

// PatchAction handles PATCH /indicators/actions/{id}
// @Summary Partially update action
// @Description Partially update an existing action (only provided fields will be updated)
// @Tags indicators
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Action ID" format(uuid)
// @Param action body PatchActionRequest true "Partial action data"
// @Success 200 {object} models.Action "Updated action"
// @Failure 400 {string} string "Invalid action ID or request body"
// @Failure 404 {string} string "Action not found"
// @Failure 500 {string} string "Failed to update action"
// @Router /indicators/actions/{id} [patch]
func (h *IndicatorHandlers) PatchAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid action ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	existingAction, err := h.actionUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Action not found", http.StatusNotFound)
		return
	}

	var req PatchActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Description != nil {
		existingAction.Description = *req.Description
	}

	if req.Status != nil {
		existingAction.Status = models.StatusEnum(*req.Status)
	}

	if req.StartAt != nil {
		existingAction.StartAt = *req.StartAt
	}

	if req.EndAt != nil {
		existingAction.EndAt = *req.EndAt
	}

	if req.AssigneeID != nil {
		existingAction.Assignee.ID = *req.AssigneeID
	}

	err = h.actionUseCase.Update(ctx, existingAction)
	if err != nil {
		http.Error(w, "Failed to update action", http.StatusInternalServerError)
		return
	}

	updatedAction, err := h.actionUseCase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "Failed to retrieve updated action", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAction)
}
