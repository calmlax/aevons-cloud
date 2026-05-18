<script setup lang="ts">
import { IconDelete, IconFilter, IconCheck, IconClose, IconEye } from '@arco-design/web-vue/es/icon'
import operLogApi from '@/api/monitor/oper_log'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const proxy = getCurrentInstance()!.proxy as any
const dicts = proxy.$useDict('sys_success_status')

const isMobile = ref(window.innerWidth < 576)
function onResize() { isMobile.value = window.innerWidth < 576 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

// ── 分页 ──────────────────────────────────────────────
const total = ref(0)
const loading = ref(false)
const dataList = ref<any[]>([])
const queryParams = ref({ pageNum: 1, pageSize: 20, username: '', module: '', status: '',direction: 'desc',field:'operAt' })

async function loadPage() {
  loading.value = true
  try {
    const res: any = await operLogApi.page(queryParams.value)
    dataList.value = res?.rows ?? []
    total.value = res?.total ?? 0
  } finally {
    loading.value = false
  }
}

function handleSearch() { queryParams.value.pageNum = 1; loadPage() }
function handleReset() { queryParams.value.username = ''; queryParams.value.module = ''; queryParams.value.status = ''; handleSearch() }

const advancedVisible = ref(false)
const filterModule = computed(() => queryParams.value.username !== '' || queryParams.value.module !== '' || queryParams.value.status !== '')

// ── 多选 ──────────────────────────────────────────────
const selectedIds = ref<string[]>([])
const rowSelection = reactive({ type: 'checkbox', showCheckedAll: true, onlyCurrent: false })

// ── 移动端 ────────────────────────────────────────────
const mobilePageSize = 10
const mobilePage = ref(mobilePageSize)
const sentinelRef = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null
const mobileData = computed(() => dataList.value.slice(0, mobilePage.value))
const mobileHasMore = computed(() => mobilePage.value < dataList.value.length)
const mobileSelectMode = ref(false)

function toggleMobileSelectMode() { mobileSelectMode.value = !mobileSelectMode.value; if (!mobileSelectMode.value) selectedIds.value = [] }
function toggleCardSelect(id: string) { const i = selectedIds.value.indexOf(id); i === -1 ? selectedIds.value.push(id) : selectedIds.value.splice(i, 1) }
function isCardSelected(id: string) { return selectedIds.value.includes(id) }
const mobileAllSelected = computed(() => mobileData.value.length > 0 && mobileData.value.every((i: any) => selectedIds.value.includes(i.id)))
const mobileIndeterminate = computed(() => selectedIds.value.length > 0 && !mobileAllSelected.value)
function toggleMobileSelectAll() {
  const ids = mobileData.value.map((i: any) => i.id)
  if (mobileAllSelected.value) { selectedIds.value = selectedIds.value.filter(id => !ids.includes(id)) }
  else { selectedIds.value.push(...ids.filter(id => !selectedIds.value.includes(id))) }
}

function setupObserver() {
  if (!sentinelRef.value) return
  observer?.disconnect()
  observer = new IntersectionObserver(entries => { if (entries[0].isIntersecting && mobileHasMore.value) mobilePage.value += mobilePageSize }, { rootMargin: '80px' })
  observer.observe(sentinelRef.value)
}
onMounted(() => setupObserver())
onUnmounted(() => observer?.disconnect())
watch(dataList, async () => { mobilePage.value = mobilePageSize; await nextTick(); setupObserver() })

// ── 详情弹窗 ──────────────────────────────────────────
const detailVisible = ref(false)
const detailRow = ref<any>(null)
function handleDetail(row: any) { detailRow.value = row; detailVisible.value = true }

// ── 删除 ──────────────────────────────────────────────
async function batchDelete() {
  if (!selectedIds.value.length) { Message.warning(t('common.confirmDelete')); return }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await operLogApi.delete(selectedIds.value)
      Message.success(t('common.deleteSuccess'))
      selectedIds.value = []
      loadPage()
    },
  })
}

async function handleClear() {
  Modal.confirm({
    title: t('common.confirmClear'),
    content: t('common.confirmClearContent'),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await operLogApi.clear()
      Message.success(t('common.clearSuccess'))
      loadPage()
    },
  })
}

// 操作类型标签色
function typeColor(tp: string) {
  const m: Record<string, string> = { INSERT: 'green', UPDATE: 'arcoblue', DELETE: 'red', QUERY: 'gray', CLEAN: 'orange' }
  return m[tp] || 'gray'
}

