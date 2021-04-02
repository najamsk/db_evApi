package viewmodels

// import(
// 	"github.com/satori/go.uuid"
// )

type ConferenceAgendaVM struct {
	StartDate    string
	Title string
	Sessions	[]SessionVM
}