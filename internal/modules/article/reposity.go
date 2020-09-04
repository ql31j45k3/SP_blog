package article

import (
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
)

func (uca *useCaseArticle) create(article Article) (uint, error) {
	result := uca.db.Create(&article)
	if result.Error != nil {
		return 0, result.Error
	}

	return article.ID, nil
}

func (uca *useCaseArticle) updateID(cond *articleCond, article Article) error {
	result := uca.db.Model(Article{}).Where("`id` = ?", cond.ID).
		Updates(map[string]interface{}{
			"title":   article.Title,
			"desc":    article.Desc,
			"content": article.Content,
			"status":  article.Status,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (uca *useCaseArticle) getID(cond *articleCond) (Article, error) {
	var article Article

	result := uca.db.First(&article, cond.ID)
	if result.Error != nil {
		return article, result.Error
	}

	return article, nil
}

func (uca *useCaseArticle) get(cond *articleCond) ([]Article, error) {
	var articles []Article

	uca.db = tools.SQLAppend(uca.db, tools.IsNotZero(int(cond.ID)), "`id` = ?", cond.ID)

	uca.db = tools.SQLAppend(uca.db, tools.IsNotEmpty(cond.title), "`title` like ?", "%"+cond.title+"%")
	uca.db = tools.SQLAppend(uca.db, tools.IsNotEmpty(cond.desc), "`desc` like ?", "%"+cond.desc+"%")
	uca.db = tools.SQLAppend(uca.db, tools.IsNotEmpty(cond.content), "`content` like ?", "%"+cond.content+"%")

	uca.db = tools.SQLAppend(uca.db, tools.IsNotNegativeOne(cond.status), "`status` = ?", cond.status)

	result := uca.db.Find(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}
