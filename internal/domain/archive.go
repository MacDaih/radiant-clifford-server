package domain

import (
	"fmt"
	"time"
)

type Archive struct {
	Ref    string    `json:"ref" bson:"ref"`
	Period TimeRange `json:"period" bson:"period"`
	Report `json:"report" bson:"report"`
}

func FormatArchive(archiveRange TimeRange, reports []Report) Archive {

	var temps, humids, pressures []float64

	for _, r := range reports {
		temps = append(temps, r.Temp)
		humids = append(humids, r.Hum)
		pressures = append(pressures, r.Press)
	}

	from := time.Unix(archiveRange.From, 0)

	label := FormatRef(from)

	return Archive{
		Ref:    label,
		Period: archiveRange,
		Report: Report{
			ReportedAt: archiveRange.From,
			Temp:       average(temps),
			Hum:        average(humids),
			Press:      average(pressures),
		},
	}
}

func FormatRef(t time.Time) string {
	return fmt.Sprintf("%s_%d", t.Month().String(), t.Year())
}
