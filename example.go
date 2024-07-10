/**
 * @author lin.tan
 * @date 2024/1/31
 * @description 使用说明和例子
 */

package lang

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

// BaseUse 基础用法
func BaseUse(ctx context.Context) {

	// 设置 ctx中的目标语言
	ctx = SaveLg(ctx, language.English)
	// 设置 gin ctx中的目标语言
	SaveGinLg(&gin.Context{}, language.English)
	// 获取ctx中的目标语言
	_ = GetLg(ctx)
	// 获取ctx中的目标语言，如果没有设置则返回默认
	_ = GetLgWithDefault(ctx, language.English)

	// 初始化翻译器和解析器
	// 默认使用i18n作为翻译器，对应的译本path:conf/locales/*.json
	// 使用者如需定制可自行重写，也可直接使用默认
	MustInit()

	// 使用翻译器进行翻译，支持多个，会返回最后一个翻译错误
	// 对于多个翻译，若某个翻译出错会返回原文本，不会阻断其他翻译，对应错误会返回
	_, _ = I18nTranslator.Translate(ctx, "翻译1", "翻译2")
	// 使用翻译器进行单个翻译，出错返回err
	_, _ = I18nTranslator.TranslateOne(ctx, "翻译1")

	// 过滤某些语言不进行翻译
	_, _ = I18nTranslator.Filter(language.Chinese, language.French).Translate(ctx, "翻译1", "翻译2")
	_, _ = I18nTranslator.Filter(language.Chinese, language.French).TranslateOne(ctx, "翻译1")

	// 该方法包装了I18nTranslator的方法，增加了中文过滤，如果context中解出来是language.chinese，则不会进行翻译
	// 同时也过滤了错误，无法翻译/翻译出错的返回原文本
	_ = Translate(ctx, "翻译1", "翻译2")
	// 对应的单个翻译
	_ = TranslateOne(ctx, "翻译1")

	// 使用analyze进行翻译
	// I18nAnalyzer组合了翻译器I18nTranslator和默认的结构体解析器DefaultTag，可以适用全场景
	// 下面的方法等同于I18nTranslator.Translate
	_, _ = I18nAnalyzer.Translate(ctx, "翻译1", "翻译2")
	// 下面的方法等同于I18nTranslator.TranslateOne
	_, _ = I18nAnalyzer.TranslateOne(ctx, "翻译1")

	// 更多的使用场景
	// 1. 简单字符串
	I18nAnalyzer.Analyze(ctx, "翻译1")
	// 2. array
	I18nAnalyzer.Analyze(ctx, [1]string{"翻译1"})
	// 3. slice
	I18nAnalyzer.Analyze(ctx, []string{"翻译1"})
	// 4. map
	I18nAnalyzer.Analyze(ctx, map[int]string{
		1: "翻译1",
	})
	// 5. 复杂嵌套的slice和map
	I18nAnalyzer.Analyze(ctx, map[string][]map[string][]string{
		"name": []map[string][]string{
			{
				"inner": []string{"翻译"},
			},
		},
	})
	// 6. 结构体
	// 结构体需要添加annotations标注需要翻译，否则不进行翻译
	// lang:"mod=line" 或者 lang:"mod=template"
	// 注意对于结构体嵌套，上层的annotations是可以被下层继承的，默认会进行翻译，如果需要截断翻译行为，可以使用lang:"-"来进行截断
	// line模式会进行全文本匹配翻译，需要在译本中添加全文本和对应的翻译
	// template模式会先进行全文本匹配翻译，如果没有匹配成功则会进行正则匹配，后面赘述
	type InnerS struct {
		// 不进行翻译
		No string
		// line模式进行全匹配翻译
		Name string `lang:"mod=line"`
		// template模板匹配
		InnerName []string `lang:"mod=template"`
	}

	// level层级说明
	// 在初始化Analyzer的时候，可以传递level来指定翻译最深层级，如下指定一层
	analyzer := NewAnalyzer(I18nTranslator, &DefaultTag{}, WithLevel(1))
	analyzer.Analyze(ctx, []string{"翻译1"})     //会进行翻译
	analyzer.Analyze(ctx, [][]string{{"翻译1"}}) //不会进行翻译，因为这里有两层深度
	// 参与层级叠加的有: 嵌套的array,slice，map和struct
	// 当前项目默认初始化的时候设定层级为20，可以自定义

	// template模板说明
	// template模式下，会进行正则匹配，对应的正则表达式 - [^\s\p{Han}\p{P}\p{S}]+
	// 该正则会匹配出所有非中文/空格/unicode中英文标点符号并替换成@$，然后按照替换后的全文去译本中全匹配
	// 匹配成功后，在根据序号替换回来，比如译本中的某个模板如下：
	// "仅计算@$/@$回传的_@$事件数据": "Only calculate _@3$ event from @1$/@2$",
	// 其中value中的@1$/@2$/@3$对应key中从左到右的三个@$，即你可以人为调整翻译后的回写顺序
	// 如果要翻译"仅计算SDK/APP回传的_event事件数据"，就可以使用上述的template进行模版翻译

	// 高级主题
	// 当前的Analyzer默认使用i18n作为翻译器，DefaultTag作为结构体解析器
	// 你可以通过NewAnalyzer定制自己的解析器，只要传递的翻译器实现了接口Translator，结构体解析器实现了接口Tag即可
}
