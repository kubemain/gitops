<!-- 数据导入页面 -->
<template>
  <div class="data-import">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="title-icon"><Upload /></el-icon>
            <h2>数据导入</h2>
            <el-tag type="info" size="small">批量导入主机资源数据</el-tag>
          </div>
        </div>
      </template>

      <div class="import-container">
        <!-- 导入方式选择 -->
        <div class="import-methods">
          <div 
            class="method-card" 
            :class="{ active: importMethod === 'excel' }"
            @click="importMethod = 'excel'"
          >
            <div class="method-icon excel">
              <el-icon :size="32"><Document /></el-icon>
            </div>
            <div class="method-info">
              <h4>Excel导入</h4>
              <p>支持 .xlsx, .xls 格式</p>
            </div>
          </div>

          <div 
            class="method-card" 
            :class="{ active: importMethod === 'csv' }"
            @click="importMethod = 'csv'"
          >
            <div class="method-icon csv">
              <el-icon :size="32"><List /></el-icon>
            </div>
            <div class="method-info">
              <h4>CSV导入</h4>
              <p>支持 .csv 格式</p>
            </div>
          </div>

          <div 
            class="method-card" 
            :class="{ active: importMethod === 'json' }"
            @click="importMethod = 'json'"
          >
            <div class="method-icon json">
              <el-icon :size="32"><Tickets /></el-icon>
            </div>
            <div class="method-info">
              <h4>JSON导入</h4>
              <p>支持 .json 格式</p>
            </div>
          </div>

          <div 
            class="method-card" 
            :class="{ active: importMethod === 'api' }"
            @click="importMethod = 'api'"
          >
            <div class="method-icon api">
              <el-icon :size="32"><Connection /></el-icon>
            </div>
            <div class="method-info">
              <h4>API同步</h4>
              <p>从外部系统同步</p>
            </div>
          </div>
        </div>

        <!-- 文件上传区域 -->
        <div v-if="importMethod !== 'api'" class="upload-area">
          <el-upload
            ref="uploadRef"
            class="upload-dragger"
            drag
            :auto-upload="false"
            :limit="1"
            :accept="acceptTypes"
            :on-change="handleFileChange"
            :on-exceed="handleExceed"
          >
            <el-icon class="upload-icon"><UploadFilled /></el-icon>
            <div class="upload-text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="upload-tip">
                支持 {{ acceptTypes }} 格式，单个文件不超过 10MB
              </div>
            </template>
          </el-upload>

          <div class="template-download">
            <el-button type="primary" link @click="downloadTemplate">
              <el-icon><Download /></el-icon>
              下载导入模板
            </el-button>
          </div>
        </div>

        <!-- API同步配置 -->
        <div v-else class="api-config">
          <el-form :model="apiConfig" label-width="100px">
            <el-form-item label="数据源">
              <el-select v-model="apiConfig.source" placeholder="选择数据源" style="width: 100%">
                <el-option label="云平台API" value="cloud" />
                <el-option label="监控系统" value="monitor" />
                <el-option label="自定义API" value="custom" />
              </el-select>
            </el-form-item>
            <el-form-item v-if="apiConfig.source === 'custom'" label="API地址">
              <el-input v-model="apiConfig.url" placeholder="请输入API地址" />
            </el-form-item>
            <el-form-item label="认证方式">
              <el-radio-group v-model="apiConfig.authType">
                <el-radio label="none">无认证</el-radio>
                <el-radio label="token">Token</el-radio>
                <el-radio label="basic">Basic Auth</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-if="apiConfig.authType === 'token'" label="Token">
              <el-input v-model="apiConfig.token" type="password" show-password placeholder="请输入Token" />
            </el-form-item>
          </el-form>
        </div>

        <!-- 数据预览 -->
        <div v-if="previewData.length > 0" class="preview-area">
          <div class="preview-header">
            <h3>数据预览</h3>
            <el-tag type="success">共 {{ previewData.length }} 条记录</el-tag>
          </div>
          <el-table :data="previewData.slice(0, 10)" style="width: 100%" max-height="300">
            <el-table-column prop="hostname" label="主机名" width="150" />
            <el-table-column prop="ip" label="IP地址" width="140" />
            <el-table-column prop="os" label="操作系统" width="120" />
            <el-table-column prop="cpu" label="CPU" width="80" />
            <el-table-column prop="memory" label="内存" width="100" />
            <el-table-column prop="environment" label="环境" width="100" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.valid ? 'success' : 'danger'" size="small">
                  {{ row.valid ? '有效' : '无效' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
          <div v-if="previewData.length > 10" class="preview-more">
            还有 {{ previewData.length - 10 }} 条数据未显示...
          </div>
        </div>

        <!-- 导入选项 -->
        <div class="import-options">
          <h3>导入选项</h3>
          <el-form :model="importOptions" label-width="140px">
            <el-form-item label="重复数据处理">
              <el-radio-group v-model="importOptions.duplicateAction">
                <el-radio label="skip">跳过</el-radio>
                <el-radio label="update">更新</el-radio>
                <el-radio label="error">报错</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="默认主机组">
              <el-select v-model="importOptions.defaultGroupId" placeholder="选择默认分组" clearable style="width: 200px">
                <el-option v-for="group in hostGroups" :key="group.id" :label="group.name" :value="group.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="默认环境">
              <el-select v-model="importOptions.defaultEnvironment" placeholder="选择默认环境" clearable style="width: 200px">
                <el-option label="生产环境" value="production" />
                <el-option label="预发环境" value="staging" />
                <el-option label="开发环境" value="development" />
              </el-select>
            </el-form-item>
          </el-form>
        </div>

        <!-- 操作按钮 -->
        <div class="import-actions">
          <el-button @click="resetImport">重置</el-button>
          <el-button type="primary" :loading="importing" :disabled="!canImport" @click="handleImport">
            <el-icon><Upload /></el-icon>
            开始导入
          </el-button>
        </div>

        <!-- 导入结果 -->
        <div v-if="importResult" class="import-result">
          <el-result
            :icon="importResult.success ? 'success' : 'warning'"
            :title="importResult.success ? '导入完成' : '导入完成（有警告）'"
          >
            <template #sub-title>
              <div class="result-stats">
                <span class="stat success">成功: {{ importResult.successCount }}</span>
                <span class="stat fail">失败: {{ importResult.failCount }}</span>
                <span class="stat skip">跳过: {{ importResult.skipCount }}</span>
              </div>
            </template>
            <template #extra>
              <el-button type="primary" @click="$router.push('/cmdb/hosts')">查看主机列表</el-button>
              <el-button @click="resetImport">继续导入</el-button>
            </template>
          </el-result>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAllHostGroups } from '@/api/host-group'
import {
  Upload, Document, List, Tickets, Connection,
  UploadFilled, Download
} from '@element-plus/icons-vue'

const importMethod = ref('excel')
const uploadRef = ref(null)
const selectedFile = ref(null)
const previewData = ref([])
const hostGroups = ref([])
const importing = ref(false)
const importResult = ref(null)

const apiConfig = ref({
  source: '',
  url: '',
  authType: 'none',
  token: ''
})

const importOptions = ref({
  duplicateAction: 'skip',
  defaultGroupId: '',
  defaultEnvironment: ''
})

const acceptTypes = computed(() => {
  const types = {
    excel: '.xlsx,.xls',
    csv: '.csv',
    json: '.json'
  }
  return types[importMethod.value] || ''
})

const canImport = computed(() => {
  if (importMethod.value === 'api') {
    return apiConfig.value.source !== ''
  }
  return selectedFile.value !== null
})

// 获取主机分组
const fetchHostGroups = async () => {
  try {
    const data = await getAllHostGroups()
    hostGroups.value = data || []
  } catch (error) {
    console.error('获取主机分组失败:', error)
  }
}

const handleFileChange = (file) => {
  selectedFile.value = file.raw
  // 模拟解析预览数据
  previewData.value = [
    { hostname: 'new-server-01', ip: '192.168.2.10', os: 'CentOS 7.9', cpu: 8, memory: 16384, environment: 'production', valid: true },
    { hostname: 'new-server-02', ip: '192.168.2.11', os: 'Ubuntu 20.04', cpu: 4, memory: 8192, environment: 'staging', valid: true },
    { hostname: 'new-server-03', ip: '192.168.2.12', os: 'CentOS 8', cpu: 16, memory: 32768, environment: 'production', valid: true }
  ]
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const downloadTemplate = () => {
  ElMessage.success('模板下载中...')
  // 实际实现下载逻辑
}

const handleImport = async () => {
  importing.value = true
  try {
    // 模拟导入过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    importResult.value = {
      success: true,
      successCount: previewData.value.length,
      failCount: 0,
      skipCount: 0
    }
    ElMessage.success('导入成功')
  } catch (error) {
    ElMessage.error('导入失败')
  } finally {
    importing.value = false
  }
}

const resetImport = () => {
  selectedFile.value = null
  previewData.value = []
  importResult.value = null
  uploadRef.value?.clearFiles()
}

onMounted(() => {
  fetchHostGroups()
})
</script>

<style scoped>
.data-import {
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
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 24px;
  color: #fa8c16;
}

.header-left h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.import-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 导入方式 */
.import-methods {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.method-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  border: 2px solid #e8e8e8;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.method-card:hover {
  border-color: #1890ff;
}

.method-card.active {
  border-color: #1890ff;
  background: #e6f7ff;
}

.method-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.method-icon.excel { background: linear-gradient(135deg, #52c41a, #95de64); }
.method-icon.csv { background: linear-gradient(135deg, #1890ff, #69c0ff); }
.method-icon.json { background: linear-gradient(135deg, #722ed1, #b37feb); }
.method-icon.api { background: linear-gradient(135deg, #fa8c16, #ffc069); }

.method-info h4 {
  margin: 0 0 4px;
  font-size: 15px;
  font-weight: 600;
}

.method-info p {
  margin: 0;
  font-size: 12px;
  color: #999;
}

/* 上传区域 */
.upload-area {
  text-align: center;
}

.upload-dragger {
  width: 100%;
}

.upload-dragger :deep(.el-upload-dragger) {
  padding: 40px;
  border-radius: 12px;
}

.upload-icon {
  font-size: 48px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.upload-text {
  font-size: 14px;
  color: #666;
}

.upload-text em {
  color: #1890ff;
  font-style: normal;
}

.upload-tip {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}

.template-download {
  margin-top: 16px;
}

/* API配置 */
.api-config {
  max-width: 600px;
  padding: 20px;
  background: #fafafa;
  border-radius: 12px;
}

/* 预览区域 */
.preview-area {
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  padding: 16px;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.preview-header h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
}

.preview-more {
  text-align: center;
  padding: 12px;
  color: #999;
  font-size: 13px;
}

/* 导入选项 */
.import-options h3 {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 16px;
}

/* 操作按钮 */
.import-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

/* 导入结果 */
.result-stats {
  display: flex;
  gap: 24px;
  justify-content: center;
}

.result-stats .stat {
  font-size: 14px;
}

.result-stats .success { color: #52c41a; }
.result-stats .fail { color: #ff4d4f; }
.result-stats .skip { color: #faad14; }

@media (max-width: 1200px) {
  .import-methods {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .import-methods {
    grid-template-columns: 1fr;
  }
}
</style>
