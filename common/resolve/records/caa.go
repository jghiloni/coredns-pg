package records

import (
	"encoding/json"
	"errors"

	"github.com/miekg/dns"
)

type CAARecord struct {
	Flag  uint8  `json:"flag"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (*CAARecord) RecordType() RecordType {
	return CAA
}

func (c *CAARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CAARecord) UnmarshalJSON(data []byte) error {
	if c == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, c)
}

func (c *CAARecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = c.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			err = nil
		}

		return nil, false, err
	}

	return &dns.CAA{
		Hdr: getHeader(fqdn, recordTypes[c.RecordType()], ttl),
		Value: c.Value,
		Tag:   c.Tag,
		Flag:  c.Flag,
	}, false, nil
}

func (c *CAARecord) validate() error {
	if c.Value == "" || c.Tag == "" {
		return ErrRecordNotFound
	}

	return nil
}
