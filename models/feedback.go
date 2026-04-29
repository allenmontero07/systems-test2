package models

import (
	"database/sql"
	"time"
)

// Feedback represents a feedback record in the database
type Feedback struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateFeedback inserts a new feedback record into the database
func CreateFeedback(db *sql.DB, name, email, subject, message string) (string, error) {
	query := `INSERT INTO feedback (name, email, subject, message) VALUES ($1, $2, $3, $4) RETURNING id`
	var id string
	err := db.QueryRow(query, name, email, subject, message).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetAllFeedback retrieves all feedback records from the database
func GetAllFeedback(db *sql.DB) ([]Feedback, error) {
	query := `SELECT id, name, email, subject, message, created_at FROM feedback ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbacks []Feedback
	for rows.Next() {
		var f Feedback
		err := rows.Scan(&f.ID, &f.Name, &f.Email, &f.Subject, &f.Message, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, f)
	}
	return feedbacks, nil
}

// GetFeedbackByID retrieves a single feedback record by ID
func GetFeedbackByID(db *sql.DB, id string) (*Feedback, error) {
	query := `SELECT id, name, email, subject, message, created_at FROM feedback WHERE id = $1`
	var f Feedback
	err := db.QueryRow(query, id).Scan(&f.ID, &f.Name, &f.Email, &f.Subject, &f.Message, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// DeleteFeedback deletes a feedback record by ID
func DeleteFeedback(db *sql.DB, id string) error {
	query := `DELETE FROM feedback WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}