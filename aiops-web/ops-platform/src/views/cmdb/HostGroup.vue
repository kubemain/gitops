<template>
  <div class="host-group-page">
    <el-card class="page-card">
      <!-- 页面标题 -->
      <div class="page-header">
        <div class="header-title">
          <el-icon class="title-icon"><FolderOpened /></el-icon>
          <h2>主机分组</h2>
          <el-tag type="info" size="small">共 {{ pagination.total }} 个分组</el-tag>
        </div>
        <div class="header-actions">
          <el-radio-group v-model="viewMode" size="small">
            <el-radio-button label="tree">
              <el-icon><Share /></el-icon>
              树形
            </el-radio-button>
            <el-radio-button label="table">
              <el-icon><List /></el-icon>
              列表
            </el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
          v-model="searchForm.name"
          placeholder="搜索分组名称"
          style="width: 220px"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          新增分组
        </el-button>
      </div>

      <!-- 树形视图 -->
      <template v-if="viewMode === 'tree'">
        <div class="tree-container" v-loading="loading">
          <div class="tree-grid">
            <div 
              v-for="group in tableData" 
              :key="group.id" 
              class="group-card"
              @click="handleViewHosts(group)"
            >
              <div class="group-icon">
                <el-icon :size="32"><Folder /></el-icon>
              </div>
              <div class="group-info">
                <h4>{{ group.name }}</h4>
                <p class="group-desc">{{ group.description || '暂无描述' }}</p>
              </div>
              <div class="group-stats">
                <div class="stat-item">
                  <span class="stat-value">{{ group.host_count || 0 }}</span>
                  <span class="stat-label">主机数</span>
                </div>
              </div>
              <div class="group-actions" @click.stop>
                <el-button link type="primary" size="small" @click="handleEdit(group)">
                  <el-icon><Edit /></el-icon>
                </el-button>
                <el-button link type="danger" size="small" @click="handleDelete(group)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
            <el-empty v-if="!loading && tableData.length === 0" description="暂无分组" />
          </div>
        </div>
      </template>

      <!-- 表格视图 -->
      <template v-else>
        <el-table
          :data="tableData"
          style="width: 100%"
          v-loading="loading"
          :header-cell-style="{ background: '#fafafa', color: '#333' }"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="name" label="分组名称" width="200">
            <template #default="{ row }">
              <div class="group-name-cell">
                <el-icon class="folder-icon"><Folder /></el-icon>
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="250">
            <template #default="{ row }">
              <span class="desc-text">{{ row.description || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="host_count" label="主机数量" width="120">
            <template #default="{ row }">
              <el-tag 
                size="small" 
                :type="row.host_count > 0 ? 'success' : 'info'"
                class="count-tag"
              >
                {{ row.host_count || 0 }} 台
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="handleViewHosts(row)">
                查看主机
              </el-button>
              <el-button link type="primary" size="small" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button link type="danger" size="small" @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </template>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="pagination.total"
          :page-sizes="[12, 24, 48]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分组名称" />
        </el-form-item>

        <el-form-item label="父级分组" prop="parent_id">
          <el-select 
            v-model="formData.parent_id" 
            placeholder="选择父级分组（可选）" 
            style="width: 100%"
            clearable
          >
            <el-option
              v-for="group in parentGroups"
              :key="group.id"
              :label="group.name"
              :value="group.id"
              :disabled="group.id === formData.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="999" style="width: 100%" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            rows="3"
            placeholder="请输入分组描述"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { Search, Plus, FolderOpened, Folder, Share, List, Edit, Delete } from '@element-plus/icons-vue'
import {
  getHostGroupList,
  getAllHostGroups,
  createHostGroup,
  updateHostGroup,
  deleteHostGroup
} from '@/api/host-group'

const router = useRouter()

// 视图模式
const viewMode = ref('tree')

// 表格数据
const tableData = ref([])
const allGroups = ref([])
const loading = ref(false)

// 分页
const pagination = reactive({
  page: 1,
  page_size: 12,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  name: ''
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('新增分组')
const isEdit = ref(false)
const submitLoading = ref(false)

const formRef = ref(null)
const formData = reactive({
  id: null,
  name: '',
  description: '',
  parent_id: null,
  sort: 0
})

// 可选的父级分组
const parentGroups = computed(() => {
  return allGroups.value.filter(g => g.id !== formData.id)
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入分组名称', trigger: 'blur' },
    { min: 2, max: 50, message: '名称长度在 2 到 50 个字符', trigger: 'blur' }
  ]
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取分组列表
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      name: searchForm.name || undefined
    }
    const data = await getHostGroupList(params)
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取分组列表失败:', error)
    ElMessage.error('获取分组列表失败')
  } finally {
    loading.value = false
  }
}

// 获取所有分组（用于父级选择）
const fetchAllGroups = async () => {
  try {
    const data = await getAllHostGroups()
    allGroups.value = data || []
  } catch (error) {
    console.error('获取所有分组失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 新增
const handleCreate = () => {
  isEdit.value = false
  dialogTitle.value = '新增分组'
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑分组'
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    description: row.description || '',
    parent_id: row.parent_id || null,
    sort: row.sort || 0
  })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    if (row.host_count && row.host_count > 0) {
      ElMessage.warning('该分组下还有主机，无法删除')
      return
    }

    await ElMessageBox.confirm(`确定要删除分组「${row.name}」吗？`, '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await deleteHostGroup(row.id)
    ElMessage.success('删除成功')
    fetchData()
    fetchAllGroups()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 查看主机
const handleViewHosts = (row) => {
  router.push({
    path: '/cmdb/hosts',
    query: { group_id: row.id }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      const submitData = {
        name: formData.name,
        description: formData.description,
        parent_id: formData.parent_id || 0,
        sort: formData.sort
      }

      if (isEdit.value) {
        await updateHostGroup(formData.id, submitData)
        ElMessage.success('更新成功')
      } else {
        await createHostGroup(submitData)
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      fetchData()
      fetchAllGroups()
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
    id: null,
    name: '',
    description: '',
    parent_id: null,
    sort: 0
  })
}

onMounted(() => {
  fetchData()
  fetchAllGroups()
})
</script>

<style scoped>
.host-group-page {
  height: 100%;
}

.page-card {
  height: 100%;
  border-radius: 12px;
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
  color: #faad14;
}

.header-title h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

/* ========== 搜索栏 ========== */
.search-bar {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 16px;
}

/* ========== 树形视图 ========== */
.tree-container {
  flex: 1;
  overflow-y: auto;
}

.tree-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
  padding: 4px;
}

.group-card {
  background: #fff;
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
}

.group-card:hover {
  border-color: #faad14;
  box-shadow: 0 4px 12px rgba(250, 173, 20, 0.15);
}

.group-icon {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #faad14, #ffc069);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.group-info h4 {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.group-desc {
  margin: 0;
  font-size: 13px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-stats {
  display: flex;
  gap: 20px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1890ff;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.group-actions {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s;
}

.group-card:hover .group-actions {
  opacity: 1;
}

/* ========== 表格视图 ========== */
.group-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.folder-icon {
  color: #faad14;
  font-size: 18px;
}

.desc-text {
  color: #666;
}

.count-tag {
  font-weight: 500;
}

/* ========== 分页 ========== */
.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

/* ========== 响应式 ========== */
@media (max-width: 768px) {
  .tree-grid {
    grid-template-columns: 1fr;
  }
}
</style>
