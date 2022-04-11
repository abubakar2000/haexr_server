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

	// User unregister API
	server.Post("/unregister", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		if UnRegUser(client.Database(currentDB), userData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("getuserinfo", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
		return c.JSON(GetUserDetails(client.Database(currentDB), userData.Email))

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
		if FindUser(client.Database(currentDB), userData) {
			// Token will be revoked
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/updateuser", func(c *fiber.Ctx) error {
		userData := &User{}
		json.Unmarshal(c.Body(), userData)
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

	server.Get("/enlistgames", func(c *fiber.Ctx) error {
		gameInfo := &Game{}
		json.Unmarshal(c.Body(), gameInfo)
		GamesList := GetGame(client.Database(currentDB))
		return c.JSON(GamesList)

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

	server.Post("/createteam", func(c *fiber.Ctx) error {
		var tempData Team
		json.Unmarshal(c.Body(), &tempData)
		if CreateTeams(client.Database(currentDB), tempData) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Get("/getteams_whole", func(c *fiber.Ctx) error {
		return c.JSON(GetTeamsWhole(client.Database(currentDB)))
	})

	server.Get("/getteam", func(c *fiber.Ctx) error {
		type TeamNameHolder struct {
			TeamName string
		}
		var teamname TeamNameHolder
		json.Unmarshal(c.Body(), &teamname)
		return c.JSON(GetTeamByName(client.Database(currentDB), teamname.TeamName))
	})

	server.Get("/getteambygameid", func(c *fiber.Ctx) error {
		type GameIDHolder struct {
			GameID string
		}
		var gameid GameIDHolder
		json.Unmarshal(c.Body(), &gameid)
		return c.JSON(GetTeamsByGameID(client.Database(currentDB), gameid.GameID))
	})

	server.Post("/addmembertoteam", func(c *fiber.Ctx) error {
		type UserAndTeam struct {
			User string
			Team string
		}
		var userandteam UserAndTeam
		json.Unmarshal(c.Body(), &userandteam)
		if AddUserToTeam(client.Database(currentDB), userandteam.User, userandteam.Team) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/removemembertoteam", func(c *fiber.Ctx) error {
		type UserAndTeam struct {
			User string
			Team string
		}
		var userandteam UserAndTeam
		json.Unmarshal(c.Body(), &userandteam)

		if RemoveUserFromTeam(client.Database(currentDB), userandteam.User, userandteam.Team) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/addtournament", func(c *fiber.Ctx) error {
		var tournaments Tournaments
		json.Unmarshal(c.Body(), &tournaments)
		if AddTournament(client.Database(currentDB), tournaments) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/addstreamlinkintournament", func(c *fiber.Ctx) error {
		type StreamLinkAndTornament struct {
			Tournament string
			Link       StreamLink
		}
		var streamLinkAndTournament StreamLinkAndTornament
		json.Unmarshal(c.Body(), &streamLinkAndTournament)
		if AddStreamingLinksToTournament(client.Database(currentDB), streamLinkAndTournament.Tournament, streamLinkAndTournament.Link) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	// should add team to qualifier only
	// server.Post("/addteamintournament", func(c *fiber.Ctx) error {
	// 	type TeamAndTornament struct {
	// 		Tournament string
	// 		Team       Team
	// 	}
	// 	var teamAndTournament TeamAndTornament
	// 	json.Unmarshal(c.Body(), &teamAndTournament)
	// 	if AddTeamToTournament(client.Database(currentDB), teamAndTournament.Tournament, teamAndTournament.Team) {
	// 		return c.SendStatus(Success)
	// 	}
	// 	return c.SendStatus(NotAcceptable)
	// })

	server.Get("/gettournaments", func(c *fiber.Ctx) error {
		return c.JSON(GetTournaments(client.Database(currentDB)))
	})

	server.Get("/gettournament", func(c *fiber.Ctx) error {
		type TournamentBody struct {
			Tournament string
		}
		tournament := TournamentBody{}
		json.Unmarshal(c.Body(), &tournament)
		return c.JSON(GetTournament(client.Database(currentDB), tournament.Tournament))
	})

	server.Get("/gettournamentbygame", func(c *fiber.Ctx) error {
		type TournamentBody struct {
			GameID string
		}
		tournament := TournamentBody{}
		json.Unmarshal(c.Body(), &tournament)
		return c.JSON(GetTournamentByGame(client.Database(currentDB), tournament.GameID))
	})

	server.Post("/addqualifierroundintournament", func(c *fiber.Ctx) error {
		type QualifierAndRounds struct {
			Tournament string
			Qualifier  Rounds
		}
		qualifierAndRound := QualifierAndRounds{}
		json.Unmarshal(c.Body(), &qualifierAndRound)
		if AddQualifierRoundInTournament(client.Database(currentDB), qualifierAndRound.Tournament,
			qualifierAndRound.Qualifier) {
			return c.SendStatus(Success)
		}
		return c.SendStatus(NotAcceptable)
	})

	server.Post("/addteamintournamentgroup", func(c *fiber.Ctx) error {
		type Body struct {
			Tournament string
			Qualifier  string
			Group      Groups
			Team       Team
		}
		teamInQualOfTournament := Body{}
		json.Unmarshal(c.Body(), &teamInQualOfTournament)
		return c.JSON(AddTeamInTournamentGroup(client.Database(currentDB), teamInQualOfTournament.Tournament,
			teamInQualOfTournament.Qualifier, teamInQualOfTournament.Group, teamInQualOfTournament.Team))
	})

	server.Static("/", "./public")
	server.Listen(":3000")

}
