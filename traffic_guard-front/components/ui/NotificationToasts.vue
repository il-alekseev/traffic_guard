<template>
  <div class="notifications-wrapper">
    <transition-group name="toast" tag="div">
      <div
        v-for="n in notificationStore.notifications"
        :key="n.id"
        class="toast"
        :class="n.type"
      >
        <div class="toast__body">
          <div v-if="n.type === 'success'" class="toast__icon-container toast__icon-container_success">
            <CheckIcon class="toast__icon" />
          </div>
          <p class="toast__text">{{ n.message }}</p>
        </div>
        
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { useNotificationStore } from '~/stores/notification'
import CheckIcon from "~/assets/img/check.svg";


const notificationStore = useNotificationStore();

</script>

<style scoped>
.notifications-wrapper {
  position: fixed;
  bottom: 4rem;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  align-items: center;
}

.toast {
  padding: 1rem 2rem;
  border-radius: 0.5rem;
  box-shadow: 0px 4px 4px rgba(0, 32, 51, 0.04), 0px 8px 24px rgba(0, 32, 51, 0.12);
  animation: fade-in-out 2s ease forwards;
}

.toast__body {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.toast__icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
}

.toast__icon-container_success {
  color: #09B37B;
}

.toast__icon {
  width: 1.5rem;
  height: 1.5rem;
}

.toast__text {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #FFFFFF;
}

.toast.success {
  background-color: #002033;
}
.toast.error {
  background-color: #e53e3e;
}

.toast.info {
  background-color: #3182ce;
}

@keyframes fade-in-out {
  0% {
    opacity: 0;
    transform: translateY(100%);
  }
  10% {
    opacity: 1;
    transform: translateY(0);
  }
  90% {
    opacity: 1;
    transform: translateY(0);
  }
  100% {
    opacity: 0;
    transform: translateY(100%);
  }
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(100%);
}
</style>
