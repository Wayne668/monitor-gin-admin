/**
 * 素材删除记录接口
 */
import request from '@/utils/request'
// 获取素材删除记录列表
export const getDeleteUnauditedMaterialList = (params) => request.basic.get('/api/v1/delete-unaudited-material-records', params)
// 获取素材删除记录详情
export const getDeleteUnauditedMaterial = (id) => request.basic.get(`/api/v1/delete-unaudited-material-records/${id}`)