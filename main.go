package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connectionString = "mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"
const currentDB = "haexrdb"

func main() {
	var ServerOK = false
	log.Print("> Starting the Haexr Servers...")
	server := fiber.New()
	log.Print("> Server Loaded")

	log.Print("> Connecting to Databases...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Print("> Connection Failed")
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Print("> Connection attempt failed, Disconnecting.")
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Print("> Cannot Ping")
	} else {
		ServerOK = true
	}

	if ServerOK {
		fmt.Println("Successfully connected and pinged.")
	}

	// API

	// Root API
	server.Get("/", func(c *fiber.Ctx) error {
		if ServerOK {
			return c.Render("./index.html", nil, "")
		}
		return c.SendString("System not OK ...")
	})

	// User SignUp API
	server.Post("/signup/:apikey", func(c *fiber.Ctx) error {
		res := Users{}
		json.Unmarshal(c.Body(), &res)
		return c.Send(c.Body())
	})

	// SignUpUser(client.Database(currentDB), bson.D{{"name", "ahmed"}})
	server.Listen(":3000")
}

func SignUpUser(db *mongo.Database, user bson.D) {
	_, err := db.Collection("PersonalDetails").InsertOne(
		context.TODO(), user,
	)
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("Success")
	}
}

type PersonalDetails struct {
	name string
	age  string
}
