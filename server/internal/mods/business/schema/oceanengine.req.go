package schema

type CustomReportReq struct {
	AdvertiserID int64                 `json:"advertiser_id"`
	DataTopic    string                `json:"data_topic"`
	Dimensions   []string              `json:"dimensions"`
	Metrics      []string              `json:"metrics"`
	Filters      []CustomReportFilter  `json:"filters"`
	StartTime    string                `json:"start_time"`
	EndTime      string                `json:"end_time"`
	OrderBy      []CustomReportOrderBy `json:"order_by"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
}

type CustomReportFilter struct {
	Field    string   `json:"field"`
	Type     int      `json:"type"`
	Operator int      `json:"operator"`
	Values   []string `json:"values"`
}

type CustomReportOrderBy struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

type TimeFiltering struct {
	CreateStartTime string `json:"create_start_time,omitempty"`
	CreateEndTime   string `json:"create_end_time,omitempty"`
}
