/**
 * @author lin.tan
 * @date 2024/1/24
 * @description 预置语言到context中，也提供获取能力
 */

package lang

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const (
	LanguageKey = "language"
)

// SaveLg 将语言保存到context中
func SaveLg(ctx context.Context, lg string) context.Context {
	return context.WithValue(ctx, LanguageKey, language.Make(lg))
}

// SaveGinLg 将语言保存到gin的context中
func SaveGinLg(ctx *gin.Context, lg string) {
	ctx.Set(LanguageKey, language.Make(lg))
}

// GetLg 获取上下文语言
func GetLg(ctx context.Context) language.Tag {
	v, _ := ctx.Value(LanguageKey).(language.Tag)
	return v
}

// GetLgWithDefault 获取上下文语言，没有则返回传递的默认语言
func GetLgWithDefault(ctx context.Context, def language.Tag) language.Tag {
	v, ok := ctx.Value(LanguageKey).(language.Tag)
	if !ok || v.IsRoot() {
		return def
	}
	return v
}
