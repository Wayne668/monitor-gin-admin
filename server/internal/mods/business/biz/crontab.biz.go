package biz

import (
	"context"
	"fmt"
	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"time"
)

type Crontab struct {
	AccountInfo *dal.AccountInfo
	Oceanengine *Oceanengine
	AgentToken  *dal.AgentToken
}

func (s *Crontab) SyncAccounts(StartDate, EndDate string) error {
	accessToken, advertiserID := "", ""
	if accessToken == "" {
		return fmt.Errorf("remarks=%s 未找到已授权的Token", "")
	}

	filtering := &schema.TimeFiltering{
		CreateStartTime: StartDate + " 00:00:00",
		CreateEndTime:   EndDate + " 23:59:59",
	}

	advertiserIDs, err := s.Oceanengine.GetAdvertiserIDs(accessToken, advertiserID, filtering)
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

		details, err := s.Oceanengine.GetAdvertiserInfo(accessToken, chunk)
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
