package transaction

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"fts.local/utils"
)

// CreateTransaction godoc
// @Summary Create an transaction
// @Description Transfer money from one account to another
// @Tags transaction
// @Accept  json
// @Produce  json
// @Param transaction body transaction.Transaction true "JSON structure"
// @Success 200 {object} transaction.Transaction
// @Failure 400 {object} utils.HTTPError
// @Failure 401 {object} utils.HTTPError
// @Failure 403 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /auth/transaction [post]
// @Security ApiKeyAuth
func (m *Mgr) CreateTransaction(c *gin.Context) {
	var t Transaction
	var err error
	var code int
	if err = c.ShouldBindJSON(&t); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		c.Abort()
		return
	}
	// deal with some race conditions first
	// 1. check if the current login user is the sender
	if t.SrcName != utils.CurrentUser {
		utils.NewError(c, http.StatusBadRequest, errors.New("you cannot tranfer others' money"))
		c.Abort()
		return
	}
	// 2. check if the sender and receiver are same person
	if t.SrcName == t.DstName {
		utils.NewError(c, http.StatusBadRequest, errors.New("you tranfer money to yourself, are you kidding?"))
		c.Abort()
		return
	}
	// 3. check if the receiver is admin user
	if t.DstName == "admin" {
		utils.NewError(c, http.StatusBadRequest, errors.New("please do not transfer money to system admin"))
		c.Abort()
		return
	}
	if code, err = m.Add(&t); err != nil {
		utils.NewError(c, code, err)
		return
	}
	c.JSON(http.StatusOK, &t)
}

// GetTransactions godoc
// @Summary Get transactions
// @Description Search transactions based on query parameters
// @Tags transaction
// @Accept  json
// @Produce  json
// @Param name query string true "Account user name"
// @Param type query uint false "Outgoing(0)/Incoming(1)/Both(2) transactions for the account"
// @Param start query string false "The start of the query time range, format like this: 2006-01-02 15:04:05"
// @Param end query string false "The end of the query time range, format like this: 2006-01-02 15:04:05"
// @Success 200 {object} transaction.Transaction
// @Failure 400 {object} utils.HTTPError
// @Failure 401 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /auth/transaction [get]
// @Security ApiKeyAuth
func (m *Mgr) GetTransactions(c *gin.Context) {
	var (
		transactions       []Transaction
		err, err1, err2    error
		startTime, endTime time.Time
	)
	userName := c.Query("name")
	direction := c.DefaultQuery("type", "2")
	startTimeString := c.Query("start")
	endTimeString := c.Query("end")
	if startTimeString != "" {
		// There are various layout, choose better understanding string instead of RFC3339
		startTime, err1 = time.ParseInLocation("2006-01-02 15:04:05", startTimeString, time.Local)
		if err1 != nil {
			utils.NewError(c, http.StatusBadRequest, errors.New("'start' parameter is invalid"))
			fmt.Println("start parameter is invalid")
			return
		}
	}
	if endTimeString != "" {
		// There are various layout, choose better understanding string instead of RFC3339
		endTime, err2 = time.ParseInLocation("2006-01-02 15:04:05", endTimeString, time.Local)
		if err2 != nil {
			utils.NewError(c, http.StatusBadRequest, errors.New("'end' parameter is invalid"))
			fmt.Println("end parameter is invalid")
			return
		}
	}
	if userName == "" {
		utils.NewError(c, http.StatusBadRequest, errors.New("account unknown, query failed"))
		fmt.Println("User unknown, query failed!")
		return
	}

	q := Query{Name: userName, Type: direction, Start: startTime, End: endTime}

	transactions, err = m.Get(q)
	if err != nil {
		utils.NewError(c, http.StatusNotFound, err)
		c.Abort()
		fmt.Println(err.Error())
	} else {
		c.JSON(http.StatusOK, &transactions)
	}
}
