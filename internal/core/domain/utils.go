package domain

import "math"

const (
	TWENTY_FOUR = 86200
	TWELVE      = 43200
	SIX         = 21600
	HOUR        = 3600
	THIRTY      = 1800
	FIFHTEEN    = 900

	STR_TWENTY_FOUR = "last_day"
	STR_TWELVE      = "last_twelve_hours"
	STR_SIX         = "last_six_hrs"
	STR_HOUR        = "last_hour"
	STR_THIRTY      = "last_thirthy_minutes"
	STR_FIFHTEEN    = "last_fifthteen_minutes"
)

func ToStamp(input string) int64 {
	switch input {
	case STR_FIFHTEEN:
		return FIFHTEEN
	case STR_THIRTY:
		return THIRTY
	case STR_HOUR:
		return HOUR
	case STR_SIX:
		return SIX
	case STR_TWELVE:
		return TWELVE
	default:
		return TWENTY_FOUR
	}
}

func average(input []float64) float64 {
	value := 0.0

	for _, i := range input {
		value += i
	}

	return math.Round((value/float64(len(input)))*100) / 100
}
