package utils

func QueueFormat(n int64) string {
	if n >= 11 && n <= 19 {
		return "th"
	}

	lastDigit := n % 10

	var s string

	switch lastDigit {
	case 1:
		s = "st"
	case 2:
		s = "nd"
	case 3:
		s = "rd"
	default:
		s = "th"
	}

	return s
}
