package repository

import (
	"context"
	"fmt"
	"log"
	"webservice/internal/domain"
	"webservice/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	reportCollection  = "report"
	archiveCollection = "archive"
)

type Repository interface {
	GetReports(context.Context, int64) ([]domain.Report, error)
	GetReportsFromRange(context.Context, domain.TimeRange) ([]domain.Report, error)
	InsertReport(context.Context, domain.Report) error
	DeleteReports(ctx context.Context, rge domain.TimeRange) error

	InsertArchive(ctx context.Context, archive domain.Archive) error
	GetArchive(ctx context.Context, ref string) (domain.Archive, error)
}

type reportsRepo struct {
	dbname string
	dbHost string
	dbPort string
}

func NewReportRepository(name, dbHost, dbPort string) Repository {
	return &reportsRepo{
		dbname: name,
		dbHost: dbHost,
		dbPort: dbPort,
	}
}

func (r *reportsRepo) GetReports(ctx context.Context, elapse int64) ([]domain.Report, error) {
	client, err := database.ConnectDB(r.dbHost, r.dbPort)

	if err != nil {
		return nil, err
	}

	defer client.Disconnect(ctx)

	var reports []domain.Report

	filter := bson.M{"report_time": bson.M{"$gte": elapse}}

	coll := client.Database(r.dbname).Collection(reportCollection)

	res, err := coll.Find(ctx, filter)
	if err != nil {
		log.Println("read err : ", err)
		return nil, err
	}

	defer res.Close(ctx)

	for res.Next(ctx) {
		var r domain.Report
		if err = res.Decode(&r); err != nil {
			log.Println("decoding err ", err)
			continue
		}

		reports = append(reports, r)
	}

	return reports, err
}

func (r *reportsRepo) GetReportsFromRange(ctx context.Context, rge domain.TimeRange) ([]domain.Report, error) {
	client, err := database.ConnectDB(r.dbHost, r.dbPort)
	if err != nil {
		return nil, err
	}

	defer client.Disconnect(ctx)

	var reports []domain.Report

	filter := bson.M{"report_time": bson.M{"$gte": rge.From, "$lte": rge.To}}

	coll := client.Database(r.dbname).Collection(reportCollection)

	res, err := coll.Find(ctx, filter)
	if err != nil {
		log.Println("read err : ", err)
		return nil, err
	}

	defer res.Close(ctx)

	for res.Next(ctx) {
		var r domain.Report
		if err = res.Decode(&r); err != nil {
			log.Println("decoding err ", err)
			continue
		}

		reports = append(reports, r)
	}
	return reports, err
}

func (r *reportsRepo) InsertReport(ctx context.Context, report domain.Report) (err error) {
	return database.Write(ctx, r.dbname, reportCollection, report)
}

func (r *reportsRepo) InsertArchive(ctx context.Context, archive domain.Archive) error {
	return database.Write(ctx, r.dbname, archiveCollection, archive)
}

func (r *reportsRepo) DeleteReports(ctx context.Context, rge domain.TimeRange) error {

	client, err := database.ConnectDB(r.dbHost, r.dbPort)
	if err != nil {
		return err
	}

	filter := bson.M{"report_time": bson.M{"$gte": rge.From, "$lte": rge.To}}
	coll := client.Database(r.dbname).Collection(reportCollection)

	if _, err := coll.DeleteMany(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (r *reportsRepo) GetArchive(ctx context.Context, ref string) (domain.Archive, error) {
	client, err := database.ConnectDB(r.dbHost, r.dbPort)

	if err != nil {
		return domain.Archive{}, err
	}

	coll := client.Database(r.dbname).Collection(archiveCollection)

	filter := bson.M{"ref": ref}

	res := coll.FindOne(ctx, filter)

	if res.Err() != nil {
		return domain.Archive{}, domain.ErrNotFound{
			Msg: fmt.Sprintf("archive with ref. %s not found", ref),
		}
	}

	var arch domain.Archive
	if err := res.Decode(&arch); err != nil {
		return domain.Archive{}, err
	}

	return arch, nil
}
