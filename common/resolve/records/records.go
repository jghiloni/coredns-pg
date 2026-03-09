package records

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/miekg/dns"
)

var (
	ErrCannotUnmarshalContents = errors.New("cannot unmarshal contents; target is nil")
	ErrRecordNotFound          = errors.New("record not found")
	ErrInvalidARecord          = errors.New("dns A records must contain IPv4 addresses")
	ErrInvalidAAAARecord       = errors.New("dns AAAA records must contain IPv6 addresses")
	ErrInvalidCNAMERecord      = errors.New("dns CNAME records must be fully qualified domain names")
)

var recordTypes = map[RecordType]uint16{
	A:     dns.TypeA,
	AAAA:  dns.TypeAAAA,
	CAA:   dns.TypeCAA,
	CNAME: dns.TypeCNAME,
	MX:    dns.TypeMX,
	NS:    dns.TypeNS,
	SOA:   dns.TypeSOA,
	SRV:   dns.TypeSRV,
	TXT:   dns.TypeTXT,
}

const (
	minTTL           = 30
	defaultTTL       = 1800
	oneHourInSeconds = 3600
)

type DNSRecord struct {
	Name    string
	Zone    string
	TTL     uint32
	Content DNSRecordContent
	Type    RecordType
}

type DNSRecordContent interface {
	json.Marshaler
	json.Unmarshaler
	RecordType() RecordType
	partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error)
	validate() error
}

func (d *DNSRecord) AsResourceRecord() (dns.RR, bool, error) {
	fqdn := d.fqdn()
	return d.Content.partialRecord(fqdn, d.GetTTL())
}

func (d *DNSRecord) GetTTL() uint32 {
	if d.TTL < minTTL {
		d.TTL = defaultTTL
	}

	return d.TTL
}

func (d *DNSRecord) fqdn() string {
	if d.Name == "" || d.Name == "@" { // apex records
		return dns.Fqdn(d.Zone)
	}

	return dns.Fqdn(fmt.Sprintf("%s.%s", d.Name, d.Zone))
}

func getHeader(fqdn string, rrtype uint16, ttl uint32) dns.RR_Header {
	return dns.RR_Header{
		Name:   fqdn,
		Rrtype: rrtype,
		Class:  dns.ClassINET,
		Ttl:    ttl,
	}
}
