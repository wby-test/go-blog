package dao

import (
	"github.com/fishblog/internal/model"
)

func (d *Dao) GetArticleTagByAID(articleID uint32) (model.ArticleTag, error) {
	articleTag := model.ArticleTag{ArticleID: articleID}

	return articleTag.GetByAID(d.engine)
}

func (d *Dao) GetArticleTagListByTID(tagID uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{TagId: tagID}

	return articleTag.ListByTID(d.engine)
}

func (d *Dao) GetArticleTagListByAID(articleIDs []uint32) ([]*model.ArticleTag, error) {
	articleTag := model.ArticleTag{}

	return articleTag.ListByAIDs(d.engine, articleIDs)
}

func (d *Dao) CreateArticleTag(articleID, tagID uint32, createdBy string) error {
	articleTag := model.ArticleTag{
		ArticleID: articleID,
		TagId: tagID,
		Model: &model.Model{
			CreatedBy: createdBy,
		},
	}

	return articleTag.Create(d.engine)
}


func(d *Dao) UpdateArticleTag(articleID, tagId uint32, modifiedBy string) error {
	articleTag := model.ArticleTag{ArticleID: articleID}
	values := make(map[string]interface{},2)
	values["article_id"] = articleID
	values["tag_id"] = tagId
	if modifiedBy != "" {
		values["modified_by"] = modifiedBy
	}

	return articleTag.UpdateOne(d.engine, values)
}

func (d *Dao) DeleteArticleTag(articleID uint32) error {
	articleTag := model.ArticleTag{
		ArticleID: articleID,
	}

	return articleTag.DeleteOne(d.engine)
}



