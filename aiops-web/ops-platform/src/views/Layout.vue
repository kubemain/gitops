<!-- src/views/Layout.vue -->
<template>
  <div class="layout-container">
    <!-- 顶部导航栏 -->
    <header class="layout-header">
      <div class="header-left">
        <div class="logo">
          <span class="logo-icon">⚙️</span>
          <span class="logo-text">系统管理</span>
        </div>
        <el-button type="primary" link @click="$router.push('/')">
          <el-icon><HomeFilled /></el-icon>
          返回工作台
        </el-button>
      </div>

      <div class="header-right">
        <el-dropdown @command="handleCommand">
          <div class="user-info">
            <el-avatar :src="userStore.avatar" :size="32">
              {{ userStore.username?.charAt(0).toUpperCase() }}
            </el-avatar>
            <span class="username">{{ userStore.nickname || userStore.username }}</span>
            <el-icon><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                个人中心
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <!-- 主内容区 -->
    <div class="layout-main">
      <!-- 侧边菜单 -->
      <aside class="layout-aside">
        <el-menu
            :default-active="activeMenu"
            router
            class="sidebar-menu"
        >
          <el-menu-item
              index="/system/user"
              v-if="userStore.hasPermission('system:user:view')"
          >
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>

          <el-menu-item
              index="/system/role"
              v-if="userStore.hasPermission('system:role:view')"
          >
            <el-icon><Lock /></el-icon>
            <span>角色管理</span>
          </el-menu-item>
        </el-menu>
      </aside>

      <!-- 内容区域 -->
      <main class="layout-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  ArrowDown,
  User,
  SwitchButton,
  HomeFilled,
  Lock
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const handleCommand = async (command) => {
  if (command === 'logout') {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await userStore.logout()
    router.push('/login')
  } else if (command === 'profile') {
    ElMessage.info('个人中心功能开发中...')
  }
}
</script>

<style scoped>
.layout-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

/* ========== 顶部导航栏 ========== */
.layout-header {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  flex-shrink: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 28px;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: all 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

/* ========== 主内容区 ========== */
.layout-main {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
}

/* ========== 侧边菜单 ========== */
.layout-aside {
  width: 220px;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  overflow-y: auto;
  flex-shrink: 0;
}

.sidebar-menu {
  border-right: none;
  height: 100%;
}

.sidebar-menu .el-menu-item {
  height: 50px;
  line-height: 50px;
  font-size: 14px;
  color: #333;
  transition: all 0.3s;
}

.sidebar-menu .el-menu-item:hover {
  background: #f0f5ff;
  color: #1890ff;
}

.sidebar-menu .el-menu-item.is-active {
  background: linear-gradient(90deg, #e6f7ff 0%, transparent 100%);
  color: #1890ff;
  border-right: 3px solid #1890ff;
}

.sidebar-menu .el-menu-item .el-icon {
  margin-right: 8px;
  font-size: 16px;
}

/* ========== 内容区域 ========== */
.layout-content {
  flex: 1;
  padding: 20px;
  background: #f5f7fa;
  overflow-y: auto;
  min-height: 0;
}

/* ========== 滚动条样式 ========== */
.layout-aside::-webkit-scrollbar,
.layout-content::-webkit-scrollbar {
  width: 6px;
}

.layout-aside::-webkit-scrollbar-track,
.layout-content::-webkit-scrollbar-track {
  background: transparent;
}

.layout-aside::-webkit-scrollbar-thumb,
.layout-content::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 3px;
}

.layout-aside::-webkit-scrollbar-thumb:hover,
.layout-content::-webkit-scrollbar-thumb:hover {
  background: #bfbfbf;
}

/* ========== 响应式 ========== */
@media (max-width: 768px) {
  .layout-aside {
    width: 180px;
  }

  .layout-content {
    padding: 16px;
  }
}
</style>
