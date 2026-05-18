<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconPlayArrow, IconRefresh } from '@arco-design/web-vue/es/icon'
import jobApi from '@/api/job'
import CronPicker from '@/components/Cron/index.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const isMobile = ref(window.innerWidth < 576)
function onResize() { isMobile.value = window.innerWidth < 576 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

// ── 搜索 ────────────────────────────────────────────────
const filterModule = ref(false)
const advancedVisible = ref(false)
const total = ref(0)
const loading = ref(false)
const dataList = ref<any[]>([])
const queryParams = ref<{ pageNum: number; pageSize: number; jobName?: string; status?: number | null }>({
  pageNum: 1, pageSize: 20
})

async function loadPage() {
  loading.value = true
  try {
    const res: any = await jobApi.page(queryParams.value)
    dataList.value = res?.rows ?? []
    total.value = res?.total ?? 0
  } finally {
    loading.value = false
  }
}
function handleSearch() { queryParams.value.pageNum = 1; loadPage() }
function handleReset() {
  queryParams.value = { pageNum: 1, pageSize: queryParams.value.pageSize }
  handleSearch()
}
watch(() => ({ ...queryParams.value }), () => {
  const q = queryParams.value
  filterModule.value = !!(q.jobName || q.status != null)
}, { deep: true })

// ── 多选 ───────────────────────────────────────────────
const selectedIds = ref<string[]>([])
const rowSelection = reactive({ type: 'checkbox', showCheckedAll: true, onlyCurrent: false })

async function batchDelete() {
  if (!selectedIds.value.length) { Message.warning(t('common.confirmDeleteContent', { count: 0 }).replace('0 ', '')); return }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await jobApi.delete(selectedIds.value.join(','))
      Message.success(t('common.deleteSuccess'))
      selectedIds.value = []
      loadPage()
    },
  })
}

// ── 新增 / 修改 弹窗 ────────────────────────────────────
const modalVisible = ref(false)
const formRef = ref()
const isEdit = ref(false)
const submitLoading = ref(false)
const form = reactive<any>({
  jobName: '', jobGroup: '', jobKey: '', cronExpr: '',
  invokeTarget: '', status: 0, concurrent: 0, retryCount: 0, timeout: 30, remark: ''
})

function handleAdd() {
  isEdit.value = false
  Object.assign(form, { id: undefined, jobName: '', jobGroup: '', jobKey: '', cronExpr: '', invokeTarget: '', status: 0, concurrent: 0, retryCount: 0, timeout: 30, remark: '' })
  nextTick(() => formRef.value?.clearValidate())
  modalVisible.value = true
}

async function handleEdit(row: any) {
  isEdit.value = true
  const data: any = await jobApi.getById(row.id)
  Object.assign(form, data)
  nextTick(() => formRef.value?.clearValidate())
  modalVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value?.validate()
  if (valid) return
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await jobApi.update(form.id, form)
      Message.success(t('common.editSuccess'))
    } else {
      await jobApi.add(form)
      Message.success(t('common.addSuccess'))
    }
    modalVisible.value = false
    loadPage()
  } finally {
    submitLoading.value = false
  }
}

// ── 手动触发 ───────────────────────────────────────────
const triggerLoading = ref<Record<string, boolean>>({})
async function handleTrigger(row: any) {
  triggerLoading.value[row.id] = true
  try {
    await jobApi.trigger(row.id)
    Message.success(t('tool.job.triggerSuccess', { name: row.jobName }))
  } finally {
    triggerLoading.value[row.id] = false
  }
}

// ── 启停状态 ───────────────────────────────────────────
async function handleStatusChange(record: any, val: boolean) {
  const status = val ? 0 : 1
  try {
    await jobApi.changeStatus(record.id, status)
    record.status = status
    Message.success(status === 0 ? t('tool.job.statusStarted') : t('tool.job.statusStopped'))
  } catch {
    Message.error(t('tool.job.operFailed'))
  }
}

