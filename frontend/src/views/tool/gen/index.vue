<script setup lang="ts">
import { IconDelete, IconEdit, IconImport, IconSettings, IconFilter, IconPlus, IconDownload, IconRefresh, IconEye, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import type { TableColumnData } from '@arco-design/web-vue';
import { saveAs } from 'file-saver';
import { getHighlighter, resolveLanguage } from '@/utils/shiki';
import {
  getTablePage, updateTable, deleteTable,
  getColumnList, batchUpdateColumn,
  getDBTableList, importTables,
  previewCode, synchTable, downloadCode,
  type GenTable, type GenTableColumn, type DBTableInfo, type PreviewFile,
} from '@/api/gen';
import dictApi from "@/api/system/dict"
import menuApi from "@/api/system/menu"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const proxy = getCurrentInstance()!.proxy as any


const dicts = proxy.$useDict('sys_condition','xcode_component')

// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

// ── 搜索 & 筛选 ────────────────────────────────────────
const keyword = ref('');
const filterModule = ref('');
const advancedVisible = ref(false);

// ── PC 分页 ────────────────────────────────────────────
const current = ref(1);
const pageSize = ref(10);
const tableTotal = ref(0);
const tableLoading = ref(false);
const tableList = ref<GenTable[]>([]);

async function loadTablePage() {
  tableLoading.value = true;
  try {
    const res = await getTablePage({
      pageNum: current.value,
      pageSize: pageSize.value,
      tableName: keyword.value.trim(),
      tableComment: filterModule.value.trim(),
    });
    tableList.value = res?.rows ?? [];
    tableTotal.value = res?.total ?? 0;
  } finally {
    tableLoading.value = false;
  }
}

function onSearch() { current.value = 1; mobilePage.value = mobilePageSize; loadTablePage(); }
function onReset() { keyword.value = ''; filterModule.value = ''; onSearch(); }

// ── 移动端无限滚动 ─────────────────────────────────────
const mobilePageSize = 10;
const mobilePage = ref(mobilePageSize);
const mobileLoadingMore = ref(false);
const sentinelRef = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

// 移动端直接用 tableList 切片（已从服务端拿全量当页数据）
const mobileData = computed(() => tableList.value.slice(0, mobilePage.value));
const mobileHasMore = computed(() => mobilePage.value < tableList.value.length);

function loadMore() {
  if (!mobileHasMore.value || mobileLoadingMore.value) return;
  mobileLoadingMore.value = true;
  setTimeout(() => { mobilePage.value += mobilePageSize; mobileLoadingMore.value = false; }, 300);
}

function setupObserver() {
  if (!sentinelRef.value) return;
  observer?.disconnect();
  observer = new IntersectionObserver(entries => { if (entries[0].isIntersecting) loadMore(); }, { rootMargin: '80px' });
  observer.observe(sentinelRef.value);
}

onMounted(() => setupObserver());
onUnmounted(() => observer?.disconnect());
watch(tableList, async () => {
  mobilePage.value = mobilePageSize;
  await nextTick();
  setupObserver();
});
// ── 多选 ───────────────────────────────────────────────
const selectedTableIds = ref<string[]>([]);
const tableRowSelection = reactive({ type: 'checkbox', showCheckedAll: true, onlyCurrent: false });

// ── 移动端多选模式 ─────────────────────────────────────
const mobileSelectMode = ref(false);

// 切换移动端选择模式
function toggleMobileSelectMode() {
  mobileSelectMode.value = !mobileSelectMode.value;
  if (!mobileSelectMode.value) {
    selectedTableIds.value = [];
  }
}

// 切换单个卡片选中状态
function toggleCardSelect(id: string) {
  const idx = selectedTableIds.value.indexOf(id);
  if (idx === -1) {
    selectedTableIds.value.push(id);
  } else {
    selectedTableIds.value.splice(idx, 1);
  }
}

// 判断卡片是否选中
function isCardSelected(id: string) {
  return selectedTableIds.value.includes(id);
}

// 全选/取消全选（移动端）
const mobileAllSelected = computed(() =>
  mobileData.value.length > 0 && mobileData.value.every(item => selectedTableIds.value.includes(item.id))
);
const mobileIndeterminate = computed(() =>
  selectedTableIds.value.length > 0 && !mobileAllSelected.value
);

function toggleMobileSelectAll() {
  if (mobileAllSelected.value) {
    // 取消全选：移除当前显示的所有项
    const currentIds = mobileData.value.map(item => item.id);
    selectedTableIds.value = selectedTableIds.value.filter(id => !currentIds.includes(id));
  } else {
    // 全选：添加当前显示的所有项
    const currentIds = mobileData.value.map(item => item.id);
    const newIds = currentIds.filter(id => !selectedTableIds.value.includes(id));
    selectedTableIds.value.push(...newIds);
  }
}

async function batchDeleteTable() {
  if (!selectedTableIds.value.length) { Message.warning(t('common.confirmDeleteContent', { count: 0 }).replace('0 ', '')); return; }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedTableIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await deleteTable(selectedTableIds.value);
      Message.success(t('common.deleteSuccess'));
      selectedTableIds.value = [];
      loadTablePage();
    },
  });
}

