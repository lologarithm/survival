package server

// User maps a connection to a list of accounts
type User struct {
	Accounts []*Account
	Client   *Client
}

// Account is mostly a container for characters and has a password to use them.
type Account struct {
	ID         uint32
	Name       string
	Password   string
	Characters []*Character
}

// Character is a single entity in the game.
type Character struct {
	Name  string
	Items []*Item
}
