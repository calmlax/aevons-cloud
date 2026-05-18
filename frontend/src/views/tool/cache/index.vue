<script setup lang="ts">
import { IconDelete, IconRefresh, IconEye } from '@arco-design/web-vue/es/icon'
import cacheApi from '@/api/monitor/cache'

// ── 数据 ──────────────────────────────────────────────
interface CacheEntry {
  key: string
  type: string
  ttl: number
  size: number
  expireAt: string
}

interface CacheGroup {
  name: string
  prefix: string
  count: number
  keys: CacheEntry[]
}

const loading = ref(false)
const groups = ref<CacheGroup[]>([])
const activeGroup = ref<string>('')   // currently selected group prefix
const searchKey = ref('')

const currentGroup = computed(() =>
  groups.value.find(g => g.prefix === activeGroup.value)
)

const filteredKeys = computed(() => {
  const list = currentGroup.value?.keys ?? []
  if (!searchKey.value) return list
  const q = searchKey.value.toLowerCase()
  return list.filter(k => k.key.toLowerCase().includes(q))
})

async function loadData() {
  loading.value = true
  try {
    const res: any = await cacheApi.list()
    groups.value = res ?? []
    // auto-select first non-empty group
    if (!activeGroup.value) {
      const first = groups.value.find(g => g.count > 0)
      activeGroup.value = first?.prefix ?? (groups.value[0]?.prefix ?? '')
    }
  } finally {
    loading.value = false
  }
}

// ── 详情弹窗 ──────────────────────────────────────────
const detailVisible = ref(false)
const detailLoading = ref(false)
const detail = ref<any>(null)

async function handleDetail(entry: CacheEntry) {
  detailVisible.value = true
  detailLoading.value = true
  detail.value = null
  try {
    detail.value = await cacheApi.detail(entry.key)
  } finally {
    detailLoading.value = false
  }
}

// ── 删除单个 key ──────────────────────────────────────
async function handleDelete(entry: CacheEntry) {
  Modal.confirm({
    title: '确认删除',
    content: `确认删除缓存键「${entry.key}」吗？`,
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await cacheApi.delete([entry.key])
      Message.success('删除成功')
      loadData()
    },
  })
}

// ── 清空当前分组 ──────────────────────────────────────
async function handleClearGroup() {
  const group = currentGroup.value
  if (!group || group.count === 0) return
  Modal.confirm({
    title: '确认清空',
    content: `确认清空「${group.name}」下的全部 ${group.count} 个缓存键吗？`,
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      if (group.prefix) {
        await cacheApi.deleteByPrefix(group.prefix)
      } else {
        // "其他" group — delete each key individually
        const keys = group.keys.map(k => k.key)
        if (keys.length) await cacheApi.delete(keys)
      }
      Message.success('清空成功')
      loadData()
    },
  })
}

function ttlLabel(ttl: number) {
  if (ttl === -1) return '永久'
  if (ttl < 0) return '-'
  if (ttl < 60) return `${ttl}s`
  if (ttl < 3600) return `${Math.floor(ttl / 60)}m ${ttl % 60}s`
  return `${Math.floor(ttl / 3600)}h ${Math.floor((ttl % 3600) / 60)}m`
}

function typeColor(type: string) {
  const map: Record<string, string> = {
    string: 'arcoblue', hash: 'green', list: 'orange',
    set: 'purple', zset: 'red', stream: 'gray',
  }
  return map[type] ?? 'gray'
}

loadData()
</script>

