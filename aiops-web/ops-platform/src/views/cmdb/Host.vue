<template>
  <div class="host-page">
    <el-card class="page-card">
      <!-- 页面标题 -->
      <div class="page-header">
        <div class="header-title">
          <el-icon class="title-icon"><Monitor /></el-icon>
          <h2>主机管理</h2>
          <el-tag type="info" size="small">共 {{ pagination.total }} 台主机</el-tag>
        </div>
        <div class="header-actions">
          <el-radio-group v-model="viewMode" size="small">
            <el-radio-button label="table">
              <el-icon><List /></el-icon>
            </el-radio-button>
            <el-radio-button label="card">
              <el-icon><Grid /></el-icon>
            </el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <!-- 高级搜索栏 -->
      <div class="search-bar">
        <div class="search-row">
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索主机名/IP"
            style="width: 220px"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          
          <el-select
            v-model="searchForm.group_id"
            placeholder="主机组"
            style="width: 150px"
            clearable
          >
            <el-option
              v-for="group in hostGroups"
              :key="group.id"
              :label="group.name"
              :value="group.id"
            />
          </el-select>

          <el-select
            v-model="searchForm.status"
            placeholder="状态"
            style="width: 120px"
            clearable
          >
            <el-option label="在线" value="online" />
            <el-option label="离线" value="offline" />
            <el-option label="维护中" value="maintenance" />
          </el-select>

          <el-select
            v-model="searchForm.environment"
            placeholder="环境"
            style="width: 120px"
            clearable
          >
            <el-option label="生产环境" value="production" />
            <el-option label="预发环境" value="staging" />
            <el-option label="开发环境" value="development" />
          </el-select>

          <el-select
            v-model="searchForm.region"
            placeholder="区域"
            style="width: 120px"
            clearable
          >
            <el-option v-for="r in regions" :key="r" :label="r" :value="r" />
          </el-select>

          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">重置</el-button>
        </div>

        <div class="action-row">
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>
            新增主机
          </el-button>
          <el-button @click="$router.push('/cmdb/import')">
            <el-icon><Upload /></el-icon>
            批量导入
          </el-button>
          <el-button :disabled="selectedRows.length === 0" @click="handleBatchDelete">
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </div>

      <!-- 表格视图 -->
      <template v-if="viewMode === 'table'">
        <el-table
          :data="tableData"
          style="width: 100%"
          v-loading="loading"
          @selection-change="handleSelectionChange"
          :header-cell-style="{ background: '#fafafa', color: '#333' }"
        >
          <el-table-column type="selection" width="50" />
          <el-table-column prop="hostname" label="主机名" width="160" fixed>
            <template #default="{ row }">
              <div class="host-name-cell">
                <span class="status-dot" :class="row.status"></span>
                <span class="hostname">{{ row.hostname }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="ip" label="IP地址" width="140" />
          <el-table-column prop="os" label="操作系统" width="120" />
          <el-table-column label="配置" width="180">
            <template #default="{ row }">
              <div class="config-cell">
                <span><el-icon><Cpu /></el-icon> {{ row.cpu }}核</span>
                <span><el-icon><Coin /></el-icon> {{ formatMemory(row.memory) }}</span>
                <span><el-icon><Box /></el-icon> {{ row.disk }}GB</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="主机组" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.host_group" size="small" type="success">
                {{ row.host_group.name }}
              </el-tag>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column label="环境" width="100">
            <template #default="{ row }">
              <el-tag :type="getEnvType(row.environment)" size="small">
                {{ getEnvName(row.environment) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="region" label="区域" width="100" />
          <el-table-column label="标签" min-width="150">
            <template #default="{ row }">
              <div class="tags-cell">
                <el-tag
                  v-for="tag in (row.tags || []).slice(0, 3)"
                  :key="tag"
                  size="small"
                  type="info"
                  class="tag-item"
                >
                  {{ tag }}
                </el-tag>
                <el-tag v-if="(row.tags || []).length > 3" size="small" type="info">
                  +{{ row.tags.length - 3 }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusName(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button link type="warning" size="small" @click="handleToggleStatus(row)">
                {{ row.status === 'online' ? '下线' : '上线' }}
              </el-button>
              <el-button link type="danger" size="small" @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>

      <!-- 卡片视图 -->
      <template v-else>
        <div class="card-grid" v-loading="loading">
          <div v-for="host in tableData" :key="host.id" class="host-card">
            <div class="card-header">
              <div class="host-status" :class="host.status"></div>
              <div class="host-info">
                <h4>{{ host.hostname }}</h4>
                <span class="host-ip">{{ host.ip }}</span>
              </div>
              <el-dropdown trigger="click">
                <el-button link>
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="handleEdit(host)">编辑</el-dropdown-item>
                    <el-dropdown-item @click="handleToggleStatus(host)">
                      {{ host.status === 'online' ? '下线' : '上线' }}
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="handleDelete(host)">
                      <span style="color: #ff4d4f">删除</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>

            <div class="card-body">
              <div class="config-info">
                <div class="config-item">
                  <span class="label">操作系统</span>
                  <span class="value">{{ host.os || '-' }}</span>
                </div>
                <div class="config-item">
                  <span class="label">CPU</span>
                  <span class="value">{{ host.cpu }}核</span>
                </div>
                <div class="config-item">
                  <span class="label">内存</span>
                  <span class="value">{{ formatMemory(host.memory) }}</span>
                </div>
                <div class="config-item">
                  <span class="label">磁盘</span>
                  <span class="value">{{ host.disk }}GB</span>
                </div>
              </div>

              <div class="card-tags">
                <el-tag v-if="host.host_group" size="small" type="success">
                  {{ host.host_group.name }}
                </el-tag>
                <el-tag :type="getEnvType(host.environment)" size="small">
                  {{ getEnvName(host.environment) }}
                </el-tag>
              </div>
            </div>

            <div class="card-footer">
              <el-tag :type="getStatusType(host.status)" size="small">
                {{ getStatusName(host.status) }}
              </el-tag>
              <span class="update-time">{{ formatDate(host.updated_at) }}</span>
            </div>
          </div>
          <el-empty v-if="!loading && tableData.length === 0" description="暂无主机数据" />
        </div>
      </template>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="pagination.total"
          :page-sizes="[12, 24, 48, 96]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="750px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主机名" prop="hostname">
              <el-input v-model="formData.hostname" placeholder="请输入主机名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IP地址" prop="ip">
              <el-input v-model="formData.ip" placeholder="请输入IP地址" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="公网IP" prop="public_ip">
              <el-input v-model="formData.public_ip" placeholder="请输入公网IP（可选）" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作系统" prop="os">
              <el-input v-model="formData.os" placeholder="如：CentOS 7.9" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="CPU核心" prop="cpu">
              <el-input-number v-model="formData.cpu" :min="1" :max="256" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="内存(MB)" prop="memory">
              <el-input-number v-model="formData.memory" :min="512" :max="1048576" :step="1024" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="磁盘(GB)" prop="disk">
              <el-input-number v-model="formData.disk" :min="10" :max="102400" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主机组" prop="group_id">
              <el-select v-model="formData.group_id" placeholder="请选择主机组" style="width: 100%" clearable>
                <el-option v-for="group in hostGroups" :key="group.id" :label="group.name" :value="group.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="环境" prop="environment">
              <el-select v-model="formData.environment" placeholder="请选择环境" style="width: 100%">
                <el-option label="生产环境" value="production" />
                <el-option label="预发环境" value="staging" />
                <el-option label="开发环境" value="development" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="区域" prop="region">
              <el-input v-model="formData.region" placeholder="如：cn-north-1" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IDC机房" prop="idc">
              <el-input v-model="formData.idc" placeholder="如：IDC-BJ-01" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="标签" prop="tags">
          <el-select
            v-model="formData.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="输入标签后回车添加"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" rows="3" placeholder="请输入备注信息" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio label="online">在线</el-radio>
            <el-radio label="offline">离线</el-radio>
            <el-radio label="maintenance">维护中</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'
import {
  Search, Plus, Monitor, List, Grid, Upload, Download, Delete,
  MoreFilled, Cpu, Coin, Box
} from '@element-plus/icons-vue'
import {
  getHostList, createHost, updateHost, deleteHost, updateHostStatus, batchDeleteHost
} from '@/api/host'
import { getAllHostGroups } from '@/api/host-group'

const route = useRoute()

// 视图模式
const viewMode = ref('table')

// 表格数据
const tableData = ref([])
const loading = ref(false)
const selectedRows = ref([])

// 分页
const pagination = reactive({
  page: 1,
  page_size: 12,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  keyword: '',
  group_id: '',
  status: '',
  environment: '',
  region: ''
})

// 主机组和区域列表
const hostGroups = ref([])
const regions = computed(() => {
  const set = new Set(tableData.value.map(h => h.region).filter(Boolean))
  return Array.from(set)
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('新增主机')
const isEdit = ref(false)
const submitLoading = ref(false)

const formRef = ref(null)
const formData = reactive({
  id: null,
  hostname: '',
  ip: '',
  public_ip: '',
  os: '',
  cpu: 4,
  memory: 8192,
  disk: 100,
  group_id: '',
  environment: 'production',
  region: '',
  idc: '',
  tags: [],
  status: 'online',
  remark: ''
})

// 表单验证规则
const formRules = {
  hostname: [{ required: true, message: '请输入主机名', trigger: 'blur' }],
  ip: [
    { required: true, message: '请输入IP地址', trigger: 'blur' },
    { pattern: /^(\d{1,3}\.){3}\d{1,3}$/, message: '请输入正确的IP地址', trigger: 'blur' }
  ],
  os: [{ required: true, message: '请输入操作系统', trigger: 'blur' }],
  cpu: [{ required: true, message: '请输入CPU核心数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  disk: [{ required: true, message: '请输入磁盘大小', trigger: 'blur' }]
}

// 格式化函数
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}

const formatMemory = (mb) => {
  if (!mb) return '-'
  if (mb >= 1024) return `${(mb / 1024).toFixed(0)}GB`
  return `${mb}MB`
}

const getStatusType = (status) => {
  const types = { online: 'success', offline: 'danger', maintenance: 'warning' }
  return types[status] || 'info'
}

const getStatusName = (status) => {
  const names = { online: '在线', offline: '离线', maintenance: '维护中' }
  return names[status] || status
}

const getEnvType = (env) => {
  const types = { production: 'success', staging: 'warning', development: '' }
  return types[env] || 'info'
}

const getEnvName = (env) => {
  const names = { production: '生产', staging: '预发', development: '开发' }
  return names[env] || env || '-'
}

// 获取主机列表
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      hostname: searchForm.keyword || undefined,
      ip: searchForm.keyword || undefined,
      group_id: searchForm.group_id || undefined,
      status: searchForm.status || undefined,
      environment: searchForm.environment || undefined,
      region: searchForm.region || undefined
    }
    const data = await getHostList(params)
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error('获取主机列表失败')
  } finally {
    loading.value = false
  }
}

// 获取主机组列表
const fetchHostGroups = async () => {
  try {
    const data = await getAllHostGroups()
    hostGroups.value = data || []
  } catch (error) {
    console.error('获取主机组列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const resetSearch = () => {
  Object.assign(searchForm, { keyword: '', group_id: '', status: '', environment: '', region: '' })
  handleSearch()
}

// 选择变化
const handleSelectionChange = (rows) => {
  selectedRows.value = rows
}

// 新增
const handleCreate = () => {
  isEdit.value = false
  dialogTitle.value = '新增主机'
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑主机'
  Object.assign(formData, {
    id: row.id,
    hostname: row.hostname,
    ip: row.ip,
    public_ip: row.public_ip || '',
    os: row.os,
    cpu: row.cpu,
    memory: row.memory,
    disk: row.disk,
    group_id: row.host_group?.id || '',
    environment: row.environment || 'production',
    region: row.region || '',
    idc: row.idc || '',
    tags: row.tags || [],
    status: row.status,
    remark: row.remark || ''
  })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除主机「${row.hostname}」吗？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteHost(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') console.error('删除失败:', error)
  }
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) return
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 台主机吗？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await batchDeleteHost({ ids: selectedRows.value.map(r => r.id) })
    ElMessage.success('批量删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') console.error('批量删除失败:', error)
  }
}

// 切换状态
const handleToggleStatus = async (row) => {
  try {
    const newStatus = row.status === 'online' ? 'offline' : 'online'
    await updateHostStatus(row.id, { status: newStatus })
    ElMessage.success(`主机已${newStatus === 'online' ? '上线' : '下线'}`)
    fetchData()
  } catch (error) {
    console.error('状态更新失败:', error)
    ElMessage.error('状态更新失败')
  }
}

// 导出
const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const submitData = { ...formData }
      delete submitData.id
      if (isEdit.value) {
        await updateHost(formData.id, submitData)
        ElMessage.success('更新成功')
      } else {
        await createHost(submitData)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchData()
    } catch (error) {
      console.error('提交失败:', error)
      ElMessage.error('提交失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// 关闭对话框
const handleDialogClose = () => {
  formRef.value?.resetFields()
  Object.assign(formData, {
    id: null, hostname: '', ip: '', public_ip: '', os: '',
    cpu: 4, memory: 8192, disk: 100, group_id: '',
    environment: 'production', region: '', idc: '',
    tags: [], status: 'online', remark: ''
  })
}

// 监听路由参数
watch(() => route.query, (query) => {
  if (query.group_id) searchForm.group_id = parseInt(query.group_id)
  if (query.environment) searchForm.environment = query.environment
  if (query.region) searchForm.region = query.region
  if (query.search) searchForm.keyword = query.search
  handleSearch()
}, { immediate: true })

onMounted(() => {
  fetchData()
  fetchHostGroups()
})
</script>


<style scoped>
.host-page {
  height: 100%;
}

.page-card {
  height: 100%;
  border-radius: 12px;
}

.page-card :deep(.el-card__body) {
  display: flex;
  flex-direction: column;
  height: calc(100% - 60px);
}

/* ========== 页面标题 ========== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 24px;
  color: #1890ff;
}

.header-title h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

/* ========== 搜索栏 ========== */
.search-bar {
  margin-bottom: 16px;
}

.search-row {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.action-row {
  display: flex;
  gap: 10px;
}

/* ========== 表格样式 ========== */
.host-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-dot.online { background: #52c41a; }
.status-dot.offline { background: #ff4d4f; }
.status-dot.maintenance { background: #faad14; }

.hostname {
  font-weight: 500;
  color: #333;
}

.config-cell {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #666;
}

.config-cell span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.tags-cell {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tag-item {
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.text-muted {
  color: #999;
}

/* ========== 卡片视图 ========== */
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  flex: 1;
  overflow-y: auto;
  padding: 4px;
}

.host-card {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  padding: 16px;
  transition: all 0.3s;
}

.host-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #1890ff;
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}

.host-status {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-top: 6px;
  flex-shrink: 0;
}

.host-status.online { background: #52c41a; }
.host-status.offline { background: #ff4d4f; }
.host-status.maintenance { background: #faad14; }

.host-info {
  flex: 1;
}

.host-info h4 {
  margin: 0 0 4px;
  font-size: 15px;
  font-weight: 600;
  color: #333;
}

.host-ip {
  font-size: 13px;
  color: #999;
}

.card-body {
  margin-bottom: 12px;
}

.config-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-bottom: 12px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.config-item .label {
  color: #999;
}

.config-item .value {
  color: #333;
  font-weight: 500;
}

.card-tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.update-time {
  font-size: 12px;
  color: #999;
}

/* ========== 分页 ========== */
.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
  flex-shrink: 0;
}

/* ========== 响应式 ========== */
@media (max-width: 1200px) {
  .search-row {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .card-grid {
    grid-template-columns: 1fr;
  }
}
</style>
