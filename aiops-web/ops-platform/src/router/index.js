// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import Layout from '@/views/Layout.vue'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue'),
        meta: { title: '登录', hidden: true }
    },
    {
        path: '/',
        name: 'Home',
        component: () => import('@/components/HomePage.vue'),
        meta: { title: '首页', hidden: true }
    },
    {
        path: '/system',
        component: Layout,
        redirect: '/system/user',
        meta: { title: '系统管理', icon: 'Setting' },
        children: [
            {
                path: 'user',
                name: 'User',
                component: () => import('@/views/system/User.vue'),
                meta: { title: '用户管理', icon: 'User' }
            },
            {
                path: 'role',
                name: 'Role',
                component: () => import('@/views/system/Role.vue'),
                meta: { title: '角色管理', icon: 'UserFilled' }
            }
        ]
    },
    {
        path: '/404',
        name: '404',
        component: () => import('@/views/Login.vue'),
        meta: { title: '404', hidden: true }
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/404',
        meta: { hidden: true }
    }
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
    scrollBehavior() {
        return { top: 0 }
    }
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()

    document.title = to.meta.title
        ? `${to.meta.title} - AIOps智能运维平台`
        : 'AIOps智能运维平台'

    const whiteList = ['/login']

    if (userStore.token) {
        if (to.path === '/login') {
            next({ path: '/' })
        } else {
            if (!userStore.userInfo || !userStore.userInfo.username) {
                try {
                    await userStore.fetchUserInfo()
                    next()
                } catch (error) {
                    userStore.logout()
                    next(`/login?redirect=${to.path}`)
                }
            } else {
                next()
            }
        }
    } else {
        if (whiteList.includes(to.path)) {
            next()
        } else {
            next(`/login?redirect=${to.path}`)
        }
    }
})

export default router
