<!-- 区域视图 - 按区域/IDC维度展示主机 -->
<template>
  <div class="region-view">
    <el-card class="view-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="title-icon"><Location /></el-icon>
            <h2>区域视图</h2>
            <el-tag type="info" size="small">按区域/IDC分类查看资源</el-tag>
          </div>
          <div class="header-right">
            <el-select v-model="groupBy" size="small" style="width: 120px">
              <el-option label="按区域" value="region" />
              <el-option label="按IDC" value="idc" />
            </el-select>
          </div>
        </div>
      </template>

      <div class="region-container" v-loading="loading">
        <!-- 地图概览 -->
        <div class="map-overview">
          <div ref="mapChartRef" class="map-chart"></div>
        </div>

        <!-- 区域/IDC列表 -->
        <div class="region-list">
          <div 
            v-for="item in groupedData" 
            :key="item.name" 
            class="region-card"
            @click="showRegionHosts(item.name)"
          >
            <div class="region-header">
              <div class="region-icon">
                <el-icon :size="20"><OfficeBuilding /></el-icon>
              </div>
              <div class="region-info">
                <h4>{{ item.name || '未分配' }}</h4>
                <span class="region-count">{{ item.hosts.length }} 台主机</span>
              </div>
            </div>

            <div class="region-stats">
              <div class="stat-bar">
                <div 
                  class="stat-fill online" 
                  :style="{ width: `${item.onlinePercent}%` }"
                ></div>
              </div>
              <div class="stat-text">
                <span class="online">{{ item.onlineCount }} 在线</span>
                <span class="offline">{{ item.offlineCount }} 离线</span>
              </div>
            </div>

            <div class="region-tags">
              <el-tag 
                v-for="env in item.environments" 
                :key="env" 
                size="small" 
                :type="getEnvType(env)"
              >
                {{ getEnvName(env) }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { getHostList } from '@/api/host'
import { Location, OfficeBuilding } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const groupBy = ref('region')
const allHosts = ref([])
const mapChartRef = ref(null)
let mapChart = null

// 按区域/IDC分组
const groupedData = computed(() => {
  const field = groupBy.value
  const groups = {}
  
  allHosts.value.forEach(host => {
    const key = host[field] || '未分配'
    if (!groups[key]) {
      groups[key] = {
        name: key,
        hosts: [],
        environments: new Set()
      }
    }
    groups[key].hosts.push(host)
    if (host.environment) {
      groups[key].environments.add(host.environment)
    }
  })

  return Object.values(groups).map(g => ({
    ...g,
    environments: Array.from(g.environments),
    onlineCount: g.hosts.filter(h => h.status === 'online').length,
    offlineCount: g.hosts.filter(h => h.status !== 'online').length,
    onlinePercent: g.hosts.length ? Math.round(g.hosts.filter(h => h.status === 'online').length / g.hosts.length * 100) : 0
  })).sort((a, b) => b.hosts.length - a.hosts.length)
})

// 获取主机列表
const fetchHosts = async () => {
  loading.value = true
  try {
    const data = await getHostList({ page: 1, page_size: 1000 })
    allHosts.value = data.list || []
    initChart()
  } catch (error) {
    console.error('获取主机列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化图表
const initChart = () => {
  if (!mapChartRef.value) return
  
  mapChart = echarts.init(mapChartRef.value)
  
  const chartData = groupedData.value.map(item => ({
    name: item.name,
    value: item.hosts.length,
    online: item.onlineCount,
    offline: item.offlineCount
  }))

  mapChart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: (params) => {
        const data = params.data
        return `${data.name}<br/>总计: ${data.value} 台<br/>在线: ${data.online} 台<br/>离线: ${data.offline} 台`
      }
    },
    series: [{
      type: 'treemap',
      data: chartData,
      label: {
        show: true,
        formatter: '{b}\n{c} 台'
      },
      itemStyle: {
        borderColor: '#fff',
        borderWidth: 2,
        gapWidth: 2
      },
      levels: [{
        itemStyle: {
          borderColor: '#fff',
          borderWidth: 2,
          gapWidth: 2
        }
      }]
    }]
  })
}

const getEnvType = (env) => {
  const types = { production: 'success', staging: 'warning', development: '' }
  return types[env] || 'info'
}

const getEnvName = (env) => {
  const names = { production: '生产', staging: '预发', development: '开发' }
  return names[env] || env
}

const showRegionHosts = (name) => {
  const query = groupBy.value === 'region' ? { region: name } : { idc: name }
  router.push({ path: '/cmdb/hosts', query })
}

const handleResize = () => {
  mapChart?.resize()
}

watch(groupBy, () => {
  initChart()
})

onMounted(() => {
  fetchHosts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  mapChart?.dispose()
})
</script>

<style scoped>
.region-view {
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

.region-container {
  display: flex;
  gap: 20px;
  min-height: 500px;
}

.map-overview {
  flex: 1;
  background: #fafafa;
  border-radius: 12px;
  padding: 16px;
}

.map-chart {
  width: 100%;
  height: 400px;
}

.region-list {
  width: 360px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 500px;
  overflow-y: auto;
}

.region-card {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 10px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s;
}

.region-card:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

.region-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.region-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #1890ff, #69c0ff);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.region-info h4 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.region-count {
  font-size: 12px;
  color: #999;
}

.region-stats {
  margin-bottom: 12px;
}

.stat-bar {
  height: 6px;
  background: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 6px;
}

.stat-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.3s;
}

.stat-fill.online {
  background: linear-gradient(90deg, #52c41a, #95de64);
}

.stat-text {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.stat-text .online { color: #52c41a; }
.stat-text .offline { color: #ff4d4f; }

.region-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

@media (max-width: 1200px) {
  .region-container {
    flex-direction: column;
  }

  .region-list {
    width: 100%;
    flex-direction: row;
    flex-wrap: wrap;
    max-height: none;
  }

  .region-card {
    flex: 1;
    min-width: 280px;
  }
}
</style>
