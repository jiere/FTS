package account

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"fts.local/utils"
)

func TestCreateAccount(t *testing.T) {
	testData := []Account{
		{Username: "A", Password: "aa", Email: "", Phone: "", Balance: 100},
		{Username: "B", Password: "bb", Email: "", Phone: "", Balance: 100},
	}

	cases := []struct {
		Name string
		Data Account
	}{
		// case 1: success
		{"success", testData[0]},
		// case 2: account aleady exists
		{"account exists", testData[1]},
	}
	// case set up
	router := utils.Router()
	if router == nil {
		utils.UTInit()
		router = utils.Router()
	}
	m := InitMgr(utils.DB())
	router.POST("/postAcc", m.CreateAccount)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			param := make(map[string]interface{})
			param["username"] = c.Data.Username
			param["password"] = c.Data.Password
			param["email"] = c.Data.Email
			param["phone"] = c.Data.Phone
			param["balance"] = c.Data.Balance

			switch i {
			case 0:
				patches := ApplyMethod(reflect.TypeOf(m), "Add", func(_ *Mgr, _ *Account) error {
					return nil
				})
				defer patches.Reset()
				_, code := utils.PostJSON("/postAcc", param, utils.Router())
				assert.Equal(t, http.StatusOK, code)

			case 1:
				patches := ApplyMethod(reflect.TypeOf(m), "Add", func(_ *Mgr, _ *Account) error {
					return errors.New("bad request")
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				_, code := utils.PostJSON("/postAcc", param, utils.Router())
				assert.Equal(t, http.StatusBadRequest, code)
			}
		})
	}
	// case tear down, do nothing
}

