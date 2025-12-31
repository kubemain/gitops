<!-- 变更历史页面 -->
<template>
  <div class="change-history">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="title-icon"><Clock /></el-icon>
            <h2>变更历史</h2>
            <el-tag type="info" size="small">追踪所有资源变更记录</el-tag>
          </div>
          <div class="header-right">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              size="small"
              style="width: 260px"
              @change="handleSearch"
            />
            <el-select v-model="changeType" placeholder="变更类型" size="small" clearable style="width: 120px" @change="handleSearch">
              <el-option label="全部" value="" />
              <el-option label="新增" value="create" />
              <el-option label="更新" value="update" />
              <el-option label="删除" value="delete" />
              <el-option label="状态变更" value="status" />
            </el-select>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索主机名"
              size="small"
              style="width: 180px"
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
        </div>
      </template>

      <div class="history-container" v-loading="loading">
        <!-- 变更统计 -->
        <div class="change-stats">
          <div class="stat-item create">
            <el-icon><Plus /></el-icon>
            <span class="stat-value">{{ stats.create || 0 }}</span>
            <span class="stat-label">新增</span>
          </div>
          <div class="stat-item update">
            <el-icon><Edit /></el-icon>
            <span class="stat-value">{{ stats.update || 0 }}</span>
            <span class="stat-label">更新</span>
          </div>
          <div class="stat-item delete">
            <el-icon><Delete /></el-icon>
            <span class="stat-value">{{ stats.delete || 0 }}</span>
            <span class="stat-label">删除</span>
          </div>
          <div class="stat-item status">
            <el-icon><Switch /></el-icon>
            <span class="stat-value">{{ stats.status || 0 }}</span>
            <span class="stat-label">状态变更</span>
          </div>
        </div>

        <!-- 时间线视图 -->
        <div class="timeline-container">
          <el-timeline>
            <el-timeline-item
              v-for="change in changeList"
              :key="change.id"
              :timestamp="formatTime(change.created_at)"
              placement="top"
              :type="getChangeColor(change.change_type)"
              :hollow="true"
              size="large"
            >
              <el-card class="change-card" shadow="hover">
                <div class="change-header">
                  <div class="change-info">
                    <el-tag :type="getChangeColor(change.change_type)" size="small">
                      {{ getChangeLabel(change.change_type) }}
                    </el-tag>
                    <span class="change-host" @click="goToHost(change.host_id)">
                      {{ change.hostname || `主机#${change.host_id}` }}
                    </span>
                  </div>
                  <span class="change-operator">
                    <el-icon><User /></el-icon>
                    {{ change.operator || '系统' }}
                  </span>
                </div>

                <div class="change-content">
                  <div v-if="change.old_value || change.new_value" class="change-diff">
                    <div v-if="change.old_value" class="diff-item old">
                      <span class="diff-label">变更前:</span>
                      <code>{{ change.old_value }}</code>
                    </div>
                    <div v-if="change.new_value" class="diff-item new">
                      <span class="diff-label">变更后:</span>
                      <code>{{ change.new_value }}</code>
                    </div>
                  </div>
                  <div v-if="change.remark" class="change-remark">
                    <el-icon><ChatDotRound /></el-icon>
                    {{ change.remark }}
                  </div>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>

          <el-empty v-if="!loading && changeList.length === 0" description="暂无变更记录" />
        </div>

        <!-- 分页 -->
        <div class="pagination" v-if="pagination.total > 0">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.page_size"
            :total="pagination.total"
            :page-sizes="[20, 50, 100]"
            layout="total, sizes, prev, pager, next"
            @size-change="fetchChanges"
            @current-change="fetchChanges"
          />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Clock, Search, Plus, Edit, Delete, Switch, User, ChatDotRound } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const dateRange = ref([])
const changeType = ref('')
const searchKeyword = ref('')

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const stats = ref({
  create: 0,
  update: 0,
  delete: 0,
  status: 0
})

// 模拟变更数据
const changeList = ref([])

const fetchChanges = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    changeList.value = [
      { id: 1, host_id: 1, hostname: 'web-server-01', change_type: 'update', old_value: 'cpu: 4', new_value: 'cpu: 8', operator: 'admin', remark: '升级CPU配置', created_at: new Date().toISOString() },
      { id: 2, host_id: 2, hostname: 'db-master-01', change_type: 'status', old_value: 'offline', new_value: 'online', operator: 'system', remark: '主机恢复在线', created_at: new Date(Date.now() - 3600000).toISOString() },
      { id: 3, host_id: 3, hostname: 'redis-01', change_type: 'create', new_value: '新增主机', operator: 'admin', remark: '新增缓存服务器', created_at: new Date(Date.now() - 7200000).toISOString() },
      { id: 4, host_id: 4, hostname: 'mq-01', change_type: 'update', old_value: 'memory: 8192', new_value: 'memory: 16384', operator: 'admin', remark: '扩容内存', created_at: new Date(Date.now() - 86400000).toISOString() },
      { id: 5, host_id: 5, hostname: 'monitor-01', change_type: 'delete', old_value: '主机信息', operator: 'admin', remark: '下线旧监控服务器', created_at: new Date(Date.now() - 172800000).toISOString() }
    ]

    stats.value = { create: 12, update: 45, delete: 3, status: 28 }
    pagination.total = 88
  } catch (error) {
    console.error('获取变更记录失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchChanges()
}

const formatTime = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getChangeColor = (type) => {
  const colors = { create: 'success', update: 'primary', delete: 'danger', status: 'warning' }
  return colors[type] || 'info'
}

const getChangeLabel = (type) => {
  const labels = { create: '新增', update: '更新', delete: '删除', status: '状态变更' }
  return labels[type] || type
}

const goToHost = (id) => {
  router.push(`/cmdb/hosts?id=${id}`)
}

onMounted(() => {
  fetchChanges()
})
</script>

<style scoped>
.change-history {
  height: 100%;
}

.page-card {
  height: 100%;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
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
  align-items: center;
}

.history-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 变更统计 */
.change-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 8px;
  background: #f5f7fa;
}

.stat-item.create { background: #f6ffed; color: #52c41a; }
.stat-item.update { background: #e6f7ff; color: #1890ff; }
.stat-item.delete { background: #fff2f0; color: #ff4d4f; }
.stat-item.status { background: #fffbe6; color: #faad14; }

.stat-item .el-icon {
  font-size: 18px;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
}

.stat-label {
  font-size: 13px;
  opacity: 0.85;
}

/* 时间线 */
.timeline-container {
  max-height: calc(100vh - 380px);
  overflow-y: auto;
  padding-right: 8px;
}

.change-card {
  margin-left: 12px;
  border-radius: 8px;
}

.change-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.change-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.change-host {
  font-size: 15px;
  font-weight: 500;
  color: #1890ff;
  cursor: pointer;
}

.change-host:hover {
  text-decoration: underline;
}

.change-operator {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #999;
}

.change-content {
  font-size: 14px;
}

.change-diff {
  background: #fafafa;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 8px;
}

.diff-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.diff-item:last-child {
  margin-bottom: 0;
}

.diff-label {
  font-size: 12px;
  color: #999;
  width: 60px;
}

.diff-item.old code {
  background: #fff2f0;
  color: #ff4d4f;
  padding: 2px 8px;
  border-radius: 4px;
}

.diff-item.new code {
  background: #f6ffed;
  color: #52c41a;
  padding: 2px 8px;
  border-radius: 4px;
}

.change-remark {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 13px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}
</style>
