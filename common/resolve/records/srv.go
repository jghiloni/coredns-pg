package records

import (
	"encoding/json"
	"errors"

	"github.com/miekg/dns"
)

type SRVRecord struct {
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}

func (*SRVRecord) RecordType() RecordType {
	return SRV
}

func (s *SRVRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SRVRecord) UnmarshalJSON(data []byte) error {
	if s == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, s)
}

func (s *SRVRecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = s.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &dns.SRV{
		Hdr:      getHeader(fqdn, recordTypes[s.RecordType()], ttl),
		Target:   s.Target,
		Port:     s.Port,
		Weight:   s.Weight,
		Priority: s.Priority,
	}, false, nil
}

func (s *SRVRecord) validate() error {
	if s.Target == "" {
		return ErrRecordNotFound
	}

	return nil
}
