package resolve

type ResolutionResponse struct {
	Query   string           `json:"query"`
	Content DNSRecordContent `json:"contents"`
	TTL     uint32           `json:"ttl"`
}
