package models

import (
	"github.com/jinzhu/gorm"
	"log"
	u "simple-api/utils"
)

type Contacts struct {
	gorm.Model
	UserId uint `json:"user_id"`
	PhoneNum string `json:"phoneNum"`
	Name     string `json:"name"`
}

func (contact *Contacts) Validate() (map[string]interface{}, bool) {
	if contact.Name == "" {
		return u.Message(false, "Contact name should be on the payload"), false
	}
	if contact.PhoneNum == "" {
		return u.Message(false, "Contact phone number should be on the payload"), false
	}
	if contact.UserId <= 0 {
		return u.Message(false, "Contact user id should be on the payload"), false
	}

	return u.Message(true, "Success"), true
}

func (contact *Contacts) Create() map[string]interface{} {
	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := u.Message(true, "Success")
	resp["contact"] = contact
	return resp
}

func GetContact(id uint) *Contacts {
	contact := &Contacts{}
	err := GetDB().Table("contacts").Where("id = ?", id).Find(contact).Error
	if err != nil {
		return nil
	}
	return contact
}

func GetContacts(id uint) []*Contacts {
	contact := make([]*Contacts, 0)
	err := GetDB().Table("contacts").Where("id = ?", id).Find(contact).Error
	if err != nil {
		log.Println(err)
		return nil
	}
	response := u.Message(true, "Success")
	response["contacts"] = contact
	return contact

}

func AddContact(phoneNum, name, email string) map[string]interface{} {
	contact := &Contacts{}
	contact.PhoneNum = phoneNum
	contact.Name = name
	//contact.Email = email
	GetDB().Create(contact)
	log.Println("Contact has been added:\nphoneNum: " + contact.PhoneNum + "\nname: " + contact.Name)

	response := u.Message(true, "Contact has been added")
	response["contacts"] = contact
	return response
}
