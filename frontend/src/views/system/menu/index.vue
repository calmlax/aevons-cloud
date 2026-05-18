<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import menuApi from '@/api/system/menu'
import langApi from '@/api/system/lang'
import IconPicker from '@/components/IconPicker/index.vue';

const { t } = useI18n()

const proxy = getCurrentInstance()!.proxy as any

const dicts = proxy.$useDict('sys_is','sys_visible','sys_common_status','sys_menu_type')


// ── 数据加载 ───────────────────────────────────────────
const loading = ref(false)
const rawList = ref<any[]>([])
const keyword = ref('')
const langList = ref<any[]>([])
const isMobile = ref(window.innerWidth < 576)

async function loadList() {
  loading.value = true
  try {
    const res: any = await menuApi.list()
    rawList.value = res ?? []
    await nextTick()
    refreshExpandedKeys()
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  loadList()
  langList.value = await langApi.availableList() as any
})

// ── 树形构建 ───────────────────────────────────────────
function buildTree(items: any[], parentId = 0): any[] {
  return items
    .filter(i => Number(i.parentId) === parentId)
    .sort((a, b) => Number(a.sort) - Number(b.sort))
    .map(i => {
      const children = buildTree(items, Number(i.id))
      return children.length ? { ...i, children } : { ...i }
    })
}

const treeData = computed(() => {
  const kw = keyword.value.trim().toLowerCase()
  if (!kw) return buildTree(rawList.value)
  const matched = rawList.value.filter(i =>
    (i.title ?? '').toLowerCase().includes(kw) ||
    (i.path ?? '').includes(kw) ||
    (i.permission ?? '').includes(kw)
  )
  return buildTree(matched)
})

// ── 展开控制（默认展开一级）──────────────────────────────
const expandedKeys = ref<string[]>([])
const isFirstLoad = ref(true)

function refreshExpandedKeys() {
  if (!isFirstLoad.value) return
  expandedKeys.value = treeData.value.map(node => String(node.id))
  isFirstLoad.value = false
}

// ── 父级选项（树形 select）─────────────────────────────
const parentOptions = computed(() => [
  ...buildSelectTree(rawList.value.filter(i => i.type !== 3)),
])

function buildSelectTree(items: any[], parentId = 0): any[] {
  return items
    .filter(i => Number(i.parentId) === parentId)
    .sort((a, b) => Number(a.sort) - Number(b.sort))
    .map(i => {
      const children = buildSelectTree(items, Number(i.id))
      return { key: String(i.id), value: String(i.id), title: i.title, ...(children.length ? { children } : {}) }
    })
}

// ── 弹窗 ───────────────────────────────────────────────
const modalVisible = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref()
const form = reactive<any>({})

function emptyForm() {
  return {
    parentId: '0', title: '', type: 2, sort: 1,
    path: '', component: '', query: '',
    visible: 1, status: 0, isFrame: 0,
    permission: '', icon: '', activeId: undefined,
  }
}

function initTranslations() {
  if (!form.translations) form.translations = {}
  langList.value.forEach(lang => {
    if (!form.translations[lang.langCode]) {
      form.translations[lang.langCode] = ''
    }
  })
}

function handleAdd(parentId = '0') {
  isEdit.value = false
  const base = emptyForm()
  Object.keys(base).forEach(k => { form[k] = (base as any)[k] })
  form.parentId = parentId
  form.type = parentId === '0' ? 1 : 2
  form.translations = {}
  initTranslations()
  nextTick(() => { formRef.value?.clearValidate?.() })
  modalVisible.value = true
}

async function handleEdit(record: any) {
  isEdit.value = true
  const base = emptyForm()
  Object.keys(base).forEach(k => { form[k] = (base as any)[k] })
  const res: any = await menuApi.getById(record.id)
  Object.keys(res).forEach(k => {
    if (k !== 'translations') form[k] = res[k]
  })
  form.translations = {}
  initTranslations()
  if (res.translations) {
    Object.entries(res.translations).forEach(([code, val]: any) => {
      form.translations[code] = typeof val === 'string' ? val : (val?.title ?? '')
    })
  }
  nextTick(() => { formRef.value?.clearValidate?.() })
  modalVisible.value = true
}

async function submitForm() {
  const valid = await formRef.value.validate()
  if (valid) return
  submitLoading.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) {
      await menuApi.update(form.id, payload)
      Message.success(t('common.editSuccess'))
    } else {
      await menuApi.add(payload)
      Message.success(t('common.addSuccess'))
    }
    modalVisible.value = false
    loadList()
  } finally {
    submitLoading.value = false
  }
}

function handleDelete(record: any) {
  Modal.confirm({
    title: t('common.confirmDelete'),
    content: t('system.menu.confirmDelete', { name: record.title }),
    okButtonProps: { status: 'danger' },
    onOk: async () => {
      await menuApi.delete(record.id)
      Message.success(t('common.deleteSuccess'))
      loadList()
    },
  })
}
</script>

