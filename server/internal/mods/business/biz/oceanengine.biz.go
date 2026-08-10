package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"monitor-gin-admin/pkg/logging"
	"monitor-gin-admin/pkg/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Oceanengine management for business
type Oceanengine struct {
	AgentTokenDAL *dal.AgentToken
}

// getAccessToken 根据 account_id 获取 access_token（authstatus=已授权）
func (o *Oceanengine) getAccessToken(ctx context.Context, accountID int64) (string, error) {
	token, err := o.AgentTokenDAL.GetByAccountID(ctx, strconv.FormatInt(accountID, 10))
	if err != nil {
		return "", fmt.Errorf("查询代理商token失败: %w", err)
	}
	if token == nil {
		return "", fmt.Errorf("未找到 account_id=%d 的代理商token", accountID)
	}
	if token.AuthStatus != "已授权" {
		return "", fmt.Errorf("account_id=%d 的代理商未授权", accountID)
	}
	if token.AccessToken == "" {
		return "", fmt.Errorf("account_id=%d 的 AccessToken 为空", accountID)
	}
	return token.AccessToken, nil
}

var customReportBucket = util.NewTokenBucket(100*time.Millisecond, 10)

// RefreshToken 刷新token
func (o *Oceanengine) RefreshToken(token *schema.AgentToken) (*schema.AgentToken, error) {
	// 构建请求参数
	params := url.Values{}
	params.Add("app_id", token.AppID)
	params.Add("secret", token.AppSecret)
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", token.RefreshToken)

	// 发送请求
	resp, err := http.Post(
		util.RefreshToken,
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var refreshResp schema.RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查响应状态
	if refreshResp.Code != 0 {
		return nil, fmt.Errorf("刷新token失败: %s", refreshResp.Message)
	}

	// 更新token信息
	token.AccessToken = refreshResp.Data.AccessToken
	token.RefreshToken = refreshResp.Data.RefreshToken
	token.TokenTime = int64(time.Now().Unix())

	return token, nil
}

func (o *Oceanengine) QueryCustomReport(ctx context.Context, accountID int64, req schema.CustomReportReq) (*schema.CustomReportResp, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return nil, err
	}

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
	err = util.DoGetRequestWithJsonParams(accessToken, util.APICustomReportData, params, &resp)
	if err != nil {
		logging.Context(ctx).Error("QueryCustomReport API request failed", zap.Error(err), zap.Int64("account_id", accountID), zap.Int64("advertiser_id", req.AdvertiserID))
		return nil, err
	}
	if resp.Code != 0 {
		logging.Context(ctx).Error("QueryCustomReport API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("account_id", accountID))
		return nil, fmt.Errorf("查询自定义报表失败: code=%d, message=%s", resp.Code, resp.Message)
	}
	return &resp, nil
}

func (o *Oceanengine) GetAdvertiserIDs(ctx context.Context, accountID int64, advertiserID string, filtering *schema.TimeFiltering) ([]int64, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return nil, err
	}

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
			logging.Context(ctx).Error("GetAdvertiserIDs API request failed", zap.Error(err), zap.Int64("account_id", accountID), zap.String("advertiser_id", advertiserID))
			return nil, err
		}

		if resp.Code != 0 {
			logging.Context(ctx).Error("GetAdvertiserIDs API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("account_id", accountID))
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

func (o *Oceanengine) GetAdvertiserInfo(ctx context.Context, accountID int64, accountIDs []int64) ([]schema.AccountDetail, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if len(accountIDs) == 0 {
		return nil, nil
	}

	params := url.Values{}
	idsJSON, _ := json.Marshal(accountIDs)
	params.Add("account_ids", string(idsJSON))

	var resp2 schema.AdvertiserInfoResponse
	err = util.DoGetRequest(accessToken, util.APIAdvertiserInfo, params, &resp2)
	if err != nil {
		logging.Context(ctx).Error("GetAdvertiserInfo API request failed", zap.Error(err), zap.Int64("account_id", accountID))
		return nil, err
	}

	if resp2.Code != 0 {
		logging.Context(ctx).Error("GetAdvertiserInfo API returned error", zap.Int("code", resp2.Code), zap.String("message", resp2.Message), zap.Int64("account_id", accountID))
		return nil, fmt.Errorf("获取账户详情失败: code=%d, message=%s", resp2.Code, resp2.Message)
	}

	return resp2.Data.AccountDetailList, nil
}

func (o *Oceanengine) GetRefPromotionData(ctx context.Context, agentID int64, advertiserId int64, filtering map[string]interface{}, fields []string) ([]schema.PromotionListItem, error) {
	accessToken, err := o.getAccessToken(ctx, agentID)
	if err != nil {
		return nil, err
	}

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
			logging.Context(ctx).Error("GetRefPromotionData API request failed", zap.Error(err), zap.Int64("agent_id", agentID), zap.Int64("advertiser_id", advertiserId))
			return nil, fmt.Errorf("获取单元列表失败: %w", err)
		}
		if resp.Code != 0 {
			logging.Context(ctx).Error("GetRefPromotionData API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("agent_id", agentID))
			return nil, fmt.Errorf("获取单元列表返回错误: code=%d msg=%s", resp.Code, resp.Message)
		}

		allList = append(allList, resp.Data.List...)

		if resp.Data.PageInfo.TotalPage <= page || len(resp.Data.List) < pageSize {
			break
		}
		page++
		time.Sleep(100 * time.Millisecond)
	}

	return allList, nil
}

func (o *Oceanengine) GetVideoMaterial(ctx context.Context, accountID int64, advertiserId int64, startDate, endDate string) ([]schema.MaterialVideo, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return nil, err
	}

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
			logging.Context(ctx).Error("GetVideoMaterial API request failed", zap.Error(err), zap.Int64("account_id", accountID), zap.Int64("advertiser_id", advertiserId))
			return nil, fmt.Errorf("请求AD视频列表失败: %w", err)
		}
		if resp.Code != 0 {
			logging.Context(ctx).Error("GetVideoMaterial API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("account_id", accountID))
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

// getADVideoTempUrl 获取AD视频临时URL：保存到浏览器本地缓存
func (o *Oceanengine) GetADVideoTempUrl(ctx context.Context, accountID int64, videoIds []string, advertiserId int64) (map[string]string, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	req := map[string]interface{}{
		"advertiser_id": advertiserId,
		"filtering": map[string]interface{}{
			"video_ids": videoIds,
		},
		"page":      1,
		"page_size": 100,
	}

	var resp schema.ADVideoListResp
	if err := util.DoGetRequestWithJsonParams(accessToken, util.APIFileVideoGet, req, &resp); err != nil {
		logging.Context(ctx).Error("GetADVideoTempUrl API request failed", zap.Error(err), zap.Int64("account_id", accountID), zap.Int64("advertiser_id", advertiserId))
		return nil, fmt.Errorf("请求AD视频URL失败: %w", err)
	}
	if resp.Code != 0 {
		logging.Context(ctx).Error("GetADVideoTempUrl API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("account_id", accountID))
		return nil, fmt.Errorf("AD视频URL返回错误: code=%d msg=%s", resp.Code, resp.Message)
	}
	for _, item := range resp.Data.List {
		if item.VideoURL != "" {
			result[item.VideoID] = item.VideoURL
		}
	}
	return result, nil
}

func (o *Oceanengine) UpdatePromotionStatus(ctx context.Context, accountID int64, optStatus string, promotionIDs []int64, advertiserId int64) (string, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return "", err
	}

	if len(promotionIDs) == 0 {
		return "", fmt.Errorf("营销ID列表不能为空")
	}
	// 校验操作状态，仅允许 ENABLE / DISABLE
	if optStatus != "ENABLE" && optStatus != "DISABLE" {
		return "", fmt.Errorf("opt_status 仅允许 ENABLE 或 DISABLE，当前值: %s", optStatus)
	}

	// 接口限制 data 长度 1～10，按 10 个一批分批调用
	const batchSize = 10
	var (
		lastReqID string
		batchID   = 1
	)

	for start := 0; start < len(promotionIDs); start += batchSize {
		end := start + batchSize
		if end > len(promotionIDs) {
			end = len(promotionIDs)
		}
		batch := promotionIDs[start:end]

		// 构造 data 列表
		data := make([]map[string]interface{}, 0, len(batch))
		for _, pid := range batch {
			data = append(data, map[string]interface{}{
				"promotion_id": pid,
				"opt_status":   optStatus,
			})
		}

		// 构造请求体
		reqBody := map[string]interface{}{
			"advertiser_id": advertiserId,
			"data":          data,
		}

		var resp struct {
			Code      int    `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
			Data      struct {
				PromotionIDs []int64 `json:"promotion_ids"`
				Errors       []struct {
					PromotionID  int64  `json:"promotion_id"`
					ErrorMessage string `json:"error_message"`
				} `json:"errors"`
			} `json:"data"`
		}

		if err := util.DoPostJSONRequest(accessToken, util.APIPromotionStatusUpdate, reqBody, &resp); err != nil {
			logging.Context(ctx).Error("UpdatePromotionStatus API request failed", zap.Error(err), zap.Int64("account_id", accountID), zap.Int64("advertiser_id", advertiserId))
			return lastReqID, fmt.Errorf("第 %d 批更新营销状态请求失败: %w", batchID, err)
		}
		lastReqID = resp.RequestID

		if resp.Code != 0 {
			logging.Context(ctx).Error("UpdatePromotionStatus API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("account_id", accountID))
			return lastReqID, fmt.Errorf("第 %d 批更新营销状态失败: code=%d message=%s", batchID, resp.Code, resp.Message)
		}

		// 记录失败的营销，便于上层定位
		if len(resp.Data.Errors) > 0 {
			failedIDs := make([]int64, 0, len(resp.Data.Errors))
			for _, e := range resp.Data.Errors {
				failedIDs = append(failedIDs, e.PromotionID)
			}
			return lastReqID, fmt.Errorf("第 %d 批部分营销更新失败: %v", batchID, failedIDs)
		}

		batchID++
		// 批次之间小睡，避免触发频控
		if start+batchSize < len(promotionIDs) {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return lastReqID, nil
}

func (o *Oceanengine) UpdateMaterialStatus(ctx context.Context, accountID int64, optStatus string, materialIDs []int64, promotionID, advertiserId int64) (string, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return "", err
	}

	if len(materialIDs) == 0 {
		return "", fmt.Errorf("素材ID列表不能为空")
	}
	// 校验操作状态，仅允许 DISABLE / ENABLE
	if optStatus != "DISABLE" && optStatus != "ENABLE" {
		return "", fmt.Errorf("opt_status 仅允许 DISABLE 或 ENABLE，当前值: %s", optStatus)
	}

	// 接口限制 data 长度 1～10，按 10 个一批分批调用
	const batchSize = 10
	var (
		lastReqID string
		batchID   = 1
	)

	for start := 0; start < len(materialIDs); start += batchSize {
		end := start + batchSize
		if end > len(materialIDs) {
			end = len(materialIDs)
		}
		batch := materialIDs[start:end]

		// 构造 data 列表
		data := make([]map[string]interface{}, 0, len(batch))
		for _, mid := range batch {
			data = append(data, map[string]interface{}{
				"material_id": mid,
				"opt_status":  optStatus,
			})
		}

		// 构造请求体：advertiser_id 使用字符串形式（与官方 Python 示例保持一致）
		reqBody := map[string]interface{}{
			"advertiser_id": strconv.FormatInt(advertiserId, 10),
			"promotion_id":  strconv.FormatInt(promotionID, 10),
			"data":          data,
		}

		var resp struct {
			Code      int    `json:"code"`
			Message   string `json:"message"`
			RequestID string `json:"request_id"`
			Data      struct {
				PromotionID int64   `json:"promotion_id"`
				MaterialIDs []int64 `json:"material_ids"`
				Errors      []struct {
					MaterialID int64  `json:"material_id"`
					Message    string `json:"message"`
				} `json:"errors"`
			} `json:"data"`
		}

		if err := util.DoPostJSONRequest(accessToken, util.APIPromotionMaterialStatusUpdate, reqBody, &resp); err != nil {
			logging.Context(ctx).Error("UpdateMaterialStatus API request failed", zap.Error(err), zap.Int64("account_id", accountID), zap.Int64("advertiser_id", advertiserId))
			return lastReqID, fmt.Errorf("第 %d 批更新素材状态请求失败: %w", batchID, err)
		}
		lastReqID = resp.RequestID

		if resp.Code != 0 {
			logging.Context(ctx).Error("UpdateMaterialStatus API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("account_id", accountID))
			return lastReqID, fmt.Errorf("第 %d 批更新素材状态失败: code=%d message=%s", batchID, resp.Code, resp.Message)
		}

		// 记录失败的素材，便于上层定位
		if len(resp.Data.Errors) > 0 {
			failedIDs := make([]int64, 0, len(resp.Data.Errors))
			for _, e := range resp.Data.Errors {
				failedIDs = append(failedIDs, e.MaterialID)
			}
			return lastReqID, fmt.Errorf("第 %d 批部分素材更新失败: %v", batchID, failedIDs)
		}

		batchID++
		// 批次之间小睡，避免触发频控
		if start+batchSize < len(materialIDs) {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return lastReqID, nil
}

func (o *Oceanengine) DeleteMaterialUnderPromotion(ctx context.Context, accountID int64, materialId, promotionID, advertiserId int64) (string, error) {
	accessToken, err := o.getAccessToken(ctx, accountID)
	if err != nil {
		return "", err
	}

	// 构造请求体：advertiser_id / promotion_id 使用字符串形式（与官方 Python 示例保持一致）
	// material_id 为 number[]，接口仅支持传入 1 个素材
	reqBody := map[string]interface{}{
		"advertiser_id": strconv.FormatInt(advertiserId, 10),
		"promotion_id":  strconv.FormatInt(promotionID, 10),
		"material_id":   []int64{materialId},
	}

	var resp struct {
		Code      int      `json:"code"`
		Message   string   `json:"message"`
		RequestID string   `json:"request_id"`
		Data      struct{} `json:"data"`
	}

	if err := util.DoPostJSONRequest(accessToken, util.APIPromotionMaterialDelete, reqBody, &resp); err != nil {
		logging.Context(ctx).Error("DeleteMaterialUnderPromotion API request failed", zap.Error(err), zap.Int64("account_id", accountID), zap.Int64("advertiser_id", advertiserId), zap.Int64("material_id", materialId))
		return "", fmt.Errorf("删除营销下素材请求失败: %w", err)
	}

	if resp.Code != 0 {
		logging.Context(ctx).Error("DeleteMaterialUnderPromotion API returned error", zap.Int("code", resp.Code), zap.String("message", resp.Message), zap.Int64("account_id", accountID))
		return resp.RequestID, fmt.Errorf("删除营销下素材失败: code=%d message=%s", resp.Code, resp.Message)
	}

	return resp.RequestID, nil
}
