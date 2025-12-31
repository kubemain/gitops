<!-- 拓扑视图 - 展示主机关系拓扑 -->
<template>
  <div class="topology-view">
    <el-card class="view-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="title-icon"><Share /></el-icon>
            <h2>拓扑视图</h2>
            <el-tag type="info" size="small">可视化资源关系</el-tag>
          </div>
          <div class="header-right">
            <el-button-group size="small">
              <el-button :type="layoutType === 'force' ? 'primary' : ''" @click="changeLayout('force')">
                力导向
              </el-button>
              <el-button :type="layoutType === 'circular' ? 'primary' : ''" @click="changeLayout('circular')">
                环形
              </el-button>
            </el-button-group>
            <el-button size="small" @click="resetZoom">
              <el-icon><RefreshRight /></el-icon>
              重置
            </el-button>
          </div>
        </div>
      </template>

      <div class="topology-container" v-loading="loading">
        <div ref="topologyChartRef" class="topology-chart"></div>

        <!-- 图例 -->
        <div class="topology-legend">
          <div class="legend-title">图例说明</div>
          <div class="legend-item">
            <span class="legend-dot group"></span>
            <span>主机分组</span>
          </div>
          <div class="legend-item">
            <span class="legend-dot online"></span>
            <span>在线主机</span>
          </div>
          <div class="legend-item">
            <span class="legend-dot offline"></span>
            <span>离线主机</span>
          </div>
        </div>

        <!-- 节点详情面板 -->
        <transition name="slide">
          <div v-if="selectedNode" class="node-detail-panel">
            <div class="panel-header">
              <span class="panel-title">{{ selectedNode.name }}</span>
              <el-button link @click="selectedNode = null">
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
            <div class="panel-content">
              <template v-if="selectedNode.type === 'host'">
                <div class="detail-item">
                  <span class="label">IP地址</span>
                  <span class="value">{{ selectedNode.data.ip }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">操作系统</span>
                  <span class="value">{{ selectedNode.data.os }}</span>
                </div>
                <div class="detail-item">
                  <span class="label">状态</span>
                  <el-tag :type="selectedNode.data.status === 'online' ? 'success' : 'danger'" size="small">
                    {{ selectedNode.data.status === 'online' ? '在线' : '离线' }}
                  </el-tag>
                </div>
                <div class="detail-item">
                  <span class="label">环境</span>
                  <span class="value">{{ selectedNode.data.environment || '-' }}</span>
                </div>
                <el-button type="primary" size="small" @click="goToHost(selectedNode.data.id)">
                  查看详情
                </el-button>
              </template>
              <template v-else>
                <div class="detail-item">
                  <span class="label">主机数量</span>
                  <span class="value">{{ selectedNode.data.host_count || 0 }} 台</span>
                </div>
                <div class="detail-item">
                  <span class="label">描述</span>
                  <span class="value">{{ selectedNode.data.description || '-' }}</span>
                </div>
                <el-button type="primary" size="small" @click="goToGroup(selectedNode.data.id)">
                  查看主机
                </el-button>
              </template>
            </div>
          </div>
        </transition>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { getHostList } from '@/api/host'
import { getAllHostGroups } from '@/api/host-group'
import { Share, RefreshRight, Close } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const layoutType = ref('force')
const selectedNode = ref(null)
const topologyChartRef = ref(null)
let topologyChart = null

const allHosts = ref([])
const hostGroups = ref([])

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const [hostsData, groupsData] = await Promise.all([
      getHostList({ page: 1, page_size: 1000 }),
      getAllHostGroups()
    ])
    allHosts.value = hostsData.list || []
    hostGroups.value = groupsData || []
    initTopology()
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化拓扑图
const initTopology = () => {
  if (!topologyChartRef.value) return

  topologyChart = echarts.init(topologyChartRef.value)

  // 构建节点数据
  const nodes = []
  const links = []
  const categories = [
    { name: '主机分组', itemStyle: { color: '#722ed1' } },
    { name: '在线主机', itemStyle: { color: '#52c41a' } },
    { name: '离线主机', itemStyle: { color: '#ff4d4f' } }
  ]

  // 添加分组节点
  hostGroups.value.forEach(group => {
    nodes.push({
      id: `group_${group.id}`,
      name: group.name,
      symbolSize: 50,
      category: 0,
      type: 'group',
      data: group,
      label: { show: true }
    })
  })

  // 添加主机节点和连线
  allHosts.value.forEach(host => {
    const isOnline = host.status === 'online'
    nodes.push({
      id: `host_${host.id}`,
      name: host.hostname,
      symbolSize: 30,
      category: isOnline ? 1 : 2,
      type: 'host',
      data: host,
      label: { show: false }
    })

    // 连接到分组
    if (host.group_id) {
      links.push({
        source: `group_${host.group_id}`,
        target: `host_${host.id}`
      })
    }
  })

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: (params) => {
        if (params.dataType === 'node') {
          const node = params.data
          if (node.type === 'host') {
            return `${node.name}<br/>IP: ${node.data.ip}<br/>状态: ${node.data.status === 'online' ? '在线' : '离线'}`
          }
          return `${node.name}<br/>主机数: ${node.data.host_count || 0}`
        }
        return ''
      }
    },
    legend: {
      data: categories.map(c => c.name),
      bottom: 10
    },
    series: [{
      type: 'graph',
      layout: layoutType.value,
      data: nodes,
      links: links,
      categories: categories,
      roam: true,
      draggable: true,
      force: {
        repulsion: 200,
        edgeLength: [80, 150]
      },
      circular: {
        rotateLabel: true
      },
      emphasis: {
        focus: 'adjacency',
        lineStyle: { width: 3 }
      },
      lineStyle: {
        color: 'source',
        curveness: 0.1,
        opacity: 0.6
      }
    }]
  }

  topologyChart.setOption(option)

  // 点击事件
  topologyChart.on('click', (params) => {
    if (params.dataType === 'node') {
      selectedNode.value = params.data
    }
  })
}

const changeLayout = (type) => {
  layoutType.value = type
  initTopology()
}

const resetZoom = () => {
  topologyChart?.dispatchAction({
    type: 'restore'
  })
}

const goToHost = (id) => {
  router.push(`/cmdb/hosts?id=${id}`)
}

const goToGroup = (id) => {
  router.push({ path: '/cmdb/hosts', query: { group_id: id } })
}

const handleResize = () => {
  topologyChart?.resize()
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  topologyChart?.dispose()
})
</script>

<style scoped>
.topology-view {
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

.header-right {
  display: flex;
  gap: 12px;
}

.topology-container {
  position: relative;
  height: calc(100vh - 240px);
  min-height: 500px;
}

.topology-chart {
  width: 100%;
  height: 100%;
}

/* 图例 */
.topology-legend {
  position: absolute;
  top: 16px;
  left: 16px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 8px;
  padding: 12px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.legend-title {
  font-size: 13px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #666;
  margin-bottom: 4px;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-dot.group { background: #722ed1; }
.legend-dot.online { background: #52c41a; }
.legend-dot.offline { background: #ff4d4f; }

/* 节点详情面板 */
.node-detail-panel {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 280px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.panel-title {
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.panel-content {
  padding: 16px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.detail-item .label {
  font-size: 13px;
  color: #999;
}

.detail-item .value {
  font-size: 13px;
  color: #333;
}

/* 动画 */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>
