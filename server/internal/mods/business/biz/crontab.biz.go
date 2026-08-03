package biz

import (
	"fmt"
	"monitor-gin-admin/internal/mods/business/dal"
	"monitor-gin-admin/internal/mods/business/schema"
	"time"
)

type CronTab struct {
	Acc *dal.AccountInfo
}

func (s *CronTab) SyncAccounts(StartDate, EndDate string) error {
	accessToken, advertiserID := "", ""
	if accessToken == "" {
		return fmt.Errorf("remarks=%s 未找到已授权的Token", "")
	}

	filtering := &schema.TimeFiltering{
		CreateStartTime: StartDate + " 00:00:00",
		CreateEndTime:   EndDate + " 23:59:59",
	}

	advertiserIDs, err := GetAdvertiserIDs(accessToken, advertiserID, filtering)
	if err != nil {
		return fmt.Errorf("获取账户ID列表失败: %w", err)
	}

	fmt.Printf("【SyncAccounts】API返回 %d 个账户ID", len(advertiserIDs))

	if len(advertiserIDs) == 0 {
		return nil
	}

	// 调用model过滤已存在账户ID
	newAdvertiserIDs, err := s.Acc.FilterExistingAdvertiserIDs(advertiserIDs)
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

		details, err := GetAdvertiserInfo(accessToken, chunk)
		if err != nil {
			return fmt.Errorf("【SyncAccounts】获取账户详情失败: %w", err)
		}

		err = s.Acc.SaveToTable(details)
		if err != nil {
			continue
		}

		fmt.Printf("【SyncAccounts】成功保存批次数据: 批次 %d-%d, %d 条数据", i, end, len(details))
		time.Sleep(200 * time.Millisecond)
	}

	return nil
}
