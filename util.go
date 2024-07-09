/**
 * @author lin.tan
 * @date 2024/1/25
 * @description
 */

package lang

import (
	"context"
	"golang.org/x/text/language"
)

// TranslateOne 单条翻译，调用频度较高的通用方法，内部有错误忽略和中文忽略，方便使用
func TranslateOne(ctx context.Context, line string) string {
	trans, _ := I18nTranslator.Filter(language.Chinese).TranslateOne(ctx, line)
	return trans
}

// Translate 多条翻译，调用频度较高的通用方法，内部有错误忽略和中文忽略，方便使用
func Translate(ctx context.Context, lines ...string) []string {
	trans, _ := I18nTranslator.Filter(language.Chinese).Translate(ctx, lines...)
	return trans
}
