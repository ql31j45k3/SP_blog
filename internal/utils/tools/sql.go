package tools

import "gorm.io/gorm"

// SQLAppend 依照 condition 判斷是否拼湊 SQL Where 條件
func SQLAppend(db *gorm.DB, condition bool, query, columnValue interface{}) *gorm.DB {
	if !condition {
		return db
	}

	return db.Where(query, columnValue)
}
