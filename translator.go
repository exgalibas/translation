/**
 * @author lin.tan
 * @date 2024/1/24
 * @description 翻译器，定义翻译方法，供解析器analyzer使用
 */

package lang

import (
	"context"
	"golang.org/x/text/language"
)

type Translator interface {
	Translate(context.Context, ...string) ([]string, error)
	TranslateOne(context.Context, string) (string, error)
	Filter(tags ...language.Tag) Translator
}
