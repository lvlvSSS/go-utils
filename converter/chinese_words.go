package converter

import (
	"github.com/Chain-Zhang/pinyin"
	"regexp"
)

func WordToPinyin(text string) (string, error) {
	chineseReg := regexp.MustCompile("[\u4e00-\u9fa5]")
	if !chineseReg.MatchString(text) {
		return text, nil
	}
	p := pinyin.New(text)
	return p.Convert()
}
