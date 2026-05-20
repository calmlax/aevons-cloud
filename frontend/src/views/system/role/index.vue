<script setup lang="ts">
import { IconDelete, IconEdit, IconFilter, IconPlus, IconCheck, IconClose } from '@arco-design/web-vue/es/icon';
import roleApi from "@/api/system/role"
import menuApi from "@/api/system/menu"
import deptApi from "@/api/system/dept"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any
const dicts = proxy.$useDict('sys_common_status','sys_data_scope')

// ── 响应式断点 ────────────────────────────────────────
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
  pageNum: 1,
  pageSize: 20,
  roleName: '',
  roleKey: '',
  status: '',
})

async function loadPage() {
  loading.value = true;
  try {
    let res: any = await roleApi.page(queryParams.value);
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
  queryParams.value.roleName = ''
  queryParams.value.roleKey = ''
  queryParams.value.status = ''
  handleSearch();
}
watch(
  () => ({ ...queryParams.value }),
  () => {
    filterModule.value = (queryParams.value.roleKey !== '' || queryParams.value.status !== '')
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
  if (!mobileSelectMode.value) selectedIds.value = [];
}

function toggleCardSelect(id: string) {
  const idx = selectedIds.value.indexOf(id);
  idx === -1 ? selectedIds.value.push(id) : selectedIds.value.splice(idx, 1);
}

function isCardSelected(id: string) {
  return selectedIds.value.includes(id);
}

const mobileAllSelected = computed(() =>
  mobileData.value.length > 0 && mobileData.value.every((item: any) => selectedIds.value.includes(item.id))
);
const mobileIndeterminate = computed(() =>
  selectedIds.value.length > 0 && !mobileAllSelected.value
);

function toggleMobileSelectAll() {
  const currentIds = mobileData.value.map((item: any) => item.id);
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
      await roleApi.delete(selectedIds.value);
      Message.success(t('common.deleteSuccess'));
      selectedIds.value = [];
      loadPage();
    },
  });
}

// ── 菜单树 ─────────────────────────────────────────────
const menuTreeData = ref<any[]>([])
const checkedMenuKeys = ref<string[]>([])
const menuExpandedKeys = ref<string[]>([])
const deptTreeData = ref<any[]>([])
const checkedDeptKeys = ref<string[]>([])

function toKeyStrings(values: any[]): string[] {
  return (values ?? []).map((item: any) => {
    if (item && typeof item === 'object') {
      return String(item.value ?? item.key ?? item.id ?? '')
    }
    return String(item)
  }).filter(Boolean)
}

// 将平铺菜单列表转为树形结构
function buildMenuTree(list: any[]): any[] {
  const map: Record<string, any> = {}
  const roots: any[] = []
  list.forEach(item => {
    const id = String(item.id)
    const label = (item.title && item.title.trim()) || item.permission || id
    map[id] = { key: id, title: label, children: [] }
  })
  list.forEach(item => {
    const id = String(item.id)
    const parentId = String(item.parentId ?? '0')
    if (parentId && parentId !== '0' && map[parentId]) {
      map[parentId].children.push(map[id])
    } else if (parentId === '0') {
      roots.push(map[id])
    }
    // parentId 不为 0 且父节点不存在则跳过（孤儿节点）
  })
  // 清理空 children，避免 a-tree 渲染叶节点时出错
  function clean(nodes: any[]) {
    nodes.forEach(n => {
      if (!n.children || n.children.length === 0) {
        delete n.children
      } else {
        clean(n.children)
      }
    })
  }
  clean(roots)
  return roots
}

async function loadMenuTree() {
  const list: any = await menuApi.list()
  menuTreeData.value = buildMenuTree(list ?? [])
  // 默认展开一级（根节点）
  menuExpandedKeys.value = menuTreeData.value.map((n: any) => n.key)
}

