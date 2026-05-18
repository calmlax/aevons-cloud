<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  IconCheckCircle,
  IconExclamationCircle,
  IconInfoCircle,
  IconCloseCircle,
  IconCheck,
  IconDelete,
} from '@arco-design/web-vue/es/icon'
import { storeToRefs } from 'pinia'
import { useNotificationStore } from '@/store/modules/notification'

const { t } = useI18n()
const notifStore = useNotificationStore()
const { notifications, unreadCount } = storeToRefs(notifStore)
const { markRead, markAllRead, remove } = notifStore

type NotifLevel = 'info' | 'success' | 'warning' | 'error'

const filter = ref<'all' | 'unread'>('all')
const filtered = computed(() =>
  filter.value === 'unread' ? notifications.value.filter(n => !n.read) : notifications.value
)

const levelIcon: Record<NotifLevel, unknown> = {
  info: IconInfoCircle,
  success: IconCheckCircle,
  warning: IconExclamationCircle,
  error: IconCloseCircle,
}
const levelColor: Record<NotifLevel, string> = {
  info: 'rgb(var(--arcoblue-6))',
  success: 'rgb(var(--green-6))',
  warning: 'rgb(var(--orange-6))',
  error: 'rgb(var(--red-6))',
}
</script>

<template>
  <div class="page-stack">
    <section class="notif-toolbar">
      <div class="notif-toolbar-left">
        <a-radio-group v-model="filter" type="button" size="small">
          <a-radio value="all">{{ t('notifications.filterAll') }}</a-radio>
          <a-radio value="unread">
            {{ t('notifications.filterUnread') }}
            <a-badge v-if="unreadCount" :count="unreadCount" style="margin-left:4px" />
          </a-radio>
        </a-radio-group>
      </div>
      <a-button
        v-if="unreadCount > 0"
        type="text"
        size="small"
        @click="markAllRead"
      >
        <template #icon><IconCheck /></template>
        {{ t('notifications.markAllRead') }}
      </a-button>
    </section>

    <div v-if="filtered.length === 0" class="notif-empty">
      <a-empty :description="t('notifications.empty')" />
    </div>

    <div v-else class="notif-list">
      <div
        v-for="n in filtered"
        :key="n.id"
        :class="['notif-item', { 'notif-item--unread': !n.read }]"
        @click="markRead(n.id)"
      >
        <component
          :is="levelIcon[n.level]"
          class="notif-icon"
          :style="{ color: levelColor[n.level] }"
        />
        <div class="notif-body">
          <div class="notif-title">{{ n.title }}</div>
          <div class="notif-text">{{ n.body }}</div>
          <div class="notif-time">{{ n.time }}</div>
        </div>
        <button class="notif-delete" :aria-label="t('notifications.delete')" @click.stop="remove(n.id)">
          <IconDelete />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notif-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.notif-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 280px;
}

.notif-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.notif-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 12px;
  background: var(--color-bg-2, #fff);
  border: 1px solid var(--color-border, #e5e6eb);
  cursor: pointer;
  transition: background 0.15s;
  position: relative;
}

.notif-item:hover {
  background: var(--color-fill-2, #f7f8fa);
}

.notif-item--unread {
  border-left: 3px solid rgb(var(--arcoblue-6));
}

.notif-item--unread::before {
  content: '';
  position: absolute;
  top: 14px;
  right: 40px;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: rgb(var(--arcoblue-6));
}

.notif-icon {
  font-size: 20px;
  flex-shrink: 0;
  margin-top: 1px;
}

.notif-body {
  flex: 1;
  min-width: 0;
}

.notif-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-1, #1d2129);
  margin-bottom: 2px;
}

.notif-text {
  font-size: 13px;
  color: var(--color-text-3, #86909c);
  line-height: 1.5;
}

.notif-time {
  font-size: 12px;
  color: var(--color-text-4, #c9cdd4);
  margin-top: 4px;
}

.notif-delete {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: 0;
  background: transparent;
  color: var(--color-text-4, #c9cdd4);
  border-radius: 6px;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.15s, background 0.15s, color 0.15s;
}

.notif-item:hover .notif-delete {
  opacity: 1;
}

.notif-delete:hover {
  background: var(--color-danger-light-1);
  color: rgb(var(--red-6));
}
</style>
