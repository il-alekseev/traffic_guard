import { defineStore } from 'pinia'

interface Notification {
  id: number
  message: string
  type: 'success' | 'error' | 'info'
}

const TIME_TO_SHOW = 2000; // 2 second. WARNING: If you need to change it, then change the animation value in the css class .toast

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [] as Notification[],
  }),
  actions: {
    notify(message: string, type: 'success' | 'error' | 'info' = 'success') {
      const id = Date.now()
      this.notifications.push({ id, message, type })
      setTimeout(() => {
        this.notifications = this.notifications.filter(n => n.id !== id)
      }, TIME_TO_SHOW)
    },
  },
})
