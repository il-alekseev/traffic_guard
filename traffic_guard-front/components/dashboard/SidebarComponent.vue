<template>
  <aside class="sidebar">
    <div class="sidebar__logo">
      <NuxtLink to="/dashboard" class="sidebar__logo-link">
        <div class="sidebar__logo-link-icon">
          <FazendaLogo />
        </div>
        <span class="sidebar__logo-link-name">
          Fazenda
        </span>
      </NuxtLink>
    </div>
    <div class="sidebar__content">
      <nav class="sidebar__nav">
        <template v-for="item in menuItems" :key="item.label">
          <NuxtLink
            :to="item.to"
            class="sidebar__item"
            :class="{ 'sidebar__item--active': $route.path === item.to }"
          >
            <div v-if="item.iconComponent" class="sidebar__item-icon">
              <component :is="item.iconComponent" class="sidebar__item-icon-img" />
            </div>
            <p class="sidebar__item-link">{{ item.label }}</p>
          </NuxtLink>
          <!-- <a
            v-else
            :href="item.href"
            class="sidebar__item"
            target="_blank"
            rel="noopener noreferrer"
          >
            <div v-if="item.iconComponent" class="sidebar__item-icon">
              <component :is="item.iconComponent" class="sidebar__item-icon-img" />
            </div>
            <p class="sidebar__item-link">{{ item.label }}</p>
          </a> -->
        </template>
      </nav>
    </div>
    <div class="sidebar__footer">
      <NuxtLink
        to="/settings"
        class="sidebar__item"
        :class="{ 'sidebar__item--active': $route.path === '/settings' }"
      >
        <div class="sidebar__item-icon">
          <SettingsIcon />
        </div>
        <p class="sidebar__item-link">Настройки</p>
      </NuxtLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useUserStore } from '~/stores/user';
import { useSidebarMenu } from '~/composables/useSidebarMenu';
import FazendaLogo from "~/assets/img/logo.svg"
import SettingsIcon from "~/assets/img/settings.svg"

const emit = defineEmits(['logout']);

const userStore = useUserStore();
const { getMenuForRole } = useSidebarMenu();

const role = computed(() => userStore.role);
// const menuItems = computed(() => getMenuForRole(role.value));
const menuItems = computed(() => getMenuForRole('KSF'));
</script>

<style>
.sidebar {
  width: 15%;
  min-width: 18rem;
  max-width: 18rem;
  height: 100vh;
  position: sticky;
  top: 0;
  background-color: #FFFFFF;
  display: flex;
  flex-direction: column;
  padding-inline: 1.5rem;
  padding-block: 1rem;
  z-index: 5;
  overflow: hidden;
}

.sidebar__logo {
  margin-bottom: 2.25rem;
}

.sidebar__logo-link {
  width: fit-content;
  height: 100%;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
}

.sidebar__logo-link-icon {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sidebar__logo-link-icon svg{
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.sidebar__logo-link-name {
  font-weight: 400;
  font-size: 1.125rem;
  line-height: 1.75rem;
  color: #3F3F46;
}

.sidebar__content {
  flex: 1;
  overflow-y: auto;
}

.sidebar__nav {
  display: flex;
  flex-direction: column;
}

.sidebar__item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 0rem 0.75rem 1.25rem;
  transition: all 0.2s ease;
  border-radius: 999px;
  text-decoration: none;
  background: #FFFFFF;
}

.sidebar__item:not(.sidebar__item--active):hover {
  background: #E8E8E8;
}

.sidebar__item--active {
  background: #18181B;
}

.sidebar__item-icon {
  width: 1.5rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #A1A1AA;
}

.sidebar__item--active .sidebar__item-icon {
  color: #FFFFFF;
}

.sidebar__item-icon img{
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.sidebar__item-link {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
  text-decoration: none;
}

.sidebar__item--active .sidebar__item-link {
  color: #FFFFFF;
}

.sidebar__footer {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-top: auto;
}
</style>
