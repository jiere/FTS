package transaction

import (
	"database/sql"
	//"errors"
	//"regexp"
	"testing"
	//"time"
	//"reflect"

	//. "github.com/agiledragon/gomonkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GormMockTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	mgr *Mgr
}

func (s *GormMockTestSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("mysql", db)
	require.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.mgr = InitMgr(s.DB)
}

func (s *GormMockTestSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(GormMockTestSuite))
}

func (s *GormMockTestSuite) TestAddFailureByInvalidSender () {
/*	 
	var (
		id 			uint32
		created_at	time.Time = time.Now()
		updated_at	time.Time = time.Now()
		deleted_at	*time.Time
		src_name	string = "a"
		dst_name	string = "b"
		money		uint32 = 20
		status		string
	)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "transaction" \("created_at", "updated_at", "deleted_at", "src_name", "dst_name", "money", "status"\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) RETURNING "transaction"\."id"`).
	    WithArgs(created_at, updated_at, deleted_at, src_name, dst_name, money, status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
		AddRow(id))
	s.mock.ExpectQuery(`SELECT \* FROM "account" WHERE \(Username \= \$1\) WITH UPDATE`).
		WithArgs(src_name).
		WillReturnError(errors.New("error"))
	s.mock.ExpectRollback()

	tx := Transaction {SrcName: "a", DstName: "b", Money: 20}
	code, err := s.mgr.Add(&tx);
	require.Error(s.T(), err)
	require.Equal(s.T(), http.StatusBadRequest, code)
	require.Equal(s.T(), "account not exists", err.Error())
*/	
}

func (s *GormMockTestSuite) TestAddFailureByInvalidReceiver () {

}

func (s *GormMockTestSuite) TestAddFailureByInvalidBalance () {

}

func (s *GormMockTestSuite) TestAddSuccess () {

}

func (s *GormMockTestSuite) TestGet () {

}