// ── 表格列定义 ─────────────────────────────────────────
const tableColumns = computed<TableColumnData[]>(() => [
  { title: t('tool.gen.colTableName'), dataIndex: 'tableName', ellipsis: true, tooltip: true, width: 160 },
  { title: t('tool.gen.colTableComment'), dataIndex: 'tableComment', ellipsis: true, tooltip: true, width: 140 },
  { title: t('tool.gen.colClassName'), dataIndex: 'className', ellipsis: true, width: 140 },
  { title: t('tool.gen.colModuleName'), dataIndex: 'moduleName', width: 110 },
  { title: t('tool.gen.colAuthor'), dataIndex: 'author', width: 90, ellipsis: true, tooltip: true },
  { title: t('tool.gen.colRouter'), dataIndex: 'router', ellipsis: true, width: 120 },
  { title: t('tool.gen.colAction'), slotName: 'actions', width: 340 },
]);

// ── 编辑弹窗 ───────────────────────────────────────────
const tableModalVisible = ref(false);
const tableForm = reactive<Partial<GenTable>>({});
const editingTableId = ref('');

function openTableEdit(record: GenTable) {
  editingTableId.value = record.id;
  Object.assign(tableForm, { ...record });
  tableModalVisible.value = true;
}

async function submitTableEdit() {
  await updateTable(editingTableId.value, tableForm);
  Message.success(t('tool.gen.updateSuccess'));
  tableModalVisible.value = false;
  loadTablePage();
}

async function confirmDeleteTable(record: GenTable) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('tool.gen.confirmDelete', { name: record.tableName }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await deleteTable([record.id]);
      Message.success(t('common.deleteSuccess'));
      loadTablePage();
    },
  });
}

// ── 导入弹窗 ───────────────────────────────────────────
const importModalVisible = ref(false);
const importLoading = ref(false);
const dbTableList = ref<DBTableInfo[]>([]);
const importSelectedKeys = ref<string[]>([]);
const importSearchKeyword = ref('');
const importRowSelection = reactive({ type: 'checkbox', showCheckedAll: true, onlyCurrent: false });

const filteredDBTables = computed(() => {
  const kw = importSearchKeyword.value.trim().toLowerCase();
  if (!kw) return dbTableList.value;
  return dbTableList.value.filter(t => t.tableName.toLowerCase().includes(kw) || t.tableComment.toLowerCase().includes(kw));
});

async function openImportModal() {
  importModalVisible.value = true;
  importLoading.value = true;
  importSelectedKeys.value = [];
  importSearchKeyword.value = '';
  try { dbTableList.value = await getDBTableList() ?? []; }
  finally { importLoading.value = false; }
}

async function submitImport() {
  if (!importSelectedKeys.value.length) { Message.warning(t('tool.gen.importSelectWarn')); return; }
  importLoading.value = true;
  try {
    await importTables(importSelectedKeys.value);
    Message.success(t('tool.gen.importSuccess', { count: importSelectedKeys.value.length }));
    importModalVisible.value = false;
    loadTablePage();
  } finally { importLoading.value = false; }
}

const importColumns = computed<TableColumnData[]>(() => [
  { title: t('tool.gen.importColTableName'), dataIndex: 'tableName', width: 220 },
  { title: t('tool.gen.importColTableComment'), dataIndex: 'tableComment' },
]);

