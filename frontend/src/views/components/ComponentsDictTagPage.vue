<script setup lang="ts">
import { ref } from 'vue';
import { DictTag } from '@/components/DictTag';
// 导入改造后的类型（dictValue / text）
import type { DictOption } from '@/components/DictTag/types';

// 用户状态字典
const userStatusOptions: DictOption[] = [
  { dictValue: 'active', text: '活跃', tagType: 'green', id: 1, dictType: 'user_status', langCode: 'zh-CN' },
  { dictValue: 'recall', text: '待召回', tagType: 'orange', id: 2, dictType: 'user_status', langCode: 'zh-CN' },
  { dictValue: 'frozen', text: '冻结', tagType: 'red', id: 3, dictType: 'user_status', langCode: 'zh-CN' },
];

// 订单状态字典
const orderStatusOptions: DictOption[] = [
  { dictValue: 'pending', text: '待付款', tagType: 'orange', id: 1, dictType: 'order_status', langCode: 'zh-CN' },
  { dictValue: 'paid', text: '已付款', tagType: 'blue', id: 2, dictType: 'order_status', langCode: 'zh-CN' },
  { dictValue: 'shipped', text: '已发货', tagType: 'cyan', id: 3, dictType: 'order_status', langCode: 'zh-CN' },
  { dictValue: 'completed', text: '已完成', tagType: 'green', id: 4, dictType: 'order_status', langCode: 'zh-CN' },
  { dictValue: 'cancelled', text: '已取消', tagType: 'red', id: 5, dictType: 'order_status', langCode: 'zh-CN' },
  { dictValue: 'refunding', text: '退款中', tagType: 'purple', id: 6, dictType: 'order_status', langCode: 'zh-CN' },
];

// 工单优先级字典
const priorityOptions: DictOption[] = [
  { dictValue: 'low', text: '低', tagType: 'info', id: 1, dictType: 'priority', langCode: 'zh-CN' },
  { dictValue: 'medium', text: '中', tagType: 'warning', id: 2, dictType: 'priority', langCode: 'zh-CN' },
  { dictValue: 'high', text: '高', tagType: 'danger', id: 3, dictType: 'priority', langCode: 'zh-CN' },
  { dictValue: 'critical', text: '紧急', tagType: 'red', id: 4, dictType: 'priority', langCode: 'zh-CN' },
];

// 审核状态字典
const auditOptions: DictOption[] = [
  { dictValue: 0, text: '待审核', tagType: 'orange', id: 1, dictType: 'audit_status', langCode: 'zh-CN' },
  { dictValue: 1, text: '已通过', tagType: 'green', id: 2, dictType: 'audit_status', langCode: 'zh-CN' },
  { dictValue: 2, text: '已拒绝', tagType: 'red', id: 3, dictType: 'audit_status', langCode: 'zh-CN' },
];

// 交互演示
const selectedStatus = ref<string>('active');
const selectedOrder = ref<string>('shipped');

// 表格演示数据
const tableData = [
  { id: 'T001', name: '林若琪', status: 'active', order: 'completed', priority: 'low' },
  { id: 'T002', name: '陈致远', status: 'recall', order: 'pending', priority: 'high' },
  { id: 'T003', name: '王知夏', status: 'frozen', order: 'refunding', priority: 'critical' },
  { id: 'T004', name: '赵明哲', status: 'active', order: 'shipped', priority: 'medium' },
];
</script>

