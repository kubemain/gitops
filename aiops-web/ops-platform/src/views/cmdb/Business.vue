<!-- 业务线管理 -->
<template>
  <div class="business-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-info">
        <h2><el-icon><OfficeBuilding /></el-icon> 业务线管理</h2>
        <p>定义产品或业务单元，作为CMDB关系链的顶层</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          新建业务线
        </el-button>
      </div>
    </div>

    <!-- 业务线列表 -->
    <div class="business-list">
      <div 
        v-for="item in businessList" 
        :key="item.id" 
        class="business-card"
        @click="viewDetail(item)"
      >
        <div class="card-header">
          <div class="card-icon" :style="{ background: item.color }">
            <el-icon><OfficeBuilding /></el-icon>
          </div>
          <div class="card-title">
            <h3>{{ item.name }}</h3>
            <span class="card-code">{{ item.code }}</span>
          </div>
          <el-dropdown @command="(cmd) => handleCommand(cmd, item)" @click.stop>
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
        
        <div class="card-desc">{{ item.description || '暂无描述' }}</div>
        
        <div class="card-stats">
          <div class="stat-item">
            <span class="stat-value">{{ item.appCount }}</span>
            <span class="stat-label">应用</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ item.serviceCount }}</span>
            <span class="stat-label">服务</span>
          </div>
          <div class="stat-item">
            <span class="stat-value">{{ item.hostCount }}</span>
            <span class="stat-label">主机</span>
          </div>
        </div>

        <div class="card-footer">
          <div class="card-owner">
            <span class="owner-label">负责人:</span>
            <span class="owner-name">{{ item.owner }}</span>
          </div>
          <div class="card-status" :class="item.status">
            {{ item.status === 'active' ? '运行中' : '已停用' }}
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="businessList.length === 0" class="empty-state">
        <el-icon class="empty-icon"><OfficeBuilding /></el-icon>
        <div class="empty-title">还没有业务线</div>
        <div class="empty-desc">创建第一个业务线，开始构建你的CMDB关系链</div>
        <el-button type="primary" @click="showAddDialog">
          <el-icon><Plus /></el-icon>
          新建业务线
        </el-button>
      </div>
    </div>

    <!-- 新建/编辑对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑业务线' : '新建业务线'"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：电商平台、金融系统" />
        </el-form-item>
        <el-form-item label="编码" prop="code">
          <el-input v-model="form.code" placeholder="如：ecommerce、finance" />
        </el-form-item>
        <el-form-item label="负责人" prop="owner">
          <el-input v-model="form.owner" placeholder="业务负责人" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="业务线描述" />
        </el-form-item>
        <el-form-item label="图标">
          <div class="icon-picker">
            <span 
              v-for="iconName in icons" 
              :key="iconName" 
              class="icon-option"
              :class="{ active: form.icon === iconName }"
              @click="form.icon = iconName"
            >
              <el-icon><component :is="iconName" /></el-icon>
            </span>
          </div>
        </el-form-item>
        <el-form-item label="颜色">
          <div class="color-picker">
            <span 
              v-for="color in colors" 
              :key="color" 
              class="color-option"
              :class="{ active: form.color === color }"
              :style="{ background: color }"
              @click="form.color = color"
            ></span>
          </div>
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
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, MoreFilled, OfficeBuilding } from '@element-plus/icons-vue'

const router = useRouter()

// 业务线列表（模拟数据）
const businessList = ref([
  {
    id: 1,
    name: '电商平台',
    code: 'ecommerce',
    description: '公司核心电商业务，包含商城、订单、支付等系统',
    icon: 'Shop',
    color: '#667eea',
    owner: '张三',
    status: 'active',
    appCount: 5,
    serviceCount: 18,
    hostCount: 32
  },
  {
    id: 2,
    name: '金融系统',
    code: 'finance',
    description: '金融相关业务，包含账户、清算、风控等',
    icon: 'Money',
    color: '#52c41a',
    owner: '李四',
    status: 'active',
    appCount: 3,
    serviceCount: 12,
    hostCount: 20
  },
  {
    id: 3,
    name: '数据平台',
    code: 'data-platform',
    description: '大数据分析平台，提供数据采集、处理、分析能力',
    icon: 'DataAnalysis',
    color: '#faad14',
    owner: '王五',
    status: 'active',
    appCount: 2,
    serviceCount: 8,
    hostCount: 15
  }
])

// 对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const form = reactive({
  id: null,
  name: '',
  code: '',
  owner: '',
  description: '',
  icon: 'OfficeBuilding',
  color: '#667eea'
})

const rules = {
  name: [{ required: true, message: '请输入业务线名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入业务线编码', trigger: 'blur' }],
  owner: [{ required: true, message: '请输入负责人', trigger: 'blur' }]
}

const icons = ['OfficeBuilding', 'Shop', 'Money', 'DataAnalysis', 'Goods', 'Iphone', 'Connection', 'Setting', 'Box', 'Promotion']
const colors = ['#667eea', '#52c41a', '#faad14', '#1890ff', '#722ed1', '#eb2f96', '#13c2c2', '#fa541c']

const showAddDialog = () => {
  isEdit.value = false
  Object.assign(form, { id: null, name: '', code: '', owner: '', description: '', icon: 'OfficeBuilding', color: '#667eea' })
  dialogVisible.value = true
}

const handleCommand = (command, item) => {
  if (command === 'edit') {
    isEdit.value = true
    Object.assign(form, item)
    dialogVisible.value = true
  } else if (command === 'delete') {
    ElMessageBox.confirm(`确定删除业务线「${item.name}」吗？`, '提示', {
      type: 'warning'
    }).then(() => {
      businessList.value = businessList.value.filter(b => b.id !== item.id)
      ElMessage.success('删除成功')
    }).catch(() => {})
  }
}

const handleSubmit = async () => {
  await formRef.value.validate()
  if (isEdit.value) {
    const index = businessList.value.findIndex(b => b.id === form.id)
    if (index > -1) {
      businessList.value[index] = { ...businessList.value[index], ...form }
    }
    ElMessage.success('更新成功')
  } else {
    businessList.value.push({
      ...form,
      id: Date.now(),
      status: 'active',
      appCount: 0,
      serviceCount: 0,
      hostCount: 0
    })
    ElMessage.success('创建成功')
  }
  dialogVisible.value = false
}

const viewDetail = (item) => {
  router.push({ path: '/cmdb/applications', query: { business_id: item.id } })
}
</script>

<style scoped>
.business-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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

.business-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.business-card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
  cursor: pointer;
  transition: all 0.3s;
}

.business-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #fff;
}

.card-title {
  flex: 1;
}

.card-title h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
}

.card-code {
  font-size: 12px;
  color: #999;
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

.card-stats {
  display: flex;
  gap: 24px;
  padding: 16px 0;
  border-top: 1px solid #f0f0f0;
  border-bottom: 1px solid #f0f0f0;
  margin-bottom: 12px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-owner {
  font-size: 12px;
  color: #666;
}

.owner-label {
  color: #999;
}

.card-status {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
}

.card-status.active {
  background: #f6ffed;
  color: #52c41a;
}

.card-status.inactive {
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

/* 图标选择器 */
.icon-picker {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.icon-option {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  border: 2px solid #f0f0f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.icon-option:hover {
  border-color: #667eea;
}

.icon-option.active {
  border-color: #667eea;
  background: #f0f5ff;
}

/* 颜色选择器 */
.color-picker {
  display: flex;
  gap: 8px;
}

.color-option {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid transparent;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.active {
  border-color: #333;
  box-shadow: 0 0 0 2px #fff, 0 0 0 4px currentColor;
}
</style>
