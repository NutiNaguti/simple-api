package models

import (
	"github.com/jinzhu/gorm"
	"log"
	u "simple-api/utils"
)

type Contacts struct {
	gorm.Model
	Email    string `json:"email"`
	PhoneNum string `json:"phoneNum"`
	Name     string `json:"name"`
}

func AddContact(phoneNum, name, email string) map[string]interface{} {
	contact := &Contacts{}
	contact.PhoneNum = phoneNum
	contact.Name = name
	contact.Email = email
	GetDB().Create(contact)
	log.Println("Contact has been added:\nphoneNum: " + contact.PhoneNum + "\nname: " + contact.Name)

	response := u.Message(true, "Contact has been added")
	response["contacts"] = contact
	return response
}

func GetContact(email string) map[string]interface{} {
	contact := &Contacts{}
	GetDB().Table("contacts").Where("email = ?", email).Find(contact)
	if contact.Name == "" || contact.PhoneNum == "" {
		return nil
	}
	response := u.Message(true, "")
	response["contacts"] = contact
	return response
}
