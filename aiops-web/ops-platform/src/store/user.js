// src/store/user.js
import { defineStore } from 'pinia'
import { login as loginApi, logout as logoutApi, getUserInfo as getUserInfoApi } from '@/api/auth'

export const useUserStore = defineStore('user', {
    state: () => ({
        token: localStorage.getItem('token') || '',
        userInfo: JSON.parse(localStorage.getItem('userInfo') || '{}'),
        permissions: JSON.parse(localStorage.getItem('permissions') || '[]')
    }),

    getters: {
        username: (state) => state.userInfo.username || '',
        nickname: (state) => state.userInfo.nickname || '',
        avatar: (state) => state.userInfo.avatar || '',
        roles: (state) => state.userInfo.roles || []
    },

    actions: {
        // 登录
        async login(loginForm) {
            try {
                const data = await loginApi(loginForm)

                this.token = data.token
                this.userInfo = data.user
                this.permissions = data.permissions || []

                localStorage.setItem('token', data.token)
                localStorage.setItem('userInfo', JSON.stringify(data.user))
                localStorage.setItem('permissions', JSON.stringify(data.permissions || []))

                return data
            } catch (error) {
                throw error
            }
        },

        // 登出
        async logout() {
            try {
                await logoutApi()
            } catch (error) {
                console.error('登出失败:', error)
            } finally {
                this.token = ''
                this.userInfo = {}
                this.permissions = []
                localStorage.removeItem('token')
                localStorage.removeItem('userInfo')
                localStorage.removeItem('permissions')
            }
        },

        // 获取用户信息
        async fetchUserInfo() {
            try {
                const data = await getUserInfoApi()
                this.userInfo = data.user || {}
                this.permissions = data.permissions || []

                localStorage.setItem('userInfo', JSON.stringify(this.userInfo))
                localStorage.setItem('permissions', JSON.stringify(this.permissions))

                return data
            } catch (error) {
                throw error
            }
        },

        // 检查权限
        hasPermission(permCode) {
            if (!permCode) return true
            return this.permissions.includes(permCode) || this.permissions.includes('*:*')
        }
    }
})
