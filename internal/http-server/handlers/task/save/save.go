package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	resp "todolist/internal/api/responce" //при использовании либ, выводил бы так, тут считаю неудобным
	"todolist/internal/storage"
)



type request struct {
	Task        string `json:"task"`
	Discription string `json:"discription"`
}

// если бы мы хотели еще что-то возвращать, например гет методы
type Responce struct {
	resp.Responce
}

// интерфейс стораджа, по месту использования)
type TaskSaver interface {
	SaveTask(taskToSave string, discription string) (int64, error)
}

func NewSave(log *slog.Logger, taskSave TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.task.save.new"

		log.Info("save started", slog.String("op", op))

		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("decode failed",
				slog.String("error", err.Error()),
				slog.String("op", op),
			)	

			http.Error(w, "decode failed", http.StatusInternalServerError)

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validateSaveRequest(&req); err != nil {
			 http.Error(w, "need task name", http.StatusUnavailableForLegalReasons)
			 return 
		}

		log.Info("valid request",
		slog.String("op", op))

		id, err := taskSave.SaveTask(req.Task, req.Discription)
		if errors.Is(err, storage.ErrTaskExists) {
			log.Info("task already exists", slog.String("task", req.Task))

			http.Error(w, "task already exists", http.StatusConflict)

			return
		}

		log.Info("task saved", slog.Int64("id", id))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"message":"задача добавлена: %s","id":%d}`, req.Task, id)
	}
}
// при использоыании сторонних библиотек, было бы намного удобнее валидировать даже на уровне струтуры
//Task        string `json:"task" validator:required` - например так
// а так же для больших проектов валидировать с пакетом validator
func validateSaveRequest(req *request) error {
	if strings.TrimSpace(req.Task) == "" {
		return storage.ErrEmptyTitle
	}

	return nil
}
