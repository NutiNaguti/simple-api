package controllers

import (
	"encoding/json"
	"net/http"
	"simple-api/models"
	u "simple-api/utils"
)

var AddContacts = func(w http.ResponseWriter, r *http.Request) {
	contacts := &models.Contacts{}
	err := json.NewDecoder(r.Body).Decode(contacts)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
	}

	resp := models.AddContact(contacts.PhoneNum, contacts.Name, contacts.Email)
	u.Respond(w, resp)
}

var GetContacts = func(w http.ResponseWriter, r *http.Request) {
	contacts := &models.Contacts{}
	err := json.NewDecoder(r.Body).Decode(contacts)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
	}

	resp := models.GetContact(contacts.Email)
	u.Respond(w, resp)
}
