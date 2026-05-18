<script setup lang="ts">
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { GaugeChart } from 'echarts/charts'
import VChart, { THEME_KEY } from 'vue-echarts'
import { provide } from 'vue'
import { useI18n } from 'vue-i18n'
import serverApi from '@/api/monitor/server'

use([CanvasRenderer, GaugeChart])

const { t } = useI18n()

const isDark = ref(document.body.getAttribute('arco-theme') === 'dark')
const mo = new MutationObserver(() => { isDark.value = document.body.getAttribute('arco-theme') === 'dark' })
onMounted(() => mo.observe(document.body, { attributes: true, attributeFilter: ['arco-theme'] }))
onUnmounted(() => mo.disconnect())
provide(THEME_KEY, computed(() => isDark.value ? 'dark' : ''))

const loading = ref(false)
const info = ref<any>(null)
const autoRefresh = ref(false)
let timer: ReturnType<typeof setInterval> | null = null
watch(autoRefresh, v => {
  if (v) { timer = setInterval(load, 5000) }
  else { if (timer) { clearInterval(timer); timer = null } }
})
onUnmounted(() => { if (timer) clearInterval(timer) })

async function load() {
  loading.value = true
  try { info.value = await serverApi.getInfo() as any }
  finally { loading.value = false }
}

