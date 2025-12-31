<!-- 应用管理 -->
<template>
  <div class="application-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-info">
        <h2><el-icon><Grid /></el-icon> 应用管理</h2>
        <p>管理业务线下的应用系统，关联服务与基础设施</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          新建应用
        </el-button>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-select v-model="filters.businessId" placeholder="选择业务线" clearable style="width: 200px">
        <el-option 
          v-for="item in businessList" 
          :key="item.id" 
          :label="item.name" 
          :value="item.id"
        />
      </el-select>
      <el-input 
        v-model="filters.keyword" 
        placeholder="搜索应用名称" 
        style="width: 200px"
        clearable
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 应用列表 -->
    <div class="app-list">
      <div 
        v-for="app in filteredApps" 
        :key="app.id" 
        class="app-card"
        @click="viewDetail(app)"
      >
        <div class="card-header">
          <div class="card-icon">
            <el-icon><Grid /></el-icon>
          </div>
          <div class="card-title">
            <h3>{{ app.name }}</h3>
            <div class="card-tags">
              <span class="tag business">{{ app.businessName }}</span>
              <span class="tag env" :class="app.environment">{{ envMap[app.environment] }}</span>
            </div>
          </div>
          <el-dropdown @command="(cmd) => handleCommand(cmd, app)" @click.stop>
            <el-button type="primary" link>
              <el-icon><MoreFilled /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit">编辑</el-dropdown-item>
                <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        
        <div class="card-desc">{{ app.description || '暂无描述' }}</div>
        
        <!-- 关联资源 -->
        <div class="card-resources">
          <div class="resource-item">
            <el-icon class="resource-icon"><SetUp /></el-icon>
            <span class="resource-count">{{ app.serviceCount }}</span>
            <span class="resource-label">服务</span>
          </div>
          <div class="resource-item">
            <el-icon class="resource-icon"><Monitor /></el-icon>
            <span class="resource-count">{{ app.hostCount }}</span>
            <span class="resource-label">主机</span>
          </div>
          <div class="resource-item">
            <el-icon class="resource-icon"><Coin /></el-icon>
            <span class="resource-count">{{ app.dbCount }}</span>
            <span class="resource-label">数据库</span>
          </div>
        </div>

        <div class="card-footer">
          <div class="card-owner">
            <el-icon class="owner-avatar"><User /></el-icon>
            <span class="owner-name">{{ app.owner }}</span>
          </div>
          <div class="card-level" :class="app.level">
            {{ levelMap[app.level] }}
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="filteredApps.length === 0" class="empty-state">
        <el-icon class="empty-icon"><Grid /></el-icon>
        <div class="empty-title">还没有应用</div>
        <div class="empty-desc">创建应用，关联到业务线</div>
        <el-button type="primary" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          新建应用
        </el-button>
      </div>
    </div>

    <!-- 新建/编辑对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑应用' : '新建应用'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="应用名称" prop="name">
          <el-input v-model="form.name" placeholder="如：订单系统、用户中心" />
        </el-form-item>
        <el-form-item label="所属业务" prop="businessId">
          <el-select v-model="form.businessId" placeholder="选择业务线" style="width: 100%">
            <el-option 
              v-for="item in businessList" 
              :key="item.id" 
              :label="item.name" 
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="环境" prop="environment">
          <el-select v-model="form.environment" placeholder="选择环境" style="width: 100%">
            <el-option label="生产环境" value="production" />
            <el-option label="预发环境" value="staging" />
            <el-option label="开发环境" value="development" />
          </el-select>
        </el-form-item>
        <el-form-item label="重要级别" prop="level">
          <el-select v-model="form.level" placeholder="选择级别" style="width: 100%">
            <el-option label="核心应用" value="core" />
            <el-option label="重要应用" value="important" />
            <el-option label="一般应用" value="normal" />
          </el-select>
        </el-form-item>
        <el-form-item label="负责人" prop="owner">
          <el-input v-model="form.owner" placeholder="应用负责人" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="应用描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, MoreFilled, Search, Grid, SetUp, Monitor, Coin, User } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const envMap = {
  production: '生产',
  staging: '预发',
  development: '开发'
}

const levelMap = {
  core: '核心',
  important: '重要',
  normal: '一般'
}

// 业务线列表
const businessList = ref([
  { id: 1, name: '电商平台' },
  { id: 2, name: '金融系统' },
  { id: 3, name: '数据平台' }
])

