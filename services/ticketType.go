package services

import (
	"github.com/najamsk/eventvisor/eventvisor.api/utils"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	"fmt"
    "github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	_"time"
	"github.com/satori/go.uuid"
	"strings"
	"strconv"
)
type TicketType struct {}
func (tType TicketType) GetAvailableTicketTypes(conferenceId uuid.UUID, clientId uuid.UUID) ([]models.TicketType, error) {
	db := utils.GetDb()
	fmt.Println("dbConnection2")
	fmt.Println(&db)
	var ticketTypes []models.TicketType
	//var err = db.Where("is_active = true and conference_id = ?", conferenceId).Find(&ticketTypes).Association("Friends").Count().Error
	var err = db.Where(` is_active = true and conference_id = ? and 
							EXISTS(select 1 from tickets 
							where 
							conference_id = ?
							and client_id = ?
							and deleted = false 
						 	and deleted_at is null 
							and tickets.valid_to >= now() 
							and tickets.ticket_type_id = ticket_types.id 
							and is_active = true 
							and (booked_by is null or booked_by = '00000000-0000-0000-0000-000000000000') ) `, conferenceId, conferenceId, clientId).Order("amount DESC").Find(&ticketTypes).Error
	fmt.Println(err)
	if err != nil{
		if gorm.IsRecordNotFoundError(err){
			return nil, nil
		}else{
		return nil, err}
	}
		
	return ticketTypes, nil
}

func (tType TicketType) GetTicketTypeBooking(userId uuid.UUID, conferenceId uuid.UUID) ([]viewmodels.TicketBookedVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var bookings []viewmodels.TicketBookedVM
	
	var query string = ` select tt.title as ticket_type, count(distinct t.id) as booked_count, tt.amount as unit_price, sum(tt.amount) as total_price, 
						 tt.ammount_currency as currency, tt.valid_from, tt.valid_to
				from ticket_types tt
				inner join tickets t on tt.id = t.ticket_type_id
				where 	tt.is_active = true 
				and 	tt.deleted = false 
				and 	tt.deleted_at is null
				and 	tt.conference_id = ?
				and 	t.is_active=true and  booked_by=?
				and 	t.deleted = false and t.deleted_at is null	
				and 	t.is_consumed = false 
				and 	t.consumed_by is null
				and 	t.valid_to >= now() 
				group by tt.title, tt.amount, tt.ammount_currency, tt.valid_from, tt.valid_to
				having count(distinct t.id) >0 ` 

	db.Raw(query, conferenceId, userId).Scan(&bookings)
	fmt.Println(bookings)	
	return bookings, nil
}
func (tType TicketType) GetConsumedTicketCount(userId uuid.UUID, conferenceId uuid.UUID) (int) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	type Result struct {
		ConsumedCount int
	   }

	   var result Result
	
	var query string = ` select count(distinct t.id) as consumed_count
				from ticket_types tt
				inner join tickets t on tt.id = t.ticket_type_id
				where 	tt.is_active = true 
				and 	tt.deleted = false 
				and 	tt.deleted_at is null
				and 	tt.conference_id = ?
				and 	t.is_active=true 
				and  	t.booked_by=?
				and 	t.deleted = false and t.deleted_at is null	
				and 	t.is_consumed = true 
				and 	t.consumed_by = ?  ` 

	db.Raw(query, conferenceId, userId, userId).Scan(&result)
	fmt.Println("GetConsumedTicketCount:",result)	
	return result.ConsumedCount
}

func (tType TicketType) GetUserTicketStat(userId uuid.UUID, conferenceId uuid.UUID) (viewmodels.TicketStatVM) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var ticketStat viewmodels.TicketStatVM
	
	var query string = ` select 
							sum(cast(t.is_consumed as int8)) as consumed_count,
							count(CASE WHEN not t.is_active and not t.is_consumed THEN 1 END) as inactive_count,
							count(CASE WHEN t.valid_to < now() and not t.is_consumed THEN 1 END) as expired_count
				from ticket_types tt
				inner join tickets t on tt.id = t.ticket_type_id
				where 	tt.is_active = true 
				and 	tt.deleted = false 
				and 	tt.deleted_at is null
				and 	tt.conference_id = ?
				--and 	t.is_active=true 
				and  	t.booked_by=?
				and 	t.deleted = false and t.deleted_at is null	
				--and 	t.is_consumed = true 
				--and 	t.consumed_by = ?  ` 

	db.Raw(query, conferenceId, userId).Scan(&ticketStat)
	fmt.Println("GetUserTicketStat:",ticketStat)	
	return ticketStat
}



