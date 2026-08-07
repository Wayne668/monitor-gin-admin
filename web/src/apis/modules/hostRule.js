/**
 * 托管规则接口
 */
import request from '@/utils/request'

// 根据账户ID列表和目标类型查询目标列表
export const getTargetByAccount = (params) => request.basic.post('/api/v1/get-target-by-account', params)

// 获取未审核的素材列表
export const getUnauditedMaterial = (params) => request.basic.get('/api/v1/get-unaudited-material', params)

// 批量删除未审核素材
export const deleteUnauditedMaterial = (params) => request.basic.post('/api/v1/delete-unaudited-material', params)
