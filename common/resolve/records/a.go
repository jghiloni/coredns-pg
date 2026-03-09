package records

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/miekg/dns"
)

type ARecord struct {
	IP net.IP `json:"ip"`
}

func (*ARecord) RecordType() RecordType {
	return A
}
func (a *ARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *ARecord) UnmarshalJSON(data []byte) error {
	if a == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, a)
}

func (a *ARecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = a.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			err = nil
		}

		return nil, false, err
	}

	return &dns.A{
		Hdr: getHeader(fqdn, recordTypes[a.RecordType()], ttl),
		A:   a.IP,
	}, false, nil
}

func (a *ARecord) validate() error {
	if a.IP == nil {
		return ErrRecordNotFound
	}

	if a.IP.To4() == nil {
		return ErrInvalidARecord
	}

	return nil
}
