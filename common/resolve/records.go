package resolve

import (
	"encoding/json"
	"errors"
	"net"
)

// Adapted from https://github.com/cloud66-oss/coredns_mysql

type DNSRecordContent interface {
	json.Marshaler
	json.Unmarshaler
	RecordType() RecordType
}

type ARecord struct {
	IP net.IP `json:"ip"`
}

type AAAARecord struct {
	IP net.IP `json:"ip"`
}

type CAARecord struct {
	Flag  uint8  `json:"flag"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type CNAMERecord struct {
	Host string `json:"host"`
}

type MXRecord struct {
	Host       string `json:"host"`
	Preference uint16 `json:"preference"`
}

type NSRecord struct {
	Host string `json:"host"`
}

type SOARecord struct {
	Ns      string `json:"ns"`
	MBox    string `json:"MBox"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	MinTtl  uint32 `json:"minttl"`
}

type SRVRecord struct {
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
	Port     uint16 `json:"port"`
	Target   string `json:"target"`
}

type TXTRecord struct {
	Text string `json:"text"`
}

func (*ARecord) RecordType() RecordType {
	return A
}
func (a *ARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *ARecord) UnmarshalJSON(data []byte) error {
	if a == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, a)
}
func (*AAAARecord) RecordType() RecordType {
	return AAAA
}
func (a *AAAARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AAAARecord) UnmarshalJSON(data []byte) error {
	if a == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, a)
}
func (*CAARecord) RecordType() RecordType {
	return CAA
}
func (c *CAARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CAARecord) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, c)
}
func (*CNAMERecord) RecordType() RecordType {
	return CNAME
}
func (c *CNAMERecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CNAMERecord) UnmarshalJSON(data []byte) error {
	if c == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, c)
}
func (*MXRecord) RecordType() RecordType {
	return MX
}
func (m *MXRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MXRecord) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, m)
}
func (*NSRecord) RecordType() RecordType {
	return NS
}
func (n *NSRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NSRecord) UnmarshalJSON(data []byte) error {
	if n == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, n)
}
func (*SOARecord) RecordType() RecordType {
	return SOA
}
func (s *SOARecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SOARecord) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, s)
}
func (*SRVRecord) RecordType() RecordType {
	return SRV
}
func (s *SRVRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SRVRecord) UnmarshalJSON(data []byte) error {
	if s == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, s)
}
func (*TXTRecord) RecordType() RecordType {
	return TXT
}
func (t *TXTRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}
func (t *TXTRecord) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("target is nil")
	}
	return json.Unmarshal(data, t)
}

var _ DNSRecordContent = new(ARecord)
var _ DNSRecordContent = new(AAAARecord)
var _ DNSRecordContent = new(CAARecord)
var _ DNSRecordContent = new(CNAMERecord)
var _ DNSRecordContent = new(MXRecord)
var _ DNSRecordContent = new(NSRecord)
var _ DNSRecordContent = new(SOARecord)
var _ DNSRecordContent = new(SRVRecord)
var _ DNSRecordContent = new(TXTRecord)
