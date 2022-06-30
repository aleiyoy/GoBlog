package model


type Tag struct {
	*Model    // 嵌入公共字段结构体
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

func (t Tag) TableName() string {
	return "blog_tag"
}
