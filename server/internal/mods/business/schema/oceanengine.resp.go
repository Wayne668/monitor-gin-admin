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
	PromotionId        int64               `json:"promotion_id"`
	PromotionName      string              `json:"promotion_name"`
	PromotionMaterials *PromotionMaterials `json:"promotion_materials"`
	AdvertiserID       int64               `json:"advertiser_id"`
	StatusFirst        string              `json:"status_first"`
	StatusSecond       []string            `json:"status_second"`
	OptStatus          string              `json:"opt_status"`
}

type PromotionMaterials struct {
	VideoMaterialList []VideoMaterialItem `json:"video_material_list"`
}

type VideoMaterialItem struct {
	MaterialID     int64  `json:"material_id"`
	MaterialStatus string `json:"material_status"`
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

type ADVideoListResp struct {
	Code      int             `json:"code"`
	Message   string          `json:"message"`
	Data      ADVideoListData `json:"data"`
	RequestID string          `json:"request_id"`
}

type ADVideoListData struct {
	List     []ADVideoItem     `json:"list"`
	PageInfo VideoListPageInfo `json:"page_info"`
}

type ADVideoItem struct {
	MaterialID int64    `json:"material_id"`
	Signature  string   `json:"signature"`
	Filename   string   `json:"filename"`
	PosterURL  string   `json:"poster_url"`
	Labels     []string `json:"labels"`
	Size       int      `json:"size"`
	Duration   float64  `json:"duration"`
	Source     string   `json:"source"`
	CreateTime string   `json:"create_time"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Url        string   `json:"url"`
}

type VideoListPageInfo struct {
	Page        int `json:"page"`
	PageSize    int `json:"page_size"`
	TotalPage   int `json:"total_page"`
	TotalNumber int `json:"total_number"`
}

// RefreshTokenResponse 刷新token响应结构
type RefreshTokenResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpireTime   int    `json:"expire_time"`
	} `json:"data"`
}
