<!-- 环境视图 - 按环境维度展示主机 -->
<template>
  <div class="environment-view">
    <el-card class="view-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="title-icon"><OfficeBuilding /></el-icon>
            <h2>环境视图</h2>
            <el-tag type="info" size="small">按环境分类查看资源</el-tag>
          </div>
          <div class="header-right">
            <el-radio-group v-model="viewMode" size="small">
              <el-radio-button label="card">卡片</el-radio-button>
              <el-radio-button label="table">表格</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </template>

      <div class="env-container" v-loading="loading">
        <!-- 卡片视图 -->
        <template v-if="viewMode === 'card'">
          <div class="env-grid">
            <div 
              v-for="env in environments" 
              :key="env.key" 
              class="env-card"
              :class="env.key"
            >
              <div class="env-header">
                <div class="env-icon" :class="env.key">
                  <el-icon :size="24">
                    <component :is="env.icon" />
                  </el-icon>
                </div>
                <div class="env-info">
                  <h3>{{ env.name }}</h3>
                  <span class="env-desc">{{ env.description }}</span>
                </div>
                <div class="env-count">
                  <span class="count-value">{{ env.hosts.length }}</span>
                  <span class="count-label">台主机</span>
                </div>
              </div>

              <div class="env-stats">
                <div class="stat-item">
                  <span class="stat-dot online"></span>
                  <span>在线 {{ env.onlineCount }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-dot offline"></span>
                  <span>离线 {{ env.offlineCount }}</span>
                </div>
              </div>

              <div class="env-hosts">
                <div 
                  v-for="host in env.hosts.slice(0, 5)" 
                  :key="host.id" 
                  class="host-item"
                  @click="goToHost(host.id)"
                >
                  <div class="host-status" :class="host.status"></div>
                  <div class="host-info">
                    <span class="host-name">{{ host.hostname }}</span>
                    <span class="host-ip">{{ host.ip }}</span>
                  </div>
                  <el-tag size="small" :type="host.status === 'online' ? 'success' : 'danger'">
                    {{ host.status === 'online' ? '在线' : '离线' }}
                  </el-tag>
                </div>
                <div 
                  v-if="env.hosts.length > 5" 
                  class="more-hosts"
                  @click="showEnvHosts(env.key)"
                >
                  查看全部 {{ env.hosts.length }} 台主机
                  <el-icon><ArrowRight /></el-icon>
                </div>
                <el-empty v-if="env.hosts.length === 0" description="暂无主机" :image-size="40" />
              </div>
            </div>
          </div>
        </template>

        <!-- 表格视图 -->
        <template v-else>
          <el-tabs v-model="activeEnv" type="card">
            <el-tab-pane 
              v-for="env in environments" 
              :key="env.key" 
              :label="`${env.name} (${env.hosts.length})`" 
              :name="env.key"
            >
              <el-table :data="env.hosts" style="width: 100%">
                <el-table-column prop="hostname" label="主机名" width="180" />
                <el-table-column prop="ip" label="IP地址" width="150" />
                <el-table-column prop="os" label="操作系统" width="120" />
                <el-table-column prop="cpu" label="CPU" width="80" />
                <el-table-column prop="memory" label="内存(MB)" width="100" />
                <el-table-column label="状态" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 'online' ? 'success' : 'danger'" size="small">
                      {{ row.status === 'online' ? '在线' : '离线' }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="120">
                  <template #default="{ row }">
                    <el-button link type="primary" size="small" @click="goToHost(row.id)">
                      详情
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
          </el-tabs>
        </template>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getHostList } from '@/api/host'
import { OfficeBuilding, Promotion, SetUp, ArrowRight } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const viewMode = ref('card')
const activeEnv = ref('production')
const allHosts = ref([])

// 环境配置
const envConfig = [
  { key: 'production', name: '生产环境', description: '线上正式环境', icon: 'Promotion' },
  { key: 'staging', name: '预发环境', description: '预发布测试环境', icon: 'SetUp' },
  { key: 'development', name: '开发环境', description: '开发调试环境', icon: 'OfficeBuilding' }
]

// 按环境分组的主机
const environments = computed(() => {
  return envConfig.map(env => {
    const hosts = allHosts.value.filter(h => h.environment === env.key)
    return {
      ...env,
      hosts,
      onlineCount: hosts.filter(h => h.status === 'online').length,
      offlineCount: hosts.filter(h => h.status !== 'online').length
    }
  })
})

// 获取主机列表
const fetchHosts = async () => {
  loading.value = true
  try {
    const data = await getHostList({ page: 1, page_size: 1000 })
    allHosts.value = data.list || []
  } catch (error) {
    console.error('获取主机列表失败:', error)
  } finally {
    loading.value = false
  }
}

const goToHost = (id) => {
  router.push(`/cmdb/hosts?id=${id}`)
}

const showEnvHosts = (env) => {
  router.push({ path: '/cmdb/hosts', query: { environment: env } })
}

onMounted(() => {
  fetchHosts()
})
</script>

<style scoped>
.environment-view {
  height: 100%;
}

.view-card {
  height: 100%;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 24px;
  color: #1890ff;
}

.header-left h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.env-container {
  min-height: 400px;
}

/* 卡片视图 */
.env-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.env-card {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  padding: 20px;
  transition: all 0.3s;
}

.env-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.env-card.production { border-top: 3px solid #52c41a; }
.env-card.staging { border-top: 3px solid #faad14; }
.env-card.development { border-top: 3px solid #1890ff; }

.env-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.env-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.env-icon.production { background: linear-gradient(135deg, #52c41a, #95de64); }
.env-icon.staging { background: linear-gradient(135deg, #faad14, #ffc069); }
.env-icon.development { background: linear-gradient(135deg, #1890ff, #69c0ff); }

.env-info h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.env-desc {
  font-size: 12px;
  color: #999;
}

.env-count {
  margin-left: auto;
  text-align: right;
}

.count-value {
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.count-label {
  display: block;
  font-size: 12px;
  color: #999;
}

.env-stats {
  display: flex;
  gap: 20px;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #666;
}

.stat-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.stat-dot.online { background: #52c41a; }
.stat-dot.offline { background: #ff4d4f; }

.env-hosts {
  max-height: 240px;
  overflow-y: auto;
}

.host-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.host-item:hover {
  background: #f5f7fa;
}

.host-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.host-status.online { background: #52c41a; }
.host-status.offline { background: #ff4d4f; }

.host-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.host-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.host-ip {
  font-size: 12px;
  color: #999;
}

.more-hosts {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 12px;
  color: #1890ff;
  font-size: 13px;
  cursor: pointer;
  border-top: 1px dashed #e8e8e8;
  margin-top: 8px;
}

.more-hosts:hover {
  background: #f0f7ff;
  border-radius: 8px;
}

@media (max-width: 1200px) {
  .env-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .env-grid {
    grid-template-columns: 1fr;
  }
}
</style>
