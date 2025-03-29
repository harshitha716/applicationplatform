package helper

func IsValidShortInput(input string) bool {
	return len(input) > 0 && len(input) <= 24
}

func IsValidMediumInput(input string) bool {
	return len(input) > 0 && len(input) <= 64
}

func IsValidLongInput(input string) bool {
	return len(input) > 0 && len(input) <= 255
}

func IsValidHexCode(input string) bool {
	if len(input) != 7 || input[0] != '#' {
		return false
	}

	for i := 1; i < len(input); i++ {
		c := input[i]
		isHexDigit := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !isHexDigit {
			return false
		}
	}

	return true
}
