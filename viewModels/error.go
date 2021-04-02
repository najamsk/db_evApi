package viewmodels

type Error struct {
	DisplayErrors []string
	InnerErrors []string
	ApiStatusCode int
	Warnings []string
}