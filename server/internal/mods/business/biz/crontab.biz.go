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
	"monitor-gin-admin/pkg/util"
)

type Crontab struct {
	AccountInfo         *dal.AccountInfo
	Oceanengine         *Oceanengine
	AgentToken          *dal.AgentToken
	HostRule            *dal.HostRule
	HostAccount         *dal.HostAccount
	PromotionMaterial   *dal.PromotionMaterial
	MaterialVideo       *dal.MaterialVideo
	HostTriggerRecord   *dal.HostTriggerRecord
	HostField           *dal.HostField
	UpdateAccountBudget *dal.UpdateAccountBudget
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

	fmt.Printf("[HandleHostRule] 查询到 %d 条启用规则\n", len(rules))

	now := time.Now()
	for _, rule := range rules {
		// 判断当前时间是否在生效日期范围内
		if now.Before(rule.TriggerStartDate) || now.After(rule.TriggerEndDate) {
			fmt.Printf("[HandleHostRule] 规则[%d] %s 不在生效日期范围内(%s ~ %s)，跳过\n",
				rule.ID, rule.RuleName, rule.TriggerStartDate.Format("2006-01-02"), rule.TriggerEndDate.Format("2006-01-02"))
			continue
		}

		fmt.Printf("[HandleHostRule] 开始处理规则[%d] %s target=%s action=%s operateMethod=%d\n",
			rule.ID, rule.RuleName, rule.Target, rule.ExecuteAction, rule.OperateMethod)

		type conditionItem struct {
			Metric   string  `json:"metric"`
			Operator string  `json:"operator"`
			Time     string  `json:"time"`
			Unit     string  `json:"unit"`
			Value    float64 `json:"value"`
		}
		// 解析 trigger_condition JSON 为 CustomReportReq
		// {"conditions":[{"metric":"stat_cost","operator":"\u003c","time":"today","unit":"元","value":20}]}
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
			fmt.Printf("[HandleHostRule] 规则[%d] %s 未识别的time范围，跳过\n", rule.ID, rule.RuleName)
			continue
		}

		fmt.Printf("[HandleHostRule] 规则[%d] 时间范围: %s ~ %s\n", rule.ID, timeRange.CreateStartTime, timeRange.CreateEndTime)

		var triggerCondition struct {
			Conditions []conditionItem `json:"conditions"`
		}

		if err := json.Unmarshal([]byte(rule.TriggerCondition), &triggerCondition); err != nil {
			fmt.Printf("[HandleHostRule] 规则[%d] 解析 trigger_condition 失败: %v\n", rule.ID, err)
			continue
		}

		fmt.Printf("[HandleHostRule] 规则[%d] 解析到 %d 个条件\n", rule.ID, len(triggerCondition.Conditions))

		// 遍历获取所有metric并添加到metrics中
		var metrics []string
		metricWithFormula := make(map[string]string)
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

			// 查询该 metric 的 formula，按 / 分割后添加到 metrics
			field, err := s.HostField.FindByField(ctx, condition.Metric)
			if err == nil && field != nil && field.Formula != "" {
				metricWithFormula[condition.Metric] = field.Formula
				fmt.Printf("[HandleHostRule] 规则[%d] metric=%s 查找到 formula=%s\n", rule.ID, condition.Metric, field.Formula)
				parts := strings.Split(field.Formula, "/")
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if p != "" {
						metrics = append(metrics, p)
					}
				}
			}
		}

		// 对 metrics 去重，并保证一定包含 stat_cost
		seen := make(map[string]bool)
		var deduplicated []string
		for _, m := range metrics {
			if !seen[m] {
				seen[m] = true
				deduplicated = append(deduplicated, m)
			}
		}
		if !seen["stat_cost"] {
			deduplicated = append(deduplicated, "stat_cost")
		}
		metrics = deduplicated

		fmt.Printf("[HandleHostRule] 规则[%d] 最终 metrics(%d个): %v\n", rule.ID, len(metrics), metrics)

		// 解析 TargetAccounts JSON 获取 account IDs
		var materialIDs []string
		var promotionIDs []string
		var targetAccountIDs []string
		if err := json.Unmarshal([]byte(rule.TargetAccounts), &targetAccountIDs); err != nil {
			fmt.Printf("解析规则 %d target_accounts 失败: %v\n", rule.ID, err)
			continue
		}

		// 根据目标获取相应的accountIDs、promotion、materialID
		var accountIDs []int64
		for _, id := range targetAccountIDs {
			advertiserID, _ := strconv.ParseInt(id, 10, 64)
			accountIDs = append(accountIDs, advertiserID)
		}

		var isExcludeTarget bool
		var isExcludeAccount bool
		switch rule.ScopeType {
		case "promotion": // 指定广告
			if err := json.Unmarshal([]byte(rule.TargetPromotion), &promotionIDs); err != nil {
				fmt.Printf("解析规则 %d target_promotion 失败: %v\n", rule.ID, err)
				continue
			}
		case "exclude_account_promotion": // 排除指定账户的广告
			isExcludeAccount = true
		case "exclude_promotion": // 排除指定广告
			isExcludeTarget = true
		case "creative": // 指定创意
			if err := json.Unmarshal([]byte(rule.TargetMaterial), &materialIDs); err != nil {
				fmt.Printf("解析规则 %d target_material 失败: %v\n", rule.ID, err)
				continue
			}
		case "exclude_account_creative": // 排除指定账户的创意
			isExcludeAccount = true
		case "exclude_creative": // 排除指定创意
			isExcludeTarget = true
		default:
			fmt.Printf("[HandleHostRule] 规则[%d] 未知目标: %s\n", rule.ID, rule.Target)
			continue
		}

		if isExcludeAccount {
			// 查询nb_host_account表
			accountIDInts := make([]int64, len(targetAccountIDs))
			for i, id := range targetAccountIDs {
				accountIDInts[i], _ = strconv.ParseInt(id, 10, 64)
			}
			excludeAccounts, err := s.HostAccount.FindExcludeAccount(ctx, accountIDInts)
			if err != nil {
				fmt.Printf("解析规则 %d target_accounts 失败: %v\n", rule.ID, err)
				continue
			}
			accountIDs = accountIDs[:0]
			for _, account := range excludeAccounts {
				accountIDs = append(accountIDs, account.AdvertiserID)
			}
			fmt.Printf("[HandleHostRule] 规则[%d] 排除账户 %v\n", rule.ID, excludeAccounts)
		}

		fmt.Printf("[HandleHostRule] 规则[%d] target=%s, 解析到 %d 个账户, %d 个promotion, %d 个material\n",
			rule.ID, rule.Target, len(accountIDs), len(promotionIDs), len(materialIDs))

		for _, advertiserID := range accountIDs {
			fmt.Printf("[HandleHostRule] 规则[%d] 处理账户 advertiserID=%d\n", rule.ID, advertiserID)

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
				existingPromotionIDs, err = s.PromotionMaterial.FindExistingTargetIDs(ctx, promotionIDInts, advertiserID, "promotion", isExcludeTarget)
				if err != nil {
					fmt.Printf("[HandleHostRule] 规则[%d] 查询已存在的 promotion_id 失败(account=%d): %v\n", rule.ID, advertiserID, err)
					continue
				}
				fmt.Printf("[HandleHostRule] 规则[%d] account=%d 找到 %d 个已存在的 promotion\n", rule.ID, advertiserID, len(existingPromotionIDs))
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
				existingMaterialIDs, err = s.PromotionMaterial.FindExistingTargetIDs(ctx, materialIDInts, advertiserID, "material", isExcludeTarget)
				if err != nil {
					fmt.Printf("[HandleHostRule] 规则[%d] 查询已存在的 material_id 失败(account=%d): %v\n", rule.ID, advertiserID, err)
					continue
				}
				fmt.Printf("[HandleHostRule] 规则[%d] account=%d 找到 %d 个已存在的 material\n", rule.ID, advertiserID, len(existingMaterialIDs))
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
					fmt.Printf("[HandleHostRule] 规则[%d] 查询自定义报表失败(account=%d, page=%d): %v\n", rule.ID, advertiserID, page, err)
					break
				}

				fmt.Printf("[HandleHostRule] 规则[%d] account=%d page=%d 返回 %d 行, totalPage=%d\n",
					rule.ID, advertiserID, page, len(resp.Data.Rows), resp.Data.PageInfo.TotalPage)

				allRows = append(allRows, resp.Data.Rows...)

				if page >= resp.Data.PageInfo.TotalPage || len(resp.Data.Rows) == 0 {
					break
				}
				page++
				time.Sleep(100 * time.Millisecond)
			}

			if len(allRows) == 0 {
				fmt.Printf("[HandleHostRule] 规则[%d] account=%d 无数据返回，跳过\n", rule.ID, advertiserID)
				continue
			}

			fmt.Printf("[HandleHostRule] 规则[%d] account=%d 共收集 %d 行数据\n", rule.ID, advertiserID, len(allRows))

			// 确定维度中用于区分目标的key
			targetDimKey := "cdp_promotion_id"
			if rule.Target == "creative" {
				targetDimKey = "material_id"
			} else {
				targetDimKey = "cdp_project_id"
			}

			// 按目标分别聚合每个metric的数据，仅在stat_time_day处于metricTime区间内累加
			targetMetrics := make(map[string]map[string]float64) // targetID -> metric -> total
			for _, row := range allRows {
				targetID, _ := row.Dimensions[targetDimKey].(string)
				if targetID == "" {
					continue
				}
				if _, ok := targetMetrics[targetID]; !ok {
					targetMetrics[targetID] = make(map[string]float64)
				}
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
						targetMetrics[targetID][condition.Metric] += v
					case json.Number:
						if f, err := v.Float64(); err == nil {
							targetMetrics[targetID][condition.Metric] += f
						}
					case string:
						if f, err := strconv.ParseFloat(v, 64); err == nil {
							targetMetrics[targetID][condition.Metric] += f
						}
					}
				}
			}

			fmt.Printf("[HandleHostRule] 规则[%d] account=%d 聚合完成，共 %d 个目标\n", rule.ID, advertiserID, len(targetMetrics))

			// 根据logic比较每个目标的聚合结果
			for targetID, metricTotals := range targetMetrics {
				targetTriggered := false
				triggerReason := ""
				if rule.OperateMethod == 1 { // or
					for _, condition := range triggerCondition.Conditions {
						var total float64
						if formula, ok := metricWithFormula[condition.Metric]; ok && formula != "" {
							// 使用公式算法计算total值，因为媒体返回的带公式的指标值不准
							formulaArr := strings.Split(formula, "/")
							if len(formulaArr) == 2 && metricTotals[formulaArr[1]] != 0 {
								total = metricTotals[formulaArr[0]] / metricTotals[formulaArr[1]]
							} else {
								total = metricTotals[condition.Metric]
							}
						} else {
							total = metricTotals[condition.Metric]
						}
						if compareValue(total, condition.Operator, float64(condition.Value)) {
							targetTriggered = true
							triggerReason = fmt.Sprintf("%s %s(%f) %s %f", condition.Time, condition.Metric, total, condition.Operator, condition.Value)
							fmt.Printf("[HandleHostRule] 规则[%d] target=%s(%s) OR触发: %s\n", rule.ID, rule.Target, targetID, triggerReason)
							break
						}
					}
					if !targetTriggered {
						fmt.Printf("[HandleHostRule] 规则[%d] target=%s(%s) OR未触发, totals=%v\n", rule.ID, rule.Target, targetID, metricTotals)
					}
				} else { // "and"
					targetTriggered = true
					for _, condition := range triggerCondition.Conditions {
						var total float64
						if formula, ok := metricWithFormula[condition.Metric]; ok && formula != "" {
							// 使用公式算法计算total值，因为媒体返回的带公式的指标值不准
							formulaArr := strings.Split(formula, "/")
							if len(formulaArr) == 2 && metricTotals[formulaArr[1]] != 0 {
								total = metricTotals[formulaArr[0]] / metricTotals[formulaArr[1]]
							} else {
								total = metricTotals[condition.Metric]
							}
						} else {
							total = metricTotals[condition.Metric]
						}
						if !compareValue(total, condition.Operator, float64(condition.Value)) {
							targetTriggered = false
							triggerReason = fmt.Sprintf("%s %s(%f) %s %f", condition.Time, condition.Metric, total, condition.Operator, condition.Value)
							fmt.Printf("[HandleHostRule] 规则[%d] target=%s(%s) AND未触发: %s\n", rule.ID, rule.Target, targetID, triggerReason)
							break
						}
					}
				}

				if targetTriggered {
					fmt.Printf("[HandleHostRule] 规则[%d] %s target=%s(%s) 触发预警! reason=%s\n", rule.ID, rule.RuleName, rule.Target, targetID, triggerReason)
					// 写入触发记录
					record := &schema.HostTriggerRecord{
						RuleID:        int(rule.ID),
						AdvertiserID:  advertiserID,
						Target:        rule.Target,
						ExecuteAction: rule.ExecuteAction,
						ExecuteStatus: "pending",
						TriggerReason: triggerReason,
					}
					// 根据target类型填充对应的ID字段
					targetIDInt, _ := strconv.ParseInt(targetID, 10, 64)
					switch rule.Target {
					case "creative":
						record.MaterialID = targetIDInt
					default:
						record.PromotionID = targetIDInt
					}
					if err := s.HostTriggerRecord.Create(ctx, record); err != nil {
						fmt.Printf("[HandleHostRule] 规则[%d] %s 写入触发记录失败: %v\n", rule.ID, rule.RuleName, err)
					} else {
						fmt.Printf("[HandleHostRule] 规则[%d] %s 触发记录写入成功\n", rule.ID, rule.RuleName)
					}
				}
			}
		}
	}

	fmt.Printf("[HandleHostRule] 处理完成\n")
	return nil
}

