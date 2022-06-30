package model

import "GoBlog/pkg/app"

type Tag struct {
	*Model    // 嵌入公共字段结构体
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}

type TagSwagger struct {
	List  []*Tag
	Pager *app.Pager
}