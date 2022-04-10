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
	TotalTeams            int // total number of teams that can join this tournament
	StreamLinks           []StreamLink
	Teams                 []Team //to be considered, also will be broken down into groups
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
// can be of 3 hours maybe for 3 group added by kundan himself
type Rounds struct {
	QualifierName                 string
	Date                          []string //spans over Group StartingAtDate
	Time                          []string //spans over Group StartingAtTime
	Groups                        []Groups
	NumOfQualifyingTeamsThisRound int //how many teams will be qualifying for this round
	MapName                       string
	IsLocked                      bool
	NumberOfTeamsPerGroup         int
}

type Groups struct {
	GroupID        string
	MatchID        string //BGMI MATCH #7768
	Group          string
	Teams          []Team
	Rounds         []string // the rounds to be played in beetween the pool of teams coming froma action sheet from below
	Results        []string // will conatain the screenshots and some data
	StartingAtTime string
	StartingAtDate string
	RoomID         string
	Password       string
}
