package transaction

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"fts.local/utils"
)

func TestCreateTransaction(t *testing.T) {
	testData := []Transaction {
		{SrcName: "A", DstName: "B", Money: 20},
		{SrcName: "B", DstName: "C", Money: 50},
		{SrcName: "A", DstName: "A", Money: 100},
		{SrcName: "A", DstName: "admin", Money: 200},
		{SrcName: "A", DstName: "C", Money: 200},
	}

	cases := []struct {
		Name	string
		Data    Transaction
	}{
		// case 1: success
		{"success", testData[0]},
		// case 2: currentuser != sender
		{"wrong sender", testData[1]},
		// case 3: sender == receiver
		{"meaningless transfer", testData[2]},
		// case 4: transfer to "admin"
		{"wrong receiver", testData[3]},
		// case 5: DB ops return error
		{"model API returns failure", testData[4]},
	}
	// case set up
	currentUserBak := utils.CurrentUser
	utils.CurrentUser = "A"
	router := utils.Router()
	if router == nil {
		utils.UTInit()
		router = utils.Router()
	}
	m := InitMgr(utils.DB())
	router.POST("/postTx", m.CreateTransaction)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			param := make(map[string]interface{})
			param["src_name"] = c.Data.SrcName
			param["dst_name"] = c.Data.DstName
			param["money"] = c.Data.Money

			switch i {
				case 0: 
					patches := ApplyMethod(reflect.TypeOf(m), "Add", func(_ *Mgr, _ *Transaction) (int, error) {
                		return http.StatusOK, nil
            		})
					defer patches.Reset()
					_, code := utils.PostJSON("/postTx", param, utils.Router())
					assert.Equal(t, http.StatusOK, code)
				case 1, 2, 3:
					patches := ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
						e := utils.HTTPError{
							Code:    status,
							Message: err.Error(),
						}
						ctx.JSON(status, e)
					})
					defer patches.Reset()					
					_, code := utils.PostJSON("/postTx", param, utils.Router())
					assert.Equal(t, http.StatusBadRequest, code)
				case 4: 
					patches := ApplyMethod(reflect.TypeOf(m), "Add", func(_ *Mgr, _ *Transaction) (int, error) {
        	        	return http.StatusBadRequest, errors.New("bad request")
					})
					patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
						e := utils.HTTPError{
							Code:    status,
							Message: err.Error(),
						}
						ctx.JSON(status, e)
					})					
					defer patches.Reset()
					_, code := utils.PostJSON("/postTx", param, utils.Router())
					assert.Equal(t, http.StatusBadRequest, code)
			}
		})
	}
	// case tear down
	utils.CurrentUser = currentUserBak
}

type RawQuery struct {
	Name  string
	Type  string
	Start string
	End   string
}

func TestGetTransactions(t *testing.T) {
	testData := []RawQuery {
		{Name: "A", Type: "", Start: "2012/05/18", End: ""},
		{Name: "B", Type: "", Start: "", End: "2012/05/18"},
		{Name: "", Type: "2", Start: "", End: ""},
		{Name: "A", Type: "1", Start: "", End: "2020-06-29%2015%3A00%3A00"},
		{Name: "A", Type: "0", Start: "2020-01-01%2000%3A00%3A00", End: ""},
	}

	cases := []struct {
		Name	string
		Data    RawQuery
	}{
		// case 1: start time format invalid
		{"start time parameter invalid", testData[0]},
		// case 2: end time format invalid
		{"end time parameter invalid", testData[1]},
		// case 3: name unknown
		{"name parameter missing", testData[2]},
		// case 4: success
		{"success", testData[3]},
		// case 5: not found
		{"not found any record", testData[4]},
	}
	// case set up	
	router := utils.Router()
	if router == nil {
		utils.UTInit()
		router = utils.Router()
	}
	m := InitMgr(utils.DB())
	router.GET("/getTx", m.GetTransactions)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			param := make(map[string]string)
			if c.Data.Name != "" {
				param["name"] = c.Data.Name
			}
			if c.Data.Type != "" {
				param["type"] = c.Data.Type
			}
			if c.Data.Start != "" {
				param["start"] = c.Data.Start
			}
			if c.Data.End != "" {
				param["end"] = c.Data.End
			}			
			uri := "/getTx" + utils.ParseToStr(param)
			switch i {
				case 0, 1, 2: 
				patches := ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()			
					_, code := utils.Get(uri, utils.Router())
					assert.Equal(t, http.StatusBadRequest, code)
				case 3: 
					patches := ApplyMethod(reflect.TypeOf(m), "Get", func(_ *Mgr, _ Query) ([]Transaction, error) {
						return []Transaction{}, nil
					})
					defer patches.Reset()
					_, code := utils.Get(uri, utils.Router())
					assert.Equal(t, http.StatusOK, code)
				case 4:
					patches := ApplyMethod(reflect.TypeOf(m), "Get", func(_ *Mgr, _ Query) ([]Transaction, error) {
						return []Transaction{}, errors.New("not found")
					})
					patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
						e := utils.HTTPError{
							Code:    status,
							Message: err.Error(),
						}
						ctx.JSON(status, e)
					})
					defer patches.Reset()
					_, code := utils.Get(uri, utils.Router())
					assert.Equal(t, http.StatusNotFound, code)
			}
		})
	}
	// case tear down, for this case, do nothing.
}
