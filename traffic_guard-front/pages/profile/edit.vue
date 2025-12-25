<template>
  <div class="profile">
    <header class="profile__header">
      <h1 class="profile__title page-title">Профиль</h1>
    </header>

    <section class="profile__content">
      <div class="profile-card profile-card--avatar">
        <div class="profile-card__data profile-card__data--row">
          <img src="/img/profile-icon.png" alt="profile-icon">
          <div class="prorile-card__data-block">
            <p class="profile-field__value_blacked-bolded">{{ fullName }}</p>
            <p class="profile-field__value_blacked">{{ roleDisplayName }}</p>
            <p
              :selectedTextStyle="{
                backgroundColor: generateColor(deviceName).background,
                color: generateColor(deviceName).color
              }"
              >{{ deviceName }}</p>
          </div>
        </div>

        <button class="profile-card__button" :class="{'profile-card__button_disabled': true}" :disabled="true">
          Загрузить аватар
        </button>
      </div>

      <div class="profile-card">
        <div class="profile-card__data">
          <div class="profile-card__header">
            <p class="profile-card__title">Персональные данные</p>

            <button class="profile-card__button" @click="showChangeUserDataModal">
              Изменить данные
            </button>
          </div>

          <div class="profile-card__fields">
            <div class="profile-field">
              <p class="profile-field__label">Имя</p>
              <p class="profile-field__value">{{ userData?.first_name || '–' }}</p>
            </div>

            <div class="profile-field">
              <p class="profile-field__label">Фамилия</p>
              <p class="profile-field__value">{{ userData?.last_name || '–' }}</p>
            </div>

            <div class="profile-field">
              <p class="profile-field__label">Отчество</p>
              <p class="profile-field__value">{{ userData?.patronymic || '–' }}</p>
            </div>

            <div class="profile-field">
              <p class="profile-field__label">Email</p>
              <p class="profile-field__value">{{ userData?.email || '–' }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="profile-card">
        <div class="profile-card__data">
          <div class="profile-card__header">
            <p class="profile-card__title">Безопасность</p>

            <button class="profile-card__button" @click="showChangePasswordModal">
              Сбросить пароль
            </button>
          </div>

          <div class="profile-card__fields">
            <div class="profile-field">
              <p class="profile-field__label">Пароль</p>
              <p class="profile-field__value">***********</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <SideModal
      v-model="isShowChangeUserDataModal"
      title="Редактировать профиль"
      @close="closeChangeUserDataModal"
    >
      <ManageUserForm
        :manageType="'update'"
        :userData="userData"
        :show-extra-buttons="false"
        :is-profile-form="true"
        @close="closeChangeUserDataModal"
        @success="onSuccessChangeUserData"
      />
    </SideModal>

    <BaseModal
      v-if="isShowChangePasswordModal"
      @close="closeChangePasswordModal"
    >
      <template #main>
        <ChangePasswordForm
          @close="closeChangePasswordModal"
          @success="onSuccessChangePassword"
        />
      </template>

    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { definePageMeta } from '#imports'
import { ref, onMounted } from 'vue';
import { getOnlyDeviceName, getRoleDisplayName } from '~/helpers';
import { useUserStore } from '~/stores/user';
import SideModal from '~/components/ui/SideModal.vue';
import ManageUserForm from '~/components/users/ManageUserForm.vue';
import BaseModal from '~/components/ui/BaseModal.vue';
import ChangePasswordForm from '~/components/profile/ChangePasswordForm.vue';
import { useDeviceColors } from '~/composables/useDeviceColors';


definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const { generateColor } = useDeviceColors();

const userStore = useUserStore();
const userData = computed(() => userStore.user);

const fullName = computed(() => userData.value?.last_name + ' ' + userData.value?.first_name);
const roleDisplayName = computed(() => getRoleDisplayName(userData.value?.role || '–-'));
const deviceName = computed(() => getOnlyDeviceName(userData.value?.role || 'SA'))

const isShowChangeUserDataModal = ref(false);
const isShowChangePasswordModal = ref(false);

const showChangePasswordModal = () => {
  isShowChangePasswordModal.value = true;
};

const showChangeUserDataModal = () => {
  isShowChangeUserDataModal.value = true;
};

const closeChangePasswordModal = () => {
  isShowChangePasswordModal.value = false;
};

const closeChangeUserDataModal = () => {
  isShowChangeUserDataModal.value = false;
};

const onSuccessChangePassword = () => {
  closeChangePasswordModal();
}

const onSuccessChangeUserData = () => {
  closeChangeUserDataModal();
}

const fetchUserData = async () => {
  if (!userStore.user) {
    await userStore.fetchUserInfo();
  }
};

fetchUserData();
</script>

<style scoped lang="scss">
.profile {
  &__content {
    margin-top: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
}

.profile-card {
  width: 100%;
  padding: 1.5rem;
  background: #fafafa;
  border: 1px solid #e4e4e7;
  border-radius: 28px;
  display: flex;
  justify-content: space-between;

  &--avatar {
    align-items: center;
  }

  &__data {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;

    &--row {
      flex-direction: row;
      width: 80%;
    }

    &-block {
      display: flex;
      flex-direction: column;
      gap: 0.375rem;
    }
  }

  &__header {
    display: flex;
    justify-content: space-between;
    width: 100%;
  }

  &__title {
    font-size: 1rem;
    font-weight: 500;
    color: #18181b;
    align-self: center;
  }

  &__fields {
    display: flex;
    flex-wrap: wrap;
    column-gap: 1.5rem;
    row-gap: 0.75rem;
  }

  &__button {
    padding: 0.5rem 0.75rem;
    background-color: #2563eb;
    color: #fff;
    font-size: 0.875rem;
    line-height: 1.25rem;
    border: 1px solid transparent;
    border-radius: 999px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    cursor: pointer;
    font-weight: 600;

    &_disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
}

.profile-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  width: 40%;

  &__label {
    font-size: 1rem;
    font-weight: 500;
    line-height: 1.5rem;
    color: #3f3f46;
  }

  &__value {
    font-size: 1rem;
    font-weight: 400;
    line-height: 1.25rem;
    color: #71717a;

    &_blacked {
      color: #3F3F46;

      &-bolded {
        font-weight: 600;
      }
    }

    
  }
}

</style>