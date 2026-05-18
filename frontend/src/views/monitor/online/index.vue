<script setup lang="ts">
import { IconDelete, IconRefresh } from '@arco-design/web-vue/es/icon'
import { useI18n } from 'vue-i18n'
import onlineApi from '@/api/monitor/online'

const { t } = useI18n()

// ── 响应式断点 ────────────────────────────────────────
const isMobile = ref(window.innerWidth < 576)
function onResize() { isMobile.value = window.innerWidth < 576 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

// ── 列表 ──────────────────────────────────────────────
const loading = ref(false)
const dataList = ref<any[]>([])
const keyword = ref('')

const filteredList = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return dataList.value
  return dataList.value.filter(i =>
    i.username?.toLowerCase().includes(kw) ||
    i.nickname?.toLowerCase().includes(kw) ||
    i.clientId?.toLowerCase().includes(kw)
  )
})

async function loadList() {
  loading.value = true
  try {
    dataList.value = (await onlineApi.list() as any) ?? []
  } finally {
    loading.value = false
  }
}

// ── 强制下线 ──────────────────────────────────────────
async function handleForceLogout(row: any) {
  Modal.confirm({
    title: t('monitor.online.confirmForceLogout'),
    content: t('monitor.online.confirmForceLogoutContent', { name: row.nickname || row.username }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await onlineApi.forceLogout(row.token)
      Message.success(t('monitor.online.forceLogoutSuccess'))
      loadList()
    },
  })
}

// ── 移动端无限滚动 ─────────────────────────────────────
const mobilePageSize = 10
const mobilePage = ref(mobilePageSize)
const sentinelRef = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

const mobileData = computed(() => filteredList.value.slice(0, mobilePage.value))
const mobileHasMore = computed(() => mobilePage.value < filteredList.value.length)

function setupObserver() {
  if (!sentinelRef.value) return
  observer?.disconnect()
  observer = new IntersectionObserver(entries => {
    if (entries[0].isIntersecting && mobileHasMore.value) {
      mobilePage.value += mobilePageSize
    }
  }, { rootMargin: '80px' })
  observer.observe(sentinelRef.value)
}
onMounted(() => setupObserver())
onUnmounted(() => observer?.disconnect())
watch(dataList, async () => {
  mobilePage.value = mobilePageSize
  await nextTick()
  setupObserver()
})

loadList()
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-tag color="arcoblue" size="large">
            {{ t('monitor.online.count', { count: dataList.length }) }}
          </a-tag>
        </a-space>
        <div class="cl-toolbar-right">
          <a-input-search
            v-model="keyword"
          :placeholder="t('monitor.online.searchPlaceholder')"
            allow-clear
            class="cl-toolbar-search"
          />
          <a-button :loading="loading" @click="loadList">
            <template #icon><IconRefresh /></template>
          </a-button>
        </div>
      </div>

      <!-- PC 表格 -->
      <div class="cl-table-wrap">
        <a-skeleton v-if="loading" :animation="true">
          <a-skeleton-line :rows="8" />
        </a-skeleton>
        <a-table
          v-else
          :bordered="false"
          :data="filteredList"
          row-key="token"
          :scroll="{ x: 900 }"
          :pagination="{ pageSize: 20, showTotal: true, showPageSize: true }"
        >
          <template #columns>
            <a-table-column :title="t('monitor.online.username')" data-index="username" :width="120" />
            <a-table-column :title="t('monitor.online.nickname')" data-index="nickname" :width="120" />
            <a-table-column :title="t('monitor.online.clientId')" data-index="clientId" :width="130" />
            <a-table-column title="Access Token" data-index="token" :width="200" :ellipsis="true" :tooltip="true" />
            <a-table-column title="Refresh Token" data-index="refreshToken" :width="200" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('monitor.online.expireAt')" data-index="expireAt" :width="170" align="center" />
            <a-table-column :title="t('monitor.online.ttl')" data-index="ttl" :width="100" align="center">
              <template #cell="{ record }">
                <a-tag :color="record.ttl > 300 ? 'green' : 'orangered'" size="small">
                  {{ record.ttl }}
                </a-tag>
              </template>
            </a-table-column>
            <a-table-column :title="t('common.action')" :width="110" align="center">
              <template #cell="{ record }">
                <a-button
                  v-permission="'monitor:online$forceLogout'"
                  size="mini"
                  type="text"
                  status="danger"
                  @click.stop="handleForceLogout(record)"
                >
                  <template #icon><IconDelete /></template>{{ t('monitor.online.forceLogout') }}
                </a-button>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>

      <!-- 移动端卡片 -->
      <div class="cl-card-list">
        <a-skeleton v-if="loading" :animation="true">
          <a-skeleton-line :rows="6" />
        </a-skeleton>
        <template v-else>
          <div
            v-for="(item, index) in mobileData"
            :key="item.token"
            class="cl-card stagger-item"
            :style="{ '--stagger-index': index % mobilePageSize }"
          >
            <div class="cl-card-header">
              <div class="cl-card-identity">
                <strong>{{ item.nickname || item.username }}</strong>
                <span class="cl-card-sub">{{ item.username }}</span>
              </div>
              <a-tag :color="item.ttl > 300 ? 'green' : 'orangered'" size="small">
                {{ item.ttl }}s
              </a-tag>
            </div>
            <div class="cl-card-meta">
              <span>{{ t('monitor.online.clientId') }}：{{ item.clientId || '-' }}</span>
              <span>{{ t('monitor.online.expireAt') }}：{{ item.expireAt || '-' }}</span>
            </div>
            <div class="cl-card-footer">
              <a-button
                v-permission="'monitor:online$forceLogout'"
                size="mini"
                type="outline"
                status="danger"
                @click="handleForceLogout(item)"
              >
                <template #icon><IconDelete /></template>{{ t('monitor.online.forceLogout') }}
              </a-button>
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