// ── 设计抽屉 ───────────────────────────────────────────
const designDrawerVisible = ref(false);
const designTableName = ref('');
const designLoading = ref(false);
const designSaving = ref(false);
const designColumns = ref<GenTableColumn[]>([]);
const dictList = ref<any[]>([]);

async function openDesignDrawer(record: GenTable) {
  designTableName.value = record.tableName;
  designDrawerVisible.value = true;
  designLoading.value = true;
  try {
    const res = await getColumnList({ tableId: record.id });
    designColumns.value = (res ?? []).map(c => ({ ...c }));
  } finally { designLoading.value = false; }
}

async function saveDesign() {
  designSaving.value = true;
  try {
    await batchUpdateColumn(designColumns.value);
    Message.success(t('tool.gen.saveSuccess'));
    designDrawerVisible.value = false;
  } finally { designSaving.value = false; }
}

// ── 代码预览 ───────────────────────────────────────────
const previewModalVisible = ref(false);
const previewLoading = ref(false);
const previewTableName = ref('');
const previewFiles = ref<PreviewFile[]>([]);
const activePreviewFile = ref('');
const highlightedHtml = ref('');
const highlightLoading = ref(false);
const isDark = ref(ref(document.documentElement.dataset.theme === 'dark'));

const activeFile = computed(() =>
  previewFiles.value.find(f => f.name === activePreviewFile.value)
);

async function renderHighlight(file: PreviewFile) {
  highlightLoading.value = true;
  highlightedHtml.value = '';
  try {
    const hl = await getHighlighter();
    const lang = resolveLanguage(file.fileType);
    highlightedHtml.value = hl.codeToHtml(file.code, {
      lang,
      theme: isDark.value ? 'github-dark' : 'github-light',
    });
  } finally {
    highlightLoading.value = false;
  }
}

watch(activeFile, (f) => { if (f) renderHighlight(f); });

async function openPreview(record: GenTable) {
  previewTableName.value = record.tableName;
  previewModalVisible.value = true;
  previewLoading.value = true;
  previewFiles.value = [];
  activePreviewFile.value = '';
  highlightedHtml.value = '';
  try {
    const res = await previewCode(record.id);
    previewFiles.value = res ?? [];
    if (previewFiles.value.length) activePreviewFile.value = previewFiles.value[0].name;
  } finally {
    previewLoading.value = false;
  }
}

function copyCode() {
  const code = activeFile.value?.code ?? '';
  navigator.clipboard.writeText(code);
  Message.success(t('tool.gen.copySuccess'));
}

// ── 同步表结构 ─────────────────────────────────────────
async function syncTable(record: GenTable) {
  Modal.confirm({
    title: t('tool.gen.syncTitle'),
    content: t('tool.gen.syncContent', { name: record.tableName }),
    okButtonProps: { status: 'warning' },
    onOk: async () => {
      await synchTable(record.id);
      Message.success(t('tool.gen.syncSuccess'));
      loadTablePage();
    },
  });
}

// ── 代码下载 ───────────────────────────────────────────
const downloadLoading = ref(false);

async function handleDownload(ids: string[]) {
  if (!ids.length) { Message.warning(t('tool.gen.downloadSelectWarn')); return; }
  downloadLoading.value = true;
  try {
    const blob = await downloadCode(ids) as unknown as Blob;
    saveAs(blob, 'aevo-gen-code.zip');
    Message.success(t('tool.gen.downloadSuccess'));
  } catch {
    Message.error(t('tool.gen.downloadFailed'));
  } finally {
    downloadLoading.value = false;
  }
}

async function loadDictList(){
  dictList.value = await dictApi.list({}) as any
}

// ── 菜单树（上级菜单选择）─────────────────────────────
const menuRawList = ref<any[]>([])
const menuTreeOptions = computed(() => [
  ...buildMenuSelectTree(menuRawList.value.filter(i => i.type !== 3)),
])

function buildMenuSelectTree(items: any[], parentId = 0): any[] {
  return items
    .filter(i => Number(i.parentId) === parentId)
    .sort((a, b) => Number(a.sort) - Number(b.sort))
    .map(i => {
      const children = buildMenuSelectTree(items, Number(i.id))
      return { key: String(i.id), value: String(i.id), title: i.title, ...(children.length ? { children } : {}) }
    })
}

