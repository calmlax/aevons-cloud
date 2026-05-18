<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus } from '@arco-design/web-vue/es/icon'
import langApi from '@/api/system/lang'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const proxy = getCurrentInstance()!.proxy as any
const isMobile = ref(window.innerWidth < 576)
function onResize() { isMobile.value = window.innerWidth < 576 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

const dicts = proxy.$useDict('sys_is', 'sys_common_status')

const total = ref(0)
const loading = ref(false)
const dataList = ref<any[]>([])
const selectedIds = ref<string[]>([])
const rowSelection = reactive({ type: 'checkbox', showCheckedAll: true, onlyCurrent: false })
const advancedVisible = ref(false)
const filterActive = ref(false)
const queryParams = ref<any>({
  pageNum: 1, pageSize: 20,
  direction: 'ascend',field:'sort'
})

async function loadPage() {
  loading.value = true
  try {
    const res: any = await langApi.page(queryParams.value)
    dataList.value = res?.rows ?? []
    total.value = res?.total ?? 0
  } finally {
    loading.value = false
  }
}
function handleSearch() { queryParams.value.pageNum = 1; loadPage() }
function handleReset() { queryParams.value = { pageNum: 1, pageSize: queryParams.value.pageSize }; handleSearch() }
watch(() => ({ ...queryParams.value }), () => {
  filterActive.value = !!(queryParams.value.langName || queryParams.value.status != null)
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

async function batchDelete() {
  if (!selectedIds.value.length) { Message.warning(t('common.confirmDeleteContent', { count: 0 }).replace('0 ', '')); return }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => { await langApi.delete(selectedIds.value); Message.success(t('common.deleteSuccess')); selectedIds.value = []; loadPage() },
  })
}

// ── 新增 / 修改 弹窗 ──────────────────────────────────
const modalVisible = ref(false)
const formRef = ref()
const langForm = reactive<any>({})
const isEdit = ref(false)
const submitLoading = ref(false)

function handleAdd() {
  isEdit.value = false
  Object.assign(langForm, { id: undefined, langCode: '', langName: '', isDefault: 0, sort: 0, status: 0, remark: '' })
  nextTick(() => formRef.value?.clearValidate())
  modalVisible.value = true
}

async function handleEdit(row: any) {
  isEdit.value = true
  const res: any = await langApi.getById(row.id)
  Object.assign(langForm, res)
  nextTick(() => formRef.value?.clearValidate())
  modalVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value?.validate()
  if (valid) return
  submitLoading.value = true
  try {
    if (isEdit.value) {
      await langApi.update(langForm.id, langForm)
      Message.success(t('common.editSuccess'))
    } else {
      await langApi.add(langForm)
      Message.success(t('common.addSuccess'))
    }
    modalVisible.value = false
    loadPage()
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.lang.confirmDelete', { name: row.langName }),
    okButtonProps: { status: 'danger' },
    onOk: async () => { await langApi.delete([row.id]); Message.success(t('common.deleteSuccess')); loadPage() },
  })
}

