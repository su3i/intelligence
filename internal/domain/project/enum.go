package project

type ProjectStatus string

const (
	Active  ProjectStatus = "ACTIVE"
	Paused ProjectStatus = "PAUSED"
	Archived ProjectStatus = "ARCHIVED"
)

type ProjectStage string

const (
	Sandbox ProjectStage = "SANDBOX"
	Production ProjectStage = "PRODUCTION"
)

type ProjectBusinessDomain string

const (
	Marketplace ProjectBusinessDomain = "MARKETPLACE"
)