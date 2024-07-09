/**
 * @author lin.tan
 * @date 2024/1/24
 * @description
 */

package lang

import (
	"context"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18n struct {
	Bundle  *i18n.Bundle
	Filters []language.Tag
}

func NewI18n(bundle *i18n.Bundle) *I18n {
	return &I18n{
		Bundle: bundle,
	}
}

// Local 通过该方法可以拿到标准的i18n的localize，用户可以通过该方法进行灵活定制行为
func (rec *I18n) Local(ctx context.Context) *i18n.Localizer {
	lg := GetLg(ctx)
	return i18n.NewLocalizer(rec.Bundle, lg.String())
}

// Filter 过滤某些语言不进行翻译
func (rec *I18n) Filter(tags ...language.Tag) *I18n {
	rec.Filters = tags
	return rec
}

func (rec *I18n) checkFilter(ctx context.Context) bool {
	lg := GetLg(ctx)
	for _, tag := range rec.Filters {
		if tag == lg {
			return true
		}
	}

	return false
}

// Translate 封装的通过msg获取对应的语言信息，没有获取到或者出错则返回入参msg
// 该方法适用于直接将msg作为id的翻译文本
// 实现接口方法 func(context.Context, ...string) ([]string, error)
// 返回的err是翻译过程中最后的一个报错，如果翻译出错则返回自己
func (rec *I18n) Translate(ctx context.Context, lines ...string) (ret []string, err error) {
	if rec.checkFilter(ctx) {
		return lines, nil
	}
	ret = make([]string, len(lines))
	for key, line := range lines {
		if line == "" {
			continue
		}
		translate, e := rec.Local(ctx).Localize(&i18n.LocalizeConfig{
			MessageID: line,
		})
		if e != nil {
			err = e
			translate = line
		}
		ret[key] = translate
	}

	return
}

// TranslateOne 单条翻译
func (rec *I18n) TranslateOne(ctx context.Context, line string) (string, error) {
	ret, err := rec.Translate(ctx, line)
	return ret[0], err
}
