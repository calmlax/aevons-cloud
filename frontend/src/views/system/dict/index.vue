<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter,IconSettings,  IconPlus, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n'
import dictApi from "@/api/system/dict"
import dictData from './dict-data.vue';

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any

const dicts = proxy.$useDict('sys_is', 'sys_common_status')
// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

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
  dictType: '',
  dictName: '',
  status: '',
  isSys: '',
  direction: 'desc',field:'id'
})

async function loadPage() {
  loading.value = true;
  try {
    let res:any = await dictApi.page(queryParams.value);
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
    
  queryParams.value.dictType = ''
  queryParams.value.dictName = ''
  queryParams.value.status = ''
  queryParams.value.isSys = ''
  handleSearch();
}

watch(
  () => ({ ...queryParams.value }),
  () => {
    filterModule.value=(queryParams.value.dictName != '' || queryParams.value.status != '' ||  queryParams.value.isSys != '')
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

// ── 新增 / 修改 共用弹窗 ────────────────────────────────
const modalVisible = ref(false);
const formRef = ref()
const form = reactive<any>({});
const isEdit = ref(false); // 标记是否为编辑模式
const submitLoading = ref(false)
// 打开新增
function handleAdd() {
  isEdit.value = false;
  Object.assign(form, {
    dictType:'',
    dictName:'',
    status:0,
    isSys:0,
    remark:'',
  });
  modalVisible.value = true;
}


// 打开编辑
async function handleEdit(record:any) {
  isEdit.value = true;
  let res = await dictApi.getById(record.id)
  Object.assign(form, res);
  modalVisible.value = true;
}

// 统一提交
async function submitForm() {
  submitLoading.value = true
  try {
    const valid = await formRef.value.validate()
    if (valid) return false
    if (isEdit.value) {
      await dictApi.update(form.id, form);
      Message.success(t('common.editSuccess'));
    } else {
      await dictApi.add(form);
      Message.success(t('common.addSuccess'));
    }
    modalVisible.value = false;
    loadPage();
  } finally {
    submitLoading.value = false
  }
}

// 删除单条
async function handleDelete(record: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.dict.confirmDelete', { type: record.dictType }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await dictApi.delete([record.id]);
      Message.success(t('common.deleteSuccess'));
      loadPage();
    },
  });
}

async function handleRefreshCache() {
  Modal.confirm({
    title: t('common.hint'),
    content: t('system.dict.refreshCacheConfirm'),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await dictApi.refreshCache();
      Message.success(t('system.dict.refreshCacheSuccess'));
    },
  });
}

