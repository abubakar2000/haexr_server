package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SignUpUser(db *mongo.Database, user *User) bool {
	status := true
	_, err := db.Collection("PersonalDetails").InsertOne(
		context.TODO(), user,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func SignUpWithCode(db *mongo.Database, user *User, code string) bool {
	status := true
	resp := db.Collection("ReferenceInfo").FindOne(context.TODO(), bson.M{"code": code})
	var refer Refer
	signeeReward := 100
	referrerReward := 200
	resp.Decode(&refer)
	if refer.Code == code {
		user.UserWallet.Bonus_cash = signeeReward

		_, err := db.Collection("PersonalDetails").InsertOne(
			context.TODO(), user,
		)
		if err != nil {
			log.Printf(err.Error())
			status = false
		} else {
			log.Printf("Success")
			status = true
			db.Collection("PersonalDetails").UpdateOne(context.TODO(),
				bson.M{"user_uuid": refer.Produce_user_uuid}, bson.M{"$inc": bson.M{"userwallet.bonus_cash": referrerReward}})
		}

	}
	return status
}

func UnRegUser(db *mongo.Database, user *User) bool {
	status := true
	log.Println("Unregister User")
	_, err := db.Collection("PersonalDetails").DeleteOne(
		context.TODO(), &fiber.Map{
			"email":    user.Email,
			"password": user.Password,
		},
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func UpdateUser(db *mongo.Database, user *User) bool {
	status := true
	_, err := db.Collection("PersonalDetails").UpdateOne(context.TODO(), bson.M{
		"email": user.Email,
	}, bson.M{"$set": user}, options.Update().SetUpsert(true))
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func FindUser(db *mongo.Database, user *User) bool {

	status := true
	resp, err := db.Collection("PersonalDetails").
		Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		status = true
	}
	for resp.Next(context.TODO()) {
		var userTemp User
		err := resp.Decode(&userTemp)
		if err != nil {
			log.Fatal(err)
		}
		if userTemp.Password == user.Password && user.Email == userTemp.Email {
			status = true
			return status
		}
	}
	status = false
	return status
}

func GetUserDetails(db *mongo.Database, email string) User {
	resp := db.Collection("PersonalDetails").FindOne(context.TODO(), bson.M{"email": email})
	var userInformation User
	resp.Decode(&userInformation)
	userInformation.Password = ""
	return userInformation
}

func AddTeam(db *mongo.Database, team *Team) bool {
	status := true
	_, err := db.Collection("Teams").InsertOne(
		context.TODO(), team,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func AddTeamMember(db *mongo.Database, teamMember *User, teamid string) bool {
	status := true
	resp, err := db.Collection("Teams").Find(context.TODO(), bson.M{"teamid": teamid})
	for resp.Next(context.TODO()) {
		var teamTemp Team
		resp.Decode(&teamTemp)
		teamTemp.UsersInTeam = append(teamTemp.UsersInTeam, teamMember.User_uuid)
		_, err := db.Collection("Teams").UpdateOne(context.TODO(),
			bson.M{"teamid": teamid}, bson.M{"$set": teamTemp})
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Added user " + teamMember.User_uuid + " to team " + teamid)
		}
	}
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func splice(array []string, nameToRem string) []string {
	index := 0
	for i := 0; i < len(array); i++ {
		if nameToRem == array[i] {
			index = i
			break
		}
	}
	return append(array[:index], array[index+1:]...)
}

func DelTeamMember(db *mongo.Database, teamMember *User, teamid string) bool {
	status := true
	resp, err := db.Collection("Teams").Find(context.TODO(), bson.M{"teamid": teamid})
	for resp.Next(context.TODO()) {
		var teamTemp Team
		resp.Decode(&teamTemp)

		teamTemp.UsersInTeam = splice(teamTemp.UsersInTeam, teamMember.User_uuid)
		_, err := db.Collection("Teams").UpdateOne(context.TODO(),
			bson.M{"teamid": teamid}, bson.M{"$set": teamTemp})
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Added user " + teamMember.User_uuid + " to team " + teamid)
		}
	}
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func AddUsersGameInfo(db *mongo.Database, gameInformationOfUser *GameInformationOfUser) bool {
	status := true
	_, err := db.Collection("UsersGameInformation").InsertOne(
		context.TODO(), gameInformationOfUser,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func AddGame(db *mongo.Database, gameInfo *Game) bool {
	status := true
	_, err := db.Collection("GameInformation").InsertOne(
		context.TODO(), gameInfo,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func GetGame(db *mongo.Database) []Game {
	GamesList := []Game{}
	list, err := db.Collection("GameInformation").Find(
		context.TODO(), bson.M{},
	)
	if err != nil {
		log.Printf(err.Error())

	} else {

		for list.Next(context.TODO()) {
			var game Game
			list.Decode(&game)
			GamesList = append(GamesList, game)
		}

	}
	return GamesList
}

func addTransaction(db *mongo.Database, transactionInfo *Transaction) bool {
	status := true
	_, err := db.Collection("TransactionInfo").InsertOne(
		context.TODO(), transactionInfo,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func addReference(db *mongo.Database, reference *Refer) bool {
	status := true
	_, err := db.Collection("ReferenceInfo").InsertOne(
		context.TODO(), reference,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func CreateTeams(db *mongo.Database, newTeam Team) bool {
	status := true
	_, err := db.Collection("Teams").InsertOne(
		context.TODO(), newTeam,
	)
	if err != nil {
		log.Printf(err.Error())
		status = false
	} else {
		log.Printf("Success")
		status = true
	}
	return status
}

func GetTeamByName(db *mongo.Database, teamName string) []string {
	res := db.Collection("Teams").FindOne(
		context.TODO(), bson.M{"teamname": teamName},
	)
	var team Team
	res.Decode(&team)
	return team.UsersInTeam
}

func GetTeamsByGameID(db *mongo.Database, gameid string) []Team {
	var teams []Team
	res, err := db.Collection("Teams").Find(
		context.TODO(), bson.M{"gameid": gameid},
	)
	if err != nil {
		return nil
	}
	for res.Next(context.TODO()) {
		var team Team
		res.Decode(&team)
		teams = append(teams, team)
	}
	return teams
}

func GetTeamsWhole(db *mongo.Database) []Team {
	TeamsList := []Team{}
	list, err := db.Collection("Teams").Find(
		context.TODO(), bson.M{},
	)
	if err != nil {
		log.Printf(err.Error())
	} else {
		for list.Next(context.TODO()) {
			var team Team
			list.Decode(&team)
			TeamsList = append(TeamsList, team)
		}
	}
	return TeamsList
}

func AddUserToTeam(db *mongo.Database, user string, team string) bool {
	_, err := db.Collection("Teams").UpdateOne(context.TODO(),
		bson.M{"teamid": team}, bson.M{"$push": bson.M{"usersinteam": user}})
	if err != nil {
		return false
	}
	return true
}

func RemoveUserFromTeam(db *mongo.Database, user string, team string) bool {
	_, err := db.Collection("Teams").UpdateOne(context.TODO(),
		bson.M{"teamid": team}, bson.M{"$pull": bson.M{"usersinteam": user}})
	if err != nil {
		return false
	}
	return true
}

func AddTournament(db *mongo.Database, tournament Tournaments) bool {
	_, err := db.Collection("Tournaments").InsertOne(context.TODO(),
		tournament)
	if err != nil {
		return false
	}
	return true
}

func AddStreamingLinksToTournament(db *mongo.Database, tournament string, steamLink StreamLink) bool {
	_, err := db.Collection("Tournaments").UpdateOne(
		context.TODO(), bson.M{"title": tournament},
		bson.M{"$push": bson.M{"streamlinks": steamLink}},
	)
	if err != nil {
		return false
	}
	return true
}

func AddTeamToTournament(db *mongo.Database, tournament string, team Team) bool {
	_, err := db.Collection("Tournaments").UpdateOne(
		context.TODO(), bson.M{"title": tournament},
		bson.M{"$push": bson.M{"teams": team}},
	)
	if err != nil {
		return false
	}
	return true
}

func GetTournament(db *mongo.Database, tournament string) Tournaments {
	res := db.Collection("Tournaments").FindOne(context.TODO(),
		bson.M{"title": tournament})
	var result Tournaments
	res.Decode(&result)
	return result
}

func GetTournaments(db *mongo.Database) []Tournaments {
	res, err := db.Collection("Tournaments").Find(context.TODO(),
		bson.M{})
	tournaments := []Tournaments{}
	tempHolder := Tournaments{}
	if err != nil {
		return nil
	}
	for res.Next(context.TODO()) {
		res.Decode(&tempHolder)
		tournaments = append(tournaments, tempHolder)
	}
	return tournaments
}

func AddQualifierRoundInTournament(db *mongo.Database, tournament string, qualifier Rounds) bool {
	res, err := db.Collection("Tournaments").UpdateOne(context.TODO(),
		bson.M{"title": tournament}, bson.M{"$push": bson.M{"rounds": qualifier}})
	if err != nil {
		return false
	}

	if res.ModifiedCount == 0 {
		return false
	}

	return true
}
