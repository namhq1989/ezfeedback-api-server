package manipulation

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/publicsuffix"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func ToTsQueryTerms(text string) string {
	text = Slugify(text)

	var (
		parts      = strings.Split(text, "-")
		totalParts = len(parts)
	)

	if totalParts == 1 {
		return fmt.Sprintf("%s:*", text)
	}

	terms := ""
	for i, part := range parts {
		if i == totalParts-1 {
			terms += fmt.Sprintf("%s:*", part)
		} else {
			terms += fmt.Sprintf("%s:* & ", part)
		}
	}

	return terms
}

func Slugify(text string) string {
	noDiacritics := removeDiacritics(text)

	noDiacritics = strings.ReplaceAll(noDiacritics, "Đ", "D")
	noDiacritics = strings.ReplaceAll(noDiacritics, "đ", "d")

	lower := strings.ToLower(noDiacritics)
	hyphens := strings.ReplaceAll(lower, " ", "-")
	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")
	if err != nil {
		return ""
	}
	safe := reg.ReplaceAllString(hyphens, "")
	return strings.Trim(safe, "-")
}

func removeDiacritics(input string) string {
	t := norm.NFD.String(input)
	var b strings.Builder
	b.Grow(len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func ValidateAndCleanHumanName(s string) string {
	cleanedName := strings.Join(strings.Fields(s), " ")

	for _, ch := range cleanedName {
		if !unicode.IsLetter(ch) && !unicode.IsSpace(ch) {
			return ""
		}
	}

	return cleanedName
}

func RemoveSpacesAndValidateNumber(s string) string {
	digitsOnly := strings.ReplaceAll(s, " ", "")

	for _, ch := range digitsOnly {
		if !unicode.IsDigit(ch) {
			return ""
		}
	}

	return digitsOnly
}

func NonAccentVietnamese(str string) string {
	if str == "" {
		return ""
	}

	str = strings.ToLower(str)
	str = replaceStringWithRegex(str, `đ`, "d")

	// Decompose the Unicode string and remove diacritical marks
	t := transform.Chain(norm.NFD, removeDiacritic{}, norm.NFC)
	result, _, err := transform.String(t, str)
	if err != nil {
		return str
	}

	// Remove any character that is not alphanumeric or a space
	result = replaceStringWithRegex(result, `[^a-zA-Z0-9\s]`, "")

	return result
}

type removeDiacritic struct{}

func (removeDiacritic) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for _, r := range string(src) {
		if !unicode.Is(unicode.Mn, r) { // Skip marks
			n := copy(dst[nDst:], []byte(string(r)))
			nDst += n
		}
		nSrc += len(string(r))
	}
	return nDst, nSrc, nil
}

func (removeDiacritic) Reset() {}

func replaceStringWithRegex(src string, pattern string, replaceText string) string {
	reg := regexp.MustCompile(pattern)
	return reg.ReplaceAllString(src, replaceText)
}

func ExtractOSFromUserAgent(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	switch {
	case strings.Contains(userAgent, "windows"):
		return "windows"
	case strings.Contains(userAgent, "mac os"):
		return "macos"
	case strings.Contains(userAgent, "ios"):
		return "ios"
	case strings.Contains(userAgent, "android"):
		return "android"
	case strings.Contains(userAgent, "linux"):
		return "linux"
	default:
		return "unknown"
	}
}

func GetRootDomain(origin string) (string, error) {
	parsedURL, err := url.Parse(origin)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host

	domain := host
	if colonIndex := strings.IndexByte(host, ':'); colonIndex != -1 {
		domain = host[:colonIndex]
	}

	if domain == "localhost" || strings.HasSuffix(domain, ".localhost") ||
		domain == "127.0.0.1" || strings.HasPrefix(domain, "192.168.") ||
		strings.HasPrefix(domain, "10.") {
		return domain, nil
	}

	if net.ParseIP(domain) != nil {
		return domain, nil
	}

	registrableDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return domain, nil
	}

	return registrableDomain, nil
}
