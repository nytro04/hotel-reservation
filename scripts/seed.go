package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nytro04/hotel-reservation/api"
	"github.com/nytro04/hotel-reservation/db"
	"github.com/nytro04/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func SeedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})

	user.IsAdmin = isAdmin

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))

	return insertedUser
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}

	if _, err := bookingStore.InsertBooking(context.Background(), booking); err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	// rooms := []types.Room{
	// 	{
	// 		Size:  "small",
	// 		Price: 99.9,
	// 	},
	// 	{
	// 		Size:  "normal",
	// 		Price: 199.9,
	// 	},
	// 	{
	// 		Size:  "kingsize",
	// 		Price: 299.9,
	// 	},
	// }

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	// for _, room := range rooms {
	// 	room.HotelID = insertedHotel.ID
	// 	_, err := roomStore.InsertRoom(ctx, &room)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	return insertedHotel
}

func main() {
	kwaku := SeedUser(false, "Kwaku", "Ansah", "kwaku@gmail.com", "P@ssw0rd")
	SeedUser(true, "isAdmin", "isAdmin", "admin@gmail.com", "P@ssw0rd")
	seedHotel("Kpoga", "Accra", 4)
	seedHotel("Humii", "Aburi", 5)
	hotel := seedHotel("Dzooye", "Pampram", 5)
	seedRoom("small", false, 99.99, hotel.ID)
	seedRoom("medium", true, 199.99, hotel.ID)
	room := seedRoom("large", true, 299.99, hotel.ID)
	seedBooking(kwaku.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

func init() {
	var err error

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
