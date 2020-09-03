package article

const (
	stateEnable = 1
)

func (uca *useCaseArticle) post(article Article) (uint, error) {
	article.State = stateEnable

	result := uca.db.Create(&article)
	if result.Error != nil {
		return 0, result.Error
	}

	return article.ID, nil
}

func (uca *useCaseArticle) getID(cond articleCond) (Article, error) {
	var article Article

	result := uca.db.First(&article, cond.ID)
	if result.Error != nil {
		return article, result.Error
	}

	return article, nil
}
