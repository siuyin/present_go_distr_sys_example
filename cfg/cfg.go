package cfg

// IMPORTANT -- ALL entries in this file are append only.
//
// CHANGING / EDITING entries require retiring the generation that used the old config.
// Think DNA. Our DNA is mostly the result responses to past pathogen attacks.
//
// Once a response is coded it stays in our DNA

// Application Constants and Service Names
const (
	App       = "EgA"
	IDOffice  = App + ".IDOffice"
	HeartBeat = App + ".HeartBeat"
)

// Ranks or Position or Function
const (
	Unassigned  = App + ".Unassigned"
	FileWatcher = App + ".FileWatcher"
	IDOfficer   = App + ".IDOfficer"
	HBListener  = App + ".HeartBeatListener"
	ManagerA    = App + ".ManagerA"
)

// Stable Files Pool
const (
	StableFilesA = App + ".StableFilesPool.A"
)

// NRS Name, Rank and Serial Number (ID)
type NRS struct {
	Name, Rank, ID string
}
