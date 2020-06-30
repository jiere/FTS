package transaction

import (
	"errors"
	"net/http"
	"time"
	"log"
	"fmt"
	"fts.local/pkg/account"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Transaction is the table definition of transaction records
type Transaction struct {
	gorm.Model
	SrcName string `gorm:"size:255;NOT NULL" json:"src_name"`
	DstName string `gorm:"size:255;NOT NULL" json:"dst_name"`
	Money   uint32 `gorm:"size:32;NOT NULL" json:"money"`
	Status  string `gorm:"size:12;DEFAULT 'active'" json:"status"`
}

// Query is the structure of query paremeters for transactions
type Query struct {
	Name  string
	Type  string
	Start time.Time
	End   time.Time
}

// Repository :defines interface for Transaction model operations
type Repository interface {
	Add(t *Transaction) (code int, e error)
	Get(q Query) (t []Transaction, e error)
}

// CreateTable :create transaction table
func (m *Mgr) CreateTable() {
	var err error
	if err = m.DB.AutoMigrate(&Transaction{}).Error; err != nil {
		log.Fatalf("create transaction table failure")
		panic("create transaction table failure")
	} else {
		fmt.Println("create transaction table success")
	}
}

// Add :Transfer money, if success, add a record of the transaction.
func (m *Mgr) Add(t *Transaction) (code int, e error) {
	var a1, a2 account.Account

	tx := m.DB.Begin()
	tx.Set("gorm:query_option", "FOR UPDATE")

	tx.Create(&t)
	err1 := tx.Where("Username = ?", t.SrcName).First(&a1).Error
	err2 := tx.Where("Username = ?", t.DstName).First(&a2).Error
	if err1 != nil || err2 != nil {
		tx.Rollback()
		return http.StatusBadRequest, errors.New("account not exists")
	}
	if t.Money > a1.Balance {
		tx.Rollback()
		return http.StatusBadRequest, errors.New("Transfer money is larger than sender's balance")
	}

	a1.Balance -= t.Money
	a2.Balance += t.Money
	tx.Save(&a1)
	tx.Save(&a2)
	t.Status = "closed"
	tx.Save(&t)
	tx.Commit()
	return http.StatusOK, nil
}

// Get :return transactions based on "q" parameters.
func (m *Mgr) Get(q Query) (t []Transaction, e error) {
	var err error
	switch q.Type {
	case "0": //sender
		{
			if q.Start.IsZero() && q.End.IsZero() {
				return m.GetOutgoingWithoutLimit(q.Name)
			}
			return m.GetOutgoingWithLimit(q)
		}
	case "1": //receiver
		{
			if q.Start.IsZero() && q.End.IsZero() {
				return m.GetIncomingWithoutLimit(q.Name)
			}
			return m.GetIncomingWithLimit(q)
		}
	case "2": //sender or receiver
		{
			if q.Start.IsZero() && q.End.IsZero() {
				return m.GetBothWithoutLimit(q.Name)
			}
			return m.GetBothWithLimit(q)
		}
	default:
		err = errors.New("Type parameter invalid")
		return []Transaction{}, err
	}
}

// GetOutgoingWithoutLimit :return all the "user" sent transactions.
func (m *Mgr) GetOutgoingWithoutLimit(user string) (t []Transaction, e error) {
	var transactions []Transaction
	err := m.DB.Where("src_name = ?", user).Find(&transactions).Error
	return transactions, err
}

// GetOutgoingWithLimit :returns "user" sent transactions happened between "start" and "end" interval
func (m *Mgr) GetOutgoingWithLimit(q Query) (t []Transaction, e error) {
	var transactions []Transaction
	if q.End.IsZero() {
		q.End = time.Now()
	}
	err := m.DB.Where("src_name = ?", q.Name).Where("updated_at < ? and updated_at > ?", q.End, q.Start).Find(&transactions).Error
	return transactions, err
}

// GetIncomingWithoutLimit :return all the "user" received transactions.
func (m *Mgr) GetIncomingWithoutLimit(user string) (t []Transaction, e error) {
	var transactions []Transaction
	err := m.DB.Where("dst_name = ?", user).Find(&transactions).Error
	return transactions, err
}

// GetIncomingWithLimit :returns "user" received transactions happened between "start" and "end" interval
func (m *Mgr) GetIncomingWithLimit(q Query) (t []Transaction, e error) {
	var transactions []Transaction
	if q.End.IsZero() {
		q.End = time.Now()
	}
	err := m.DB.Where("dst_name = ?", q.Name).Where("updated_at < ? and updated_at > ?", q.End, q.Start).Find(&transactions).Error
	return transactions, err
}

// GetBothWithoutLimit :return all the "user" related transactions.
func (m *Mgr) GetBothWithoutLimit(user string) (t []Transaction, e error) {
	var transactions []Transaction
	err := m.DB.Where("src_name = ? or dst_name = ?", user, user).Find(&transactions).Error
	return transactions, err
}

// GetBothWithLimit :returns "user" related transactions happened between "start" and "end" interval
func (m *Mgr) GetBothWithLimit(q Query) (t []Transaction, e error) {
	var transactions []Transaction
	if q.End.IsZero() {
		q.End = time.Now()
	}
	err := m.DB.Where("src_name = ? or dst_name = ?", q.Name, q.Name).Where("updated_at < ? and updated_at > ?", q.End, q.Start).Find(&transactions).Error
	return transactions, err
}