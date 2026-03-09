package records

import (
	"encoding/json"
	"errors"

	"github.com/miekg/dns"
)

type MXRecord struct {
	Host       string `json:"host"`
	Preference uint16 `json:"preference"`
}

func (*MXRecord) afr() {}

func (*MXRecord) RecordType() RecordType {
	return MX
}

func (m *MXRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MXRecord) UnmarshalJSON(data []byte) error {
	if m == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, m)
}

func (m *MXRecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = m.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &dns.MX{
		Hdr:        getHeader(fqdn, recordTypes[m.RecordType()], ttl),
		Mx:         m.Host,
		Preference: m.Preference,
	}, true, nil
}

func (m *MXRecord) validate() error {
	if m.Host == "" {
		return ErrRecordNotFound
	}

	return nil
}