func (s *Crontab) SyncPromotionMaterial(agentID, advertiserID int64) error {
	ctx := context.Background()

	// 1. 查询 nb_host_account 表获取 status=1 的账户
	accounts, err := s.HostAccount.FindAllEnabled(ctx)
	if err != nil {
		return fmt.Errorf("查询托管账户失败: %w", err)
	}

	fmt.Printf("[SyncPromotionMaterial] 查询到 %d 个启用账户\n", len(accounts))

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
		fmt.Printf("[SyncPromotionMaterial] 过滤后 %d 个账户\n", len(accounts))
	}
	if len(accounts) == 0 {
		fmt.Println("[SyncPromotionMaterial] 没有匹配的托管账户")
		return nil
	}

	filtering := map[string]interface{}{
		"status_first": schema.PromotionStatusEnable,
	}
	fields := []string{"promotion_id", "promotion_materials", "advertiser_id", "promotion_name", "opt_status", "status_first"}

	for _, account := range accounts {
		fmt.Printf("[SyncPromotionMaterial] 处理账户 agentID=%d advertiserID=%d\n", account.AgentID, account.AdvertiserID)

		// 2. 查询 promotion data
		items, err := s.Oceanengine.GetRefPromotionData(ctx, account.AgentID, account.AdvertiserID, filtering, fields)
		if err != nil {
			fmt.Printf("[SyncPromotionMaterial] 拉取账户 %d 广告失败: %v\n", account.AdvertiserID, err)
			continue
		}

		fmt.Printf("[SyncPromotionMaterial] 账户 %d 拉取到 %d 条 promotion\n", account.AdvertiserID, len(items))

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
			fmt.Printf("[SyncPromotionMaterial] 账户 %d 没有素材数据\n", account.AdvertiserID)
			continue
		}

		fmt.Printf("[SyncPromotionMaterial] 账户 %d 收集到 %d 条 promotion_material, %d 个唯一 material_id\n",
			account.AdvertiserID, len(pmItems), len(materialIDSet))

		// 3. 写入 nb_promotion_material 表（存在更新，不存在插入）
		if err := s.PromotionMaterial.UpsertBatch(ctx, pmItems); err != nil {
			fmt.Printf("[SyncPromotionMaterial] 写入 nb_promotion_material 失败(advertiserID=%d): %v\n", account.AdvertiserID, err)
			continue
		}

		fmt.Printf("[SyncPromotionMaterial] 账户 %d 写入 %d 条 promotion_material 成功\n", account.AdvertiserID, len(pmItems))

		// 4. 对于 nb_material_video 中不存在的 material_id，调用 GetVideoMaterial 写入
		allMaterialIDs := make([]int64, 0, len(materialIDSet))
		for mid := range materialIDSet {
			allMaterialIDs = append(allMaterialIDs, mid)
		}
		existingIDs, err := s.MaterialVideo.FindExistingMaterialIDs(ctx, allMaterialIDs)
		if err != nil {
			fmt.Printf("[SyncPromotionMaterial] 查询已存在 material_id 失败: %v\n", err)
			continue
		}

		var newMaterialIDs []int64
		for _, mid := range allMaterialIDs {
			if !existingIDs[mid] {
				newMaterialIDs = append(newMaterialIDs, mid)
			}
		}

		fmt.Printf("[SyncPromotionMaterial] 账户 %d 需要新增 %d 个视频素材\n", account.AdvertiserID, len(newMaterialIDs))

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

				fmt.Printf("[SyncPromotionMaterial] 账户 %d 拉取视频素材 batch %d-%d\n", account.AdvertiserID, i, end)

				videoFiltering := map[string]interface{}{
					"material_ids": batch,
				}
				videos, err := s.Oceanengine.GetVideoMaterial(ctx, account.AgentID, account.AdvertiserID, videoFiltering)
				if err != nil {
					fmt.Printf("[SyncPromotionMaterial] 拉取视频素材失败(advertiserID=%d, batch=%d-%d): %v\n", account.AdvertiserID, i, end, err)
					continue
				}
				allVideos = append(allVideos, videos...)
				time.Sleep(100 * time.Millisecond)
			}

			if len(allVideos) > 0 {
				if err := s.MaterialVideo.SaveBatch(ctx, allVideos); err != nil {
					fmt.Printf("[SyncPromotionMaterial] 保存视频素材失败(advertiserID=%d): %v\n", account.AdvertiserID, err)
				} else {
					fmt.Printf("[SyncPromotionMaterial] 保存 %d 条视频素材\n", len(allVideos))
				}
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("[SyncPromotionMaterial] 处理完成\n")
	return nil
}