function fmtBytes(b: number): string {
  if (!b) return '0 B'
  const u = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0, v = b
  while (v >= 1024 && i < u.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(2)} ${u[i]}`
}

function pctColor(p: number) {
  return p >= 90 ? '#f53f3f' : p >= 70 ? '#ff7d00' : '#00b42a'
}

function gauge(pct: number, label: string) {
  const c = pctColor(pct)
  const bg = isDark.value ? '#2a2a3a' : '#f0f0f5'
  return {
    backgroundColor: 'transparent',
    series: [{
      type: 'gauge',
      startAngle: 210, endAngle: -30, min: 0, max: 100, radius: '90%',
      axisLine: { lineStyle: { width: 12, color: [[pct / 100, c], [1, bg]] } },
      pointer: { show: false }, axisTick: { show: false },
      splitLine: { show: false }, axisLabel: { show: false },
      detail: { valueAnimation: true, formatter: '{value}%', color: c, fontSize: 20, fontWeight: 700, offsetCenter: [0, '5%'] },
      title: { text: label, color: isDark.value ? '#8b949e' : '#86909c', fontSize: 11, offsetCenter: [0, '35%'] },
      data: [{ value: pct }]
    }]
  }
}

const cpuPct   = computed(() => parseFloat(info.value?.cpu?.usedPercent || '0'))
const memPct   = computed(() => parseFloat(info.value?.mem?.usedPercent || '0'))
const goPct    = computed(() => { const g = info.value?.go; return g?.heapSys ? Math.round((g.heapInuse / g.heapSys) * 100) : 0 })
const redisPct = computed(() => {
  const r = info.value?.redis; if (!r) return 0
  const hits = Number(r.keyspaceHits) || 0, misses = Number(r.keyspaceMisses) || 0
  const total = hits + misses
  return total === 0 ? 100 : Math.round((hits / total) * 100)
})

load()
</script>

<template>
  <div class="page-stack">
    <!-- 工具栏 -->
    <div style="display:flex;justify-content:flex-end;align-items:center;gap:12px;margin-bottom:16px">
      <a-switch v-model="autoRefresh" :checked-text="t('monitor.server.autoRefresh')" :unchecked-text="t('monitor.server.manual')" />
      <a-button type="primary" :loading="loading" @click="load">{{ t('monitor.server.refresh') }}</a-button>
    </div>

    <template v-if="loading && !info">
      <a-row :gutter="[16, 16]">
        <!-- 系统信息整行 -->
        <a-col :span="24">
          <a-card class="panel-card" :bordered="false">
            <template #title><a-skeleton-line :rows="1" style="width:120px" /></template>
            <a-row :gutter="[16, 8]">
              <a-col :span="24" :sm="12" :md="4" v-for="i in 6" :key="i">
                <div class="sys-item">
                  <a-skeleton-line :rows="1" style="width:60%;margin-bottom:6px" />
                  <a-skeleton-line :rows="1" style="width:90%" />
                </div>
              </a-col>
            </a-row>
          </a-card>
        </a-col>

        <!-- CPU / 内存 / Go 三列 -->
        <a-col :xs="24" :md="8" v-for="i in 3" :key="`top-${i}`" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title><a-skeleton-line :rows="1" style="width:80px" /></template>
            <a-skeleton :animation="true" style="height:160px;display:flex;align-items:center;justify-content:center">
              <a-skeleton-shape shape="circle" style="width:120px;height:120px;margin:0 auto" />
            </a-skeleton>
            <div v-for="j in 3" :key="j" class="kv">
              <a-skeleton-line :rows="1" style="width:40%" />
              <a-skeleton-line :rows="1" style="width:35%" />
            </div>
          </a-card>
        </a-col>

        <!-- Redis / DB / RocketMQ 三列 -->
        <a-col :xs="24" :md="8" v-for="i in 3" :key="`mid-${i}`" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title><a-skeleton-line :rows="1" style="width:80px" /></template>
            <a-skeleton :animation="true" style="height:160px;display:flex;align-items:center;justify-content:center">
              <a-skeleton-shape shape="circle" style="width:120px;height:120px;margin:0 auto" />
            </a-skeleton>
            <div v-for="j in 5" :key="j" class="kv">
              <a-skeleton-line :rows="1" style="width:40%" />
              <a-skeleton-line :rows="1" style="width:35%" />
            </div>
          </a-card>
        </a-col>

        <!-- 磁盘整行 -->
        <a-col :span="24">
          <a-card class="panel-card" :bordered="false">
            <template #title><a-skeleton-line :rows="1" style="width:80px" /></template>
            <a-row :gutter="[16, 16]">
              <a-col :xs="24" :sm="12" :lg="8" v-for="i in 3" :key="i">
                <div class="disk-card">
                  <a-row align="center" style="margin-bottom:8px">
                    <a-col :span="16"><a-skeleton-line :rows="1" style="width:70%" /></a-col>
                    <a-col :span="8" style="text-align:right"><a-skeleton-line :rows="1" style="width:40%;margin-left:auto" /></a-col>
                  </a-row>
                  <a-skeleton-line :rows="1" style="width:100%;margin-bottom:10px" />
                  <a-row :gutter="[4, 0]">
                    <a-col :span="8" v-for="j in 3" :key="j">
                      <a-skeleton-line :rows="2" />
                    </a-col>
                  </a-row>
                </div>
              </a-col>
            </a-row>
          </a-card>
        </a-col>
      </a-row>
    </template>

    <template v-else-if="info">
      <a-row :gutter="[16, 16]" align="stretch">

        <!-- 服务器系统信息（整行） -->
        <a-col :span="24">
          <a-card class="panel-card" :bordered="false">
            <template #title>🖥️ {{ t('monitor.server.sysInfo') }}</template>
            <a-row :gutter="[16, 8]">
              <a-col :span="24" :sm="12" :md="4" v-for="item in [
                { label: t('monitor.server.hostname'),  val: info.sys?.hostname },
                { label: t('monitor.server.os'),        val: `${info.sys?.platform} ${info.sys?.platformVersion}` },
                { label: t('monitor.server.kernel'),    val: info.sys?.kernelVersion },
                { label: t('monitor.server.arch'),      val: info.sys?.kernelArch },
                { label: t('monitor.server.bootTime'),  val: info.sys?.bootTime },
                { label: t('monitor.server.uptime'),    val: info.sys?.uptime },
              ]" :key="item.label">
                <div class="sys-item">
                  <div class="sys-label">{{ item.label }}</div>
                  <div class="sys-val" :title="item.val">{{ item.val || '-' }}</div>
                </div>
              </a-col>
            </a-row>
          </a-card>
        </a-col>

        <!-- CPU -->
        <a-col :xs="24" :md="8" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title>🖥️ {{ t('monitor.server.cpu') }}</template>
            <v-chart style="height:160px" :option="gauge(Math.round(cpuPct), t('monitor.server.cpuUsage'))" autoresize />
            <div class="kv"><span class="kl">{{ t('monitor.server.cpuModel') }}</span><span class="kv2" :title="info.cpu?.modelName">{{ info.cpu?.modelName || '-' }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.physicalCores') }} / {{ t('monitor.server.logicalCores') }}</span><span class="kv2">{{ info.cpu?.physicalCores }} / {{ info.cpu?.logicalCores }}</span></div>
          </a-card>
        </a-col>

        <!-- 内存 -->
        <a-col :xs="24" :md="8" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title>💾 {{ t('monitor.server.mem') }}</template>
            <v-chart style="height:160px" :option="gauge(Math.round(memPct), t('monitor.server.memUsage'))" autoresize />
            <div class="kv"><span class="kl">{{ t('monitor.server.memTotal') }}</span><span class="kv2">{{ fmtBytes(info.mem?.total) }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.memUsed') }}</span><span class="kv2" :style="{ color: pctColor(memPct) }">{{ fmtBytes(info.mem?.used) }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.memAvailable') }}</span><span class="kv2" style="color:#00b42a">{{ fmtBytes(info.mem?.available) }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.memCached') }}</span><span class="kv2">{{ fmtBytes(info.mem?.cached) }}</span></div>
          </a-card>
        </a-col>

        <!-- Go 运行时 -->
        <a-col :xs="24" :md="8" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title>
              🚀 {{ t('monitor.server.go') }}
              <a-tag color="arcoblue" size="small" style="margin-left:8px">{{ info.go?.version }}</a-tag>
            </template>
            <v-chart style="height:160px" :option="gauge(goPct, t('monitor.server.goHeapUsage'))" autoresize />
            <div class="kv"><span class="kl">{{ t('monitor.server.goOS') }}</span><span class="kv2">{{ info.go?.os }} / {{ info.go?.arch }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.goCPU') }}</span><span class="kv2">{{ info.go?.numCPU }}</span></div>
            <div class="kv"><span class="kl">Goroutine</span><span class="kv2">{{ info.go?.numGoroutine }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.goHeapAlloc') }}</span><span class="kv2">{{ fmtBytes(info.go?.heapAlloc) }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.goHeapInuse') }}</span><span class="kv2">{{ fmtBytes(info.go?.heapInuse) }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.goHeapSys') }}</span><span class="kv2">{{ fmtBytes(info.go?.heapSys) }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.goGcNum') }}</span><span class="kv2">{{ info.go?.gcNum }}</span></div>
            <div class="kv"><span class="kl">{{ t('monitor.server.goGcPause') }}</span><span class="kv2">{{ ((info.go?.gcPauseTotal || 0) / 1e6).toFixed(1) }} ms</span></div>
          </a-card>
        </a-col>

        <!-- Redis -->
        <a-col :xs="24" :md="8" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title>
              ⚡ {{ t('monitor.server.redis') }}
              <a-tag color="red" size="small" style="margin-left:8px">v{{ info.redis?.version }}</a-tag>
            </template>
            <a-alert v-if="info.redis?.error" type="error" :message="info.redis.error" />
            <template v-else>
              <v-chart style="height:160px" :option="gauge(redisPct, t('monitor.server.redisHitRate'))" autoresize />
              <div class="kv"><span class="kl">{{ t('monitor.server.redisOnline') }}</span><span class="kv2" style="color:#165dff;font-size:18px;font-weight:700">{{ info.redis?.onlineTokens }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.redisKeys') }}</span><span class="kv2">{{ info.redis?.dbSize }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.redisClients') }}</span><span class="kv2">{{ info.redis?.connectedClients }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.redisMemUsed') }}</span><span class="kv2">{{ info.redis?.usedMemory }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.redisMemPeak') }}</span><span class="kv2">{{ info.redis?.usedMemoryPeak }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.redisHits') }}</span><span class="kv2">{{ info.redis?.keyspaceHits }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.redisUptime') }}</span><span class="kv2">{{ info.redis?.uptimeSeconds }}s</span></div>
            </template>
          </a-card>
        </a-col>

        <!-- 数据库 -->
        <a-col :xs="24" :md="8" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title>
              🗄️ {{ t('monitor.server.db') }}
              <a-tag color="green" size="small" style="margin-left:8px">{{ info.db?.version }}</a-tag>
            </template>
            <a-alert v-if="info.db?.error" type="error" :message="info.db.error" />
            <template v-else>
              <div class="kv"><span class="kl">{{ t('monitor.server.dbConns') }}</span><span class="kv2">{{ info.db?.openConns }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.dbInUse') }}</span><span class="kv2" style="color:#f53f3f">{{ info.db?.inUseConns }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.dbIdle') }}</span><span class="kv2" style="color:#00b42a">{{ info.db?.idleConns }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.dbSize') }}</span><span class="kv2">{{ info.db?.dbSizeMB }} MB</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.dbProcess') }}</span><span class="kv2">{{ info.db?.processCount }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.dbUptime') }}</span><span class="kv2">{{ info.db?.uptime }}</span></div>
            </template>
          </a-card>
        </a-col>

        <!-- RocketMQ -->
        <a-col :xs="24" :md="8" style="display:flex;flex-direction:column">
          <a-card class="panel-card full-card" :bordered="false">
            <template #title>
              📨 {{ t('monitor.server.rocketmq') }}
              <template v-if="info.rocketmq?.enabled">
                <a-tag
                  :color="info.rocketmq.status === 'online' ? 'green' : info.rocketmq.status === 'degraded' ? 'orange' : 'red'"
                  size="small"
                  style="margin-left:8px"
                >{{ info.rocketmq.status }}</a-tag>
              </template>
              <a-tag v-else color="gray" size="small" style="margin-left:8px">{{ t('monitor.server.rmqDisabled') }}</a-tag>
            </template>
            <template v-if="!info.rocketmq?.enabled">
              <a-empty :description="t('monitor.server.rmqDisabled')" style="padding:24px 0" />
            </template>
            <template v-else>
              <!-- NameServer 连通性 -->
              <div class="kv" v-for="ns in info.rocketmq.nameServers" :key="ns.addr">
                <span class="kl">{{ t('monitor.server.rmqNameServer') }}</span>
                <span class="kv2" style="display:flex;align-items:center;gap:6px">
                  <a-badge :status="ns.reachable ? 'success' : 'danger'" />
                  {{ ns.addr }}
                </span>
              </div>
              <a-divider :margin="8" />
              <!-- 生产者 -->
              <div class="rmq-section">{{ t('monitor.server.rmqProducer') }}</div>
              <div class="kv"><span class="kl">{{ t('monitor.server.rmqGroup') }}</span><span class="kv2">{{ info.rocketmq.producer?.group }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.rmqRetry') }}</span><span class="kv2">{{ info.rocketmq.producer?.retryTimes }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.rmqTimeout') }}</span><span class="kv2">{{ info.rocketmq.producer?.sendTimeout }}s</span></div>
              <a-divider :margin="8" />
              <!-- 消费者 -->
              <div class="rmq-section">{{ t('monitor.server.rmqConsumer') }}</div>
              <div class="kv"><span class="kl">{{ t('monitor.server.rmqGroup') }}</span><span class="kv2">{{ info.rocketmq.consumer?.group }}</span></div>
              <div class="kv"><span class="kl">{{ t('monitor.server.rmqConcurrency') }}</span><span class="kv2">{{ info.rocketmq.consumer?.concurrency }}</span></div>
            </template>
          </a-card>
        </a-col>

        <!-- 磁盘（最后整行，多分区） -->
        <a-col :span="24">
          <a-card class="panel-card" :bordered="false">
            <template #title>💿 {{ t('monitor.server.disk') }}</template>
            <a-alert v-if="info.disk?.error" type="error" :message="info.disk.error" />
            <a-row v-else :gutter="[16, 16]">
              <a-col v-for="d in info.disk?.list" :key="d.path" :xs="24" :sm="12" :lg="8">
                <div class="disk-card">
                  <!-- 路径 + 类型 + 百分比 -->
                  <a-row align="center" style="margin-bottom:8px">
                    <a-col :span="16" style="display:flex;align-items:center;gap:6px">
                      <span class="disk-path">{{ d.path }}</span>
                      <a-tag size="small">{{ d.fstype }}</a-tag>
                    </a-col>
                    <a-col :span="8" style="text-align:right">
                      <span class="disk-pct" :style="{ color: pctColor(parseFloat(d.usedPercent)) }">{{ d.usedPercent }}%</span>
                    </a-col>
                  </a-row>
                  <!-- 进度条 -->
                  <a-progress :percent="parseFloat(d.usedPercent) / 100" :color="pctColor(parseFloat(d.usedPercent))" :show-text="false" size="small" animation style="margin-bottom:10px" />
                  <!-- 三列数据 -->
                  <a-row :gutter="[4, 0]" style="font-size:12px">
                    <a-col :span="8">
                      <div class="disk-stat-label">{{ t('monitor.server.diskTotal') }}</div>
                      <div class="disk-stat-val">{{ fmtBytes(d.total) }}</div>
                    </a-col>
                    <a-col :span="8" style="text-align:center">
                      <div class="disk-stat-label">{{ t('monitor.server.diskUsed') }}</div>
                      <div class="disk-stat-val" :style="{ color: pctColor(parseFloat(d.usedPercent)) }">{{ fmtBytes(d.used) }}</div>
                    </a-col>
                    <a-col :span="8" style="text-align:right">
                      <div class="disk-stat-label">{{ t('monitor.server.diskFree') }}</div>
                      <div class="disk-stat-val" style="color:#00b42a">{{ fmtBytes(d.free) }}</div>
                    </a-col>
                  </a-row>
                </div>
              </a-col>
            </a-row>
          </a-card>
        </a-col>

      </a-row>
    </template>

    <a-empty v-else :description="t('common.noData')" />
  </div>
</template>

<style scoped>
.full-card { flex: 1; }
.sys-item { padding: 6px 0; }
.sys-label { font-size: 11px; color: var(--color-text-3); margin-bottom: 2px; }
.sys-val { font-size: 13px; font-weight: 500; color: var(--color-text-1); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.kv { display: flex; justify-content: space-between; align-items: center; padding: 5px 0; border-bottom: 1px solid var(--color-border-1); font-size: 13px; gap: 8px; }
.kv:last-child { border-bottom: none; }
.kl  { color: var(--color-text-3); flex-shrink: 0; }
.kv2 { font-weight: 500; color: var(--color-text-1); text-align: right; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.disk-card { padding: 12px; border-radius: 8px; border: 1px solid var(--color-border-1); }
.disk-path { font-weight: 600; font-size: 14px; color: var(--color-text-1); }
.disk-pct  { font-weight: 700; font-size: 14px; }
.disk-stat-label { color: var(--color-text-3); margin-bottom: 2px; }
.disk-stat-val { font-weight: 600; color: var(--color-text-1); white-space: nowrap; }
.rmq-section { font-size: 12px; font-weight: 600; color: var(--color-text-3); text-transform: uppercase; letter-spacing: 0.5px; padding: 4px 0 2px; }
</style>
