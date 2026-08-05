import request from '@/utils/request'

// 获取状态为 STATUS_ENABLE 的账户列表
export const getEnabledAccountList = () => request.basic.get('/api/v1/account-list')
