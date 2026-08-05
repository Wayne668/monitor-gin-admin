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

func (o *Oceanengine) QueryCustomReport(accessToken string, req schema.CustomReportReq) (*schema.CustomReportResp, error) {
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

func (o *Oceanengine) GetAdvertiserIDs(accessToken, advertiserID string, filtering *schema.TimeFiltering) ([]int64, error) {
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

func (o *Oceanengine) GetAdvertiserInfo(accessToken string, accountIDs []int64) ([]schema.AccountDetail, error) {
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

func (o *Oceanengine) GetRefPromotionData(accessToken string, advertiserId int64, filtering map[string]interface{}, fields []string) ([]schema.PromotionListItem, error) {
	allList := make([]schema.PromotionListItem, 0)
	page := 1
	pageSize := 20

	for {
		var resp schema.PromotionListResp
		params := map[string]interface{}{
			"advertiser_id": advertiserId,
			"filtering":     filtering,
			"fields":        fields,
			"page":          page,
			"page_size":     pageSize,
		}
		err := util.DoGetRequestWithJsonParams(accessToken, util.APIPromotionListGet, params, &resp)
		if err != nil {
			return nil, fmt.Errorf("获取单元列表失败: %w", err)
		}
		if resp.Code != 0 {
			return nil, fmt.Errorf("获取单元列表返回错误: code=%d msg=%s", resp.Code, resp.Message)
		}

		allList = append(allList, resp.Data.List...)

		if resp.Data.PageInfo.TotalPage <= page || len(resp.Data.List) < pageSize {
			break
		}
		page++
		time.Sleep(100 * time.Millisecond)
	}

	if len(allList) == 0 {
		return nil, fmt.Errorf("没有单元")
	}
	return allList, nil
}

func (o *Oceanengine) GetVideoMaterial(accessToken string, advertiserId int64, startDate, endDate string) ([]schema.MaterialVideo, error) {
	pageSize := 50
	req := map[string]interface{}{
		"advertiser_id": advertiserId,
		"filtering": map[string]string{
			"start_time": startDate,
			"end_time":   endDate,
		},
		"page_size": pageSize,
	}

	page := 1
	record := make([]schema.MaterialVideo, 0)
	for {
		req["page"] = page
		var resp schema.ADVideoListResp
		if err := util.DoGetRequestWithJsonParams(accessToken, util.APIFileVideoGet, req, &resp); err != nil {
			return nil, fmt.Errorf("请求AD视频列表失败: %w", err)
		}
		if resp.Code != 0 {
			return nil, fmt.Errorf("AD视频列表返回错误: code=%d msg=%s", resp.Code, resp.Message)
		}

		for _, item := range resp.Data.List {
			labels := ""
			if len(item.Labels) > 0 {
				for i, p := range item.Labels {
					if i > 0 {
						labels += ","
					}
					labels += p
				}
			}

			info := schema.MaterialVideo{
				VideoID:      item.VideoID,
				AdvertiserID: advertiserId,
				MaterialID:   item.MaterialID,
				Signature:    item.Signature,
				FileName:     item.FileName,
				PosterURL:    item.PosterURL,
				Labels:       labels,
			}
			if item.CreateTime != "" {
				if t, err := time.Parse("2006-01-02 15:04:05", item.CreateTime); err == nil {
					info.CreatedAt = t
				}
			}

			record = append(record, info)
		}

		if page*pageSize >= resp.Data.PageInfo.TotalNumber || len(resp.Data.List) == 0 {
			break
		}
		page++
		time.Sleep(150 * time.Millisecond)
	}
	return record, nil
}