func (tType TicketType) BookTicket(userId uuid.UUID, conferenceId uuid.UUID, clientId uuid.UUID, memberId uuid.UUID, amountPaid float64, ticketTypeId uuid.UUID) (*uuid.UUID, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	type Result struct {ID *uuid.UUID}

	var result Result	
	var query string = ` update tickets 
						 set booked_by = ?, sold_by=?, amount_paid=?, booked_at = now()
						 where ticket_type_id = ? 
						 and conference_id = ? 
						 and client_id = ? 
						 and is_active = true 
						 and deleted = false 
						 and deleted_at is null 
						 and (booked_by is null or booked_by = '00000000-0000-0000-0000-000000000000')
						 and valid_to >= now()
						 limit 1 returning id `

	db.Raw(query, memberId, userId, amountPaid, ticketTypeId, conferenceId, clientId).Scan(&result)
	fmt.Println("BookTicket:",result)	
	return result.ID, nil
}

func (tType TicketType) ConsumeTicket(userId uuid.UUID, conferenceId uuid.UUID, clientId uuid.UUID, memberId uuid.UUID) (int, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	type Result struct {ID *uuid.UUID}

	var result []Result	
	var query string = ` update tickets 
						 set is_consumed = true,  consumed_by = ?, consumed_at = now()
						 where conference_id = ? 
						 and client_id = ? 
						 and is_active = true 
						 and deleted = false 
						 and deleted_at is null 
						 and booked_by = ?
						 and valid_to >= now()
						 and  is_consumed = false 
						 returning id `

	err := db.Raw(query, memberId, conferenceId, clientId, memberId).Scan(&result).Error
	fmt.Println("Consumed.Error:",err)
	fmt.Println("Consumed:",result)	
	return len(result), nil
}

func (tType TicketType) GetTicket(id *uuid.UUID) (*models.Ticket, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var ticketdb models.Ticket
	
	var err = db.Where(` id = ? `, id).Find(&ticketdb).Error//db.Find(&ticketdb, "ID = ?", id).Error
	fmt.Println(err)
	if err != nil{
		return nil, err
	}
	
	return &ticketdb, nil
}

func (tType TicketType) GetUserBookedTicket(memberId uuid.UUID) ([]*models.Ticket, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var ticketdb []*models.Ticket
	
	var err = db.Where(` booked_by = ? and 
						 is_consumed = false and 
						 consumed_by is null and 
						 is_active = true and 
						 deleted = false and 
						 deleted_at is null and 
						 and valid_to >= now() `, memberId).Find(&ticketdb).Error
	fmt.Println(err)
	if err != nil{
		return nil, err
	}
	
	return ticketdb, nil
}

func (tType TicketType) CheckTicketAvailability(userId uuid.UUID, clientId uuid.UUID, conferenceId uuid.UUID, tickets []viewmodels.Ticket) bool {
	db := utils.GetDb()
	
	fmt.Println(&db)

	type Result struct {
		TicketsAvailable int
	   }

	   var result Result
	   //var ticket_required int = 0;
	var query strings.Builder
	query.WriteString(" select count(count) as tickets_available from ( ")
	for i, ticket := range tickets {
		query.WriteString(" select count(ticket_type_id) from tickets where ");
		query.WriteString(" conference_id = '"+conferenceId.String()+"'");
		query.WriteString(" and client_id = '"+clientId.String()+"'");
		query.WriteString(" and tickets.ticket_type_id = '"+ticket.TicketTypeId+"'");
		query.WriteString(" and deleted = false and deleted_at is null and tickets.valid_to >= now() and is_active = true ");
		query.WriteString(" and (booked_by is null or booked_by = '00000000-0000-0000-0000-000000000000')");
		query.WriteString(" and reserved_by is null");
		query.WriteString(" group by ticket_type_id having count(ticket_type_id)> "+ strconv.FormatInt(int64(ticket.Quantity), 10));
		if((i+1)<len(tickets)){
			query.WriteString(" union ");
		}

		//ticket_required = ticket_required + ticket.Quantity
	}
	
	query.WriteString(" ) ");
	fmt.Println("query:", query.String());
	db.Raw(query.String()).Scan(&result)
	fmt.Println("tickets_available:", result);
	return result.TicketsAvailable == len(tickets) && result.TicketsAvailable >0
}

func (tType TicketType) ReserveTickets(userId uuid.UUID, clientId uuid.UUID, conferenceId uuid.UUID, tickets []viewmodels.Ticket, ticketExpiryInSec int) ([]string, int) {
	db := utils.GetDb()
	
	fmt.Println(&db)

	type Result struct {
		ID uuid.UUID
	   }

	   var result []Result
	   var ticket_required int = 0;
	   var ticketIds []string

	   

	   for i, ticket := range tickets {
		var query strings.Builder
		fmt.Println("i:", i);
		query.WriteString(" update tickets set  reserved_by = '"+userId.String()+"' ");
		query.WriteString(", reserved_at = now(),  reserve_expire_at = now() +'"+strconv.FormatInt(int64(ticketExpiryInSec), 10)+"s'");
		query.WriteString(" where  conference_id = '"+conferenceId.String()+"'");
		query.WriteString(" and client_id = '"+clientId.String()+"'");
		query.WriteString(" and tickets.ticket_type_id = '"+ticket.TicketTypeId+"'");
		query.WriteString(" and deleted = false and deleted_at is null and tickets.valid_to >= now() ");
		query.WriteString(" and is_active = true  ");
		query.WriteString(" and (booked_by is null or booked_by = '00000000-0000-0000-0000-000000000000')  ");
		query.WriteString(" and reserved_by is null ");
		query.WriteString(" limit "+ strconv.FormatInt(int64(ticket.Quantity), 10)+" returning id")
		db.Raw(query.String()).Scan(&result)

		ticket_required = ticket_required + ticket.Quantity

		for _, res := range result {
			ticketIds = append(ticketIds, res.ID.String())
		}

	   }
		
		//fmt.Println("reserved result:", len(result));

	   return ticketIds, ticket_required
}

