package validation

import (
	"regexp"
	"strings"
)

var (
	emailPattern = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// IsValidEmail ...
func IsValidEmail(value string) bool {
	return emailPattern.MatchString(value)
}

func IsValidHexColor(color string) bool {
	if len(color) != 7 {
		return false
	}

	if !strings.HasPrefix(color, "#") {
		return false
	}

	hexCode := color[1:]

	match, _ := regexp.MatchString("^[0-9A-Fa-f]+$", hexCode)
	return match
}

func IsValidDomain(domain string) bool {
	if len(domain) == 0 {
		return false
	}

	domain = strings.ToLower(domain)
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")
	if len(domain) > 253 {
		return false
	}

	if !strings.Contains(domain, ".") {
		return false
	}

	labels := strings.Split(domain, ".")
	tld := labels[len(labels)-1]
	if len(tld) < 2 {
		return false
	}

	labelRegex := regexp.MustCompile("^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$")

	for _, label := range labels {
		if len(label) == 0 {
			return false
		}

		if len(label) > 63 {
			return false
		}

		if !labelRegex.MatchString(label) {
			return false
		}
	}

	return true
}
