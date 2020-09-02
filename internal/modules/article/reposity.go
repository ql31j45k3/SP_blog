package article

func (uca *useCaseArticle) getID(cond articleCond) (Article, error) {
	var article Article

	result := uca.db.First(&article, cond.ID)
	if result.Error != nil {
		return article, result.Error
	}

	return article, nil
}
