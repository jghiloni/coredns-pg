package records

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/miekg/dns"
)

type TXTRecord struct {
	Text string `json:"text"`
}

func (*TXTRecord) RecordType() RecordType {
	return TXT
}

func (t *TXTRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TXTRecord) UnmarshalJSON(data []byte) error {
	if t == nil {
		return ErrCannotUnmarshalContents
	}

	return json.Unmarshal(data, t)
}
func (t *TXTRecord) partialRecord(fqdn string, ttl uint32) (record dns.RR, fetchExtras bool, err error) {
	if err = t.validate(); err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &dns.TXT{
		Hdr: getHeader(fqdn, recordTypes[t.RecordType()], ttl),
		Txt: t.wrapText(),
	}, false, nil
}

func (t *TXTRecord) validate() error {
	if t.Text == "" {
		return ErrRecordNotFound
	}

	return nil
}

func (t *TXTRecord) wrapText() []string {
	var lines []string

	lineLen := 255
	reader := strings.NewReader(t.Text)

	line := make([]byte, lineLen)

	var (
		readLen int
		err     error
	)
	for readLen, err = reader.Read(line); readLen == lineLen && err != nil; {
		lines = append(lines, string(line))
	}

	if err != nil {
		return lines
	}

	return append(lines, string(line))
}
