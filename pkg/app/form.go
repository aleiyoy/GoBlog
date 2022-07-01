package app

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	"strings"
)

/*
针对入参校验的方法进行了二次封装，在 BindAndValid 方法中，
通过 ShouldBind 进行参数绑定和入参校验，当发生错误后，再通
过上一步在中间件 Translations 设置的 Translator 来对错
误消息体进行具体的翻译行为。
*/


type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

// 看 error 接口的定义 中只有一个Error方法
//在 Go 语言中，如果一个类型实现了某个 interface 中的所有方法，
// 那么编译器就会认为该类型实现了此 interface，它们是”一样“的。
func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	return true, nil
}
