package account

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Account is the table definition of user account
type Account struct {
	gorm.Model
	Username string `gorm:"size:255;index;NOT NULL" json:"username"`
	Password string `gorm:"size:255;NOT NULL" json:"password"`
	Email    string `gorm:"type:varchar(32);unique_index;" json:"email"`
	Phone    string `gorm:"type:varchar(20);DEFAULT NULL" json:"phone"`
	Balance  uint32 `gorm:"size:32;DEFAULT NULL" json:"balance"`
}

// LoginReq is the struct of login parameters
type LoginReq struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

// LoginRsp is the data server send back to client when login success
type LoginRsp struct {
	Username string `json:"name"`
	Token    string `json:"token"`
}

// Repository :defines interface for Account model operations
type Repository interface {
	Add(a *Account) (e error)
	GetByID(id string) (a *Account, e error)
	GetByName(name string) (a *Account, e error)
	GetIDByName(name string) (id string, e error)
	GetNameByID(id string) (name string, e error)
	List(line string) (a []*Account, e error)
	Save(id string, a *Account) (e error)
	Delete(id string) (e error)
}

// CreateTable :create account table and add admin account.
func (m *Mgr) CreateTable() {
	var err error
	if err = m.DB.AutoMigrate(&Account{}).Error; err != nil {
		log.Fatalln("create account table failure")
		panic("create account table failure")
	} else {
		log.Println("create account table success")
	}
	// add admin account to do some privilege operations such as delete account.
	var adminAccount Account
	if err = m.DB.Where("Username = ?", "admin").First(&adminAccount).Error; err == nil {
		log.Println("admin account already exists, do nothing")
		return
	}
	adminAccount = Account{Username: "admin", Password: "mobiwallet.com", Email: "admin@mobiwallet.com"}
	if err = m.DB.Create(&adminAccount).Error; err != nil {
		log.Fatalln("create admin account failure")
	} else {
		log.Println("create admin account success")
	}

}

// Add :insert a row into account table
func (m *Mgr) Add(a *Account) (e error) {
	var err error
	name := a.Username
	if _, err = m.GetByName(name); err == nil {
		// account already exists
		err = errors.New("user already exists, please directly login")
		return err
	}
	return m.DB.Create(a).Error
}

// GetByID :select a row from account table based on account ID
func (m *Mgr) GetByID(id string) (a Account, e error) {
	var account Account
	err := m.DB.Where("ID = ?", id).Find(&account).Error
	return account, err
}

// GetByName :select a row from account table based on user name
func (m *Mgr) GetByName(name string) (a Account, e error) {
	var account Account
	err := m.DB.Where("Username = ?", name).Find(&account).Error
	return account, err
}

// GetIDByName :returns account ID corresponding to username
func (m *Mgr) GetIDByName(name string) (id string, e error) {
	var account Account
	err := m.DB.Where("Username = ?", name).Find(&account).Error
	return string(account.ID), err
}

// GetNameByID :returns username corresponding to account ID
func (m *Mgr) GetNameByID(id string) (name string, e error) {
	var account Account
	err := m.DB.Where("ID = ?", id).Find(&account).Error
	return account.Username, err
}

// List :list all the accounts, if line != 0, select at most "line" rows in account table
func (m *Mgr) List(line string) (a []Account, e error) {
	var accounts []Account
	var err error
	if line == "" {
		err = m.DB.Find(&accounts).Error
	} else {
		err = m.DB.Limit(line).Find(&accounts).Error
	}
	return accounts, err
}

// Save :update a row.
func (m *Mgr) Save(id string, a *Account) (e error) {
	return m.DB.Set("gorm:query_option", "FOR UPDATE").Where("ID = ?", id).Save(a).Error
}

// Delete :delete a row.
func (m *Mgr) Delete(id string) (e error) {
	return m.DB.Set("gorm:query_option", "FOR UPDATE").Where("ID = ?", id).Delete(Account{}).Error
}
