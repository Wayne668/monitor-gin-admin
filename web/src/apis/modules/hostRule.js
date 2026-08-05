/**
 * 托管规则接口
 */
import request from '@/utils/request'

// 根据账户ID列表和目标类型查询目标列表
export const getTargetByAccount = (params) => request.basic.post('/api/v1/get-target-by-account', params)
