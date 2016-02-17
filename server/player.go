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
	Stats        Stats // Stats of the character when unaltered
	CurrentStats Stats // Current state of the character.

	ID             uint32
	Name           string
	EquippedItems  []*Item // Item by slot.
	InventoryItems []*Item // Items held in backpack.
}

// Stats are the stats that all entities in the game use.
type Stats struct {
	HP               int32 // HP is how much damage the person can take before dying
	Stamina          int16 // Stamina is the resource used to perform physical abilities
	Concentration    int16 // Concentration is the resource used to perform mental abilities
	Speed            int16 // Speed is how far an entity can travel
	MagicStrength    int16 // MagicStrength is how much power is added to magical abilities
	PhysicalStrength int16 // PhysicalStrength is how much power is added to physical abilities
}
