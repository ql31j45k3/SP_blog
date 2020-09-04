package tools

import "gorm.io/gorm"

func SQLAppend(db *gorm.DB, condition bool, query, columnValue interface{}) *gorm.DB {
	if !condition {
		return db
	}

	return db.Where(query, columnValue)
}
