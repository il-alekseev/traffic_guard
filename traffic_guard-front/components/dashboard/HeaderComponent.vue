<template>
  <div class="header">
    <div class="header__content">
      <div class="header__user">
        <div class="header__notification">
          <div class="header__notification-icon">
            <NotificationIcon />
          </div>
        </div>
        <div class="header__user-data">
          <div class="header__user-avatar">
            <UserCircleIcon />
          </div>
          <div class="header__user-info">
            <p class="header__user-name">{{ fullName }}</p>
          </div>
          <div class="header__user-arrow"  @click="toggleIsMenuOpen">
            <ArrowIcon />
          </div>
        </div>

        <div
          v-if="isMenuOpen"
          class="header__user-menu user-menu">
          <ul>
            <li>
              <NuxtLink to="/profile/edit" @click="toggleIsMenuOpen">
                <div class="user-menu__item">
                  <div class="user-menu__icon">
                    <UserIcon />
                  </div>
                  <span class="user-menu__link">Профиль</span>
                </div>
              </NuxtLink>
            </li>
            <li @click="emit('logout')">
              <div class="user-menu__item user-menu__item_red">
                <div class="user-menu__icon">
                  <LogoutIcon />
                </div>
                <span class="user-menu__link">Выход</span>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import { useUserStore } from '~/stores/user';
import UserCircleIcon from "~/assets/img/user-circle.svg";
import ArrowIcon from "~/assets/img/arrow-down.svg";
import UserIcon from "~/assets/img/user.svg";
import LogoutIcon from "~/assets/img/logout.svg"
import NotificationIcon from "~/assets/img/notification-active.svg"


const emit = defineEmits(['logout']);

const userStore = useUserStore();

const fullName = computed(() => {
  if (userStore.firstName && userStore.lastName) return userStore.fullName;
  else return "Admin Admin"
})

const isMenuOpen = ref(false);

const toggleIsMenuOpen = () => {
  isMenuOpen.value = !isMenuOpen.value;
}

</script>

<style>
.header {
  height: 4rem;
  position: sticky;
  background-color: #FFFFFF;
  z-index: 999999;
}

.header__content {
  width: 100%;
  padding-inline: 2rem;
  padding-block: 1rem;
  display: flex;
  justify-content: flex-end;  
}

.header__user {
  display: flex;
  gap: 1.5rem;
}

.header__notification {
  padding-right: 1.5rem;
  border-right: 1px solid #E4E4E7;
  display: flex;
  justify-content: center;
  align-items: center;
}

.header__notification-icon {
  display: flex;
  justify-content: center;
  align-items: center;
}

.header__user-data {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  position: relative;
}

.header__user-avatar {
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header__user-avatar img{
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.header__user-name {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
}

.header__user-arrow {
  width: 1.25rem;
  height: 1.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.header__user-arrow img{
  width: 100%;
  height: 100%;
  object-fit: contain;
}


.header__user-menu {
  z-index: 999999;
  position: absolute;
  top: 3.5rem;
  right: 0;
  width: 15rem;
  background-color: #FFFFFF;
  border-radius: 0.375rem;
  box-shadow: 0px 1px 2px 0px #0000000D;
  border: 1px solid #E4E4E7;
}

.header__user-menu ul {
  list-style-type:  none;
  display: flex;
  flex-direction: column;
}

.header__user-menu ul li {
  padding-left: 1.25rem;
  padding-block: 0.75rem;
  border-top: 1px solid #E4E4E7;
}

.header__user-menu ul li:first-child {
  border-top: none;
}

.user-menu__item {
  display: flex;
  gap: 0.5rem;
}

.user-menu__item_red {
  color: #C10007;
}

.user-menu__icon {
  width: 1.5rem;
  height: 1.5rem;
  object-fit: cover;
  color: #9CA3AF;
}

.user-menu__item_red .user-menu__icon {
  color: #C10007;
}

.user-menu__link {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
  text-decoration: none;
}

.user-menu__item_red .user-menu__link {
  color: #C10007;
}

.header__user-menu li {
  cursor: pointer;
}

.header__user-menu li a{
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
  text-decoration: none;
}

</style>