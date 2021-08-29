package article

import (
	"strings"

	"gorm.io/gorm"

	"github.com/ql31j45k3/SP_blog/internal/utils/tools"
)

func newRepositoryArticle() repositoryArticle {
	return &articleMysql{}
}

type repositoryArticle interface {
	Create(db *gorm.DB, article articles) (uint, error)
	UpdateID(db *gorm.DB, cond articleCond, article articles) error
	GetID(db *gorm.DB, cond articleCond) (articles, error)
	Get(db *gorm.DB, cond articleCond) ([]articles, error)
	Search(db *gorm.DB, cond searchCond) ([]articles, error)
}

type articleMysql struct {
	_ struct{}
}

func (am *articleMysql) Create(db *gorm.DB, article articles) (uint, error) {
	tx := db.Begin()

	result := db.Create(&article)
	if result.Error != nil {
		tx.Rollback()
		return 0, result.Error
	}

	for i := range article.ArticleLabel {
		articleLabels := article.ArticleLabel[i]
		//nolint:typecheck
		articleLabels.ArticlesID = article.ID

		if _, err := am.createLabel(db, articleLabels); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	//nolint:typecheck
	return article.ID, nil
}

func (am *articleMysql) UpdateID(db *gorm.DB, cond articleCond, article articles) error {
	tx := db.Begin()

	result := db.Model(articles{}).Where("`id` = ?", cond.ID).
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

	if err := am.deleteLabel(db, cond.ID); err != nil {
		tx.Rollback()
		return err
	}

	for i := range article.ArticleLabel {
		articleLabels := article.ArticleLabel[i]
		articleLabels.ArticlesID = cond.ID

		if _, err := am.createLabel(db, articleLabels); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (am *articleMysql) createLabel(db *gorm.DB, articleLabel articleLabels) (uint, error) {
	result := db.Create(&articleLabel)
	if result.Error != nil {
		return 0, result.Error
	}

	//nolint:typecheck
	return articleLabel.ID, nil
}

func (am *articleMysql) deleteLabel(db *gorm.DB, articlesID uint) error {
	return db.Where("`articles_id` = ?", articlesID).Delete(articleLabels{}).Error
}

func (am *articleMysql) GetID(db *gorm.DB, cond articleCond) (articles, error) {
	var article articles

	result := db.First(&article, cond.ID)
	if result.Error != nil {
		return article, result.Error
	}

	return article, nil
}

func (am *articleMysql) Get(db *gorm.DB, cond articleCond) ([]articles, error) {
	var articles []articles

	db = tools.SQLAppend(db, tools.IsNotZero(int(cond.ID)), "`id` = ?", cond.ID)

	db = tools.SQLAppend(db, tools.IsNotEmpty(cond.title), "`title` like ?", "%"+cond.title+"%")
	db = tools.SQLAppend(db, tools.IsNotEmpty(cond.desc), "`desc` like ?", "%"+cond.desc+"%")
	db = tools.SQLAppend(db, tools.IsNotEmpty(cond.content), "`content` like ?", "%"+cond.content+"%")

	db = tools.SQLAppend(db, tools.IsNotNegativeOne(cond.status), "`status` = ?", cond.status)

	//nolint:typecheck
	db = tools.SQLPagination(db, cond.GetRowCount(), cond.GetOffset())

	result := db.Find(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}

func (am *articleMysql) Search(db *gorm.DB, cond searchCond) ([]articles, error) {
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
	//nolint:typecheck
	values = append(values, cond.GetRowCount(), cond.GetOffset())

	var articles []articles
	result := db.Raw(sql.String(), values...).Scan(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}
