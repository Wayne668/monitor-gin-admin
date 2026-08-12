/**
 * 托管触发记录接口
 */
import request from '@/utils/request'

// 查询触发记录列表
export const getHostTriggerRecordList = (params) => request.basic.get('/api/v1/host-trigger-records', params)