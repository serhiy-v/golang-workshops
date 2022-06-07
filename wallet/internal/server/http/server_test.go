package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/workshops/wallet/internal/middleware/auth"
	"github.com/workshops/wallet/internal/repository/postgre"
	"github.com/workshops/wallet/internal/services/validator"
	"github.com/workshops/wallet/internal/services/wallet"
)

//nolint
func TestGetUsers(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Unable to connect")
	}
	defer db.Close()
	repo := postgre.NewRepository(db)
	validate := validator.NewValidator()
	service := wallet.NewService(repo)
	wrapper := auth.NewJwtWrapper("verysecretkey", 999)
	srv := NewServer(service, wrapper, validate)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	srv.GetUsers(w, req)
	assert.Equal(t, w.Code, 200)
}
