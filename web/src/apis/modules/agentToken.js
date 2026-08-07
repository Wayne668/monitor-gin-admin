/**
 * 账户token接口
 */
import request from '@/utils/request'
// 获取账户token列表
export const getAgentTokenList = (params) => request.basic.get('/api/v1/agent-tokens', params)
// 获取账户token详情
export const getAgentToken = (id) => request.basic.get(`/api/v1/agent-tokens/${id}`)
// 新增账户token
export const createAgentToken = (params) => request.basic.post('/api/v1/agent-tokens', params)
// 更新账户token
export const updateAgentToken = (id, params) => request.basic.put(`/api/v1/agent-tokens/${id}`, params)
// 删除账户token
export const delAgentToken = (id) => request.basic.delete(`/api/v1/agent-tokens/${id}`)
