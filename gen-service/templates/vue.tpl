<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import {{.ClassName | toLowerCamel}}Api from "@/api/{{.ModuleName}}/{{.CleanName}}"

const proxy = getCurrentInstance()!.proxy as any
// ── 响应式断点（576px 以下为移动端）────────────────────
const isMobile = ref(window.innerWidth < 576);
function onResize() { isMobile.value = window.innerWidth < 576; }
onMounted(() => window.addEventListener('resize', onResize));
onUnmounted(() => window.removeEventListener('resize', onResize));

const dicts = proxy.$useDict( 
  {{- $first := true -}}
  {{- range .Fields -}}
    {{- if ne .DictType "" -}}
      {{- if not $first -}}, {{- end -}}
      '{{- .DictType -}}'
      {{- $first = false -}}
    {{- end -}}
  {{- end -}}
)
// ── 搜索 & 筛选 ────────────────────────────────────────
const filterModule = ref(false);
const advancedVisible = ref(false);

// ── PC 分页 ────────────────────────────────────────────
const total = ref(0);
const loading = ref(false);
const dataList = ref<any[]>([]);
const queryParams = ref({
  pageNum:1,
  pageSize:20,{{range .Fields -}}
  {{if ne .Condition ""}}
  {{.FieldName | toLowerCamel}}: '',{{end}}{{end}}
})

