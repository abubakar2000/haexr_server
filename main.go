package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/qinains/fastergoding"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connectionString = "mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"
const currentDB = "haexrdb"
const Success = 200
const NotAcceptable = 406

func main() {
	fastergoding.Run()

	var ServerOK = false
	log.Print("> Starting the Haexr Servers...")

	server := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("\nShutdown Command in progress...")
		_ = server.Shutdown()
	}()

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

	// Root API
	server.Get("/", func(c *fiber.Ctx) error {
		if ServerOK {
			return c.Render("./index.html", nil, "")
		}
		return c.SendString("System not OK ...")
	})
	server.Get("/test", func(c *fiber.Ctx) error {
		log.Println("Tested")
		return c.SendString("Test Successful")
	})

	// -----------------------------------------------------------------

	// User SignUp API
	server.Post("/register", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), &userData)
		if c.Query("code") != "" {
			// ?code=REFERCODE
			if SignUpWithCode(client.Database(currentDB), userData, c.Query("code")) {
				return c.SendStatus(Success)
			}
			return c.SendStatus(NotAcceptable)
		} else {
			if SignUpUser(client.Database(currentDB), userData) {
				return c.SendStatus(Success)
			}
			return c.SendStatus(NotAcceptable)
		}

	})

	server.Static("/", "./public")

	// User unregister API
	server.Post("/unregister", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		if UnRegUser(client.Database(currentDB), userData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/login", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)

		if FindUser(client.Database(currentDB), userData) {
			// Token will be given
			return c.SendStatus(Success)
		} else {
			return c.SendStatus(NotAcceptable)
		}
	})

	server.Post("/logout", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		fmt.Println(userData.Email)
		fmt.Println(userData.Password)
		if FindUser(client.Database(currentDB), userData) {
			// Token will be revoked
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/updateuser", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		fmt.Println(userData.Email)
		fmt.Println(userData.Password)
		if UpdateUser(client.Database(currentDB), userData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Team
	server.Post("/addteam", func(c *fiber.Ctx) error {
		teamData := &Team{}
		json.Unmarshal(c.Body(), teamData)
		if AddTeam(client.Database(currentDB), teamData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Team Member
	server.Post("/addteammember", func(c *fiber.Ctx) error {
		newMemberData := &User{}
		json.Unmarshal(c.Body(), newMemberData)
		if AddTeamMember(client.Database(currentDB), newMemberData, c.Query("teamid")) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// remove teammember
	server.Post("/delteammember", func(c *fiber.Ctx) error {
		newMemberData := &User{}
		json.Unmarshal(c.Body(), newMemberData)
		if DelTeamMember(client.Database(currentDB), newMemberData, c.Query("teamid")) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Users GameInformation
	server.Post("/addgameinformationofuser", func(c *fiber.Ctx) error {
		gameInformationOfUser := &GameInformationOfUser{}
		json.Unmarshal(c.Body(), gameInformationOfUser)
		if AddUsersGameInfo(client.Database(currentDB), gameInformationOfUser) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Users Game
	server.Post("/addgame", func(c *fiber.Ctx) error {
		gameInfo := &Game{}
		json.Unmarshal(c.Body(), gameInfo)
		if AddGame(client.Database(currentDB), gameInfo) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// Add Users Game
	server.Post("/addtransaction", func(c *fiber.Ctx) error {
		transactionInfo := &Transaction{}
		json.Unmarshal(c.Body(), transactionInfo)
		if addTransaction(client.Database(currentDB), transactionInfo) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/addreference", func(c *fiber.Ctx) error {
		referenceInfo := &Refer{}
		json.Unmarshal(c.Body(), referenceInfo)
		if addReference(client.Database(currentDB), referenceInfo) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Listen(":3000")

}
