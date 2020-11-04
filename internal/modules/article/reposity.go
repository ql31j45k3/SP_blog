package article

import (
	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
	"strings"
)

func (uca *useCaseArticle) create(article Article) (uint, error) {
	tx := uca.db.Begin()

	result := uca.db.Create(&article)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	for i, _ := range article.ArticleLabel {
		articleLabels := article.ArticleLabel[i]
		articleLabels.ArticlesID = article.ID

		if _, err := uca.createLabel(articleLabels); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	return article.ID, nil
}

func (uca *useCaseArticle) updateID(cond *articleCond, article Article) error {
	tx := uca.db.Begin()

	result := uca.db.Model(Article{}).Where("`id` = ?", cond.ID).
		Updates(map[string]interface{}{
			"title":   article.Title,
			"desc":    article.Desc,
			"content": article.Content,
			"status":  article.Status,
		})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if err := uca.deleteLabel(cond.ID); err != nil {
		tx.Rollback()
		return err
	}

	for i, _ := range article.ArticleLabel {
		articleLabels := article.ArticleLabel[i]
		articleLabels.ArticlesID = cond.ID

		if _, err := uca.createLabel(articleLabels); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (uca *useCaseArticle) createLabel(articleLabel ArticleLabel) (uint, error) {
	result := uca.db.Create(&articleLabel)
	if result.Error != nil {
		return 0, result.Error
	}

	return articleLabel.ID, nil
}

func (uca *useCaseArticle) deleteLabel(articlesID uint) error {
	return uca.db.Where("`articles_id` = ?", articlesID).Delete(ArticleLabel{}).Error
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

	uca.db = tools.SQLPagination(uca.db, cond.GetRowCount(), cond.GetOffset())

	result := uca.db.Find(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}

func (uca *useCaseArticle) search(cond *searchCond) ([]Article, error) {
	var sql strings.Builder
	var values []interface{}

	sql.WriteString("SELECT `articles`.`id`,")
	sql.WriteString("       `articles`.`created_at`,")
	sql.WriteString("       `articles`.`updated_at`,")
	sql.WriteString("       `articles`.`title`,")
	sql.WriteString("       `articles`.`desc`,")
	sql.WriteString("       `articles`.`content`,")
	sql.WriteString("       `articles`.`status`")
	sql.WriteString("  FROM `articles` LEFT JOIN `article_labels` ON `articles`.`id` = `article_labels`.`articles_id`")

	sql.WriteString(" WHERE `articles`.`status` = ?")
	values = append(values, tools.StatusEnable)

	keywordSqlStr := " AND (`articles`.`title` LIKE ? OR `articles`.`desc` LIKE ? OR `articles`.`content` LIKE ?)"
	keywordValues := []string{"%"+cond.keyword+"%", "%"+cond.keyword+"%", "%"+cond.keyword+"%"}
	values = tools.SQLRawAppend(tools.IsNotEmpty(cond.keyword), &sql, keywordSqlStr, values, keywordValues)

	tagsSqlStr := tools.SQLArrayToString(cond.tags, "`article_labels`.`tag`")
	values = tools.SQLRawAppend(len(cond.tags) > 0, &sql, tagsSqlStr, values, cond.tags)

	sql.WriteString(" LIMIT ? OFFSET ?")
	values = append(values, cond.GetRowCount(), cond.GetOffset())

	var articles []Article
	result := uca.db.Raw(sql.String(), values...).Scan(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}
