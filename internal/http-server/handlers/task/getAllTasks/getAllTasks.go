package getalltasks

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "todolist/internal/storage/sqlite"
)

type TaskResponse struct {
    ID          int64  `json:"id"`
    Title       string `json:"title"`
    Discription string `json:"discription"`
    Achievement bool   `json:"achievement"`
    TimeCreated string `json:"timeCreated,omitempty"`
    TimeCompleted string `json:"timeCompleted,omitempty"`
}

type TaskGetterAll interface {
    GetAllTasks() ([]sqlite.TaskDB, error)
}

func NewGetAllTasks(log *slog.Logger, getter TaskGetterAll) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        const op = "handler.task.getalltasks.new"

        log.Info("Get all tasks started", slog.String("op", op))

        tasksDB, err := getter.GetAllTasks()
        if err != nil {
            log.Error("get all tasks failed", 
                slog.String("error", err.Error()),
                slog.String("op", op))
            
            http.Error(w, "get all tasks failed", http.StatusInternalServerError)
            return
        }

        
        var tasksResp []TaskResponse
        for _, task := range tasksDB {
            tasksResp = append(tasksResp, TaskResponse{
                ID:          task.ID,
                Title:       task.Title,
                Discription: task.Discription,
                Achievement: task.Achievement,
                TimeCreated: task.TimeCreated,
                TimeCompleted: task.TimeCompleted,
            })
        }

        log.Info("all tasks found", 
            slog.Int("count", len(tasksResp)),
            slog.String("op", op))

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(tasksResp)
    }
}
