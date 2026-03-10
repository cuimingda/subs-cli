package subtitles

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	srtIndexLineRE  = regexp.MustCompile(`^\d+$`)
	srtTimingLineRE = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}[,.]\d{3}\s+-->\s+\d{2}:\d{2}:\d{2}[,.]\d{3}`)
	assTagRE        = regexp.MustCompile(`\{[^}]*\}`)
	htmlTagRE       = regexp.MustCompile(`<[^>]+>`)
)

func IsLikelyChineseEnglishBilingual(file string) bool {
	if err := validateSubtitleFileSize(file); err != nil {
		return false
	}

	content, err := os.ReadFile(file)
	if err != nil || len(content) == 0 {
		return false
	}

	if !utf8.Valid(content) {
		encoding := detectFileEncoding(file)
		if isUTF8Encoding(encoding) {
			encoding = UnknownEncoding
		}
		if converted, err := convertToUTF8(content, encoding); err == nil {
			content = converted
		}
	}

	text := extractSubtitleText(filepath.Ext(file), string(content))
	return containsChinese(text) && containsEnglish(text)
}

func extractSubtitleText(ext, content string) string {
	switch strings.ToLower(ext) {
	case ".ass", ".ssa":
		return extractASSText(content)
	case ".srt":
		return extractSRTText(content)
	default:
		return cleanSubtitleText(content)
	}
}

func extractSRTText(content string) string {
	lines := strings.Split(content, "\n")
	textLines := make([]string, 0, len(lines))

	for _, rawLine := range lines {
		line := strings.TrimSpace(strings.TrimRight(rawLine, "\r"))
		if line == "" || srtIndexLineRE.MatchString(line) || srtTimingLineRE.MatchString(line) {
			continue
		}
		textLines = append(textLines, cleanSubtitleText(line))
	}

	return strings.Join(textLines, "\n")
}

func extractASSText(content string) string {
	lines := strings.Split(content, "\n")
	textLines := make([]string, 0, len(lines))

	for _, rawLine := range lines {
		line := strings.TrimSpace(strings.TrimRight(rawLine, "\r"))
		if !strings.HasPrefix(strings.ToLower(line), "dialogue:") {
			continue
		}

		dialogue := strings.TrimSpace(line[len("Dialogue:"):])
		fields := strings.SplitN(dialogue, ",", 10)
		if len(fields) < 10 {
			continue
		}

		textLines = append(textLines, cleanSubtitleText(fields[9]))
	}

	return strings.Join(textLines, "\n")
}

func cleanSubtitleText(text string) string {
	cleaned := assTagRE.ReplaceAllString(text, " ")
	cleaned = htmlTagRE.ReplaceAllString(cleaned, " ")
	cleaned = strings.NewReplacer(`\N`, " ", `\n`, " ", `\h`, " ").Replace(cleaned)
	return cleaned
}

func containsChinese(text string) bool {
	for _, r := range text {
		if unicode.In(r, unicode.Han) {
			return true
		}
	}
	return false
}

func containsEnglish(text string) bool {
	for _, r := range text {
		if unicode.IsLetter(r) && unicode.In(r, unicode.Latin) {
			return true
		}
	}
	return false
}
