package domain

func GetDaysOfMonth(input string) int {
	switch input {
	case "january", "march", "may", "july", "august", "october", "november", "december":
		return 31
	case "april", "june", "september":
		return 30
	default:
		return 29
	}
}
