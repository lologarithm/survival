package server

// User maps a connection to a list of accounts
type User struct {
	Accounts []*Account // List of authenticated accounts
	Client   *Client    // Client connection
	GameID   uint32     // Currently connected game ID
}

// Account is mostly a container for character and has a password to use them.
type Account struct {
	ID        uint32
	Name      string
	Password  string
	Character *Character
}

// Character is a single entity in the game.
type Character struct {
	ID    uint32
	Name  string
	Items []*Item
}