async function handleDesignDrawer(record: any) {
  proxy.$refs['dictDataRef'].open(record, '字典数据-'+record.dictType)
}
loadPage();
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'sys:dict$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            {{ t('common.add') }}
          </a-button>
          <a-button v-permission="'sys:dict$refresh'" type="primary" status="warning" @click="handleRefreshCache">
            {{ t('system.dict.refreshCache') }}
          </a-button>

          <!-- <a-button
            v-permission="'sys:dict$delete'"
            status="danger"
            :disabled="!selectedIds.length"
            @click="batchDelete"
          >
            <template #icon><IconDelete /></template>
            批量删除
          </a-button> -->
        </a-space>

        <!-- 搜索 -->
        <div class="cl-toolbar-right">
          <a-input-search
            v-model="queryParams.dictType"
            :placeholder="t('system.dict.searchPlaceholder')"
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
                <a-input v-model="queryParams.dictName" :placeholder="t('system.dict.dictNameFilterPlaceholder')" allow-clear />
                <a-select v-model="queryParams.status" :placeholder="t('common.status')" allow-clear >
                  <a-option v-for="dict in dicts.sys_common_status" :key="dict.dictValue" :value="dict.dictValue">
                    {{ dict.label }}
                  </a-option>
                </a-select>
                <a-select v-model="queryParams.isSys" :placeholder="t('system.dict.isSysFilterPlaceholder')" allow-clear >
                  <a-option v-for="dict in dicts.sys_is" :key="dict.dictValue" :value="dict.dictValue">
                    {{ dict.label }}
                  </a-option>
                </a-select>
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
            <a-table-column :title="t('system.dict.dictType')" data-index="dictType" :ellipsis="true" :tooltip="true">
              <template #cell="{ record }">
                <code class="mm-code mm-code-sign">{{ record.dictType }}</code>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.dict.dictName')" data-index="dictName" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('common.status')" :width="100" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_common_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('system.dict.isSys')" :width="100" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_is" :value="record.isSys"/>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.dict.createdAt')" data-index="createdAt" :width="230" :ellipsis="true" :tooltip="true" align="center"/>
            <a-table-column :title="t('common.action')" :width="260" align="center">
              <template #cell="{ record }">
                <a-space size="mini">
                <a-button v-permission="'sys:dict$design'" size="mini" type="text" @click.stop="handleDesignDrawer(record)">
                  <template #icon><SvgIcon name="dict" /></template>{{ t('system.dict.configure') }}
                </a-button>
                  <a-button v-permission="'sys:dict$edit'"
                    size="mini"
                    type="text"
                    @click.stop="handleEdit(record)"
                  >
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button
                    v-permission="'sys:dict$delete'"
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
                <strong>{{ item.dictType }}</strong>
                <span class="cl-card-sub">{{ item.dictName || t('common.noData') }}</span>
              </div>
              <DictTag :options="dicts.sys_common_status" :value="item.status" dot />
            </div>
            <div class="cl-card-meta">
              <span>{{ t('system.dict.isSys') }}：{{ item.isSys }}</span>
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'sys:dict$design'" size="mini" type="outline" @click.stop="handleDesignDrawer(item)">
                  <template #icon><IconSettings /></template>{{ t('system.dict.configure') }}
                </a-button>
                <a-button v-permission="'sys:dict$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button
                  v-permission="'sys:dict$delete'"
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
    :title="isEdit ? t('system.dict.editTitle') : t('system.dict.addTitle')"
      :width="isMobile ? '100%' : 560"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.dict.dictType')" field="dictType" 
            :rules="[{ required: true, message: t('system.dict.dictTypeRequired') },{ maxLength: 32, message: t('common.maxLength', { max: 32 }) },]" :validate-trigger="['blur']">
              <a-input v-model="form.dictType" :placeholder="t('system.dict.dictTypePlaceholder')" />
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.dict.dictName')" field="dictName" 
            :rules="[{ required: true, message: t('system.dict.dictNameRequired') },{ maxLength: 50, message: t('common.maxLength', { max: 50 }) },]" :validate-trigger="['blur']">
              <a-input v-model="form.dictName" :placeholder="t('system.dict.dictNamePlaceholder')" />
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.status')" field="status" 
            :rules="[{ required: true, message: t('system.dict.statusRequired') }]" :validate-trigger="['blur']">
              <a-radio-group v-model="form.status" type="button">
                <a-radio v-for="dict in dicts.sys_common_status" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.dict.isSys')" field="isSys" 
            :rules="[{ required: true, message: t('system.dict.isSysRequired') }]" :validate-trigger="['blur']">
              <a-switch v-model="form.isSys" :checked-value="1" :unchecked-value="0"
                :checked-text="t('system.dict.isSysYes')" :unchecked-text="t('system.dict.isSysNo')" />
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('common.remark')" field="remark" 
            :rules="[{ maxLength: 500, message: t('common.maxLength', { max: 500 }) },]" :validate-trigger="['blur']">
              <a-textarea v-model="form.remark" :placeholder="t('system.dict.remarkPlaceholder')" allow-clear/>
            </a-form-item>
          </a-col>
        <a-col :span="24" style="margin-top: 30px;">
         <div style="display: flex; justify-content: center; gap: 12px;">
       <a-button @click="modalVisible = false">{{ t('common.cancel') }}</a-button>
      <a-button type="primary" @click="submitForm()" :loading="submitLoading">{{ t('common.save') }}</a-button>
    </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <dictData ref="dictDataRef" :isMobile="isMobile"></dictData>
  </div>
</template>