<template>
  <div class="page-stack">
    <a-card class="panel-card" :bordered="false">

      <!-- 工具栏 -->
      <div class="cl-toolbar">
        <a-space>
          <a-button v-permission="'sys:menu$add'" type="primary" @click="handleAdd()">
            <template #icon><SvgIcon name="plus" /></template>
            {{ t('common.add') }}
          </a-button>
        </a-space>
        <div class="cl-toolbar-right">
          <a-input-search
            v-model="keyword"
            :placeholder="t('system.menu.searchPlaceholder')"
            allow-clear
            class="cl-toolbar-search"
          />
        </div>
      </div>

      <!-- 树形表格 -->
      <div>
        <a-skeleton v-if="loading" :animation="true">
          <a-skeleton-line :rows="8" />
        </a-skeleton>
        <a-table
          v-else
          :bordered="false"
          :data="treeData"
          row-key="id"
          :pagination="false"
          :scroll="{ x: 1000 }"
          v-model:expanded-keys="expandedKeys"
          size="medium"
        >
          <template #columns>
            <a-table-column :title="t('system.menu.colName')" data-index="title" :width="300">
              <template #cell="{ record }">
                <span class="mm-name-cell">
                  <SvgIcon v-if="record.icon" :name="record.icon" :size="14" class="mm-icon" />
                  <span v-if="record.title">{{ record.title }}</span>
                  <span v-else class="mm-dash">—</span>
                </span>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.menu.colType')" :width="80" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_menu_type" :value="record.type" />
              </template>
            </a-table-column>
            <a-table-column :title="t('system.menu.colPath')" data-index="path" :width="200" :ellipsis="true" :tooltip="true">
              <template #cell="{ record }">
                <code v-if="record.path" class="mm-code">{{ record.path }}</code>
                <span v-else class="mm-dash">—</span>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.menu.colComponent')" data-index="component" :width="200" :ellipsis="true" :tooltip="true">
              <template #cell="{ record }">
                <code v-if="record.component" class="mm-code">{{ record.component }}</code>
                <span v-else class="mm-dash">—</span>
              </template>
            </a-table-column>
            <a-table-column :title="t('system.menu.colPermission')" data-index="permission" :width="200" :ellipsis="true" :tooltip="true">
              <template #cell="{ record }">
                <code v-if="record.permission" class="mm-code mm-code-sign">{{ record.permission }}</code>
                <span v-else class="mm-dash">—</span>
              </template>
            </a-table-column>
            <a-table-column :title="t('common.sort')" data-index="sort" :width="80" align="center" />
            <a-table-column :title="t('system.menu.colVisible')" :width="80" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_visible" :value="record.visible" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.status')" :width="80" align="center">
              <template #cell="{ record }">
                <DictTag :options="dicts.sys_common_status" :value="record.status" dot />
              </template>
            </a-table-column>
            <a-table-column :title="t('common.action')" :width="100" align="center" fixed="right">
              <template #cell="{ record }">
                <a-space size="mini">
                  <a-button
                    v-if="record.type !== 3"
                    v-permission="'sys:menu$add'"
                    size="mini"
                    type="text"
                    :title="t('system.menu.addChild')"
                    @click.stop="handleAdd(String(record.id))"
                  >
                    <template #icon><SvgIcon name="plus" /></template>
                  </a-button>
                  <a-button
                    v-permission="'sys:menu$edit'"
                    size="mini"
                    type="text"
                    :title="t('system.menu.editMenu')"
                    @click.stop="handleEdit(record)"
                  >
                    <template #icon><SvgIcon name="edit" /></template>
                  </a-button>
                  <a-button
                    v-permission="'sys:menu$delete'"
                    size="mini"
                    type="text"
                    status="danger"
                    :title="t('system.menu.deleteMenu')"
                    @click.stop="handleDelete(record)"
                  >
                    <template #icon><SvgIcon name="delete" /></template>
                  </a-button>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>
    </a-card>

    <!-- 新增 / 编辑弹窗 -->
    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? t('system.menu.editTitle') : t('system.menu.addTitle')"
      :width="isMobile ? '100%' : 580"
      :fullscreen="isMobile"
      :footer="false"
    >
      <a-form :model="form" ref="formRef" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item :label="t('system.menu.parentMenu')" field="parentId">
              <a-tree-select
                v-model="form.parentId"
                :data="parentOptions"
                :placeholder="t('system.menu.parentMenuPlaceholder')"
                allow-clear
                :fallback-option="false"
              />
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.menu.menuType')" field="type" :rules="[{ required: true }]">
              <a-radio-group v-model="form.type" type="button">
                <a-radio v-for="dict in dicts.sys_menu_type" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 12">
            <a-form-item :label="t('common.sort')" field="sort">
              <a-input-number v-model="form.sort" :placeholder="t('common.pleaseInput')" />
            </a-form-item>
          </a-col>

          <a-col :span="24">
            <a-form-item :label="t('system.menu.menuName')" :content-flex="false" :merge-props="false">
              <a-space direction="vertical" fill>
                <a-form-item
                  hide-label
                  v-for="lang in langList"
                  :key="lang.langCode"
                  :field="`translations.${lang.langCode}`"
                  :rules="[{ required: true, message: t('system.menu.menuNameRequired', { lang: lang.langName }) },{ maxLength: 50, message: t('system.menu.menuNameMaxLength') }]"
                  :validate-trigger="['blur']"
                >
                  <a-input
                    v-if="form.translations?.[lang.langCode] !== undefined"
                    v-model="form.translations[lang.langCode]"
                    :placeholder="t('system.menu.menuNamePlaceholder')"
                  >
                    <template #prepend>
                      <div style="width: 100px; text-align: center;">{{ lang.langName }}</div>
                    </template>
                  </a-input>
                </a-form-item>
              </a-space>
            </a-form-item>
          </a-col>

          <a-col v-if="form.type !== 3" :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.menu.icon')" field="icon">
                 <IconPicker
                v-model="form.icon"
                placeholder=""
              />
            </a-form-item>
          </a-col>
          <a-col v-if="form.type !== 3" :span="isMobile ? 24 : 12">
            <a-form-item :label="t('system.menu.isFrame')" field="isFrame">
              <a-switch v-model="form.isFrame" :checked-value="1" :unchecked-value="0" :checked-text="t('system.menu.isFrameYes')" :unchecked-text="t('system.menu.isFrameNo')" />
            </a-form-item>
          </a-col>

          <template v-if="form.type === 1 || form.type === 2">
            <a-col :span="24">
              <a-form-item :label="t('system.menu.routePath')" field="path"
                :rules="[{ required: form.type === 2, message: t('system.menu.routePathRequired') },{ maxLength: 100, message: t('system.menu.maxLength100') }]"
                :validate-trigger="['blur']"
              >
                <a-input v-model="form.path" :placeholder="t('system.menu.routePathPlaceholder')" allow-clear />
              </a-form-item>
            </a-col>
          </template>

          <template v-if="form.type === 2">
            <a-col :span="24">
              <a-form-item :label="t('system.menu.componentPath')" field="component"
                :rules="[{ required: true, message: t('system.menu.componentPathRequired') },{ maxLength: 100, message: t('system.menu.maxLength100') }]"
                :validate-trigger="['blur']"
              >
                <a-input v-model="form.component" :placeholder="t('system.menu.componentPathPlaceholder')" allow-clear />
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item :label="t('system.menu.routeQuery')" field="query">
                <a-input v-model="form.query" placeholder='{"id":1}' allow-clear />
              </a-form-item>
            </a-col>
          </template>

          <a-col :span="24">
            <a-form-item :label="t('system.menu.permission')" field="permission"
              :rules="form.type === 3 ? [{ required: true, message: t('system.menu.permissionRequired') },{ maxLength: 32, message: t('system.menu.maxLength32') }] : []"
              :validate-trigger="['blur']"
            >
              <a-input v-model="form.permission" placeholder="system:menu$list" allow-clear />
            </a-form-item>
          </a-col>

          <a-col v-if="form.type === 2" :span="24">
            <a-form-item :label="t('system.menu.activeMenu')" field="activeId" :extra="t('system.menu.activeMenuExtra')">
              <a-tree-select
                v-model="form.activeId"
                :data="parentOptions"
                :placeholder="t('system.menu.activeMenuPlaceholder')"
                allow-clear
                :fallback-option="false"
                style="width:100%"
              />
            </a-form-item>
          </a-col>

          <a-col :span="isMobile ? 24 : 8">
            <a-form-item v-if="form.type !== 3" :label="t('system.menu.visible')" field="visible">
              <a-radio-group v-model="form.visible" type="button">
                <a-radio v-for="dict in dicts.sys_visible" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
          <a-col :span="isMobile ? 24 : 8">
            <a-form-item :label="t('common.status')" field="status">
              <a-radio-group v-model="form.status" type="button">
                <a-radio v-for="dict in dicts.sys_common_status" :value="Number(dict.dictValue)">{{ dict.label }}</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>

          <a-col :span="24" style="margin-top: 30px;">
            <div style="display:flex;justify-content:center;gap:12px;">
              <a-button @click="modalVisible = false">{{ t('common.cancel') }}</a-button>
              <a-button type="primary" :loading="submitLoading" @click="submitForm()">{{ t('common.save') }}</a-button>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>
