package model

import (
	"GoBlog/pkg/app"
	"github.com/jinzhu/gorm"
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

func (a Article) TableName() string {
	return "blog_article"
}

type ArticleSwagger struct {
	List 	[]*Article
	Pagger 	*app.Pager
}


func(a Article) Create(db *gorm.DB)(*Article , error){
	var err error
	if err = db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}


func (a Article) Update(db *gorm.DB, values interface{}) error {
	// 这里不用结构体更新数据库，是因为gorm很难判断结构体中state字段的0，是空值还是真实并有意义的值
	if err := db.Model(a).Where("id = ? AND is_del = ?", a.ID, 0).Updates(values).Error; err != nil {
		return err
	}
	return nil

	//return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Update(t).Error
}

func(a Article)Get(db *gorm.DB) (Article, error) {

	var article Article
	var err error
	db = db.Where("id = ? AND is_del = ? AND state = ?", a.Model.ID, 0, a.State)

	err = db.First(&article).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return article, err
	}
	return article, nil
}

func(a Article) Delete(db *gorm.DB) error{

	err := db.Where("id = ? AND is_del = ?", a.Model.ID, 0).Delete(&a).Error

	if err != nil{
		return err
	}
	return nil
}

type ArticleRow struct {
	ArticleID     uint32
	TagID         uint32
	TagName       string
	ArticleTitle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
}

func(a Article) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int)([]*ArticleRow, error){
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


func (a Article) CountByTagID(db *gorm.DB, tagID uint32) (int, error) {
	var count int
	err := db.Table(ArticleTag{}.TableName()+" AS at").
		Joins("LEFT JOIN `"+Tag{}.TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+Article{}.TableName()+"` AS ar ON at.article_id = ar.id").
		Where("at.`tag_id` = ? AND ar.state = ? AND ar.is_del = ?", tagID, a.State, 0).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
