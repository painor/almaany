package almaany

import (
	"github.com/grokify/html-strip-tags-go"
	"strconv"
)

func FormatMaany(maana Manaa) string {
	title := "📕الكلمة : " + "<code>" + maana.Word + "</code>"

	wordType := "📒النوع : " + "<code>" + maana.WordType + "</code>"
	explanations := "🔍التفاسير : \n"

	for index, explain := range maana.Explanations {
		temp := strconv.Itoa(index+1) + " - " + explain
		explanations += "<code>" + strip.StripTags(temp) + "</code>"
		explanations += "\n"
	}
	return title + "\n" + wordType + "\n" + explanations
}
