<script setup lang="ts">
import { IconDelete, IconDownload, IconEdit, IconFilter, IconPlus, IconCheck, IconClose, IconLock, IconUpload } from '@arco-design/web-vue/es/icon'
import userApi from '@/api/system/user'
import roleApi from '@/api/system/role'
import deptApi from '@/api/system/dept'
import postApi from '@/api/system/post'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const proxy = getCurrentInstance()!.proxy as any
const isMobile = ref(window.innerWidth < 576)
function onResize() { isMobile.value = window.innerWidth < 576 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

const dicts = proxy.$useDict('sys_user_type', 'sys_user_sex', 'sys_common_status')

// ── 搜索 ────────────────────────────────────────────────
const filterModule = ref(false)
const advancedVisible = ref(false)
const total = ref(0)
const loading = ref(false)
const dataList = ref<any[]>([])
const queryParams = ref<{
  pageNum: number; pageSize: number
  username?: string; nickname?: string; type?: number
  email?: string; mobile?: string; sex?: number; status?: number
}>({ pageNum: 1, pageSize: 20 })

async function loadPage() {
  loading.value = true
  try {
    const res: any = await userApi.page(queryParams.value)
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
  filterModule.value = !!(q.username || q.nickname || q.type != null || q.email || q.mobile || q.sex != null || q.status != null)
}, { deep: true })

// ── 移动端无限滚动 ─────────────────────────────────────
const mobilePageSize = 10
const mobilePage = ref(mobilePageSize)
const mobileLoadingMore = ref(false)
const sentinelRef = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null
const mobileData = computed(() => dataList.value.slice(0, mobilePage.value))
const mobileHasMore = computed(() => mobilePage.value < dataList.value.length)
function loadMore() {
  if (!mobileHasMore.value || mobileLoadingMore.value) return
  mobileLoadingMore.value = true
  setTimeout(() => { mobilePage.value += mobilePageSize; mobileLoadingMore.value = false }, 300)
}
function setupObserver() {
  if (!sentinelRef.value) return
  observer?.disconnect()
  observer = new IntersectionObserver(entries => { if (entries[0].isIntersecting) loadMore() }, { rootMargin: '80px' })
  observer.observe(sentinelRef.value)
}
onMounted(() => setupObserver())
onUnmounted(() => observer?.disconnect())
watch(dataList, async () => { mobilePage.value = mobilePageSize; await nextTick(); setupObserver() })

// ── 多选 ───────────────────────────────────────────────
const selectedIds = ref<string[]>([])
const rowSelection = reactive({ type: 'checkbox', showCheckedAll: true, onlyCurrent: false })
const mobileSelectMode = ref(false)
function toggleMobileSelectMode() { mobileSelectMode.value = !mobileSelectMode.value; if (!mobileSelectMode.value) selectedIds.value = [] }
function toggleCardSelect(id: string) { const idx = selectedIds.value.indexOf(id); idx === -1 ? selectedIds.value.push(id) : selectedIds.value.splice(idx, 1) }
function isCardSelected(id: string) { return selectedIds.value.includes(id) }
const mobileAllSelected = computed(() => mobileData.value.length > 0 && mobileData.value.every((item: any) => selectedIds.value.includes(item.id)))
const mobileIndeterminate = computed(() => selectedIds.value.length > 0 && !mobileAllSelected.value)
function toggleMobileSelectAll() {
  const currentIds = mobileData.value.map((item: any) => item.id)
  if (mobileAllSelected.value) selectedIds.value = selectedIds.value.filter(id => !currentIds.includes(id))
  else selectedIds.value.push(...currentIds.filter((id: string) => !selectedIds.value.includes(id)))
}

async function batchDelete() {
  if (!selectedIds.value.length) { Message.warning(t('common.confirmDeleteContent', { count: 0 }).replace('0 ', '')); return }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => { await userApi.delete(selectedIds.value); Message.success(t('common.deleteSuccess')); selectedIds.value = []; loadPage() },
  })
}

// ── Excel 导入导出 ────────────────────────────────────
const exportLoading = ref(false)
const templateLoading = ref(false)
const importLoading = ref(false)
const fileInputRef = ref<HTMLInputElement | null>(null)
const importModalVisible = ref(false)
const importResultVisible = ref(false)
const importResult = ref<any>({
  successCount: 0,
  failedCount: 0,
  errors: [],
})

