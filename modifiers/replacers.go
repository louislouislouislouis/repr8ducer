package modifiers

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/louislouislouislouis/repr8ducer/ui/message"
	"github.com/louislouislouislouis/repr8ducer/utils"
)

// TODO add goroutine to improve perf

type urlReplacer struct {
	urls map[string]string
}

func NewUrlReplacer() *urlReplacer {
	return &urlReplacer{
		urls: make(map[string]string, 0),
	}
}

var urlRegex *regexp.Regexp = regexp.MustCompile(
	`(?:https?|jdbc):\/\/(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,6}(?::\d+)?(?:\/[^\s]*)?`,
)

func (u *urlReplacer) SearchUrlsInDir(dirPath string) (message.UrlDetectionMsgValue, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return message.UrlDetectionMsgValue{}, fmt.Errorf("cannot read dir : %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			u.SearchUrlsInDir(fmt.Sprintf("%s/%s", dirPath, file.Name()))
		} else {

			filePath := filepath.Join(dirPath, file.Name())

			content, err := os.ReadFile(filePath)
			if err != nil {
				utils.Log.Err(err).Msg("cannot read file")
				continue
			}

			matches := urlRegex.FindAllString(string(content), -1)
			for _, v := range matches {
				u.urls[v] = ""
			}
		}
	}

	keys := make([]string, 0, len(u.urls))
	for key := range u.urls {
		keys = append(keys, key)
	}

	return message.UrlDetectionMsgValue{
		Urls: keys,
	}, nil
}

func (u *urlReplacer) ReplaceUrlsInDir(dirPath string, urls map[string]string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("cannot read dir : %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			return u.ReplaceUrlsInDir(fmt.Sprintf("%s/%s", dirPath, file.Name()), urls)
		} else {
			filePath := filepath.Join(dirPath, file.Name())

			content, err := os.ReadFile(filePath)
			if err != nil {
				utils.Log.Err(err).Msg("cannot read file")
				continue
			}

			fileContent := string(content)

			for url, replacement := range urls {
				fileContent = strings.ReplaceAll(fileContent, url, replacement)
			}
			err = os.WriteFile(filePath, []byte(fileContent), 0644)
			if err != nil {
				utils.Log.Err(err).Msg(fmt.Sprintf("cannot write %s", filePath))
				continue
			}

			utils.Log.Debug().Msg(fmt.Sprintf("Modifications appliquées dans : %s\n", filePath))
		}
	}
	return nil
}