<template>
  <div class="page-stack">
    <!-- Hero -->
    <section class="hero-panel">
      <div>
        <p class="hero-kicker">组件目录</p>
        <h2>DictTag 字典标签</h2>
        <p class="hero-copy">
          基于字典选项渲染语义化标签，支持预设颜色、圆点指示器和三种尺寸，适合状态列、审核结果、优先级等场景。
        </p>
      </div>
    </section>

    <!-- 基础用法 -->
    <a-card class="panel-card" :bordered="false" title="基础用法">
      <div class="aevo-section">
        <p class="aevo-label">用户状态</p>
        <div class="aevo-row">
          <DictTag v-for="opt in userStatusOptions" :key="String(opt.dictValue)" :options="userStatusOptions" :value="opt.dictValue" />
        </div>
      </div>

      <a-divider />

      <div class="aevo-section">
        <p class="aevo-label">订单状态</p>
        <div class="aevo-row">
          <DictTag v-for="opt in orderStatusOptions" :key="String(opt.dictValue)" :options="orderStatusOptions" :value="opt.dictValue" />
        </div>
      </div>

      <a-divider />

      <div class="aevo-section">
        <p class="aevo-label">工单优先级</p>
        <div class="aevo-row">
          <DictTag v-for="opt in priorityOptions" :key="String(opt.dictValue)" :options="priorityOptions" :value="opt.dictValue" />
        </div>
      </div>

      <a-divider />

      <div class="aevo-section">
        <p class="aevo-label">审核状态（数字值）</p>
        <div class="aevo-row">
          <DictTag v-for="opt in auditOptions" :key="String(opt.dictValue)" :options="auditOptions" :value="opt.dictValue" />
        </div>
      </div>
    </a-card>

    <!-- 圆点模式 -->
    <a-card class="panel-card" :bordered="false" title="圆点指示器">
      <div class="aevo-section">
        <p class="aevo-label">dot 模式，适合列表行内状态</p>
        <div class="aevo-row">
          <DictTag v-for="opt in userStatusOptions" :key="String(opt.dictValue)" :options="userStatusOptions" :value="opt.dictValue" dot />
        </div>
      </div>
      <a-divider />
      <div class="aevo-section">
        <p class="aevo-label">订单状态 dot</p>
        <div class="aevo-row">
          <DictTag v-for="opt in orderStatusOptions" :key="String(opt.dictValue)" :options="orderStatusOptions" :value="opt.dictValue" dot />
        </div>
      </div>
    </a-card>

    <!-- 尺寸 -->
    <a-card class="panel-card" :bordered="false" title="尺寸">
      <div class="aevo-section">
        <div class="aevo-size-row">
          <div class="aevo-size-item">
            <span class="aevo-size-label">small</span>
            <DictTag :options="orderStatusOptions" value="completed" size="small" />
            <DictTag :options="orderStatusOptions" value="pending" size="small" dot />
          </div>
          <div class="aevo-size-item">
            <span class="aevo-size-label">medium（默认）</span>
            <DictTag :options="orderStatusOptions" value="completed" size="medium" />
            <DictTag :options="orderStatusOptions" value="pending" size="medium" dot />
          </div>
          <div class="aevo-size-item">
            <span class="aevo-size-label">large</span>
            <DictTag :options="orderStatusOptions" value="completed" size="large" />
            <DictTag :options="orderStatusOptions" value="pending" size="large" dot />
          </div>
        </div>
      </div>
    </a-card>

    <!-- 回退值 -->
    <a-card class="panel-card" :bordered="false" title="未匹配回退">
      <div class="aevo-section">
        <p class="aevo-label">值不在字典中时的处理</p>
        <div class="aevo-row">
          <DictTag :options="userStatusOptions" value="unknown" />
          <DictTag :options="userStatusOptions" value="unknown" fallback="未知状态" />
          <DictTag :options="userStatusOptions" :value="null" fallback="—" />
        </div>
        <p class="aevo-hint">第一个：直接显示原始值；第二个：使用 fallback 文字；第三个：null 值回退</p>
      </div>
    </a-card>

    <!-- 交互演示 -->
    <a-card class="panel-card" :bordered="false" title="交互演示">
      <div class="aevo-section">
        <div class="aevo-interactive">
          <div class="aevo-interactive-item">
            <p class="aevo-label">切换用户状态</p>
            <a-radio-group v-model="selectedStatus" type="button" size="small">
              <a-radio v-for="opt in userStatusOptions" :key="String(opt.dictValue)" :value="String(opt.dictValue)">
                {{ opt.text }}
              </a-radio>
            </a-radio-group>
            <div class="aevo-row" style="margin-top: 12px">
              <DictTag :options="userStatusOptions" :value="selectedStatus" size="large" dot />
            </div>
          </div>

          <div class="aevo-interactive-item">
            <p class="aevo-label">切换订单状态</p>
            <a-select v-model="selectedOrder" style="width: 160px" size="small">
              <a-option v-for="opt in orderStatusOptions" :key="String(opt.dictValue)" :value="String(opt.dictValue)">
                {{ opt.text }}
              </a-option>
            </a-select>
            <div class="aevo-row" style="margin-top: 12px">
              <DictTag :options="orderStatusOptions" :value="selectedOrder" size="large" />
            </div>
          </div>
        </div>
      </div>
    </a-card>

    <!-- 表格场景 -->
    <a-card class="panel-card" :bordered="false" title="表格场景">
      <a-table :data="tableData" :pagination="false" row-key="id" :scroll="{ x: 600 }">
        <template #columns>
          <a-table-column title="ID" data-index="id" :width="80" />
          <a-table-column title="姓名" data-index="name" />
          <a-table-column title="用户状态">
            <template #cell="{ record }">
              <DictTag :options="userStatusOptions" :value="record.status" dot />
            </template>
          </a-table-column>
          <a-table-column title="订单状态">
            <template #cell="{ record }">
              <DictTag :options="orderStatusOptions" :value="record.order" />
            </template>
          </a-table-column>
          <a-table-column title="优先级">
            <template #cell="{ record }">
              <DictTag :options="priorityOptions" :value="record.priority" size="small" />
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<style scoped>
.aevo-section {
  padding: 4px 0;
}

.aevo-label {
  margin: 0 0 10px;
  font-size: 13px;
  color: var(--text-subtle);
}

.aevo-hint {
  margin: 8px 0 0;
  font-size: 12px;
  color: var(--text-subtle);
  opacity: 0.7;
}

.aevo-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.aevo-size-row {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
}

.aevo-size-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.aevo-size-label {
  font-size: 12px;
  color: var(--text-subtle);
  min-width: 80px;
}

.aevo-interactive {
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
}

.aevo-interactive-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
