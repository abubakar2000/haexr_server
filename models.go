package main

type User struct {
	User_uuid             string //haexr_id
	Email                 string
	Password              string
	Fname                 string
	Lname                 string
	Telephone             string
	Address               string
	Country               string
	ProfileImage          string
	PreferredGames        []Game
	UserWallet            Wallet
	UsersGamesInformation []GameInformationOfUser
}

type Team struct {
	TeamID      string
	TeamName    string
	GameID      string
	UsersInTeam []string
}

type GameInformationOfUser struct {
	GameID     string
	Total_time string
	IGN        string // in game name
	ID         string // in game id
	TeamId     string // which game
}

type Game struct {
	GameID       string
	GameName     string
	GameTeamType []string
	GameLogo     string
}

// Composed in the user
type Wallet struct {
	Wallet_id    string
	Deposit_cash int
	Winning_cash int
	Bonus_cash   int
}

type Transaction struct {
	Transaction_id string
	Wallet_id      string
	Source         string
	Timestamp      string
}

type Refer struct {
	Refer_id          string
	Produce_user_uuid string // who generated this reference
	Validity          string
	Timestamp         string
	Code              string
}
