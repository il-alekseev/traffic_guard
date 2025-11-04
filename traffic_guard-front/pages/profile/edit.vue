<template>
  <div class="profile-control">
    <div class="profile-control__header">
      <h1 class="profile-control__title page-title">Управление профилем</h1>
    </div>

    <div class="profile-card">
      <div class="profile-card__title"><p>Личные данные</p></div>
      <div v-if="userData" class="profile-card__fake-form">
        <BaseInput
          id="last_name"
          name="last_name"
          label="Фамилия"
          v-model="userData.last_name"
          placeholder="Фамилия"
          :disabled="true"
        />

        <BaseInput
          id="first_name"
          name="first_name"
          label="Имя"
          v-model="userData.first_name"
          placeholder="Имя"
          :disabled="true"
        />

        <BaseInput
          id="patronymic"
          name="patronymic"
          label="Отчество"
          v-model="userData.patronymic"
          placeholder="Отчество"
          :disabled="true"
        />

        <BaseInput
          id="login"
          name="login"
          label="Логин"
          v-model="userData.login"
          placeholder="Username"
          extraInfo="Используется для входа в систему"
          :disabled="true"
        />

        <BaseInput
          id="email"
          name="email"
          label="Email"
          v-model="userData.email"
          placeholder="email"
          extraInfo="Используется для связи и восстановления доступа"
          :disabled="true"
        />
      </div>

      <div class="profile-card__buttons">
        <!-- <SubButton @click="showChangePasswordModal">
          <template #icon>
            <ChangePasswordIcon />
          </template>
          Сменить пароль
        </SubButton> -->
        <!-- <SubButton @click="showChangeUserDataModal">
          <template #icon>
            <EditDataIcon />
          </template>
          Изменить данные
        </SubButton> -->
      </div>
    </div>
    <SideModal
      v-model="isShowChangePasswordModal"
      title="Смена пароля"
      @close="isShowChangePasswordModal = false"
    >
      <ChangePasswordForm
        @close="isShowChangePasswordModal = false"
        @success="onSuccessChangePassword"
      />
    </SideModal>
    <!-- <SideModal
      v-model="isShowChangeUserDataModal"
      title="Изменить данные"
      @close="isShowChangeUserDataModal = false"
    >
      <ChangeUserDataForm
        :userData="userData"
        @close="isShowChangeUserDataModal = false"
        @success="onSuccessChangeUserData"
      />
    </SideModal> -->
  </div>
</template>

<script setup lang="ts">
import { definePageMeta } from '#imports'
import { ref, onMounted } from 'vue';
import { useUserStore } from '~/stores/user';
import type { User } from '~/types/user';
import BaseInput from '~/components/ui/BaseInput.vue';
// import SubButton from '~/components/UI/SubButton.vue';
import SideModal from '~/components/ui/SideModal.vue';
import ChangePasswordForm from '~/components/profile/ChangePasswordForm.vue';
// import ChangePasswordIcon from '~/assets/img/change-password.svg';
// import EditDataIcon from '~/assets/img/note-pencil.svg';
// import ChangeUserDataForm from '~/components/profile/ChangeUserDataForm.vue';

definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const userStore = useUserStore();
const userData = ref<User>();

const isShowChangePasswordModal = ref(false);
const isShowChangeUserDataModal = ref(false);

const showChangePasswordModal = () => {
  isShowChangePasswordModal.value = true;
};

// const showChangeUserDataModal = () => {
//   isShowChangeUserDataModal.value = true;
// };

const fetchUserData = async () => {
  if (!userStore.user) {
    await userStore.fetchUserInfo();
  }

  if (userStore.user) {
    userData.value = { ...userStore.user };
  }
};

const setNotification = (title: string) => {
  
}

const onSuccessChangePassword = () => {
  setNotification('Пароль сменён');
  isShowChangePasswordModal.value = false;
}

// const onSuccessChangeUserData = () => {
//   setNotification('Данные изменены');
//   isShowChangeUserDataModal.value = false;
// }

onMounted(async () => {
  await fetchUserData();
});
</script>

<style>
.profile-control__header {
  margin-bottom: var(--size-6);
}

.profile-control__title {
  font-size: var(--font-size-3xl);
  line-height: 2rem;
  font-weight: 600;
  color: var(--color-typo-primary);
  margin-block: var(--size-1);
}

.profile-card {
  padding: var(--size-6);
  background-color: var(--color-white-tertiary);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  max-width: 27.5rem;
}

.profile-card__title {
  font-size: var(--font-size-lg);
  line-height: 1.5rem;
  font-weight: 600;
  color: var(--color-typo-primary);
  margin-bottom: var(--size-4);
}

.profile-card__fake-form {
  display: flex;
  flex-direction: column;
  gap: var(--size-4);
  margin-bottom: var(--size-8);
}

.profile-card__buttons {
  display: flex;
  gap: var(--size-2);
}
</style>