package swagger

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	sleepInterval int              = 10
	testRounds    int              = 20
	transferTimes int              = 10
	testForever   bool             = true
	accounts      []AccountAccount = []AccountAccount{
		{Username: "Augusto", Password: "augusto", Email: "augusto@160.com", Phone: "13012345678", Balance: 100},
		{Username: "Baggio", Password: "baggio", Email: "baggio@161.com", Phone: "13112345678", Balance: 100},
		{Username: "Cassano", Password: "cassano", Email: "cassano@162.com", Phone: "13212345678", Balance: 100},
		{Username: "Djokovic", Password: "djokovic", Email: "djokovic@163.com", Phone: "13312345678", Balance: 100},
		{Username: "Elkeson", Password: "elkeson", Email: "elkeson@164.com", Phone: "13412345678", Balance: 100},
		{Username: "Federer", Password: "federer", Email: "federer@165.com", Phone: "13512345678", Balance: 100},
		{Username: "Gerrard", Password: "gerrard", Email: "gerrard@166.com", Phone: "13612345678", Balance: 100},
		{Username: "Hulk", Password: "hulk", Email: "hulk@167.com", Phone: "13712345678", Balance: 100},
		{Username: "Iniesta", Password: "iniesta", Email: "iniesta@168.com", Phone: "13812345678", Balance: 100},
		{Username: "Jimenez", Password: "jimenez", Email: "jimenez@169.com", Phone: "13912345678", Balance: 100},
		{Username: "Kane", Password: "kane", Email: "kane@170.com", Phone: "14012345678", Balance: 100},
		{Username: "Lampard", Password: "lampard", Email: "lampard@171.com", Phone: "14112345678", Balance: 100},
		{Username: "Messi", Password: "messi", Email: "messi@172.com", Phone: "14212345678", Balance: 100},
		{Username: "Neymar", Password: "neymar", Email: "neymar@173.com", Phone: "14312345678", Balance: 100},
		{Username: "Oscar", Password: "oscar", Email: "oscar@174.com", Phone: "14412345678", Balance: 100},
		{Username: "Paulinho", Password: "paulinho", Email: "paulinho@175.com", Phone: "14512345678", Balance: 100},
		{Username: "Quain", Password: "quain", Email: "quain@176.com", Phone: "14612345678", Balance: 100},
		{Username: "Ronaldo", Password: "ronaldo", Email: "ronaldo@177.com", Phone: "14712345678", Balance: 100},
		{Username: "Sarah", Password: "sarah", Email: "sarah@178.com", Phone: "14812345678", Balance: 100},
		{Username: "Tores", Password: "tores", Email: "tores@179.com", Phone: "14912345678", Balance: 100},
		{Username: "Utaka", Password: "utaka", Email: "utaka@180.com", Phone: "15012345678", Balance: 100},
		{Username: "Vieri", Password: "vieri", Email: "vieri@181.com", Phone: "15112345678", Balance: 100},
		{Username: "Wulei", Password: "wulei", Email: "wulei@182.com", Phone: "15212345678", Balance: 100},
		{Username: "Xavi", Password: "xavi", Email: "xavi@183.com", Phone: "15312345678", Balance: 100},
		{Username: "Yepes", Password: "yepes", Email: "yepes@184.com", Phone: "15412345678", Balance: 100},
		{Username: "Zidane", Password: "zidane", Email: "zidane@185.com", Phone: "15512345678", Balance: 100},
	}

	sessions []SessionInfo
)

// SessionInfo :records those necessary info during session.
type SessionInfo struct {
	id         uint //set value after /reg, also could be set in /auth/account
	name       string
	pwd        string
	token      string //set value after /auth
	registered bool   //set to true after /reg
	logined    bool   //set to true after /auth, set to false if token expires
	balance    int32
	expired    bool
	users      []string
}

func main() {
	var wg sync.WaitGroup

	sessions = make([]SessionInfo, len(accounts))

	maxUsers := len(accounts)
	wg.Add(maxUsers)
	for i := 0; i < maxUsers; i++ {
		for j := 0; j < maxUsers; j++ {
			sessions[i].users = append(sessions[i].users, accounts[i].Username)
		}
		go func(i int, w *sync.WaitGroup) {
			defer w.Done()
			ftsTrip(i)
		}(i, &wg)
	}
	wg.Wait()
}

func ftsTrip(i int) {
	myAccount := accounts[i]
	mySession := sessions[i]
	mySession.name = myAccount.Username
	mySession.pwd = myAccount.Password
	mySession.balance = myAccount.Balance
	rand.Seed(time.Now().Unix())

	if testForever {
		// Repeat below step 1-6 on and on;
		for {
			RunOneRound(&myAccount, &mySession)
			time.Sleep(time.Duration(rand.Intn(sleepInterval)) * time.Millisecond)
		}
	} else {
		// Repeat 'testTimes' times and quit
		for i := 0; i < testRounds; i++ {
			RunOneRound(&myAccount, &mySession)
			time.Sleep(time.Duration(rand.Intn(sleepInterval)) * time.Millisecond)
		}
	}
}