loadPage()
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'sys:lang$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>{{ t('common.add') }}
          </a-button>
          <a-button v-permission="'sys:lang$delete'" status="danger" :disabled="!selectedIds.length" @click="batchDelete">
            <template #icon><IconDelete /></template>{{ t('common.batchDelete') }}
          </a-button>
        </a-space>
        <div class="cl-toolbar-right">
          <a-input-search v-model="queryParams.langName" :placeholder="t('system.lang.searchPlaceholder')" allow-clear class="cl-toolbar-search"
            @search="handleSearch" @press-enter="handleSearch" />
          <a-popover v-model:popup-visible="advancedVisible" trigger="click" position="br" popup-container="body">
            <a-button shape="circle" :type="filterActive ? 'primary' : 'secondary'">
              <template #icon><IconFilter /></template>
            </a-button>
            <template #content>
              <div class="cl-filter-panel">
                <p class="cl-filter-title">{{ t('common.advancedFilter') }}</p>
                <a-input v-model="queryParams.langName" :placeholder="t('system.lang.langNameFilterPlaceholder')" allow-clear />
                <a-select v-model="queryParams.status" :placeholder="t('common.status')" allow-clear>
                  <a-option v-for="dict in dicts.sys_common_status" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-option>
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
        <a-skeleton v-if="loading" :animation="true"><a-skeleton-line :rows="6" /></a-skeleton>
        <a-table v-else
          :bordered="false" :data="dataList" row-key="id" :scroll="{ x: 700 }"
          :row-selection="rowSelection" v-model:selectedKeys="selectedIds"
          :pagination="{ current: queryParams.pageNum, pageSize: queryParams.pageSize, total, showTotal: true, showPageSize: true, pageSizeOptions: [10, 20, 50] }"
          @page-change="(p: number) => { queryParams.pageNum = p; loadPage() }"
          @page-size-change="(s: number) => { queryParams.pageSize = s; queryParams.pageNum = 1; loadPage() }">
          <template #columns>
            <a-table-column :title="t('system.lang.langCode')" data-index="langCode" align="center" :width="160" />
            <a-table-column :title="t('system.lang.langName')" data-index="langName" align="center" :width="160" />
            <a-table-column :title="t('system.lang.isDefault')" data-index="isDefault" align="center" :width="100">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_is" :value="record.isDefault" />
              </template>
            </a-table-column>
            <a-table-column :title="t('system.lang.sort')" data-index="sort" align="center" :width="80" />
            <a-table-column :title="t('common.status')" data-index="status" align="center" :width="90">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_common_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.remark')" data-index="remark" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('common.action')" :width="160" align="center" fixed="right">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button v-permission="'sys:lang$edit'" size="mini" type="text" @click.stop="handleEdit(record)">
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button v-permission="'sys:lang$delete'" size="mini" type="text" status="danger" @click.stop="handleDelete(record)">
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
        <a-skeleton v-if="loading" :animation="true"><a-skeleton-line :rows="6" /></a-skeleton>
        <template v-else>
          <div v-for="(item, index) in mobileData" :key="item.id" class="cl-card stagger-item"
            :style="{ '--stagger-index': index % mobilePageSize }">
            <div class="cl-card-header">
              <div class="cl-card-identity">
                <strong>{{ item.langName }}</strong>
                <span class="cl-card-sub">{{ item.langCode }}</span>
              </div>
              <DictTag :options="dicts.sys_common_status" :value="item.status" dot />
            </div>
            <div class="cl-card-meta">
              <span>{{ t('system.lang.isDefault') }}：<DictTag :options="dicts.sys_is" :value="item.isDefault" /></span>
              <span>{{ t('system.lang.sort') }}：{{ item.sort }}</span>
            </div>
            <div class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'sys:lang$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button v-permission="'sys:lang$delete'" size="mini" type="outline" status="danger" @click="handleDelete(item)">
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
    <a-modal v-model:visible="modalVisible" :title="isEdit ? t('system.lang.editTitle') : t('system.lang.addTitle')"
      :width="isMobile ? '100%' : 520" :fullscreen="isMobile" :footer="false">
      <a-form :model="langForm" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.lang.langCode')" field="langCode"
              :rules="[{ required: true, message: t('system.lang.langCodeRequired') }, { maxLength: 10, message: t('common.maxLength', { max: 10 }) }]"
              validate-trigger="blur">
              <a-input v-model="langForm.langCode" :placeholder="t('system.lang.langCodePlaceholder')" :disabled="isEdit" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.lang.langName')" field="langName"
              :rules="[{ required: true, message: t('system.lang.langNameRequired') }, { maxLength: 50, message: t('common.maxLength', { max: 50 }) }]"
              validate-trigger="blur">
              <a-input v-model="langForm.langName" :placeholder="t('system.lang.langNamePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.lang.sort')" field="sort">
              <a-input-number v-model="langForm.sort" :min="0" style="width:100%" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.lang.isDefault')" field="isDefault">
              <a-switch v-model="langForm.isDefault" :checked-value="1" :unchecked-value="0" :checked-text="t('system.lang.isDefaultYes')" :unchecked-text="t('system.lang.isDefaultNo')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.status')" field="status">
              <a-radio-group v-model="langForm.status" type="button">
                <a-radio v-for="dict in dicts.sys_common_status" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('common.remark')" field="remark" :rules="[{ maxLength: 200, message: t('common.maxLength', { max: 200 }) }]">
              <a-textarea v-model="langForm.remark" :placeholder="t('system.lang.remarkPlaceholder')" allow-clear />
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
  </div>
</template>
