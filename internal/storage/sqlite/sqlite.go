package sqlite

import (
	"database/sql"
	"fmt"
	_ "os/exec"
	"todolist/internal/storage"

	"github.com/mattn/go-sqlite3" //init sql3 draiver
)

// поля только бд, потом конвертируем
type TaskDB struct {
	ID            int64
	Title         string
	Discription   string
	Achievement   bool
	TimeCreated   string
	TimeCompleted string
}

type Storage struct {
	db *sql.DB
}

// типо миграции
func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS tasks(
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL UNIQUE,
	discription TEXT,
	achievement BOOLEAN NOT NULL DEFAULT false,
	timeCreated TEXT,
	timeCompleted TEXT);
	CREATE INDEX IF NOT EXISTS idx_tasks ON tasks(title);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// сохранение
func (s *Storage) SaveTask(taskToSave string, discription string) (int64, error) {
	const op = "storage.sqlite.SaveTask"

	if taskToSave == "" {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrEmptyTitle)
	}

	stmt, err := s.db.Prepare(`INSERT INTO tasks(title, discription, timeCreated)
	VALUES(?, ?, CURRENT_TIMESTAMP)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(taskToSave, discription)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrTaskExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err) // Какие то бд не имею такого метода, тут может быть ошибка
	}

	return id, nil
}

// удаление
func (s *Storage) DelTask(id int64) (int64, error) {
	const op = "storage.sqlite.DelTask"

	stmt, err := s.db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	affected, err := res.RowsAffected() //тоже самое что и ластинсерт
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return affected, nil
}

// выполнение
func (s *Storage) CompliteTask(id int64) (int64, error) {
	const op = "storage.sqlite.CompliteTask"

	var currentStatus bool
	err := s.db.QueryRow("SELECT achievement FROM tasks WHERE id = ?", id).Scan(&currentStatus)
	// как вариант
	// if err == sql.ErrNoRows {
	//     return 0, fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
	// }
	if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrNoExtended(sqlite3.ErrNotFound) {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("%s: select failed: %w", op, err)
	}

	if currentStatus {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrTaskAlreadyComplited)
	}

	stmt, err := s.db.Prepare(`UPDATE tasks SET achievement = true,
	timeCompleted = CURRENT_TIMESTAMP
	WHERE id = ?`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare failed %w", op, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return affected, nil
}

// гетеры
func (s *Storage) GetTaskByID(id int64) (*TaskDB, error) {
	const op = "storage.sqlite.GetTaskByID"

	task := &TaskDB{}
	//падает в панику при скане, если нулл
	err := s.db.QueryRow(`
    SELECT id, 
           COALESCE(title, ''),           -- NULL → ""
           COALESCE(discription, ''),     -- NULL → ""
           achievement, 
           COALESCE(timeCreated, ''),     -- NULL → ""
           COALESCE(timeCompleted, '')    -- NULL → ""
    FROM tasks WHERE id = ?
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Discription,
		&task.Achievement,
		&task.TimeCreated,
		&task.TimeCompleted,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrTaskNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return task, nil
}

func (s *Storage) GetAllTasks() ([]TaskDB, error) {
    const op = "storage.sqlite.GetAllTasks"

    rows, err := s.db.Query(`
        SELECT id, 
               COALESCE(title, ''), 
               COALESCE(discription, ''), 
               achievement, 
               COALESCE(timeCreated, ''), 
               COALESCE(timeCompleted, '') 
        FROM tasks
    `)
    if err != nil {
        return nil, fmt.Errorf("%s: query failed: %w", op, err)
    }

    tasks := make([]TaskDB, 0)

    for rows.Next() {
        task := TaskDB{}
        err := rows.Scan(
            &task.ID,
            &task.Title,
            &task.Discription,
            &task.Achievement,
            &task.TimeCreated,
            &task.TimeCompleted,
        )
        if err != nil {
            return nil, fmt.Errorf("%s: scan failed: %w", op, err)
        }
        tasks = append(tasks, task)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: rows error: %w", op, err)
    }

    return tasks, nil
}
