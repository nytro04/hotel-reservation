package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nytro04/hotel-reservation/api"
	"github.com/nytro04/hotel-reservation/db"
	"github.com/nytro04/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongHotelStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "kwaku", "ansah", false)

	fmt.Println("kwaku", user)
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin", admin)
	fmt.Println("admin token", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "Stay hyped", "Kumasi", 5, nil)
	fmt.Println(hotel)
	room := fixtures.AddRoom(store, "large", true, 99.99, hotel.ID)
	fmt.Println(room)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println(booking)

}
