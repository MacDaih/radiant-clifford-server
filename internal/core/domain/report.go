package domain

type Report struct {
	ReportedAt int64   `bson:"report_time" json:"time"`
	Temp       float64 `bson:"temp" json:"temp"`
	Hum        float64 `bson:"hum" json:"hum"`
	Light      float64 `bson:"ligth" json:"light"`
	Press      float64 `bson:"press" json:"press"`
}
type ReportSample struct {
	Metrics Overview
	Reports []Report
}
type Overview struct {
	TempAverage float64 `json:"temp_av"`
	HumAverage  float64 `json:"hum_av"`
	MaxTemp     float64 `json:"max_temp"`
	MinTemp     float64 `json:"min_temp"`
	MaxHum      float64 `json:"max_hum"`
	MinHum      float64 `json:"min_hum"`
}

type TimeRange struct {
	From int64
	To   int64
}

func FormatSample(reports []Report) ReportSample {
	if len(reports) == 0 {
		return ReportSample{}
	}

	var maxTemp float64 = reports[0].Temp
	var minTemp float64 = reports[0].Temp
	var maxHum float64 = reports[0].Hum
	var minHum float64 = reports[0].Hum

	for _, j := range reports {
		if maxTemp < j.Temp {
			maxTemp = j.Temp
		}
		if minTemp > j.Temp {
			minTemp = j.Temp
		}
		if maxHum < j.Hum {
			maxHum = j.Temp
		}
		if minHum > j.Hum {
			minHum = j.Hum
		}
	}

	var tempReports, humReports []float64

	for _, r := range reports {
		tempReports = append(tempReports, r.Temp)
		humReports = append(tempReports, r.Hum)
	}

	return ReportSample{
		Metrics: Overview{
			TempAverage: average(tempReports),
			HumAverage:  average(humReports),
			MaxTemp:     maxTemp,
			MinTemp:     minTemp,
			MaxHum:      maxHum,
			MinHum:      minHum,
		},
		Reports: reports,
	}
}
