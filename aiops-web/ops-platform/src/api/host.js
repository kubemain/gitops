import request from './request'

// 获取主机列表
export const getHostList = (params) => {
  return request({
    url: '/api/v1/cmdb/hosts',
    method: 'get',
    params
  })
}

// 获取主机详情
export const getHostById = (id) => {
  return request({
    url: `/api/v1/cmdb/hosts/${id}`,
    method: 'get'
  })
}

// 创建主机
export const createHost = (data) => {
  return request({
    url: '/api/v1/cmdb/hosts',
    method: 'post',
    data
  })
}

// 更新主机
export const updateHost = (id, data) => {
  return request({
    url: `/api/v1/cmdb/hosts/${id}`,
    method: 'put',
    data
  })
}

// 删除主机
export const deleteHost = (id) => {
  return request({
    url: `/api/v1/cmdb/hosts/${id}`,
    method: 'delete'
  })
}

// 批量删除主机
export const batchDeleteHost = (data) => {
  return request({
    url: '/api/v1/cmdb/hosts/batch-delete',
    method: 'post',
    data
  })
}

// 更新主机状态
export const updateHostStatus = (id, data) => {
  return request({
    url: `/api/v1/cmdb/hosts/${id}/status`,
    method: 'patch',
    data
  })
}

// 获取主机统计信息
export const getHostStats = () => {
  return request({
    url: '/api/v1/cmdb/hosts/stats',
    method: 'get'
  })
}
