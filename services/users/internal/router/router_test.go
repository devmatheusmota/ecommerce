package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "test-secret-for-handlers"

func TestRouter_Health(t *testing.T) {
	repo := repository.NewMockUserRepository()
	h := NewWithRepository(repo)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got, want := rec.Code, http.StatusOK; got != want {
		t.Errorf("status: got %d, want %d", got, want)
	}
	var body struct {
		Data struct {
			Status  string `json:"status"`
			Service string `json:"service"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got, want := body.Data.Status, "ok"; got != want {
		t.Errorf("data.status: got %q, want %q", got, want)
	}
	if got, want := body.Data.Service, "users"; got != want {
		t.Errorf("data.service: got %q, want %q", got, want)
	}
}

func TestRouter_Register(t *testing.T) {
	repo := repository.NewMockUserRepository()
	h := NewWithRepository(repo)
	validCPF := "529.982.247-25"

	t.Run("success", func(t *testing.T) {
		body := map[string]string{
			"email": "register@example.com", "password": "pass1234",
			"name": "Register User", "phone": "11999999999", "cpf": validCPF,
		}
		req := httptest.NewRequest(http.MethodPost, "/register", jsonBody(body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if got, want := rec.Code, http.StatusCreated; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data struct {
				ID    string `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got := res.Data.ID; got == "" {
			t.Errorf("data.id: got %q, want non-empty", got)
		}
		if got, want := res.Data.Email, "register@example.com"; got != want {
			t.Errorf("data.email: got %q, want %q", got, want)
		}
		if got, want := res.Data.Name, "Register User"; got != want {
			t.Errorf("data.name: got %q, want %q", got, want)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("not json")))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusBadRequest; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		body := map[string]string{"email": "bad", "password": "short", "name": "A", "phone": "1", "cpf": validCPF}
		req := httptest.NewRequest(http.MethodPost, "/register", jsonBody(body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusBadRequest; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		body := map[string]string{
			"email": "dup@example.com", "password": "pass1234",
			"name": "First", "phone": "11999999999", "cpf": validCPF,
		}
		req1 := httptest.NewRequest(http.MethodPost, "/register", jsonBody(body))
		rec1 := httptest.NewRecorder()
		h.ServeHTTP(rec1, req1)
		if got, want := rec1.Code, http.StatusCreated; got != want {
			t.Fatalf("first register status: got %d, want %d", got, want)
		}
		req2 := httptest.NewRequest(http.MethodPost, "/register", jsonBody(body))
		rec2 := httptest.NewRecorder()
		h.ServeHTTP(rec2, req2)
		if got, want := rec2.Code, http.StatusConflict; got != want {
			t.Errorf("second register status: got %d, want %d", got, want)
		}
	})
}

func TestRouter_Login(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", testJWTSecret)
	t.Cleanup(func() { os.Setenv("JWT_SECRET", originalSecret) })

	repo := repository.NewMockUserRepository()
	hash, _ := bcrypt.GenerateFromPassword([]byte("mypass"), bcrypt.DefaultCost)
	repo.SetUser(&domain.User{
		ID: "uid-1", Email: "login@example.com", Name: "Login", Phone: "11999999999",
		CPF: "529.982.247-25", PasswordHash: string(hash),
	})
	h := NewWithRepository(repo)

	t.Run("success", func(t *testing.T) {
		body := map[string]string{"email": "login@example.com", "password": "mypass"}
		req := httptest.NewRequest(http.MethodPost, "/login", jsonBody(body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data struct {
				Token    string `json:"token"`
				ExpireAt string `json:"expire_at"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got := res.Data.Token; got == "" {
			t.Errorf("data.token: got %q, want non-empty", got)
		}
		if got := res.Data.ExpireAt; got == "" {
			t.Errorf("data.expire_at: got %q, want non-empty", got)
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		body := map[string]string{"email": "login@example.com", "password": "wrong"}
		req := httptest.NewRequest(http.MethodPost, "/login", jsonBody(body))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusUnauthorized; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})
}

func TestRouter_Me(t *testing.T) {
	repo := repository.NewMockUserRepository()
	repo.SetUser(&domain.User{
		ID: "me-user-id", Email: "me@example.com", Name: "Me Name", Phone: "11888887777",
		CPF: "529.982.247-25", PasswordHash: "hash",
	})
	h := NewWithRepository(repo)

	t.Run("missing X-User-ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusUnauthorized; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("X-User-ID", "me-user-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data struct {
				ID    string `json:"id"`
				Email string `json:"email"`
				Name  string `json:"name"`
				Phone string `json:"phone"`
				CPF   string `json:"cpf"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got, want := res.Data.ID, "me-user-id"; got != want {
			t.Errorf("data.id: got %q, want %q", got, want)
		}
		if got, want := res.Data.Email, "me@example.com"; got != want {
			t.Errorf("data.email: got %q, want %q", got, want)
		}
		if got, want := res.Data.Name, "Me Name"; got != want {
			t.Errorf("data.name: got %q, want %q", got, want)
		}
		if got, want := res.Data.Phone, "11888887777"; got != want {
			t.Errorf("data.phone: got %q, want %q", got, want)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("X-User-ID", "nonexistent-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusNotFound; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})
}

func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
