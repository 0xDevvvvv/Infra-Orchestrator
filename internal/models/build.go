package models

import "time"

type BuildStatus string

const (
	Pending BuildStatus = "PENDING"
)

type Build struct {
	ID        string
	RepoURL   string
	Branch    string
	Status    BuildStatus
	CreatedAt time.Time
}
