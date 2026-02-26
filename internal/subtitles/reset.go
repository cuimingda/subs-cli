package subtitles

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	unicodeenc "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type ResetResult struct {
	Total   int
	Updated int
}

func ResetCurrentDirSubtitleFilesToUTF8() (ResetResult, error) {
	files, err := ListCurrentDirSubtitleFiles()
	if err != nil {
		return ResetResult{}, err
	}

	result := ResetResult{Total: len(files)}

	for _, file := range files {
		updated, err := resetFileToUTF8(file)
		if err != nil {
			return result, err
		}

		if updated {
			result.Updated++
		}
	}

	return result, nil
}

func resetFileToUTF8(file string) (bool, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return false, err
	}

	encoding := detectFileEncoding(file)
	if isUTF8Encoding(encoding) {
		if utf8.Valid(content) {
			return false, nil
		}

		encoding = UnknownEncoding
	}

	converted, err := convertToUTF8(content, encoding)
	if err != nil {
		return false, err
	}

	if bytes.Equal(content, converted) {
		return false, nil
	}

	if err := os.WriteFile(file, converted, 0o644); err != nil {
		return false, err
	}

	return true, nil
}

func isUTF8Encoding(encoding string) bool {
	return normalizeEncoding(encoding) == "UTF-8"
}

func convertToUTF8(content []byte, encoding string) ([]byte, error) {
	normalized := normalizeEncoding(encoding)
	if normalized == "UTF-8" {
		return content, nil
	}

	decoder, err := decoderByEncoding(normalized)
	if err != nil {
		return nil, err
	}

	decoded, err := io.ReadAll(transform.NewReader(bytes.NewReader(content), decoder.NewDecoder()))
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func decoderByEncoding(encoding string) (encoding.Encoding, error) {
	switch normalizeEncoding(encoding) {
	case "GBK":
		return simplifiedchinese.GBK, nil
	case "GB2312":
		return simplifiedchinese.HZGB2312, nil
	case "GB18030":
		return simplifiedchinese.GB18030, nil
	case "GB-18030":
		return simplifiedchinese.GB18030, nil
	case "UTF-16":
		return unicodeenc.UTF16(unicodeenc.LittleEndian, unicodeenc.UseBOM), nil
	case "UTF-16LE":
		return unicodeenc.UTF16(unicodeenc.LittleEndian, unicodeenc.UseBOM), nil
	case "UTF-16BE":
		return unicodeenc.UTF16(unicodeenc.BigEndian, unicodeenc.UseBOM), nil
	case "UTF16":
		return unicodeenc.UTF16(unicodeenc.LittleEndian, unicodeenc.UseBOM), nil
	case "UTF16LE":
		return unicodeenc.UTF16(unicodeenc.LittleEndian, unicodeenc.UseBOM), nil
	case "UTF16BE":
		return unicodeenc.UTF16(unicodeenc.BigEndian, unicodeenc.UseBOM), nil
	case "ANSI", "WINDOWS-1252", "WINDOWS1252":
		return charmap.Windows1252, nil
	case "ISO-8859-1", "ISO-8859", "ISO8859-1", "LATIN1":
		return charmap.Windows1252, nil
	}

	return nil, fmt.Errorf("unsupported encoding: %s", encoding)
}

func normalizeEncoding(encoding string) string {
	encoding = strings.ToUpper(strings.TrimSpace(encoding))
	encoding = strings.ReplaceAll(encoding, "_", "-")
	if encoding == "UTF8" {
		encoding = "UTF-8"
	}

	return encoding
}
