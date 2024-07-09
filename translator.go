/**
 * @author lin.tan
 * @date 2024/1/24
 * @description 翻译器，定义翻译方法，供解析器analyzer使用
 */

package lang

import (
	"context"
)

// Translator 翻译器接口
type Translator interface {
	Translate(context.Context, ...string) ([]string, error)
	TranslateOne(context.Context, string) (string, error)
}
