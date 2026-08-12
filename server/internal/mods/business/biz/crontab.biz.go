package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
)

type Crontab struct {
	AccountInfo       *dal.AccountInfo
	Oceanengine       *Oceanengine
	AgentToken        *dal.AgentToken
	HostRule          *dal.HostRule
	HostAccount       *dal.HostAccount
	PromotionMaterial *dal.PromotionMaterial
	MaterialVideo     *dal.MaterialVideo
}

func (s *Crontab) SyncAccounts(ctx context.Context, accountID int64, StartDate, EndDate string) error {
	advertiserID := ""
	if advertiserID == "" {
		return fmt.Errorf("remarks=%s 未找到已授权的Token", "")
	}

	filtering := &schema.TimeFiltering{
		CreateStartTime: StartDate + " 00:00:00",
		CreateEndTime:   EndDate + " 23:59:59",
	}

	advertiserIDs, err := s.Oceanengine.GetAdvertiserIDs(ctx, accountID, advertiserID, filtering)
	if err != nil {
		return fmt.Errorf("获取账户ID列表失败: %w", err)
	}

	fmt.Printf("【SyncAccounts】API返回 %d 个账户ID", len(advertiserIDs))

	if len(advertiserIDs) == 0 {
		return nil
	}

	// 调用model过滤已存在账户ID
	newAdvertiserIDs, err := s.AccountInfo.FilterExistingAdvertiserIDs(advertiserIDs)
	if err != nil {
		return fmt.Errorf("【SyncAccounts】过滤已存在账户失败: %w", err)
	}

	fmt.Printf("【SyncAccounts】过滤后需要同步 %d 个新账户ID", len(newAdvertiserIDs))

	if len(newAdvertiserIDs) == 0 {
		return nil
	}

	chunkSize := 50
	for i := 0; i < len(newAdvertiserIDs); i += chunkSize {
		end := i + chunkSize
		if end > len(newAdvertiserIDs) {
			end = len(newAdvertiserIDs)
		}
		chunk := newAdvertiserIDs[i:end]

		details, err := s.Oceanengine.GetAdvertiserInfo(ctx, accountID, chunk)
		if err != nil {
			return fmt.Errorf("【SyncAccounts】获取账户详情失败: %w", err)
		}

		err = s.AccountInfo.SaveToTable(details)
		if err != nil {
			continue
		}

		fmt.Printf("【SyncAccounts】成功保存批次数据: 批次 %d-%d, %d 条数据", i, end, len(details))
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

func (s *Crontab) RefreshToken() error {
	// 查询所有已授权的token
	result, err := s.AgentToken.Query(context.Background(), schema.AgentTokenQueryParam{})
	if err != nil {
		return fmt.Errorf("查询token失败: %w", err)
	}
	tokens := result.Data

	// 筛选需要刷新的token
	var expiredTokens []*schema.AgentToken
	for _, token := range tokens {
		if s.NeedRefresh(token) {
			expiredTokens = append(expiredTokens, token)
		}
	}

	// 逐个刷新token
	for _, token := range expiredTokens {
		fmt.Printf("开始刷新账号 %s 的token\n", token.AccountName)

		updatedToken, err := s.Oceanengine.RefreshToken(token)
		if err != nil {
			fmt.Printf("刷新账号 %s 的token失败: %v\n", token.AccountName, err)
			continue
		}

		// 保存到数据库
		err = s.AgentToken.Update(context.Background(), updatedToken)
		if err != nil {
			fmt.Printf("刷新账号 %s 的token失败: %v\n", token.AccountName, err)
			continue
		}

		fmt.Printf("刷新账号 %s 的token成功\n", token.AccountName)
		// 避免请求过于频繁
		time.Sleep(1 * time.Second)
	}

	return nil
}

// NeedRefresh 检查是否需要刷新token
func (s *Crontab) NeedRefresh(token *schema.AgentToken) bool {
	// 检查auth是否为已授权
	if token.AuthStatus != "已授权" {
		return false
	}

	// 检查token是否过期（当前时间戳 - tokentime > 86400），此处修改为提前400s，因为定时任务是每5分钟执行一次，所以提前400s确保在5分钟内刷新
	currentTime := time.Now().Unix()
	return currentTime-int64(token.TokenTime) > 86000
}

func (s *Crontab) HandleHostRule() error {
	ctx := context.Background()
	rules, err := s.HostRule.QueryAllEnabled(ctx)
	if err != nil {
		return fmt.Errorf("查询启用托管规则失败: %w", err)
	}

	now := time.Now()
	for _, rule := range rules {
		// 判断当前时间是否在生效日期范围内
		if now.Before(rule.TriggerStartDate) || now.After(rule.TriggerEndDate) {
			continue
		}

		type conditionItem struct {
			Metric   string `json:"metric"`
			Operator string `json:"operator"`
			Time     string `json:"time"`
			Unit     string `json:"unit"`
			Value    int64  `json:"value"`
		}
		// 解析 trigger_condition JSON 为 CustomReportReq
		// {"conditions":[{"metric":"stat_cost","operator":"\u003c","time":"today","unit":"元","value":20}],"logic":"and"}
		var timeRange schema.TimeFiltering
		if strings.Contains(rule.TriggerCondition, "last_7days") {
			timeRange = schema.TimeFiltering{
				CreateStartTime: now.AddDate(0, 0, -7).Format("2006-01-02"),
				CreateEndTime:   now.Format("2006-01-02"),
			}
		} else if strings.Contains(rule.TriggerCondition, "last_5days") {
			timeRange = schema.TimeFiltering{
				CreateStartTime: now.AddDate(0, 0, -5).Format("2006-01-02"),
				CreateEndTime:   now.Format("2006-01-02"),
			}
		} else if strings.Contains(rule.TriggerCondition, "last_3days") {
			timeRange = schema.TimeFiltering{
				CreateStartTime: now.AddDate(0, 0, -3).Format("2006-01-02"),
				CreateEndTime:   now.Format("2006-01-02"),
			}
		} else if strings.Contains(rule.TriggerCondition, "today") {
			timeRange = schema.TimeFiltering{
				CreateStartTime: now.Format("2006-01-02"),
				CreateEndTime:   now.Format("2006-01-02"),
			}
		} else {
			continue
		}

		var triggerCondition struct {
			Logic      string          `json:"logic"`
			Conditions []conditionItem `json:"conditions"`
		}

		if err := json.Unmarshal([]byte(rule.TriggerCondition), &triggerCondition); err != nil {
			fmt.Printf("解析规则 %d trigger_condition 失败: %v\n", rule.ID, err)
			continue
		}

		// 遍历获取所有metric并添加到metrics中
		var metrics []string

		metricTime := make(map[string]schema.TimeFiltering)
		for _, condition := range triggerCondition.Conditions {
			metrics = append(metrics, condition.Metric)
			switch condition.Time {
			case "today":
				metricTime[condition.Metric] = schema.TimeFiltering{
					CreateStartTime: now.Format("2006-01-02"),
					CreateEndTime:   now.Format("2006-01-02"),
				}
			case "last_3days":
				metricTime[condition.Metric] = schema.TimeFiltering{
					CreateStartTime: now.AddDate(0, 0, -3).Format("2006-01-02"),
					CreateEndTime:   now.Format("2006-01-02"),
				}
			case "last_5days":
				metricTime[condition.Metric] = schema.TimeFiltering{
					CreateStartTime: now.AddDate(0, 0, -5).Format("2006-01-02"),
					CreateEndTime:   now.Format("2006-01-02"),
				}
			case "last_7days":
				metricTime[condition.Metric] = schema.TimeFiltering{
					CreateStartTime: now.AddDate(0, 0, -7).Format("2006-01-02"),
					CreateEndTime:   now.Format("2006-01-02"),
				}
			}
		}

		// 解析 TargetAccounts JSON 获取 account IDs
		var accountIDs []string
		if err := json.Unmarshal([]byte(rule.TargetAccounts), &accountIDs); err != nil {
			fmt.Printf("解析规则 %d target_accounts 失败: %v\n", rule.ID, err)
			continue
		}

		var materialIDs []string
		if err := json.Unmarshal([]byte(rule.TargetMaterial), &materialIDs); err != nil {
			fmt.Printf("解析规则 %d target_material 失败: %v\n", rule.ID, err)
			continue
		}
		var promotionIDs []string
		if err := json.Unmarshal([]byte(rule.TargetPromotion), &promotionIDs); err != nil {
			fmt.Printf("解析规则 %d target_promotion 失败: %v\n", rule.ID, err)
			continue
		}

		for _, accountIDStr := range accountIDs {
			advertiserID, err := strconv.ParseInt(accountIDStr, 10, 64)
			if err != nil {
				continue
			}

			// 查询nb_promotion_material表, 使用promotionIDs或者projectIDs或者materialIDs作为in条件结合advertiserID筛选出符合目标的记录，用于filtering
			var dimensions []string
			var topic string
			var existingPromotionIDs []int64
			var existingMaterialIDs []int64
			switch rule.Target {
			case "promotion":
				dimensions = []string{"stat_time_day", "cdp_promotion_id"}
				topic = "BASIC_DATA"
				// 筛选出已存在的 promotion_id
				promotionIDInts := make([]int64, len(promotionIDs))
				for i, id := range promotionIDs {
					promotionIDInts[i], _ = strconv.ParseInt(id, 10, 64)
				}
				existingPromotionIDs, err = s.PromotionMaterial.FindExistingTargetIDs(ctx, promotionIDInts, advertiserID, "promotion")
				if err != nil {
					fmt.Printf("规则 %d 查询已存在的 promotion_id 失败(account=%s): %v\n", rule.ID, accountIDStr, err)
					continue
				}
				if len(existingPromotionIDs) == 0 {
					continue
				}
			case "creative":
				dimensions = []string{"stat_time_day", "material_id", "cdp_promotion_id"}
				topic = "MATERIAL_DATA"
				// 筛选出已存在的 material_id
				materialIDInts := make([]int64, len(materialIDs))
				for i, id := range materialIDs {
					materialIDInts[i], _ = strconv.ParseInt(id, 10, 64)
				}
				existingMaterialIDs, err = s.PromotionMaterial.FindExistingTargetIDs(ctx, materialIDInts, advertiserID, "material")
				if err != nil {
					fmt.Printf("规则 %d 查询已存在的 material_id 失败(account=%s): %v\n", rule.ID, accountIDStr, err)
					continue
				}
				if len(existingMaterialIDs) == 0 {
					continue
				}
			}

			// 分页循环请求，收集所有数据
			var allRows []schema.CustomReportRow
			page := 1
			for {
				req := schema.CustomReportReq{
					AdvertiserID: advertiserID,
					DataTopic:    topic,
					Dimensions:   dimensions,
					Metrics:      metrics,
					Filters: []schema.CustomReportFilter{
						{Field: "image_mode", Type: 2, Operator: 7, Values: []string{"5", "15"}},
					},
					StartTime: timeRange.CreateStartTime,
					EndTime:   timeRange.CreateEndTime,
					Page:      page,
					PageSize:  100,
				}

				// 根据target类型添加对应的筛选条件
				if rule.Target == "creative" {
					values := make([]string, len(existingMaterialIDs))
					for i, id := range existingMaterialIDs {
						values[i] = strconv.FormatInt(id, 10)
					}
					req.Filters = append(req.Filters, schema.CustomReportFilter{
						Field:    "material_id",
						Type:     2,
						Operator: 7,
						Values:   values,
					})
				} else {
					values := make([]string, len(existingPromotionIDs))
					for i, id := range existingPromotionIDs {
						values[i] = strconv.FormatInt(id, 10)
					}
					req.Filters = append(req.Filters, schema.CustomReportFilter{
						Field:    "cdp_promotion_id",
						Type:     2,
						Operator: 7,
						Values:   values,
					})
				}

				resp, err := s.Oceanengine.QueryCustomReport(ctx, rule.AgentID, req)
				if err != nil {
					fmt.Printf("规则 %d 查询自定义报表失败(account=%s, page=%d): %v\n", rule.ID, accountIDStr, page, err)
					break
				}

				allRows = append(allRows, resp.Data.Rows...)

				if page >= resp.Data.PageInfo.TotalPage || len(resp.Data.Rows) == 0 {
					break
				}
				page++
				time.Sleep(100 * time.Millisecond)
			}

			if len(allRows) == 0 {
				continue
			}

			// 聚合每个metric的数据，仅在stat_time_day处于metricTime区间内累加
			metricTotals := make(map[string]float64)
			for _, row := range allRows {
				statDay, _ := row.Dimensions["stat_time_day"].(string)
				for _, condition := range triggerCondition.Conditions {
					// 检查stat_time_day是否在metric时间区间内
					mt, ok := metricTime[condition.Metric]
					if ok && statDay != "" {
						if statDay < mt.CreateStartTime || statDay > mt.CreateEndTime {
							continue
						}
					}
					val, ok := row.Metrics[condition.Metric]
					if !ok {
						continue
					}
					switch v := val.(type) {
					case float64:
						metricTotals[condition.Metric] += v
					case json.Number:
						if f, err := v.Float64(); err == nil {
							metricTotals[condition.Metric] += f
						}
					case string:
						if f, err := strconv.ParseFloat(v, 64); err == nil {
							metricTotals[condition.Metric] += f
						}
					}
				}
			}

			// 根据logic比较聚合结果：or任一满足即预警，and全部满足才预警
			triggered := false
			if triggerCondition.Logic == "or" {
				for _, condition := range triggerCondition.Conditions {
					total := metricTotals[condition.Metric]
					if compareValue(total, condition.Operator, float64(condition.Value)) {
						triggered = true
						break
					}
				}
			} else { // "and"
				triggered = true
				for _, condition := range triggerCondition.Conditions {
					total := metricTotals[condition.Metric]
					if !compareValue(total, condition.Operator, float64(condition.Value)) {
						triggered = false
						break
					}
				}
			}

			if triggered {
				fmt.Printf("规则 %d 触发预警(account=%s)\n", rule.ID, accountIDStr)
				// todo 执行预警动作
			}
		}
	}

	return nil
}

func (s *Crontab) SyncPromotionMaterial(agentID, advertiserID int64) error {
	ctx := context.Background()

	// 1. 查询 nb_host_account 表获取 status=1 的账户
	accounts, err := s.HostAccount.FindAllEnabled(ctx)
	if err != nil {
		return fmt.Errorf("查询托管账户失败: %w", err)
	}

	// 根据参数过滤账户
	if agentID > 0 || advertiserID > 0 {
		var filtered schema.HostAccounts
		for _, a := range accounts {
			match := true
			if agentID > 0 && a.AgentID != agentID {
				match = false
			}
			if advertiserID > 0 && a.AdvertiserID != advertiserID {
				match = false
			}
			if match {
				filtered = append(filtered, a)
			}
		}
		accounts = filtered
	}
	if len(accounts) == 0 {
		fmt.Println("【SyncPromotionMaterial】没有匹配的托管账户")
		return nil
	}

	filtering := map[string]interface{}{
		"status_first": schema.PromotionStatusEnable,
	}
	fields := []string{"promotion_id", "promotion_materials", "advertiser_id", "promotion_name", "opt_status", "status_first"}

	for _, account := range accounts {
		// 2. 查询 promotion data
		items, err := s.Oceanengine.GetRefPromotionData(ctx, account.AgentID, account.AdvertiserID, filtering, fields)
		if err != nil {
			fmt.Printf("【SyncPromotionMaterial】拉取账户 %d 广告失败: %v\n", account.AdvertiserID, err)
			continue
		}

		// 收集所有 material_id 和 promotion_material 记录
		var pmItems []schema.PromotionMaterial
		materialIDSet := make(map[int64]bool)

		for _, item := range items {
			if item.PromotionMaterials == nil {
				continue
			}
			for _, videoMaterial := range item.PromotionMaterials.VideoMaterialList {
				pm := schema.PromotionMaterial{
					AdvertiserID:   account.AdvertiserID,
					PromotionID:    item.PromotionId,
					PromotionName:  item.PromotionName,
					StatusFirst:    item.StatusFirst,
					OptStatus:      item.OptStatus,
					MaterialStatus: videoMaterial.MaterialStatus,
					MaterialID:     videoMaterial.MaterialID,
				}
				pmItems = append(pmItems, pm)
				materialIDSet[videoMaterial.MaterialID] = true
			}
		}

		if len(pmItems) == 0 {
			fmt.Printf("【SyncPromotionMaterial】账户 %d 没有素材数据\n", account.AdvertiserID)
			continue
		}

		// 3. 写入 nb_promotion_material 表（存在更新，不存在插入）
		if err := s.PromotionMaterial.UpsertBatch(ctx, pmItems); err != nil {
			fmt.Printf("【SyncPromotionMaterial】写入 nb_promotion_material 失败(advertiserID=%d): %v\n", account.AdvertiserID, err)
			continue
		}

		// 4. 对于 nb_material_video 中不存在的 material_id，调用 GetVideoMaterial 写入
		allMaterialIDs := make([]int64, 0, len(materialIDSet))
		for mid := range materialIDSet {
			allMaterialIDs = append(allMaterialIDs, mid)
		}
		existingIDs, err := s.MaterialVideo.FindExistingMaterialIDs(ctx, allMaterialIDs)
		if err != nil {
			fmt.Printf("【SyncPromotionMaterial】查询已存在 material_id 失败: %v\n", err)
			continue
		}

		var newMaterialIDs []int64
		for _, mid := range allMaterialIDs {
			if !existingIDs[mid] {
				newMaterialIDs = append(newMaterialIDs, mid)
			}
		}

		if len(newMaterialIDs) > 0 {
			// 分批拉取，单次最多 100 个
			var allVideos []schema.MaterialVideo
			batchSize := 100
			for i := 0; i < len(newMaterialIDs); i += batchSize {
				end := i + batchSize
				if end > len(newMaterialIDs) {
					end = len(newMaterialIDs)
				}
				batch := newMaterialIDs[i:end]

				videoFiltering := map[string]interface{}{
					"material_ids": batch,
				}
				videos, err := s.Oceanengine.GetVideoMaterial(ctx, account.AgentID, account.AdvertiserID, videoFiltering)
				if err != nil {
					fmt.Printf("【SyncPromotionMaterial】拉取视频素材失败(advertiserID=%d, batch=%d-%d): %v\n", account.AdvertiserID, i, end, err)
					continue
				}
				allVideos = append(allVideos, videos...)
				time.Sleep(100 * time.Millisecond)
			}

			if len(allVideos) > 0 {
				if err := s.MaterialVideo.SaveBatch(ctx, allVideos); err != nil {
					fmt.Printf("【SyncPromotionMaterial】保存视频素材失败(advertiserID=%d): %v\n", account.AdvertiserID, err)
				} else {
					fmt.Printf("【SyncPromotionMaterial】保存 %d 条视频素材\n", len(allVideos))
				}
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

	return nil
}

// compareValue 比较两个float64值，支持 <=, <, >=, >, ==, !=
func compareValue(actual float64, operator string, expected float64) bool {
	switch operator {
	case "<=":
		return actual <= expected
	case "<":
		return actual < expected
	case ">=":
		return actual >= expected
	case ">":
		return actual > expected
	case "==":
		return actual == expected
	case "!=":
		return actual != expected
	default:
		return false
	}
}
