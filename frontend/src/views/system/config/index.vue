<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import confApi from "@/api/system/conf"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any
// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

const dicts = proxy.$useDict('sys_is','sys_conf_scope')
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
  name: '',
  confKey: '',
  isSys: '',
  scope: '',
})

async function loadPage() {
  loading.value = true;
  try {
    let res:any = await confApi.page(queryParams.value);
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
    
  queryParams.value.name = ''
  queryParams.value.confKey = ''
  queryParams.value.isSys = ''
  queryParams.value.scope = ''
  handleSearch();
}
watch(
  () => ({ ...queryParams.value }),
  () => {
    filterModule.value=(queryParams.value.name !== ''||queryParams.value.isSys !== ''||queryParams.value.scope !== '')
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
    name:'',
    confKey:'',
    confValue:'',
    isSys:0,
    scope:0,
    isSecret:0,
    remark:'',
  });
  modalVisible.value = true;
}

// 打开编辑
async function handleEdit(row:any) {
  isEdit.value = true;
  let res = await confApi.getById(row.id)
  Object.assign(form, res);
  modalVisible.value = true;
}

// 统一提交
async function submitForm() {
  try {
    const valid = await formRef.value.validate()
    if (valid) return false
    submitLoading.value = true
    if (isEdit.value) {
      await confApi.update(form.id!, form);
      Message.success(t('common.editSuccess'));
    } else {
      await confApi.add(form);
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
    content: t('system.config.confirmDelete', { id: row.id }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await confApi.delete([row.id]);
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
          <a-button v-permission="'sys:conf$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            {{ t('common.add') }}
          </a-button>
        </a-space>

        <!-- 搜索 -->
        <div class="cl-toolbar-right">
          <a-input-search
            v-model="queryParams.confKey"
            :placeholder="t('system.config.searchPlaceholder')"
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
                <a-input v-model="queryParams.name" :placeholder="t('system.config.namePlaceholder')" allow-clear />
                <a-select v-model="queryParams.isSys" :placeholder="t('system.config.isSysPlaceholder')" allow-clear>
                  <a-option v-for="dict in dicts.sys_is" :value="dict.value">{{ dict.label }}</a-option>
                </a-select>
                <a-select v-model="queryParams.scope" :placeholder="t('system.config.scopePlaceholder')" allow-clear>
                  <a-option v-for="dict in dicts.sys_conf_scope" :value="dict.value">{{ dict.label }}</a-option>
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
            <a-table-column :title="t('system.config.name')" data-index="name" :ellipsis="true" :tooltip="true"/>
            <a-table-column :title="t('system.config.confKey')" data-index="confKey">
              <template #cell="{ record }">
                <code class="mm-code mm-code-sign">{{ record.confKey }}</code>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.config.confValue')" data-index="confValue" :ellipsis="true" :tooltip="true"/>
            <a-table-column :title="t('system.config.isSys')" data-index="isSys" align="center" width="100">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_is" :value="record.isSys"/>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.config.scope')" data-index="scope" align="center" width="150"  :ellipsis="true" :tooltip="true">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_conf_scope" :value="record.scope"/>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.config.isSecret')" data-index="isSecret" align="center" width="100">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_is" :value="record.isSecret"/>
              </template>
            </a-table-column>
            <a-table-column :title="t('common.remark')" data-index="remark" :ellipsis="true" :tooltip="true"/>
            <a-table-column :title="t('system.config.updatedAt')" data-index="updatedAt" align="center"/>
            <a-table-column :title="t('common.action')" :width="160" align="center">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button
                    v-permission="'sys:conf$edit'"
                    size="mini"
                    type="text"
                    @click.stop="handleEdit(record)"
                  >
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button
                    v-permission="'sys:conf$delete'"
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

        <a-skeleton v-if="loading" :animation="true">
          <a-skeleton-line :rows="6" />
        </a-skeleton>
        <template v-else>
          <div
            v-for="(item, index) in mobileData"
            :key="item.id"
            class="cl-card stagger-item"
            :style="{ '--stagger-index': index % mobilePageSize }"
          >
            <div class="cl-card-header">
              <div class="cl-card-identity">
                <strong>{{ item.confKey }}</strong>
                <span class="cl-card-sub">{{ item.name }}</span>
              </div>
            </div>
            <div class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'sys:conf$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button
                  v-permission="'sys:conf$delete'"
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
      :title="isEdit ? t('system.config.editTitle') : t('system.config.addTitle')"
      :width="isMobile ? '100%' : 560"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
        <a-col :span="24">
            <a-form-item :label="t('system.config.name')" field="name" 
            :rules="[{ required: true, message: t('system.config.nameRequired') },{ maxLength: 255, message: t('common.maxLength', { max: 255 }) },]" :validate-trigger="['blur']">
              <a-input v-model="form.name" :placeholder="t('system.config.namePlaceholderInput')" />
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('system.config.confKey')" field="confKey" 
            :rules="[{ required: true, message: t('system.config.confKeyRequired') },{ maxLength: 64, message: t('common.maxLength', { max: 64 }) },]" :validate-trigger="['blur']">
              <a-input v-model="form.confKey" :placeholder="t('system.config.confKeyPlaceholder')" />
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('system.config.confValue')" field="confValue" 
            :rules="[{ required: true, message: t('system.config.confValueRequired') },{ maxLength: 2000, message: t('common.maxLength', { max: 2000 }) },]" :validate-trigger="['blur']">
              <a-input v-model="form.confValue" :placeholder="t('system.config.confValuePlaceholder')" />
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('system.config.scope')" field="scope" 
            :rules="[{ required: true, message: t('system.config.scopeRequired') }]" :validate-trigger="['blur']">
              <a-radio-group v-model="form.scope" type="button">
                <a-radio v-for="dict in dicts.sys_conf_scope" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.config.isSys')" field="isSys" 
            :rules="[{ required: true, message: t('system.config.isSysRequired') }]" :validate-trigger="['blur']">
              <a-switch v-model="form.isSys" :checked-value="1" :unchecked-value="0" :checked-text="t('system.config.isSysYes')" :unchecked-text="t('system.config.isSysNo')" />
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.config.isSecret')" field="isSecret" 
            :rules="[{ required: true, message: t('system.config.isSecretRequired') }]" :validate-trigger="['blur']">
              <a-switch v-model="form.isSecret" :checked-value="1" :unchecked-value="0" :checked-text="t('system.config.isSecretYes')" :unchecked-text="t('system.config.isSecretNo')" />
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('common.remark')" field="remark" 
            :rules="[{ maxLength: 500, message: t('common.maxLength', { max: 500 }) },]" :validate-trigger="['blur']">
              <a-textarea v-model="form.remark" :placeholder="t('system.config.remarkPlaceholder')" allow-clear/>
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

  </div>
</template>