async function handleExport() {
  exportLoading.value = true
  try {
    await userApi.exportFile(queryParams.value, t('system.user.exportFileName'))
    Message.success(t('system.user.exportSuccess'))
  } finally {
    exportLoading.value = false
  }
}

async function handleDownloadTemplate() {
  templateLoading.value = true
  try {
    await userApi.downloadTemplate(t('system.user.importTemplateFileName'))
    Message.success(t('system.user.downloadTemplateSuccess'))
  } finally {
    templateLoading.value = false
  }
}

function openImportPicker() {
  fileInputRef.value?.click()
}

async function handleImportChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return

  const fileName = file.name.toLowerCase()
  if (!fileName.endsWith('.xlsx')) {
    Message.error(t('system.user.importFileTypeError'))
    return
  }

  const formData = new FormData()
  formData.append('file', file)

  importLoading.value = true
  try {
    const res: any = await userApi.importExcel(formData)
    importResult.value = {
      successCount: res?.successCount ?? 0,
      failedCount: res?.failedCount ?? 0,
      errors: res?.errors ?? [],
    }
    importModalVisible.value = false
    importResultVisible.value = true
    Message.success(t('system.user.importParsed'))
  } finally {
    importLoading.value = false
  }
}

// ── 选项数据 ───────────────────────────────────────────
const roleOptions = ref<any[]>([])
const deptOptions = ref<any[]>([])
const postOptions = ref<any[]>([])

async function loadOptions() {
  const [roles, depts, posts] = await Promise.all([
    roleApi.list({ status: 0 }),
    deptApi.list({ status: 0 }),
    postApi.list({ status: 0 }),
  ])
  roleOptions.value = (roles as any) ?? []
  deptOptions.value = (depts as any) ?? []
  postOptions.value = (posts as any) ?? []
}

// ── 新增 / 修改 弹窗 ────────────────────────────────────
const modalVisible = ref(false)
const formRef = ref()
const isEdit = ref(false)
const submitLoading = ref(false)
const form = reactive<any>({
  username: '', nickname: '', type: undefined, email: '', mobile: '',
  sex: undefined, avatar: '', autograph: '', password: '', status: 0,
  roleIds: [], deptPosts: [],
})

function handleAdd() {
  isEdit.value = false
  resetForm('')
  loadOptions()
  modalVisible.value = true
}

async function handleEdit(row: any) {
  isEdit.value = true
  resetForm(row.id)
  loadOptions()
  modalVisible.value = true
}

async function resetForm(id:string) {
  await nextTick()
  formRef.value?.clearValidate()
  if (isEdit.value) {
    const [user, relations]: any = await Promise.all([
      userApi.getById(id),
      userApi.getRelations(id),
    ])
    Object.assign(form, user, {
      password: '',
      roleIds: relations?.roleIds ?? [],
      deptPosts: (relations?.deptPosts ?? []).map((d: any) => ({ deptId: d.deptId, postId: d.postId })),
    })
  } else {
    Object.assign(form, {
      id: undefined, username: '', nickname: '', type: 0, email: '',
      mobile: '', sex: 0, avatar: '', autograph: '', password: '',
      status: 0, roleIds: [], deptPosts: [{ deptId: undefined, postId: undefined }],
    })
  }
}

async function submitForm() {
  const valid = await formRef.value?.validate()
  if (valid) return
  submitLoading.value = true
  try {
    // 清理提交数据：roleIds 转 number，deptPosts 过滤空行，只保留必要字段
    const payload: any = {
      username: form.username,
      nickname: form.nickname,
      type: form.type,
      email: form.email,
      mobile: form.mobile,
      sex: form.sex,
      avatar: form.avatar,
      autograph: form.autograph,
      password: form.password,
      status: form.status,
      roleIds: form.roleIds ?? [],
      deptPosts: (form.deptPosts ?? [])
        .filter((dp: any) => dp.deptId != null && dp.postId != null)
        .map((dp: any) => ({ deptId: dp.deptId, postId: dp.postId })),
    }
    if (isEdit.value) {
      await userApi.update(form.id, payload)
      Message.success(t('common.editSuccess'))
    } else {
      await userApi.add(payload)
      Message.success(t('common.addSuccess'))
    }
    modalVisible.value = false
    loadPage()
  } finally {
    submitLoading.value = false
  }
}

