package records

import (
	"encoding/json"
	"errors"

	"github.com/miekg/dns"
)

type NSRecord struct {
	Host string `json:"host"`
}

func (*NSRecord) afr() {}

func (*NSRecord) RecordType() RecordType {
	return NS
}

func (n *NSRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NSRecord) UnmarshalJSON(data []byte) error {
	if n == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, n)
}

func (n *NSRecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = n.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &dns.NS{
		Hdr: getHeader(fqdn, recordTypes[n.RecordType()], ttl),
		Ns:  n.Host,
	}, true, nil
}

func (n *NSRecord) validate() error {
	if n.Host == "" {
		return ErrRecordNotFound
	}

	return nil
}
