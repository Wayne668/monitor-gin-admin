package schema

type CustomReportResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Rows         []CustomReportRow      `json:"rows"`
		TotalMetrics map[string]interface{} `json:"total_metrics"`
		PageInfo     struct {
			Page        int `json:"page"`
			PageSize    int `json:"page_size"`
			TotalNumber int `json:"total_number"`
			TotalPage   int `json:"total_page"`
		} `json:"page_info"`
	} `json:"data"`
	RequestID string `json:"request_id"`
}

type CustomReportRowMetrics map[string]interface{}
type CustomReportRowDimensions map[string]interface{}

type CustomReportRow struct {
	Metrics    CustomReportRowMetrics    `json:"metrics"`
	Dimensions CustomReportRowDimensions `json:"dimensions"`
}

type AdvertiserSelectResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List           []int64 `json:"list"`
		CursorPageInfo struct {
			TotalNumber int64 `json:"total_number"`
			HasMore     bool  `json:"has_more"`
			Count       int   `json:"count"`
			Cursor      int64 `json:"cursor"`
		} `json:"cursor_page_info"`
	} `json:"data"`
	RequestID string `json:"request_id"`
}

type AdvertiserInfoResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccountDetailList []AccountDetail `json:"account_detail_list"`
	} `json:"data"`
	RequestID string `json:"request_id"`
}

type AccountDetail struct {
	AdvertiserID       int64  `json:"advertiser_id"`
	AdvertiserName     string `json:"advertiser_name"`
	AdvertiserStatus   string `json:"advertiser_status"`
	FirstIndustryName  string `json:"first_industry_name"`
	SecondIndustryName string `json:"second_industry_name"`
	CreateTime         string `json:"create_time"`
	AdvCompanyID       int64  `json:"adv_company_id"`
	AdvCompanyName     string `json:"adv_company_name"`
	SelfOperationTag   string `json:"self_operation_tag"`
}

type PromotionListItem struct {
	PromotionId int64  `json:"promotion_id"`
	Status      string `json:"status_first"`
}

type PromotionListResp struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	Data      PromotionListData `json:"data"`
	RequestId string            `json:"request_id"`
}

type PromotionListData struct {
	List     []PromotionListItem `json:"list"`
	PageInfo struct {
		Page        int `json:"page"`
		PageSize    int `json:"page_size"`
		TotalNumber int `json:"total_number"`
		TotalPage   int `json:"total_page"`
	} `json:"page_info"`
}
