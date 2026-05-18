import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export type NotifLevel = 'info' | 'success' | 'warning' | 'error';

export interface Notification {
  id: number;
  level: NotifLevel;
  title: string;
  body: string;
  time: string;
  read: boolean;
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([]);

  const unreadCount = computed(() => notifications.value.filter((n) => !n.read).length);

  function setNotifications(list: Notification[]) {
    notifications.value = list;
  }

  function markRead(id: number) {
    const n = notifications.value.find((n) => n.id === id);
    if (n) n.read = true;
  }

  function markAllRead() {
    notifications.value.forEach((n) => (n.read = true));
  }

  function remove(id: number) {
    notifications.value = notifications.value.filter((n) => n.id !== id);
  }

  return { notifications, unreadCount, setNotifications, markRead, markAllRead, remove };
});
