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

func TestRouter_UpdateMe(t *testing.T) {
	validCPF := "529.982.247-25"
	repo := repository.NewMockUserRepository()
	repo.SetUser(&domain.User{
		ID: "update-me-id", Email: "update@example.com", Name: "Before", Phone: "11000000000",
		CPF: validCPF, PasswordHash: "hash",
	})
	h := NewWithRepository(repo)

	t.Run("missing X-User-ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/me", jsonBody(map[string]string{"name": "After"}))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusUnauthorized; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})

	t.Run("validation error no fields", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/me", jsonBody(map[string]string{}))
		req.Header.Set("X-User-ID", "update-me-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusBadRequest; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})

	t.Run("success update name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/me", jsonBody(map[string]string{"name": "After Name"}))
		req.Header.Set("X-User-ID", "update-me-id")
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
		if got, want := res.Data.Name, "After Name"; got != want {
			t.Errorf("data.name: got %q, want %q", got, want)
		}
		if got, want := res.Data.Email, "update@example.com"; got != want {
			t.Errorf("data.email: got %q, want %q", got, want)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/me", jsonBody(map[string]string{"name": "Any"}))
		req.Header.Set("X-User-ID", "nonexistent-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusNotFound; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})
}

func TestRouter_Addresses(t *testing.T) {
	validCPF := "529.982.247-25"
	userRepo := repository.NewMockUserRepository()
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	userRepo.SetUser(&domain.User{
		ID: "addr-user-id", Email: "addr@example.com", Name: "User", Phone: "11999999999",
		CPF: validCPF, PasswordHash: string(hash),
	})
	addrRepo := repository.NewMockAddressRepository()
	h := NewWithRepositories(userRepo, addrRepo, repository.NewMockPasswordResetRepository())

	t.Run("create address success", func(t *testing.T) {
		body := map[string]any{
			"street": "Rua A", "number": "1", "complement": "", "neighborhood": "Centro", "city": "São Paulo", "state": "SP", "zip_code": "01310100", "type": "shipping", "is_default_shipping": true,
		}
		req := httptest.NewRequest(http.MethodPost, "/me/addresses", jsonBody(body))
		req.Header.Set("X-User-ID", "addr-user-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusCreated; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data struct {
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if res.Data.ID == "" || res.Data.Type != "shipping" {
			t.Errorf("got %+v", res.Data)
		}
	})

	t.Run("list addresses", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me/addresses", nil)
		req.Header.Set("X-User-ID", "addr-user-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data []struct {
				ID                string `json:"id"`
				Street            string `json:"street"`
				ZipCode           string `json:"zip_code"`
				Type              string `json:"type"`
				IsDefaultShipping bool   `json:"is_default_shipping"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if len(res.Data) < 1 {
			t.Errorf("expected at least 1 address, got %d", len(res.Data))
		}
		if res.Data[0].Street != "Rua A" || res.Data[0].ZipCode != "01310100" || !res.Data[0].IsDefaultShipping {
			t.Errorf("got %+v", res.Data[0])
		}
	})

	t.Run("get address by id", func(t *testing.T) {
		createReq := httptest.NewRequest(http.MethodPost, "/me/addresses", jsonBody(map[string]any{
			"street": "Rua B", "number": "2", "complement": "Apto 1", "neighborhood": "Jardim", "city": "Rio", "state": "RJ", "zip_code": "22041080", "type": "billing", "is_default_billing": false, "is_default_shipping": false,
		}))
		createReq.Header.Set("X-User-ID", "addr-user-id")
		createRec := httptest.NewRecorder()
		h.ServeHTTP(createRec, createReq)
		var createRes struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if json.NewDecoder(createRec.Body).Decode(&createRes) != nil || createRes.Data.ID == "" {
			t.Fatalf("need created address id")
		}
		req := httptest.NewRequest(http.MethodGet, "/me/addresses/"+createRes.Data.ID, nil)
		req.Header.Set("X-User-ID", "addr-user-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data struct {
				ID         string `json:"id"`
				Street     string `json:"street"`
				Complement string `json:"complement"`
				Type       string `json:"type"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if res.Data.Street != "Rua B" || res.Data.Complement != "Apto 1" || res.Data.Type != "billing" {
			t.Errorf("got %+v", res.Data)
		}
	})

	t.Run("update address", func(t *testing.T) {
		createReq := httptest.NewRequest(http.MethodPost, "/me/addresses", jsonBody(map[string]any{
			"street": "Rua C", "number": "3", "complement": "", "neighborhood": "Centro", "city": "SP", "state": "SP", "zip_code": "01310100", "type": "shipping", "is_default_shipping": false, "is_default_billing": false,
		}))
		createReq.Header.Set("X-User-ID", "addr-user-id")
		createRec := httptest.NewRecorder()
		h.ServeHTTP(createRec, createReq)
		var createRes struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if json.NewDecoder(createRec.Body).Decode(&createRes) != nil || createRes.Data.ID == "" {
			t.Fatalf("need created address id")
		}
		updateBody := map[string]any{
			"street": "Rua C Updated", "number": "3", "complement": "Sala 10", "neighborhood": "Centro", "city": "SP", "state": "SP", "zip_code": "01310100", "type": "shipping", "is_default_shipping": true, "is_default_billing": false,
		}
		req := httptest.NewRequest(http.MethodPatch, "/me/addresses/"+createRes.Data.ID, jsonBody(updateBody))
		req.Header.Set("X-User-ID", "addr-user-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data struct {
				Street            string `json:"street"`
				Complement        string `json:"complement"`
				IsDefaultShipping bool   `json:"is_default_shipping"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if res.Data.Street != "Rua C Updated" || res.Data.Complement != "Sala 10" || !res.Data.IsDefaultShipping {
			t.Errorf("got %+v", res.Data)
		}
	})

	t.Run("delete address", func(t *testing.T) {
		createReq := httptest.NewRequest(http.MethodPost, "/me/addresses", jsonBody(map[string]any{
			"street": "Rua D", "number": "4", "complement": "", "neighborhood": "Norte", "city": "BH", "state": "MG", "zip_code": "30130000", "type": "billing", "is_default_billing": false, "is_default_shipping": false,
		}))
		createReq.Header.Set("X-User-ID", "addr-user-id")
		createRec := httptest.NewRecorder()
		h.ServeHTTP(createRec, createReq)
		var createRes struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if json.NewDecoder(createRec.Body).Decode(&createRes) != nil || createRes.Data.ID == "" {
			t.Fatalf("need created address id")
		}
		req := httptest.NewRequest(http.MethodDelete, "/me/addresses/"+createRes.Data.ID, nil)
		req.Header.Set("X-User-ID", "addr-user-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusNoContent; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		getReq := httptest.NewRequest(http.MethodGet, "/me/addresses/"+createRes.Data.ID, nil)
		getReq.Header.Set("X-User-ID", "addr-user-id")
		getRec := httptest.NewRecorder()
		h.ServeHTTP(getRec, getReq)
		if got, want := getRec.Code, http.StatusNotFound; got != want {
			t.Errorf("after delete get status: got %d, want %d", got, want)
		}
	})

	t.Run("get address wrong user returns 404", func(t *testing.T) {
		listReq := httptest.NewRequest(http.MethodGet, "/me/addresses", nil)
		listReq.Header.Set("X-User-ID", "addr-user-id")
		listRec := httptest.NewRecorder()
		h.ServeHTTP(listRec, listReq)
		var listRes struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if json.NewDecoder(listRec.Body).Decode(&listRes) != nil || len(listRes.Data) == 0 {
			t.Fatalf("need at least one address")
		}
		addrID := listRes.Data[0].ID
		req := httptest.NewRequest(http.MethodGet, "/me/addresses/"+addrID, nil)
		req.Header.Set("X-User-ID", "other-user-id")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusNotFound; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})
}

func TestRouter_PasswordReset(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", testJWTSecret)
	t.Cleanup(func() { os.Setenv("JWT_SECRET", originalSecret) })

	validCPF := "529.982.247-25"
	userRepo := repository.NewMockUserRepository()
	hash, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
	userRepo.SetUser(&domain.User{
		ID: "reset-user-id", Email: "reset@example.com", Name: "Reset", Phone: "11999999999",
		CPF: validCPF, PasswordHash: string(hash),
	})
	passResetRepo := repository.NewMockPasswordResetRepository()
	h := NewWithRepositories(userRepo, repository.NewMockAddressRepository(), passResetRepo)

	t.Run("request success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/password-reset/request", jsonBody(map[string]string{"email": "reset@example.com"}))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		var res struct {
			Data struct {
				Message    string `json:"message"`
				ResetToken string `json:"reset_token"`
			} `json:"data"`
		}
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if res.Data.ResetToken == "" {
			t.Error("expected non-empty reset_token")
		}
	})

	t.Run("request email not found still 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/password-reset/request", jsonBody(map[string]string{"email": "nonexistent@example.com"}))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})

	t.Run("confirm invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/password-reset/confirm", jsonBody(map[string]string{"token": "invalid", "new_password": "newpass123"}))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got, want := rec.Code, http.StatusBadRequest; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
	})

	t.Run("confirm success", func(t *testing.T) {
		requestReq := httptest.NewRequest(http.MethodPost, "/password-reset/request", jsonBody(map[string]string{"email": "reset@example.com"}))
		requestRec := httptest.NewRecorder()
		h.ServeHTTP(requestRec, requestReq)
		var requestRes struct {
			Data struct {
				ResetToken string `json:"reset_token"`
			} `json:"data"`
		}
		if err := json.NewDecoder(requestRec.Body).Decode(&requestRes); err != nil || requestRes.Data.ResetToken == "" {
			t.Skip("need token from request")
		}
		confirmReq := httptest.NewRequest(http.MethodPost, "/password-reset/confirm", jsonBody(map[string]string{"token": requestRes.Data.ResetToken, "new_password": "newpass123"}))
		confirmRec := httptest.NewRecorder()
		h.ServeHTTP(confirmRec, confirmReq)
		if got, want := confirmRec.Code, http.StatusNoContent; got != want {
			t.Errorf("status: got %d, want %d", got, want)
		}
		loginReq := httptest.NewRequest(http.MethodPost, "/login", jsonBody(map[string]string{"email": "reset@example.com", "password": "newpass123"}))
		loginRec := httptest.NewRecorder()
		h.ServeHTTP(loginRec, loginReq)
		if got, want := loginRec.Code, http.StatusOK; got != want {
			t.Errorf("login after reset: got %d, want %d", got, want)
		}
	})
}

func jsonBody(v any) *bytes.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
