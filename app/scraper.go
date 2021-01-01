package app

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
)

// Scrape a word from translator website
func (app *Application) ScrapeWord(dictionary string, output_type string, word string) (string, error) {

	word = uniformWord(dictionary, word)
	var wordUrl = getWordUrl(app.Config.BaseUrl, dictionary, word)
	// for iframe solution does not need to do request call
	if output_type == OUTPUT_TYPE_IFRAME {
		iframeMarkup := "<body style=\"top:0;left: 0;width:100%;height: 100%; position: absolute; border: none;overflow:hidden;\">\n<iframe src=\"IFRAME_SOURCE\" frameborder=\"0\" style=\"top:0;left: 0;width:100%;height: 100%; position: absolute; border: none\"></iframe>\n</body>"
		return iframeMarkup, nil
	}
	// if user want the full html of the page
	res, err := http.Get(wordUrl)
	if err != nil || res.StatusCode != 200 {
		app.Logger.Error("err to get the word", zap.Error(err), zap.Any("word", word), zap.Any("status_code", res.StatusCode))
		return "", err
	}

	defer res.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		app.Logger.Error("err on load html", zap.Error(err), zap.Any("word", word))
		return "", err
	}

	// if user want the full html of the page
	if output_type == OUTPUT_TYPE_FULL_HTML {
		fullPageMarkup, err := doc.Html()
		if err != nil {
			app.Logger.Error("err on get full page markup", zap.Error(err), zap.Any("wordUrl", wordUrl))
		}
		return fullPageMarkup, nil
	}
	//TODO for content output type
	//doc.Find(".entry_content").Each(func(i int, s *goquery.Selection) {
	//	wordTranslate = s.Text()
	//
	//})

	return word, nil
}

func (app *Application) ScrapeAll(dictionary string, output_type string) error {
	inputFile, err := os.Open(app.Config.CSVWordInputFile)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()
	reader := csv.NewReader(inputFile)
	reader.Comma = '|'

	// File writer
	outputFile, err := os.Create(app.Config.CSVWordOutputFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	writer := csv.NewWriter(outputFile)

	defer writer.Flush()
	// Process CSV file line by line
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			app.Logger.Error("err on load a word from csv", zap.Error(err), zap.Any("word", record[0]))
			continue
		}
		app.Logger.Info("start to translate", zap.Any("word", record[0]))
		translatedWorld, err := app.ScrapeWord(dictionary, output_type, record[0])
		if err != nil {
			app.Logger.Error("error on get word from translator website", zap.Error(err), zap.Any("word", record[0]))
			continue
		}
		err = writer.Write([]string{record[0], translatedWorld})
		if err != nil {
			return err
		}
	}
	return nil
}

func uniformWord(dictionary string, word string) string {
	switch dictionary {
	case LONGMAN:
		return strings.ReplaceAll(word, " ", "-")
	case LINGUEE:
		word = strings.ReplaceAll(word, "der ", "")
		word = strings.ReplaceAll(word, "die ", "")
		word = strings.ReplaceAll(word, "das ", "")
		return word
	}
	return ""
}
func getWordUrl(dictionary_url string, dictionary string, word string) string {
	switch dictionary {
	case LONGMAN:
		return fmt.Sprintf("%s/%s", dictionary_url, word)
	case LINGUEE:
		return fmt.Sprintf("%s%s", dictionary_url, word)
	}
	return ""
}
