<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import postApi from "@/api/system/post"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any
// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

const dicts = proxy.$useDict('sys_common_status')
// ── 搜索 & 筛选 ────────────────────────────────────────
const filterModule = ref(false);
const advancedVisible = ref(false);

// ── PC 分页 ────────────────────────────────────────────
const total = ref(0);
const loading = ref(false);
const dataList = ref<any[]>([]);
const queryParams = ref({
  pageNum:1,
  pageSize:20,
  postName: undefined,
  status: undefined,
})

async function loadPage() {
  loading.value = true;
  try {
    let res:any = await postApi.page(queryParams.value);
    dataList.value = res?.rows ?? [];
    total.value = res?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function handleSearch() {
  queryParams.value.pageNum = 1;
  loadPage();
}
function handleReset() {
    
  queryParams.value.postName = undefined
  queryParams.value.status = undefined
  handleSearch();
}
watch(
  () => ({ ...queryParams.value }),
  () => {
    filterModule.value=(queryParams.value.status !== '')
  },
  { deep: true }
);

// ── 移动端无限滚动 ─────────────────────────────────────
const mobilePageSize = 10;
const mobilePage = ref(mobilePageSize);
const mobileLoadingMore = ref(false);
const sentinelRef = ref<HTMLElement | null>(null);
let observer: IntersectionObserver | null = null;

const mobileData = computed(() => dataList.value.slice(0, mobilePage.value));
const mobileHasMore = computed(() => mobilePage.value < dataList.value.length);

function loadMore() {
  if (!mobileHasMore.value || mobileLoadingMore.value) return;
  mobileLoadingMore.value = true;
  setTimeout(() => {
    mobilePage.value += mobilePageSize;
    mobileLoadingMore.value = false;
  }, 300);
}

function setupObserver() {
  if (!sentinelRef.value) return;
  observer?.disconnect();
  observer = new IntersectionObserver(entries => {
    if (entries[0].isIntersecting) loadMore();
  }, { rootMargin: '80px' });
  observer.observe(sentinelRef.value);
}

onMounted(() => setupObserver());
onUnmounted(() => observer?.disconnect());
watch(dataList, async () => {
  mobilePage.value = mobilePageSize;
  await nextTick();
  setupObserver();
});

// ── 多选 ───────────────────────────────────────────────
const selectedIds = ref<string[]>([]);
const rowSelection = reactive({ type: 'checkbox', showCheckedAll: true, onlyCurrent: false });

// ── 移动端多选模式 ─────────────────────────────────────
const mobileSelectMode = ref(false);

function toggleMobileSelectMode() {
  mobileSelectMode.value = !mobileSelectMode.value;
  if (!mobileSelectMode.value) {
    selectedIds.value = [];
  }
}

function toggleCardSelect(id: string) {
  const idx = selectedIds.value.indexOf(id);
  idx === -1 ? selectedIds.value.push(id) : selectedIds.value.splice(idx, 1);
}

function isCardSelected(id: string) {
  return selectedIds.value.includes(id);
}

const mobileAllSelected = computed(() =>
  mobileData.value.length > 0 && mobileData.value.every((item:any) => selectedIds.value.includes(item.id))
);
const mobileIndeterminate = computed(() =>
  selectedIds.value.length > 0 && !mobileAllSelected.value
);

function toggleMobileSelectAll() {
  const currentIds = mobileData.value.map((item:any) => item.id);
  if (mobileAllSelected.value) {
    selectedIds.value = selectedIds.value.filter(id => !currentIds.includes(id));
  } else {
    selectedIds.value.push(...currentIds.filter(id => !selectedIds.value.includes(id)));
  }
}

async function batchDelete() {
  if (!selectedIds.value.length) {
    Message.warning(t('common.confirmDelete'));
    return;
  }
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('common.confirmDeleteContent', { count: selectedIds.value.length }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await postApi.delete(selectedIds.value);
      Message.success(t('common.deleteSuccess'));
      selectedIds.value = [];
      loadPage();
    },
  });
}

// ── 新增 / 修改 共用弹窗 ────────────────────────────────
const modalVisible = ref(false);
const formRef = ref()
const form = reactive<any>({});
const isEdit = ref(false); // 标记是否为编辑模式
const submitLoading = ref(false)

// 打开新增
function handleAdd() {
  isEdit.value = false;
  resetForm('');
  modalVisible.value = true;
}

// 打开编辑
function handleEdit(row:any) {
  isEdit.value = true;
  resetForm(row.id);
  modalVisible.value = true;
}

async function resetForm(id:string){
  await nextTick()
  formRef.value?.clearValidate()
  if(isEdit.value){
    let res = await postApi.getById(id)
    Object.assign(form, res);
  }else{
    Object.assign(form, {
      postKey: '',
      postName: '',
      sort: 0,
      status: 0,
      remark: '',
    });
  }
}
// 统一提交
async function submitForm() {
  try {
    const valid = await formRef.value.validate()
    if (valid) return false
    submitLoading.value = true
    if (isEdit.value) {
      await postApi.update(form.id!, form);
      Message.success(t('common.editSuccess'));
    } else {
      await postApi.add(form);
      Message.success(t('common.addSuccess'));
    }
    modalVisible.value = false;
    loadPage();
  } finally {
    submitLoading.value = false
  }
}

// 删除单条
async function handleDelete(row: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.post.confirmDelete', { name: row.postName }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await postApi.delete([row.id]);
      Message.success(t('common.deleteSuccess'));
      loadPage();
    },
  });
}

