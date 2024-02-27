package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nytro04/hotel-reservation/types"
)

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()

	userHandler := NewUserHandler(tdb.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "kwaku@gmail.com",
		FirstName: "kwaku",
		LastName:  "Ansah",
		Password:  "P@ssw0rd",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "Application/json")

	res, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User

	json.NewDecoder(res.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("Expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Expected the Encrypted Password not to be included in the json response")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected LastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected Email %s but got %s", params.Email, user.Email)
	}
}
