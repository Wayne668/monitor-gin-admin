/**
 * 托管规则接口
 */
import request from '@/utils/request'
// 查询托管规则列表
export const getHostRuleList = (params) => request.basic.get('/api/v1/host-rules', params)
// 获取托管规则详情
export const getHostRule = (id) => request.basic.get(`/api/v1/host-rules/${id}`)
// 保存托管规则（新增/编辑）
export const saveHostRule = (data) => request.basic.post('/api/v1/host-rules', data)
// 更新规则状态
export const updateHostRuleStatus = (id, data) => request.basic.patch(`/api/v1/host-rules/${id}/status`, data)

// 根据账户ID列表和目标类型查询目标列表
export const getTargetByAccount = (params) => request.basic.post('/api/v1/get-target-by-account', params)

// 获取未审核的素材列表
export const getUnauditedMaterial = (params) => request.basic.get('/api/v1/get-unaudited-material', params)

// 批量删除未审核素材
export const deleteUnauditedMaterial = (params) => request.basic.post('/api/v1/delete-unaudited-material', params)

// 根据代理商ID获取托管账户列表
export const getHostAccountList = (params) => request.basic.get('/api/v1/host-account-list', params)