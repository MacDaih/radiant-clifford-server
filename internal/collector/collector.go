package collector

import (
	"context"
	"encoding/json"
	"net"
	"time"
	"webservice/internal/domain"
	"webservice/internal/repository"
)

type serviceCollector struct {
	repository repository.Repository
}

type Collector interface {
	ReadSock(conn net.Conn) error
	CleanUpWithArchive() error
}

func NewCollector(repository repository.Repository) Collector {
	return &serviceCollector{
		repository: repository,
	}
}

func (s *serviceCollector) takeAction(r domain.Report) error {
	if r.Light <= 50 {
		//TODO implement light control flow
        // e.g **lights on** send cmd
		func() {}()
	}
	return nil
}

func (s *serviceCollector) ReadSock(conn net.Conn) error {
	var r domain.Report
	buf := make([]byte, 128)
	rb, err := conn.Read(buf)

	if err != nil {
		return err
	}

	index := 0
	for i, v := range buf[:rb] {
		if string(v) == "}" {
			index = i + 1
		}
	}

	if err := json.Unmarshal(buf[:index], &r); err != nil {
		return err
	}

	r.ReportedAt = time.Now().Unix()

	if err := s.takeAction(r); err != nil {
		return err
	}

	return s.repository.InsertReport(context.Background(), r)
}

func (s *serviceCollector) CleanUpWithArchive() error {

	ctx := context.Background()

	today := time.Now()

	ref := domain.FormatRef(today)

	_, err := s.repository.GetArchive(ctx, ref)

	switch err.(type) {
	case domain.ErrNotFound:
		goto proceed
	default:
		return err
	}

proceed:
	prevmonth := today.Month() - 1
	if prevmonth == 0 {
		prevmonth = 12
	}

	delta := time.Now().Day() - 1
	days := domain.GetDaysOfMonth(prevmonth.String())

	from := time.Now().Add(-(time.Duration(days) * (24 * time.Hour)))
	to := time.Now().Add(-((24 * time.Hour) - time.Duration(delta)))

	tr := domain.TimeRange{
		From: from.Unix(),
		To:   to.Unix(),
	}
	reports, err := s.repository.GetReportsFromRange(ctx, tr)

	if err != nil {
		return err
	}

	if len(reports) == 0 {
		return nil
	}

	arch := domain.FormatArchive(tr, reports)

	if err := s.repository.InsertArchive(ctx, arch); err != nil {
		return err
	}

	return s.repository.DeleteReports(ctx, tr)
}
