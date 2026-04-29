package handlers

import (
	"database/sql"
	"net/http"

	"test2-api/helpers"
	"test2-api/models"
)

// FeedbackHandler holds the database connection
type FeedbackHandler struct {
	DB *sql.DB
}

// NewFeedbackHandler creates a new FeedbackHandler
func NewFeedbackHandler(db *sql.DB) *FeedbackHandler {
	return &FeedbackHandler{DB: db}
}

// CreateFeedback handles POST /feedback
func (h *FeedbackHandler) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	
	helpers.ReadJSON(w, r, &input)

	id, err := models.CreateFeedback(h.DB, input.Name, input.Email, input.Subject, input.Message)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create feedback"})
		return
	}

	createdFeedback := models.Feedback{
		ID:      id,
		Name:    input.Name,
		Email:   input.Email,
		Subject: input.Subject,
		Message: input.Message,
	}

	helpers.WriteJSON(w, http.StatusCreated, createdFeedback)
}

// GetAllFeedback handles GET /feedback
func (h *FeedbackHandler) GetAllFeedback(w http.ResponseWriter, r *http.Request) {
	feedbacks, err := models.GetAllFeedback(h.DB)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve feedback"})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, feedbacks)
}

// GetFeedbackByID handles GET /feedback/{id}
func (h *FeedbackHandler) GetFeedbackByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	feedback, err := models.GetFeedbackByID(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Feedback not found"})
			return
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve feedback"})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, feedback)
}

// DeleteFeedback handles DELETE /feedback/{id}
func (h *FeedbackHandler) DeleteFeedback(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := models.DeleteFeedback(h.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			helpers.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Feedback not found"})
			return
		}
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete feedback"})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Feedback deleted successfully"})
}