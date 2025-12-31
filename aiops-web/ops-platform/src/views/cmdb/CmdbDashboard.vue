<!-- CMDB总览 - 业务数据模型为核心 -->
<template>
  <div class="cmdb-dashboard">
    <!-- 欢迎区 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <div class="welcome-text">
          <h1>
            <el-icon><Coin /></el-icon>
            CMDB 配置管理数据库
          </h1>
          <p>构建业务与基础设施的关系链，实现故障影响分析、变更评估、容量规划</p>
        </div>
        <div class="welcome-actions">
          <button class="action-btn primary" @click="$router.push('/cmdb/business')">
            <el-icon><OfficeBuilding /></el-icon>
            <span class="btn-text">管理业务</span>
          </button>
          <button class="action-btn" @click="$router.push('/cmdb/view/topology')">
            <el-icon><Share /></el-icon>
            <span class="btn-text">查看拓扑</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 业务关系链概览 -->
    <div class="relation-chain">
      <div class="chain-title">
        <el-icon><Connection /></el-icon>
        <span>业务关系链</span>
        <span class="title-tip">点击卡片进入管理</span>
      </div>
      <div class="chain-flow">
        <div class="chain-card" @click="$router.push('/cmdb/business')">
          <div class="card-icon business">
            <el-icon><OfficeBuilding /></el-icon>
          </div>
          <div class="card-info">
            <div class="card-count">{{ modelStats.business }}</div>
            <div class="card-label">业务线</div>
          </div>
          <div class="card-desc">产品/业务单元</div>
        </div>
        <div class="chain-arrow">
          <el-icon><Right /></el-icon>
        </div>
        <div class="chain-card" @click="$router.push('/cmdb/applications')">
          <div class="card-icon app">
            <el-icon><Grid /></el-icon>
          </div>
          <div class="card-info">
            <div class="card-count">{{ modelStats.applications }}</div>
            <div class="card-label">应用</div>
          </div>
          <div class="card-desc">系统/应用程序</div>
        </div>
        <div class="chain-arrow">
          <el-icon><Right /></el-icon>
        </div>
        <div class="chain-card" @click="$router.push('/cmdb/services')">
          <div class="card-icon service">
            <el-icon><SetUp /></el-icon>
          </div>
          <div class="card-info">
            <div class="card-count">{{ modelStats.services }}</div>
            <div class="card-label">服务</div>
          </div>
          <div class="card-desc">微服务/模块</div>
        </div>
        <div class="chain-arrow">
          <el-icon><Right /></el-icon>
        </div>
        <div class="chain-card infra" @click="$router.push('/cmdb/hosts')">
          <div class="card-icon host">
            <el-icon><Monitor /></el-icon>
          </div>
          <div class="card-info">
            <div class="card-count">{{ hostStats.total }}</div>
            <div class="card-label">基础设施</div>
          </div>
          <div class="card-desc">主机/数据库/中间件</div>
        </div>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-section">
      <!-- 左侧：CMDB价值 + 快速入口 -->
      <div class="left-panel">
        <!-- CMDB能解决什么 -->
        <div class="panel-card value-card">
          <div class="card-title">
            <el-icon><InfoFilled /></el-icon>
            CMDB能解决什么
          </div>
          <div class="value-list">
            <div class="value-item">
              <el-icon class="value-icon"><Search /></el-icon>
              <div class="value-content">
                <div class="value-name">故障影响分析</div>
                <div class="value-desc">主机故障时，快速定位影响哪些业务</div>
              </div>
            </div>
            <div class="value-item">
              <el-icon class="value-icon"><Document /></el-icon>
              <div class="value-content">
                <div class="value-name">变更评估</div>
                <div class="value-desc">变更前评估影响范围，降低风险</div>
              </div>
            </div>
            <div class="value-item">
              <el-icon class="value-icon"><TrendCharts /></el-icon>
              <div class="value-content">
                <div class="value-name">容量规划</div>
                <div class="value-desc">基于业务增长预测资源需求</div>
              </div>
            </div>
            <div class="value-item">
              <el-icon class="value-icon"><Wallet /></el-icon>
              <div class="value-content">
                <div class="value-name">成本分摊</div>
                <div class="value-desc">按业务线精确分摊IT成本</div>
              </div>
            </div>
            <div class="value-item">
              <el-icon class="value-icon"><Cpu /></el-icon>
              <div class="value-content">
                <div class="value-name">AI上下文</div>
                <div class="value-desc">为AI分析提供业务关系数据</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 快速开始 -->
        <div class="panel-card">
          <div class="card-title">
            <el-icon><Promotion /></el-icon>
            快速开始
          </div>
          <div class="quick-list">
            <div class="quick-item" @click="$router.push('/cmdb/business')">
              <span class="quick-step">1</span>
              <div class="quick-content">
                <div class="quick-name">创建业务线</div>
                <div class="quick-desc">定义你的产品或业务单元</div>
              </div>
              <el-icon class="quick-arrow"><Right /></el-icon>
            </div>
            <div class="quick-item" @click="$router.push('/cmdb/applications')">
              <span class="quick-step">2</span>
              <div class="quick-content">
                <div class="quick-name">添加应用</div>
                <div class="quick-desc">关联业务线下的应用系统</div>
              </div>
              <el-icon class="quick-arrow"><Right /></el-icon>
            </div>
            <div class="quick-item" @click="$router.push('/cmdb/services')">
              <span class="quick-step">3</span>
              <div class="quick-content">
                <div class="quick-name">配置服务</div>
                <div class="quick-desc">定义应用的微服务/模块</div>
              </div>
              <el-icon class="quick-arrow"><Right /></el-icon>
            </div>
            <div class="quick-item" @click="$router.push('/cmdb/hosts')">
              <span class="quick-step">4</span>
              <div class="quick-content">
                <div class="quick-name">绑定资源</div>
                <div class="quick-desc">将主机/数据库绑定到服务</div>
              </div>
              <el-icon class="quick-arrow"><Right /></el-icon>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：基础设施概览 + 最近变更 -->
      <div class="right-panel">
        <!-- 基础设施状态 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">
              <el-icon><Monitor /></el-icon>
              基础设施状态
            </span>
            <span class="card-action" @click="$router.push('/cmdb/hosts')">查看全部 ></span>
          </div>
          <div class="infra-stats">
            <div class="infra-item" @click="$router.push('/cmdb/hosts')">
              <div class="infra-icon host">
                <el-icon><Monitor /></el-icon>
              </div>
              <div class="infra-info">
                <div class="infra-count">{{ hostStats.total }}</div>
                <div class="infra-label">主机</div>
              </div>
              <div class="infra-status">
                <span class="status online">{{ hostStats.online }} 在线</span>
                <span class="status offline" v-if="hostStats.offline > 0">{{ hostStats.offline }} 离线</span>
              </div>
            </div>
            <div class="infra-item">
              <div class="infra-icon db">
                <el-icon><Coin /></el-icon>
              </div>
              <div class="infra-info">
                <div class="infra-count">{{ infraStats.databases }}</div>
                <div class="infra-label">数据库</div>
              </div>
              <div class="infra-status">
                <span class="status online">{{ infraStats.databases }} 正常</span>
              </div>
            </div>
            <div class="infra-item">
              <div class="infra-icon cache">
                <el-icon><Lightning /></el-icon>
              </div>
              <div class="infra-info">
                <div class="infra-count">{{ infraStats.caches }}</div>
                <div class="infra-label">缓存</div>
              </div>
              <div class="infra-status">
                <span class="status online">{{ infraStats.caches }} 正常</span>
              </div>
            </div>
            <div class="infra-item">
              <div class="infra-icon mq">
                <el-icon><Message /></el-icon>
              </div>
              <div class="infra-info">
                <div class="infra-count">{{ infraStats.mqs }}</div>
                <div class="infra-label">消息队列</div>
              </div>
              <div class="infra-status">
                <span class="status online">{{ infraStats.mqs }} 正常</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 环境分布 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">
              <el-icon><PieChart /></el-icon>
              环境分布
            </span>
          </div>
          <div ref="chartRef" class="chart-container"></div>
        </div>

        <!-- 最近变更 -->
        <div class="panel-card">
          <div class="card-header">
            <span class="card-title">
              <el-icon><List /></el-icon>
              最近变更
            </span>
            <span class="card-action" @click="$router.push('/cmdb/changes')">查看全部 ></span>
          </div>
          <div class="change-list">
            <div v-for="item in recentChanges" :key="item.id" class="change-item">
              <el-icon class="change-icon" :class="item.type"><component :is="item.iconComponent" /></el-icon>
              <div class="change-content">
                <div class="change-text">{{ item.text }}</div>
                <div class="change-time">{{ item.time }}</div>
              </div>
            </div>
            <div v-if="recentChanges.length === 0" class="empty-tip">暂无变更记录</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { getHostStats } from '@/api/host'
