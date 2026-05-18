<script setup lang="ts">
/**
 * Cron 表达式生成器
 * 适配 robfig/cron/v3（WithSeconds 模式）：秒 分 时 日 月 周
 * 格式：S M H D Month Weekday
 */

const props = defineProps<{ modelValue?: string }>()
const emit = defineEmits<{ (e: 'update:modelValue', val: string): void }>()

const visible = ref(false)
const activeTab = ref('periodic')

// ── 当前预览表达式 ─────────────────────────────────────
const preview = ref(props.modelValue ?? '')
watch(() => props.modelValue, v => { if (v) preview.value = v })

// ── 周期触发 ───────────────────────────────────────────
const periodic = reactive({
  unit: 'minutes' as 'seconds' | 'minutes' | 'hours' | 'days',
  interval: 1,
  atMinute: 0,   // hours 模式：在第几分
  atTime: '00:00', // days 模式：在几点
})

function buildPeriodic(): string {
  const n = Math.max(1, periodic.interval)
  switch (periodic.unit) {
    case 'seconds': return `*/${n} * * * * *`
    case 'minutes': return `0 */${n} * * * *`
    case 'hours':   return `0 ${periodic.atMinute} */${n} * * *`
    case 'days': {
      const [h, m] = periodic.atTime.split(':').map(Number)
      return `0 ${m} ${h} */${n} * *`
    }
  }
}

// ── 日历触发 ───────────────────────────────────────────
const calendar = reactive({
  unit: 'daily' as 'daily' | 'weekly' | 'monthly',
  weekday: 1,   // 0=Sun … 6=Sat
  monthDay: 1,  // 1-31
  time: '00:00',
})

const weekdayOptions = [
  { label: '周日', value: 0 }, { label: '周一', value: 1 },
  { label: '周二', value: 2 }, { label: '周三', value: 3 },
  { label: '周四', value: 4 }, { label: '周五', value: 5 },
  { label: '周六', value: 6 },
]

// robfig/cron 周日=0，周一=1 … 周六=6（与 Quartz 不同，无需 +1）
function buildCalendar(): string {
  const [h, m] = calendar.time.split(':').map(Number)
  switch (calendar.unit) {
    case 'daily':   return `0 ${m} ${h} * * *`
    case 'weekly':  return `0 ${m} ${h} * * ${calendar.weekday}`
    case 'monthly': return `0 ${m} ${h} ${calendar.monthDay} * *`
  }
}

// ── 单次触发 ───────────────────────────────────────────
const single = reactive({ datetime: '' })

function buildSingle(): string {
  if (!single.datetime) return ''
  const d = new Date(single.datetime)
  if (isNaN(d.getTime())) return ''
  return `${d.getSeconds()} ${d.getMinutes()} ${d.getHours()} ${d.getDate()} ${d.getMonth() + 1} *`
}

// ── 自定义 ─────────────────────────────────────────────
const custom = reactive({ expr: '' })

// ── 生成预览 ───────────────────────────────────────────
function generate(): string {
  switch (activeTab.value) {
    case 'periodic':  return buildPeriodic()
    case 'calendar':  return buildCalendar()
    case 'single':    return buildSingle()
    case 'custom':    return custom.expr.trim()
    default:          return ''
  }
}

function handlePreview() {
  const expr = generate()
  if (!expr) { Message.warning('请先完善配置'); return }
  preview.value = expr
}

function handleConfirm() {
  const expr = generate()
  if (!expr) { Message.warning('请先完善配置'); return }
  preview.value = expr
  emit('update:modelValue', expr)
  visible.value = false
}

// ── 对外暴露 open ──────────────────────────────────────
function open(initExpr?: string) {
  if (initExpr) {
    preview.value = initExpr
    custom.expr = initExpr
    activeTab.value = 'custom'
  }
  visible.value = true
}

defineExpose({ open })
</script>

