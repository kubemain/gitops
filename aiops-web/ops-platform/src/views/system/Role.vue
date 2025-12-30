<template>



  <div class="role-container">
    <el-card>
      <div class="header">
        <h3>角色列表</h3>
      </div>

      <el-table
          :data="tableData"
          style="width: 100%; margin-top: 20px"
          v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="code" label="角色编码" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="系统角色" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.is_system === 1" type="warning" size="small">
              是
            </el-tag>
            <span v-else>否</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button
                link
                type="primary"
                size="small"
                @click="handlePermission(row)"
                v-if="userStore.hasPermission('system:role:assign')"
            >
              分配权限
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 权限分配对话框 -->
    <el-dialog
        v-model="permissionDialogVisible"
        title="分配权限"
        width="600px"
    >
      <div style="margin-bottom: 10px">
        <strong>角色：</strong>{{ currentRole?.name }}
      </div>

      <el-tree
          ref="treeRef"
          :data="permissionTree"
          show-checkbox
          node-key="id"
          :default-checked-keys="checkedPermissions"
          :props="{ label: 'name', children: 'children' }"
      />

      <template #footer>
        <el-button @click="permissionDialogVisible = false">取消</el-button>
        <el-button
            type="primary"
            @click="handleSavePermission"
            :loading="submitLoading"
        >
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import {
  getRoleList,
  getRoleDetail,
  assignPermissions,
  getAllPermissions
} from '@/api/role'

const userStore = useUserStore()

const tableData = ref([])
const loading = ref(false)

// 权限相关
const permissionDialogVisible = ref(false)
const currentRole = ref(null)
const permissionTree = ref([])
const checkedPermissions = ref([])
const submitLoading = ref(false)
const treeRef = ref(null)

// 获取角色列表
const fetchData = async () => {
  loading.value = true
  try {
    const data = await getRoleList()
    tableData.value = data || []
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  } finally {
    loading.value = false
  }
}

// 获取权限树
const fetchPermissions = async () => {
  try {
    const data = await getAllPermissions()
    permissionTree.value = buildTree(data || [])
  } catch (error) {
    console.error('获取权限列表失败:', error)
  }
}

// 构建树形结构
const buildTree = (list) => {
  const map = {}
  const tree = []

  // 先建立映射
  list.forEach(item => {
    map[item.id] = { ...item, children: [] }
  })

  // 构建树
  list.forEach(item => {
    if (item.parent_id === 0) {
      tree.push(map[item.id])
    } else if (map[item.parent_id]) {
      map[item.parent_id].children.push(map[item.id])
    }
  })

  return tree
}

// 分配权限
const handlePermission = async (row) => {
  currentRole.value = row
  permissionDialogVisible.value = true

  try {
    const data = await getRoleDetail(row.id)
    checkedPermissions.value = data.permissions?.map(p => p.id) || []
  } catch (error) {
    console.error('获取角色权限失败:', error)
  }
}

// 保存权限
const handleSavePermission = async () => {
  if (!currentRole.value) return

  submitLoading.value = true
  try {
    // 获取选中的节点（包括半选中的父节点）
    const checkedKeys = treeRef.value.getCheckedKeys()
    const halfCheckedKeys = treeRef.value.getHalfCheckedKeys()
    const allKeys = [...checkedKeys, ...halfCheckedKeys]

    await assignPermissions(currentRole.value.id, {
      permission_ids: allKeys
    })

    ElMessage.success('权限分配成功')
    permissionDialogVisible.value = false
  } catch (error) {
    console.error('分配权限失败:', error)
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  fetchData()
  fetchPermissions()
})
</script>

<style scoped>
.role-container {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h3 {
  margin: 0;
}
</style>