import {
  Coin, OfficeBuilding, Share, Connection, Right, Grid, SetUp, Monitor,
  InfoFilled, Search, Document, TrendCharts, Wallet, Cpu, Promotion,
  Lightning, Message, PieChart, List, Plus, Edit, Refresh
} from '@element-plus/icons-vue'

const router = useRouter()

// 业务模型统计
const modelStats = ref({
  business: 3,
  applications: 8,
  services: 24
})

// 主机统计
const hostStats = ref({ total: 0, online: 0, offline: 0 })

// 基础设施统计
const infraStats = ref({
  databases: 3,
  caches: 2,
  mqs: 1
})

// 图表
const chartRef = ref(null)
let chart = null

// 最近变更
const recentChanges = ref([
  { id: 1, type: 'add', iconComponent: markRaw(Plus), text: '新增业务线「电商平台」', time: '10分钟前' },
  { id: 2, type: 'edit', iconComponent: markRaw(Edit), text: '应用「订单系统」关联服务', time: '1小时前' },
  { id: 3, type: 'link', iconComponent: markRaw(Connection), text: '主机 web-server-01 绑定到「用户服务」', time: '2小时前' },
  { id: 4, type: 'update', iconComponent: markRaw(Refresh), text: '服务「支付网关」配置更新', time: '3小时前' }
])

