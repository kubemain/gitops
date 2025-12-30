<template>
  <div class="home-page">
    <!-- 顶部导航栏 -->
    <header class="home-header">
      <div class="header-left">
        <div class="logo">
          <span class="logo-icon">⚡</span>
          <span class="logo-text">智能运维平台</span>
        </div>
      </div>

      <div class="header-right">
        <!-- 快捷入口 -->
        <el-button
            type="primary"
            link
            @click="$router.push('/system/user')"
            v-if="userStore.hasPermission('system:user:view')"
        >
          <el-icon><Setting /></el-icon>
          系统管理
        </el-button>

        <!-- 用户信息 -->
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
              <el-dropdown-item command="settings">
                <el-icon><Setting /></el-icon>
                系统设置
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
    <div class="page-content">
      <!-- 左侧：节点图谱 -->
      <div class="graph-container">
        <div class="graph-toolbar">
          <button class="tool-btn" @click="resetGraph" title="重置图谱布局">
            <i>🔄</i> 重置
          </button>
          <button class="tool-btn" @click="fitView" title="自适应画布大小">
            <i>⛶</i> 适配
          </button>
          <button
              class="tool-btn"
              :class="{ active: showLabels }"
              @click="toggleLabels"
              title="切换标签显示"
          >
            <i>🏷️</i> 标签
          </button>
          <div class="toolbar-divider"></div>
          <span class="toolbar-info">节点: {{ nodeCount }} | 连线: {{ edgeCount }}</span>
        </div>
        <div ref="chartDom" class="chart-canvas"></div>

        <!-- 图例 -->
        <div class="graph-legend">
          <div class="legend-item">
            <span class="legend-dot" style="background: linear-gradient(135deg, #667eea, #764ba2)"></span>
            <span>AI大脑</span>
          </div>
          <div class="legend-item">
            <span class="legend-dot" style="background: #5B8FF9"></span>
            <span>核心功能</span>
          </div>
          <div class="legend-item">
            <span class="legend-dot" style="background: #91d5ff"></span>
            <span>子功能</span>
          </div>
        </div>

        <!-- AI对话浮动按钮 -->
        <div class="ai-float-button" @click="openAIDialog" title="打开AI智能助手">
          <div class="ai-float-icon">
            <span class="ai-icon">🤖</span>
            <span class="ai-pulse"></span>
          </div>
          <span class="ai-float-text">AI助手</span>
        </div>
      </div>

      <!-- 右侧面板 -->
      <aside class="data-panel">
        <!-- 欢迎卡片 -->
        <div class="panel-card welcome-card">
          <div class="welcome-content">
            <div class="welcome-greeting">
              <h2>👋 {{ greeting }}，{{ userStore.nickname || userStore.username }}！</h2>
              <p class="welcome-time">{{ currentTime }}</p>
            </div>
            <div class="welcome-weather">
              <span class="weather-icon">☀️</span>
              <span class="weather-temp">28°C</span>
            </div>
          </div>
          <div class="quick-stats">
            <div class="stat-item">
              <span class="stat-number">3</span>
              <span class="stat-label">待处理</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">8</span>
              <span class="stat-label">今日任务</span>
            </div>
            <div class="stat-item">
              <span class="stat-number">98%</span>
              <span class="stat-label">健康度</span>
            </div>
          </div>
        </div>

        <!-- 快捷操作 -->
        <div class="panel-card quick-actions-card">
          <h3 class="card-title">
            <i>⚡</i> 快捷操作
          </h3>
          <div class="quick-actions">
            <button class="quick-btn primary">
              <i>🖥️</i>
              <span>SSH登录</span>
            </button>
            <button class="quick-btn success">
              <i>📊</i>
              <span>查看监控</span>
            </button>
            <button class="quick-btn warning">
              <i>⚡</i>
              <span>批量执行</span>
            </button>
            <button class="quick-btn info">
              <i>🤖</i>
              <span>AI问答</span>
            </button>
          </div>
        </div>

        <!-- 核心指标 -->
        <div class="panel-card metrics-card">
          <h3 class="card-title">
            <i>📈</i> 核心指标
            <span class="refresh-btn" @click="refreshMetrics">🔄</span>
          </h3>
          <div class="metrics-list">
            <div class="metric-row">
              <div class="metric-info">
                <span class="metric-icon danger">⚠️</span>
                <div class="metric-text">
                  <span class="metric-label">严重告警</span>
                  <span class="metric-value danger">20</span>
                </div>
              </div>
              <div class="metric-trend up">+5</div>
            </div>

            <div class="metric-row">
              <div class="metric-info">
                <span class="metric-icon success">🖥️</span>
                <div class="metric-text">
                  <span class="metric-label">在线主机</span>
                  <span class="metric-value">78</span>
                </div>
              </div>
              <div class="metric-trend down">-2</div>
            </div>

            <div class="metric-row">
              <div class="metric-info">
                <span class="metric-icon info">💾</span>
                <div class="metric-text">
                  <span class="metric-label">存储使用</span>
                  <span class="metric-value">20/100T</span>
                </div>
              </div>
              <div class="metric-progress">
                <div class="progress-bar" style="width: 20%"></div>
              </div>
            </div>

            <div class="metric-row">
              <div class="metric-info">
                <span class="metric-icon warning">🔧</span>
                <div class="metric-text">
                  <span class="metric-label">待修复漏洞</span>
                  <span class="metric-value">13</span>
                </div>
              </div>
              <div class="metric-status warning">中危</div>
            </div>

            <div class="metric-row">
              <div class="metric-info">
                <span class="metric-icon">✅</span>
                <div class="metric-text">
                  <span class="metric-label">任务完成率</span>
                  <span class="metric-value">3/8</span>
                </div>
              </div>
              <div class="metric-progress">
                <div class="progress-bar success" style="width: 37.5%"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- 实时动态 -->
        <div class="panel-card activity-card">
          <h3 class="card-title">
            <i>🔥</i> 实时动态
          </h3>
          <div class="activity-list">
            <div class="activity-item">
              <div class="activity-icon ai">🤖</div>
              <div class="activity-content">
                <div class="activity-text">
                  <strong>MySQL Agent</strong> 自动优化了3条慢SQL
                </div>
                <div class="activity-time">刚刚</div>
              </div>
            </div>

            <div class="activity-item">
              <div class="activity-icon success">✅</div>
              <div class="activity-content">
                <div class="activity-text">
                  <strong>自动扩容</strong> K8s Pod 5→8
                </div>
                <div class="activity-time">1分钟前</div>
              </div>
            </div>

            <div class="activity-item">
              <div class="activity-icon user">👤</div>
              <div class="activity-content">
                <div class="activity-text">
                  <strong>李四</strong> 解决了告警 #1234
                </div>
                <div class="activity-time">3分钟前</div>
              </div>
            </div>

            <div class="activity-item">
              <div class="activity-icon warning">⚠️</div>
              <div class="activity-content">
                <div class="activity-text">
                  <strong>新告警</strong> Redis内存使用率>90%
                </div>
                <div class="activity-time">5分钟前</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 待办任务 -->
        <div class="panel-card tasks-card">
          <h3 class="card-title">
            <i>📋</i> 待办任务
            <span class="task-count">5</span>
          </h3>
          <div class="task-list">
            <div class="task-item">
              <input type="checkbox" class="task-checkbox" />
              <div class="task-content">
                <div class="task-title">
                  <span class="task-priority high">高</span>
                  处理MySQL磁盘告警
                </div>
                <div class="task-meta">
                  <span class="task-tag">运维</span>
                  <span class="task-deadline">今天 18:00</span>
                </div>
              </div>
            </div>

            <div class="task-item">
              <input type="checkbox" class="task-checkbox" />
              <div class="task-content">
                <div class="task-title">
                  <span class="task-priority medium">中</span>
                  审批财务部门迁移申请
                </div>
                <div class="task-meta">
                  <span class="task-tag">审批</span>
                  <span class="task-deadline">明天</span>
                </div>
              </div>
            </div>

            <div class="task-item completed">
              <input type="checkbox" class="task-checkbox" checked />
              <div class="task-content">
                <div class="task-title">
                  数据库日常巡检
                </div>
                <div class="task-meta">
                  <span class="task-tag">巡检</span>
                  <span class="task-done">✓ 已完成</span>
                </div>
              </div>
            </div>
          </div>
          <button class="view-all-btn">查看全部任务 →</button>
        </div>
      </aside>
    </div>

    <!-- AI对话弹窗 -->
    <div v-if="showAIDialog" class="ai-dialog-overlay" @click="closeAIDialog">
      <div class="ai-dialog" @click.stop>
        <div class="ai-dialog-header">
          <div class="ai-dialog-title">
            <span class="ai-dialog-icon">🤖</span>
            <span>AI智能助手</span>
            <span class="ai-status online">● 在线</span>
          </div>
          <button class="ai-dialog-close" @click="closeAIDialog">✕</button>
        </div>

        <div class="ai-dialog-body">
          <div class="ai-messages" ref="messagesContainer">
            <div v-for="(msg, index) in aiMessages" :key="index"
                 :class="['ai-message', msg.type]">
              <div class="message-avatar">
                {{ msg.type === 'user' ? '👤' : '🤖' }}
              </div>
              <div class="message-content">
                <div class="message-text">{{ msg.content }}</div>
                <div class="message-time">{{ msg.time }}</div>
              </div>
            </div>
          </div>

          <div class="ai-input-area">
            <input
                type="text"
                v-model="aiQuestion"
                @keyup.enter="sendAIMessage"
                placeholder="问我任何运维问题..."
                class="ai-input"
            />
            <button class="ai-send-btn" @click="sendAIMessage">
              <span>📤</span>
            </button>
          </div>

          <div class="ai-quick-questions">
            <span class="quick-label">💡 快捷问题：</span>
            <div class="quick-tags">
              <span class="quick-tag" @click="askQuickQuestion('如何查看告警？')">查看告警</span>
              <span class="quick-tag" @click="askQuickQuestion('主机状态监控')">主机状态</span>
              <span class="quick-tag" @click="askQuickQuestion('慢查询分析')">慢查询分析</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessageBox, ElMessage } from 'element-plus'
