package getbyid

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"todolist/internal/storage"
	"todolist/internal/storage/sqlite"
)

// поля для json, тестить сторадж без http, и другие плюшки
type TaskResponse struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Discription   string `json:"discription"`
	Achievement   bool   `json:"achievement"`
	TimeCreated   string `json:"timeCreated,omitempty"`
	TimeCompleted string `json:"timeCompleted,omitempty"`
}

type TaskGetterById interface {
	GetTaskByID(id int64) (*sqlite.TaskDB, error)
}

func NewGetterByiD(log *slog.Logger, ID TaskGetterById) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.task.getbyid.new"

		log.Info("Getter by id started", slog.String("op", op))

		
		idStr := r.PathValue("id")
		if idStr == "" {
			log.Info("id is required parametr")
			http.Error(w, "id parameter is required", http.StatusBadRequest)
			return 
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Info("invalid id")
			http.Error(w, "invalid id", http.StatusBadRequest)
			return 
		}

		log.Info("request body decoded", slog.Any("request", id))

		taskDB, err := ID.GetTaskByID(id)
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Info("task not found", slog.Int64("id", (id)))

			http.Error(w, "task not found", http.StatusNotFound)

			return
		}
		taskResp := TaskResponse{
			ID:            taskDB.ID,
			Title:         taskDB.Title,
			Discription:   taskDB.Discription,
			Achievement:   taskDB.Achievement,
			TimeCreated:   taskDB.TimeCreated,
			TimeCompleted: taskDB.TimeCompleted,
		}

		log.Info("task found", slog.Int64("id", id), slog.String("op", op))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(taskResp)

	}
}