const fetchData = async () => {
  try {
    const data = await getHostStats()
    hostStats.value = data || { total: 10, online: 8, offline: 2 }
    initChart()
  } catch (error) {
    console.error('获取数据失败:', error)
    hostStats.value = { total: 10, online: 8, offline: 2 }
    initChart()
  }
}

const initChart = () => {
  if (!chartRef.value) return
  chart = echarts.init(chartRef.value)

  const data = [
    { name: '生产环境', value: 7, itemStyle: { color: '#52c41a' } },
    { name: '预发环境', value: 2, itemStyle: { color: '#faad14' } },
    { name: '开发环境', value: 1, itemStyle: { color: '#1890ff' } }
  ]

  chart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    series: [{
      type: 'pie',
      radius: ['45%', '70%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 4, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, position: 'outside', formatter: '{b}\n{c}台' },
      data
    }]
  })
}

const handleResize = () => chart?.resize()

onMounted(() => {
  fetchData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
})
</script>

<style scoped>
.cmdb-dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
}

/* ========== 欢迎区 ========== */
.welcome-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  padding: 24px 32px;
  color: #fff;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.welcome-text h1 {
  margin: 0 0 8px;
  font-size: 22px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 10px;
}

.welcome-text p {
  margin: 0;
  opacity: 0.9;
  font-size: 14px;
}

.welcome-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border: 1px solid rgba(255,255,255,0.3);
  background: rgba(255,255,255,0.1);
  border-radius: 8px;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn:hover {
  background: rgba(255,255,255,0.2);
  transform: translateY(-2px);
}

.action-btn.primary {
  background: #fff;
  color: #667eea;
  border-color: #fff;
  font-weight: 500;
}

.action-btn.primary:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

/* ========== 业务关系链 ========== */
.relation-chain {
  background: #fff;
  border-radius: 12px;
  padding: 20px 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.chain-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.chain-title .el-icon {
  color: #667eea;
}

.title-tip {
  margin-left: auto;
  font-size: 12px;
  color: #999;
  font-weight: normal;
}

.chain-flow {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.chain-card {
  flex: 1;
  max-width: 200px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
}

.chain-card:hover {
  background: #f0f5ff;
  border-color: #667eea;
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
}

.chain-card.infra {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08) 0%, rgba(118, 75, 162, 0.08) 100%);
}

