/**
 * 托管账户接口
 */
import request from '@/utils/request'

// 查询托管账户列表
export const getHostAccountList = (params) => request.basic.get('/api/v1/host-accounts', params)

// 获取托管账户详情
export const getHostAccount = (id) => request.basic.get(`/api/v1/host-accounts/${id}`)

// 新增托管账户
export const createHostAccount = (data) => request.basic.post('/api/v1/host-accounts', data)

// 更新托管账户
export const updateHostAccount = (id, data) => request.basic.put(`/api/v1/host-accounts/${id}`, data)

// 删除托管账户
export const delHostAccount = (id) => request.basic.delete(`/api/v1/host-accounts/${id}`)