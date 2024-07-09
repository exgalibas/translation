/**
 * @author lin.tan
 * @date 2024/1/29
 * @description
 */

package lang

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// LineParse 行模式翻译
func LineParse(ctx context.Context, translator Translator, input ...string) ([]string, error) {
	return translator.Translate(ctx, input...)
}

// TemplateParse 模板模式翻译
func TemplateParse(ctx context.Context, translator Translator, input ...string) ([]string, error) {
	if len(input) <= 0 {
		return input, nil
	}

	if trans, err := LineParse(ctx, translator, input...); err == nil {
		return trans, err
	}

	var lastErr error
	// 正则匹配非空格非中文非Unicode标点符号
	reg := regexp.MustCompile(`[^\s\p{Han}\p{P}\p{S}]+`)
	for k, in := range input {
		matches := reg.FindAllString(in, -1)
		if len(matches) > 0 {
			in = reg.ReplaceAllString(in, `@$`)
		}

		if strings.TrimSpace(in) == "" {
			continue
		}
		trans, err := translator.TranslateOne(ctx, in)
		if err != nil {
			lastErr = err
			continue
		}
		on := make([]string, 0)
		for k, match := range matches {
			on = append(on, fmt.Sprintf("@%d$", k+1), match)
		}
		if len(on) <= 0 {
			continue
		}
		input[k] = strings.NewReplacer(on...).Replace(trans)
	}
	return input, lastErr
}
