package del

import (
	"log/slog"
	"net/http"
	"todolist/internal/storage"
	"errors"
	"fmt"
	"strconv"
)

type TaskDeleter interface {
	DelTask(id int64) (int64, error)
}

func NewDel(log *slog.Logger, ID TaskDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.task.del.new"

		log.Info("Delete started", slog.String("op", op))
		

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


		affected, err := ID.DelTask(id)
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Info("task not found", slog.Int64("id", (id)))

			http.Error(w, "task not found", http.StatusNotFound)

			return
		}

		log.Info("task deleted",
		slog.Int64("id", id),
		slog.Int64("row affected", affected))

		if affected == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"message":"задачи с таким айди нет ","id":%d}`, id)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"задача удалена: ","id":%d}`, id)
	}
}
