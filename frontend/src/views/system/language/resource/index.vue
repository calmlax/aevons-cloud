<script setup lang="ts">
import { IconDelete, IconEdit, IconPlus, IconRefresh, IconSearch } from '@arco-design/web-vue/es/icon'
import langResourceApi from '@/api/system/lang_resource'
import langApi from '@/api/system/lang'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const proxy = getCurrentInstance()!.proxy as any

const isMobile = ref(window.innerWidth < 576)
function onResize() { isMobile.value = window.innerWidth < 576 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

const dicts = proxy.$useDict('sys_locale_namespace')

// ── 语言列表 ──────────────────────────────────────────
const langOptions = ref<any[]>([])
async function loadLangOptions() {
  const res: any = await langApi.list({})
  langOptions.value = res ?? []
}

// ── 命名空间 ──────────────────────────────────────────
const activeNs = ref('')
const nsLoading = ref(false)

async function loadNamespaces() {
  nsLoading.value = true
  try {
    await selectNamespace('default')
  } finally {
    nsLoading.value = false
  }
}

async function selectNamespace(ns: string) {
  if (dirty.value) {
    const ok = await confirmDiscard()
    if (!ok) return
  }
  activeNs.value = ns
  activeKey.value = ''
  isAdding.value = false
  dirty.value = false
  translations.value = []
  await Promise.all([loadKeysForNs(ns), loadRightList(true)])
}

// ── 左侧 Key 树 ───────────────────────────────────────
const keyCache = ref<Record<string, string[]>>({})
const treeData = ref<any[]>([])
const treeLoading = ref(false)
const treeSearch = ref('')
const expandedKeys = ref<string[]>([])

async function loadKeysForNs(ns: string) {
  treeLoading.value = true
  try {
    if (!keyCache.value[ns]) {
      const res: any = await langResourceApi.getKeysByNamespace(ns)
      keyCache.value[ns] = res ?? []
    }
    rebuildTree(keyCache.value[ns])
  } finally {
    treeLoading.value = false
  }
}

function buildKeyTree(keys: string[]): any[] {
  const root: Record<string, any> = {}
  for (const key of keys) {
    const parts = key.split('.')
    let node = root
    let path = ''
    for (let i = 0; i < parts.length; i++) {
      const part = parts[i]
      path = path ? `${path}.${part}` : part
      if (!node[part]) {
        node[part] = { key: path, title: part, isLeaf: false, _fullKey: undefined, children: {} }
      }
      if (i === parts.length - 1) {
        node[part]._fullKey = key
        node[part].isLeaf = Object.keys(node[part].children).length === 0
      }
      node = node[part].children
    }
  }
  function toArray(obj: Record<string, any>): any[] {
    return Object.values(obj).map(n => {
      const children = toArray(n.children)
      return { key: n.key, title: n.title, isLeaf: children.length === 0, _fullKey: n._fullKey, children: children.length ? children : undefined }
    }).sort((a, b) => a.title.localeCompare(b.title))
  }
  return toArray(root)
}

function rebuildTree(keys: string[]) {
  const kw = treeSearch.value.trim().toLowerCase()
  const filtered = kw ? keys.filter(k => k.toLowerCase().includes(kw)) : keys
  treeData.value = buildKeyTree(filtered)
  expandedKeys.value = kw ? collectAllKeys(treeData.value) : treeData.value.map(n => n.key)
}

function collectAllKeys(nodes: any[]): string[] {
  const result: string[] = []
  for (const n of nodes) { result.push(n.key); if (n.children) result.push(...collectAllKeys(n.children)) }
  return result
}

let treeSearchTimer: ReturnType<typeof setTimeout> | null = null
function onTreeSearch() {
  if (treeSearchTimer) clearTimeout(treeSearchTimer)
  treeSearchTimer = setTimeout(() => {
    if (activeNs.value && keyCache.value[activeNs.value]) rebuildTree(keyCache.value[activeNs.value])
  }, 200)
}

function onTreeSelect(_keys: (string | number)[], data: { node?: any }) {
  const node = data.node
  if (!node?._fullKey) return
  selectKey(node._fullKey)
}

// ── 右侧列表模式（懒加载分页）────────────────────────
const rightList = ref<string[]>([])
const rightListTotal = ref(0)
const rightListLoading = ref(false)
const rightListPage = ref(1)
const rightListPageSize = 20
const rightSearch = reactive({ resourceKey: '', content: '' })
let rightSearchTimer: ReturnType<typeof setTimeout> | null = null

async function loadRightList(reset = false) {
  if (reset) { rightListPage.value = 1; rightList.value = [] }
  rightListLoading.value = true
  try {
    const res: any = await langResourceApi.pageKeys({
      namespace: activeNs.value,
      resourceKey: rightSearch.resourceKey || undefined,
      content: rightSearch.content || undefined,
      pageNum: rightListPage.value,
      pageSize: rightListPageSize,
    })
    const rows: string[] = res?.rows ?? []
    rightListTotal.value = res?.total ?? 0
    if (reset) rightList.value = rows
    else rightList.value.push(...rows)
  } finally {
    rightListLoading.value = false
  }
}

function onRightSearch() {
  if (rightSearchTimer) clearTimeout(rightSearchTimer)
  rightSearchTimer = setTimeout(() => loadRightList(true), 250)
}

function loadMoreRight() {
  if (rightList.value.length >= rightListTotal.value) return
  rightListPage.value++
  loadRightList(false)
}

// 右侧列表哨兵
const rightSentinelRef = ref<HTMLElement | null>(null)
let rightObserver: IntersectionObserver | null = null
function setupRightObserver() {
  rightObserver?.disconnect()
  if (!rightSentinelRef.value) return
  rightObserver = new IntersectionObserver(entries => {
    if (entries[0].isIntersecting) loadMoreRight()
  }, { rootMargin: '60px' })
  rightObserver.observe(rightSentinelRef.value)
}
watch(rightSentinelRef, () => setupRightObserver())
onUnmounted(() => rightObserver?.disconnect())

// 右侧列表中的 key 展开翻译（内联展示）
const expandedListKey = ref('')
const listTranslations = ref<{ langCode: string; content: string }[]>([])
const listTransLoading = ref(false)

async function toggleListKey(key: string) {
  if (expandedListKey.value === key) { expandedListKey.value = ''; return }
  expandedListKey.value = key
  listTransLoading.value = true
  try {
    const res: any = await langResourceApi.getTranslations(activeNs.value, key)
    const existing: any[] = res ?? []
    listTranslations.value = langOptions.value.map(l => {
      const found = existing.find((e: any) => e.langCode === l.langCode)
      return { langCode: l.langCode, content: found?.content ?? '' }
    })
  } finally {
    listTransLoading.value = false
  }
}

// ── 右侧编辑模式 ──────────────────────────────────────
const activeKey = ref('')
const translations = ref<{ langCode: string; content: string }[]>([])
const rightLoading = ref(false)
const saveLoading = ref(false)
const dirty = ref(false)

// 右侧视图：'list' | 'edit' | 'add'
const rightView = ref<'list' | 'edit' | 'add'>('list')

async function confirmDiscard(): Promise<boolean> {
  return new Promise(resolve => {
    Modal.confirm({
      title: t('system.langResource.confirmDiscard'),
      content: t('system.langResource.confirmDiscardContent'),
      okText: t('system.langResource.discardOk'),
      cancelText: t('common.cancel'),
      onOk: () => resolve(true), onCancel: () => resolve(false),
    })
  })
}

async function selectKey(key: string) {
  if (dirty.value) { if (!(await confirmDiscard())) return }
  activeKey.value = key
  dirty.value = false
  isAdding.value = false
  rightView.value = 'edit'
  rightLoading.value = true
  try {
    const res: any = await langResourceApi.getTranslations(activeNs.value, key)
    const existing: any[] = res ?? []
    translations.value = langOptions.value.map(l => {
      const found = existing.find((e: any) => e.langCode === l.langCode)
      return { langCode: l.langCode, content: found?.content ?? '' }
    })
  } finally {
    rightLoading.value = false
  }
}

function backToList() {
  activeKey.value = ''
  isAdding.value = false
  dirty.value = false
  rightView.value = 'list'
}

// ── 新增模式 ──────────────────────────────────────────
const isAdding = ref(false)
const addForm = reactive({ namespace: '', resourceKey: '' })
const addFormRef = ref()

function validateResourceKey(value: string, callback: (err?: string) => void) {
  if (!value) { callback(); return }
  const parts = value.split('.')
  if (parts.length > 4) { callback(t('system.langResource.resourceKeyMaxDepth')); return }
  const segReg = /^[a-zA-Z0-9_-]+$/
  for (const p of parts) {
    if (!p) { callback(t('system.langResource.resourceKeyEmptySegment')); return }
    if (!segReg.test(p)) { callback(t('system.langResource.resourceKeyInvalidChar')); return }
  }
  callback()
}

async function openAddPanel() {
  if (dirty.value) { if (!(await confirmDiscard())) return }
  isAdding.value = true
  activeKey.value = ''
  dirty.value = false
  rightView.value = 'add'
  addForm.namespace = activeNs.value
  addForm.resourceKey = ''
  translations.value = langOptions.value.map(l => ({ langCode: l.langCode, content: '' }))
}

async function openAddPanelWithPrefix(prefix: string) {
  if (dirty.value) { if (!(await confirmDiscard())) return }
  isAdding.value = true
  activeKey.value = ''
  dirty.value = false
  rightView.value = 'add'
  addForm.namespace = activeNs.value
  addForm.resourceKey = prefix + '.'
  translations.value = langOptions.value.map(l => ({ langCode: l.langCode, content: '' }))
}

async function saveTranslations() {
  if (rightView.value === 'add') {
    const valid = await addFormRef.value?.validate()
    if (valid) return
  }
  saveLoading.value = true
  const ns = rightView.value === 'add' ? addForm.namespace : activeNs.value
  const key = rightView.value === 'add' ? addForm.resourceKey : activeKey.value
  try {
    await langResourceApi.saveTranslations({
      namespace: ns,
      resourceKey: key,
      items: translations.value.map(t => ({ langCode: t.langCode, content: t.content })),
    })
    Message.success(t('system.langResource.saveSuccess'))
    dirty.value = false
    delete keyCache.value[ns]
    if (activeNs.value !== ns) await selectNamespace(ns)
    else { await loadKeysForNs(ns); await loadRightList(true) }
    activeKey.value = key
    rightView.value = 'edit'
    isAdding.value = false
  } finally {
    saveLoading.value = false
  }
}

async function deleteKey() {
  Modal.confirm({
    title: t('system.langResource.deleteTitle'),
    content: t('system.langResource.deleteContent', { ns: activeNs.value, key: activeKey.value }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      const res: any = await langResourceApi.getTranslations(activeNs.value, activeKey.value)
      const ids = (res ?? []).map((r: any) => r.id).filter(Boolean)
      if (ids.length) await langResourceApi.delete(ids)
      Message.success(t('system.langResource.deleteSuccess'))
      delete keyCache.value[activeNs.value]
      await loadKeysForNs(activeNs.value)
      await loadRightList(true)
      backToList()
    },
  })
}

async function deleteListKey(key: string) {
  Modal.confirm({
    title: t('system.langResource.deleteTitle'),
    content: t('system.langResource.deleteContent', { ns: activeNs.value, key }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      const res: any = await langResourceApi.getTranslations(activeNs.value, key)
      const ids = (res ?? []).map((r: any) => r.id).filter(Boolean)
      if (ids.length) await langResourceApi.delete(ids)
      Message.success(t('system.langResource.deleteSuccess'))
      delete keyCache.value[activeNs.value]
      await loadKeysForNs(activeNs.value)
      await loadRightList(true)
    },
  })
}

// ── 初始化 ────────────────────────────────────────────
onMounted(async () => {
  await loadLangOptions()
  await loadNamespaces()
})
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">
      <div class="res-layout" :class="{ 'res-layout--mobile': isMobile }">

        <!-- ══ 左侧 ══ -->
        <div class="res-left">
          <div class="res-ns-bar">
            <a-spin :loading="nsLoading" style="width:100%">
              <div class="res-ns-list">
                <div v-for="dict in dicts.sys_locale_namespace" :key="dict.dictValue"
                  class="res-ns-item" :class="{ active: activeNs === dict.dictValue }"
                  @click="selectNamespace(dict.dictValue)">{{ dict.label }}</div>
              </div>
            </a-spin>
          </div>
          <div class="res-left-toolbar">
            <a-input v-model="treeSearch" :placeholder="t('system.langResource.searchKeyPlaceholder')" allow-clear size="small"
              @input="onTreeSearch" @clear="onTreeSearch">
              <template #prefix><IconSearch /></template>
            </a-input>
            <a-button size="mini" type="text" :title="t('system.langResource.refreshTooltip')"
              @click="() => { delete keyCache[activeNs]; loadKeysForNs(activeNs); loadRightList(true) }">
              <template #icon><IconRefresh /></template>
            </a-button>
            <a-button size="mini" type="primary" @click="openAddPanel">
              <template #icon><IconPlus /></template>
            </a-button>
          </div>
          <div class="res-tree-wrap">
            <a-spin :loading="treeLoading" style="width:100%;height:100%">
              <a-tree v-if="treeData.length" :data="treeData"
                :virtual-list-props="{ height: '100%' }"
                v-model:expanded-keys="expandedKeys"
                :show-line="true" size="small"
                @select="onTreeSelect">
                <template #title="node">
                  <span class="res-tree-node">
                    <span :class="{ 'res-key-leaf': node.isLeaf, 'res-key-active': node._fullKey === activeKey }">
                      {{ node.title }}
                    </span>
                    <span class="res-node-add" @click.stop="openAddPanelWithPrefix(node.key)">+</span>
                  </span>
                </template>
              </a-tree>
              <a-empty v-else :description="t('system.langResource.noResource')" style="padding:24px" />
            </a-spin>
          </div>
        </div>

        <!-- ══ 右侧 ══ -->
        <div class="res-right">

          <!-- ── 列表模式 ── -->
          <template v-if="rightView === 'list'">
            <div class="res-list-toolbar">
              <a-input v-model="rightSearch.resourceKey" :placeholder="t('system.langResource.searchKeyPlaceholder')" allow-clear size="small"
                style="flex:1" @input="onRightSearch" @clear="onRightSearch">
                <template #prefix><IconSearch /></template>
              </a-input>
              <a-input v-model="rightSearch.content" :placeholder="t('system.langResource.searchContentPlaceholder')" allow-clear size="small"
                style="flex:1" @input="onRightSearch" @clear="onRightSearch">
                <template #prefix><IconSearch /></template>
              </a-input>
            </div>
            <a-spin :loading="rightListLoading && rightList.length === 0" style="flex:1;overflow-y:auto">
              <div class="res-list">
                <div v-for="item in rightList" :key="item" class="res-list-item">
                  <div class="res-list-item-header" @click="toggleListKey(item)">
                    <span class="res-list-key">{{ item }}</span>
                    <a-space size="mini">
                      <a-button size="mini" type="text" @click.stop="selectKey(item)">
                        <template #icon><IconEdit /></template>
                      </a-button>
                      <a-button size="mini" type="text" status="danger" @click.stop="deleteListKey(item)">
                        <template #icon><IconDelete /></template>
                      </a-button>
                    </a-space>
                  </div>
                  <!-- 内联展开翻译预览 -->
                  <div v-if="expandedListKey === item" class="res-list-trans">
                    <a-spin :loading="listTransLoading" style="width:100%">
                      <div v-for="t in listTranslations" :key="t.langCode" class="res-list-trans-row">
                        <span class="res-list-lang">{{ t.langCode }}</span>
                        <span class="res-list-content">{{ t.content || '—' }}</span>
                      </div>
                    </a-spin>
                  </div>
                </div>
                <a-empty v-if="!rightListLoading && !rightList.length" :description="t('system.langResource.noResource')" style="padding:40px 0" />
              </div>
              <!-- 懒加载哨兵 -->
              <div ref="rightSentinelRef" class="res-sentinel">
                <a-spin v-if="rightListLoading && rightList.length > 0" />
                <span v-else-if="rightList.length >= rightListTotal && rightList.length > 0" class="res-no-more">{{ t('system.langResource.allLoaded') }}</span>
              </div>
            </a-spin>
          </template>

          <!-- ── 编辑 / 新增模式 ── -->
          <template v-else>
            <!-- 返回按钮 -->
            <div class="res-back-bar">
              <a-button size="small" type="text" @click="backToList">{{ t('system.langResource.backToList') }}</a-button>
            </div>

            <!-- 新增模式头部 -->
            <template v-if="rightView === 'add'">
              <div class="res-right-title">{{ t('system.langResource.addTitle') }}</div>
              <a-divider style="margin:10px 0" />
              <a-form :model="addForm" ref="addFormRef" layout="horizontal"
                :label-col-props="{ span: 4 }" :wrapper-col-props="{ span: 20 }">
                <a-form-item :label="t('system.langResource.namespace')" field="namespace"
                  :rules="[{ required: true, message: t('system.langResource.namespaceRequired') }, { maxLength: 32, message: t('common.maxLength', { max: 32 }) }]"
                  validate-trigger="blur">
                  <a-input v-model="addForm.namespace" :placeholder="t('system.langResource.namespacePlaceholder')" allow-clear @input="dirty = true" />
                </a-form-item>
                <a-form-item :label="t('system.langResource.resourceKey')" field="resourceKey"
                  :rules="[
                    { required: true, message: t('system.langResource.resourceKeyRequired') },
                    { maxLength: 100, message: t('system.langResource.resourceKeyMaxLength') },
                    { validator: validateResourceKey }
                  ]"
                  validate-trigger="blur"
                  :extra="t('system.langResource.resourceKeyExtra')">
                  <a-input v-model="addForm.resourceKey" :placeholder="t('system.langResource.resourceKeyPlaceholder')" allow-clear @input="dirty = true" />
                </a-form-item>
              </a-form>
              <a-divider style="margin:10px 0" />
            </template>

            <!-- 编辑模式头部 -->
            <template v-else>
              <div class="res-right-header">
                <div>
                  <div class="res-right-key">{{ activeKey }}</div>
                  <div class="res-right-ns">{{ activeNs }}</div>
                </div>
                <a-button status="danger" size="small" @click="deleteKey">
                  <template #icon><IconDelete /></template>{{ t('system.langResource.deleteBtn') }}
                </a-button>
              </div>
              <a-divider style="margin:10px 0" />
            </template>

            <!-- 翻译输入区 -->
            <a-spin :loading="rightLoading" style="width:100%">
              <div class="res-translations">
                <div v-for="t in translations" :key="t.langCode" class="res-trans-row">
                  <div class="res-trans-lang">
                    <span class="res-lang-code">{{ t.langCode }}</span>
                    <span class="res-lang-name">{{ langOptions.find(l => l.langCode === t.langCode)?.langName }}</span>
                  </div>
                  <a-textarea v-model="t.content" :placeholder="$t('system.langResource.translationPlaceholder', { lang: t.langCode })" allow-clear
                    :auto-size="{ minRows: 1, maxRows: 4 }" class="res-trans-input"
                    @input="dirty = true" />
                </div>
              </div>
            </a-spin>

            <div class="res-right-footer">
              <a-button @click="backToList">{{ t('common.cancel') }}</a-button>
              <a-button type="primary" :loading="saveLoading" @click="saveTranslations">{{ t('common.save') }}</a-button>
            </div>
          </template>

        </div>
      </div>
    </a-card>
  </div>
