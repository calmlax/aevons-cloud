<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n';
import noticeApi from "@/api/system/notice"

const { t } = useI18n();
const proxy = getCurrentInstance()!.proxy as any

// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

const dicts = proxy.$useDict('sys_notice_type', 'sys_common_status')

// ── 搜索 & 筛选 ────────────────────────────────────────
const filterModule = ref(false);
const advancedVisible = ref(false);

// ── 分页 & 数据（PC / 移动端共用查询参数）─────────────
const total = ref(0);
const loading = ref(false);
const dataList = ref<any[]>([]);
const queryParams = ref({
  pageNum: 1,
  pageSize: 10,
  title: '',
  type: '',
  status: '',
})

// 移动端滚动分页：追加模式
const mobileList = ref<any[]>([]);   // 累积所有已加载的卡片
const mobileLoadingMore = ref(false);
const sentinelRef = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const mobileHasMore = computed(() => mobileList.value.length < total.value);

// 加载一页数据
async function loadPage() {
  loading.value = true;
  try {
    const res: any = await noticeApi.page(queryParams.value);
    dataList.value = res?.rows ?? [];
    total.value = res?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

// 移动端：加载下一页并追加
async function loadMoreMobile() {
  if (!mobileHasMore.value || mobileLoadingMore.value) return;
  mobileLoadingMore.value = true;
  try {
    queryParams.value.pageNum += 1;
    const res: any = await noticeApi.page(queryParams.value);
    mobileList.value.push(...(res?.rows ?? []));
    total.value = res?.total ?? 0;
  } finally {
    mobileLoadingMore.value = false;
  }
}

// 初始搜索（重置到第1页）
async function handleSearch() {
  queryParams.value.pageNum = 1;
  if (isMobile.value) {
    // 移动端：清空累积列表，重新加载第1页
    loading.value = true;
    try {
      const res: any = await noticeApi.page(queryParams.value);
      mobileList.value = res?.rows ?? [];
      total.value = res?.total ?? 0;
    } finally {
      loading.value = false;
    }
  } else {
    loadPage();
  }
}

function handleReset() {
  queryParams.value.title = '';
  queryParams.value.type = '';
  queryParams.value.status = '';
  handleSearch();
}

watch(
  () => ({ ...queryParams.value }),
  () => {
    filterModule.value = (
      queryParams.value.title !== '' ||
      queryParams.value.type !== '' ||
      queryParams.value.status !== ''
    );
  },
  { deep: true }
);

// ── IntersectionObserver（移动端滚动加载）─────────────
function setupObserver() {
  observer?.disconnect();
  if (!sentinelRef.value) return;
  observer = new IntersectionObserver(
    (entries) => { if (entries[0].isIntersecting) loadMoreMobile(); },
    { rootMargin: '80px' }
  );
  observer.observe(sentinelRef.value);
}

onMounted(async () => {
  // 初始加载
  if (isMobile.value) {
    await handleSearch();
    await nextTick();
    setupObserver();
  } else {
    loadPage();
  }
});
onUnmounted(() => observer?.disconnect());

// 切换到移动端时重建 observer
watch(isMobile, async (val) => {
  if (val) {
    await handleSearch();
    await nextTick();
    setupObserver();
  } else {
    observer?.disconnect();
  }
});

// ── 多选 ───────────────────────────────────────────────
const selectedIds = ref<string[]>([]);
const mobileSelectMode = ref(false);

function toggleMobileSelectMode() {
  mobileSelectMode.value = !mobileSelectMode.value;
  if (!mobileSelectMode.value) selectedIds.value = [];
}

function toggleCardSelect(id: string) {
  const idx = selectedIds.value.indexOf(id);
  idx === -1 ? selectedIds.value.push(id) : selectedIds.value.splice(idx, 1);
}

function isCardSelected(id: string) {
  return selectedIds.value.includes(id);
}

// 当前视图数据（PC 用 dataList，移动端用 mobileList）
const currentList = computed(() => isMobile.value ? mobileList.value : dataList.value);

const allSelected = computed(() =>
  currentList.value.length > 0 &&
  currentList.value.every((item: any) => selectedIds.value.includes(item.id))
);
const indeterminate = computed(() =>
  selectedIds.value.length > 0 && !allSelected.value
);

function toggleSelectAll() {
  const ids = currentList.value.map((item: any) => item.id);
  if (allSelected.value) {
    selectedIds.value = selectedIds.value.filter(id => !ids.includes(id));
  } else {
    selectedIds.value.push(...ids.filter(id => !selectedIds.value.includes(id)));
  }
}

async function batchDelete() {
  if (!selectedIds.value.length) {
    Message.warning(t('system.notice.selectFirst'));
    return;
  }
  Modal.confirm({
    title: t('system.notice.confirmDeleteTitle'),
    content: t('system.notice.confirmDeleteContent', { count: selectedIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await noticeApi.delete(selectedIds.value);
      Message.success(t('system.notice.deleteSuccess'));
      selectedIds.value = [];
      handleSearch();
    },
  });
}

// ── 新增 / 修改 共用弹窗 ────────────────────────────────
const modalVisible = ref(false);
const formRef = ref();
const form = reactive<any>({});
const isEdit = ref(false);
const submitLoading = ref(false);

function handleAdd() {
  isEdit.value = false;
  resetForm('');
  modalVisible.value = true;
}

function handleEdit(row: any) {
  isEdit.value = true;
  resetForm(row.id);
  modalVisible.value = true;
}

async function resetForm(id:string) {
  formRef.value?.clearValidate();
  if (isEdit.value) {
    const res = await noticeApi.getById(id);
    Object.assign(form, res);
  } else {
    Object.assign(form, {
      title: '',
      type: 1,
      content: '',
      status: 0,
      remark: '',
    });
  }
}

async function submitForm() {
  try {
    const valid = await formRef.value.validate();
    if (valid) return;
    submitLoading.value = true;
    if (isEdit.value) {
      await noticeApi.update(form.id, form);
      Message.success(t('system.notice.editSuccess'));
    } else {
      await noticeApi.add(form);
      Message.success(t('system.notice.addSuccess'));
    }
    modalVisible.value = false;
    handleSearch();
  } finally {
    submitLoading.value = false;
  }
}

async function handleDelete(row: any) {
  Modal.confirm({
    title: t('system.notice.confirmDeleteTitle'),
    content: t('system.notice.confirmDeleteSingle', { title: row.title }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await noticeApi.delete([row.id]);
      Message.success(t('system.notice.deleteSuccess'));
      handleSearch();
    },
  });
}
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'sys:notice$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            {{ t('system.notice.add') }}
          </a-button>
          <a-button
            v-permission="'sys:notice$delete'"
            status="danger"
            :disabled="!selectedIds.length"
            @click="batchDelete"
          >
            <template #icon><IconDelete /></template>
            {{ t('system.notice.batchDelete') }}
          </a-button>
        </a-space>

        <div class="cl-toolbar-right">
          <a-input-search
            v-model="queryParams.title"
            :placeholder="t('system.notice.searchPlaceholder')"
            allow-clear
            class="cl-toolbar-search"
            @search="handleSearch"
            @press-enter="handleSearch"
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
                <p class="cl-filter-title">{{ t('system.notice.advancedFilter') }}</p>
                <a-select v-model="queryParams.type" :placeholder="t('system.notice.noticeType')" allow-clear>
                  <a-option v-for="dict in dicts.sys_notice_type" :key="dict.dictValue" :value="Number(dict.dictValue)">
                    {{ dict.label }}
                  </a-option>
                </a-select>
                <a-select v-model="queryParams.status" :placeholder="t('system.notice.status')" allow-clear>
                  <a-option v-for="dict in dicts.sys_common_status" :key="dict.dictValue" :value="Number(dict.dictValue)">
                    {{ dict.label }}
                  </a-option>
                </a-select>
                <div class="cl-filter-actions">
                  <a-button size="small" @click="handleReset">{{ t('system.notice.reset') }}</a-button>
                  <a-button size="small" type="primary" @click="() => { handleSearch(); advancedVisible = false }">
                    {{ t('system.notice.search') }}
                  </a-button>
                </div>
              </div>
            </template>
          </a-popover>
        </div>
      </div>

      <!-- 多选工具栏（始终显示） -->
      <div class="cl-mobile-select-bar">
        <a-button size="small" :type="mobileSelectMode ? 'primary' : 'secondary'" @click="toggleMobileSelectMode">
          <template #icon>
            <IconClose v-if="mobileSelectMode" />
            <IconCheck v-else />
          </template>
          {{ mobileSelectMode ? t('system.notice.cancel') : t('system.notice.select') }}
        </a-button>
        <template v-if="mobileSelectMode">
          <a-checkbox
            :model-value="allSelected"
            :indeterminate="indeterminate"
            @change="toggleSelectAll"
          >
            {{ t('system.notice.selectAll') }}
          </a-checkbox>
          <span class="cl-select-count">{{ t('system.notice.selectedCount', { count: selectedIds.length }) }}</span>
        </template>
      </div>

      <!-- 列表（分割线样式，PC + 移动端统一） -->
      <a-skeleton v-if="loading" :animation="true">
        <a-skeleton-line :rows="6" />
      </a-skeleton>
      <template v-else>
        <div class="notice-list">
          <template v-for="(item, index) in currentList" :key="item.id">
            <div
              class="notice-item stagger-item"
              :class="{ 'notice-item--selected': mobileSelectMode && isCardSelected(item.id) }"
              :style="{ '--stagger-index': index % queryParams.pageSize }"
              @click="mobileSelectMode ? toggleCardSelect(item.id) : undefined"
            >
              <!-- 标题行 -->
              <div class="notice-item-header">
                <a-checkbox
                  v-if="mobileSelectMode"
                  :model-value="isCardSelected(item.id)"
                  @click.stop
                  @change="toggleCardSelect(item.id)"
                  style="margin-right: 8px; flex-shrink: 0;"
                />
                <span class="notice-item-title">{{ item.title }}</span>
                <DictTag :options="dicts.sys_common_status" :value="item.status" dot style="margin-left: 8px; flex-shrink: 0;" />
              </div>
              <!-- 元信息行 -->
              <div class="notice-item-meta">
                <DictTag :options="dicts.sys_notice_type" :value="item.type" />
                <span v-if="item.remark" class="notice-item-remark">{{ item.remark }}</span>
                <span class="notice-item-time">{{ item.createdAt }}</span>
              </div>
              <!-- 内容摘要 -->
              <div v-if="item.content" class="notice-item-content">{{ item.content }}</div>
              <!-- 操作（右下） -->
              <div v-if="!mobileSelectMode" class="notice-item-actions">
                <a-button v-permission="'sys:notice$edit'" size="mini" type="text" @click.stop="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('system.notice.edit') }}
                </a-button>
                <a-button
                  v-permission="'sys:notice$delete'"
                  size="mini"
                  type="text"
                  status="danger"
                  @click.stop="handleDelete(item)"
                >
                  <template #icon><IconDelete /></template>{{ t('system.notice.delete') }}
                </a-button>
              </div>
            </div>
            <a-divider v-if="index < currentList.length - 1" :margin="0" />
          </template>
        </div>

        <a-empty v-if="!currentList.length" :description="t('system.notice.noData')" />

        <!-- 移动端：滚动加载哨兵 -->
        <div v-if="isMobile" ref="sentinelRef" class="cl-sentinel">
          <a-spin v-if="mobileLoadingMore" />
          <span v-else-if="!mobileHasMore && mobileList.length > 0" class="cl-no-more">{{ t('system.notice.allLoaded') }}</span>
        </div>

        <!-- PC：分页组件 -->
        <div v-if="!isMobile && total > 0" class="notice-pagination">
          <a-pagination
            :total="total"
            :current="queryParams.pageNum"
            :page-size="queryParams.pageSize"
            :page-size-options="[10, 20, 50, 100]"
            show-total
            show-page-size
            show-jumper
            @change="(page: number) => { queryParams.pageNum = page; loadPage() }"
            @page-size-change="(size: number) => { queryParams.pageSize = size; queryParams.pageNum = 1; loadPage() }"
          />
        </div>
      </template>
    </a-card>

    <!-- ── 新增 / 修改 共用弹窗 ── -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? t('system.notice.editTitle') : t('system.notice.addTitle')"
      :width="isMobile ? '100%' : 600"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item :label="t('system.notice.noticeTitle')" field="title"
              :rules="[
                { required: true, message: t('system.notice.noticeTitleRequired') }, 
                { maxLength: 50, message: t('system.notice.noticeTitleMaxLength') }
              ]"
              :validate-trigger="['blur']">
              <a-input v-model="form.title" :placeholder="t('system.notice.noticeTitlePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.notice.noticeTypeLabel')" field="type">
              <a-radio-group v-model="form.type" type="button">
                <a-radio v-for="dict in dicts.sys_notice_type" :key="dict.dictValue" :value="Number(dict.dictValue)">
                  {{ dict.label }}
                </a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('system.notice.noticeContent')" field="content"
              :rules="[
                { required: true, message: t('system.notice.noticeContentRequired') }, 
                { maxLength: 3000, message: t('system.notice.noticeContentMaxLength') }
              ]"
              :validate-trigger="['blur']">
              <a-textarea v-model="form.content" :placeholder="t('system.notice.noticeContentPlaceholder')" :auto-size="{ minRows: 4, maxRows: 10 }" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.notice.statusLabel')" field="status">
              <a-radio-group v-model="form.status" type="button">
                <a-radio v-for="dict in dicts.sys_common_status" :key="dict.dictValue" :value="Number(dict.dictValue)">
                  {{ dict.label }}
                </a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="24" style="margin-top: 16px;">
            <div style="display: flex; justify-content: center; gap: 12px;">
       <a-button @click="modalVisible = false">{{ t('common.cancel') }}</a-button>
              <a-button type="primary" @click="submitForm()" :loading="submitLoading">{{ t('system.notice.save') }}</a-button>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped>
.notice-list {
  margin-top: 4px;
}

.notice-item {
  padding: 12px 4px;
  cursor: default;
  transition: background 0.15s;
}

.notice-item:hover {
  background: var(--color-fill-1);
}

.notice-item--selected {
  background: var(--color-primary-light-1);
}

.notice-item-header {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 6px;
}

.notice-item-title {
  flex: 1;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1);
  line-height: 1.4;
}

.notice-item-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  color: var(--color-text-3);
  margin-bottom: 4px;
  flex-wrap: wrap;
}

.notice-item-remark {
  color: var(--color-text-3);
}

.notice-item-time {
  margin-left: auto;
}

.notice-item-content {
  font-size: 13px;
  color: var(--color-text-2);
  line-height: 1.6;
  margin: 4px 0 6px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notice-item-actions {
  display: flex;
  justify-content: flex-end;
  gap: 4px;
  margin-top: 6px;
}

/* PC 端始终显示多选栏 */
@media (min-width: 576px) {
  .cl-mobile-select-bar {
    display: flex !important;
    margin-bottom: 8px;
  }
}

.notice-pagination {
  display: flex;
  justify-content: center;
  padding: 16px 0 4px;
}
</style>
