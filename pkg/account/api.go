package account

import (
	"errors"
	"log"
	"net/http"
	"time"

	"fts.local/utils"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CreateAccount godoc
// @Summary Create an account
// @Description Create an account by JSON format parameters
// @Tags account
// @Accept  json
// @Produce  json
// @Param account body account.Account true "account.Account struct JSON"
// @Success 200 {object} account.Account
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /reg [post]
func (m *Mgr) CreateAccount(c *gin.Context) {
	var account Account
	var err error
	if err = c.ShouldBindJSON(&account); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		c.Abort()
		return
	}

	if err = m.Add(&account); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, &account)
	}
}

// Login godoc
// @Summary Login an account
// @Description Login using username and password
// @Tags account
// @Accept  json
// @Produce  json
// @Param req body account.LoginReq true "account.LoginReq struct JSON"
// @Success 200 {object} account.LoginRsp
// @Failure 400 {object} utils.HTTPError
// @Failure 401 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /auth [post]
func (m *Mgr) Login(c *gin.Context) {
	var (
		req LoginReq
		acc Account
		err error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.NewError(c, http.StatusBadRequest, errors.New("Login paremeter parse failure"))
		c.Abort()
		return
	}
	acc, err = m.GetByName(req.Username)
	if err != nil {
		utils.NewError(c, http.StatusUnauthorized, errors.New("username incorrect"))
		c.Abort()
	} else if req.Password != acc.Password {
		utils.NewError(c, http.StatusUnauthorized, errors.New("password incorrect"))
		c.Abort()
	} else {
		m.GenerateToken(c, acc)
	}
}

// GenerateToken :generate JWT token
func (m *Mgr) GenerateToken(c *gin.Context, account Account) {
	j := utils.NewJWT()
	claims := utils.CustomClaims{
		account.Username,
		account.Email,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),
			ExpiresAt: int64(time.Now().Unix() + 3600),
			Issuer:    "mobilewallet.com",
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		return
	}

	log.Println(token)

	data := LoginRsp{
		Username: account.Username,
		Token:    token,
	}
	c.JSON(http.StatusOK, &data)
}

func (m *Mgr) hidePrivilegeInfo(a *Account) {
	a.Password = "***"
	a.Email = "***"
	a.Phone = "***"
	a.Balance = 0
}

// GetAccount godoc
// @Summary Show an account
// @Description Get by account ID
// @Tags account
// @Accept  json
// @Produce  json
// @Param  id path uint32 true "Account ID"
// @Success 200 {object} account.Account
// @Failure 401 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /auth/account/{id} [get]
// @Security ApiKeyAuth
func (m *Mgr) GetAccount(c *gin.Context) {
	accID := c.Param("id")
	account, err := m.GetByID(accID)
	if err != nil {
		utils.NewError(c, http.StatusNotFound, err)
		c.Abort()
	} else {
		if account.Username != utils.CurrentUser {
			m.hidePrivilegeInfo(&account)
		}
		c.JSON(http.StatusOK, &account)
	}
}

// ListAccounts godoc
// @Summary List all accounts
// @Description List all the accounts
// @Tags account
// @Accept  json
// @Produce  json
// @Param line query uint false "List result number no more than the value of 'line'"
// @Success 200 {object} account.Account
// @Failure 401 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /auth/account [get]
// @Security ApiKeyAuth
func (m *Mgr) ListAccounts(c *gin.Context) {
	line := c.Query("line")
	accounts, err := m.List(line)
	if err != nil {
		utils.NewError(c, http.StatusNotFound, err)
		c.Abort()
	} else {
		for i := 0; i < len(accounts); i++ {
			if accounts[i].Username != utils.CurrentUser {
				m.hidePrivilegeInfo(&accounts[i])
			}
		}
		c.JSON(http.StatusOK, &accounts)
	}
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update by account ID
// @Tags account
// @Accept  json
// @Produce  json
// @Param  id path uint32 true "Account ID"
// @Param fields body account.Account true "Update fields/values as JSON"
// @Success 200 {object} account.Account
// @Failure 400 {object} utils.HTTPError
// @Failure 401 {object} utils.HTTPError
// @Failure 403 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /auth/account/{id} [put]
// @Security ApiKeyAuth
func (m *Mgr) UpdateAccount(c *gin.Context) {
	var account Account
	var err error
	var targetUser string
	if err = c.ShouldBindJSON(&account); err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		c.Abort()
		return
	}
	accID := c.Param("id")
	targetUser, err = m.GetNameByID(accID)
	if err != nil {
		utils.NewError(c, http.StatusInternalServerError, err)
		c.Abort()
		return
	} else if targetUser != utils.CurrentUser {
		utils.NewError(c, http.StatusForbidden, errors.New("you cannot modify others' account"))
		c.Abort()
		return
	}
	err = m.Save(accID, &account)
	if err != nil {
		utils.NewError(c, http.StatusBadRequest, err)
		c.Abort()
	} else {
		c.JSON(http.StatusOK, &account)
	}
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete by account ID
// @Tags account
// @Accept  json
// @Produce  json
// @Param  id path uint32 true "Account ID"
// @Success 204 {object} account.Account
// @Failure 401 {object} utils.HTTPError
// @Failure 403 {object} utils.HTTPError
// @Failure 404 {object} utils.HTTPError
// @Router /auth/account/{id} [delete]
// @Security ApiKeyAuth
func (m *Mgr) DeleteAccount(c *gin.Context) {
	if utils.CurrentUser != "admin" {
		utils.NewError(c, http.StatusForbidden, errors.New("only admin account could do this operation"))
		c.Abort()
		return
	}
	id := c.Param("id")
	err := m.Delete(id)
	if err != nil {
		utils.NewError(c, http.StatusNotFound, err)
		c.Abort()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "The account has been deleted",
		})
	}
}
