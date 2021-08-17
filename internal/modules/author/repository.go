package author

import (
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"gorm.io/gorm"
)

func newRepositoryAuthor() repositoryAuthor {
	return &authorMysql{}
}

type repositoryAuthor interface {
	Create(db *gorm.DB, author authors) (uint, error)
	UpdateID(db *gorm.DB, cond *authorCond, author authors) error
	GetID(db *gorm.DB, cond *authorCond) (authors, error)
	Get(db *gorm.DB, cond *authorCond) ([]authors, error)
}

type authorMysql struct {
	_ struct{}
}

func (am *authorMysql) Create(db *gorm.DB, author authors) (uint, error) {
	result := db.Create(&author)
	if result.Error != nil {
		return 0, result.Error
	}

	return author.ID, nil
}

func (am *authorMysql) UpdateID(db *gorm.DB, cond *authorCond, author authors) error {
	result := db.Model(authors{}).Where("`id` = ?", cond.ID).
		Updates(map[string]interface{}{
			"title":   author.Title,
			"content": author.Content,
			"status":  author.Status,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (am *authorMysql) GetID(db *gorm.DB, cond *authorCond) (authors, error) {
	var author authors

	result := db.First(&author, cond.ID)
	if result.Error != nil {
		return author, result.Error
	}

	return author, nil
}

func (am *authorMysql) Get(db *gorm.DB, cond *authorCond) ([]authors, error) {
	var authors []authors

	db = tools.SQLAppend(db, tools.IsNotZero(int(cond.ID)), "`id` = ?", cond.ID)

	db = tools.SQLAppend(db, tools.IsNotEmpty(cond.title), "`title` like ?", "%"+cond.title+"%")
	db = tools.SQLAppend(db, tools.IsNotEmpty(cond.content), "`content` like ?", "%"+cond.content+"%")

	db = tools.SQLAppend(db, tools.IsNotNegativeOne(cond.status), "`status` = ?", cond.status)

	db = tools.SQLPagination(db, cond.GetRowCount(), cond.GetOffset())

	result := db.Find(&authors)
	if result.Error != nil {
		return authors, result.Error
	}

	return authors, nil
}
