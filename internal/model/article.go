package model

import (
	"github.com/fishblog/pkg/app"
	"github.com/jinzhu/gorm"
)

type Article struct {
	*Model
	Title 			string `json:"title"`
	Desc 			string `json:"desc"`
	Content			string `json:"content"`
	CoverImageUrl	string `json:"cover_image_url"`
	State 			uint8 `json:"state"`
}

type ArticleRow struct {
	ArticleID    	uint32
	TagID			uint32
	TagName			string
	ArticleTitle	string
	ArticleDesc		string
	CoverImageUrl   string
	Content			string
}

type ArticleSwagger struct {
	List  []*Article
	Pager   *app.Pager
}

func (a Article) TableName() string {
	return "blog_article"
}

func (a Article) Create(db *gorm.DB) (*Article, error) {
	if err := db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (a Article) Update(db *gorm.DB, values interface{}) error {
	return  db.Model(&a).Where("id = ? AND is_del = ?", a.ID, 0).Update(values).Error
}

func (a Article) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", a.ID, 0).Delete(&a).Error
}


func (a Article) Count(db *gorm.DB) (int, error) {
	var count int
	if a.Title != "" {
		db = db.Where("name = ?", a.Title)
	}
	db = db.Where("state = ?", a.State)
	if err := db.Model(&a).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (a Article) Get(db *gorm.DB) (Article, error) {
	var article Article
	db = db.Where("id = ? AND is_del = ? AND state = ?", a.ID, 0, a.State)
	//单条记录查询，所以改Find方法为First
	//err := db.Find(&article).Error
	err := db.First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}
	return article, nil
}

//直接查询文章列表表不就可以吗?? no, blog_article no tag_id
func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName() + " AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"`AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (a Article) List(db *gorm.DB, pageOffset, pageSize int) ([]*Article, error ) {
	var articles []*Article
	var err error

	if pageOffset > 0 && pageSize > 0 {
		db.Offset(pageOffset).Limit(pageSize)
	}

	if a.Title != "" {
		db = db.Where("title = ?", a.Title)
	}
	db = db.Where("is_del = ? ", 0)
	if err = db.Model(&a).Where("state = ?", a.State).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, err
}

// conjunctive query / join query
func (a Article) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int) ([]*ArticleRow, error){
	fields := []string{"ar.id AS article_id", "ar.title AS article_title", "ar.desc AS article_desc", "ar.cover_image_url", "ar.content"}
	fields = append(fields, []string{"t.id AS tag_id", "t.name AS tag_name"}...)

	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	rows, err := db.Select(fields).Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*ArticleRow
	for rows.Next() {
		r := &ArticleRow{}
		if err := rows.Scan(&r.ArticleID, &r.ArticleTitle, &r.ArticleDesc, &r.CoverImageUrl, &r.Content, &r.TagID, &r.TagName); err != nil {
			return nil, err
		}

		articles = append(articles, r)
	}

	return articles, nil
}








