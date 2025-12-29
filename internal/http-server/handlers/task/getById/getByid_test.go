package getbyid

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "todolist/internal/storage"
    "todolist/internal/storage/sqlite"
)

type mockGetterByID struct {
    idToTask map[int64]*sqlite.TaskDB
    idToErr  map[int64]error
}

func (m *mockGetterByID) GetTaskByID(id int64) (*sqlite.TaskDB, error) {
    if err, ok := m.idToErr[id]; ok {
        return nil, err
    }
    if t, ok := m.idToTask[id]; ok {
        return t, nil
    }
    return nil, storage.ErrTaskNotFound
}

func newTestLogger() *slog.Logger {
    return slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
}

func TestGetByID_InvalidID(t *testing.T) {
    log := newTestLogger()
    m := &mockGetterByID{}
    h := NewGetterByiD(log, m)

    r := httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
    r.SetPathValue("id", "abc")

    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
    }
}

func TestGetByID_NotFound(t *testing.T) {
    log := newTestLogger()
    m := &mockGetterByID{
        idToErr: map[int64]error{
            2: storage.ErrTaskNotFound,
        },
    }
    h := NewGetterByiD(log, m)

    r := httptest.NewRequest(http.MethodGet, "/todos/2", nil)
    r.SetPathValue("id", "2")

    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
    }
}

func TestGetByID_Success(t *testing.T) {
    log := newTestLogger()
    m := &mockGetterByID{
        idToTask: map[int64]*sqlite.TaskDB{
            3: {
                ID:          3,
                Title:       "Task 3",
                Discription: "Desc 3",
                Achievement: true,
                TimeCreated: "2025-12-30T00:00:00Z",
                TimeCompleted: "2025-12-30T01:00:00Z",
            },
        },
    }

    h := NewGetterByiD(log, m)

    r := httptest.NewRequest(http.MethodGet, "/todos/3", nil)
    r.SetPathValue("id", "3")

    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
    }

    if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
        t.Fatalf("expected application/json, got %q", ct)
    }

    var resp TaskResponse
    if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }

    if resp.ID != 3 || resp.Title != "Task 3" || !resp.Achievement {
        t.Fatalf("unexpected response: %#v", resp)
    }
}
