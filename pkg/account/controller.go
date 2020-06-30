package account

import "github.com/jinzhu/gorm"

// Mgr :manage account issues
type Mgr struct {
	DB *gorm.DB
}

// InitMgr :init account mgr
func InitMgr(db *gorm.DB) *Mgr {
	return &Mgr{DB: db}
}
