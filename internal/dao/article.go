package dao

import (
	"github.com/fishblog/internal/model"
	"github.com/fishblog/pkg/app"
)

type Article struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	IsDel         uint8  `json:"is_del"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         uint8  `json:"state"`
}

func (d *Dao) CreateArticle(param *Article) (*model.Article, error) {
	article := model.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		Model:         &model.Model{CreatedBy: param.CreatedBy},
	}

	return article.Create(d.engine)
}

func (d *Dao) UpdateArticle(param *Article) error {
	article := model.Article{
		Model: &model.Model{ID: param.ID},
	}
	values := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}
	if param.Title != "" {
		values["title"] = param.Title
	}
	if param.CoverImageUrl != "" {
		values["CoverImageUrl"] = param.CoverImageUrl
	}
	if param.Content != "" {
		values["content"] = param.Content
	}
	if param.Desc != "" {
		values["desc"] = param.Desc
	}

	return article.Update(d.engine, values)
}

func (d *Dao) DeleteArticle(id uint32) error {
	article := model.Article{
		Model: &model.Model{ID: id},
	}
	return article.Delete(d.engine)
}

func (d *Dao) GetArticle(id uint32, state uint8) (model.Article, error) {
	article := model.Article{
		State: state,
		Model: &model.Model{
			ID: id,
		},
	}
	return article.Get(d.engine)
}

func (d *Dao) CountArticleListByTID(tid uint32, state uint8) (int, error) {
	article := model.Article{
		State: state,
	}

	return article.CountByTagID(d.engine, tid)
}

func (d *Dao) GetArticleListByTID(tid uint32, state uint8, page, pageSize int) ([]*model.ArticleRow, error) {
	article := model.Article{State: state}

	return article.ListByTagID(d.engine, tid, app.GetPageOffset(page, pageSize), pageSize)
}
