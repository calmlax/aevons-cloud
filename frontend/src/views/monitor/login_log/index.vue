<script setup lang="ts">
import { IconDelete, IconFilter, IconCheck, IconClose } from '@arco-design/web-vue/es/icon'
import loginLogApi from '@/api/monitor/login_log'
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
const queryParams = ref({ pageNum: 1, pageSize: 20, username: '', status: '' as any,direction: 'desc',field:'loginAt'})

async function loadPage() {
  loading.value = true
  try {
    const res: any = await loginLogApi.page(queryParams.value)
    dataList.value = res?.rows ?? []
    total.value = res?.total ?? 0
  } finally {
    loading.value = false
  }
}

function handleSearch() { queryParams.value.pageNum = 1; loadPage() }
function handleReset() { queryParams.value.username = ''; queryParams.value.status = ''; handleSearch() }

// ── 筛选 ──────────────────────────────────────────────
const advancedVisible = ref(false)
const filterModule = computed(() => queryParams.value.username !== '' || queryParams.value.status !== '')

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

// ── 删除 ──────────────────────────────────────────────
    async function batchDelete() {
  if (!selectedIds.value.length) { Message.warning(t('common.confirmDelete')); return }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await loginLogApi.delete(selectedIds.value)
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
      await loginLogApi.clear()
      Message.success(t('common.clearSuccess'))
      loadPage()
    },
  })
}

loadPage()
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">
      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'monitor:login:log$delete'" status="danger" :disabled="!selectedIds.length" @click="batchDelete">
            <template #icon><IconDelete /></template>{{ t('common.batchDelete') }}
          </a-button>
          <a-button v-permission="'monitor:login:log$clear'" status="danger" @click="handleClear">{{ t('common.clear') }}</a-button>
        </a-space>
        <div class="cl-toolbar-right">
          <a-input-search v-model="queryParams.username" :placeholder="t('monitor.loginLog.searchPlaceholder')" allow-clear class="cl-toolbar-search" @search="handleSearch" @press-enter="handleSearch" />
          <a-popover v-model:popup-visible="advancedVisible" trigger="click" position="br" popup-container="body">
            <a-button shape="circle" :type="filterModule ? 'primary' : 'secondary'">
              <template #icon><IconFilter /></template>
            </a-button>
            <template #content>
              <div class="cl-filter-panel">
                <p class="cl-filter-title">{{ t('common.advancedFilter') }}</p>
                <a-select v-model="queryParams.status" placeholder="状态" allow-clear>
                  <a-option v-for="dict in dicts.sys_success_status" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-option>
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
          :bordered="false" :data="dataList" row-key="id" :scroll="{ x: 900 }" :row-selection="rowSelection"
          v-model:selectedKeys="selectedIds"
          :pagination="{ current: queryParams.pageNum, pageSize: queryParams.pageSize, total, showTotal: true, showPageSize: true, showJumper: true }"
          @page-change="(p: number) => { queryParams.pageNum = p; loadPage() }"
          @page-size-change="(s: number) => { queryParams.pageSize = s; queryParams.pageNum = 1; loadPage() }">
          <template #columns>
            <a-table-column :title="t('monitor.loginLog.username')" data-index="username" :width="200" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('monitor.loginLog.clientId')" data-index="clientId" :width="120" align="center" />
            <a-table-column :title="t('monitor.loginLog.grantType')" data-index="grantType" :width="180" align="center" />
            <a-table-column :title="t('monitor.loginLog.ip')" data-index="ip" :width="140" />
            <a-table-column :title="t('monitor.loginLog.os')" data-index="os" :width="100" />
            <a-table-column :title="t('monitor.loginLog.browser')" data-index="browser" :width="100" />
            <a-table-column :title="t('common.status')" data-index="status" :width="90" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_success_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('monitor.loginLog.msg')" data-index="msg" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('monitor.loginLog.loginAt')" data-index="loginAt" :width="240" align="center" />
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
                <strong>{{ item.username }}</strong>
                <span class="cl-card-sub">{{ item.ip }}</span>
              </div>
              <DictTag :options="dicts.sys_success_status" :value="item.status" dot />
            </div>
            <div class="cl-card-meta">
              <span>{{ item.os }} · {{ item.browser }}</span>
              <span>{{ item.loginAt }}</span>
            </div>
          </div>
          <a-empty v-if="!mobileData.length" :description="t('common.noData')" />
          <div ref="sentinelRef" class="cl-sentinel">
            <span v-if="!mobileHasMore && mobileData.length > 0" class="cl-no-more">{{ t('common.noMore') }}</span>
          </div>
        </template>
      </div>
    </a-card>
  </div>
</template>
