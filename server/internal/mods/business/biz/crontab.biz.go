package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"strings"
	"time"
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

		// 解析 trigger_condition JSON 为 CustomReportReq
		var condition struct {
			DataTopic  string                       `json:"dataTopic"`
			Dimensions []string                     `json:"dimensions"`
			Metrics    []string                     `json:"metrics"`
			Filters    []schema.CustomReportFilter  `json:"filters"`
			OrderBy    []schema.CustomReportOrderBy `json:"order_by"`
		}
		if err := json.Unmarshal([]byte(rule.TriggerCondition), &condition); err != nil {
			fmt.Printf("解析规则 %d trigger_condition 失败: %v\n", rule.ID, err)
			continue
		}

		// 解析 TargetAccounts JSON 获取 account IDs
		var accountIDs []string
		if err := json.Unmarshal([]byte(rule.TargetAccounts), &accountIDs); err != nil {
			fmt.Printf("解析规则 %d target_accounts 失败: %v\n", rule.ID, err)
			continue
		}

		for _, accountIDStr := range accountIDs {
			advertiserID, err := fmt.Sscanf(accountIDStr, "%d", new(int64))
			if err != nil {
				continue
			}
			_ = advertiserID

			req := schema.CustomReportReq{
				DataTopic:  condition.DataTopic,
				Dimensions: condition.Dimensions,
				Metrics:    condition.Metrics,
				Filters:    condition.Filters,
				OrderBy:    condition.OrderBy,
				StartTime:  now.Format("2006-01-02"),
				EndTime:    now.Format("2006-01-02"),
				Page:       1,
				PageSize:   10,
			}

			_, err = s.Oceanengine.QueryCustomReport(ctx, rule.AgentID, req)
			if err != nil {
				fmt.Printf("规则 %d 查询自定义报表失败(account=%s): %v\n", rule.ID, accountIDStr, err)
				continue
			}
		}
	}

	return nil
}

func (s *Crontab) SyncPromotionMaterial() error {
	ctx := context.Background()

	// 1. 查询 nb_host_account 表获取 status=1 的账户
	accounts, err := s.HostAccount.FindAllEnabled(ctx)
	if err != nil {
		return fmt.Errorf("查询托管账户失败: %w", err)
	}
	if len(accounts) == 0 {
		fmt.Println("【SyncPromotionMaterial】没有启用的托管账户")
		return nil
	}

	filtering := map[string]interface{}{
		"status_first": schema.PromotionStatusEnable,
	}
	fields := []string{"promotion_id", "promotion_materials", "advertiser_id", "promotion_name"}

	for _, account := range accounts {
		fmt.Printf("【SyncPromotionMaterial】开始同步账户: agentID=%d, advertiserID=%d\n", account.AgentID, account.AdvertiserID)

		// 2. 查询 promotion data
		items, err := s.Oceanengine.GetRefPromotionData(ctx, account.AgentID, account.AdvertiserID, filtering, fields)
		if err != nil {
			fmt.Printf("【SyncPromotionMaterial】拉取账户 %d 广告失败: %v\n", account.AdvertiserID, err)
			continue
		}

		// 收集所有 material_id 和 promotion_material 记录
		var pmItems []schema.PromotionMaterial
		var allMaterialIDs []int64

		for _, item := range items {
			if item.PromotionMaterials == nil {
				continue
			}
			for _, videoMaterial := range item.PromotionMaterials.VideoMaterialList {
				statusSecond := ""
				if len(item.StatusSecond) > 0 {
					statusSecond = strings.Join(item.StatusSecond, ",")
				}
				pm := schema.PromotionMaterial{
					AdvertiserID:   account.AdvertiserID,
					PromotionID:    item.PromotionId,
					PromotionName:  item.PromotionName,
					StatusFirst:    item.StatusFirst,
					StatusSecond:   statusSecond,
					OptStatus:      item.OptStatus,
					MaterialStatus: videoMaterial.MaterialStatus,
					MaterialID:     videoMaterial.MaterialID,
				}
				pmItems = append(pmItems, pm)
				allMaterialIDs = append(allMaterialIDs, videoMaterial.MaterialID)
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
		fmt.Printf("【SyncPromotionMaterial】写入/更新 %d 条 promotion_material 记录\n", len(pmItems))

		// 4. 对于 nb_promotion_material 中不存在的 material_id，调用 GetVideoMaterial 写入 nb_material_video
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
			fmt.Printf("【SyncPromotionMaterial】账户 %d 有 %d 个新素材，拉取详情...\n", account.AdvertiserID, len(newMaterialIDs))
			now := time.Now()
			startDate := now.AddDate(0, 0, -30).Format("2006-01-02")
			endDate := now.Format("2006-01-02")

			videos, err := s.Oceanengine.GetVideoMaterial(ctx, account.AgentID, account.AdvertiserID, startDate, endDate)
			if err != nil {
				fmt.Printf("【SyncPromotionMaterial】拉取视频素材失败(advertiserID=%d): %v\n", account.AdvertiserID, err)
				continue
			}

			if len(videos) > 0 {
				if err := s.MaterialVideo.SaveBatch(ctx, videos); err != nil {
					fmt.Printf("【SyncPromotionMaterial】保存视频素材失败(advertiserID=%d): %v\n", account.AdvertiserID, err)
				} else {
					fmt.Printf("【SyncPromotionMaterial】保存 %d 条视频素材\n", len(videos))
				}
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

	return nil
}
