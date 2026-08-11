/**
 * 素材删除记录接口
 */
import request from '@/utils/request'
// 获取素材删除记录列表
export const getDeleteUnauditedMaterialList = (params) => request.basic.get('/api/v1/delete-unaudited-material-records', params)
// 获取素材删除记录详情
export const getDeleteUnauditedMaterial = (id) => request.basic.get(`/api/v1/delete-unaudited-material-records/${id}`)
// 重试删除失败的记录
export const retryFailedDelete = (accountId) => request.basic.post(`/api/v1/retry-failed-delete?accountId=${accountId}`)
// 获取未审核的素材列表
export const getUnauditedMaterial = (params) => request.basic.get('/api/v1/get-unaudited-material', params)
// 批量删除未审核素材
export const deleteUnauditedMaterial = (params) => request.basic.post('/api/v1/delete-unaudited-material', params)
