/**
 * @author lin.tan
 * @date 2023/10/9
 * @description
 */

package lang

import (
	"context"
	"math"
	"reflect"

	"golang.org/x/text/language"
)

type Analyzer struct {
	// 可插拔
	*baseAnalyzer
	// 解析方法
	parse Parse
	level int
}

// baseAnalyzer 基础解析器，里面包含了翻译器和结构体解析器，这两个都是定义的接口，可进行灵活插拔定制
type baseAnalyzer struct {
	Tag
	Translator
}

type Option func(*Analyzer)

// WithParse 自定义解析方法
func WithParse(parse Parse) Option {
	return func(analyzer *Analyzer) {
		analyzer.parse = parse
	}
}

// WithLevel 自定义解析层级
func WithLevel(level int) Option {
	return func(analyzer *Analyzer) {
		analyzer.level = level
	}
}

func NewAnalyzer(translator Translator, tag Tag, ops ...Option) *Analyzer {
	analyzer := &Analyzer{
		level: math.MaxInt,
		baseAnalyzer: &baseAnalyzer{
			tag,
			translator,
		},
		parse: LineParse,
	}

	for _, op := range ops {
		op(analyzer)
	}
	return analyzer
}

func (analyzer *Analyzer) Analyze(ctx context.Context, value any) any {
	// 中文不进行翻译，跳过提效
	if GetLgWithDefault(ctx, language.Chinese) == language.Chinese {
		return value
	}
	sv := reflect.ValueOf(value)
	if !sv.IsValid() {
		return value
	}
	wrap := false
	if sv.Kind() != reflect.Ptr {
		nsv := reflect.New(sv.Type())
		nsv.Elem().Set(sv)
		sv = nsv
		wrap = true
	}
	analyzer.recursiveAnalyze(ctx, sv, analyzer.parse, analyzer.level)
	if wrap {
		return sv.Elem().Interface()
	}
	return sv.Interface()
}

type Parse func(context.Context, Translator, ...string) ([]string, error)

func (base *baseAnalyzer) recursiveAnalyze(ctx context.Context, v reflect.Value, parse Parse, level int) {
	if !v.IsValid() || level < 0 {
		return
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		base.recursiveAnalyze(ctx, v.Elem(), parse, level)
	case reflect.Struct:
		level--
		var vType reflect.Type
		if v.CanInterface() {
			vType = reflect.TypeOf(v.Interface())
		}
		for i := 0; i < v.NumField(); i++ {
			if vType != nil {
				p, err := base.AnalyzeTag(vType.Field(i).Tag.Get("lang"))
				if err != nil {
					// 结构体解析报错，截断
					continue
				}
				base.recursiveAnalyze(ctx, v.Field(i), p, level)
			}
		}
	case reflect.String:
		if v.CanSet() && parse != nil {
			str, err := parse(ctx, base.Translator, v.String())
			if err == nil && len(str) > 0 {
				v.SetString(str[0])
			}
		}
	case reflect.Slice, reflect.Array:
		level--
		for i := 0; i < v.Len(); i++ {
			base.recursiveAnalyze(ctx, v.Index(i), parse, level)
		}
	case reflect.Map:
		level--
		keys := v.MapKeys()
		for _, k := range keys {
			base.recursiveAnalyze(ctx, v.MapIndex(k), parse, level)
		}
	}
}