// ── 删除 ───────────────────────────────────────────────
async function handleDelete(row: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('tool.job.confirmDelete', { name: row.jobName }),
    okButtonProps: { status: 'danger' },
    onOk: async () => { await jobApi.delete(row.id); Message.success(t('common.deleteSuccess')); loadPage() },
  })
}

// ── 日志抽屉 ───────────────────────────────────────────
const logDrawerVisible = ref(false)
const logJobId = ref<string>('')
const logJobName = ref<string>('')
const logTotal = ref(0)
const logLoading = ref(false)
const logList = ref<any[]>([])
const logQuery = ref({ pageNum: 1, pageSize: 20, jobId: '',direction: 'desc',field:'createdAt' })

function handleViewLog(row: any) {
  logJobId.value = row.id
  logJobName.value = row.jobName
  logQuery.value = { pageNum: 1, pageSize: 20, jobId: row.id,direction: 'desc',field:'createdAt' }
  logDrawerVisible.value = true
  loadLogPage()
}

async function loadLogPage() {
  logLoading.value = true
  try {
    const res: any = await jobApi.logPage(logQuery.value)
    logList.value = res?.rows ?? []
    logTotal.value = res?.total ?? 0
  } finally {
    logLoading.value = false
  }
}

const selectedLogIds = ref<string[]>([])
async function batchDeleteLog() {
  if (!selectedLogIds.value.length) { Message.warning(t('common.confirmDeleteContent', { count: 0 }).replace('0 ', '')); return }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedLogIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await jobApi.logDelete(selectedLogIds.value.join(','))
      Message.success(t('common.deleteSuccess'))
      selectedLogIds.value = []
      loadLogPage()
    },
  })
}

function logStatusColor(status: number) {
  if (status === 0) return 'green'
  if (status === 1) return 'red'
  return 'blue'
}
function logStatusText(status: number) {
  if (status === 0) return t('tool.job.logStatusSuccess')
  if (status === 1) return t('tool.job.logStatusFailed')
  return t('tool.job.logStatusRunning')
}

