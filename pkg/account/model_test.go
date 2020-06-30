package account

import (
	"database/sql"
	"time"
	//"time"
	//"fts.local/utils"
	//"regexp"
	//"errors"
	"testing"

	//"github.com/prashantv/gostub"
	"github.com/DATA-DOG/go-sqlmock"
	//"github.com/go-test/deep"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	mgr  *Mgr
	//person     *model.Person
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)
	s.DB.LogMode(true)
	s.DB.SingularTable(true)
	s.mgr = InitMgr(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

/*
func (s *Suite) TestAdd() {
	var (
		id			uint32 = 1
		created_at   time.Time = time.Now()
		updated_at   time.Time = time.Now()
		deleted_at   *time.Time
		name 		string = "test_name"
		password	string
		email		string
		phone		string
		balance 	uint32 = 10
	)
	// mock: m.GetByName(name)
	s.mock.ExpectQuery(`SELECT \* FROM "account" WHERE "account"\."deleted_at" IS NULL AND \(\(Username \= \$1\)\)`).
		WithArgs(name).
		WillReturnError(errors.New("error"))
	// mock: m.DB.Create(a)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "account" \("created_at","updated_at","deleted_at","username","password","email","phone","balance"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\) RETURNING "account"\."id"`).
		WithArgs(created_at, updated_at, deleted_at, name, password, email, phone, balance).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(string(id)))
	s.mock.ExpectCommit()

	err := s.mgr.Add(&Account{Username: "test_name"})

	require.NoError(s.T(), err)
}
*/

func (s *Suite) TestGetByID() {
	var (
		id         int    = 1
		created_at time.Time = time.Now()
		updated_at time.Time = time.Now()
		deleted_at *time.Time
		name    string = "test-name"
		balance uint32 = 100
	)
	// mock: m.DB.Where("ID = ?", id).Find(&account)
	s.mock.ExpectQuery(`SELECT \* FROM "account" WHERE "account"\."deleted_at" IS NULL AND \(\(ID \= \$1\)\)`).
		WithArgs(string(id)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username","pasword", "email", "phone", "balance"}).
			AddRow(id, created_at, updated_at, deleted_at, name, "", "", "", balance))

	res, err := s.mgr.GetByID(string(id))

	require.NoError(s.T(), err)
	require.Equal(s.T(), name, res.Username)
	require.Equal(s.T(), balance, res.Balance)
}

func (s *Suite) TestGetByName() {
	var (
		id         int    = 1
		created_at time.Time = time.Now()
		updated_at time.Time = time.Now()
		deleted_at *time.Time
		name    string = "test-name"
		balance uint32 = 100
	)
	// mock: m.DB.Where("Username = ?", name).Find(&account)
	s.mock.ExpectQuery(`SELECT \* FROM "account" WHERE "account"\."deleted_at" IS NULL AND \(\(Username \= \$1\)\)`).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username","pasword", "email", "phone", "balance"}).
			AddRow(id, created_at, updated_at, deleted_at, name, "", "", "", balance))

	res, err := s.mgr.GetByName(name)

	require.NoError(s.T(), err)
	require.Equal(s.T(), name, res.Username)
	require.Equal(s.T(), balance, res.Balance)
}


func (s *Suite) TestGetIDByName() {
	var (
		id         int    = 1
		created_at time.Time = time.Now()
		updated_at time.Time = time.Now()
		deleted_at *time.Time
		name    string = "test-name"
		balance uint32 = 100
	)
	// mock: m.DB.Where("Username = ?", name).Find(&account)
	s.mock.ExpectQuery(`SELECT \* FROM "account" WHERE "account"\."deleted_at" IS NULL AND \(\(Username \= \$1\)\)`).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username","pasword", "email", "phone", "balance"}).
			AddRow(id, created_at, updated_at, deleted_at, name, "", "", "", balance))

	res, err := s.mgr.GetIDByName(name)

	require.NoError(s.T(), err)
	require.Equal(s.T(), string(id), res)

}

func (s *Suite) TestGetNameByID() {
	var (
		id         int    = 1
		created_at time.Time = time.Now()
		updated_at time.Time = time.Now()
		deleted_at *time.Time
		name    string = "test-name"
		balance uint32 = 100
	)
	// mock: m.DB.Where("ID = ?", id).Find(&account)
	s.mock.ExpectQuery(`SELECT \* FROM "account" WHERE "account"\."deleted_at" IS NULL AND \(\(ID \= \$1\)\)`).
		WithArgs(string(id)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username","pasword", "email", "phone", "balance"}).
			AddRow(id, created_at, updated_at, deleted_at, name, "", "", "", balance))

	res, err := s.mgr.GetNameByID(string(id))

	require.NoError(s.T(), err)
	require.Equal(s.T(), name, res)
}

func (s *Suite) TestList() {
	var (
		id1         int    = 1
		created_at1 time.Time = time.Now()
		updated_at1 time.Time = time.Now()
		deleted_at1 *time.Time
		name1    string = "test-name1"
		balance1 uint32 = 100

		id2         int    = 1
		created_at2 time.Time = time.Now()
		updated_at2 time.Time = time.Now()
		deleted_at2 *time.Time
		name2    string = "test-name2"
		balance2 uint32 = 200
	)
	// case 1: without line limit
	s.mock.ExpectQuery(`SELECT \* FROM "account"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username","pasword", "email", "phone", "balance"}).
			AddRow(id1, created_at1, updated_at1, deleted_at1, name1, "", "", "", balance1).
			AddRow(id2, created_at2, updated_at2, deleted_at2, name2, "", "", "", balance2))
	
	res1, err := s.mgr.List("")
	require.NoError(s.T(), err)
	require.Equal(s.T(), 2, len(res1))
/*
	// case 2: with line limit
	line := "1"
	s.mock.ExpectQuery(`SELECT \* FROM "account" WHERE "account"\."deleted_at" IS NULL LIMIT \$1`).
		WithArgs(line).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username","pasword", "email", "phone", "balance"}).
			AddRow(id1, created_at1, updated_at1, deleted_at1, name1, "", "", "", balance1))
	
	res2, err := s.mgr.List(line)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, len(res2))	
*/
}
/*
func (s *Suite) TestSave() {
	
}

func (s *Suite) TestDelete() {
	
}
*/