<template>
  <div class="page-stack">
    <a-row :gutter="12" style="height: 100%">

      <!-- 左侧分组列表 -->
      <a-col :xs="24" :sm="7" :md="6">
        <a-card :bordered="false" class="panel-card cache-group-card">
          <div class="cache-group-header">
            <span class="cache-group-title">缓存分组</span>
            <a-button size="mini" shape="circle" :loading="loading" @click="loadData">
              <template #icon><IconRefresh /></template>
            </a-button>
          </div>
          <a-skeleton v-if="loading" :animation="true">
            <a-skeleton-line :rows="8" />
          </a-skeleton>
          <div v-else class="cache-group-list">
            <div
              v-for="g in groups"
              :key="g.prefix"
              class="cache-group-item"
              :class="{ 'cache-group-item--active': activeGroup === g.prefix }"
              @click="activeGroup = g.prefix; searchKey = ''"
            >
              <span class="cache-group-name">{{ g.name }}</span>
              <a-badge :count="g.count" :max-count="999" :dot="false"
                :color="activeGroup === g.prefix ? '#165dff' : '#86909c'" />
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 右侧 key 列表 -->
      <a-col :xs="24" :sm="17" :md="18">
        <a-card :bordered="false" class="panel-card">
          <!-- 工具栏 -->
          <div class="cl-toolbar">
            <a-space>
              <span class="cache-group-label">{{ currentGroup?.name }}</span>
              <a-button
                v-permission="'cache$delete'"
                status="danger"
                size="small"
                :disabled="!currentGroup || currentGroup.count === 0"
                @click="handleClearGroup"
              >
                <template #icon><IconDelete /></template>
                清空分组
              </a-button>
            </a-space>
            <div class="cl-toolbar-right">
              <a-input-search
                v-model="searchKey"
                placeholder="过滤 key"
                allow-clear
                class="cl-toolbar-search"
              />
            </div>
          </div>

          <!-- key 表格 -->
          <a-skeleton v-if="loading" :animation="true">
            <a-skeleton-line :rows="8" />
          </a-skeleton>
          <a-table
            v-else
            :bordered="false"
            :data="filteredKeys"
            row-key="key"
            :pagination="{ pageSize: 20, showTotal: true, showPageSize: true }"
            :scroll="{ x: 700 }"
          >
            <template #columns>
              <a-table-column title="Key" data-index="key" :width="190" :ellipsis="true" :tooltip="true" />
              <a-table-column title="类型" data-index="type" :width="90" align="center">
                <template #cell="{ record }">
                  <a-tag :color="typeColor(record.type)" size="small">{{ record.type }}</a-tag>
                </template>
              </a-table-column>
              <a-table-column title="剩余 TTL" :width="110" align="center">
                <template #cell="{ record }">
                  <span :style="{ color: record.ttl > 0 && record.ttl < 60 ? 'var(--color-danger-6)' : '' }">
                    {{ ttlLabel(record.ttl) }}
                  </span>
                </template>
              </a-table-column>
              <a-table-column title="大小" :width="90" align="center">
                <template #cell="{ record }">
                  {{ record.size > 0 ? `${record.size} B` : '-' }}
                </template>
              </a-table-column>
              <a-table-column title="过期时间" data-index="expireAt" :width="180" align="center">
                <template #cell="{ record }">
                  {{ record.expireAt || '永久' }}
                </template>
              </a-table-column>
              <a-table-column title="操作" :width="130" align="center">
                <template #cell="{ record }">
                  <a-space size="mini">
                    <a-button size="mini" type="text" @click="handleDetail(record)">
                      <template #icon><IconEye /></template>详情
                    </a-button>
                    <a-button
                      v-permission="'cache$delete'"
                      size="mini"
                      type="text"
                      status="danger"
                      @click="handleDelete(record)"
                    >
                      <template #icon><IconDelete /></template>删除
                    </a-button>
                  </a-space>
                </template>
              </a-table-column>
            </template>
          </a-table>

          <a-empty v-if="!loading && !filteredKeys.length" description="暂无缓存数据" style="padding: 40px 0" />
        </a-card>
      </a-col>
    </a-row>

    <!-- 详情弹窗 -->
    <a-modal v-model:visible="detailVisible" title="缓存详情" :width="640" :footer="false">
      <a-spin :loading="detailLoading">
        <template v-if="detail">
          <a-descriptions :column="2" bordered size="medium">
            <a-descriptions-item label="Key" :span="2">
              <a-typography-text code copyable>{{ detail.key }}</a-typography-text>
            </a-descriptions-item>
            <a-descriptions-item label="类型">
              <a-tag :color="typeColor(detail.type)" size="small">{{ detail.type }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="大小">
              {{ detail.size > 0 ? `${detail.size} B` : '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="剩余 TTL">{{ ttlLabel(detail.ttl) }}</a-descriptions-item>
            <a-descriptions-item label="过期时间">{{ detail.expireAt || '永久' }}</a-descriptions-item>
            <a-descriptions-item v-if="detail.value" label="Value" :span="2">
              <a-typography-paragraph
                :ellipsis="{ rows: 10, expandable: true }"
                style="font-family: monospace; font-size: 12px; white-space: pre-wrap; word-break: break-all; margin: 0"
              >{{ detail.value }}</a-typography-paragraph>
            </a-descriptions-item>
          </a-descriptions>
        </template>
      </a-spin>
    </a-modal>
  </div>
</template>

<style scoped>
.cache-group-card {
  height: 100%;
}

.cache-group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.cache-group-title {
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text-1);
}

.cache-group-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cache-group-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s;
  color: var(--color-text-2);
}

.cache-group-item:hover {
  background: var(--color-fill-2);
}

.cache-group-item--active {
  background: var(--color-primary-light-1);
  color: rgb(var(--primary-6));
  font-weight: 500;
}

.cache-group-name {
  font-size: 13px;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cache-group-label {
  font-weight: 600;
  font-size: 14px;
  color: var(--color-text-1);
}
</style>