func TestLogin(t *testing.T) {
	testData := []LoginReq{
		{Username: "A", Password: "aa"},
		{Username: "B", Password: "bb"},
		{Username: "C", Password: "cc"},
	}

	cases := []struct {
		Name string
		Data LoginReq
	}{
		// case 1: success
		{"success", testData[0]},
		// case 2: invalid username
		{"invalid username", testData[1]},
		// case 3: invalid password
		{"invalid password", testData[2]},
	}
	// case set up
	router := utils.Router()
	if router == nil {
		utils.UTInit()
		router = utils.Router()
	}
	m := InitMgr(utils.DB())
	router.POST("/postLogin", m.Login)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			param := make(map[string]interface{})
			param["username"] = c.Data.Username
			param["password"] = c.Data.Password

			switch i {
			case 0:
				patches := ApplyMethod(reflect.TypeOf(m), "GetByName", func(_ *Mgr, _ string) (Account, error) {
					return Account{Username: "A", Password: "aa"}, nil
				})
				patches.ApplyMethod(reflect.TypeOf(m), "GenerateToken", func(_ *Mgr, _ *gin.Context, _ Account) {
					return
				})
				defer patches.Reset()
				_, code := utils.PostJSON("/postLogin", param, utils.Router())
				assert.Equal(t, http.StatusOK, code)

			case 1:
				patches := ApplyMethod(reflect.TypeOf(m), "GetByName", func(_ *Mgr, _ string) (Account, error) {
					return Account{}, errors.New("bad request")
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				_, code := utils.PostJSON("/postLogin", param, utils.Router())
				assert.Equal(t, http.StatusUnauthorized, code)
			case 2:
				patches := ApplyMethod(reflect.TypeOf(m), "GetByName", func(_ *Mgr, _ string) (Account, error) {
					return Account{Username: "C", Password: "ccc"}, nil
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				_, code := utils.PostJSON("/postLogin", param, utils.Router())
				assert.Equal(t, http.StatusUnauthorized, code)
			}
		})
	}
	// case tear down, do nothing
}
func TestGetAccount(t *testing.T) {
	testData := []Account{
		{Username: "A", Password: "aa", Email: "aaa@abc.com", Phone: "13933338888", Balance: 100},
		{Username: "B", Password: "bb", Email: "bbb@abc.com", Phone: "13712345678", Balance: 200},
	}

	cases := []struct {
		Name string
		Data Account
	}{
		// case 1: get personal account info
		{"search personal account", testData[0]},
		// case 2: account aleady exists
		{"search other's account", testData[1]},
		// case 3: account not exists
		{"search invalid account", Account{}},
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
	router.GET("/getAcc", m.GetAccount)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			switch i {
			case 0:
				patches := ApplyMethod(reflect.TypeOf(m), "GetByID", func(_ *Mgr, _ string) (Account, error) {
					return testData[0], nil
				})
				defer patches.Reset()
				uri := "/getAcc?id=1"
				body, code := utils.Get(uri, utils.Router())
				var v Account
				json.Unmarshal(body, &v)
				assert.Equal(t, http.StatusOK, code)
				assert.Equal(t, testData[0], v)
			case 1:
				patches := ApplyMethod(reflect.TypeOf(m), "GetByID", func(_ *Mgr, _ string) (Account, error) {
					return testData[1], nil
				})
				defer patches.Reset()
				uri := "/getAcc?id=1"
				body, code := utils.Get(uri, utils.Router())
				var v Account
				json.Unmarshal(body, &v)
				assert.Equal(t, http.StatusOK, code)
				assert.Equal(t, testData[1].Username, v.Username)
				assert.EqualValues(t, 0, v.Balance)
				assert.Equal(t, "***", v.Password)
			case 2:
				patches := ApplyMethod(reflect.TypeOf(m), "GetByID", func(_ *Mgr, _ string) (Account, error) {
					return Account{}, errors.New("not found")
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				uri := "/getAcc?id=1"
				_, code := utils.Get(uri, utils.Router())
				assert.Equal(t, http.StatusNotFound, code)
			}
		})
	}
	// case tear down, recover the real current user
	utils.CurrentUser = currentUserBak
}

func TestListAccounts(t *testing.T) {
	testData := []Account{
		{Username: "A", Password: "aa", Email: "aaa@abc.com", Phone: "13933338888", Balance: 100},
		{Username: "B", Password: "bb", Email: "bbb@abc.com", Phone: "13712345678", Balance: 200},
	}

	cases := []struct {
		Name string
		Line int
	}{
		// case 1: list accounts w/o limit
		{"list all accounts", 0},
		// case 2: list accounts w/ limit by line
		{"list at most 'line' accounts", 1},
		// case 3: fail case.
		{"accounts not found", 0},
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
	router.GET("/List", m.ListAccounts)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			switch i {
			case 0:
				patches := ApplyMethod(reflect.TypeOf(m), "List", func(_ *Mgr, _ string) ([]Account, error) {
					return testData, nil
				})
				defer patches.Reset()
				uri := "/List"
				body, code := utils.Get(uri, utils.Router())
				var v []Account
				json.Unmarshal(body, &v)
				assert.Equal(t, http.StatusOK, code)
				assert.Equal(t, 2, len(v))
				assert.Equal(t, testData[0].Username, v[0].Username)
				assert.Equal(t, "***", v[1].Password)
				
			case 1:
				patches := ApplyMethod(reflect.TypeOf(m), "List", func(_ *Mgr, _ string) ([]Account, error) {
					return []Account{testData[0]}, nil
				})
				defer patches.Reset()
				uri := "/List?line=1"
				body, code := utils.Get(uri, utils.Router())
				var v []Account
				json.Unmarshal(body, &v)
				assert.Equal(t, http.StatusOK, code)
				assert.Equal(t, 1, len(v))
			case 2:
				patches := ApplyMethod(reflect.TypeOf(m), "List", func(_ *Mgr, _ string) ([]Account, error) {
					return []Account{}, errors.New("not found")
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				uri := "/List"
				_, code := utils.Get(uri, utils.Router())
				assert.Equal(t, http.StatusNotFound, code)
			}
		})
	}
	// case tear down, recover the real current user
	utils.CurrentUser = currentUserBak
}

func TestUpdateAccount(t *testing.T) {
	param := make(map[string]interface{})
	param["password"] = "abc"
	param["email"] = "aaa@bbb.com"
	
	cases := []struct {
		Name 		string
		Curruser 	string
		Updateuser  string
	}{
		// case 1: GetNameByID() return error
		{"internal error by GetNameByID()", "A", "A"},
		// case 2: update other's account
		{"update other's account", "A", "B"},
		// case 3: Save() return error
		{"bad request leads to Save() failure", "A", "A"},
		// case 4: success case
		{"success case", "A", "A"},
	}
	// case set up
	currentUserBak := utils.CurrentUser
	router := utils.Router()
	if router == nil {
		utils.UTInit()
		router = utils.Router()
	}
	m := InitMgr(utils.DB())
	router.PUT("/Put", m.UpdateAccount)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			utils.CurrentUser = c.Curruser
			switch i {
			case 0:
				patches := ApplyMethod(reflect.TypeOf(m), "GetNameByID", func(_ *Mgr, _ string) (string, error) {
					return c.Updateuser, errors.New("error")
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				uri := "/Put?id=2"
				_, code := utils.PutJSON(uri, param, utils.Router())
				assert.Equal(t, http.StatusInternalServerError, code)
			case 1:
				patches := ApplyMethod(reflect.TypeOf(m), "GetNameByID", func(_ *Mgr, _ string) (string, error) {
					return c.Updateuser, nil
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})				
				defer patches.Reset()
				uri := "/Put?id=2"
				_, code := utils.PutJSON(uri, param, utils.Router())
				assert.Equal(t, http.StatusForbidden, code)
			case 2:
				patches := ApplyMethod(reflect.TypeOf(m), "GetNameByID", func(_ *Mgr, _ string) (string, error) {
					return c.Updateuser, nil
				})
				patches.ApplyMethod(reflect.TypeOf(m), "Save", func(_ *Mgr, _ string, _ *Account) error {
					return errors.New("error")
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})				
				defer patches.Reset()
				uri := "/Put?id=2"
				_, code := utils.PutJSON(uri, param, utils.Router())
				assert.Equal(t, http.StatusBadRequest, code)
			case 3:
				patches := ApplyMethod(reflect.TypeOf(m), "GetNameByID", func(_ *Mgr, _ string) (string, error) {
					return c.Updateuser, nil
				})
				patches.ApplyMethod(reflect.TypeOf(m), "Save", func(_ *Mgr, _ string, _ *Account) error {
					return nil
				})				
				defer patches.Reset()
				uri := "/Put?id=2"
				_, code := utils.PutJSON(uri, param, utils.Router())
				assert.Equal(t, http.StatusOK, code)
			}
		})
	}
	// case tear down, recover the real current user
	utils.CurrentUser = currentUserBak
}

