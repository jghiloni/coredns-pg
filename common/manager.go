package common

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jghiloni/coredns-pg/common/config"
	"github.com/jghiloni/coredns-pg/common/dns"
	"github.com/jghiloni/coredns-pg/common/generated/db"
	"github.com/jghiloni/coredns-pg/common/generated/tables"
	"gorm.io/gorm"
)

var (
	validateRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?$`)
)

var (
	ErrRecordTypeMismatch = errors.New("unexpected record type")
	ErrMissingTrailingDot = errors.New("requests must end with a trailing dot")
	ErrNameSegmentInvalid = errors.New("every segment of a name must only include ascii letters, numbers, and dashes")
	ErrNameSegmentTooLong = errors.New("a name segment must not be longer than 63 characters")
	ErrNameSegmentEmpty   = errors.New("a name may not have consecutive dots")
	ErrCouldNotLookupZone = errors.New("could not lookup zone data")
	ErrZoneNotFound       = errors.New("zone not found or has been deleted")
)

type UpdateRecordRequest struct {
	NewName     string
	NewTTL      uint
	NewContents dns.DNSRecordContent
}

type DNSQuerier interface {
	GetRecord(ctx context.Context, id string) (tables.Record, error)
	IsZoneValid(ctx context.Context, fqdn string) (bool, error)
	ResolveRequest(ctx context.Context, request string, recordType dns.RecordType) (dns.DNSRecordContent, error)
	ListZoneRecords(ctx context.Context, zoneFQDN string) ([]tables.Record, error)
	ListZones(ctx context.Context) ([]tables.Zone, error)
	ListRecentlyDeletedRecords(ctx context.Context, since time.Duration) ([]tables.Record, error)
	ListRecentlyDeletedZones(ctx context.Context, since time.Duration) ([]tables.Zone, error)
}

type DNSManager interface {
	DNSQuerier
	CreateZone(ctx context.Context, fqdn string) error
	DeleteZone(ctx context.Context, fqdn string) error
	CreateRecord(ctx context.Context, name string, zone string, ttl uint, recordType dns.RecordType, contents dns.DNSRecordContent) (tables.Record, error)
	UpdateRecord(ctx context.Context, id string, record UpdateRecordRequest) error
	DeleteRecord(ctx context.Context, id string) error
}

type dnsOrmManager struct{}

func NewDNSQuerier(cfg config.DatabaseConfig) (DNSQuerier, error) {
	return newORMManager(cfg)
}

func NewDNSManager(cfg config.DatabaseConfig) (DNSManager, error) {
	return newORMManager(cfg)
}

func newORMManager(cfg config.DatabaseConfig) (*dnsOrmManager, error) {
	gdb, err := cfg.OpenDB()
	if err != nil {
		return nil, err
	}

	db.SetDefault(gdb)

	return &dnsOrmManager{}, nil
}

func (d *dnsOrmManager) GetRecord(ctx context.Context, id string) (tables.Record, error) {
	return gorm.G[tables.Record](db.Q.UnderlyingDB()).Where("id = ? AND deleted_at is null", id).First(ctx)
}

func (d *dnsOrmManager) ResolveRequest(ctx context.Context, request string, recordType dns.RecordType) (dns.DNSRecordContent, error) {
	record, err := db.Record.WithContext(ctx).ResolveRequest(request, recordType)
	if err != nil {
		return nil, err
	}

	content := record.Content
	if content.RecordType() != recordType {
		return nil, fmt.Errorf("%w: expected type %s, got %s", ErrRecordTypeMismatch, recordType, content.RecordType())
	}

	return content, nil
}

func (d *dnsOrmManager) ListZoneRecords(ctx context.Context, zoneFQDN string) ([]tables.Record, error) {
	if err := d.validateName(zoneFQDN); err != nil {
		return nil, err
	}
	return db.Record.WithContext(ctx).GetZoneRecords(zoneFQDN)
}

func (d *dnsOrmManager) IsZoneValid(ctx context.Context, fqdn string) (bool, error) {
	z, err := db.Zone.WithContext(ctx).Where(db.Zone.Fqdn.Eq(fqdn), db.Zone.DeletedAt.IsNotNull()).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		z = nil
	}

	return z != nil, err
}

func (d *dnsOrmManager) ListZones(ctx context.Context) ([]tables.Zone, error) {
	zonePtrs, err := db.Zone.WithContext(ctx).Where(db.Zone.DeletedAt.IsNotNull()).Find()
	var zones []tables.Zone
	for _, p := range zonePtrs {
		if p != nil {
			zones = append(zones, *p)
		}
	}

	return zones, err
}

func (d *dnsOrmManager) ListRecentlyDeletedRecords(ctx context.Context, since time.Duration) ([]tables.Record, error) {
	return db.Record.WithContext(ctx).GetRecentlyDeleted(time.Now().Add(-since))
}

func (d *dnsOrmManager) ListRecentlyDeletedZones(ctx context.Context, since time.Duration) ([]tables.Zone, error) {
	return db.Zone.WithContext(ctx).GetRecentlyDeleted(time.Now().Add(-since))

}

func (d *dnsOrmManager) CreateZone(ctx context.Context, fqdn string) error {
	if err := d.validateName(fqdn); err != nil {
		return err
	}

	return db.Zone.WithContext(ctx).Create(&tables.Zone{Fqdn: fqdn})
}

func (d *dnsOrmManager) DeleteZone(ctx context.Context, fqdn string) error {
	_, err := db.Zone.WithContext(ctx).Delete(&tables.Zone{Fqdn: fqdn})
	return err
}

func (d *dnsOrmManager) CreateRecord(ctx context.Context, name string, zone string, ttl uint, recordType dns.RecordType, contents dns.DNSRecordContent) (tables.Record, error) {
	isValid, err := d.IsZoneValid(ctx, zone)
	if !isValid {
		if err != nil {
			return tables.Record{}, fmt.Errorf("%w: %w", ErrCouldNotLookupZone, err)
		}

		return tables.Record{}, fmt.Errorf("%w: %s", ErrZoneNotFound, zone)
	}

	if err = d.validateName(fmt.Sprintf("%s.%s", name, zone)); err != nil {
		return tables.Record{}, err
	}

	if recordType != contents.RecordType() {
		return tables.Record{}, fmt.Errorf("%w: request was for a %s record but contents were for a %s record", ErrRecordTypeMismatch, recordType, contents.RecordType())
	}

	record := tables.Record{
		Name:       name,
		Zone:       zone,
		TTL:        uint(ttl),
		Content:    contents,
		RecordType: recordType,
	}

	err = db.Record.WithContext(ctx).Create(&record)
	return record, err
}

func (d *dnsOrmManager) UpdateRecord(ctx context.Context, id string, request UpdateRecordRequest) error {
	existingRecord, err := d.GetRecord(ctx, id)
	if err != nil {
		return err
	}

	if request.NewName != "" {
		if err = d.validateName(fmt.Sprintf("%s.%s", request.NewName, existingRecord.Zone)); err != nil {
			return err
		}

		existingRecord.Name = request.NewName
	}

	if request.NewTTL > 0 {
		existingRecord.TTL = request.NewTTL
	}

	if request.NewContents != nil {
		if request.NewContents.RecordType() != existingRecord.RecordType {
			return fmt.Errorf("%w: record is of type %s but the new contents are type %s", ErrRecordTypeMismatch, existingRecord.RecordType, request.NewContents.RecordType())
		}

		existingRecord.Content = request.NewContents
	}

	return db.Record.WithContext(ctx).Save(&existingRecord)
}

func (d *dnsOrmManager) DeleteRecord(ctx context.Context, id string) error {
	_, err := db.Record.WithContext(ctx).Delete(&tables.Record{ID: id})
	return err
}

func (d *dnsOrmManager) validateName(fqdn string) error {
	var errs []error

	if !strings.HasSuffix(fqdn, ".") {
		errs = append(errs, ErrMissingTrailingDot)
	}

	parts := strings.Split(strings.TrimSuffix(fqdn, "."), ".")
	for i, part := range parts {
		switch {
		case len(part) > 63:
			errs = append(errs, fmt.Errorf("%w: %s is %d characters", ErrNameSegmentTooLong, part, len(part)))
		case len(part) == 0:
			errs = append(errs, fmt.Errorf("%w: segment %d is empty", ErrNameSegmentEmpty, i))
		case !validateRegex.MatchString(part):
			errs = append(errs, fmt.Errorf("%w: %s is invalid", ErrNameSegmentInvalid, part))
		}
	}

	return errors.Join(errs...)
}
