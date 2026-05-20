<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconCheck, IconClose, IconRefresh } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n'
import oauthClientApi from "@/api/system/oauth_client"

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any

// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

const dicts = proxy.$useDict('sys_is','sys_authorized_scope','sys_authorized_grant_type')

// ── 搜索 & 筛选 ────────────────────────────────────────
const filterModule = ref(false);
const advancedVisible = ref(false);

// ── PC 分页 ────────────────────────────────────────────
const total = ref(0);
const loading = ref(false);
const dataList = ref<any[]>([]);


const data = reactive({
  queryParams:{
    pageNum:1,
    pageSize:20,
    clientId: '',
    clientName: '',
    authorizedGrantTypes: '',
    createdAt: '',
  }
})

const { queryParams } = toRefs(data)

async function loadPage() {
  loading.value = true;
  try {
    let res:any = await oauthClientApi.page(queryParams.value);
    dataList.value = res?.rows ?? [];
    console.log(dataList.value)
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
  queryParams.value.clientId = ''
  queryParams.value.clientName = ''
  queryParams.value.authorizedGrantTypes = ''
  queryParams.value.createdAt = ''
  handleSearch();
}
watch(
  () => ({ ...queryParams.value }),
  () => {
    filterModule.value=(queryParams.value.clientId !== ''||queryParams.value.authorizedGrantTypes !== '')
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
const refreshGatewayCacheLoading = ref(false);

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
      await oauthClientApi.delete(selectedIds.value);
      Message.success(t('common.deleteSuccess'));
      selectedIds.value = [];
      loadPage();
    },
  });
}

