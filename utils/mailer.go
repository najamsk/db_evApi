package utils

import (
	"gopkg.in/gomail.v2"
	_"strings"
	"fmt"
	_"bytes"
	"github.com/najamsk/eventvisor/eventvisor.api/internal/config"
	"html/template"
	"bytes"
	"crypto/tls"
	"github.com/najamsk/eventvisor/eventvisor.api/models"
	"github.com/satori/go.uuid"
	"github.com/najamsk/eventvisor/eventvisor.api/viewModels"
	
)

func SendEmailWithDBTemplate(templateName string, data interface{}, clientId uuid.UUID, conferenceId uuid.UUID, tolist []string, fromlist []string, config *config.Config){
	db := GetDb();

	var smtpSetting viewmodels.EmailSmtpVM
	var smtp models.ClientSmtpSetting
	var smtp_err = db.Where("client_id = ? and conference_id = ? and is_active = ?", clientId, conferenceId, true).First(&smtp).Error
	if(smtp_err == nil){
		smtpSetting.Host = smtp.Host;
		smtpSetting.Port = smtp.Port;
		smtpSetting.UserName = smtp.UserName;
		smtpSetting.Password = smtp.Password;
		smtpSetting.EmailFrom = smtp.EmailFrom;
	}
	var template models.EmailTemplate
	var _ = db.Where("client_id = ? and conference_id = ? and name = ? and is_active = ?", clientId, conferenceId, templateName, true).First(&template).Error

	if(smtp_err ==nil && len(template.EmailBody)>0){
			var body string
				body, _ = ParseEmailTemplateData(data, templateName, template.EmailBody);
				SendEmail(tolist, fromlist, template.Subject, body, smtpSetting);
			
			
	}
}
//BasicAuth ead these from config. and setup a func in utils to export this func
func SendEmail(emailto []string, emailcc []string, subject string, emailbody string, smtp viewmodels.EmailSmtpVM) (error){
	fmt.Println("emailto:", emailto)
	fmt.Println("emailcc:", emailcc)
	m := gomail.NewMessage()

	//format to email address list
	addresses_to := make([]string, len(emailto))
	for i, recipient := range emailto {
		addresses_to[i] = m.FormatAddress(recipient, "")
	}
	
	m.SetHeader("To", addresses_to...)

	//format cc email address list
	fmt.Println("len(emailcc):", len(emailcc))
	if(len(emailcc)>0){
	addresses_cc := make([]string, len(emailcc))
	for i, recipient := range emailcc {
		addresses_cc[i] = m.FormatAddress(recipient, "")
	}
	
	m.SetHeader("Cc", addresses_cc...)
}

	m.SetHeader("From", smtp.EmailFrom)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailbody)
		
	d := gomail.NewDialer(smtp.Host, smtp.Port, smtp.UserName, smtp.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		//panic(err)
		return err
	}
	return nil
}

func ParseEmailTemplate(data interface{}, templateName string) (string, error){

	t := template.New(templateName)
	var tBasePath = "./emailTemplates/";

	var err error
	t, err = t.ParseFiles(tBasePath + templateName)
	if err != nil {
		fmt.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()
	fmt.Println(result)
	return result, err;
}

func ParseEmailTemplateData(data interface{}, templateName string, html string) (string, error){

	t := template.New(templateName)
	//var tBasePath = "./emailTemplates/";

	var err error
	t, err = t.Parse(html)
	if err != nil {
		fmt.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()
	fmt.Println(result)
	return result, err;
}
