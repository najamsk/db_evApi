package utils

import(
	"fmt"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"encoding/json"
	"github.com/satori/go.uuid"
	"time"
	//"strings"
)

func IsInRole(userRole []*models.Role, configRole []string) bool {

	var userInRole bool = false;

	for i, role := range configRole {
		fmt.Println("i:", i)
		fmt.Println("role:", role)

		for j, urole := range userRole {
			fmt.Println("j:", j)
			fmt.Println("urole:", urole.Name)
			if(role == urole.Name){
				userInRole = true;
				break;
			}
		}
		if(userInRole == true){
			break;
		}
    }
	return userInRole
  }

  func ConvertToJSON(s interface{}) string {

	b, err := json.Marshal(s)
    if err == nil {
        return string(b);
	}
	return ""
  }

  func SendTicketPaymentConfirmationEmail(toUser *models.User, conference models.Conference, clientID uuid.UUID, amountPaid float64, currency string, contact *models.Conferences_contacts,  tickets string){

	if(contact == nil){
		contact = &models.Conferences_contacts{}
	}

	type TemplateData struct {FirstName string; LastName string; Email string; Phone string; SupportEmail string; AmountPaid float64; Currency string; Date string; SerialNo string; ConferenceTitle string; WebsiteLink string; ContactWebLink string; Facebook string; Twitter string;}	
		data := TemplateData{toUser.FirstName, toUser.LastName, toUser.Email, contact.PhoneNumber, contact.EmailSupport, amountPaid, currency, time.Now().Format("02 Jan 2006"), tickets, conference.Title, contact.Web, contact.ContactWebLink,conference.SocialMedia.Facebook, conference.SocialMedia.Twitter}
		SendEmailWithDBTemplate("ticket_payment_confirmation", data, clientID, conference.ID, []string{toUser.Email}, []string{}, nil)

  }