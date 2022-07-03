package dao

import "GoBlog/internal/model"
import "GoBlog/pkg/app"

// 在 dao 层进行了数据访问对象的封装，并针对业务所需的字段进行了处理。


func (d *Dao) CountTag(name string, state uint8) (int, error) {
	tag := model.Tag{Name: name, State: state}
	return tag.Count(d.engine)
}

func (d *Dao) GetTag(id uint32, state uint8) (model.Tag, error) {
	tag := model.Tag{Model: &model.Model{ID: id}, State: state}
	return tag.Get(d.engine)
}

func (d *Dao) GetTagList(name string, state uint8, page, pageSize int) ([]*model.Tag, error) {
	tag := model.Tag{Name: name, State: state}
	pageOffset := app.GetPageOffset(page, pageSize)
	return tag.List(d.engine, pageOffset, pageSize)
}

func (d *Dao) CreateTag(name string, state uint8, createdBy string) error {
	tag := model.Tag{
		Name:  name,
		State: state,
		Model: &model.Model{CreatedBy: createdBy},
	}

	return tag.Create(d.engine)
}

func (d *Dao) UpdateTag(id uint32, name string, state uint8, modifiedBy string) error {
	tag := model.Tag{
		//Name:  name,
		//State: state,
		//Model: &model.Model{ID: id, ModifiedBy: modifiedBy},
		Model: &model.Model{ID: id},
	}

	// 这里不用结构体更新数据库，是因为gorm很难判断结构体中state字段的0，是空值还是真实并有意义的值
	values := map[string]interface{}{
		"state" : 		state,
		"modifiedBy":	modifiedBy,
	}

	if name != ""{
		values["name"] = name
	}

	return tag.Update(d.engine, values)
}

func (d *Dao) DeleteTag(id uint32) error {
	tag := model.Tag{Model: &model.Model{ID: id}}
	return tag.Delete(d.engine)
}