// RunOneRound :Test all use cases in one round.
func RunOneRound(a *AccountAccount, s *SessionInfo) {
	var err error
	ctx := context.Background()
	//step 1: Register account
	doReg(ctx, a, s)

	//step 2: Login account
	err = doLogin(ctx, s)
	if err != nil {
		log.Fatalf("%s Login unexpected failure with error: %v", a.Username, err)
		return
	}

	//step 3: List accounts
	err = doList(ctx, s)
	if err != nil && err.Error() == "token is expired" {
		s.logined = false
		doLogin(ctx, s)
	}

	//step 4: Transfer money
	for i := 0; i < transferTimes; i++ {
		err = doTransfer(ctx, s)
		if err != nil && err.Error() == "token is expired" {
			s.logined = false
			doLogin(ctx, s)
			r := rand.New(rand.NewSource(time.Now().Unix()))
			t := r.Intn(sleepInterval)
			time.Sleep(time.Duration(t) * time.Millisecond)
		}
	}

	//step 5: Check personal account balance
	err = doCheckBalance(ctx, s)
	if err != nil && err.Error() == "token is expired" {
		s.logined = false
		doLogin(ctx, s)
	}

	//step 6: Check personal transaction records
	err = doCheckTransactions(ctx, s)
	if err != nil && err.Error() == "token is expired" {
		s.logined = false
		doLogin(ctx, s)
	}
}

//Register account --> /api/v1/reg, POST
func doReg(c context.Context, a *AccountAccount, s *SessionInfo) error {
	if s.registered {
		return nil
	}
	cfg := NewConfiguration()
	client := NewAPIClient(cfg)
	rsp, err := client.AccountApi.RegPost(c, *a)
	s.registered = true
	if err == nil {
		s.id = rsp.ID
	} else {
		fmt.Println(err.Error())
	}
	return err
}

//Login account
// --> /api/v1/auth
// POST
func doLogin(c context.Context, s *SessionInfo) error {
	if s.logined {
		return nil
	}
	cfg := NewConfiguration()
	client := NewAPIClient(cfg)
	var req AccountLoginReq
	req.Name = s.name
	req.Password = s.pwd
	rsp, err := client.AccountApi.AuthPost(c, req)
	if err == nil {
		s.token = rsp.Token
		fmt.Printf("%s successfully login, his/her token is:\n%s\n\n", s.name, s.token)
		s.logined = true
	}
	return err
}

//List accounts
// --> /api/v1/auth/account
// GET
func doList(c context.Context, s *SessionInfo) error {
	cfg := NewConfiguration()
	client := NewAPIClient(cfg)
	var accounts []AccountAccount
	ctx := context.WithValue(c, ContextAPIKey, s.token)
	accounts, err := client.AccountApi.AuthAccountGet(ctx, nil)
	if err != nil {
		log.Println(err.Error())
	} else {
		fmt.Printf("Currently there are %d users registered in FTS.\n", len(accounts))
		fmt.Println("They are: ")
		for _, v := range accounts {
			if v.Username == s.name {
				// update account id again here to fix the clear issue in doReg()
				s.id = v.ID
			}
			if len(accounts) > len(s.users) {
				// new user registered, add him/her into current session user list.
				var found bool = false
				for _, u := range s.users {
					if v.Username == u {
						found = true
						break
					}
				}
				if !found {
					s.users = append(s.users, v.Username)
				}
			}
			fmt.Printf("ID: %d, Name: %s\n", v.ID, v.Username)
		}
	}
	return err
}

//Transfer money
// -- > /api/v1/transaction
// param: json, TransactionTransaction{}
// POST
// Do 'transferTimes' times,
// In each time, if fails, sleep "sleepInterval" miliseconds
func doTransfer(c context.Context, s *SessionInfo) error {
	cfg := NewConfiguration()
	client := NewAPIClient(cfg)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	index := r.Intn(len(s.users))
	name := s.users[index]
	money := int32(r.Intn(10))
	tx := TransactionTransaction{SrcName: s.name, DstName: name, Money: money}
	ctx := context.WithValue(c, ContextAPIKey, s.token)
	_, err := client.TransactionApi.AuthTransactionPost(ctx, tx)
	if err != nil {
		fmt.Printf("%s tried to transfer money, but failed just now with such error: %v\n", s.name, err)
	} else {
		fmt.Printf("%s just successfully transferred money!\n", s.name)
	}
	return err
}

//Check personal account balance--> /api/v1/auth/account?id=xxx, GET
func doCheckBalance(c context.Context, s *SessionInfo) error {
	cfg := NewConfiguration()
	client := NewAPIClient(cfg)
	ctx := context.WithValue(c, ContextAPIKey, s.token)
	account, err := client.AccountApi.AuthAccountIdGet(ctx, s.id)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		s.balance = account.Balance
		fmt.Printf("Currently %s's account balance is %d.\n", s.name, s.balance)
	}
	return err
}

//Check personal transaction records --> /api/v1/auth/transaction?name=xxx, GET
func doCheckTransactions(c context.Context, s *SessionInfo) error {
	cfg := NewConfiguration()
	client := NewAPIClient(cfg)
	var transactions []TransactionTransaction
	ctx := context.WithValue(c, ContextAPIKey, s.token)
	transactions, err := client.TransactionApi.AuthTransactionGet(ctx, s.name, nil)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Found %d transactions in which user %s involved.\n", len(transactions), s.name)
	}
	return err
}
