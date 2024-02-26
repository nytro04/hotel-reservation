package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nytro04/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// this works, but i will prefer to send a PATCH request with cancel params set to true instead of a get
func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{
		Type: "Msg",
		Msg:  "updated",
	})
}

// TODO: This needs to be Admin authorized. admin get all bookings
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

// TODO: This needs to be user authorized. i.e user get their bookings/booking
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "unauthorized",
		})
	}
	return c.JSON(booking)
}
