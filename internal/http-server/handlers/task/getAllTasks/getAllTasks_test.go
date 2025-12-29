package getalltasks

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "errors"
    "todolist/internal/storage/sqlite"
)

type mockGetterAll struct {
    tasks []sqlite.TaskDB
    err   error
}

func (m *mockGetterAll) GetAllTasks() ([]sqlite.TaskDB, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.tasks, nil
}

func newTestLogger() *slog.Logger {
    return slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
}

func TestGetAllTasks_Handler_Error(t *testing.T) {
    log := newTestLogger()
    m := &mockGetterAll{
        err: errors.New("storage error"),
    }

    h := NewGetAllTasks(log, m)

    r := httptest.NewRequest(http.MethodGet, "/todos", nil)
    w := httptest.NewRecorder()

    h.ServeHTTP(w, r)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected %d, got %d", http.StatusInternalServerError, w.Code)
    }
}

func TestGetAllTasks_Handler_Success(t *testing.T) {
    log := newTestLogger()
    m := &mockGetterAll{
        tasks: []sqlite.TaskDB{
            {
                ID:          1,
                Title:       "Task 1",
                Discription: "Desc 1",
                Achievement: false,
                TimeCreated: "2025-12-30T00:00:00Z",
            },
            {
                ID:          2,
                Title:       "Task 2",
                Discription: "Desc 2",
                Achievement: true,
                TimeCreated: "2025-12-30T01:00:00Z",
                TimeCompleted: "2025-12-30T02:00:00Z",
            },
        },
    }

    h := NewGetAllTasks(log, m)

    r := httptest.NewRequest(http.MethodGet, "/todos", nil)
    w := httptest.NewRecorder()

    h.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
    }

    if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
        t.Fatalf("expected application/json, got %q", ct)
    }

    var resp []TaskResponse
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal response: %v", err)
    }

    if len(resp) != 2 {
        t.Fatalf("expected 2 tasks, got %d", len(resp))
    }

    if resp[0].ID != 1 || resp[0].Title != "Task 1" {
        t.Fatalf("unexpected first task: %#v", resp[0])
    }
    if resp[1].ID != 2 || !resp[1].Achievement {
        t.Fatalf("unexpected second task: %#v", resp[1])
    }
}

func TestGetAllTasks_Handler_Empty(t *testing.T) {
    log := newTestLogger()
    m := &mockGetterAll{
        tasks: []sqlite.TaskDB{},
    }

    h := NewGetAllTasks(log, m)

    r := httptest.NewRequest(http.MethodGet, "/todos", nil)
    w := httptest.NewRecorder()

    h.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
    }

    var resp []TaskResponse
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal response: %v", err)
    }

    if len(resp) != 0 {
        t.Fatalf("expected 0 tasks, got %d", len(resp))
    }
}
