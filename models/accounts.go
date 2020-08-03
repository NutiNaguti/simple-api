package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	u "simple-api/utils"
	"strings"
)

/*
Структура прав доступа jWT
*/
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// структура для учетной записи пользователя
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

// проверить входящие данные пользователя
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	// Email должен быть уникальным
	temp := &Account{}

	// проверка на наличие ошибок и дубликатов электронных писем
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	// создать новый токен для новой учетной записи
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "******" // удаление пароля

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return u.Message(false, "Invalid login creditials. Please try again")
	}

	// Все работает! Войти в систему
	account.Password = "******"

	// Создать токен
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	account.Token = tokenString // сохранить токен в ответе
	log.Println(account.Token)

	resp := u.Message(true, "Logged in")
	resp["account"] = account
	return resp
}

func GetUser(u uint) *Account {
	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" {
		return nil
	}

	acc.Password = "******"
	return acc
}
