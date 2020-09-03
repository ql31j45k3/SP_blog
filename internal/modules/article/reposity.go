package article

func (uca *useCaseArticle) create(article Article) (uint, error) {
	result := uca.db.Create(&article)
	if result.Error != nil {
		return 0, result.Error
	}

	return article.ID, nil
}

func (uca *useCaseArticle) updateID(cond articleCond, article Article) error {
	result := uca.db.Model(Article{}).Where("id = ?", cond.ID).
		Updates(map[string]interface{}{
			"title": article.Title,
			"desc": article.Desc,
			"content": article.Content,
			"status": article.Status,
		})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (uca *useCaseArticle) getID(cond articleCond) (Article, error) {
	var article Article

	result := uca.db.First(&article, cond.ID)
	if result.Error != nil {
		return article, result.Error
	}

	return article, nil
}