function buildDeptTree(list: any[], parentId = '0'): any[] {
  return (list ?? [])
    .filter((item: any) => String(item.parentId ?? '0') === String(parentId))
    .sort((a: any, b: any) => Number(a.sort ?? 0) - Number(b.sort ?? 0))
    .map((item: any) => {
      const children = buildDeptTree(list, item.id)
      const node = {
        key: String(item.id),
        title: item.deptName || String(item.id),
      } as any
      if (children.length) node.children = children
      return node
    })
}

async function loadDeptTree() {
  const list: any = await deptApi.list({})
  deptTreeData.value = buildDeptTree(list ?? [])
}

function findNode(nodes: any[], key: string): any {
  for (const n of nodes) {
    if (n.key === key) return n
    if (n.children) {
      const found = findNode(n.children, key)
      if (found) return found
    }
  }
  return null
}

// ── 新增 / 修改 共用弹窗 ────────────────────────────────
const modalVisible = ref(false);
const formRef = ref()
const form = reactive<any>({});
const isEdit = ref(false);
const submitLoading = ref(false)

function handleAdd() {
  isEdit.value = false;
  resetForm('');
  modalVisible.value = true;
}

function handleEdit(row: any) {
  isEdit.value = true;
  Object.assign(form, { id: row.id });
  resetForm(row.id);
  modalVisible.value = true;
}

async function resetForm(id:string) {
  await nextTick()
  formRef.value?.clearValidate()
  checkedMenuKeys.value = []
  checkedDeptKeys.value = []
  if (isEdit.value) {
    const res: any = await roleApi.getById(id)
    Object.assign(form, res);
    form.menuCheckStrictly ??= 1
    form.deptCheckStrictly ??= 1
    // 父子联动开启时只回显叶子节点；关闭时保留后端返回的完整勾选集合。
    const menuIds: any = await roleApi.getMenuIds(String(form.id))
    const allIds = (menuIds ?? []).map((id: any) => String(id))
    if (form.menuCheckStrictly === 1) {
      const idSet = new Set(allIds)
      checkedMenuKeys.value = allIds.filter((id: string) => {
        const node = findNode(menuTreeData.value, id)
        return !node?.children?.length || !node.children.some((c: any) => idSet.has(c.key))
      })
    } else {
      checkedMenuKeys.value = allIds
    }
    const deptIds: any = await roleApi.getDeptIds(String(form.id))
    const allDeptIds = (deptIds ?? []).map((deptId: any) => String(deptId))
    if (form.deptCheckStrictly === 1) {
      const deptIdSet = new Set(allDeptIds)
      checkedDeptKeys.value = allDeptIds.filter((id: string) => {
        const node = findNode(deptTreeData.value, id)
        return !node?.children?.length || !node.children.some((c: any) => deptIdSet.has(c.key))
      })
    } else {
      checkedDeptKeys.value = allDeptIds
    }
  } else {
    Object.assign(form, {
      roleName: undefined,
      roleKey: undefined,
      sort: 0,
      dataScope: undefined,
      menuCheckStrictly: 1,
      deptCheckStrictly: 1,
      status: 0,
      remark: undefined,
      deptIds: [],
    });
  }
}

async function submitForm() {
  try {
    const valid = await formRef.value.validate()
    if (valid) return false
    submitLoading.value = true
    // 父子联动开启时只传叶节点，后端自动补齐祖先；关闭时保留全部选中值。
    const normalizedMenuKeys = toKeyStrings(checkedMenuKeys.value)
    const menuIds = (form.menuCheckStrictly === 1
      ? normalizedMenuKeys.filter((key: string) => {
          const node = findNode(menuTreeData.value, key)
          return !node?.children?.length
        })
      : normalizedMenuKeys
    ).map(Number)
    const deptIds = form.dataScope === 1 ? toKeyStrings(checkedDeptKeys.value).map(Number) : []
    if (isEdit.value) {
      await roleApi.update(form.id!, { ...form, menuIds, deptIds })
      Message.success(t('common.editSuccess'));
    } else {
      await roleApi.add({ ...form, menuIds, deptIds })
      Message.success(t('common.addSuccess'));
    }
    modalVisible.value = false;
    loadPage();
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.role.confirmDelete', { name: row.roleName }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await roleApi.delete([row.id]);
      Message.success(t('common.deleteSuccess'));
      loadPage();
    },
  });
}

