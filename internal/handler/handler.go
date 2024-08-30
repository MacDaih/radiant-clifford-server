package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"webservice/internal/domain"
	"webservice/internal/repository"

	"github.com/gorilla/mux"
)

// For dev purpose only
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

type serviceHandler struct {
	repository repository.Repository
}

type Handler interface {
	GetReportsFrom(w http.ResponseWriter, r *http.Request)
	GetReportsByDate(w http.ResponseWriter, r *http.Request)
}

func NewServiceHandler(repository repository.Repository) Handler {
	return &serviceHandler{
		repository: repository,
	}
}

func (s *serviceHandler) GetReportsFrom(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	t := time.Now().Unix()

	params := mux.Vars(r)

	var rge int64
	if v, ok := params["range"]; !ok {
		log.Println("report handler err : called for reports with no time range")
		w.WriteHeader(http.StatusBadRequest)
        return
	} else {
		rge = domain.ToStamp(v)
	}

	last := t - rge
	reports, err := s.repository.GetReports(r.Context(), last)

	if err != nil {
		log.Println("report handler err : ", err)
		w.WriteHeader(http.StatusServiceUnavailable)
        return
	}

	sample := domain.FormatSample(reports)
	json.NewEncoder(w).Encode(sample)
}

func (s *serviceHandler) GetReportsByDate(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	params := mux.Vars(r)
	var trge domain.TimeRange
	if v, ok := params["date"]; !ok {
		log.Println("report handler err : called for reports with no date")
		w.WriteHeader(http.StatusBadRequest)
        return
	} else {
		p := strings.Split(v, "-")
		if len(p) != 3 {
			w.WriteHeader(http.StatusBadRequest)
		}

		var ints []int
		for _, i := range p {
			res, err := strconv.Atoi(i)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			ints = append(ints, res)
		}
		t := time.Date(ints[0], time.Month(ints[1]), ints[2], 0, 0, 0, 0, time.UTC).Unix()
		trge.From = t
		trge.To = t + domain.TWENTY_FOUR
	}
	reports, err := s.repository.GetReportsFromRange(r.Context(), trge)

	if err != nil {
		log.Println("failed fetching reports", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(reports)
}
