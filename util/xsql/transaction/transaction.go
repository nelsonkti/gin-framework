package transaction

import (
	"gorm.io/gorm"
)

type Transaction struct {
	Tx map[string]*gorm.DB
}

// Commit 提交事务
func (r *Transaction) Commit() error {
	if r.Tx == nil {
		return nil
	}

	for _, db := range r.Tx {
		err := db.Commit().Error
		if err != nil {
			return err
		}
	}
	return nil
}

// Rollback 回滚事务
func (r *Transaction) Rollback() error {
	if r.Tx == nil {
		return nil
	}

	for _, db := range r.Tx {
		err := db.Rollback().Error
		if err != nil {
			return err
		}
	}

	return nil
}
