package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nytro04/hotel-reservation/db"
	"github.com/nytro04/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbname = "hotel-reservation"

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Good hope",
		Location: "Accra",
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 199.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 299.9,
		},
	}

	for _, room := range rooms {

		room.HotelID = insertedHotel.ID

		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(insertedRoom)
	}
	fmt.Println(insertedHotel)

}
