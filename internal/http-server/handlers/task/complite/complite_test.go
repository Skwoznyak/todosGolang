package complite

import (
    "log/slog"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "todolist/internal/storage"
)

type mockCompliter struct {
    idToErr      map[int64]error
    idToAffected map[int64]int64
}

func (m *mockCompliter) CompliteTask(id int64) (int64, error) {
    if err, ok := m.idToErr[id]; ok {
        return 0, err
    }
    if a, ok := m.idToAffected[id]; ok {
        return a, nil
    }
    
    return 0, storage.ErrTaskNotFound
}

func newTestLogger() *slog.Logger {
    
    return slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
}

func TestCompliteHandler_AlreadyCompleted(t *testing.T) {
    log := newTestLogger()
    m := &mockCompliter{
        idToErr: map[int64]error{
            1: storage.ErrTaskAlreadyComplited,
        },
    }

    h := NewComplite(log, m)

    r := httptest.NewRequest(http.MethodPut, "/todos/1", nil)
    mux := http.NewServeMux()
    mux.Handle("PUT /todos/{id}", h)

    w := httptest.NewRecorder()
    mux.ServeHTTP(w, r)

    if w.Code != http.StatusConflict {
        t.Fatalf("expected status %d, got %d", http.StatusConflict, w.Code)
    }
}

func TestCompliteHandler_NotFound(t *testing.T) {
    log := newTestLogger()
    m := &mockCompliter{
        idToErr: map[int64]error{
            2: storage.ErrTaskNotFound,
        },
    }

    h := NewComplite(log, m)

    r := httptest.NewRequest(http.MethodPut, "/todos/2", nil)
    mux := http.NewServeMux()
    mux.Handle("PUT /todos/{id}", h)

    w := httptest.NewRecorder()
    mux.ServeHTTP(w, r)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
    }
}

func TestCompliteHandler_Success(t *testing.T) {
    log := newTestLogger()
    m := &mockCompliter{
        idToAffected: map[int64]int64{
            3: 1,
        },
    }

    h := NewComplite(log, m)

    r := httptest.NewRequest(http.MethodPut, "/todos/3", nil)
    mux := http.NewServeMux()
    mux.Handle("PUT /todos/{id}", h)

    w := httptest.NewRecorder()
    mux.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
    }

    body := w.Body.String()
    if !strings.Contains(body, `"id":3`) {
        t.Fatalf("unexpected body: %s", body)
    }
}
