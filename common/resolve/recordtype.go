package resolve

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"slices"
	"strings"
)

var ErrInvalidRecordType = errors.New("unsupported resolve record type")

type RecordType string

func isValidRecordType(r RecordType) func(RecordType) bool {
	return func(s RecordType) bool {
		return strings.EqualFold(string(r), string(s))
	}
}

func (r RecordType) Value() (driver.Value, error) {
	if slices.IndexFunc(allRecordTypes, isValidRecordType(r)) == -1 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRecordType, r)
	}

	return string(r), nil
}

func (r *RecordType) Scan(src any) error {
	rt := RecordType(fmt.Sprintf("%v", src))

	if slices.IndexFunc(allRecordTypes, isValidRecordType(rt)) == -1 {
		return fmt.Errorf("%w: %s", ErrInvalidRecordType, rt)
	}

	*r = rt
	return nil
}

const (
	A     RecordType = "A"
	AAAA  RecordType = "AAAA"
	CAA   RecordType = "CAA"
	CNAME RecordType = "CNAME"
	MX    RecordType = "MX"
	NS    RecordType = "NS"
	SOA   RecordType = "SOA"
	SRV   RecordType = "SRV"
	TXT   RecordType = "TXT"
)

var allRecordTypes = []RecordType{A, AAAA, CAA, CNAME, MX, NS, SOA, SRV, TXT}
