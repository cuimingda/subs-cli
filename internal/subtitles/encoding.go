package subtitles

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/saintfish/chardet"
)

const UnknownEncoding = "UNKNOWN"

type FileEncoding struct {
	FileName string
	Encoding string
}

const nonUTF8Instruction = "Please run `subs encoding reset` to convert subtitle files to UTF-8 first."

var detectTextEncoding = defaultDetectTextEncoding

func ListCurrentDirSubtitleFileEncodings() ([]FileEncoding, error) {
	files, err := ListCurrentDirSubtitleFiles()
	if err != nil {
		return nil, err
	}

	encodings := make([]FileEncoding, 0, len(files))
	for _, file := range files {
		encoding := detectFileEncoding(file)
		encodings = append(encodings, FileEncoding{
			FileName: file,
			Encoding: encoding,
		})
	}

	return encodings, nil
}

func EnsureCurrentDirAssFilesUTF8(files []string) error {
	if len(files) == 0 {
		return nil
	}

	return ensureFilesUTF8(files)
}

func ensureFilesUTF8(files []string) error {
	nonUTF8 := make([]string, 0)
	for _, file := range files {
		if err := validateSubtitleFileSize(file); err != nil {
			return err
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if utf8.Valid(content) {
			continue
		}

		encoding := strings.ToUpper(detectFileEncoding(file))
		if encoding == "" || encoding == UnknownEncoding {
			nonUTF8 = append(nonUTF8, fmt.Sprintf("%s (UNKNOWN)", file))
			continue
		}

		if encoding != "UTF-8" {
			nonUTF8 = append(nonUTF8, fmt.Sprintf("%s (%s)", file, encoding))
		}
	}

	if len(nonUTF8) == 0 {
		return nil
	}

	return fmt.Errorf("%s. %s", strings.Join(nonUTF8, ", "), nonUTF8Instruction)
}

func detectFileEncoding(file string) string {
	if err := validateSubtitleFileSize(file); err != nil {
		return UnknownEncoding
	}

	content, err := os.ReadFile(file)
	if err != nil || len(content) == 0 {
		return UnknownEncoding
	}

	encoding, err := detectTextEncoding(content)
	if err != nil || encoding == "" {
		return UnknownEncoding
	}

	return strings.ToUpper(encoding)
}

func defaultDetectTextEncoding(content []byte) (string, error) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(content)
	if err != nil {
		return "", err
	}

	if result == nil || result.Charset == "" {
		return "", errors.New("encoding detection failed")
	}

	return result.Charset, nil
}
