package subtitles

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type AssDialogueFonts struct {
	FileName string
	Fonts    []string
}

type DialogueFontPruneResult struct {
	TotalAssFiles int
	RemovedTags   int
}

func ListDialogueFontsByAssFiles() ([]AssDialogueFonts, error) {
	results := make([]AssDialogueFonts, 0)
	files, err := listCurrentDirAssFiles()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fonts, err := listDialogueFontsInAssFile(file)
		if err != nil {
			return nil, err
		}

		results = append(results, AssDialogueFonts{
			FileName: file,
			Fonts:    fonts,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].FileName < results[j].FileName
	})

	return results, nil
}

func PruneDialogueFontTagsFromAssFiles() (DialogueFontPruneResult, error) {
	result := DialogueFontPruneResult{}
	files, err := listCurrentDirAssFiles()
	if err != nil {
		return DialogueFontPruneResult{}, err
	}

	for _, file := range files {
		result.TotalAssFiles++
		if err := validateSubtitleFileSize(file); err != nil {
			return DialogueFontPruneResult{}, err
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return DialogueFontPruneResult{}, err
		}

		prunedContent, removedTags := pruneDialogueFontTagsInText(string(content))
		if removedTags == 0 {
			continue
		}

		if err := writeFilePreserveMode(file, []byte(prunedContent)); err != nil {
			return DialogueFontPruneResult{}, err
		}
		result.RemovedTags += removedTags
	}

	return result, nil
}

func listCurrentDirAssFiles() ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		suffix := strings.ToLower(filepath.Ext(entry.Name()))
		if suffix != ".ass" {
			continue
		}

		files = append(files, entry.Name())
	}

	if err := EnsureCurrentDirAssFilesUTF8(files); err != nil {
		return nil, err
	}

	return files, nil
}

func listDialogueFontsInAssFile(path string) ([]string, error) {
	if err := validateSubtitleFileSize(path); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fonts := extractDialogueFonts(string(content))

	sort.Strings(fonts)
	return fonts, nil
}

func pruneDialogueFontTagsInText(text string) (string, int) {
	pruned := bytes.Buffer{}
	removedTags := 0
	last := 0

	i := 0
	for i < len(text) {
		if isDialogueFontTagAt(text, i) {
			end := findDialogueFontTagEnd(text, i+3)
			pruned.WriteString(text[last:i])
			last = end
			i = end
			removedTags++
			continue
		}

		i++
	}

	if removedTags == 0 {
		return "", 0
	}

	pruned.WriteString(text[last:])
	return pruned.String(), removedTags
}

func extractDialogueFonts(text string) []string {
	fontSet := make(map[string]struct{})
	fonts := make([]string, 0)

	i := 0
	for i < len(text)-2 {
		if !isDialogueFontTagAt(text, i) {
			i++
			continue
		}

		end := findDialogueFontTagEnd(text, i+3)
		name := strings.TrimSpace(text[i+3 : end])
		if name != "" {
			if _, exists := fontSet[name]; !exists {
				fontSet[name] = struct{}{}
				fonts = append(fonts, name)
			}
		}
		i = end
	}

	return fonts
}

func isDialogueFontTagAt(text string, i int) bool {
	return i+2 < len(text) && text[i] == '\\' && text[i+1] == 'f' && text[i+2] == 'n'
}

func findDialogueFontTagEnd(text string, start int) int {
	j := start
	for j < len(text) {
		switch text[j] {
		case '\\', '}', '\r', '\n':
			return j
		default:
			j++
		}
	}
	return j
}
