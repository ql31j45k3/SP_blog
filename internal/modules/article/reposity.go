package article

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

	if cond.ID != 0 {
		uca.db = uca.db.Where("`id` = ?", cond.ID)
	}

	if cond.title != "" {
		uca.db = uca.db.Where("`title` like ?", "%"+cond.title+"%")
	}

	if cond.desc != "" {
		uca.db = uca.db.Where("`desc` like ?", "%"+cond.desc+"%")
	}

	if cond.content != "" {
		uca.db = uca.db.Where("`content` like ?", "%"+cond.content+"%")
	}

	if cond.status == 0 || cond.status == 1 {
		uca.db = uca.db.Where("`status` = ?", cond.status)
	}

	result := uca.db.Find(&articles)
	if result.Error != nil {
		return articles, result.Error
	}

	return articles, nil
}
