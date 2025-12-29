package del

import (
    "log/slog"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "todolist/internal/storage"
)

type mockDeleter struct {
    idToErr      map[int64]error
    idToAffected map[int64]int64
}

func (m *mockDeleter) DelTask(id int64) (int64, error) {
    if err, ok := m.idToErr[id]; ok {
        return 0, err
    }
    if a, ok := m.idToAffected[id]; ok {
        return a, nil
    }
    // по умолчанию — не найдено
    return 0, storage.ErrTaskNotFound
}

func newTestLogger() *slog.Logger {
    return slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
}

func TestDelHandler_InvalidID(t *testing.T) {
    log := newTestLogger()
    m := &mockDeleter{}
    h := NewDel(log, m)

    r := httptest.NewRequest(http.MethodDelete, "/todos/abc", nil)
    r.SetPathValue("id", "abc")

    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)

    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
    }
}

func TestDelHandler_NotFoundError(t *testing.T) {
    log := newTestLogger()
    m := &mockDeleter{
        idToErr: map[int64]error{
            1: storage.ErrTaskNotFound,
        },
    }
    h := NewDel(log, m)

    r := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
    r.SetPathValue("id", "1")

    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
    }
}

func TestDelHandler_NoRowsAffected(t *testing.T) {
    log := newTestLogger()
    m := &mockDeleter{
        idToAffected: map[int64]int64{
            2: 0,
        },
    }
    h := NewDel(log, m)

    r := httptest.NewRequest(http.MethodDelete, "/todos/2", nil)
    r.SetPathValue("id", "2")

    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)

    if w.Code != http.StatusNotFound {
        t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
    }
    body := w.Body.String()
    if !strings.Contains(body, `"id":2`) {
        t.Fatalf("unexpected body: %s", body)
    }
}

func TestDelHandler_Success(t *testing.T) {
    log := newTestLogger()
    m := &mockDeleter{
        idToAffected: map[int64]int64{
            3: 1,
        },
    }
    h := NewDel(log, m)

    r := httptest.NewRequest(http.MethodDelete, "/todos/3", nil)
    r.SetPathValue("id", "3")

    w := httptest.NewRecorder()
    h.ServeHTTP(w, r)

    if w.Code != http.StatusOK {
        t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
    }
    body := w.Body.String()
    if !strings.Contains(body, `"id":3`) {
        t.Fatalf("unexpected body: %s", body)
    }
}
