package main

import (
	"log/slog"
	"net/http"
	"os"
	"todolist/internal/config"
	"todolist/internal/http-server/handlers/task/complite"
	"todolist/internal/http-server/handlers/task/del"
	getalltasks "todolist/internal/http-server/handlers/task/getAllTasks"
	getbyid "todolist/internal/http-server/handlers/task/getById"
	"todolist/internal/http-server/handlers/task/save"
	"todolist/internal/http-server/middleware/logger"
	_ "todolist/internal/storage"
	"todolist/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//если смущает можно заупскать, но мне больше нравится так
	// >> $Env:CONFIG_PATH = ".\config\local.yaml"
	// >> go run cmd/todolist/main.go

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting todolist", slog.String("env", cfg.Env))
	log.Debug("debug message are enable")

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		msg := "Failed to migrate: " + (err.Error())
		log.Error(msg)
		os.Exit(1)
	}

	// id, err := storage.SaveTask("купить молоко", "vkusno")
	// if err != nil {
	// 	log.Info("failed to save task")
	// }
	// msg := "was created task"
	// log.Info(msg)

	// _ = id
	// qa, err := storage.DelTask(1)
	// if err != nil {
	// 	log.Info("проблемы с удалением")
	// 	os.Exit(2)
	// }

	// log.Info("Удаление по 1 айди, прошло успешно: ", slog.String("qa: ",strconv.Itoa(int(qa)) ))

	//можно было использовать mux, но мне так просто больше нравится
	http.HandleFunc("POST /todos", logger.MyMiddleware(log, save.NewSave(log, storage))) //
	http.HandleFunc("DELETE /todos/{id}", logger.MyMiddleware(log, del.NewDel(log, storage))) //
	http.HandleFunc("PUT /complite/{id}", logger.MyMiddleware(log, complite.NewComplite(log, storage)))
	http.HandleFunc("GET /todos/{id}", logger.MyMiddleware(log, getbyid.NewGetterByiD(log, storage))) //
	http.HandleFunc("GET /todos", logger.MyMiddleware(log, getalltasks.NewGetAllTasks(log, storage))) //

	if err = http.ListenAndServe(":9091", nil); err != nil {
		log.Info("failed to run http server")
		os.Exit(3)
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
