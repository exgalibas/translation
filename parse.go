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

func LineParse(ctx context.Context, translator Translator, input ...string) ([]string, error) {
	return translator.Translate(ctx, input...)
}

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
		oldnew := make([]string, 0)
		for k, match := range matches {
			oldnew = append(oldnew, fmt.Sprintf("@%d$", k+1), match)
		}
		if len(oldnew) <= 0 {
			continue
		}
		input[k] = strings.NewReplacer(oldnew...).Replace(trans)
	}
	return input, lastErr
}
