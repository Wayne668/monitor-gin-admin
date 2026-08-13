/**
 * 广告主预算接口
 */
import request from '@/utils/request'

// 立即生效：更新广告主预算
export const updateAdvertiserBudget = (data) => request.basic.post('/api/v1/advertiser-budget/update', data)

// 次日生效：保存预算记录
export const scheduleAdvertiserBudget = (data) => request.basic.post('/api/v1/advertiser-budget/schedule', data)

// 查询预算修改记录
export const getBudgetRecords = (params) => request.basic.get('/api/v1/advertiser-budget/records', params)