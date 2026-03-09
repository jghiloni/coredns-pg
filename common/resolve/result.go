package resolve

import (
	"github.com/miekg/dns"
)

type ResolutionResponse struct {
	Query        string
	Record       dns.RR
	ExtraRecords []dns.RR
}
