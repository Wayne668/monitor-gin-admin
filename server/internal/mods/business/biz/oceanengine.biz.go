package biz

import (
	"encoding/json"
	"fmt"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/util"
	"net/url"
	"strconv"
	"time"
)

// Oceanengine management for business
type Oceanengine struct{}

var customReportBucket = util.NewTokenBucket(100*time.Millisecond, 10)

func QueryCustomReport(accessToken string, req schema.CustomReportReq) (*schema.CustomReportResp, error) {
	customReportBucket.Take()

	var resp schema.CustomReportResp
	params := map[string]interface{}{
		"advertiser_id": strconv.FormatInt(req.AdvertiserID, 10),
		"data_topic":    req.DataTopic,
		"dimensions":    req.Dimensions,
		"metrics":       req.Metrics,
		"filters":       req.Filters,
		"start_time":    req.StartTime,
		"end_time":      req.EndTime,
		"order_by":      req.OrderBy,
		"page":          strconv.Itoa(req.Page),
		"page_size":     strconv.Itoa(req.PageSize),
	}
	err := util.DoGetRequestWithJsonParams(accessToken, util.APICustomReportData, params, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("查询自定义报表失败: code=%d, message=%s", resp.Code, resp.Message)
	}
	return &resp, nil
}

func GetAdvertiserIDs(accessToken, advertiserID string, filtering *schema.TimeFiltering) ([]int64, error) {
	allIDs := make([]int64, 0)
	count := 100
	var cursor int64

	for {
		params := url.Values{}
		params.Add("advertiser_id", advertiserID)
		params.Add("count", fmt.Sprintf("%d", count))
		if cursor > 0 {
			params.Add("cursor", fmt.Sprintf("%d", cursor))
		}

		if filtering != nil {
			filteringJSON, _ := json.Marshal(filtering)
			params.Add("filtering", string(filteringJSON))
		}

		var resp schema.AdvertiserSelectResponse
		err := util.DoGetRequest(accessToken, util.APIAdvertiserSelect, params, &resp)
		if err != nil {
			return nil, err
		}

		if resp.Code != 0 {
			return nil, fmt.Errorf("code=%d, message=%s", resp.Code, resp.Message)
		}

		allIDs = append(allIDs, resp.Data.List...)

		if !resp.Data.CursorPageInfo.HasMore {
			break
		}

		cursor = resp.Data.CursorPageInfo.Cursor
		time.Sleep(100 * time.Millisecond)
	}

	return allIDs, nil
}

func GetAdvertiserInfo(accessToken string, accountIDs []int64) ([]schema.AccountDetail, error) {
	if len(accountIDs) == 0 {
		return nil, nil
	}

	params := url.Values{}
	idsJSON, _ := json.Marshal(accountIDs)
	params.Add("account_ids", string(idsJSON))

	var resp schema.AdvertiserInfoResponse
	err := util.DoGetRequest(accessToken, util.APIAdvertiserInfo, params, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("获取账户详情失败: code=%d, message=%s", resp.Code, resp.Message)
	}

	return resp.Data.AccountDetailList, nil
}
