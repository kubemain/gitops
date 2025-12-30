<!-- src/views/system/User.vue -->
<template>
  <div class="user-page">
    <el-card class="page-card">
      <!-- 页面标题 -->
      <div class="page-header">
        <div class="header-title">
          <el-icon class="title-icon"><User /></el-icon>
          <h2>用户管理</h2>
        </div>
        <div class="header-breadcrumb">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>系统管理</el-breadcrumb-item>
            <el-breadcrumb-item>用户管理</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
      </div>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <el-input
            v-model="searchForm.username"
            placeholder="搜索用户名"
            style="width: 200px"
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
        <el-button
            type="primary"
            @click="handleCreate"
            v-if="userStore.hasPermission('system:user:create')"
        >
          <el-icon><Plus /></el-icon>
          新增用户
        </el-button>
      </div>

      <!-- 用户表格 -->
      <el-table
          :data="tableData"
          style="width: 100%; margin-top: 20px"
          v-loading="loading"
          :header-cell-style="{ background: '#fafafa', color: '#333' }"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column label="角色" width="200">
          <template #default="{ row }">
            <el-tag
                v-for="role in row.roles"
                :key="role.id"
                style="margin-right: 5px"
                size="small"
                type="success"
            >
              {{ role.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button
                link
                type="primary"
                size="small"
                @click="handleEdit(row)"
                v-if="userStore.hasPermission('system:user:edit')"
            >
              编辑
            </el-button>
            <el-button
                link
                type="warning"
                size="small"
                @click="handleResetPassword(row)"
                v-if="userStore.hasPermission('system:user:edit')"
            >
              重置密码
            </el-button>
            <el-button
                link
                type="danger"
                size="small"
                @click="handleDelete(row)"
                v-if="userStore.hasPermission('system:user:delete') && row.id !== 1"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.page_size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
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
        width="600px"
        @close="handleDialogClose"
    >
      <el-form
          ref="formRef"
          :model="formData"
          :rules="formRules"
          label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
              v-model="formData.username"
              placeholder="请输入用户名"
              :disabled="isEdit"
          />
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="formData.nickname" placeholder="请输入昵称" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>

        <el-form-item label="手机号" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>

        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input
              v-model="formData.password"
              type="password"
              placeholder="请输入密码（至少6位）"
              show-password
          />
        </el-form-item>

        <el-form-item label="角色" prop="role_ids">
          <el-select
              v-model="formData.role_ids"
              multiple
              placeholder="请选择角色"
              style="width: 100%"
          >
            <el-option
                v-for="role in roleList"
                :key="role.id"
                :label="role.name"
                :value="role.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Plus, User } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import {
  getUserList,
  createUser,
  updateUser,
  deleteUser,
  resetPassword
} from '@/api/user'
import { getRoleList } from '@/api/role'

const userStore = useUserStore()

// 表格数据
const tableData = ref([])
const loading = ref(false)

// 分页
const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  username: ''
})

// 角色列表
const roleList = ref([])

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const isEdit = ref(false)
const submitLoading = ref(false)

const formRef = ref(null)
const formData = reactive({
  id: null,
  username: '',
  nickname: '',
  email: '',
  phone: '',
  password: '',
  role_ids: [],
  status: 1
})

// 表单验证规则
const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取用户列表
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.page_size,
      username: searchForm.username || undefined
    }
    const data = await getUserList(params)
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 获取角色列表
const fetchRoleList = async () => {
  try {
    const data = await getRoleList()
    roleList.value = data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
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
  dialogTitle.value = '新增用户'
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑用户'

  formData.id = row.id
  formData.username = row.username
  formData.nickname = row.nickname
  formData.email = row.email
  formData.phone = row.phone
  formData.role_ids = row.roles?.map(r => r.id) || []
  formData.status = row.status

  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
        `确定要删除用户「${row.username}」吗？`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
    )

    await deleteUser(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 重置密码
const handleResetPassword = async (row) => {
  try {
    const { value } = await ElMessageBox.prompt(
        `请输入用户「${row.username}」的新密码（至少6位）`,
        '重置密码',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          inputPattern: /.{6,}/,
          inputErrorMessage: '密码长度不能少于6位'
        }
    )

    await resetPassword(row.id, { new_password: value })
    ElMessage.success('密码重置成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置密码失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (isEdit.value) {
        await updateUser(formData.id, {
          nickname: formData.nickname,
          email: formData.email,
          phone: formData.phone,
          role_ids: formData.role_ids,
          status: formData.status
        })
        ElMessage.success('更新成功')
      } else {
        await createUser({
          username: formData.username,
          nickname: formData.nickname,
          email: formData.email,
          phone: formData.phone,
          password: formData.password,
          role_ids: formData.role_ids,
          status: formData.status
        })
        ElMessage.success('创建成功')
      }

      dialogVisible.value = false
      fetchData()
    } catch (error) {
      console.error('提交失败:', error)
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
    username: '',
    nickname: '',
    email: '',
    phone: '',
    password: '',
    role_ids: [],
    status: 1
  })
}

onMounted(() => {
  fetchData()
  fetchRoleList()
})
</script>

<style scoped>
.user-page {
  width: 100%;
  height: 100%;
}

.page-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* ========== 页面标题 ========== */
.page-header {
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.title-icon {
  font-size: 24px;
  color: #1890ff;
}

.header-title h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.header-breadcrumb {
  font-size: 14px;
}

/* ========== 搜索栏 ========== */
.search-bar {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 16px;
}

/* ========== 分页 ========== */
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