async function loadMenuList() {
  const res: any = await menuApi.list()
  menuRawList.value = res ?? []
}

// ── 初始化 ─────────────────────────────────────────────
loadTablePage();
loadDictList()
loadMenuList()
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">
      <div style="padding-bottom: 15px;">
        <a-alert :show-icon="false">
          <h4>{{ t('tool.gen.cmdTitle') }}</h4>
          <pre class="language-shell" style="background-color: black;color:#fff"> 
 $ cd aevons-cloud/gen-service
 $ go run ./cmd/server/main.go gen &lt;表名1,表名2...&gt; &lt;模块名称&gt;
          </pre>
          <ul>
            <li>{{ t('tool.gen.cmdDesc1') }}</li>
            <li>{{ t('tool.gen.cmdDesc2') }}</li>
            <li>{{ t('tool.gen.cmdDesc3') }}</li>
          </ul>
        </a-alert>
      </div>
      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <!-- 左：操作按钮 -->
        <a-space>
          <a-button v-permission="'gen:table$import'" type="primary" @click="openImportModal">
            <template #icon><IconImport /></template>{{ t('tool.gen.importBtn') }}
          </a-button>
          <a-button v-permission="'gen:table$download'" type="primary" status="success"
            :disabled="!selectedTableIds.length" :loading="downloadLoading"
            @click="handleDownload(selectedTableIds)">
            <template #icon><IconDownload /></template>{{ t('tool.gen.batchDownload') }}
          </a-button>
          <a-button v-permission="'gen:table$delete'" status="danger" :disabled="!selectedTableIds.length" @click="batchDeleteTable">
            <template #icon><IconDelete /></template>{{ t('tool.gen.batchDelete') }}
          </a-button>
        </a-space>

        <!-- 右：搜索 + 筛选 -->
        <div class="cl-toolbar-right">
          <a-input-search v-model="keyword" :placeholder="t('tool.gen.searchPlaceholder')" allow-clear class="cl-toolbar-search"
            @search="onSearch"
            @press-enter="onSearch"
          />
          <a-popover
            v-model:popup-visible="advancedVisible"
            trigger="click"
            position="br"
            popup-container="body"
          >
            <a-button shape="circle" :type="filterModule ? 'primary' : 'secondary'">
              <template #icon><IconFilter /></template>
            </a-button>
            <template #content>
              <div class="cl-filter-panel">
                <p class="cl-filter-title">{{ t('common.advancedFilter') }}</p>
                <a-input v-model="filterModule" :placeholder="t('tool.gen.filterPlaceholder')" allow-clear />
                <div class="cl-filter-actions">
                  <a-button size="small" @click="onReset">{{ t('common.reset') }}</a-button>
                  <a-button size="small" type="primary" @click="() => { onSearch(); advancedVisible = false; }">{{ t('common.search') }}</a-button>
                </div>
              </div>
            </template>
          </a-popover>
        </div>
      </div>

      <!-- PC 表格（576px 以上显示） -->
      <div class="cl-table-wrap">
        <a-skeleton v-if="tableLoading" :animation="true">
          <a-skeleton-line :rows="8" />
        </a-skeleton>
        <a-table
          v-else
          :bordered="false"
          :columns="tableColumns"
          :data="tableList"
          row-key="id"
          :scroll="{ x: 900 }"
          :row-selection="tableRowSelection"
          v-model:selectedKeys="selectedTableIds"
          :pagination="false"
        >
          <template #actions="{ record }">
            <a-space size="mini">
              <a-button v-permission="'gen:table$preview'" size="mini" type="text" :title="t('tool.gen.preview')" @click.stop="openPreview(record)">
                <template #icon><IconEye /></template>{{ t('tool.gen.preview') }}
              </a-button>
              <a-button v-permission="'gen:table$import'" size="mini" type="text" :title="t('tool.gen.sync')" @click.stop="syncTable(record)">
                <template #icon><IconRefresh /></template>{{ t('tool.gen.sync') }}
              </a-button>
              <a-button v-permission="'gen:table$download'" size="mini" type="text" :title="t('tool.gen.download')" @click.stop="handleDownload([record.id])">
                <template #icon><IconDownload /></template>{{ t('tool.gen.download') }}
              </a-button>
              <a-button v-permission="'gen:table$design'" size="mini" type="text" :title="t('tool.gen.fieldDesign')" @click.stop="openDesignDrawer(record)">
                <template #icon><SvgIcon name="data-table" /></template>{{ t('tool.gen.fieldDesign') }}
              </a-button>
              <a-button v-permission="'gen:table$edit'" size="mini" type="text" @click.stop="openTableEdit(record)">
                <template #icon><IconEdit /></template>{{ t('common.edit') }}
              </a-button>
              <a-button v-permission="'gen:table$delete'" size="mini" type="text" status="danger" @click.stop="confirmDeleteTable(record)">
                <template #icon><IconDelete /></template>{{ t('common.delete') }}
              </a-button>
            </a-space>
          </template>
        </a-table>
        <!-- PC 分页 -->
        <div class="cl-pagination">
          <a-pagination
            v-model:current="current"
            :total="tableTotal"
            :page-size="pageSize"
            show-total
            show-page-size
            @change="loadTablePage"
            @page-size-change="(s: number) => { pageSize = s; current = 1; loadTablePage(); }"
          />
        </div>
      </div>

      <!-- 移动端卡片列表（576px 以下显示） -->
      <div class="cl-card-list">
        <!-- 移动端选择模式工具栏 -->
        <div class="cl-mobile-select-bar">
          <a-button
            size="small"
            :type="mobileSelectMode ? 'primary' : 'secondary'"
            @click="toggleMobileSelectMode"
          >
            <template #icon>
              <IconClose v-if="mobileSelectMode" />
              <IconCheck v-else />
            </template>
            {{ mobileSelectMode ? t('common.cancel') : t('common.select') }}
          </a-button>
          <template v-if="mobileSelectMode">
            <a-checkbox
              :model-value="mobileAllSelected"
              :indeterminate="mobileIndeterminate"
              @change="toggleMobileSelectAll"
            >{{ t('common.selectAll') }}</a-checkbox>
            <span class="cl-select-count">{{ t('common.selected', { count: selectedTableIds.length }) }}</span>
          </template>
        </div>

        <a-skeleton v-if="tableLoading" :animation="true">
          <a-skeleton-line :rows="6" />
        </a-skeleton>
        <template v-else>
          <div
            v-for="(item, index) in mobileData"
            :key="item.id"
            class="cl-card stagger-item"
            :class="{ 'cl-card--selected': mobileSelectMode && isCardSelected(item.id) }"
            :style="{ '--stagger-index': index % mobilePageSize }"
            @click="mobileSelectMode ? toggleCardSelect(item.id) : undefined"
          >
            <div class="cl-card-header">
              <!-- 选择模式下显示复选框 -->
              <a-checkbox
                v-if="mobileSelectMode"
                :model-value="isCardSelected(item.id)"
                @click.stop
                @change="toggleCardSelect(item.id)"
                class="cl-card-checkbox"
              />
              <!-- <div v-else class="cl-card-icon">
                <IconPlus style="font-size:18px;color:var(--color-text-3)" />
              </div> -->
              <div class="cl-card-identity">
                <strong>{{ item.tableName }}</strong>
                <span class="cl-card-sub">{{ item.tableComment || t('tool.gen.noComment') }}</span>
              </div>
              <a-tag color="arcoblue" size="small">{{ item.moduleName || '-' }}</a-tag>
            </div>
            <div class="cl-card-meta">
              <span>{{ t('tool.gen.cardClassName') }}：{{ item.className || '-' }}</span>
              <span>{{ t('tool.gen.cardAuthor') }}：{{ item.author || '-' }}</span>
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'gen:table$design'" size="mini" type="outline" @click="openDesignDrawer(item)">
                  <template #icon><IconSettings /></template>{{ t('tool.gen.fieldDesign') }}
                </a-button>
                <a-button v-permission="'gen:table$edit'" size="mini" type="outline" @click="openTableEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button v-permission="'gen:table$delete'" size="mini" type="outline" status="danger" @click="confirmDeleteTable(item)">
                  <template #icon><IconDelete /></template>{{ t('common.delete') }}
                </a-button>
              </a-space>
            </div>
          </div>

          <!-- 加载更多骨架 -->
          <template v-if="mobileLoadingMore">
            <div v-for="n in 3" :key="`sk-${n}`" class="cl-card cl-card-skeleton">
              <div class="cl-card-header">
                <div class="skeleton cl-sk-avatar" />
                <div class="cl-card-identity">
                  <div class="skeleton cl-sk-name" />
                  <div class="skeleton cl-sk-sub" />
                </div>
              </div>
              <div class="cl-card-meta">
                <div class="skeleton cl-sk-meta" />
                <div class="skeleton cl-sk-meta" />
              </div>
            </div>
          </template>

          <a-empty v-if="!mobileData.length && !mobileLoadingMore" :description="t('common.noData')" />
          <div ref="sentinelRef" class="cl-sentinel">
            <span v-if="!mobileLoadingMore && !mobileHasMore && mobileData.length > 0" class="cl-no-more">
              {{ t('common.noMore') }}
            </span>
          </div>
        </template>
      </div>

    </a-card>

    <!-- ── 代码预览弹窗 ── -->
    <a-modal v-model:visible="previewModalVisible" :title="t('tool.gen.previewTitle', { name: previewTableName })"
      :width="isMobile ? '100%' : '80%'"
      :fullscreen="isMobile"
      :footer="false"
      unmount-on-close
    >
      <a-spin :loading="previewLoading" style="width:100%;min-height:400px">
        <div v-if="!previewLoading && previewFiles.length" class="preview-layout">
          <!-- 左侧文件列表 -->
          <div class="preview-sidebar">
            <div
              v-for="f in previewFiles"
              :key="f.name"
              class="preview-file-item"
              :class="{ active: activePreviewFile === f.name }"
              @click="activePreviewFile = f.name"
            >
              <span class="preview-file-type">{{ f.fileType }}</span>
              {{ f.name }}
            </div>
          </div>
          <!-- 右侧代码内容 -->
          <div class="preview-content">
            <div class="preview-content-header">
              <span class="preview-filename">{{ activePreviewFile }}</span>
              <a-button size="mini" type="text" @click="copyCode">{{ t('tool.gen.previewCopy') }}</a-button>
            </div>
            <div class="preview-code-wrap">
              <a-spin :loading="highlightLoading" style="width:100%;height:100%">
                <!-- shiki 输出的带样式 HTML -->
                <div v-if="highlightedHtml" class="preview-shiki" v-html="highlightedHtml" />
              </a-spin>
            </div>
          </div>
        </div>
        <a-empty v-else-if="!previewLoading" :description="t('tool.gen.previewNoData')" />
      </a-spin>
    </a-modal>

    <a-modal v-model:visible="importModalVisible" :title="t('tool.gen.importTitle')"
      :width="isMobile ? '100%' : '40%'" :fullscreen="isMobile"
      :ok-text="t('tool.gen.importOkText')" :ok-loading="importLoading" @ok="submitImport">
      <div style="margin-bottom:12px;display:flex;align-items:center;gap:12px;flex-wrap:wrap">
        <a-input-search v-model="importSearchKeyword" :placeholder="t('tool.gen.importSearchPlaceholder')"
          allow-clear style="width:260px" />
        <span style="color:var(--color-text-3);font-size:13px">
          {{ t('tool.gen.importSelected', { selected: importSelectedKeys.length, total: filteredDBTables.length }) }}
        </span>
      </div>
      <a-table
        row-key="tableName"
        :columns="importColumns"
        :data="filteredDBTables"
        :loading="importLoading"
        :row-selection="importRowSelection"
        v-model:selectedKeys="importSelectedKeys"
        :pagination="false"
        :scroll="{ y: 360 }"
        size="small"
      />
    </a-modal>

    <a-modal v-model:visible="tableModalVisible" :title="t('tool.gen.editTitle')"
      :width="isMobile ? '100%' : 560" :fullscreen="isMobile"
      :ok-text="t('tool.gen.editOkText')" @ok="submitTableEdit">
      <a-form :model="tableForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.gen.fieldTableName')" required>
              <a-input v-model="tableForm.tableName" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.gen.fieldTableComment')">
              <a-input v-model="tableForm.tableComment" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.gen.fieldClassName')">
              <a-input v-model="tableForm.className" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.gen.fieldModuleName')">
              <a-input v-model="tableForm.moduleName" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.gen.fieldRouter')">
              <a-input v-model="tableForm.router" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('tool.gen.fieldPermission')">
              <a-input v-model="tableForm.permission" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('tool.gen.fieldParentMenu')">
              <a-tree-select v-model="tableForm.menuId" :data="menuTreeOptions"
                :placeholder="t('tool.gen.fieldParentMenuPlaceholder')"
                allow-clear :fallback-option="false" style="width:100%" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('tool.gen.fieldAuthor')">
              <a-input v-model="tableForm.author" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('tool.gen.fieldRemark')">
              <a-textarea v-model="tableForm.remark" :max-length="200" show-word-limit />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <a-drawer v-model:visible="designDrawerVisible" :title="t('tool.gen.designTitle', { name: designTableName })"
      :placement="isMobile ? 'bottom' : 'right'"
      :width="isMobile ? '100%' : '80%'"
      :height="isMobile ? '85%' : undefined"
      :footer="true"
    >
      <template #footer>
        <a-space>
          <a-button @click="designDrawerVisible = false">{{ t('common.cancel') }}</a-button>
          <a-button type="primary" :loading="designSaving" @click="saveDesign">{{ t('common.save') }}</a-button>
        </a-space>
      </template>

      <a-spin :loading="designLoading" style="width:100%">
        <a-table :data="designColumns" :pagination="false" row-key="id" size="small"
          :draggable="{ type: 'handle', width: 40 }" :scroll="{ x: 1100 }">
          <template #columns>
            <a-table-column :title="t('tool.gen.designColDbField')" data-index="columnName" :width="130" />
            <a-table-column :title="t('tool.gen.designColDbType')" data-index="columnType" :width="150" />
            <a-table-column :title="t('tool.gen.designColComment')" :width="150">
              <template #cell="{ record }"><a-input v-model="record.columnComment" size="small" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColFieldName')" :width="150">
              <template #cell="{ record }"><a-input v-model="record.fieldName" size="small" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColDataType')" :width="150">
              <template #cell="{ record }"><a-input v-model="record.dataType" size="small" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColJson')" :width="120" align="center">
              <template #cell="{ record }"><a-input v-model="record.json" size="small" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColRequired')" :width="60" align="center">
              <template #cell="{ record }"><a-checkbox v-model="record.isRequired" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColInsert')" :width="60" align="center">
              <template #cell="{ record }"><a-checkbox v-model="record.isInsert" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColEdit')" :width="60" align="center">
              <template #cell="{ record }"><a-checkbox v-model="record.isEdit" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColList')" :width="60" align="center">
              <template #cell="{ record }"><a-checkbox v-model="record.isList" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColSortable')" :width="60" align="center">
              <template #cell="{ record }"><a-checkbox v-model="record.sortable" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColCondition')" :width="120" align="center">
              <template #cell="{ record }">
                <a-select v-model="record.condition" size="small" allow-clear style="width:100%">
                  <a-option v-for="q in dicts.sys_condition" :key="q.dictValue" :value="q.dictValue">{{ q.label }}</a-option>
                </a-select>
              </template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColComponent')" :width="180" align="center">
              <template #cell="{ record }">
                <a-select v-model="record.component" size="small" allow-clear style="width:100%">
                  <a-option v-for="c in dicts.xcode_component" :key="c.dictValue" :value="c.dictValue">{{ c.label }}</a-option>
                </a-select>
              </template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColDict')" :width="180" align="center">
              <template #cell="{ record }">
                <a-select v-model="record.dictType" size="small" allow-clear style="width:100%" allow-search>
                  <a-option v-for="d in dictList" :key="d.dictType" :value="d.dictType">{{ d.dictName }}</a-option>
                </a-select>
              </template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColDefault')" :width="100" align="center">
              <template #cell="{ record }"><a-input v-model="record.defaultValue" size="small" /></template>
            </a-table-column>
            <a-table-column :title="t('tool.gen.designColSort')" :width="80" align="center">
              <template #cell="{ record }">
                <a-input-number v-model="record.sort" size="small" :min="0" style="width:100%" />
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-spin>
    </a-drawer>
  </div>
</template>