import { ArrowDown, User, Setting, SwitchButton } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const router = useRouter()
const userStore = useUserStore()

const chartDom = ref(null)
let chart = null
const aiQuestion = ref('')
const currentTime = ref('')
const nodeCount = ref(0)
const edgeCount = ref(0)
const showLabels = ref(true)
const showAIDialog = ref(false)
const aiMessages = ref([
  {
    type: 'assistant',
    content: '你好！我是AI智能助手，有什么可以帮助你的吗？',
    time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
])
const messagesContainer = ref(null)

// 计算问候语
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

// 处理用户下拉菜单
const handleCommand = async (command) => {
  if (command === 'logout') {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await userStore.logout()
    router.push('/login')
  } else if (command === 'settings') {
    router.push('/system/user')
  } else if (command === 'profile') {
    ElMessage.info('个人中心功能开发中...')
  }
}

// 打开AI对话
const openAIDialog = () => {
  showAIDialog.value = true
  nextTick(() => {
    scrollToBottom()
  })
}

// 关闭AI对话
const closeAIDialog = () => {
  showAIDialog.value = false
}

// 发送AI消息
const sendAIMessage = () => {
  if (!aiQuestion.value.trim()) return

  const userMessage = {
    type: 'user',
    content: aiQuestion.value,
    time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  }
  aiMessages.value.push(userMessage)

  // 模拟AI回复
  setTimeout(() => {
    const aiReply = {
      type: 'assistant',
      content: `我正在分析「${aiQuestion.value}」，这个问题涉及到系统监控和性能优化...`,
      time: new Date().toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
    }
    aiMessages.value.push(aiReply)
    scrollToBottom()
  }, 1000)

  aiQuestion.value = ''
  scrollToBottom()
}

// 快捷问题
const askQuickQuestion = (question) => {
  aiQuestion.value = question
  sendAIMessage()
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 节点跳转函数
const navigateTo = (nodeName) => {
  console.log('点击节点:', nodeName)

  // 系统管理相关（已开发）
  const developedRoutes = {
    '日常运维': '/system/user',
    '监控管理': '/system/user',
    '统计与分析': '/system/user',
    '流程管理': '/system/user',
    '知识管理': '/system/user'
  }

  if (developedRoutes[nodeName]) {
    router.push(developedRoutes[nodeName])
    return
  }

  // 未开发的功能提示
  ElMessage({
    message: `功能「${nodeName}」正在开发中，敬请期待...`,
    type: 'info',
    duration: 2000,
    showClose: true
  })
}

// 更新时间
const updateTime = () => {
  const now = new Date()
  const hours = now.getHours()
  const minutes = now.getMinutes().toString().padStart(2, '0')
  currentTime.value = `${hours}:${minutes}`
}

// 刷新指标
const refreshMetrics = () => {
  console.log('刷新指标数据')
  ElMessage.success('正在刷新指标数据...')
}

// 重置图谱
const resetGraph = () => {
  if (chart) {
    const option = chart.getOption()
    chart.setOption(option, true)
  }
}

// 适配视图
const fitView = () => {
  if (chart) {
    chart.resize()
  }
}

// 切换标签显示
const toggleLabels = () => {
  showLabels.value = !showLabels.value
  if (chart) {
    chart.setOption({
      series: [{
        label: {
          show: showLabels.value
        }
      }]
    })
  }
}

onMounted(() => {
  updateTime()
  setInterval(updateTime, 60000)

  chart = echarts.init(chartDom.value)

  const graphData = {
    categories: [
      { name: '中心', itemStyle: { color: '#667eea' } },
      { name: '一级', itemStyle: { color: '#5B8FF9' } },
      { name: '二级', itemStyle: { color: '#91d5ff' } }
    ],
    nodes: [
      {
        name: '智能运维大脑',
        category: 0,
        symbolSize: 120,
        label: {
          fontSize: 16,
          fontWeight: 'bold',
          color: '#fff',
          formatter: '🧠\n{b}'
        },
        itemStyle: {
          color: {
            type: 'radial',
            colorStops: [
              { offset: 0, color: '#667eea' },
              { offset: 1, color: '#764ba2' }
            ]
          },
          shadowBlur: 30,
          shadowColor: '#667eea'
        }
      },
      { name: '日常运维', category: 1, symbolSize: 80 },
      { name: '监控管理', category: 1, symbolSize: 80 },
      { name: '统计与分析', category: 1, symbolSize: 80 },
      { name: '流程管理', category: 1, symbolSize: 80 },
      { name: '知识管理', category: 1, symbolSize: 80 },
      { name: '资源管理', category: 2, symbolSize: 50 },
      { name: 'CMDB', category: 2, symbolSize: 50 },
      { name: '安全管理', category: 2, symbolSize: 50 },
      { name: '备份管理', category: 2, symbolSize: 50 },
      { name: '监控规则', category: 2, symbolSize: 50 },
      { name: '告警规则', category: 2, symbolSize: 50 },
      { name: '日志采集', category: 2, symbolSize: 50 },
      { name: '日志检索', category: 2, symbolSize: 50 },
      { name: '监控面板', category: 2, symbolSize: 50 },
      { name: '告警记录', category: 2, symbolSize: 50 },
      { name: '数据看板', category: 2, symbolSize: 50 },
      { name: '根因分析', category: 2, symbolSize: 50 },
      { name: '趋势预警', category: 2, symbolSize: 50 },
      { name: '流程看板', category: 2, symbolSize: 50 },
      { name: '服务申请', category: 2, symbolSize: 50 },
      { name: '事件处置', category: 2, symbolSize: 50 },
      { name: '问题管理', category: 2, symbolSize: 50 },
      { name: '变更管理', category: 2, symbolSize: 50 },
      { name: '文档管理', category: 2, symbolSize: 50 },
      { name: '问题库', category: 2, symbolSize: 50 },
      { name: '知识图谱', category: 2, symbolSize: 50 },
    ],
    links: []
  }

  const level1Nodes = ['日常运维', '监控管理', '统计与分析', '流程管理', '知识管理']
  level1Nodes.forEach(node => {
    graphData.links.push({
      source: '智能运维大脑',
      target: node,
      lineStyle: { width: 4, color: '#5B8FF9' }
    })
  })

  const connections = {
    '日常运维': ['资源管理', 'CMDB', '安全管理', '备份管理'],
    '监控管理': ['监控规则', '告警规则', '日志采集', '日志检索', '监控面板', '告警记录'],
    '统计与分析': ['数据看板', '根因分析', '趋势预警'],
    '流程管理': ['流程看板', '服务申请', '事件处置', '问题管理', '变更管理'],
    '知识管理': ['文档管理', '问题库', '知识图谱'],
  }

  Object.entries(connections).forEach(([parent, children]) => {
    children.forEach(child => {
      graphData.links.push({
        source: parent,
        target: child,
        lineStyle: { width: 1.5, color: '#d9d9d9' }
      })
    })
  })

  nodeCount.value = graphData.nodes.length
  edgeCount.value = graphData.links.length

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: (params) => {
        if (params.dataType === 'node') {
          return `<div style="padding: 8px;">
            <div style="font-weight: bold; margin-bottom: 4px;">${params.data.name}</div>
            <div style="font-size: 12px; color: #666;">点击进入</div>
          </div>`
        }
        return ''
      }
    },
    legend: { show: false },
    animationDuration: 1500,
    animationEasingUpdate: 'quinticInOut',
    series: [{
      type: 'graph',
      layout: 'force',
      data: graphData.nodes,
      links: graphData.links,
      categories: graphData.categories,
      roam: true,
      draggable: true,
      label: {
        show: true,
        position: 'inside',
        formatter: '{b}',
        fontSize: 11,
        color: '#333'
      },
      edgeSymbol: ['none', 'none'],
      lineStyle: {
        opacity: 0.7,
        curveness: 0.2
      },
      emphasis: {
        focus: 'adjacency',
        lineStyle: {
          width: 6,
          color: '#1890ff'
        },
        itemStyle: {
          shadowBlur: 30,
          shadowColor: '#1890ff',
          borderWidth: 4,
          borderColor: '#fff'
        },
        label: {
          fontSize: 14,
          fontWeight: 'bold'
        }
      },
      force: {
        repulsion: 800,
        gravity: 0.05,
        edgeLength: [100, 180],
        layoutAnimation: true,
        friction: 0.5
      }
    }]
  }

  chart.setOption(option)

  chart.on('click', (params) => {
    if (params.dataType === 'node') {
      navigateTo(params.data.name)
    }
  })

  const handleResize = () => {
    if (chart) {
      chart.resize()
    }
  }
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  if (chart) {
    chart.dispose()
  }
})
</script>


