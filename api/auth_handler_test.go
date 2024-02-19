package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nytro04/hotel-reservation/db"
	"github.com/nytro04/hotel-reservation/types"
)

func insertUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "kwaku@gmail.com",
		FirstName: "kwaku",
		LastName:  "Ansah",
		Password:  "P@ssw0rd",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestAuthSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "kwaku@gmail.com",
		Password: "P@ssw0rd",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatalf("expected the JWT Token to be present in the auth response")
	}
	// set the encrypted password to an empty string, because we do not return that in any JSON response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be the inserted user")
	}
}
func TestAuthWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := AuthParams{
		Email:    "kwaku@gmail.com",
		Password: "P@ssw0rdwrong",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response type to be <invalid credentials> %s", genResp.Msg)
	}

}
