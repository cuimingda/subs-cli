package subtitles

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type AssStyleFonts struct {
	FileName string
	Fonts    []string
}

type AssStyleFontResetResult struct {
	TotalAssFiles int
	UpdatedFiles  int
	UpdatedFonts  int
}

func ListStyleFontsByAssFiles() ([]AssStyleFonts, error) {
	results := make([]AssStyleFonts, 0)
	files, err := listCurrentDirAssFiles()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fonts, err := listStyleFontsInAssFile(file)
		if err != nil {
			return nil, err
		}

		results = append(results, AssStyleFonts{
			FileName: file,
			Fonts:    fonts,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].FileName < results[j].FileName
	})

	return results, nil
}

func listStyleFontsInAssFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	fontSet := make(map[string]struct{})
	fonts := make([]string, 0)
	scanner := bufio.NewScanner(file)

	inStylesSection := false
	fontNameIndex := -1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if line == "[V4+ Styles]" {
			inStylesSection = true
			fontNameIndex = -1
			continue
		}

		if !inStylesSection {
			continue
		}

		if strings.HasPrefix(line, "[") && line != "[V4+ Styles]" {
			break
		}

		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "format:") {
			fontNameIndex = parseFormatForFontIndex(line)
			continue
		}

		if !strings.HasPrefix(lower, "style:") || fontNameIndex < 0 {
			continue
		}

		styleName := strings.TrimSpace(line[len("Style:"):])
		columns := splitAssStyleFields(styleName)
		if fontNameIndex >= len(columns) {
			continue
		}

		fontName := strings.TrimSpace(columns[fontNameIndex])
		if fontName == "" {
			continue
		}

		if _, exists := fontSet[fontName]; exists {
			continue
		}

		fontSet[fontName] = struct{}{}
		fonts = append(fonts, fontName)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fonts, nil
}

func parseFormatForFontIndex(line string) int {
	prefixRemoved := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(line), "format:"))
	formatFields := splitAssStyleFields(prefixRemoved)
	for idx, field := range formatFields {
		if strings.EqualFold(strings.TrimSpace(field), "Fontname") {
			return idx
		}
	}

	return -1
}

func splitAssStyleFields(value string) []string {
	return strings.Split(value, ",")
}

func ResetCurrentDirAssStyleFontsToMicrosoftYaHei() (AssStyleFontResetResult, error) {
	result := AssStyleFontResetResult{}
	files, err := listCurrentDirAssFiles()
	if err != nil {
		return AssStyleFontResetResult{}, err
	}

	result.TotalAssFiles = len(files)

	for _, file := range files {
		updatedStyles, err := resetStyleFontsInAssFile(file)
		if err != nil {
			return AssStyleFontResetResult{}, err
		}
		if updatedStyles > 0 {
			result.UpdatedFiles++
			result.UpdatedFonts += updatedStyles
		}
	}

	return result, nil
}

func resetStyleFontsInAssFile(path string) (int, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(content), "\n")
	out := make([]string, 0, len(lines))

	inStylesSection := false
	fontNameIndex := -1
	updated := 0

	for _, line := range lines {
		rawLine := strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(rawLine)

		if trimmed == "[V4+ Styles]" {
			inStylesSection = true
			fontNameIndex = -1
			out = append(out, rawLine)
			continue
		}

		if !inStylesSection {
			out = append(out, rawLine)
			continue
		}

		if strings.HasPrefix(trimmed, "[") && trimmed != "[V4+ Styles]" {
			inStylesSection = false
			out = append(out, rawLine)
			continue
		}

		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "format:") {
			fontNameIndex = parseFormatForFontIndex(trimmed)
			out = append(out, rawLine)
			continue
		}

		if !strings.HasPrefix(lower, "style:") || fontNameIndex < 0 {
			out = append(out, rawLine)
			continue
		}

		rest := strings.TrimSpace(rawLine[len("Style:"):])
		columns := splitAssStyleFields(rest)
		if fontNameIndex >= len(columns) {
			out = append(out, rawLine)
			continue
		}

		if strings.TrimSpace(columns[fontNameIndex]) == "Microsoft YaHei" {
			out = append(out, rawLine)
			continue
		}

		columns[fontNameIndex] = "Microsoft YaHei"
		out = append(out, "Style: "+strings.Join(columns, ","))
		updated++
	}

	if updated == 0 {
		return 0, nil
	}

	returned := strings.Join(out, "\n")
	if err := os.WriteFile(path, []byte(returned), 0o644); err != nil {
		return 0, err
	}

	return updated, nil
}