// 应用列表
const appList = ref([
  {
    id: 1,
    name: '订单系统',
    businessId: 1,
    businessName: '电商平台',
    description: '处理订单创建、支付、发货等核心流程',
    environment: 'production',
    level: 'core',
    owner: '张三',
    serviceCount: 5,
    hostCount: 8,
    dbCount: 2
  },
  {
    id: 2,
    name: '用户中心',
    businessId: 1,
    businessName: '电商平台',
    description: '用户注册、登录、个人信息管理',
    environment: 'production',
    level: 'core',
    owner: '李四',
    serviceCount: 3,
    hostCount: 4,
    dbCount: 1
  },
  {
    id: 3,
    name: '商品系统',
    businessId: 1,
    businessName: '电商平台',
    description: '商品信息管理、库存管理',
    environment: 'production',
    level: 'important',
    owner: '王五',
    serviceCount: 4,
    hostCount: 6,
    dbCount: 1
  },
  {
    id: 4,
    name: '账户系统',
    businessId: 2,
    businessName: '金融系统',
    description: '用户账户、余额、交易记录',
    environment: 'production',
    level: 'core',
    owner: '赵六',
    serviceCount: 4,
    hostCount: 6,
    dbCount: 2
  },
  {
    id: 5,
    name: '数据采集',
    businessId: 3,
    businessName: '数据平台',
    description: '多源数据采集、清洗、入库',
    environment: 'production',
    level: 'important',
    owner: '钱七',
    serviceCount: 3,
    hostCount: 5,
    dbCount: 1
  }
])

// 筛选
const filters = reactive({
  businessId: route.query.business_id ? Number(route.query.business_id) : null,
  keyword: ''
})

const filteredApps = computed(() => {
  return appList.value.filter(app => {
    if (filters.businessId && app.businessId !== filters.businessId) return false
    if (filters.keyword && !app.name.includes(filters.keyword)) return false
    return true
  })
})

// 对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const form = reactive({
  id: null,
  name: '',
  businessId: null,
  environment: 'production',
  level: 'normal',
  owner: '',
  description: ''
})

const rules = {
  name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  businessId: [{ required: true, message: '请选择业务线', trigger: 'change' }],
  environment: [{ required: true, message: '请选择环境', trigger: 'change' }],
  owner: [{ required: true, message: '请输入负责人', trigger: 'blur' }]
}

const showAddDialog = () => {
  isEdit.value = false
  Object.assign(form, { 
    id: null, 
    name: '', 
    businessId: filters.businessId, 
    environment: 'production', 
    level: 'normal',
    owner: '', 
    description: '' 
  })
  dialogVisible.value = true
}

const handleCommand = (command, item) => {
  if (command === 'edit') {
    isEdit.value = true
    Object.assign(form, item)
    dialogVisible.value = true
  } else if (command === 'delete') {
    ElMessageBox.confirm(`确定删除应用「${item.name}」吗？`, '提示', {
      type: 'warning'
    }).then(() => {
      appList.value = appList.value.filter(a => a.id !== item.id)
      ElMessage.success('删除成功')
    }).catch(() => {})
  }
}

const handleSubmit = async () => {
  await formRef.value.validate()
  const business = businessList.value.find(b => b.id === form.businessId)
  if (isEdit.value) {
    const index = appList.value.findIndex(a => a.id === form.id)
    if (index > -1) {
      appList.value[index] = { 
        ...appList.value[index], 
        ...form,
        businessName: business?.name 
      }
    }
    ElMessage.success('更新成功')
  } else {
    appList.value.push({
      ...form,
      id: Date.now(),
      businessName: business?.name,
      serviceCount: 0,
      hostCount: 0,
      dbCount: 0
    })
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
}

const viewDetail = (app) => {
  router.push({ path: '/cmdb/services', query: { app_id: app.id } })
}
</script>

<style scoped>
.application-page {
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

.app-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.app-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  cursor: pointer;
  transition: all 0.3s;
}

.app-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.card-icon {
  width: 44px;
  height: 44px;
  background: #f0f5ff;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  color: #667eea;
}

.card-title {
  flex: 1;
}

.card-title h3 {
  margin: 0 0 6px;
  font-size: 16px;
  color: #333;
}

.card-tags {
  display: flex;
  gap: 6px;
}

.tag {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
}

.tag.business {
  background: #f0f5ff;
  color: #667eea;
}

.tag.env.production {
  background: #f6ffed;
  color: #52c41a;
}

.tag.env.staging {
  background: #fff7e6;
  color: #faad14;
}

.tag.env.development {
  background: #e6f7ff;
  color: #1890ff;
}

.card-desc {
  font-size: 13px;
  color: #666;
  margin-bottom: 16px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-resources {
  display: flex;
  gap: 20px;
  padding: 12px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 12px;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.resource-icon {
  font-size: 14px;
  color: #667eea;
}

.resource-count {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.resource-label {
  font-size: 12px;
  color: #999;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-owner {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
}

.owner-avatar {
  font-size: 14px;
  color: #999;
}

.card-level {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
}

.card-level.core {
  background: #fff1f0;
  color: #ff4d4f;
}

.card-level.important {
  background: #fff7e6;
  color: #faad14;
}

.card-level.normal {
  background: #f5f5f5;
  color: #999;
}

/* 空状态 */
.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 20px;
  background: #fff;
  border-radius: 12px;
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
</style>
