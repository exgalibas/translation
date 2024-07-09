/**
 * @author lin.tan
 * @date 2023/10/9
 * @description 结构体标注，用于解析结构体的多语言
 */

package lang

import (
	"errors"
	"strings"
)

type Tag interface {
	AnalyzeTag(string) (Parse, error)
}

type DefaultTag struct {
	mod string
}

const (
	ModBlock    = "-"
	ModLine     = "line"
	ModTemplate = "template"
)

var BlockErr = errors.New("tag mod block")
var ModErr = errors.New("tag mod invalid")

func (tag *DefaultTag) AnalyzeTag(tagStr string) (Parse, error) {

	if tag == nil || tagStr == "" {
		return nil, nil
	}
	// 返回error截断内层翻译
	if tagStr == ModBlock {
		return nil, BlockErr
	}

	tags := strings.Split(tagStr, ",")
	for _, tagItem := range tags {
		items := strings.Split(strings.TrimSpace(tagItem), "=")
		if len(items) < 1 {
			continue
		}

		switch items[0] {
		case "mod":
			tag.mod = items[1]
		}
	}

	switch tag.mod {
	case ModLine:
		return LineParse, nil
	case ModTemplate:
		return TemplateParse, nil
	default:
		return nil, ModErr
	}
}
