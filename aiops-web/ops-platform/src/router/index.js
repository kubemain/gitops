// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'
import Layout from '@/views/Layout.vue'
import CmdbLayout from '@/views/cmdb/CmdbLayout.vue'

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
    // CMDB管理
    {
        path: '/cmdb',
        component: CmdbLayout,
        redirect: '/cmdb/dashboard',
        meta: { title: 'CMDB管理', icon: 'Server' },
        children: [
            // 总览
            {
                path: 'dashboard',
                name: 'CmdbDashboard',
                component: () => import('@/views/cmdb/CmdbDashboard.vue'),
                meta: { title: '总览' }
            },
            // 业务线管理
            {
                path: 'business',
                name: 'Business',
                component: () => import('@/views/cmdb/Business.vue'),
                meta: { title: '业务线管理' }
            },
            // 应用管理
            {
                path: 'applications',
                name: 'Applications',
                component: () => import('@/views/cmdb/Application.vue'),
                meta: { title: '应用管理' }
            },
            // 服务管理
            {
                path: 'services',
                name: 'Services',
                component: () => import('@/views/cmdb/Service.vue'),
                meta: { title: '服务管理' }
            },
            // 主机管理
            {
                path: 'hosts',
                name: 'Hosts',
                component: () => import('@/views/cmdb/Host.vue'),
                meta: { title: '主机管理' }
            },
            // 主机分组
            {
                path: 'host-groups',
                name: 'HostGroups',
                component: () => import('@/views/cmdb/HostGroup.vue'),
                meta: { title: '主机分组' }
            },
            // 多维视图
            {
                path: 'view/topology',
                name: 'TopologyView',
                component: () => import('@/views/cmdb/views/TopologyView.vue'),
                meta: { title: '拓扑视图' }
            },
            {
                path: 'view/environment',
                name: 'EnvironmentView',
                component: () => import('@/views/cmdb/views/EnvironmentView.vue'),
                meta: { title: '环境视图' }
            },
            {
                path: 'view/region',
                name: 'RegionView',
                component: () => import('@/views/cmdb/views/RegionView.vue'),
                meta: { title: '区域视图' }
            },
            // 数据导入
            {
                path: 'import',
                name: 'DataImport',
                component: () => import('@/views/cmdb/DataImport.vue'),
                meta: { title: '数据导入' }
            },
            // 变更历史
            {
                path: 'changes',
                name: 'ChangeHistory',
                component: () => import('@/views/cmdb/ChangeHistory.vue'),
                meta: { title: '变更历史' }
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