loadPage();
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'sys:post$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            {{ t('common.add') }}
          </a-button>

          <a-button
            v-permission="'sys:post$delete'"
            status="danger"
            :disabled="!selectedIds.length"
            @click="batchDelete"
          >
            <template #icon><IconDelete /></template>
            {{ t('common.batchDelete') }}
          </a-button>
        </a-space>

        <!-- 搜索 -->
        <div class="cl-toolbar-right">
          <a-input-search
            v-model="queryParams.postName"
            :placeholder="t('system.post.searchPlaceholder')"
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
                <p class="cl-filter-title">{{ t('common.advancedFilter') }}</p>
                <a-input v-model="queryParams.status" :placeholder="t('system.post.statusFilterPlaceholder')" allow-clear />
                <div class="cl-filter-actions">
                  <a-button size="small" @click="handleReset">{{ t('common.reset') }}</a-button>
                  <a-button size="small" type="primary" @click="() => { handleSearch(); advancedVisible = false }">
                    {{ t('common.search') }}
                  </a-button>
                </div>
              </div>
            </template>
          </a-popover>
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
          :data="dataList"
          row-key="id"
          :scroll="{ x: 900 }"
          :row-selection="rowSelection"
          v-model:selectedKeys="selectedIds"
          :pagination="{
            current: queryParams.pageNum,
            pageSize: queryParams.pageSize,
            total: total,
            showTotal: true,
            showPageSize: true,
            showJumper: true,
            pageSizeOptions: [10, 20, 50, 100],
            }"        
          @page-change="(page:number) => {
            queryParams.pageNum = page
            loadPage()
          }"
          @page-size-change="(size:number) => {
            queryParams.pageSize = size
            queryParams.pageNum = 1
            loadPage()
          }"
        >
          <template #columns>
          <a-table-column :title="t('system.post.postKey')" data-index="postKey" />
          <a-table-column :title="t('system.post.postName')" data-index="postName" />
          <a-table-column :title="t('common.sort')" data-index="sort" align="center" :width="80"/>
          <a-table-column :title="t('common.status')" data-index="status" align="center" :width="100">
            <template #cell="{ record }">
              <DictTag :options="dicts.sys_common_status" :value="record.status" dot />
            </template>
          </a-table-column>
          <a-table-column :title="t('common.remark')" data-index="remark" :ellipsis="true" :tooltip="true" />
          <a-table-column :title="t('common.action')" :width="260" align="center">
            <template #cell="{ record }">
              <a-space size="mini">
                <a-button
                  v-permission="'sys:post$edit'"
                  size="mini"
                  type="text"
                  @click.stop="handleEdit(record)"
                >
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button
                  v-permission="'sys:post$delete'"
                  size="mini"
                  type="text"
                  status="danger"
                  @click.stop="handleDelete(record)"
                >
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
            >
              {{ t('common.selectAll') }}
            </a-checkbox>
            <span class="cl-select-count">{{ t('common.selected', { count: selectedIds.length }) }}</span>
          </template>
        </div>

        <a-skeleton v-if="loading" :animation="true">
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
              <a-checkbox
                v-if="mobileSelectMode"
                :model-value="isCardSelected(item.id)"
                @click.stop
                @change="toggleCardSelect(item.id)"
                class="cl-card-checkbox"
              />
              <div class="cl-card-identity">
                <strong>{{ item.postKey }}</strong>
                <span class="cl-card-sub">{{ item.postName || '-' }}</span>
              </div>
              <a-tag color="arcoblue" size="small">{{ item.postKey || '-' }}</a-tag>
            </div>
            <div class="cl-card-meta">
              <span>{{ t('common.sort') }}：{{ item.sort ?? '-' }}</span>
              <span>{{ t('common.remark') }}：{{ item.remark || '-' }}</span>
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'sys:post$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button
                  v-permission="'sys:post$delete'"
                  size="mini"
                  type="outline"
                  status="danger"
                  @click="handleDelete(item)"
                >
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

    <!-- ── 新增 / 修改 共用弹窗 ── -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? t('system.post.editTitle') : t('system.post.addTitle')"
      :width="isMobile ? '100%' : 520"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.post.postKey')" field="postKey"
              :rules="[{ required: true, message: t('system.post.postKeyRequired') }, { maxLength: 32, message: t('common.maxLength', { max: 32 }) }]"
              :validate-trigger="['blur']">
              <a-input v-model="form.postKey" :placeholder="t('system.post.postKeyPlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.post.postName')" field="postName"
              :rules="[{ required: true, message: t('system.post.postNameRequired') }, { maxLength: 50, message: t('common.maxLength', { max: 50 }) }]"
              :validate-trigger="['blur']">
              <a-input v-model="form.postName" :placeholder="t('system.post.postNamePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.sort')" field="sort">
              <a-input-number v-model="form.sort" :min="0" style="width:100%" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.status')" field="status">
              <a-switch v-model="form.status" :checked-value="0" :unchecked-value="1" :checked-text="t('system.post.statusNormal')" :unchecked-text="t('system.post.statusDisabled')" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('common.remark')" field="remark"
              :rules="[{ maxLength: 500, message: t('common.maxLength', { max: 500 }) }]"
              :validate-trigger="['blur']">
              <a-textarea v-model="form.remark" :placeholder="t('system.post.remarkPlaceholder')" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="24" style="margin-top: 16px;">
            <div style="display: flex; justify-content: center; gap: 12px;">
              <a-button @click="modalVisible = false">{{ t('common.cancel') }}</a-button>
              <a-button type="primary" @click="submitForm()" :loading="submitLoading">{{ t('common.save') }}</a-button>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

  </div>
</template>