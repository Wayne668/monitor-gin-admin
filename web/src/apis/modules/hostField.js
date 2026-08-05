/**
 * 托管字段接口
 */
import request from '@/utils/request'
// 获取字段列表
export const getHostFieldList = (params) => request.basic.get('/api/v1/host-fields', params)
// 获取全部字段（不分页，前端渲染用）
export const getHostFieldAll = () => request.basic.get('/api/v1/host-fields/all')
// 获取字段详情
export const getHostField = (id) => request.basic.get(`/api/v1/host-fields/${id}`)
// 新增字段
export const createHostField = (params) => request.basic.post('/api/v1/host-fields', params)
// 更新字段
export const updateHostField = (id, params) => request.basic.put(`/api/v1/host-fields/${id}`, params)
// 删除字段
export const delHostField = (id) => request.basic.delete(`/api/v1/host-fields/${id}`)