// SyncHostTriggerRecord 同步待执行的触发记录
func (s *Crontab) SyncHostTriggerRecord() error {
	ctx := context.Background()

	// 查询待执行的触发记录
	records, err := s.HostTriggerRecord.QueryPending(ctx)
	if err != nil {
		return fmt.Errorf("查询待执行触发记录失败: %w", err)
	}

	fmt.Printf("[SyncHostTriggerRecord] 查询到 %d 条待执行记录\n", len(records))

	if len(records) == 0 {
		return nil
	}

	for _, record := range records {
		// 映射 execute_action 到 optStatus: pause→DISABLE, restart→ENABLE
		optStatus := "DISABLE"
		if record.ExecuteAction == "restart" {
			optStatus = "ENABLE"
		}

		fmt.Printf("[SyncHostTriggerRecord] 处理记录[%d] ruleID=%d target=%s action=%s→%s advertiserID=%d promotionID=%d materialID=%d\n",
			record.ID, record.RuleID, record.Target, record.ExecuteAction, optStatus, record.AdvertiserID, record.PromotionID, record.MaterialID)

		// 获取规则信息（含钉钉 webhook）
		rule, err := s.HostRule.Get(ctx, uint(record.RuleID))
		if err != nil || rule == nil {
			fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 获取规则失败: %v\n", record.ID, err)
			_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "failed", fmt.Sprintf("获取规则失败: %v", err))
			continue
		}

		var apiResult string
		switch record.Target {
		case "promotion":
			// 检查是否已处于目标状态
			inStatus, err := s.PromotionMaterial.IsPromotionInTargetStatus(ctx, record.AdvertiserID, record.PromotionID, optStatus)
			if err != nil {
				fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 状态检查失败: %v\n", record.ID, err)
				_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "failed", fmt.Sprintf("状态检查失败: %v", err))
				continue
			}
			fmt.Printf("[SyncHostTriggerRecord] 记录[%d] promotionID=%d 目标状态=%s, 已达标=%v\n", record.ID, record.PromotionID, optStatus, inStatus)
			if inStatus {
				apiResult = "已处于目标状态，无需执行"
				_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "succeed", apiResult)
			} else {
				fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 调用 UpdatePromotionStatus promotionID=%d optStatus=%s\n", record.ID, record.PromotionID, optStatus)
				_, err = s.Oceanengine.UpdatePromotionStatus(ctx, rule.AgentID, optStatus, []int64{record.PromotionID}, record.AdvertiserID)
				if err != nil {
					apiResult = fmt.Sprintf("执行失败: %v", err)
					fmt.Printf("[SyncHostTriggerRecord] 记录[%d] UpdatePromotionStatus 失败: %v\n", record.ID, err)
					_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "failed", apiResult)
				} else {
					apiResult = "执行成功"
					fmt.Printf("[SyncHostTriggerRecord] 记录[%d] UpdatePromotionStatus 成功\n", record.ID)
					_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "succeed", apiResult)
				}
			}
		case "creative":
			// material_status 映射: DISABLE→MATERIAL_STATUS_DELETE, ENABLE→MATERIAL_STATUS_OK
			materialTargetStatus := "MATERIAL_STATUS_DELETE"
			if optStatus == "ENABLE" {
				materialTargetStatus = "MATERIAL_STATUS_OK"
			}
			inStatus, err := s.PromotionMaterial.IsMaterialInTargetStatus(ctx, record.AdvertiserID, record.MaterialID, materialTargetStatus)
			if err != nil {
				fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 素材状态检查失败: %v\n", record.ID, err)
				_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "failed", fmt.Sprintf("状态检查失败: %v", err))
				continue
			}
			fmt.Printf("[SyncHostTriggerRecord] 记录[%d] materialID=%d 目标状态=%s, 已达标=%v\n", record.ID, record.MaterialID, materialTargetStatus, inStatus)
			if inStatus {
				apiResult = "已处于目标状态，无需执行"
				_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "succeed", apiResult)
			} else {
				fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 调用 UpdateMaterialStatus materialID=%d optStatus=%s\n", record.ID, record.MaterialID, optStatus)
				_, err = s.Oceanengine.UpdateMaterialStatus(ctx, rule.AgentID, optStatus, []int64{record.MaterialID}, record.PromotionID, record.AdvertiserID)
				if err != nil {
					apiResult = fmt.Sprintf("执行失败: %v", err)
					fmt.Printf("[SyncHostTriggerRecord] 记录[%d] UpdateMaterialStatus 失败: %v\n", record.ID, err)
					_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "failed", apiResult)
				} else {
					apiResult = "执行成功"
					fmt.Printf("[SyncHostTriggerRecord] 记录[%d] UpdateMaterialStatus 成功\n", record.ID)
					_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "succeed", apiResult)
				}
			}
		default:
			fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 不支持的目标类型: %s\n", record.ID, record.Target)
			_ = s.HostTriggerRecord.UpdateResult(ctx, record.ID, "failed", fmt.Sprintf("不支持的目标类型: %s", record.Target))
			continue
		}

		fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 执行结果: %s\n", record.ID, apiResult)

		// 发送钉钉通知
		if rule.DingtalkWebhookUrl != "" {
			notifyContent := fmt.Sprintf("【托管规则执行通知】\n触发原因：%s\n执行结果：%s", record.TriggerReason, apiResult)
			if err := util.SendDingTalkText(rule.DingtalkWebhookUrl, notifyContent, nil, false); err != nil {
				fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 发送钉钉通知失败: %v\n", record.ID, err)
			} else {
				fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 钉钉通知发送成功\n", record.ID)
			}
		} else {
			fmt.Printf("[SyncHostTriggerRecord] 记录[%d] 未配置钉钉webhook，跳过通知\n", record.ID)
		}
	}

	fmt.Printf("[SyncHostTriggerRecord] 处理完成，共处理 %d 条记录\n", len(records))
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