func TestDeleteAccount(t *testing.T) {
	cases := []struct {
		Name string
		User string
	}{
		// case 1: non-admin user cannot do delete operation
		{"current user is not admin", "A"},
		// case 2: success case
		{"delete success", "admin"},
		// case 3: delete account not exists
		{"delete invalid account", "admin"},
	}
	// case set up
	currentUserBak := utils.CurrentUser
	router := utils.Router()
	if router == nil {
		utils.UTInit()
		router = utils.Router()
	}
	m := InitMgr(utils.DB())
	router.DELETE("/Delete", m.DeleteAccount)
	// start cases
	for i, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			utils.CurrentUser = c.User
			switch i {
			case 0:
				patches := ApplyMethod(reflect.TypeOf(m), "Delete", func(_ *Mgr, _ string) error {
					return nil
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				uri := "/Delete?id=2"
				_, code := utils.Delete(uri, utils.Router())
				assert.Equal(t, http.StatusForbidden, code)
			case 1:
				patches := ApplyMethod(reflect.TypeOf(m), "Delete", func(_ *Mgr, _ string) error {
					return nil
				})
				defer patches.Reset()
				uri := "/Delete?id=2"
				_, code := utils.Delete(uri, utils.Router())
				assert.Equal(t, http.StatusOK, code)
			case 2:
				patches := ApplyMethod(reflect.TypeOf(m), "Delete", func(_ *Mgr, _ string) error {
					return errors.New("not found")
				})
				patches.ApplyFunc(utils.NewError, func(ctx *gin.Context, status int, err error) {
					e := utils.HTTPError{
						Code:    status,
						Message: err.Error(),
					}
					ctx.JSON(status, e)
				})
				defer patches.Reset()
				uri := "/Delete?id=2"
				_, code := utils.Delete(uri, utils.Router())
				assert.Equal(t, http.StatusNotFound, code)
			}
		})
	}
	// case tear down, recover the real current user
	utils.CurrentUser = currentUserBak
}
