package records

import (
	"encoding/json"
	"errors"
	"net"

	"github.com/miekg/dns"
)

type AAAARecord struct {
	IP net.IP `json:"ip"`
}

func (*AAAARecord) RecordType() RecordType {
	return AAAA
}

func (a *AAAARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AAAARecord) UnmarshalJSON(data []byte) error {
	if a == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, a)
}

func (a *AAAARecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = a.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			err = nil
		}

		return nil, false, err
	}

	return &dns.AAAA{
		Hdr: getHeader(fqdn, recordTypes[a.RecordType()], ttl),
		AAAA: a.IP,
	}, false, nil
}

func (a *AAAARecord) validate() error {
	if a.IP == nil {
		return ErrRecordNotFound
	}

	if a.IP.To16() == nil {
		return ErrInvalidAAAARecord
	}

	return nil
}