</template>

<style scoped>
.res-layout {
  display: flex;
  height: calc(100vh - 180px);
  min-height: 500px;
  overflow: hidden;
}
.res-layout--mobile { flex-direction: column; height: auto; }

/* ── 左侧 ── */
.res-left {
  width: 280px;
  flex-shrink: 0;
  border-right: 1px solid var(--color-border-1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.res-layout--mobile .res-left {
  width: 100%;
  border-right: none;
  border-bottom: 1px solid var(--color-border-1);
  max-height: 320px;
}
.res-ns-bar { border-bottom: 1px solid var(--color-border-1); padding: 12px 8px; flex-shrink: 0; }
.res-ns-list { display: flex; gap: 4px; flex-wrap: wrap; }
.res-ns-item {
  padding: 2px 10px; border-radius: 12px; font-size: 12px; cursor: pointer;
  color: var(--color-text-2); border: 1px solid transparent; white-space: nowrap; transition: all 0.15s;
}
.res-ns-item:hover { background: var(--color-fill-2); }
.res-ns-item.active {
  background: var(--color-primary-light-1); color: rgb(var(--primary-6));
  border-color: rgb(var(--primary-3)); font-weight: 600;
}
.res-left-toolbar { display: flex; align-items: center; gap: 6px; padding: 8px 8px 4px; flex-shrink: 0; }
.res-left-toolbar .arco-input-wrapper { flex: 1; }
.res-tree-wrap { flex: 1; overflow: hidden; padding: 4px 0; }
.res-tree-wrap :deep(.arco-tree) { height: 100%; }
/* 让 arco tree 节点 title 区域撑满宽度 */
.res-tree-wrap :deep(.arco-tree-node-title) { flex: 1; min-width: 0; }
.res-tree-wrap :deep(.arco-tree-node-title-text) { display: flex; width: 100%; }

.res-tree-node { display: flex; align-items: center; width: 100%; }
.res-tree-node > span:first-child { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.res-node-add {
  opacity: 0;
  font-size: 16px;
  font-weight: 600;
  line-height: 1;
  color: rgb(var(--primary-6));
  cursor: pointer;
  padding: 1px 6px;
  border-radius: 3px;
  flex-shrink: 0;
}
.res-node-add:hover { background: var(--color-primary-light-2); }
.res-tree-wrap :deep(.arco-tree-node:hover) .res-node-add { opacity: 1; }

/* ── 右侧 ── */
.res-right {
  flex: 1; overflow: hidden; display: flex; flex-direction: column; min-width: 0;
}

/* 列表模式 */
.res-list-toolbar {
  display: flex; gap: 8px; padding: 12px 16px 8px; flex-shrink: 0;
  border-bottom: 1px solid var(--color-border-1);
}
.res-list { padding: 0 16px; }
.res-list-item {
  border-bottom: 1px solid var(--color-border-1);
}
.res-list-item-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 10px 0; cursor: pointer; gap: 8px;
}
.res-list-item-header:hover { background: var(--color-fill-1); margin: 0 -16px; padding: 10px 16px; }
.res-list-key { font-size: 13px; font-weight: 500; color: var(--color-text-1); flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.res-list-trans { padding: 8px 0 12px; display: flex; flex-direction: column; gap: 6px; }
.res-list-trans-row { display: flex; gap: 12px; font-size: 12px; }
.res-list-lang { width: 60px; flex-shrink: 0; color: var(--color-text-3); font-weight: 600; }
.res-list-content { color: var(--color-text-2); flex: 1; }
.res-sentinel { padding: 12px; text-align: center; }
.res-no-more { font-size: 12px; color: var(--color-text-4); }

/* 编辑/新增模式 */
.res-back-bar { padding: 8px 16px; flex-shrink: 0; border-bottom: 1px solid var(--color-border-1); }
.res-right-header {
  display: flex; align-items: flex-start; justify-content: space-between;
  gap: 12px; padding: 12px 20px 0;
}
.res-right-title { font-size: 15px; font-weight: 700; color: var(--color-text-1); padding: 12px 20px 0; }
.res-right-key { font-size: 16px; font-weight: 700; color: var(--color-text-1); }
.res-right-ns { font-size: 12px; color: var(--color-text-3); margin-top: 2px; }
.res-right-empty { flex: 1; display: flex; align-items: center; justify-content: center; }
.res-right-footer {
  display: flex; justify-content: center; gap: 12px;
  margin-top: 24px; padding: 16px 20px;
  border-top: 1px solid var(--color-border-1); flex-shrink: 0;
}
.res-translations { display: flex; flex-direction: column; gap: 16px; padding: 0 20px; overflow-y: auto; flex: 1; }
.res-trans-row { display: flex; align-items: flex-start; gap: 12px; }
.res-trans-lang { width: 90px; flex-shrink: 0; padding-top: 6px; }
.res-lang-code { font-size: 13px; font-weight: 600; color: var(--color-text-1); display: block; }
.res-lang-name { font-size: 11px; color: var(--color-text-3); }
.res-trans-input { flex: 1; }

/* form padding fix */
.res-right :deep(.arco-form) { padding: 0 20px; }
</style>
