package identifier

import (
	"strings"
	"unicode"
)

func normalizeCode(value string) string {
	value = strings.TrimSpace(value)
	replacer := strings.NewReplacer(" ", "", "-", "")
	return replacer.Replace(value)
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}

	for _, r := range value {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func hasGS1CheckDigit(value string) bool {
	if !isDigits(value) || len(value) < 2 {
		return false
	}

	sum := 0
	weight := 3
	for i := len(value) - 2; i >= 0; i-- {
		sum += int(value[i]-'0') * weight
		if weight == 3 {
			weight = 1
		} else {
			weight = 3
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit == int(value[len(value)-1]-'0')
}

func hasLuhnCheckDigit(value string) bool {
	if !isDigits(value) || len(value) < 2 {
		return false
	}

	sum := 0
	shouldDouble := false
	for i := len(value) - 1; i >= 0; i-- {
		digit := int(value[i] - '0')
		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}