// ── 初始化 ─────────────────────────────────────────────
loadPage()
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">
      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'job$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>{{ t('common.add') }}
          </a-button>
          <a-button v-permission="'job$delete'" status="danger" :disabled="!selectedIds.length" @click="batchDelete">
            <template #icon><IconDelete /></template>{{ t('common.batchDelete') }}
          </a-button>
        </a-space>
        <div class="cl-toolbar-right">
          <a-input-search v-model="queryParams.jobName" :placeholder="t('tool.job.searchPlaceholder')" allow-clear class="cl-toolbar-search"
            @search="handleSearch" @press-enter="handleSearch" />
          <a-popover v-model:popup-visible="advancedVisible" trigger="click" position="br" popup-container="body">
            <a-button shape="circle" :type="filterModule ? 'primary' : 'secondary'">
              <template #icon><IconFilter /></template>
            </a-button>
            <template #content>
              <div class="cl-filter-panel">
                <p class="cl-filter-title">{{ t('common.advancedFilter') }}</p>
                <a-input v-model="queryParams.jobName" :placeholder="t('tool.job.jobName')" allow-clear />
                <a-select v-model="queryParams.status" :placeholder="t('common.status')" allow-clear>
                  <a-option :value="0">{{ t('tool.job.statusNormal') }}</a-option>
                  <a-option :value="1">{{ t('tool.job.statusPaused') }}</a-option>
                </a-select>
                <div class="cl-filter-actions">
                  <a-button size="small" @click="handleReset">{{ t('common.reset') }}</a-button>
                  <a-button size="small" type="primary" @click="() => { handleSearch(); advancedVisible = false }">{{ t('common.search') }}</a-button>
                </div>
              </div>
            </template>
          </a-popover>
        </div>
      </div>

      <!-- 表格 -->
      <div class="cl-table-wrap">
        <a-skeleton v-if="loading" :animation="true"><a-skeleton-line :rows="8" /></a-skeleton>
        <a-table v-else
          :bordered="false" :data="dataList" row-key="id" :scroll="{ x: 1100 }" :row-selection="rowSelection"
          v-model:selectedKeys="selectedIds"
          :pagination="{ current: queryParams.pageNum, pageSize: queryParams.pageSize, total, showTotal: true, showPageSize: true, showJumper: true, pageSizeOptions: [10, 20, 50, 100] }"
          @page-change="(p: number) => { queryParams.pageNum = p; loadPage() }"
          @page-size-change="(s: number) => { queryParams.pageSize = s; queryParams.pageNum = 1; loadPage() }">
          <template #columns>
            <a-table-column :title="t('tool.job.jobName')" data-index="jobName" :width="160" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('tool.job.jobGroup')" data-index="jobGroup" :width="120" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('tool.job.jobKey')" data-index="jobKey" :width="160" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('tool.job.cronExpr')" data-index="cronExpr" :width="160" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('tool.job.invokeTarget')" data-index="invokeTarget" :width="160" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('tool.job.colConcurrent')" data-index="concurrent" align="center" :width="80">
              <template #cell="{ record }">
                <a-tag :color="record.concurrent === 1 ? 'green' : 'gray'" size="small">
                  {{ record.concurrent === 1 ? t('tool.job.concurrentAllow') : t('tool.job.concurrentForbid') }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column :title="t('tool.job.colStatus')" data-index="status" align="center" :width="90">
              <template #cell="{ record }">
                <a-switch :model-value="record.status === 0" size="small"
                  @change="(val: string | number | boolean) => handleStatusChange(record, val === true || val === 0 || val === '0')" />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.action')" :width="220" align="center" fixed="right">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button v-permission="'job$trigger'" size="mini" type="text"
                    :loading="triggerLoading[record.id]" @click.stop="handleTrigger(record)">
                    <template #icon><IconPlayArrow /></template>{{ t('tool.job.trigger') }}
                  </a-button>
                  <a-button v-permission="'job$edit'" size="mini" type="text" @click.stop="handleEdit(record)">
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button v-permission="'job:log$list'" size="mini" type="text" @click.stop="handleViewLog(record)">
                    <template #icon><IconRefresh /></template>{{ t('tool.job.log') }}
                  </a-button>
                  <a-button v-permission="'job$delete'" size="mini" type="text" status="danger" @click.stop="handleDelete(record)">
                    <template #icon><IconDelete /></template>{{ t('common.delete') }}
                  </a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>
    </a-card>

    <a-modal v-model:visible="modalVisible" :title="isEdit ? t('tool.job.editTitle') : t('tool.job.addTitle')"
      :width="isMobile ? '100%' : 640" :fullscreen="isMobile" :footer="false">
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.job.jobName')" field="jobName"
              :rules="[{ required: true, message: t('tool.job.jobNameRequired') }, { maxLength: 64, message: t('common.maxLength', { max: 64 }) }]"
              validate-trigger="blur">
              <a-input v-model="form.jobName" :placeholder="t('tool.job.jobNamePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.job.jobGroup')" field="jobGroup">
              <a-input v-model="form.jobGroup" :placeholder="t('tool.job.jobGroupPlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.job.jobKey')" field="jobKey"
              :rules="[{ required: true, message: t('tool.job.jobKeyRequired') }, { maxLength: 64, message: t('common.maxLength', { max: 64 }) }]"
              validate-trigger="blur">
              <a-input v-model="form.jobKey" :placeholder="t('tool.job.jobKeyPlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.job.cronExpr')" field="cronExpr"
              :rules="[{ required: true, message: t('tool.job.cronRequired') }]"
              validate-trigger="blur">
              <CronPicker v-model="form.cronExpr" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('tool.job.invokeTarget')" field="invokeTarget">
              <a-input v-model="form.invokeTarget" :placeholder="t('tool.job.invokeTargetPlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 8">
            <a-form-item :label="t('tool.job.timeout')" field="timeout">
              <a-input-number v-model="form.timeout" :min="0" :max="3600" style="width:100%" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 8">
            <a-form-item :label="t('tool.job.retryCount')" field="retryCount">
              <a-input-number v-model="form.retryCount" :min="0" :max="10" style="width:100%" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 8">
            <a-form-item :label="t('tool.job.concurrent')" field="concurrent">
              <a-radio-group v-model="form.concurrent" type="button">
                <a-radio :value="0">{{ t('tool.job.concurrentForbid') }}</a-radio>
                <a-radio :value="1">{{ t('tool.job.concurrentAllow') }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.job.status')" field="status">
              <a-radio-group v-model="form.status" type="button">
                <a-radio :value="0">{{ t('tool.job.statusNormal') }}</a-radio>
                <a-radio :value="1">{{ t('tool.job.statusPaused') }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('common.remark')" field="remark">
              <a-textarea v-model="form.remark" :placeholder="t('tool.job.remarkPlaceholder')" :max-length="255" show-word-limit :auto-size="{ minRows: 2, maxRows: 4 }" />
            </a-form-item>
          </a-col>
          <a-col :span="24" style="margin-top:8px">
            <div style="display:flex;justify-content:center;gap:12px">
              <a-button @click="modalVisible = false">{{ t('common.cancel') }}</a-button>
              <a-button type="primary" :loading="submitLoading" @click="submitForm">{{ t('common.save') }}</a-button>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="logDrawerVisible" :title="t('tool.job.logTitle', { name: logJobName })"
      :width="isMobile ? '100%' : 800" :footer="false">
      <div style="margin-bottom:12px;display:flex;justify-content:space-between;align-items:center">
        <a-space>
          <a-button status="danger" v-permission="'job:log$delete'" size="small" :disabled="!selectedLogIds.length" @click="batchDeleteLog">
            <template #icon><IconDelete /></template>{{ t('common.batchDelete') }}
          </a-button>
          <a-button size="small" @click="loadLogPage">
            <template #icon><IconRefresh /></template>{{ t('common.refresh') }}
          </a-button>
        </a-space>
      </div>
      <a-table :data="logList" row-key="id" size="small"
        :row-selection="{ type: 'checkbox', showCheckedAll: true }"
        v-model:selectedKeys="selectedLogIds"
        :loading="logLoading"
        :pagination="{ current: logQuery.pageNum, pageSize: logQuery.pageSize, total: logTotal, showTotal: true, showPageSize: true }"
        @page-change="(p: number) => { logQuery.pageNum = p; loadLogPage() }"
        @page-size-change="(s: number) => { logQuery.pageSize = s; logQuery.pageNum = 1; loadLogPage() }">
        <template #columns>
          <a-table-column :title="t('tool.job.logStatus')" data-index="status" align="center" :width="80">
            <template #cell="{ record }">
              <a-tag :color="logStatusColor(record.status)" size="small">{{ logStatusText(record.status) }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="t('tool.job.logTriggerType')" data-index="triggerType" align="center" :width="90">
            <template #cell="{ record }">
              <a-tag :color="record.triggerType === 'manual' ? 'arcoblue' : 'gray'" size="small">
                {{ record.triggerType === 'manual' ? t('tool.job.logTriggerManual') : t('tool.job.logTriggerAuto') }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="t('tool.job.logDuration')" data-index="duration" align="center" :width="90" />
          <a-table-column :title="t('tool.job.logMessage')" data-index="message" :ellipsis="true">
            <template #cell="{ record }">
              <a-tooltip v-if="record.message" :content="record.message" position="top">
                <span :style="record.status === 1 ? 'color:var(--color-danger-6);cursor:pointer' : ''">
                  {{ record.message }}
                </span>
              </a-tooltip>
              <span v-else style="color:var(--color-text-4)">—</span>
            </template>
          </a-table-column>
          <a-table-column :title="t('tool.job.logStartTime')" data-index="startTime" :width="230" align="center" />
        </template>
      </a-table>
    </a-drawer>
  </div>
</template>
