package records

import (
	"encoding/json"
	"errors"

	"github.com/miekg/dns"
)

type CNAMERecord struct {
	Host string `json:"host"`
}

func (*CNAMERecord) RecordType() RecordType {
	return CNAME
}

func (c *CNAMERecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CNAMERecord) UnmarshalJSON(data []byte) error {
	if c == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, c)
}

func (c *CNAMERecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = c.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &dns.CNAME{
		Hdr:    getHeader(fqdn, recordTypes[c.RecordType()], ttl),
		Target: dns.Fqdn(c.Host),
	}, false, nil
}

func (c *CNAMERecord) validate() error {
	raw := c.Host
	if raw == "" {
		return ErrRecordNotFound
	}

	if !dns.IsFqdn(raw + ".") {
		return ErrInvalidCNAMERecord
	}

	return nil
}
