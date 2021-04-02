package models

type Client struct {
	Base
	Name         string
	IsActive     bool
	Subscription Subscription

	// CreatedAt    time.Time
	// UpdatedAt    time.Time
	// DeletedAt    time.Time
}