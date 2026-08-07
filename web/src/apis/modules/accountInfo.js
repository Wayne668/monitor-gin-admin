import request from '@/utils/request'

// 获取状态为 STATUS_ENABLE 的账户列表
// @param {Object} [params] - 查询参数
//   - fields: 逗号分隔的字段名，例如 "advertiser_id,advertiser_name"，为空则返回全部字段
//   - limit:  返回条数，<=0 或不传时后端默认返回最新 100 条（按 updated_at 倒序）
export const getEnabledAccountList = (params) => request.basic.get('/api/v1/account-list', params)