// ── 部门岗位动态行 ─────────────────────────────────────
function addDeptPost() { 
    form.deptPosts.push({ deptId: undefined, postId: undefined }) 
}
function removeDeptPost(idx: number) { form.deptPosts.splice(idx, 1) }

// ── 状态开关 ───────────────────────────────────────────
async function handleStatusChange(record: any, val: boolean) {
  const status = val ? 0 : 1
  try {
    await userApi.updateStatus(record.id, status)
    record.status = status
    Message.success(t('system.user.statusUpdateSuccess'))
  } catch {
    record.status = val ? 1 : 0
    Message.error(t('system.user.statusUpdateFailed'))
  }
}

// ── 重置密码弹窗 ───────────────────────────────────────
const resetPwdVisible = ref(false)
const resetPwdFormRef = ref()
const resetPwdForm = reactive({ id: '' as string | number, password: '', confirmPassword: '' })
const resetPwdLoading = ref(false)

function handleResetPassword(row: any) {
  resetPwdForm.id = row.id
  resetPwdForm.password = ''
  resetPwdForm.confirmPassword = ''
  resetPwdVisible.value = true
}

async function submitResetPassword() {
  const valid = await resetPwdFormRef.value?.validate()
  if (valid) return
  resetPwdLoading.value = true
  try {
    await userApi.resetPassword(resetPwdForm.id, resetPwdForm.password)
    Message.success(t('system.user.resetSuccess'))
    resetPwdVisible.value = false
  } finally {
    resetPwdLoading.value = false
  }
}

// ── 删除 ───────────────────────────────────────────────
async function handleDelete(row: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.user.confirmDelete', { name: row.nickname || row.username }),
    okButtonProps: { status: 'danger' },
    onOk: async () => { await userApi.delete([row.id]); Message.success(t('common.deleteSuccess')); loadPage() },
  })
}

