<script setup lang="ts">
import { IconDelete, IconEdit, IconPlus } from '@arco-design/web-vue/es/icon';
import deptApi from "@/api/system/dept"
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any
const dicts = proxy.$useDict('sys_dept_type', 'sys_common_status')

const isMobile = ref(window.innerWidth < 576)
function onResize() { isMobile.value = window.innerWidth < 576 }
onMounted(() => window.addEventListener('resize', onResize))
onUnmounted(() => window.removeEventListener('resize', onResize))

// ── 列表 ──────────────────────────────────────────────
const loading = ref(false)
const rawList = ref<any[]>([])
const treeData = ref<any[]>([])
const expandedKeys = ref<string[]>([])
const isFirstLoad = ref(true)
const queryParams = ref({ deptName: '', status: '' as any })

function buildTree(list: any[], parentId = '0'): any[] {
  return list
    .filter(item => String(item.parentId) === String(parentId))
    .map(item => {
      const children = buildTree(list, item.id)
      return children.length ? { ...item, children } : { ...item }
    })
    .sort((a, b) => a.sort - b.sort)
}

async function loadList() {
  loading.value = true
  try {
    const res: any = await deptApi.list(queryParams.value)
    rawList.value = res ?? []
    treeData.value = buildTree(rawList.value)
    // 只在首次加载时初始化展开状态
    if (isFirstLoad.value) {
      expandedKeys.value = rawList.value
        .filter(i => String(i.parentId) === '0')
        .map(i => String(i.id))
      isFirstLoad.value = false
    }
  } finally {
    loading.value = false
  }
}

function handleSearch() { loadList() }
function handleReset() {
  queryParams.value.deptName = ''
  queryParams.value.status = ''
  loadList()
}


// ── 弹窗 ──────────────────────────────────────────────
const parentTreeData = ref<any[]>([])
const modalVisible = ref(false)
const formRef = ref()
const form = reactive<any>({})
const isEdit = ref(false)
const submitLoading = ref(false)

function handleAdd(row?: any) {
  isEdit.value = false
  parentTreeData.value = buildTree(rawList.value)
  Object.assign(form, {
    parentId: row ? String(row.id) : undefined,
    deptType: undefined,
    deptName: '',
    sort: 0,
    status: 0,
    remark: '',
  })
  nextTick(() => formRef.value?.clearValidate())
  modalVisible.value = true
}

function handleEdit(row: any) {
  isEdit.value = true
  parentTreeData.value = buildTree(rawList.value.filter(i => String(i.id) !== String(row.id)))
  Object.assign(form, {
    id: row.id,
    parentId: row.parentId === '0' ? '' : row.parentId,
    deptType: row.deptType,
    deptName: row.deptName,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  nextTick(() => formRef.value?.clearValidate())
  modalVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value.validate()
  if (valid) return
  submitLoading.value = true
  try {
    if (form.parentId == undefined || form.parentId == ''){
        form.parentId = '0'
    }
    const payload = { ...form, parentId: form.parentId }
    if (isEdit.value) {
      await deptApi.update(String(form.id), payload)
      Message.success(t('common.editSuccess'))
    } else {
      await deptApi.add(payload)
      Message.success(t('common.addSuccess'))
    }
    modalVisible.value = false
    loadList()
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.dept.confirmDelete', { name: row.deptName }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await deptApi.delete(row.id)
      Message.success(t('common.deleteSuccess'))
      loadList()
    },
  })
}

