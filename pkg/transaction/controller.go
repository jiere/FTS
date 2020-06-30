package transaction

import "github.com/jinzhu/gorm"

// Mgr :manage transaction issues
type Mgr struct {
	DB *gorm.DB
}

// InitMgr :init transaction mgr 
func InitMgr(db *gorm.DB) *Mgr {
	return &Mgr{DB: db}
}