.card-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 12px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.card-icon.business { background: #e6f7ff; color: #1890ff; }
.card-icon.app { background: #f6ffed; color: #52c41a; }
.card-icon.service { background: #fff7e6; color: #faad14; }
.card-icon.host { background: #f0f5ff; color: #667eea; }

.card-info {
  margin-bottom: 8px;
}

.card-count {
  font-size: 28px;
  font-weight: 700;
  color: #333;
}

.card-label {
  font-size: 14px;
  color: #666;
}

.card-desc {
  font-size: 12px;
  color: #999;
}

.chain-arrow {
  font-size: 20px;
  color: #d9d9d9;
  flex-shrink: 0;
}

/* ========== 主内容区 ========== */
.main-section {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  min-height: 0;
}

.left-panel,
.right-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.panel-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title .el-icon {
  color: #667eea;
}

.card-header .card-title {
  margin-bottom: 0;
}

.card-action {
  font-size: 13px;
  color: #667eea;
  cursor: pointer;
}

.card-action:hover {
  text-decoration: underline;
}

/* ========== CMDB价值 ========== */
.value-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.value-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  background: #fafafa;
  transition: all 0.3s;
}

.value-item:hover {
  background: #f0f5ff;
}

.value-icon {
  font-size: 20px;
  color: #667eea;
  flex-shrink: 0;
  margin-top: 2px;
}

.value-content {
  flex: 1;
}

.value-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.value-desc {
  font-size: 12px;
  color: #999;
}

/* ========== 快速开始 ========== */
.quick-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.quick-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #f0f0f0;
}

.quick-item:hover {
  background: #f0f5ff;
  border-color: #667eea;
}

.quick-step {
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.quick-content {
  flex: 1;
}

.quick-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.quick-desc {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}

.quick-arrow {
  color: #999;
  font-size: 16px;
  transition: all 0.3s;
}

.quick-item:hover .quick-arrow {
  color: #667eea;
  transform: translateX(4px);
}

/* ========== 基础设施状态 ========== */
.infra-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.infra-item {
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.infra-item:hover {
  background: #f0f5ff;
}

.infra-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  margin-bottom: 8px;
}

.infra-icon.host { background: #f0f5ff; color: #667eea; }
.infra-icon.db { background: #fff7e6; color: #faad14; }
.infra-icon.cache { background: #f6ffed; color: #52c41a; }
.infra-icon.mq { background: #e6f7ff; color: #1890ff; }

.infra-info {
  margin-bottom: 8px;
}

.infra-count {
  font-size: 24px;
  font-weight: 700;
  color: #333;
}

.infra-label {
  font-size: 12px;
  color: #999;
}

.infra-status {
  display: flex;
  gap: 8px;
}

.status {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
}

.status.online {
  background: #f6ffed;
  color: #52c41a;
}

.status.offline {
  background: #fff2f0;
  color: #ff4d4f;
}

/* ========== 图表 ========== */
.chart-container {
  height: 180px;
}

/* ========== 变更列表 ========== */
.change-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.change-item {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.change-icon {
  font-size: 16px;
  color: #667eea;
  margin-top: 2px;
}

.change-icon.add { color: #52c41a; }
.change-icon.edit { color: #1890ff; }
.change-icon.link { color: #667eea; }
.change-icon.update { color: #faad14; }

.change-content {
  flex: 1;
}

.change-text {
  font-size: 13px;
  color: #333;
}

.change-time {
  font-size: 11px;
  color: #999;
  margin-top: 2px;
}

.empty-tip {
  text-align: center;
  color: #999;
  font-size: 13px;
  padding: 20px;
}

/* ========== 响应式 ========== */
@media (max-width: 1200px) {
  .chain-flow {
    flex-wrap: wrap;
  }
  
  .chain-arrow {
    display: none;
  }
  
  .chain-card {
    min-width: 140px;
  }
  
  .main-section {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .welcome-content {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
}
</style>