loadPage();
loadMenuTree();
loadDeptTree();
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'sys:role$add'" type="primary" @click="handleAdd">
            <template #icon><IconPlus /></template>
            {{ t('common.add') }}
          </a-button>
          <a-button
            v-permission="'sys:role$delete'"
            status="danger"
            :disabled="!selectedIds.length"
            @click="batchDelete"
          >
            <template #icon><IconDelete /></template>
            {{ t('common.batchDelete') }}
          </a-button>
        </a-space>

        <div class="cl-toolbar-right">
          <a-input-search
            v-model="queryParams.roleName"
            :placeholder="t('system.role.searchPlaceholder')"
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
                <a-input v-model="queryParams.roleKey" :placeholder="t('system.role.roleKeyFilterPlaceholder')" allow-clear />
                <a-select v-model="queryParams.status" :placeholder="t('common.status')" allow-clear>
                  <a-option v-for="dict in dicts.sys_common_status" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-option>
                </a-select>
                <div class="cl-filter-actions">
                  <a-button size="small" @click="handleReset">{{ t('common.reset') }}</a-button>
                  <a-button size="small" type="primary" @click="() => { handleSearch(); advancedVisible = false }">{{ t('common.search') }}</a-button>
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
          @page-change="(page: number) => { queryParams.pageNum = page; loadPage() }"
          @page-size-change="(size: number) => { queryParams.pageSize = size; queryParams.pageNum = 1; loadPage() }"
        >
          <template #columns>
            <a-table-column :title="t('system.role.roleName')" data-index="roleName" />
            <a-table-column :title="t('system.role.roleKey')" data-index="roleKey" />
            <a-table-column :title="t('system.role.dataScope')" data-index="sort" align="center" :width="200" >
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_data_scope" :value="record.dataScope" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.status')" data-index="status" align="center" :width="100">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_common_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.remark')" data-index="remark" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('common.action')" :width="160" align="center">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button v-permission="'sys:role$edit'" size="mini" type="text" @click.stop="handleEdit(record)">
                    <template #icon><IconEdit /></template>{{ t('common.edit') }}
                  </a-button>
                  <a-button v-permission="'sys:role$delete'" size="mini" type="text" status="danger" @click.stop="handleDelete(record)">
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
            <a-checkbox :model-value="mobileAllSelected" :indeterminate="mobileIndeterminate" @change="toggleMobileSelectAll">{{ t('common.selectAll') }}</a-checkbox>
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
              <a-checkbox v-if="mobileSelectMode" :model-value="isCardSelected(item.id)" @click.stop @change="toggleCardSelect(item.id)" class="cl-card-checkbox" />
              <div class="cl-card-identity">
                <strong>{{ item.roleName }}</strong>
                <span class="cl-card-sub">{{ item.roleKey || '-' }}</span>
              </div>
              <DictTag :options="dicts.sys_common_status" :value="item.status" dot />
            </div>
            <div class="cl-card-meta">
              <span>{{ t('common.sort') }}：{{ item.sort }}</span>
              <span>{{ t('common.remark') }}：{{ item.remark || '-' }}</span>
            </div>
            <div v-if="!mobileSelectMode" class="cl-card-footer">
              <a-space size="mini">
                <a-button v-permission="'sys:role$edit'" size="mini" type="outline" @click="handleEdit(item)">
                  <template #icon><IconEdit /></template>{{ t('common.edit') }}
                </a-button>
                <a-button v-permission="'sys:role$delete'" size="mini" type="outline" status="danger" @click="handleDelete(item)">
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

    <!-- ── 新增 / 修改 弹窗 ── -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? t('system.role.editTitle') : t('system.role.addTitle')"
      :width="isMobile ? '100%' : 620"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.role.roleName')" field="roleName"
              :rules="[{ required: true, message: t('system.role.roleNameRequired') }, { maxLength: 50, message: t('common.maxLength', { max: 50 }) }]"
              :validate-trigger="['blur']">
              <a-input v-model="form.roleName" :placeholder="t('system.role.roleNamePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.role.roleKey')" field="roleKey"
              :rules="[{ maxLength: 32, message: t('common.maxLength', { max: 32 }) }]"
              :validate-trigger="['blur']">
              <a-input v-model="form.roleKey" :placeholder="t('system.role.roleKeyPlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.role.dataScope')" field="dataScope"
              :rules="[{ required: true, message: t('system.role.dataScopeRequired') }]"
              :validate-trigger="['blur']">
              <a-select v-model="form.dataScope" size="small" style="width:100%">
                <a-option v-for="q in dicts.sys_data_scope" :key="q.dictValue" :value="Number(q.dictValue)">{{ q.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col v-if="form.dataScope === 1" :span="24">
            <a-form-item :label="t('system.role.deptPerms')" :content-flex="false" :merge-props="false" style="width:100%">
              <div class="tree-option-head">
                <span>{{ t('system.role.deptCheckStrictly') }}</span>
                <a-switch v-model="form.deptCheckStrictly" :checked-value="1" :unchecked-value="0" />
              </div>
              <a-tree-select
                v-model="checkedDeptKeys"
                :data="deptTreeData"
                :placeholder="t('system.user.deptPlaceholder')"
                multiple
                tree-checkable
                :tree-check-strictly="form.deptCheckStrictly !== 1"
                allow-clear
                :fallback-option="false"
                style="width:100%"
              />
            </a-form-item>
          </a-col>
          <!-- <a-col :span="isMobile ? 24 : 12">
            <a-form-item label="排序" field="sort"
              :rules="[{ required: true, message: '排序不能为空' }]"
              :validate-trigger="['blur']">
              <a-input-number v-model="form.sort" placeholder="请输入排序" style="width:100%" :precision="0" />
            </a-form-item>
          </a-col> -->
          <a-col :span="24">
            <a-form-item :label="t('common.remark')" field="remark"
              :rules="[{ maxLength: 500, message: t('common.maxLength', { max: 500 }) }]"
              :validate-trigger="['blur']">
              <a-textarea v-model="form.remark" :placeholder="t('system.role.remarkPlaceholder')" allow-clear />
            </a-form-item>
          </a-col>

          <!-- 菜单功能权限树 -->
          <a-col :span="24">
            <a-form-item :label="t('system.role.menuPerms')" :content-flex="false" :merge-props="false" style="width:100%">
              <div class="tree-option-head">
                <span>{{ t('system.role.menuCheckStrictly') }}</span>
                <a-switch v-model="form.menuCheckStrictly" :checked-value="1" :unchecked-value="0" />
              </div>
              <div style="border: 1px solid var(--color-border-2); border-radius: 4px; padding: 8px; max-height: 300px; overflow-y: auto; width: 100%;">
                <a-tree
                  v-model:checked-keys="checkedMenuKeys"
                  v-model:expanded-keys="menuExpandedKeys"
                  :data="menuTreeData"
                  checkable
                  :check-strictly="form.menuCheckStrictly !== 1"
                  block-node
                  style="width: 100%"
                />
                <a-empty v-if="!menuTreeData.length" :description="t('common.noData')" />
              </div>
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.status')" field="status">
              <a-radio-group v-model="form.status" type="button">
                <a-radio v-for="dict in dicts.sys_common_status" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
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

<style scoped>
.tree-option-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
  color: var(--color-text-2);
  font-size: 12px;
}
</style>
