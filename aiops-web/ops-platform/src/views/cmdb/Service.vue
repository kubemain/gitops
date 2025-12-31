<!-- 服务管理 -->
<template>
  <div class="service-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-info">
        <h2><el-icon><SetUp /></el-icon> 服务管理</h2>
        <p>管理应用下的微服务/模块，绑定基础设施资源</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          新建服务
        </el-button>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-select v-model="filters.appId" placeholder="选择应用" clearable style="width: 200px">
        <el-option 
          v-for="item in appList" 
          :key="item.id" 
          :label="item.name" 
          :value="item.id"
        />
      </el-select>
      <el-input 
        v-model="filters.keyword" 
        placeholder="搜索服务名称" 
        style="width: 200px"
        clearable
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 服务列表 -->
    <el-table :data="filteredServices" style="width: 100%" @row-click="viewDetail">
      <el-table-column label="服务" min-width="200">
        <template #default="{ row }">
          <div class="service-info">
            <el-icon class="service-icon"><SetUp /></el-icon>
            <div class="service-text">
              <div class="service-name">{{ row.name }}</div>
              <div class="service-app">{{ row.appName }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          <span class="type-tag" :class="row.type">{{ typeMap[row.type] }}</span>
        </template>
      </el-table-column>
      <el-table-column label="端口" width="80">
        <template #default="{ row }">
          <span class="port">{{ row.port || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="绑定资源" width="200">
        <template #default="{ row }">
          <div class="resource-tags">
            <span class="resource-tag" v-if="row.hostCount">
              <el-icon><Monitor /></el-icon> {{ row.hostCount }} 主机
            </span>
            <span class="resource-tag" v-if="row.dbCount">
              <el-icon><Coin /></el-icon> {{ row.dbCount }} 数据库
            </span>
            <span class="resource-tag" v-if="row.cacheCount">
              <el-icon><Lightning /></el-icon> {{ row.cacheCount }} 缓存
            </span>
            <span v-if="!row.hostCount && !row.dbCount && !row.cacheCount" class="no-resource">
              未绑定
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="负责人" width="100" prop="owner" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <span class="status-dot" :class="row.status"></span>
          {{ row.status === 'running' ? '运行中' : '已停止' }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click.stop="editService(row)">编辑</el-button>
          <el-button type="primary" link size="small" @click.stop="bindResource(row)">绑定</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <div v-if="filteredServices.length === 0" class="empty-state">
      <el-icon class="empty-icon"><SetUp /></el-icon>
      <div class="empty-title">还没有服务</div>
      <div class="empty-desc">创建服务，绑定主机、数据库等基础设施</div>
      <el-button type="primary" @click="showAddDialog">
        <el-icon><Plus /></el-icon>
        新建服务
      </el-button>
    </div>

    <!-- 新建/编辑对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑服务' : '新建服务'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="服务名称" prop="name">
          <el-input v-model="form.name" placeholder="如：user-service、order-api" />
        </el-form-item>
        <el-form-item label="所属应用" prop="appId">
          <el-select v-model="form.appId" placeholder="选择应用" style="width: 100%">
            <el-option 
              v-for="item in appList" 
              :key="item.id" 
              :label="item.name" 
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="服务类型" prop="type">
          <el-select v-model="form.type" placeholder="选择类型" style="width: 100%">
            <el-option label="API服务" value="api" />
            <el-option label="Web服务" value="web" />
            <el-option label="后台任务" value="worker" />
            <el-option label="定时任务" value="cron" />
          </el-select>
        </el-form-item>
        <el-form-item label="端口">
          <el-input-number v-model="form.port" :min="1" :max="65535" placeholder="服务端口" />
        </el-form-item>
        <el-form-item label="负责人" prop="owner">
          <el-input v-model="form.owner" placeholder="服务负责人" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="服务描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 绑定资源对话框 -->
    <el-dialog v-model="bindDialogVisible" title="绑定资源" width="600px">
      <div class="bind-section">
        <h4><el-icon><Monitor /></el-icon> 绑定主机</h4>
        <el-select v-model="bindForm.hostIds" multiple placeholder="选择主机" style="width: 100%">
          <el-option 
            v-for="host in hostList" 
            :key="host.id" 
            :label="`${host.hostname} (${host.ip})`" 
            :value="host.id"
          />
        </el-select>
      </div>
      <div class="bind-section">
        <h4><el-icon><Coin /></el-icon> 绑定数据库</h4>
        <el-select v-model="bindForm.dbIds" multiple placeholder="选择数据库" style="width: 100%">
          <el-option label="MySQL主库 (192.168.1.20)" :value="1" />
          <el-option label="MySQL从库 (192.168.1.21)" :value="2" />
          <el-option label="PostgreSQL (192.168.1.25)" :value="3" />
        </el-select>
      </div>
      <div class="bind-section">
        <h4><el-icon><Lightning /></el-icon> 绑定缓存</h4>
        <el-select v-model="bindForm.cacheIds" multiple placeholder="选择缓存" style="width: 100%">
          <el-option label="Redis集群 (192.168.1.30)" :value="1" />
          <el-option label="Memcached (192.168.1.35)" :value="2" />
        </el-select>
      </div>
      <template #footer>
        <el-button @click="bindDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleBind">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Search, SetUp, Monitor, Coin, Lightning } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const typeMap = {
  api: 'API',
  web: 'Web',
  worker: '后台',
  cron: '定时'
}

// 应用列表
const appList = ref([
  { id: 1, name: '订单系统' },
  { id: 2, name: '用户中心' },
  { id: 3, name: '商品系统' },
  { id: 4, name: '账户系统' },
  { id: 5, name: '数据采集' }
])

// 主机列表
const hostList = ref([
  { id: 1, hostname: 'web-server-01', ip: '192.168.1.10' },
  { id: 2, hostname: 'web-server-02', ip: '192.168.1.11' },
  { id: 3, hostname: 'api-server-01', ip: '192.168.1.15' },
  { id: 4, hostname: 'api-server-02', ip: '192.168.1.16' }
])

// 服务列表
const serviceList = ref([
  {
    id: 1,
    name: 'order-api',
    appId: 1,
    appName: '订单系统',
    type: 'api',
    port: 8080,
    owner: '张三',
    status: 'running',
    hostCount: 2,
    dbCount: 1,
    cacheCount: 1,
    description: '订单核心API服务'
  },
  {
    id: 2,
    name: 'order-worker',
    appId: 1,
    appName: '订单系统',
    type: 'worker',
    port: null,
    owner: '张三',
    status: 'running',
    hostCount: 1,
    dbCount: 1,
    cacheCount: 0,
    description: '订单异步处理服务'
  },
  {
    id: 3,
    name: 'user-api',
    appId: 2,
    appName: '用户中心',
    type: 'api',
    port: 8081,
    owner: '李四',
    status: 'running',
    hostCount: 2,
    dbCount: 1,
    cacheCount: 1,
    description: '用户服务API'
  },
  {
    id: 4,
    name: 'user-web',
    appId: 2,
    appName: '用户中心',
    type: 'web',
    port: 80,
    owner: '李四',
    status: 'running',
    hostCount: 2,
    dbCount: 0,
    cacheCount: 0,
    description: '用户前端服务'
  },
  {
    id: 5,
    name: 'product-api',
    appId: 3,
    appName: '商品系统',
    type: 'api',
    port: 8082,
    owner: '王五',
    status: 'running',
    hostCount: 2,
    dbCount: 1,
    cacheCount: 1,
    description: '商品服务API'
  }
])

// 筛选
const filters = reactive({
  appId: route.query.app_id ? Number(route.query.app_id) : null,
  keyword: ''
})

const filteredServices = computed(() => {
  return serviceList.value.filter(s => {
    if (filters.appId && s.appId !== filters.appId) return false
    if (filters.keyword && !s.name.includes(filters.keyword)) return false
    return true
  })
})

// 新建/编辑对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const form = reactive({
  id: null,
  name: '',
  appId: null,
  type: 'api',
  port: null,
  owner: '',
  description: ''
})

const rules = {
  name: [{ required: true, message: '请输入服务名称', trigger: 'blur' }],
  appId: [{ required: true, message: '请选择应用', trigger: 'change' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  owner: [{ required: true, message: '请输入负责人', trigger: 'blur' }]
}

// 绑定资源对话框
const bindDialogVisible = ref(false)
const currentService = ref(null)
const bindForm = reactive({
  hostIds: [],
  dbIds: [],
  cacheIds: []
})

const showAddDialog = () => {
  isEdit.value = false
  Object.assign(form, { 
    id: null, 
    name: '', 
    appId: filters.appId, 
    type: 'api',
    port: null,
    owner: '', 
    description: '' 
  })
  dialogVisible.value = true
}

const editService = (row) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  await formRef.value.validate()
  const app = appList.value.find(a => a.id === form.appId)
  if (isEdit.value) {
    const index = serviceList.value.findIndex(s => s.id === form.id)
    if (index > -1) {
      serviceList.value[index] = { 
        ...serviceList.value[index], 
        ...form,
        appName: app?.name 
      }
    }
    ElMessage.success('更新成功')
  } else {
    serviceList.value.push({
      ...form,
      id: Date.now(),
      appName: app?.name,
      status: 'running',
      hostCount: 0,
      dbCount: 0,
      cacheCount: 0
    })
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
}

const bindResource = (row) => {
  currentService.value = row
  Object.assign(bindForm, { hostIds: [], dbIds: [], cacheIds: [] })
  bindDialogVisible.value = true
}

const handleBind = () => {
  if (currentService.value) {
    const index = serviceList.value.findIndex(s => s.id === currentService.value.id)
    if (index > -1) {
      serviceList.value[index].hostCount = bindForm.hostIds.length
      serviceList.value[index].dbCount = bindForm.dbIds.length
      serviceList.value[index].cacheCount = bindForm.cacheIds.length
    }
  }
  ElMessage.success('绑定成功')
  bindDialogVisible.value = false
}

const viewDetail = (row) => {
  router.push({ path: '/cmdb/hosts', query: { service_id: row.id } })
}
</script>

<style scoped>
.service-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-info h2 {
  margin: 0 0 4px;
  font-size: 20px;
  color: #333;
}

.header-info p {
  margin: 0;
  font-size: 13px;
  color: #999;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.service-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.service-icon {
  font-size: 20px;
  color: #667eea;
}

.service-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.service-app {
  font-size: 12px;
  color: #999;
}

.type-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
}

.type-tag.api {
  background: #e6f7ff;
  color: #1890ff;
}

.type-tag.web {
  background: #f6ffed;
  color: #52c41a;
}

.type-tag.worker {
  background: #fff7e6;
  color: #faad14;
}

.type-tag.cron {
  background: #f9f0ff;
  color: #722ed1;
}

.port {
  font-family: monospace;
  color: #666;
}

.resource-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.resource-tag {
  font-size: 11px;
  padding: 2px 6px;
  background: #f5f5f5;
  border-radius: 4px;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4px;
}

.resource-tag .el-icon {
  font-size: 12px;
  color: #667eea;
}

.no-resource {
  font-size: 12px;
  color: #999;
}

.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 6px;
}

.status-dot.running {
  background: #52c41a;
}

.status-dot.stopped {
  background: #ff4d4f;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: #fff;
  border-radius: 12px;
  margin-top: 16px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  color: #d9d9d9;
}

.empty-title {
  font-size: 18px;
  color: #333;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 14px;
  color: #999;
  margin-bottom: 24px;
}

/* 绑定资源 */
.bind-section {
  margin-bottom: 20px;
}

.bind-section h4 {
  margin: 0 0 10px;
  font-size: 14px;
  color: #333;
  display: flex;
  align-items: center;
  gap: 6px;
}

.bind-section h4 .el-icon {
  color: #667eea;
}
</style>
