package subtitles

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type AssDialogueFonts struct {
	FileName string
	Fonts    []string
}

func ListDialogueFontsByAssFiles() ([]AssDialogueFonts, error) {
	results := make([]AssDialogueFonts, 0)

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.ToLower(filepath.Ext(path)) != ".ass" {
			return nil
		}

		fonts, err := listDialogueFontsInAssFile(path)
		if err != nil {
			return err
		}

		results = append(results, AssDialogueFonts{
			FileName: path,
			Fonts:    fonts,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].FileName < results[j].FileName
	})

	return results, nil
}

func listDialogueFontsInAssFile(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fontSet := make(map[string]struct{})
	fonts := make([]string, 0)

	text := string(content)
	i := 0
	for i < len(text)-2 {
		if text[i] != '\\' || text[i+1] != 'f' || text[i+2] != 'n' {
			i++
			continue
		}

		start := i + 3
		j := start
		for j < len(text) {
			if text[j] == '\\' || text[j] == '}' || text[j] == '\r' || text[j] == '\n' {
				break
			}
			j++
		}

		name := strings.TrimSpace(text[start:j])
		if name == "" {
			i = j
			continue
		}

		if _, exists := fontSet[name]; exists {
			i = j
			continue
		}

		fontSet[name] = struct{}{}
		fonts = append(fonts, name)
		i = j
	}

	sort.Strings(fonts)
	return fonts, nil
}