async function loadPage() {
  loading.value = true;
  try {
    let res:any = await {{.ClassName | toLowerCamel}}Api.page(queryParams.value);
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
    {{range .Fields -}}
  {{ if ne .Condition ""}}
  queryParams.value.{{.FieldName | toLowerCamel}} = ''{{end}}{{end}}
  handleSearch();
}
watch(
  () => ({ ...queryParams.value }),
  () => {
    filterModule.value=({{- $first := true -}}
  {{- range .Fields -}}
    {{- if ne .Condition "" -}}
      {{- if not $first -}} || {{- end -}}
      queryParams.value.{{.FieldName | toLowerCamel}} !== ''
      {{- $first = false -}}
    {{- end -}}
  {{- end -}})
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
    Message.warning('请先选择要删除的记录');
    return;
  }
  Modal.confirm({
    title: '确认删除',
    content: `确认删除选中的 ${selectedIds.value.length} 条记录？`,
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await {{.ClassName | toLowerCamel}}Api.delete(selectedIds.value);
      Message.success('删除成功');
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
async function handleEdit(row:any) {
  isEdit.value = true;
  resetForm(row.id);
  modalVisible.value = true;
}

async function resetForm(id:string){
  formRef.value.clearValidate()
  if(isEdit.value){
    let res = await {{.ClassName | toLowerCamel}}Api.getById(id)
    Object.assign(form, res);
  }else{
    Object.assign(form, {
      {{range .Fields -}}{{- if .IsInsert -}}
      {{.FieldName | toLowerCamel}}:{{- if ne .DefaultValue "" -}}{{.DefaultValue}}{{- else -}}undefined{{- end -}},
      {{- end }}
      {{end}}
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
      await {{.ClassName | toLowerCamel}}Api.update(form.id!, form);
      Message.success('修改成功');
    } else {
      await {{.ClassName | toLowerCamel}}Api.add(form);
      Message.success('新增成功');
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
    title: '确认删除',
    content: `确认删除{{.Comment}}编号为「${row.id}」的记录吗？`,
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await {{.ClassName | toLowerCamel}}Api.delete([row.id]);
      Message.success('删除成功');
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
          <a-button v-permission="'{{.Permission}}$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            新增
          </a-button>

          <a-button
            v-permission="'{{.Permission}}$delete'"
            status="danger"
            :disabled="!selectedIds.length"
            @click="batchDelete"
          >
            <template #icon><IconDelete /></template>
            批量删除
          </a-button>
        </a-space>

        <!-- 搜索 -->
        <div class="cl-toolbar-right">
          <a-input-search
            v-model="queryParams.id"
            placeholder="搜索关键字"
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
                <p class="cl-filter-title">高级筛选</p>
                {{range .Fields -}}{{if ne .Condition ""}}
                <a-input v-model="queryParams.{{.FieldName | toLowerCamel}}" placeholder="{{.ColumnComment}}" allow-clear />{{end}}{{end}}
                <div class="cl-filter-actions">
                  <a-button size="small" @click="handleReset">重置</a-button>
                  <a-button size="small" type="primary" @click="() => { handleSearch(); advancedVisible = false }">
                    搜索
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
        {{range .Fields -}}{{- if .IsList }}
          <a-table-column title="{{.ColumnComment}}" data-index="{{.FieldName | toLowerCamel}}"
          {{- if .Sortable}} :sortable="{sortDirections: ['ascend', 'descend']}"{{end}} align="center"/>
          {{- end }}{{ end }}
          <a-table-column title="操作" :width="260" align="center">
            <template #cell="{ record }">
              <a-space size="mini">
                <a-button
                  v-permission="'{{.Permission}}$edit'"
                  size="mini"
                  type="text"
                  @click.stop="handleEdit(record)"
                >
                  <template #icon><IconEdit /></template>修改
                </a-button>
                <a-button
                  v-permission="'{{.Permission}}$delete'"
                  size="mini"
                  type="text"
                  status="danger"
                  @click.stop="handleDelete(record)"
                >
                  <template #icon><IconDelete /></template>删除
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
            {{`{{ mobileSelectMode ? '取消' : '选择' }}`}}
          </a-button>
          <template v-if="mobileSelectMode">
            <a-checkbox
              :model-value="mobileAllSelected"
              :indeterminate="mobileIndeterminate"
              @change="toggleMobileSelectAll"
            >
              全选
            </a-checkbox>
            <span class="cl-select-count">已选 {{ `{{ selectedIds.length }}` }} 项</span>
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
                <strong>{{ `{{ item.id }}` }}</strong>
                <span class="cl-card-sub">{{ `{{ item.id || '暂无数据' }}` }}</span>
              </div>
              <a-tag color="arcoblue" size="small">{{ `{{ item.id || '-' }}` }}</a-tag>
            </div>
            <div class="cl-card-meta">
              <span>xxxxx：{{ `{{ item.id || '-' }}` }}</span>
              <span>xxxxx：{{ `{{ item.id || '-' }}` }}</span>
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'{{.Permission}}$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>编辑
                </a-button>
                <a-button
                  v-permission="'{{.Permission}}$delete'"
                  size="mini"
                  type="outline"
                  status="danger"
                  @click="handleDelete(item)"
                >
                  <template #icon><IconDelete /></template>删除
                </a-button>
              </a-space>
            </div>
          </div>

          <a-empty v-if="!mobileData.length" description="暂无数据" />
          <div ref="sentinelRef" class="cl-sentinel">
            <span v-if="!mobileHasMore && mobileData.length > 0" class="cl-no-more">已全部加载</span>
          </div>
        </template>
      </div>
    </a-card>

    <!-- ── 新增 / 修改 共用弹窗 ── -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? '编辑{{.Comment}}' : '新增{{.Comment}}'"
      :width="isMobile ? '100%' : 560"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
        {{range .Fields -}}{{- if or .IsInsert .IsEdit -}}
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item label="{{.ColumnComment}}" field="{{.FieldName | toLowerCamel}}" {{- if or .IsRequired (gt .DataLength 0)}} 
            :rules="[{{- if .IsRequired -}}{ required: true, message: '{{.ColumnComment}}不允许为空' },{{- end}}
                    {{- if gt .DataLength 0 -}}{ maxLength: {{.DataLength}}, message: '{{.ColumnComment}}长度不能超过{{.DataLength}}个字符' },{{- end -}}
            ]" :validate-trigger="['blur']"{{- end -}}>
              {{if eq .Component "input-number"}}
              <a-input-number v-model="form.{{.FieldName | toLowerCamel}}" placeholder="请输入{{.ColumnComment}}" />
              {{else if eq .Component "select"}}
              <a-select v-model="form.{{.FieldName | toLowerCamel}}" placeholder="请选择{{.ColumnComment}}" allow-clear >
                {{if ne .DictType ""}}
                <a-option v-for="dict in dicts.{{.DictType}}" :key="dict.dictValue" :value="dict.dictValue">{{ `{{ dict.label }}` }}</a-option>
                {{end}}
              </a-select>
              {{else if eq .Component "checkbox"}}
              <a-checkbox-group v-model="form.{{.FieldName | toLowerCamel}}">
                {{if ne .DictType ""}}
                <a-checkbox v-for="dict in dicts.{{.DictType}}" :key="dict.dictValue" :value="dict.dictValue">{{ `{{ dict.label }}` }}</a-checkbox>
                {{end}}
              </a-checkbox-group>
              {{else if eq .Component "radio"}}
              <a-radio-group v-model="form.{{.FieldName | toLowerCamel}}" type="button">
                {{if ne .DictType ""}}
                <a-radio v-for="dict in dicts.{{.DictType}}" :key="dict.dictValue" :value="dict.dictValue">{{ `{{ dict.label }}` }}</a-radio>
                {{end}}
              </a-radio-group>
              {{else if eq .Component "switch"}}
              <a-switch v-model="form.{{.FieldName | toLowerCamel}}" :checked-value="1" :unchecked-value="0" checked-text="是"  unchecked-text="否" />
              {{else if eq .Component "textarea"}}
              <a-textarea v-model="form.{{.FieldName | toLowerCamel}}" placeholder="请输入{{.ColumnComment}}" allow-clear/>
              {{else if eq .Component "multiple-select"}}
              <a-select v-model="form.{{.FieldName | toLowerCamel}}" placeholder="请选择{{.ColumnComment}}" multiple>
                {{if ne .DictType ""}}
                <a-option v-for="dict in dicts.{{.DictType}}" :key="dict.dictValue" :value="dict.dictValue">{{ `{{ dict.label }}` }}</a-option>
                {{end}}
              </a-select>
              {{else if eq .Component "tree-select"}}
              <a-tree-select v-model="form.{{.FieldName | toLowerCamel}}" :data="treeData" :allow-search="true" :allow-clear="true"
                :disable-filter="true" :label-in-value="true"
                :default-value="{ value: 'id', label: 'name' }"
                placeholder="请选择{{.ColumnComment}}"></a-tree-select>
              {{else if eq .Component "date-picker"}}
              <a-date-picker v-model="form.{{.FieldName | toLowerCamel}}" show-time format="YYYY-MM-DD hh:mm" placeholder="请选择{{.ColumnComment}}" />
              {{else if eq .Component "time-picker"}}
              <a-time-picker v-model="form.{{.FieldName | toLowerCamel}}" placeholder="请选择{{.ColumnComment}}"/>
              {{else if eq .Component "input-tag"}}
              <a-input-tag v-model="form.{{.FieldName | toLowerCamel}}"  placeholder="请输入{{.ColumnComment}}" :max-tag-count="3" allow-clear/>
              {{else}}
              <a-input v-model="form.{{.FieldName | toLowerCamel}}" placeholder="请输入{{.ColumnComment}}" />
              {{end}}
            </a-form-item>
          </a-col>
        {{- end }}
        {{end}}
          <a-col :span="24" style="margin-top: 30px;">
            <div style="display: flex; justify-content: center; gap: 12px;">
              <a-button @click="modalVisible = false">取消</a-button>
              <a-button type="primary" @click="submitForm()" :loading="submitLoading">保存</a-button>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

  </div>
</template>