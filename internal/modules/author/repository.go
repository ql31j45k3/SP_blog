package author

import "github.com/ql31j45k3/SP_blog/internal/utils/tools"

func (uca *useCaseAuthor) create(author Author) (uint, error) {
	result := uca.db.Create(&author)
	if result.Error != nil {
		return 0, result.Error
	}

	return author.ID, nil
}

func (uca *useCaseAuthor) updateID(cond *authorCond, author Author) error {
	result := uca.db.Model(Author{}).Where("`id` = ?", cond.ID).
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

func (uca *useCaseAuthor) getID(cond *authorCond) (Author, error) {
	var author Author

	result := uca.db.First(&author, cond.ID)
	if result.Error != nil {
		return author, result.Error
	}

	return author, nil
}

func (uca *useCaseAuthor) get(cond *authorCond) ([]Author, error) {
	var authors []Author

	uca.db = tools.SQLAppend(uca.db, tools.IsNotZero(int(cond.ID)), "`id` = ?", cond.ID)

	uca.db = tools.SQLAppend(uca.db, tools.IsNotEmpty(cond.title), "`title` like ?", "%"+cond.title+"%")
	uca.db = tools.SQLAppend(uca.db, tools.IsNotEmpty(cond.content), "`content` like ?", "%"+cond.content+"%")

	uca.db = tools.SQLAppend(uca.db, tools.IsNotNegativeOne(cond.status), "`status` = ?", cond.status)

	uca.db = tools.SQLPagination(uca.db, cond.GetRowCount(), cond.GetOffset())

	result := uca.db.Find(&authors)
	if result.Error != nil {
		return authors, result.Error
	}

	return authors, nil
}