async function refreshGatewayCache() {
  refreshGatewayCacheLoading.value = true;
  try {
    const res:any = await oauthClientApi.refreshGatewayCache();
    const count = res?.count ?? 0;
    Message.success(t('oauthClient.refreshGatewayCacheSuccess', { count }));
  } catch (error) {
    console.error(error);
    Message.error((error as Error)?.message || t('oauthClient.refreshGatewayCacheFailed'));
  } finally {
    refreshGatewayCacheLoading.value = false;
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
  formRef.value.clearValidate()
  if(isEdit.value){
    let res = await oauthClientApi.getById(id) as any
    if (res?.scope) {
      res.scope = res.scope.split(',').map((s:string) => s.trim()).filter(Boolean);
    }
    if (res?.authorizedGrantTypes) {
      res.authorizedGrantTypes = res.authorizedGrantTypes.split(',').map((s:string) => s.trim()).filter(Boolean);
    }
    Object.assign(form, res);
  }else{
    Object.assign(form, {
      clientId:'',
      clientSecret:'',
      clientName:'',
      resources:'',
      logoUri:'',
      scope:[],
      authorizedGrantTypes:[],
      webServerRedirectUri:'',
      accessTokenValidity:7200,
      refreshTokenValidity:604800,
      autoapprove:0,
    });
  }
}
// 统一提交
async function submitForm() {
  try {
    const valid = await formRef.value.validate()
    if (valid) return false
    submitLoading.value = true
    var formData = {...form}
    formData.scope = form.scope.join(',')
    formData.authorizedGrantTypes = form.authorizedGrantTypes.join(',')
    if (isEdit.value) {
      await oauthClientApi.update(form.id!, formData);
      Message.success(t('common.editSuccess'));
    } else {
      await oauthClientApi.add(formData);
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
    content: t('oauthClient.confirmDelete', { name: row.clientName || row.id }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await oauthClientApi.delete([row.id]);
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
          <a-button v-permission="'sys:oauth:client$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            {{ t('common.add') }}
          </a-button>
          <a-button
            v-permission="'sys:oauth:client$edit'"
            :loading="refreshGatewayCacheLoading"
            @click="refreshGatewayCache"
          >
            <template #icon><IconRefresh /></template>
            {{ t('oauthClient.refreshGatewayCache') }}
          </a-button>

          <a-button
            v-permission="'sys:oauth:client$delete'"
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
            v-model="queryParams.clientName"
            :placeholder="t('oauthClient.searchPlaceholder')"
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
                <a-input v-model="queryParams.clientId" placeholder="AppID" allow-clear />
                <a-input v-model="queryParams.authorizedGrantTypes" :placeholder="t('oauthClient.grantType')" allow-clear />
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
            <a-table-column title="AppID" data-index="clientId"/>
            <a-table-column :title="t('oauthClient.appName')" data-index="clientName"/>
            <!-- <a-table-column title="客户端LOGO" data-index="logoUri" align="center"/> -->
            <a-table-column :title="t('oauthClient.scope')" align="center">
              <template #cell="{ record }">
                <template v-if="record.scope && record.scope.trim()">
                  <a-space size="4px" wrap>
                    <a-tag v-for="(type, index) in record.scope.split(',')" :key="index" color="blue" size="small">
                      {{ type.trim() }}
                    </a-tag>
                  </a-space>
                </template>
                <span v-else class="text-gray-400">-</span>
              </template>
            </a-table-column>
            <a-table-column :title="t('oauthClient.grantType')" align="center">
              <template #cell="{ record }">
                <template v-if="record.authorizedGrantTypes && record.authorizedGrantTypes.trim()">
                  <a-space size="4px" wrap>
                    <a-tag v-for="(type, index) in record.authorizedGrantTypes.split(',')" :key="index" color="blue" size="small">
                      {{ type.trim() }}
                    </a-tag>
                  </a-space>
                </template>
                <span v-else class="text-gray-400">-</span>
              </template>
            </a-table-column>
            <a-table-column :title="t('oauthClient.autoApprove')" data-index="autoapprove" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_is" :value="record.autoapprove"/>
              </template>
            </a-table-column>
            <a-table-column :title="t('oauthClient.createdAt')" data-index="createdAt" align="center"/>
            <a-table-column :title="t('common.action')" :width="260" align="center">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button
                    v-permission="'sys:oauth:client$edit'"
                    size="mini"
                    type="text"
                    @click.stop="handleEdit(record)"
                  >
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button
                    v-permission="'sys:oauth:client$delete'"
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
                <strong>{{ item.clientName }}</strong>
                <span class="cl-card-sub">{{ item.clientId }}</span>
              </div>
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'sys:oauth:client$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button
                  v-permission="'sys:oauth:client$delete'"
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
      :title="isEdit ? t('common.edit') : t('common.add')"
      :width="isMobile ? '100%' : 560"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item label="AppID" field="clientId"
            :rules="[{ required: true, message: t('common.required') },{ maxLength: 32, message: t('common.maxLength', { max: 32 }) }]" :validate-trigger="['blur']">
              <a-input v-model="form.clientId" :placeholder="t('oauthClient.appId')" />
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('oauthClient.appName')" field="clientName"
            :rules="[{ required: true, message: t('common.required') },{ maxLength: 255, message: t('common.maxLength', { max: 255 }) }]" :validate-trigger="['blur']">
              <a-input v-model="form.clientName" :placeholder="t('oauthClient.appName')" />
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('oauthClient.appSecret')" field="clientSecret" v-if="isEdit"
            :rules="[{ maxLength: 32, message: t('common.maxLength', { max: 32 }) }]" :validate-trigger="['blur']">
              <a-input v-model="form.clientSecret" :placeholder="t('oauthClient.appSecret')" />
            </a-form-item>
            <a-form-item :label="t('oauthClient.appSecret')" field="clientSecret" v-else
            :rules="[{ required: true, message: t('common.required') },{ maxLength: 32, message: t('common.maxLength', { max: 32 }) }]" :validate-trigger="['blur']">
              <a-input v-model="form.clientSecret" :placeholder="t('oauthClient.appSecret')" />
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('oauthClient.resources')" field="resources"
            :rules="[{ maxLength: 2048, message: t('common.maxLength', { max: 2048 }) }]" :validate-trigger="['blur']">
              <a-input
                v-model="form.resources"
                :placeholder="t('oauthClient.resourcesPlaceholder')"
              />
              <template #extra>
                {{ t('oauthClient.resourcesHint') }}
              </template>
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('oauthClient.scope')" field="scope"
            :rules="[{ required: true, message: t('common.required') }]" :validate-trigger="['blur']">
              <a-select v-model="form.scope" :placeholder="t('oauthClient.scope')" multiple>
                <a-option v-for="dict in dicts.sys_authorized_scope" :value="dict.dictValue">{{ dict.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('oauthClient.grantType')" field="authorizedGrantTypes"
            :rules="[{ required: true, message: t('common.required') }]" :validate-trigger="['blur']">
              <a-select v-model="form.authorizedGrantTypes" :placeholder="t('oauthClient.grantType')" multiple>
                <a-option v-for="dict in dicts.sys_authorized_grant_type" :value="dict.dictValue">{{ dict.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        <a-col :span="24">
            <a-form-item :label="t('oauthClient.redirectUri')" field="webServerRedirectUri"
            :rules="[{ required: true, message: t('common.required') },{ maxLength: 256, message: t('common.maxLength', { max: 256 }) }]" :validate-trigger="['blur']">
              <a-input v-model="form.webServerRedirectUri" :placeholder="t('oauthClient.redirectUri')" />
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('oauthClient.accessTokenTtl')" field="accessTokenValidity"
            :rules="[{ required: true, message: t('common.required') }]" :validate-trigger="['blur']">
              <a-input-number v-model="form.accessTokenValidity" :placeholder="t('oauthClient.accessTokenTtl')" />
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('oauthClient.refreshTokenTtl')" field="refreshTokenValidity"
            :rules="[{ required: true, message: t('common.required') }]" :validate-trigger="['blur']">
              <a-input-number v-model="form.refreshTokenValidity" :placeholder="t('oauthClient.refreshTokenTtl')" />
            </a-form-item>
          </a-col>
        <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('oauthClient.autoApprove')" field="autoapprove"
            :rules="[{ required: true, message: t('common.required') }]" :validate-trigger="['blur']">
              <a-radio-group v-model="form.autoapprove">
                <a-radio v-for="dict in dicts.sys_is" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
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
