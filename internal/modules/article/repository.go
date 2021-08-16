package article

import (
	"strings"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
)

func (a *article) create(article articles) (uint, error) {
	tx := a.db.Begin()

	result := a.db.Create(&article)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	for i, _ := range article.ArticleLabel {
		articleLabels := article.ArticleLabel[i]
		articleLabels.ArticlesID = article.ID

		if _, err := a.createLabel(articleLabels); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	return article.ID, nil
}

func (a *article) updateID(cond *articleCond, article articles) error {
	tx := a.db.Begin()

	result := a.db.Model(articles{}).Where("`id` = ?", cond.ID).
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

	if err := a.deleteLabel(cond.ID); err != nil {
		tx.Rollback()
		return err
	}

	for i, _ := range article.ArticleLabel {
		articleLabels := article.ArticleLabel[i]
		articleLabels.ArticlesID = cond.ID

		if _, err := a.createLabel(articleLabels); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (a *article) createLabel(articleLabel ArticleLabel) (uint, error) {
	result := a.db.Create(&articleLabel)
	if result.Error != nil {
		return 0, result.Error
	}

	return articleLabel.ID, nil
}

func (a *article) deleteLabel(articlesID uint) error {
	return a.db.Where("`articles_id` = ?", articlesID).Delete(ArticleLabel{}).Error
}

func (a *article) getID(cond *articleCond) (articles, error) {
	var article articles

	result := a.db.First(&article, cond.ID)
	if result.Error != nil {
		return article, result.Error
	}

	return article, nil
}

func (a *article) get(cond *articleCond) ([]articles, error) {
	var articles []articles

	a.db = tools.SQLAppend(a.db, tools.IsNotZero(int(cond.ID)), "`id` = ?", cond.ID)

	a.db = tools.SQLAppend(a.db, tools.IsNotEmpty(cond.title), "`title` like ?", "%"+cond.title+"%")
	a.db = tools.SQLAppend(a.db, tools.IsNotEmpty(cond.desc), "`desc` like ?", "%"+cond.desc+"%")
	a.db = tools.SQLAppend(a.db, tools.IsNotEmpty(cond.content), "`content` like ?", "%"+cond.content+"%")

	a.db = tools.SQLAppend(a.db, tools.IsNotNegativeOne(cond.status), "`status` = ?", cond.status)

	a.db = tools.SQLPagination(a.db, cond.GetRowCount(), cond.GetOffset())

	result := a.db.Find(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}

func (a *article) search(cond *searchCond) ([]articles, error) {
	var sql strings.Builder
	var values []interface{}

	sql.WriteString("SELECT DISTINCT")
	sql.WriteString("       `articles`.`id`,")
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
	keywordValues := []string{"%" + cond.keyword + "%", "%" + cond.keyword + "%", "%" + cond.keyword + "%"}
	values = tools.SQLRawAppend(tools.IsNotEmpty(cond.keyword), &sql, keywordSqlStr, values, keywordValues)

	tagsSqlStr := tools.SQLArrayToString(cond.tags, "`article_labels`.`tag`")
	values = tools.SQLRawAppend(len(cond.tags) > 0, &sql, tagsSqlStr, values, cond.tags)

	sql.WriteString(" LIMIT ? OFFSET ?")
	values = append(values, cond.GetRowCount(), cond.GetOffset())

	var articles []articles
	result := a.db.Raw(sql.String(), values...).Scan(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}
