// src/api/role.js
import request from './request'

// 获取角色列表
export function getRoleList() {
    return request({
        url: '/roles',
        method: 'get'
    })
}

// 获取角色详情
export function getRoleDetail(id) {
    return request({
        url: `/roles/${id}`,
        method: 'get'
    })
}

// 分配权限
export function assignPermissions(id, data) {
    return request({
        url: `/roles/${id}/permissions`,
        method: 'post',
        data
    })
}

// 获取所有权限
export function getAllPermissions() {
    return request({
        url: '/permissions',
        method: 'get'
    })
}
