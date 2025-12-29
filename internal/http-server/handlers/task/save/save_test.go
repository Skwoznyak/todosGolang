package save

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "log/slog"
    "todolist/internal/storage"
)

type mockSaver struct {
    taskToErr map[string]error
    nextID    int64
}

func (m *mockSaver) SaveTask(taskToSave string, discription string) (int64, error) {
    if err, ok := m.taskToErr[taskToSave]; ok {
        return 0, err
    }
    m.nextID++
    return m.nextID, nil
}

func newTestLogger() *slog.Logger {
    return slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
}

func TestSaveHandler_InvalidJSON(t *testing.T) {
    log := newTestLogger()
    m := &mockSaver{}
    h := NewSave(log, m)

    r := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{invalid json`))
    w := httptest.NewRecorder()

    h.ServeHTTP(w, r)

    if w.Code != http.StatusInternalServerError {
        t.Fatalf("expected %d, got %d", http.StatusInternalServerError, w.Code)
    }
}

func TestSaveHandler_EmptyTitle(t *testing.T) {
    log := newTestLogger()
    m := &mockSaver{}
    h := NewSave(log, m)

    body := `{"task":"   ","discription":"desc"}`
    r := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
    w := httptest.NewRecorder()

    h.ServeHTTP(w, r)

    if w.Code != http.StatusUnavailableForLegalReasons {
        t.Fatalf("expected %d, got %d", http.StatusUnavailableForLegalReasons, w.Code)
    }
}

func TestSaveHandler_DuplicateTask(t *testing.T) {
    log := newTestLogger()
    m := &mockSaver{
        taskToErr: map[string]error{
            "Task 1": storage.ErrTaskExists,
        },
    }
    h := NewSave(log, m)

    body := `{"task":"Task 1","discription":"desc"}`
    r := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
    w := httptest.NewRecorder()

    h.ServeHTTP(w, r)

    if w.Code != http.StatusConflict {
        t.Fatalf("expected %d, got %d", http.StatusConflict, w.Code)
    }
}

func TestSaveHandler_Success(t *testing.T) {
    log := newTestLogger()
    m := &mockSaver{
        nextID: 0,
    }
    h := NewSave(log, m)

    body := `{"task":"Buy milk","discription":"today"}`
    r := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(body))
    w := httptest.NewRecorder()

    h.ServeHTTP(w, r)

    if w.Code != http.StatusCreated {
        t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
    }

    if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
        t.Fatalf("expected application/json, got %q", ct)
    }

    respBody := w.Body.String()
    if !strings.Contains(respBody, `"id":1`) || !strings.Contains(respBody, "Buy milk") {
        t.Fatalf("unexpected body: %s", respBody)
    }
}
