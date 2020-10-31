package article

import (
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"strings"
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

func (uca *useCaseArticle) search(cond *searchCond) ([]Article, error) {
	var sql strings.Builder

	sql.WriteString("SELECT `articles`.`id`,")
	sql.WriteString("       `articles`.`created_at`,")
	sql.WriteString("       `articles`.`updated_at`,")
	sql.WriteString("       `articles`.`title`,")
	sql.WriteString("       `articles`.`desc`,")
	sql.WriteString("       `articles`.`content`,")
	sql.WriteString("       `articles`.`status`")
	sql.WriteString("  FROM `articles`")
	sql.WriteString(" WHERE `articles`.`status` = ?")
	sql.WriteString("   AND (`articles`.`title` LIKE ? OR `articles`.`desc` LIKE ? OR `articles`.`content` LIKE ?)")

	var articles []Article

	result := uca.db.Raw(sql.String(),tools.StatusEnable,
		"%"+cond.keyword+"%", "%"+cond.keyword+"%", "%"+cond.keyword+"%").Scan(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}