function confirmPasswordValidator(val: string, cb: (err?: string) => void) {
  if (val !== resetPwdForm.password) cb(t('system.user.confirmPasswordMismatch'))
  else cb()
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
          <a-button v-permission="'sys:user$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>{{ t('common.add') }}
          </a-button>
          <a-button v-permission="'sys:user$delete'" status="danger" :disabled="!selectedIds.length" @click="batchDelete">
            <template #icon><IconDelete /></template>{{ t('common.batchDelete') }}
          </a-button>
          <a-button v-permission="'sys:user$export'" :loading="exportLoading" @click="handleExport">
            <template #icon><IconDownload /></template>{{ t('common.export') }}
          </a-button>
          <a-button v-permission="'sys:user$import'" :loading="templateLoading || importLoading" @click="importModalVisible = true">
            <template #icon><IconUpload /></template>{{ t('system.user.importBtn') }}
          </a-button>
        </a-space>
        <input ref="fileInputRef" type="file" accept=".xlsx" style="display:none" @change="handleImportChange" />
        <div class="cl-toolbar-right">
          <a-input-search v-model="queryParams.username" :placeholder="t('system.user.searchPlaceholder')" allow-clear class="cl-toolbar-search"
            @search="handleSearch" @press-enter="handleSearch" />
          <a-popover v-model:popup-visible="advancedVisible" trigger="click" position="br" popup-container="body">
            <a-button shape="circle" :type="filterModule ? 'primary' : 'secondary'">
              <template #icon><IconFilter /></template>
            </a-button>
            <template #content>
              <div class="cl-filter-panel">
                <p class="cl-filter-title">{{ t('common.advancedFilter') }}</p>
                <a-input v-model="queryParams.username" :placeholder="t('system.user.username')" allow-clear />
                <a-input v-model="queryParams.nickname" :placeholder="t('system.user.nickname')" allow-clear />
                <a-input v-model="queryParams.email" :placeholder="t('system.user.email')" allow-clear />
                <a-input v-model="queryParams.mobile" :placeholder="t('system.user.mobile')" allow-clear />
                <a-select v-model="queryParams.status" :placeholder="t('common.status')" allow-clear>
                  <a-option v-for="d in dicts.sys_common_status" :key="d.value" :value="Number(d.value)">{{ d.label }}</a-option>
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

      <!-- PC 表格 -->
      <div class="cl-table-wrap">
        <a-skeleton v-if="loading" :animation="true"><a-skeleton-line :rows="8" /></a-skeleton>
        <a-table v-else
          :bordered="false" :data="dataList" row-key="id" :scroll="{ x: 1000 }" :row-selection="rowSelection"
          v-model:selectedKeys="selectedIds"
          :pagination="{ current: queryParams.pageNum, pageSize: queryParams.pageSize, total, showTotal: true, showPageSize: true, showJumper: true, pageSizeOptions: [10, 20, 50, 100] }"
          @page-change="(p: number) => { queryParams.pageNum = p; loadPage() }"
          @page-size-change="(s: number) => { queryParams.pageSize = s; queryParams.pageNum = 1; loadPage() }">
          <template #columns>
            <a-table-column :title="t('system.user.username')" data-index="username" :width="150" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('system.user.nickname')" data-index="nickname" :width="150" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('system.user.email')" data-index="email" :width="150" :ellipsis="true" :tooltip="true"/>
            <a-table-column :title="t('system.user.mobile')" data-index="mobile" :width="150" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('system.user.sex')" data-index="sex" align="center"  :width="100">
              <template #cell="{ record }">
                <dict-tag :options="dicts.sys_user_sex" :value="record.sex" />
              </template>
            </a-table-column>
            <a-table-column :title="t('system.user.status')" data-index="status" align="center" :width="100">
              <template #cell="{ record }">
                <a-switch :model-value="record.status === 0" size="small"
                  @change="(val: string | number | boolean) => handleStatusChange(record, val === true || val === 0 || val === '0')" />
              </template>
            </a-table-column>
            <a-table-column :title="t('system.user.createdAt')" data-index="createdAt" align="center"  :width="240"/>
            <a-table-column :title="t('common.action')" :width="200" align="center" fixed="right">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button v-permission="'sys:user$edit'" size="mini" type="text" @click.stop="handleEdit(record)">
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button v-permission="'sys:user$edit'" size="mini" type="text" @click.stop="handleResetPassword(record)">
                    <template #icon><IconLock /></template>{{ t('system.user.resetPassword') }}
                  </a-button>
                  <a-button v-permission="'sys:user$delete'" size="mini" type="text" status="danger" @click.stop="handleDelete(record)">
                    <template #icon><IconDelete /></template>{{ t('common.delete') }}
                  </a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>

      <!-- 移动端卡片 -->
      <div class="cl-card-list">
        <div class="cl-mobile-select-bar">
          <a-button size="small" :type="mobileSelectMode ? 'primary' : 'secondary'" @click="toggleMobileSelectMode">
            <template #icon><IconClose v-if="mobileSelectMode" /><IconCheck v-else /></template>
            {{ mobileSelectMode ? t('common.cancel') : t('common.select') }}
          </a-button>
          <template v-if="mobileSelectMode">
            <a-checkbox :model-value="mobileAllSelected" :indeterminate="mobileIndeterminate" @change="toggleMobileSelectAll">{{ t('common.selectAll') }}</a-checkbox>
            <span class="cl-select-count">{{ t('common.selected', { count: selectedIds.length }) }}</span>
          </template>
        </div>
        <a-skeleton v-if="loading" :animation="true"><a-skeleton-line :rows="6" /></a-skeleton>
        <template v-else>
          <div v-for="(item, index) in mobileData" :key="item.id" class="cl-card stagger-item"
            :class="{ 'cl-card--selected': mobileSelectMode && isCardSelected(item.id) }"
            :style="{ '--stagger-index': index % mobilePageSize }"
            @click="mobileSelectMode ? toggleCardSelect(item.id) : undefined">
            <div class="cl-card-header">
              <a-checkbox v-if="mobileSelectMode" :model-value="isCardSelected(item.id)" @click.stop @change="toggleCardSelect(item.id)" class="cl-card-checkbox" />
              <div class="cl-card-identity">
                <strong>{{ item.nickname || item.username }}</strong>
                <span class="cl-card-sub">{{ item.username }}</span>
              </div>
              <a-tag :color="item.status === 0 ? 'green' : 'red'" size="small">{{ item.status === 0 ? t('common.normal') : t('common.disabled') }}</a-tag>
            </div>
            <div class="cl-card-meta">
              <span>{{ t('system.user.email') }}：{{ item.email || '-' }}</span>
              <span>{{ t('system.user.mobile') }}：{{ item.mobile || '-' }}</span>
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'sys:user$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button v-permission="'sys:user$edit'" size="mini" type="outline" @click="handleResetPassword(item)">
                  <template #icon><IconLock /></template>{{ t('system.user.resetPassword') }}
                </a-button>
                <a-button v-permission="'sys:user$delete'" size="mini" type="outline" status="danger" @click="handleDelete(item)">
                  <template #icon><IconDelete /></template>{{ t('common.delete') }}
                </a-button>
              </a-space>
            </div>
          </div>
          <a-empty v-if="!mobileData.length" :description="t('common.noData')" />
          <div ref="sentinelRef" class="cl-sentinel">
            <span v-if="!mobileHasMore && mobileData.length > 0" class="cl-no-more">{{ t('common.noMore') }}</span>
          </div>
        </template>
      </div>
    </a-card>

    <!-- 新增 / 修改 弹窗 -->
    <a-modal v-model:visible="modalVisible" :title="isEdit ? t('system.user.editTitle') : t('system.user.addTitle')"
      :width="isMobile ? '100%' : 680" :fullscreen="isMobile" :footer="false">
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item :label="t('system.user.userType')" field="type">
              <a-radio-group v-model="form.type" type="button">
                <a-radio v-for="dict in dicts.sys_user_type" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.user.username')" field="username"
              :rules="[{ required: true, message: t('system.user.usernameRequired') }, { maxLength: 32, message: t('common.maxLength', { max: 32 }) }]"
              validate-trigger="blur">
              <a-input v-model="form.username" :placeholder="t('system.user.usernamePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="isEdit ? t('system.user.passwordEdit') : t('system.user.password')" field="password"
              :rules="isEdit ? [] : [{ required: true, message: t('system.user.passwordRequired') }]" validate-trigger="blur">
              <a-input-password v-model="form.password" :placeholder="t('system.user.passwordPlaceholder')" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.user.nickname')" field="nickname"
              :rules="[{ required: true, message: t('system.user.nicknameRequired') }, { maxLength: 50, message: t('common.maxLength', { max: 50 }) }]"
              validate-trigger="blur">
              <a-input v-model="form.nickname" :placeholder="t('system.user.nicknamePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.user.sex')" field="sex">
              <a-radio-group v-model="form.sex" type="button">
                <a-radio v-for="dict in dicts.sys_user_sex" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.user.email')" field="email" :rules="[{ maxLength: 64, message: t('common.maxLength', { max: 64 }) }]" validate-trigger="blur">
              <a-input v-model="form.email" :placeholder="t('system.user.emailPlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.user.mobile')" field="mobile" :rules="[{ maxLength: 11, message: t('common.maxLength', { max: 11 }) }]" validate-trigger="blur">
              <a-input v-model="form.mobile" :placeholder="t('system.user.mobilePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.user.status')" field="status">
              <a-radio-group v-model="form.status" type="button">
                <a-radio v-for="dict in dicts.sys_common_status" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>

          <!-- 角色 -->
          <a-col :span="24">
            <a-form-item :label="t('system.user.roleAssign')" field="roleIds">
              <a-select v-model="form.roleIds" multiple :placeholder="t('system.user.rolePlaceholder')" allow-clear style="width:100%">
                <a-option v-for="r in roleOptions" :key="r.id" :value="(r.id as any)">{{ r.roleName }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>

          <!-- 部门岗位 -->
          <a-col :span="24">
            <a-form-item :label="t('system.user.deptPost')">
              <div style="display:flex;flex-direction:column;gap:8px;width:100%">
              <div v-for="(dp, idx) in form.deptPosts" :key="idx" style="display:flex;gap:8px;align-items:center;width:100%">
                <a-select v-model="dp.deptId" :placeholder="t('system.user.deptPlaceholder')" allow-clear style="flex:1">
                  <a-option v-for="d in deptOptions" :key="String(d.id)" :value="(d.id as any)">{{ d.deptName }}</a-option>
                </a-select>
                <a-select v-model="dp.postId" :placeholder="t('system.user.postPlaceholder')" allow-clear style="flex:1">
                  <a-option v-for="p in postOptions" :key="String(p.id)" :value="(p.id as any)">{{ p.postName }}</a-option>
                </a-select>
                <a-button size="mini" status="danger" @click="removeDeptPost(idx as number)">
                  <template #icon><IconClose /></template>
                </a-button>
              </div>
              <div style="display:flex;justify-content:center;margin-top:4px">
                <a-button size="small" type="text" @click="addDeptPost">
                  <template #icon><IconPlus /></template>{{ t('system.user.addDeptPost') }}
                </a-button>
              </div>
              </div>
            </a-form-item>
          </a-col>

          <a-col :span="24" style="margin-top:16px">
            <div style="display:flex;justify-content:center;gap:12px">
              <a-button @click="modalVisible = false">{{ t('common.cancel') }}</a-button>
              <a-button type="primary" :loading="submitLoading" @click="submitForm">{{ t('common.save') }}</a-button>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- 重置密码弹窗 -->
    <a-modal v-model:visible="resetPwdVisible" :title="t('system.user.resetPasswordTitle')" :width="isMobile ? '100%' : 400" :footer="false">
      <a-form :model="resetPwdForm" ref="resetPwdFormRef" layout="vertical">
        <a-form-item :label="t('system.user.newPassword')" field="password"
          :rules="[{ required: true, message: t('system.user.newPasswordRequired') }, { minLength: 6, message: t('system.user.newPasswordMinLength') }]"
          validate-trigger="blur">
          <a-input-password v-model="resetPwdForm.password" :placeholder="t('system.user.newPasswordPlaceholder')" allow-clear />
        </a-form-item>
        <a-form-item :label="t('system.user.confirmPassword')" field="confirmPassword"
          :rules="[{ required: true, message: t('system.user.confirmPasswordRequired') }, { validator: confirmPasswordValidator }]"
          validate-trigger="blur">
          <a-input-password v-model="resetPwdForm.confirmPassword" :placeholder="t('system.user.confirmPasswordPlaceholder')" allow-clear />
        </a-form-item>
        <div style="display:flex;justify-content:center;gap:12px;margin-top:16px">
          <a-button @click="resetPwdVisible = false">{{ t('common.cancel') }}</a-button>
          <a-button type="primary" :loading="resetPwdLoading" @click="submitResetPassword">{{ t('system.user.confirmReset') }}</a-button>
        </div>
      </a-form>
    </a-modal>

    <!-- 导入弹窗 -->
    <a-modal
      v-model:visible="importModalVisible"
      :title="t('system.user.importBtn')"
      :width="isMobile ? '100%' : 560"
      :footer="false">
      <a-space direction="vertical" fill size="large">
        <a-alert type="info">
          <template #title>{{ t('system.user.importTemplateFileName') }}</template>
          <template #default>{{ t('system.user.downloadTemplate') }} / {{ t('system.user.importBtn') }}</template>
        </a-alert>
        <a-button long :loading="templateLoading" @click="handleDownloadTemplate">
          <template #icon><IconDownload /></template>{{ t('system.user.downloadTemplate') }}
        </a-button>
        <a-button long type="primary" :loading="importLoading" @click="openImportPicker">
          <template #icon><IconUpload /></template>{{ t('system.user.importBtn') }}
        </a-button>
        <div style="display:flex;justify-content:center">
          <a-button @click="importModalVisible = false">{{ t('common.cancel') }}</a-button>
        </div>
      </a-space>
    </a-modal>

    <!-- 导入结果 -->
    <a-modal v-model:visible="importResultVisible" :title="t('system.user.importResultTitle')" :width="isMobile ? '100%' : 720" :footer="false">
      <a-space direction="vertical" fill size="large">
        <a-alert type="info">
          <template #title>{{ t('system.user.importResultSummary', { success: importResult.successCount, failed: importResult.failedCount }) }}</template>
        </a-alert>
        <a-table
          v-if="importResult.errors?.length"
          :bordered="false"
          :pagination="false"
          :data="importResult.errors"
          row-key="row"
          :scroll="{ y: 320 }">
          <template #columns>
            <a-table-column :title="t('system.user.importErrorRow')" data-index="row" :width="100" />
            <a-table-column :title="t('system.user.importErrorColumn')" data-index="column" :width="180" />
            <a-table-column :title="t('system.user.importErrorMsg')" data-index="msg" />
          </template>
        </a-table>
        <a-empty v-else :description="t('system.user.importNoErrors')" />
        <div style="display:flex;justify-content:center">
          <a-button type="primary" @click="importResultVisible = false">{{ t('common.confirm') }}</a-button>
        </div>
      </a-space>
    </a-modal>
  </div>
</template>