<style scoped>
.home-page {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

/* ========== 主内容区 ========== */
.page-content {
  flex: 1;
  display: flex;
  gap: 16px;
  overflow: hidden;
  min-height: 0;
}

/* ========== 左侧图谱 ========== */
.graph-container {
  flex: 1;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  position: relative;
}

.graph-toolbar {
  height: 48px;
  padding: 0 16px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  gap: 12px;
  background: #fafafa;
  flex-shrink: 0;
}

.tool-btn {
  padding: 6px 12px;
  border: 1px solid #d9d9d9;
  background: #fff;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.3s;
}

.tool-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.tool-btn.active {
  background: #e6f7ff;
  border-color: #1890ff;
  color: #1890ff;
}

.toolbar-divider {
  width: 1px;
  height: 20px;
  background: #e8e8e8;
  margin: 0 8px;
}

.toolbar-info {
  color: #999;
  font-size: 12px;
  margin-left: auto;
}

.chart-canvas {
  flex: 1;
  width: 100%;
  min-height: 0;
}

.graph-legend {
  position: absolute;
  bottom: 16px;
  left: 16px;
  background: rgba(255, 255, 255, 0.95);
  padding: 12px 16px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  gap: 16px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
}

.legend-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* ========== 右侧面板 ========== */
.data-panel {
  width: 380px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  flex-shrink: 0;
  padding-right: 4px;
}

.panel-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s;
  flex-shrink: 0;
}

