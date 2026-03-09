package records

import (
	"encoding/json"
	"fmt"

	"github.com/miekg/dns"
)

type SOARecord struct {
	Ns      string `json:"ns"`
	MBox    string `json:"MBox"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	MinTtl  uint32 `json:"minttl"`
}

func (*SOARecord) RecordType() RecordType {
	return SOA
}

func (s *SOARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SOARecord) UnmarshalJSON(data []byte) error {
	if s == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, s)
}

func (s *SOARecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	ds := &dns.SOA{
		Hdr:     getHeader(fqdn, recordTypes[s.RecordType()], ttl),
		Ns:      s.Ns,
		Mbox:    s.MBox,
		Refresh: s.Refresh,
		Retry:   s.Retry,
		Expire:  s.Expire,
	}
	if s.Ns == "" {
		ds.Ns = fmt.Sprintf("%s.%s", "ns1", fqdn)
		ds.Mbox = fmt.Sprintf("%s.%s", "hostmaster", fqdn)
		ds.Refresh = 24 * oneHourInSeconds
		ds.Retry = 2 * oneHourInSeconds
		ds.Expire = oneHourInSeconds
	}

	return ds, false, nil
}

func (s *SOARecord) validate() error {
	return nil
}
