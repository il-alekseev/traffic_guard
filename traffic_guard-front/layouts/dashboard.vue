<template>
  <div class="dashboard">
    <div class="dashboard__container">
      <SidebarComponent />
      <main class="dashboard__main">
        <ClientOnly>
          <HeaderComponent @logout="logout"/>
        </ClientOnly>
        <div class="dashboard__content">
          <slot />
        </div>

      </main>
    </div>
    <NotificationToasts />
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '~/stores/user';
import { useDevicesStore } from '~/stores/devices';
import { useCategoriesStore } from '~/stores/categories';
import {navigateTo} from "#app";
import NotificationToasts from "~/components/ui/NotificationToasts.vue";
import SidebarComponent from "~/components/dashboard/SidebarComponent.vue";
import HeaderComponent from "~/components/dashboard/HeaderComponent.vue";

const userStore = useUserStore();
const devicesStore = useDevicesStore();
const categoriesStore = useCategoriesStore();

const checkUser = async () => {
  if (!userStore.user) {
    await userStore.fetchUserInfo();
  }
  if (userStore.user && userStore.user?.is_need_to_change_password) {
    return navigateTo('/change-password');
  }
}

const checkDevices = async () => {
  if (!devicesStore.devices) {
    await devicesStore.fetchDevices();
  }
}

const checkCategories = async () => {
  if (!categoriesStore.categories) {
    await categoriesStore.fetchCategories();
  }
}

onMounted(async () => {
  await Promise.all([
    checkUser(),
    checkDevices()
  ])
  checkCategories();
})


const logout = () => {
  userStore.logout();
  navigateTo('/auth');
}
</script>

<style>
.dashboard {
  height: 100vh;
  overflow: hidden;
}

.dashboard__container {
  display: flex;
  height: 100%;
  position: relative;
}

.dashboard__main {
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.dashboard__content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  background-color: #F4F4F5;
  border-top-left-radius: 28px;
  padding: 2rem;
}

.dashboard__content-blur-layer {
  background: #D7D7D766;
  backdrop-filter: blur(10px);
}

.page-title {
  font-weight: 500;
  font-size: 1.875rem;
  line-height: 2.25rem;
  color: #3F3F46;
}

.page-subtitle-block {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: #A1A1AA;
}

.page-subtitle {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
}

.page-refresh-btn {
  cursor: pointer;
  width: 1.25rem;
  height: 1.25rem;
}
</style>