<template>
  <!-- 触发区域 -->
  <a-input-group compact style="display:flex">
    <a-input
      :model-value="modelValue"
      placeholder="点击右侧按钮选择 Cron 表达式"
      readonly
      style="flex:1;cursor:pointer"
      @click="open(modelValue)"
    />
    <a-button type="primary" @click="open(modelValue)">选择</a-button>
  </a-input-group>

  <!-- 弹窗 -->
  <a-modal
    v-model:visible="visible"
    title="Cron 表达式生成器"
    :width="560"
    :footer="false"
    unmount-on-close
  >
    <a-tabs v-model:active-key="activeTab" type="line" size="small">

      <!-- 周期触发 -->
      <a-tab-pane key="periodic" title="周期触发">
        <div class="cron-section">
          <a-space wrap>
            <span class="cron-label">每</span>
            <a-input-number
              v-model="periodic.interval"
              :min="1" :max="999"
              style="width:90px"
            />
            <a-select v-model="periodic.unit" style="width:90px">
              <a-option value="seconds">秒</a-option>
              <a-option value="minutes">分钟</a-option>
              <a-option value="hours">小时</a-option>
              <a-option value="days">天</a-option>
            </a-select>
            <template v-if="periodic.unit === 'hours'">
              <span class="cron-label">，在第</span>
              <a-input-number v-model="periodic.atMinute" :min="0" :max="59" style="width:80px" />
              <span class="cron-label">分</span>
            </template>
            <template v-if="periodic.unit === 'days'">
              <span class="cron-label">，在</span>
              <a-time-picker
                v-model="periodic.atTime"
                format="HH:mm"
                style="width:110px"
                :allow-clear="false"
              />
            </template>
            <span class="cron-label">执行一次</span>
          </a-space>
        </div>
      </a-tab-pane>

      <!-- 日历触发 -->
      <a-tab-pane key="calendar" title="日历触发">
        <div class="cron-section">
          <a-space wrap>
            <a-select v-model="calendar.unit" style="width:110px">
              <a-option value="daily">每天</a-option>
              <a-option value="weekly">每周</a-option>
              <a-option value="monthly">每月</a-option>
            </a-select>
            <template v-if="calendar.unit === 'weekly'">
              <a-select v-model="calendar.weekday" style="width:90px">
                <a-option v-for="w in weekdayOptions" :key="w.value" :value="w.value">{{ w.label }}</a-option>
              </a-select>
            </template>
            <template v-if="calendar.unit === 'monthly'">
              <span class="cron-label">第</span>
              <a-input-number v-model="calendar.monthDay" :min="1" :max="31" style="width:80px" />
              <span class="cron-label">日</span>
            </template>
            <a-time-picker
              v-model="calendar.time"
              format="HH:mm"
              style="width:110px"
              :allow-clear="false"
            />
            <span class="cron-label">执行一次</span>
          </a-space>
        </div>
      </a-tab-pane>

      <!-- 单次触发 -->
      <a-tab-pane key="single" title="单次触发">
        <div class="cron-section">
          <a-space>
            <span class="cron-label">触发时间</span>
            <a-date-picker
              v-model="single.datetime"
              show-time
              format="YYYY-MM-DD HH:mm:ss"
              style="width:220px"
            />
          </a-space>
        </div>
      </a-tab-pane>

      <!-- 自定义 -->
      <a-tab-pane key="custom" title="自定义">
        <div class="cron-section">
          <a-input
            v-model="custom.expr"
            placeholder="秒 分 时 日 月 周  （robfig/cron 6位格式）"
            allow-clear
          />
          <div class="cron-hint">
            格式：<code>秒(0-59) 分(0-59) 时(0-23) 日(1-31) 月(1-12) 周(0-6,0=周日)</code><br/>
            示例：<code>0 */5 * * * *</code>（每5分钟）&nbsp;
            <code>0 0 9 * * 1</code>（每周一9点）
          </div>
        </div>
      </a-tab-pane>
    </a-tabs>

    <!-- 预览 + 操作 -->
    <div class="cron-footer">
      <div class="cron-preview-row">
        <span class="cron-label">预览：</span>
        <a-tag color="arcoblue" size="medium" style="font-family:monospace;font-size:13px">
          {{ preview || '—' }}
        </a-tag>
        <a-button size="small" type="outline" @click="handlePreview">生成预览</a-button>
      </div>
      <a-space style="justify-content:flex-end;width:100%">
        <a-button @click="visible = false">取消</a-button>
        <a-button type="primary" @click="handleConfirm">确定</a-button>
      </a-space>
    </div>
  </a-modal>
</template>

<style scoped>
.cron-section {
  padding: 16px 4px 8px;
  min-height: 80px;
}
.cron-label {
  font-size: 13px;
  color: var(--color-text-2);
  white-space: nowrap;
}
.cron-hint {
  margin-top: 10px;
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 1.8;
}
.cron-hint code {
  background: var(--color-fill-2);
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 12px;
}
.cron-footer {
  border-top: 1px solid var(--color-border-2);
  padding-top: 14px;
  margin-top: 8px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.cron-preview-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
