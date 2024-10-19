package helpers

import "regexp"

// Supported file extensions
var IMAGE_EXTENSIONS = []string{".jpg", ".jpeg", ".heic"}
var VIDEO_EXTENSIONS = []string{".mp4", ".mov"}

var DateParserRegex = regexp.MustCompile(`(\d{4})[._-]?(\d{2})[._-]?(\d{2})(?:[._-]?(\d{2}))?(?:[._-]?(\d{2}))?(?:[._-]?(\d{2}))?`)

var CorrectNameRegex = regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2}) (\d{2})\.(\d{2})\.(\d{2})`)