/**
 * @author lin.tan
 * @date 2023/10/9
 * @description
 */

package lang

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
)

var I18nAnalyzer *Analyzer
var I18nTranslator *I18n

func MustInit() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("init language analyzer error[%w]", err))
	}
	path := cwd + "/conf/locales/"
	accepts := []language.Tag{language.English, language.Chinese}
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, tag := range accepts {
		realPath := path + tag.String() + ".json"
		bundle.MustLoadMessageFile(realPath)
	}
	I18nTranslator = NewI18n(bundle)
	// 使用i18n作为翻译器，使用默认的tag作为结构体解析器，最多下探20层结构(包括array，slice，map和struct)
	I18nAnalyzer = NewAnalyzer(I18nTranslator, &DefaultTag{}, WithLevel(20))
}
