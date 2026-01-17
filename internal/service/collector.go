package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"


	"webservice/internal/core/domain"
	"webservice/internal/core/port"
	"webservice/internal/repository"
)

func RecordReportFunc(repo repository.Report) port.RecordReport {
	return func(ctx context.Context, payload []byte) error {
		var r domain.Report

		if err := json.Unmarshal(payload, &r); err != nil {
			return err
		}

		r.ReportedAt = time.Now().Unix()

		return repo.InsertReport(ctx, r)
	}
}

func CleanUpWithArchiveFunc(repo repository.Report) port.CleanUpWithArchive {
	return func(ctx context.Context) error {

		today := time.Now()

		ref := domain.FormatRef(today)

		_, err := repo.GetArchive(ctx, ref)
		if err != nil {
			if !errors.Is(err, domain.ErrNotFound{}) {
				return err
			}
		}

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
		reports, err := repo.GetReportsFromRange(ctx, tr)

		if err != nil {
			return err
		}

		if len(reports) == 0 {
			return nil
		}

		arch := domain.FormatArchive(tr, reports)

		if err := repo.InsertArchive(ctx, arch); err != nil {
			return err
		}

		return repo.DeleteReports(ctx, tr)
	}
}
