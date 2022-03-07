package store

type Player struct {
	DSUid   string
	Uid     uint32
	SteamId string
	Name    string
}

type PlayerStats struct {
	PlayerInfo Player
	GroupId    int8
	DonatLevel int8
}

// Groups
type GroupRoles struct {
	GroupId           uint
	DiscordRoleLeader string
	DiscordRoleMember string
}

// Vehicles
type PremiumVehicles struct {
	Classname   string
	Name        string
	Images      []string
	Description string
	Price       uint16
	Discount    uint8
}

type Vehicles struct {
	Classname   string
	Image       string
	DisplayName string
}