// SyncAdvertiserBudget 拉取 nb_update_account_budget 中 is_set=0 且 budget_mod='nextDay' 的记录，
// 调用 Oceanengine API 更新预算
func (s *Crontab) SyncAdvertiserBudget() error {
	ctx := context.Background()

	records, err := s.UpdateAccountBudget.FindPendingNextDay(ctx)
	if err != nil {
		return fmt.Errorf("查询待执行预算记录失败: %w", err)
	}

	fmt.Printf("[SyncAdvertiserBudget] 查询到 %d 条待执行预算记录\n", len(records))

	if len(records) == 0 {
		return nil
	}

	// 查询所有托管账户，建立 advertiserID -> agentID 的映射
	accounts, err := s.HostAccount.FindAllEnabled(ctx)
	if err != nil {
		return fmt.Errorf("查询托管账户失败: %w", err)
	}

	accountMap := make(map[int64]int64) // advertiserID -> agentID
	for _, acc := range accounts {
		accountMap[acc.AdvertiserID] = acc.AgentID
	}

	for _, record := range records {
		agentID, ok := accountMap[record.AdvertiserID]
		if !ok {
			errMsg := fmt.Sprintf("未找到广告主 %d 对应的代理商", record.AdvertiserID)
			fmt.Printf("[SyncAdvertiserBudget] 记录[%d] %s\n", record.ID, errMsg)
			_ = s.UpdateAccountBudget.UpdateSetStatus(ctx, record.ID, 0, errMsg)
			continue
		}

		fmt.Printf("[SyncAdvertiserBudget] 处理记录[%d] advertiserID=%d budget=%.2f\n",
			record.ID, record.AdvertiserID, record.Budget)

		_, err := s.Oceanengine.UpdateAdvertiserBudget(ctx, agentID, record.AdvertiserID, "BUDGET_MODE_DAY", record.Budget)
		if err != nil {
			errMsg := fmt.Sprintf("更新预算失败: %v", err)
			fmt.Printf("[SyncAdvertiserBudget] 记录[%d] %s\n", record.ID, errMsg)
			_ = s.UpdateAccountBudget.UpdateSetStatus(ctx, record.ID, 0, errMsg)
		} else {
			fmt.Printf("[SyncAdvertiserBudget] 记录[%d] 更新成功\n", record.ID)
			_ = s.UpdateAccountBudget.UpdateSetStatus(ctx, record.ID, 1, "")
		}
	}

	fmt.Printf("[SyncAdvertiserBudget] 处理完成，共处理 %d 条记录\n", len(records))
	return nil
}