loadPage()
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">
      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'monitor:oper:log$delete'" status="danger" :disabled="!selectedIds.length" @click="batchDelete">
            <template #icon><IconDelete /></template>{{ t('common.batchDelete') }}
          </a-button>
          <a-button v-permission="'monitor:oper:log$clear'" status="danger" @click="handleClear">{{ t('common.clear') }}</a-button>
        </a-space>
        <div class="cl-toolbar-right">
          <a-input-search v-model="queryParams.username" :placeholder="t('monitor.operLog.searchPlaceholder')" allow-clear class="cl-toolbar-search" @search="handleSearch" @press-enter="handleSearch" />
          <a-popover v-model:popup-visible="advancedVisible" trigger="click" position="br" popup-container="body">
            <a-button shape="circle" :type="filterModule ? 'primary' : 'secondary'">
              <template #icon><IconFilter /></template>
            </a-button>
            <template #content>
              <div class="cl-filter-panel">
                <p class="cl-filter-title">{{ t('common.advancedFilter') }}</p>
                <a-input v-model="queryParams.module" :placeholder="t('monitor.operLog.modulePlaceholder')" allow-clear />
                <a-select v-model="queryParams.status" :placeholder="t('common.status')" allow-clear>
                  <a-option v-for="dict in dicts.sys_success_status" :key="dict.value" :value="Number(dict.value)">{{ dict.label }}</a-option>
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
          :pagination="{ current: queryParams.pageNum, pageSize: queryParams.pageSize, total, showTotal: true, showPageSize: true, showJumper: true }"
          @page-change="(p: number) => { queryParams.pageNum = p; loadPage() }"
          @page-size-change="(s: number) => { queryParams.pageSize = s; queryParams.pageNum = 1; loadPage() }">
          <template #columns>
            <a-table-column :title="t('monitor.operLog.module')" data-index="module" :width="150" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('monitor.operLog.type')" data-index="type" :width="90" align="center">
              <template #cell="{ record }">
                <a-tag :color="typeColor(record.type)" size="small">{{ record.type }}</a-tag>
              </template>
            </a-table-column>
            <a-table-column :title="t('monitor.operLog.username')" data-index="username" :width="100" />
            <a-table-column :title="t('monitor.operLog.ip')" data-index="ip" :width="120" />
            <a-table-column :title="t('monitor.operLog.method')" data-index="method" :width="100" align="center" />
            <a-table-column :title="t('monitor.operLog.url')" data-index="url" :ellipsis="true" :width="200" :tooltip="true" />
            <a-table-column :title="t('monitor.operLog.time')" data-index="time" :width="90" align="center">
              <template #cell="{ record }">
                <a-tag :color="Number(record.time) > 1000 ? 'red' : Number(record.time) > 500 ? 'orange' : 'green'" size="small">
                  {{ record.time }}ms
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column :title="t('common.status')" data-index="status" :width="90" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_success_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('monitor.operLog.operAt')" data-index="operAt" :width="230" align="center" />
            <a-table-column :title="t('common.action')" :width="80" align="center">
              <template #cell="{ record }">
                <a-button size="mini" type="text" @click="handleDetail(record)">
                  <template #icon><IconEye /></template>{{ t('common.detail') }}
                </a-button>
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
            @click="mobileSelectMode ? toggleCardSelect(item.id) : handleDetail(item)">
            <div class="cl-card-header">
              <a-checkbox v-if="mobileSelectMode" :model-value="isCardSelected(item.id)" @click.stop @change="toggleCardSelect(item.id)" class="cl-card-checkbox" />
              <div class="cl-card-identity">
                <strong>{{ item.module }}</strong>
                <span class="cl-card-sub">{{ item.username }} · {{ item.ip }}</span>
              </div>
              <a-tag :color="typeColor(item.type)" size="small">{{ item.type }}</a-tag>
            </div>
            <div class="cl-card-meta">
              <span>{{ item.method }} {{ item.url }}</span>
              <span>{{ item.time }}ms · {{ item.operAt }}</span>
            </div>
          </div>
          <a-empty v-if="!mobileData.length" :description="t('common.noData')" />
          <div ref="sentinelRef" class="cl-sentinel">
            <span v-if="!mobileHasMore && mobileData.length > 0" class="cl-no-more">{{ t('common.noMore') }}</span>
          </div>
        </template>
      </div>
    </a-card>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" :title="t('monitor.operLog.detail')" :width="isMobile ? '100%' : 640" :fullscreen="isMobile" :footer="false">
      <a-descriptions v-if="detailRow" :column="2" bordered size="small">
        <a-descriptions-item :label="t('monitor.operLog.module')">{{ detailRow.module }}</a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.type')"><a-tag :color="typeColor(detailRow.type)" size="small">{{ detailRow.type }}</a-tag></a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.username')">{{ detailRow.username }}</a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.ip')">{{ detailRow.ip }}</a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.method')">{{ detailRow.method }}</a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.time')">{{ detailRow.time }}ms</a-descriptions-item>
        <a-descriptions-item label="URL" :span="2">{{ detailRow.url }}</a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.device')">{{ detailRow.device }}</a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.os')">{{ detailRow.os }}</a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.browser')">{{ detailRow.browser }}</a-descriptions-item>
        <a-descriptions-item :label="t('common.status')"><DictTag :options="dicts.sys_success_status" :value="detailRow.status" dot /></a-descriptions-item>
        <a-descriptions-item v-if="detailRow.error" :label="t('monitor.operLog.error')" :span="2">
          <a-typography-text type="danger">{{ detailRow.error }}</a-typography-text>
        </a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.payload')" :span="2">
          <a-typography-text code style="word-break:break-all;white-space:pre-wrap">{{ detailRow.payload || '-' }}</a-typography-text>
        </a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.result')" :span="2">
          <a-typography-text code style="word-break:break-all;white-space:pre-wrap">{{ detailRow.result || '-' }}</a-typography-text>
        </a-descriptions-item>
        <a-descriptions-item :label="t('monitor.operLog.operAt')" :span="2">{{ detailRow.operAt }}</a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>