loadList()
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'sys:dept$add'" type="primary" @click="handleAdd()">
            <template #icon><IconPlus /></template>{{ t('common.add') }}
          </a-button>
        </a-space>
        <div class="cl-toolbar-right">
          <a-input-search
            v-model="queryParams.deptName"
            :placeholder="t('system.dept.deptName')"
            allow-clear
            class="cl-toolbar-search"
            @search="handleSearch"
            @press-enter="handleSearch"
          />
          <template v-if="!isMobile">
            <a-select v-model="queryParams.status" :placeholder="t('common.status')" allow-clear style="width:100px" @change="handleSearch">
              <a-option v-for="d in dicts.sys_common_status" :key="d.value" :value="Number(d.value)">{{ d.label }}</a-option>
            </a-select>
            <a-button @click="handleReset">{{ t('common.reset') }}</a-button>
          </template>
        </div>
      </div>

      <!-- 树形表格 -->
      <div>
        <a-skeleton v-if="loading" :animation="true"><a-skeleton-line :rows="8" /></a-skeleton>
        <a-table
          v-else
          :bordered="false"
          :data="treeData"
          row-key="id"
          :pagination="false"
          :scroll="{ x: 800 }"
          v-model:expandedKeys="expandedKeys"
        >
          <template #columns>
            <a-table-column :title="t('system.dept.deptName')" data-index="deptName" :width="240" />
            <a-table-column :title="t('system.dept.deptType')" data-index="deptType" align="center" :width="100">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_dept_type" :value="record.deptType" />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.sort')" data-index="sort" align="center" :width="80" />
            <a-table-column :title="t('common.status')" data-index="status" align="center" :width="100">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_common_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.remark')" data-index="remark" :width="200" :ellipsis="true" :tooltip="true" />
            <a-table-column :title="t('common.action')" :width="100" align="center" fixed="right">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button v-permission="'sys:dept$add'" size="mini" type="text" @click.stop="handleAdd(record)">
                    <template #icon><IconPlus /></template>
                  </a-button>
                  <a-button v-permission="'sys:dept$edit'" size="mini" type="text" @click.stop="handleEdit(record)">
                    <template #icon><IconEdit /></template>
                  </a-button>
                  <a-button v-permission="'sys:dept$delete'" size="mini" type="text" status="danger" @click.stop="handleDelete(record)">
                    <template #icon><IconDelete /></template>
                  </a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>
    </a-card>

    <!-- 新增 / 修改 弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? t('system.dept.editTitle') : t('system.dept.addTitle')"
      :width="isMobile ? '100%' : 520"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item :label="t('system.dept.parentDept')" field="parentId">
              <a-tree-select
                v-model="form.parentId"
                :data="parentTreeData"
                :field-names="{ key: 'id', title: 'deptName' }"
                :placeholder="t('system.dept.parentDeptPlaceholder')"
                allow-clear
                style="width:100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.dept.deptType')" field="deptType"
              :rules="[{ required: true, message: t('system.dept.deptTypeRequired') }]"
              :validate-trigger="['change']">
              <a-select v-model="form.deptType" :placeholder="t('system.dept.deptTypePlaceholder')" style="width:100%">
                <a-option v-for="dict in dicts.sys_dept_type" :key="dict.dictValue" :value="Number(dict.dictValue)">{{ dict.label }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('system.dept.deptName')" field="deptName"
              :rules="[{ required: true, message: t('system.dept.deptNameRequired') }, { maxLength: 30, message: t('common.maxLength', { max: 30 }) }]"
              :validate-trigger="['blur']">
              <a-input v-model="form.deptName" :placeholder="t('system.dept.deptNamePlaceholder')" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('common.remark')" field="remark"
              :rules="[{ maxLength: 255, message: t('common.maxLength', { max: 255 }) }]"
              :validate-trigger="['blur']">
              <a-textarea v-model="form.remark" :placeholder="t('system.dept.remarkPlaceholder')" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.sort')" field="sort">
              <a-input-number v-model="form.sort" :min="0" style="width:100%" />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.status')" field="status">
              <a-switch v-model="form.status" :checked-value="0" :unchecked-value="1" :checked-text="t('system.dept.statusNormal')" :unchecked-text="t('system.dept.statusDisabled')" />
            </a-form-item>
          </a-col>
          <a-col :span="24" style="margin-top:16px">
            <div style="display:flex; justify-content:center; gap:12px">
              <a-button @click="modalVisible = false">{{ t('common.cancel') }}</a-button>
              <a-button type="primary" :loading="submitLoading" @click="submitForm">{{ t('common.save') }}</a-button>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>
