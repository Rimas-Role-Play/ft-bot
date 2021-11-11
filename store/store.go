package store

type Player struct {
	DSUid string
	Uid uint32
	SteamId string
	Name string
}

type PlayerStats struct {
	PlayerInfo Player
	GroupId int8
	DonatLevel int8
}

// Groups
type GroupRoles struct {
	GroupId uint
	DiscordRoleLeader string
	DiscordRoleMember string
}