.panel-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.card-title {
  margin: 0 0 16px 0;
  font-size: 15px;
  font-weight: 600;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title i {
  font-size: 18px;
}

.refresh-btn {
  margin-left: auto;
  cursor: pointer;
  opacity: 0.6;
  transition: all 0.3s;
  font-size: 14px;
}

.refresh-btn:hover {
  opacity: 1;
  transform: rotate(180deg);
}

/* ========== 欢迎卡片 ========== */
.welcome-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.welcome-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.welcome-greeting h2 {
  margin: 0 0 8px 0;
  font-size: 22px;
  font-weight: 600;
}

.welcome-time {
  opacity: 0.9;
  font-size: 14px;
}

.weather-icon {
  font-size: 32px;
}

.weather-temp {
  font-size: 20px;
  font-weight: 500;
  margin-left: 8px;
}

.quick-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.stat-item {
  text-align: center;
  padding: 12px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 8px;
  backdrop-filter: blur(10px);
}

.stat-number {
  display: block;
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  opacity: 0.9;
}

/* ========== 快捷操作 ========== */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.quick-btn {
  padding: 14px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.3s;
  color: #fff;
}

.quick-btn i {
  font-size: 24px;
}

.quick-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.quick-btn.success {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.quick-btn.warning {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.quick-btn.info {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.quick-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

/* ========== 核心指标 ========== */
.metrics-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.metric-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
  transition: all 0.3s;
}

.metric-row:hover {
  background: #f0f5ff;
  transform: translateX(4px);
}

.metric-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.metric-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-size: 18px;
}

.metric-icon.danger {
  background: #fff1f0;
}

.metric-icon.success {
  background: #f6ffed;
}

.metric-icon.warning {
  background: #fffbe6;
}

.metric-icon.info {
  background: #e6f7ff;
}

.metric-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.metric-label {
  font-size: 12px;
  color: #999;
}

.metric-value {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.metric-value.danger {
  color: #ff4d4f;
}

.metric-trend {
  font-size: 12px;
  font-weight: 500;
  padding: 2px 8px;
  border-radius: 4px;
}

.metric-trend.up {
  color: #ff4d4f;
  background: #fff1f0;
}

.metric-trend.down {
  color: #52c41a;
  background: #f6ffed;
}

.metric-progress {
  width: 80px;
  height: 6px;
  background: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: #1890ff;
  border-radius: 3px;
  transition: width 0.3s;
}

.progress-bar.success {
  background: #52c41a;
}

.metric-status {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.metric-status.warning {
  background: #fff7e6;
  color: #fa8c16;
}

/* ========== 实时动态 ========== */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
  transition: all 0.3s;
}

.activity-item:hover {
  background: #f0f5ff;
}

.activity-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 18px;
  flex-shrink: 0;
}

.activity-icon.ai {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.activity-icon.success {
  background: #f6ffed;
}

.activity-icon.user {
  background: #e6f7ff;
}

.activity-icon.warning {
  background: #fff7e6;
}

.activity-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.activity-text {
  font-size: 13px;
  color: #333;
  line-height: 1.5;
}

.activity-text strong {
  color: #1890ff;
  font-weight: 600;
}

.activity-time {
  font-size: 11px;
  color: #999;
}

/* ========== 待办任务 ========== */
.task-count {
  margin-left: auto;
  background: #ff4d4f;
  color: #fff;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.task-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #fafafa;
  border-radius: 8px;
  transition: all 0.3s;
  align-items: flex-start;
}

.task-item:hover {
  background: #f0f5ff;
}

.task-item.completed {
  opacity: 0.6;
}

.task-item.completed .task-title {
  text-decoration: line-through;
}

.task-checkbox {
  margin-top: 2px;
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.task-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.task-title {
  font-size: 13px;
  color: #333;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}

.task-priority {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 500;
}

.task-priority.high {
  background: #ff4d4f;
  color: #fff;
}

.task-priority.medium {
  background: #faad14;
  color: #fff;
}

.task-meta {
  display: flex;
  gap: 12px;
  align-items: center;
  font-size: 11px;
}

.task-tag {
  padding: 2px 8px;
  background: #e6f7ff;
  color: #1890ff;
  border-radius: 4px;
}

.task-deadline {
  color: #999;
}

.task-done {
  color: #52c41a;
}

.view-all-btn {
  width: 100%;
  padding: 10px;
  background: #fafafa;
  border: 1px dashed #d9d9d9;
  border-radius: 8px;
  color: #666;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s;
}

.view-all-btn:hover {
  background: #f0f5ff;
  border-color: #1890ff;
  color: #1890ff;
}

/* ========== AI浮动按钮 ========== */
.ai-float-button {
  position: absolute;
  bottom: 24px;
  right: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  padding: 14px 24px;
  border-radius: 50px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
  transition: all 0.3s ease;
  z-index: 10;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.ai-float-button:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.5);
}

.ai-float-icon {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ai-icon {
  font-size: 24px;
  z-index: 2;
}

.ai-pulse {
  position: absolute;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  animation: pulse 2s ease-out infinite;
}

@keyframes pulse {
  0% {
    transform: scale(0.8);
    opacity: 1;
  }
  100% {
    transform: scale(1.5);
    opacity: 0;
  }
}

.ai-float-text {
  font-size: 15px;
  font-weight: 600;
}

/* ========== AI对话弹窗 ========== */
.ai-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.ai-dialog {
  width: 600px;
  max-width: 90vw;
  height: 700px;
  max-height: 85vh;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.ai-dialog-header {
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 16px 16px 0 0;
}

.ai-dialog-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
}

.ai-dialog-icon {
  font-size: 24px;
}

.ai-status.online {
  font-size: 12px;
  margin-left: 8px;
}

.ai-dialog-close {
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  border-radius: 8px;
  cursor: pointer;
  font-size: 18px;
  transition: all 0.3s;
}

.ai-dialog-close:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(90deg);
}

.ai-dialog-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.ai-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.ai-message {
  display: flex;
  gap: 12px;
  animation: messageIn 0.3s ease;
}

@keyframes messageIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.ai-message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  background: #f0f5ff;
  flex-shrink: 0;
}

.ai-message.user .message-avatar {
  background: #e6f7ff;
}

.message-content {
  max-width: 70%;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ai-message.user .message-content {
  align-items: flex-end;
}

.message-text {
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
  color: #333;
}

.ai-message.user .message-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.message-time {
  font-size: 11px;
  color: #999;
  padding: 0 8px;
}

.ai-input-area {
  padding: 16px 20px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  gap: 12px;
  background: #fafafa;
}

.ai-input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid #e8e8e8;
  border-radius: 24px;
  font-size: 14px;
  transition: all 0.3s;
}

.ai-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.ai-send-btn {
  width: 48px;
  height: 48px;
  border: none;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  border-radius: 50%;
  cursor: pointer;
  font-size: 20px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ai-send-btn:hover {
  transform: scale(1.1) rotate(15deg);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
}

.ai-quick-questions {
  padding: 12px 20px 16px;
  background: #fafafa;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.quick-label {
  font-size: 12px;
  color: #999;
}

.quick-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.quick-tag {
  padding: 6px 14px;
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 16px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
  transition: all 0.3s;
}

.quick-tag:hover {
  background: #667eea;
  color: #fff;
  border-color: #667eea;
  transform: translateY(-2px);
}

/* ========== 滚动条样式 ========== */
.data-panel::-webkit-scrollbar,
.ai-messages::-webkit-scrollbar {
  width: 6px;
}

.data-panel::-webkit-scrollbar-track,
.ai-messages::-webkit-scrollbar-track {
  background: transparent;
}

.data-panel::-webkit-scrollbar-thumb,
.ai-messages::-webkit-scrollbar-thumb {
  background: #d9d9d9;
  border-radius: 3px;
}

.data-panel::-webkit-scrollbar-thumb:hover,
.ai-messages::-webkit-scrollbar-thumb:hover {
  background: #bfbfbf;
}

/* ========== 响应式 ========== */
@media (max-width: 1440px) {
  .data-panel {
    width: 340px;
  }
}

@media (max-width: 1200px) {
  .data-panel {
    width: 300px;
  }
}

.home-page {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

/* ========== 顶部导航栏 ========== */
.home-header {
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
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 28px;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.1); }
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
.page-content {
  flex: 1;
  display: flex;
  gap: 16px;
  padding: 16px;
  overflow: hidden;
  min-height: 0;
}


</style>
