import request from './request'

// 获取主机组列表
export const getHostGroupList = (params) => {
  return request({
    url: '/api/v1/cmdb/host-groups',
    method: 'get',
    params
  })
}

// 获取所有主机组
export const getAllHostGroups = () => {
  return request({
    url: '/api/v1/cmdb/host-groups/all',
    method: 'get'
  })
}

// 获取主机组详情
export const getHostGroupById = (id) => {
  return request({
    url: `/api/v1/cmdb/host-groups/${id}`,
    method: 'get'
  })
}

// 创建主机组
export const createHostGroup = (data) => {
  return request({
    url: '/api/v1/cmdb/host-groups',
    method: 'post',
    data
  })
}

// 更新主机组
export const updateHostGroup = (id, data) => {
  return request({
    url: `/api/v1/cmdb/host-groups/${id}`,
    method: 'put',
    data
  })
}

// 删除主机组
export const deleteHostGroup = (id) => {
  return request({
    url: `/api/v1/cmdb/host-groups/${id}`,
    method: 'delete'
  })
}
