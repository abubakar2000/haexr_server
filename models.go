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

type Tournaments struct {
	Banner                []string
	TournamentName        string
	SponseredName         string
	TeamType              string
	RegistrationStartDate string
	RegistrationLastDate  string
	TournamentStartDate   string
	TournamentEndDate     string
	EntryFees             string
	TotalSlot             string
	CountryEligible       []string
	ScheduleStages        []string
	PointTable            []string
	Winnings              []string
	Tier                  string
	StreamLink            []string
}

type Stages struct {
	StartDate                string
	EndDate                  string
	IsLocked                 string
	NoOfTeamsAvailable       string
	NoOfTeamsPerGroup        string
	DateAvailable            []string
	TimeAvailable            []string
	NoOfPassingTeams         string
	NoOfRoundsPerGroup       string
	MapName                  []string
	TimeTakenByRoundPerMatch string
	GroupSlot                string //It is type of map and we will show key (slot 3-n) in frontend on All teams page
}

type Results struct {
	// Tournament.Stage.Group.RoundNumber Id
	ResultImg          []string //UserInput
	ResultPosition     []string //UserInput
	ResultKills        []string //UserInput
	StreamerScreenshot []string //AdminInput
}
