  <template>
  <div class="dashboard">
    <div class="dashboard__container">
      <SidebarComponent />
      <main class="dashboard__main">
        <HeaderComponent @logout="logout"/>
        <!-- <div class="dashboard__content-blur-layer"></div> -->
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
import { useDevicesControlStore } from '~/stores/devicesControl';
import { useCategoriesControlStore } from '~/stores/categoriesControl';
import {navigateTo} from "#app";
import NotificationToasts from "~/components/ui/NotificationToasts.vue";
import SidebarComponent from "~/components/dashboard/SidebarComponent.vue";
import HeaderComponent from "~/components/dashboard/HeaderComponent.vue";

const userStore = useUserStore();
const deviceControlStore = useDevicesControlStore();
const categoriesControlStroe = useCategoriesControlStore();

const checkUser = async () => {
  if (!userStore.user) {
    await userStore.fetchUserInfo();
  }
  if (userStore.user && userStore.user?.is_need_to_change_password) {
    return navigateTo('/change-password');
  }
}

const checkDevices = async () => {
  if (!deviceControlStore.devices) {
    await deviceControlStore.fetchDevices();
  }
}

const checkCategories = async () => {
  if (!categoriesControlStroe.categories) {
    await categoriesControlStroe.fetchCategories();
  }
}

onMounted(async () => {
  await checkUser();
  checkDevices();
  checkCategories();
})


const logout = () => {
  userStore.logout();
  navigateTo('/auth');
}
</script>

<style>
.dashboard {
  min-height: 100vh;
  
}

.dashboard__container {
  display: flex;
  min-height: 100%;
  position: relative;
}

.dashboard__main {
  width: 80%;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.dashboard__content {
  flex: 1;
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
