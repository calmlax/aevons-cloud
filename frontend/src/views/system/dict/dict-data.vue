<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import { useI18n } from 'vue-i18n'
import dictDataApi from "@/api/system/dict_data"
import langApi from "@/api/system/lang"
import { DictTag } from '@/components/DictTag';

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any

const dicts = proxy.$useDict('sys_is', 'sys_common_status','sys_dict_tag_style')

const props = defineProps({
    isMobile: Boolean
});

const visible = ref(false)
const title = ref('')

function cancel() {
    visible.value = false
}

const dict = ref<any>({})

const langList = ref<any[]>([])

async function open(d:any) {
  visible.value = true
  title.value = d.dictType
  dict.value = {...d}
  loadPage()
  langList.value = await langApi.availableList() as any

  isEdit.value = false;
  resetForm('');
}

// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

const loading = ref(false);
const dataList = ref<any[]>([]);
const queryParams = ref({
  dictType: '',
  field: 'sort',
  direction: 'ascend',
})

async function loadPage() {
  loading.value = true;
  try {
    queryParams.value.dictType = dict.value.dictType
    let res:any = await dictDataApi.list(queryParams.value);
    dataList.value = res ?? [];
  } finally {
    loading.value = false;
  }
}

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
watch(dataList, () => {
  mobilePage.value = mobilePageSize;
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
      await dictDataApi.delete(selectedIds.value);
      Message.success(t('common.deleteSuccess'));
      selectedIds.value = [];
      loadPage();
    },
  });
}

// ── 新增 / 修改 共用弹窗 ────────────────────────────────
const modalVisible = ref(false);
const form = reactive<any>({});
const formRef = ref()
const isEdit = ref(false); // 标记是否为编辑模式
const submitLoading = ref(false)

// 打开新增
function handleAdd() {
  isEdit.value = false;
  resetForm('');
  modalVisible.value = true;
}

// 打开编辑
function handleEdit(record:any) {
  isEdit.value = true;
  resetForm(record.id);
  modalVisible.value = true;
}

async function resetForm(id:string){
  formRef.value.clearValidate()
  if(isEdit.value){
    let res = await dictDataApi.getById(id)
    form.translations = {}
    Object.assign(form, res);
    langList.value.forEach(lang => {
      if (!form.translations[lang.langCode]) {
        form.translations[lang.langCode] = {
          label: '',
          tip: ''
        }
      }
    })
  }else{
    Object.assign(form, {
      id:undefined,
      dictType:dict.value.dictType,
      dictValue:undefined,
      status:0,
      sort:'0',
      tagType:undefined,
      tagClass:undefined,
    });
    form.translations = {}
    langList.value.forEach(lang => {
      form.translations[lang.langCode] = {
        label: '',
        tip: ''
      }
    })
  }
}

// 统一提交
async function submitForm() {
  // submitLoading.value = true
  try {
    const valid = await formRef.value.validate()
    if (valid) return false
    if (isEdit.value) {
      await dictDataApi.update(form.id!, form);
      Message.success(t('common.editSuccess'));
    } else {
      await dictDataApi.add(form);
      Message.success(t('common.addSuccess'));
    }
    modalVisible.value = false;
    loadPage();
  } finally {
    // submitLoading.value = false
  }
}

// 拖拽排序
async function handleTableChange(data: any[]) {
  dataList.value = data;
  const items = data.map((item: any, index: number) => ({ id: item.id, sort: index + 1 }));
  try {
    await dictDataApi.updateSort(items);
    Message.success(t('system.dict.sortSaved'));
  } catch {
    Message.error(t('system.dict.sortFailed'));
    loadPage();
  }
}

// 删除单条
async function confirmDelete(record: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.dict.confirmDeleteData', { value: record.value }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await dictDataApi.delete([record.id]);
      Message.success(t('common.deleteSuccess'));
      loadPage();
    },
  });
}

defineExpose({
    open,
    cancel,
});
</script>