func (tType TicketType) CancelReserveTickets(userId uuid.UUID, clientId uuid.UUID, conferenceId uuid.UUID, ticketIds []string, timelimitInSec int) int {
	db := utils.GetDb()
	fmt.Println("CancelReserveTickets:");
	fmt.Println(&db)

	type Result struct {
		ID uuid.UUID
	   }

	   var result []Result
	   
	   var query strings.Builder

		query.WriteString(" update tickets set  reserved_by = null, reserved_at = null, reserve_expire_at = null ");
		query.WriteString(" where  conference_id = '"+conferenceId.String()+"'");
		query.WriteString(" and client_id = '"+clientId.String()+"'");
		query.WriteString(" and deleted = false and deleted_at is null and tickets.valid_to >= now() ");
		query.WriteString(" and is_active = true  ");
		query.WriteString(" and (booked_by is null or booked_by = '00000000-0000-0000-0000-000000000000')  ");
		query.WriteString(" and reserved_by = '"+userId.String()+"'");
		if(len(ticketIds)>0){
			query.WriteString(" and id in ( '"+strings.Join(ticketIds, "','")+"')");
		}
		if(timelimitInSec>0){
			query.WriteString(" and reserved_at + '"+strconv.FormatInt(int64(timelimitInSec), 10)+"s' < now() ");
		}
		query.WriteString(" returning id")
		fmt.Println("query:",query.String());
		db.Raw(query.String()).Scan(&result)

	   return len(result);
}

func (tType TicketType) GetTicketTypes(ids []string) ([]*models.TicketType, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	var ticketTypeList []*models.TicketType
	
	//var err = db.Where(` id in ( ? ) `, strings.Join(ids, "','")).Find(&ticketTypeList).Error//db.Find(&ticketdb, "ID = ?", id).Error
	db.Raw("SELECT * FROM ticket_types WHERE id in( '"+strings.Join(ids, "','")+"' )").Scan(&ticketTypeList)
	// fmt.Println(err)
	// if err != nil{
	// 	return nil, err
	// }
	
	return ticketTypeList, nil
}

func (tType TicketType) UpdateTicketBooking(userId uuid.UUID, clientId uuid.UUID, conferenceId uuid.UUID, ticketIds []string) ([]string, error) {
	db := utils.GetDb()
	fmt.Println(&db)

	type Result struct {
		SerialNo string
	   }

	   var result []Result
	   var tickets []string
	   var query strings.Builder

		query.WriteString(" update tickets set  booked_by = '"+userId.String()+"', booked_at = now() ");
		query.WriteString(" where  conference_id = '"+conferenceId.String()+"'");
		query.WriteString(" and client_id = '"+clientId.String()+"'");
		query.WriteString(" and deleted = false and deleted_at is null and tickets.valid_to >= now() ");
		query.WriteString(" and is_active = true  ");
		query.WriteString(" and (booked_by is null or booked_by = '00000000-0000-0000-0000-000000000000')  ");
		query.WriteString(" and reserved_by = '"+userId.String()+"'");
		if(len(ticketIds)>0){
			query.WriteString(" and id in ( '"+strings.Join(ticketIds, "','")+"')");
		}
		query.WriteString(" returning serial_no")
		fmt.Println("UpdateTicketBooking query:",query.String());
		err := db.Raw(query.String()).Scan(&result).Error
		for _, res := range result {
			tickets = append(tickets, res.SerialNo)
		}
	   return tickets, err;
}

func (tType TicketType) GetUserAllTicketsByConference(userId uuid.UUID, conferenceId uuid.UUID) ([]viewmodels.TicketVM, error) {
	db := utils.GetDb()
	
	fmt.Println(&db)
	//var conferencedb models.Conference
	var tickets []viewmodels.TicketVM
	
	var query string = ` select t.serial_no, t.valid_to,  t.is_consumed, tt.amount as price, tt.ammount_currency as currency, (t.valid_to < now()) as is_expire
	from tickets t
	inner join ticket_types tt on t.ticket_type_id = tt.id
	where t.booked_by = ? and t.conference_id = ? and t.is_active =true and tt.is_active =true order by t.booked_at desc` 

	db.Raw(query, userId, conferenceId).Scan(&tickets)
	fmt.Println(tickets)	
	return tickets, nil
}
