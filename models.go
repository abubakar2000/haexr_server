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
	Deposit_cash int // deposite case payment gateway
	Winning_cash int // withdraw
	Bonus_cash   int // coupons from referrer coupons and ads
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

// -------------

type Matches struct {
	MatchType []string
	Stages    []string
	Groups    []string
	Rounds    []string
	Results   []string
}

// ----
type Tournaments struct {
	Banner                string
	Title                 string
	Sponsor               string
	Entrancefee           string
	RegistrationStartDate string
	RegistrationLastDate  string
	TournamentStartDate   string
	TournamentEndDate     string
	TournamentsTeamType   string
	EligibleCountries     []string
	TotalTeams            int
	StreamLinks           []StreamLink
	Teams                 []Team //to be considered
	Winnings              string
	Rounds                []Rounds //the cards in it
	PointTable            string
	Tier                  string
}

type StreamLink struct {
	Platform string
	Url      string
}

// Qualifies -> semi final -> final thing //stages
type Rounds struct {
	date   []string
	time   []string
	Groups []Groups
}

type Groups struct {
	GroupID        string
	MatchID        string //BGMI MATCH #7768
	Group          string
	teams          []Team
	rounds         []string // the rounds to be played in beetween the pool of teams coming froma action sheet from below
	results        []string // will conatain the screenshots and some data
	MapName        string
	StartingAtTime string
	StartingAtDate string
	RoomID         string
	Password       string
}