<template>
  <a-drawer
      v-model:visible="visible"
      :title="t('system.dict.dataTitle', { type: title })"
      :placement="props.isMobile ? 'bottom' : 'right'"
      :width="props.isMobile ? '100%' : '50%'"
      :height="props.isMobile ? '85%' : undefined"
      :footer="false"
      append-to-body
    >

  
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            {{ t('common.add') }}
          </a-button>

          <a-button
            status="danger"
            :disabled="!selectedIds.length"
            @click="batchDelete"
          >
            <template #icon><IconDelete /></template>
            {{ t('common.batchDelete') }}
          </a-button>
        </a-space>
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
          :pagination="false"
          :draggable="{ type: 'handle', width: 40 }"
          @change="handleTableChange"
        >
          <template #columns>
            <a-table-column :title="t('system.dict.dictValue')" data-index="dictValue" :width="200"/>
            <a-table-column :title="t('system.dict.dictLabel')" data-index="label" :width="200">
              <template #cell="{ record }">
                <DictTag :options="dataList" :value="record.dictValue" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.status')" :width="100" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_common_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.action')" :width="180" align="center" fixed="right">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button
                    size="mini"
                    type="text"
                    @click.stop="handleEdit(record)"
                  >
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button
                    size="mini"
                    type="text"
                    status="danger"
                    @click.stop="confirmDelete(record)"
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
                <strong>{{ item.dictValue }}</strong>
                <span class="cl-card-sub">
                  <DictTag :options="dataList" :value="item.dictValue" dot v-if="item.label != ''" /><span v-else>{{ t('common.noData') }}</span></span>
              </div>
              <DictTag :options="dicts.sys_common_status" :value="item.status" dot />
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button
                  size="mini"
                  type="outline"
                  status="danger"
                  @click="confirmDelete(item)"
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
      :width="isMobile ? '100%' : 700"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="24">
              <a-form-item :label="t('common.status')" field="status" :validate-trigger="['blur']">
                <a-radio-group v-model="form.status" type="button">
                  <a-radio v-for="dict in dicts.sys_common_status" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
                </a-radio-group>
              </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('system.dict.tagStyle')" field="tagType">
                <a-radio-group v-model="form.tagType" type="button">
                  <a-radio v-for="dict in dicts.sys_dict_tag_style" :value="dict.dictValue"> {{ dict.label }}</a-radio>
                </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="24">
              <a-form-item :label="t('system.dict.dictValue')" field="dictValue" 
              :rules="[{ required: true, message: t('system.dict.dictValueRequired') },{ maxLength: 32, message: t('common.maxLength', { max: 32 }) },]" :validate-trigger="['blur']">
                <a-input v-model="form.dictValue" :placeholder="t('system.dict.dictValuePlaceholder')" />
              </a-form-item>
          </a-col>
          <a-col :span="24">

            <a-form-item  v-for="lang in langList" :label="lang.langName" :content-flex="false" :merge-props="false">
              <a-space direction="vertical" fill>
                <a-form-item hide-label
                  :key="lang.langCode"
                  :field="`translations.${lang.langCode}.label`"
                  :rules="[{ required: true, message: t('system.dict.dictLabelRequired', { lang: lang.langName }) },{ maxLength: 50, message: t('common.maxLength', { max: 50 }) }]"
                  :validate-trigger="['blur']"
                >
                  <a-input v-model="form.translations[lang.langCode].label" :placeholder="t('system.dict.dictLabelPlaceholder')">
                    <template #prepend>
                      <div style="width: 60px;text-align: center;">{{ t('system.dict.dictLabelPrepend') }}</div>
                    </template>
                  </a-input>
                </a-form-item>
                <a-form-item hide-label
                  :rules="[{ maxLength: 200, message: t('common.maxLength', { max: 200 }) }]"
                  :validate-trigger="['blur']"
                >
                  <a-input v-model="form.translations[lang.langCode].tip" :placeholder="t('system.dict.tipPlaceholder')">
                    <template #prepend>
                      <div style="width: 60px;text-align: center;">{{ t('system.dict.tipPrepend') }}</div>
                    </template>
                  </a-input>
                </a-form-item>
              </a-space>
            </a-form-item>
          </a-col>
          <!-- <a-col :span="24">
              <a-form-item label="样式类名" field="tagClass" 
              :rules="[{ maxLength: 10, message: '样式类名长度不能超过10个字符' },]" :validate-trigger="['blur']">
                <a-input v-model="form.tagClass" placeholder="请输入样式类名" />
              </a-form-item>
          </a-col> -->
          <a-col :span="24">
              <a-form-item :label="t('system.dict.sort')" field="sort">
                <a-input v-model="form.sort" :placeholder="t('system.dict.sortPlaceholder')" />
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
  </a-drawer>
</template>
<style scoped>
:deep(.arco-radio-group .arco-radio) {
  margin-right: 5px;
  margin-bottom: 7px;
}
.custom-radio-card {
  padding: 5px 5px;
  border: 1px solid var(--color-border-2);
  border-radius: 4px;
  width: inherit;
  box-sizing: border-box;
}

.custom-radio-card-mask {
  height: 14px;
  width: 14px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 100%;
  border: 1px solid var(--color-border-2);
  box-sizing: border-box;
}

.custom-radio-card-mask-dot {
  width: 8px;
  height: 8px;
  border-radius: 100%;
}

.custom-radio-card-title {
  color: var(--color-text-1);
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 8px;
}

.custom-radio-card:hover,
.custom-radio-card-checked,
.custom-radio-card:hover .custom-radio-card-mask,
.custom-radio-card-checked  .custom-radio-card-mask{
  border-color: rgb(var(--primary-6));
}

.custom-radio-card-checked {
  background-color: var(--color-primary-light-1);
}

.custom-radio-card:hover .custom-radio-card-title,
.custom-radio-card-checked .custom-radio-card-title {
  color: rgb(var(--primary-6));
}

.custom-radio-card-checked .custom-radio-card-mask-dot {
  background-color: rgb(var(--primary-6));
}
</style>