package author

import "github.com/ql31j45k3/SP_blog/internal/utils/tools"

func (a *author) create(author Author) (uint, error) {
	result := a.db.Create(&author)
	if result.Error != nil {
		return 0, result.Error
	}

	return author.ID, nil
}

func (a *author) updateID(cond *authorCond, author Author) error {
	result := a.db.Model(Author{}).Where("`id` = ?", cond.ID).
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

func (a *author) getID(cond *authorCond) (Author, error) {
	var author Author

	result := a.db.First(&author, cond.ID)
	if result.Error != nil {
		return author, result.Error
	}

	return author, nil
}

func (a *author) get(cond *authorCond) ([]Author, error) {
	var authors []Author

	a.db = tools.SQLAppend(a.db, tools.IsNotZero(int(cond.ID)), "`id` = ?", cond.ID)

	a.db = tools.SQLAppend(a.db, tools.IsNotEmpty(cond.title), "`title` like ?", "%"+cond.title+"%")
	a.db = tools.SQLAppend(a.db, tools.IsNotEmpty(cond.content), "`content` like ?", "%"+cond.content+"%")

	a.db = tools.SQLAppend(a.db, tools.IsNotNegativeOne(cond.status), "`status` = ?", cond.status)

	a.db = tools.SQLPagination(a.db, cond.GetRowCount(), cond.GetOffset())

	result := a.db.Find(&authors)
	if result.Error != nil {
		return authors, result.Error
	}

	return authors, nil
}
