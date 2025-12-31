<!-- CMDB管理 - 业务数据模型为核心 -->
<template>
  <div class="cmdb-page">
    <!-- 顶部导航栏 -->
    <header class="cmdb-header">
      <div class="header-left">
        <div class="logo" @click="$router.push('/')">
          <el-icon class="logo-icon"><Platform /></el-icon>
          <span class="logo-text">智能运维平台</span>
        </div>
        <span class="nav-divider">/</span>
        <span class="nav-current">
          <el-icon><Coin /></el-icon>
          CMDB管理
        </span>
      </div>

      <div class="header-right">
        <el-button type="primary" link @click="$router.push('/')">
          <el-icon><HomeFilled /></el-icon>
          返回工作台
        </el-button>
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
    <div class="cmdb-main">
      <!-- 左侧导航 -->
      <aside class="cmdb-sidebar">
        <nav class="sidebar-nav">
          <!-- 总览 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/dashboard') }"
            @click="$router.push('/cmdb/dashboard')"
          >
            <el-icon class="nav-icon"><DataAnalysis /></el-icon>
            <span class="nav-text">总览</span>
          </div>

          <div class="nav-section-title">业务模型</div>
          
          <!-- 业务线 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/business') }"
            @click="$router.push('/cmdb/business')"
          >
            <el-icon class="nav-icon"><OfficeBuilding /></el-icon>
            <span class="nav-text">业务线</span>
          </div>
          
          <!-- 应用 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/applications') }"
            @click="$router.push('/cmdb/applications')"
          >
            <el-icon class="nav-icon"><Grid /></el-icon>
            <span class="nav-text">应用</span>
          </div>
          
          <!-- 服务 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/services') }"
            @click="$router.push('/cmdb/services')"
          >
            <el-icon class="nav-icon"><SetUp /></el-icon>
            <span class="nav-text">服务</span>
          </div>

          <div class="nav-section-title">基础设施</div>
          
          <!-- 主机 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/hosts') }"
            @click="$router.push('/cmdb/hosts')"
          >
            <el-icon class="nav-icon"><Monitor /></el-icon>
            <span class="nav-text">主机</span>
            <span class="nav-badge" v-if="stats.total">{{ stats.total }}</span>
          </div>
          
          <!-- 分组 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/host-groups') }"
            @click="$router.push('/cmdb/host-groups')"
          >
            <el-icon class="nav-icon"><Folder /></el-icon>
            <span class="nav-text">分组</span>
          </div>

          <div class="nav-divider-line"></div>
          
          <!-- 拓扑视图 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/view/') }"
            @click="$router.push('/cmdb/view/topology')"
          >
            <el-icon class="nav-icon"><Share /></el-icon>
            <span class="nav-text">拓扑</span>
          </div>
          
          <!-- 导入 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/import') }"
            @click="$router.push('/cmdb/import')"
          >
            <el-icon class="nav-icon"><Upload /></el-icon>
            <span class="nav-text">导入</span>
          </div>
          
          <!-- 变更 -->
          <div 
            class="nav-item" 
            :class="{ active: isActive('/cmdb/changes') }"
            @click="$router.push('/cmdb/changes')"
          >
            <el-icon class="nav-icon"><List /></el-icon>
            <span class="nav-text">变更</span>
          </div>
        </nav>

        <!-- 底部状态 -->
        <div class="sidebar-footer">
          <div class="status-item online">
            <span class="status-dot"></span>
            <span>{{ stats.online || 0 }} 在线</span>
          </div>
          <div class="status-item offline">
            <span class="status-dot"></span>
            <span>{{ stats.offline || 0 }} 离线</span>
          </div>
        </div>
      </aside>

      <!-- 内容区域 -->
      <main class="cmdb-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessageBox, ElMessage } from 'element-plus'
import { getHostStats } from '@/api/host'
import { 
  ArrowDown, User, SwitchButton, HomeFilled, Platform, Coin,
  DataAnalysis, OfficeBuilding, Grid, SetUp, Monitor, Folder,
  Share, Upload, List
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const stats = ref({ total: 0, online: 0, offline: 0 })

const isActive = (path) => route.path === path || route.path.startsWith(path)

const fetchStats = async () => {
  try {
    const data = await getHostStats()
    stats.value = data || { total: 0, online: 0, offline: 0 }
  } catch (error) {
    console.error('获取统计失败:', error)
  }
}

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

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.cmdb-page {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

/* ========== 顶部导航 ========== */
.cmdb-header {
  height: 60px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.logo-icon {
  font-size: 24px;
  color: #667eea;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.nav-divider {
  color: #d9d9d9;
  font-size: 18px;
}

.nav-current {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 15px;
  color: #333;
  font-weight: 500;
}

.nav-current .el-icon {
  color: #667eea;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 8px;
  transition: background 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #333;
}

/* ========== 主内容区 ========== */
.cmdb-main {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* ========== 左侧导航 ========== */
.cmdb-sidebar {
  width: 90px;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
}

.nav-section-title {
  font-size: 10px;
  color: #999;
  padding: 12px 8px 4px;
  text-align: center;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.nav-item:hover {
  background: #f0f5ff;
}

.nav-item.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15) 0%, rgba(118, 75, 162, 0.15) 100%);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 0 2px 2px 0;
}

.nav-icon {
  font-size: 20px;
  color: #666;
}

.nav-item.active .nav-icon {
  color: #667eea;
}

.nav-text {
  font-size: 11px;
  color: #666;
}

.nav-item.active .nav-text {
  color: #667eea;
  font-weight: 500;
}

.nav-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  background: #667eea;
  color: #fff;
  font-size: 9px;
  padding: 1px 4px;
  border-radius: 6px;
  min-width: 16px;
  text-align: center;
}

.nav-divider-line {
  height: 1px;
  background: #f0f0f0;
  margin: 8px 12px;
}

/* ========== 底部状态 ========== */
.sidebar-footer {
  padding: 12px 8px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.status-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 10px;
  color: #666;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-item.online .status-dot { background: #52c41a; }
.status-item.offline .status-dot { background: #ff4d4f; }

/* ========== 内容区域 ========== */
.cmdb-content {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}
</style>
