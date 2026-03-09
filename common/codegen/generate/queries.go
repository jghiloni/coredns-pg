package generate

import (
	"time"

	"gorm.io/gen"

	"github.com/jghiloni/coredns-pg/common/resolve/records"
)

type DNSRecordQueries interface {
	// GetRecentlyDeleted
	//
	// SELECT * FROM @@table WHERE deleted_at IS NOT NULL AND deleted_at > @oldest ORDER BY deleted_at DESC
	GetRecentlyDeleted(oldest time.Time) ([]gen.T, error)

	// ResolveRequest
	//
	// SELECT r.* FROM records r INNER JOIN zones z ON r.zone = z.fqdn WHERE (
	// 	r.record_type = @recordType AND r.deleted_at IS NULL AND z.deleted_at IS NULL AND
	// 	(
	// 		(r.name = '@' AND r.zone = @request) OR (@request LIKE REPLACE(r.name, '*', '%') || '.' || r.zone)
	// 	)
	// ) LIMIT 1
	//
	ResolveRequest(request string, recordType records.RecordType) (gen.T, error)

	// ResolveRequests
	//
	// SELECT r.* FROM records r INNER JOIN zones z ON r.zone = z.fqdn WHERE (
	// 	r.record_type IN @recordTypes AND r.deleted_at IS NULL AND z.deleted_at IS NULL AND
	// 	(
	// 		(r.name = '@' AND r.zone = @request) OR (@request LIKE REPLACE(r.name, '*', '%') || '.' || r.zone)
	// 	)
	// ) LIMIT 1
	//
	ResolveRequests(request string, recordTypes ...records.RecordType) ([]gen.T, error)


	// GetZoneRecords
	//
	// SELECT r.* FROM records r INNER JOIN zones z ON r.zone = z.fqdn WHERE (
	//	r.zone = @zone AND r.deleted_at IS NULL AND z.deleted_at IS NULL
	// )
	GetZoneRecords(zone string) ([]gen.T, error)
}

type DNSZoneQueries interface {
	// GetRecentlyDeleted
	//
	// SELECT * FROM @@table WHERE deleted_at IS NOT NULL AND deleted_at > @oldest ORDER BY deleted_at DESC
	GetRecentlyDeleted(oldest time.Time) ([]gen.T, error)
}
