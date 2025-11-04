<template>
  <div class="profile-control">
    <BaseAlert v-if="error">{{ error }}</BaseAlert>
    <form @submit.prevent="changePassword" class="profile-control__form" >
      <div class="change-password-form__group">
        <BaseInput
          id="password"
          name="password"
          type="password"
          label="Введите текущий пароль"
          v-model="form.currentPassword"
          :error="errors.currentPassword"
          placeholder="••••••••"
          autocomplete="current-password"
        />

        <BaseInput
          id="password"
          name="password"
          type="password"
          label="Введите новый пароль"
          v-model="form.newPassword"
          :error="errors.newPassword"
          placeholder="••••••••"
          autocomplete="current-password"
        />


        <BaseInput
          id="passwordRepeat"
          name="passwordRepeat"
          type="password"
          label="Подтвердите новый пароль"
          v-model="form.confirmPassword"
          :error="errors.confirmPassword"
          placeholder="••••••••"
          autocomplete="current-password"
        />
      </div>

      <div class="change-password-form__actions">
        <BaseButton
          type="button"
          variant="secondary"
          :disabled="loading"
          loadingText="Смена пароля..."
          @click="emit('close')"
        >
          Отмена
        </BaseButton>

        <BaseButton
          type="submit"
          variant="primary"
          :loading="loading"
          loadingText="Смена пароля..."
        >
          Сменить
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import BaseInput from '~/components/ui/BaseInput.vue';
import BaseButton from '~/components/ui/BaseButton.vue';
import BaseAlert from '~/components/ui/BaseAlert.vue';
import { useUserStore } from '~/stores/user';
import { useNotificationStore } from '@/stores/notification';

const emit = defineEmits(['success', 'close']);

const userStore = useUserStore();
const loading = ref(false);
const error = ref('');

const form = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const errors = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const resetValidationErrors = () => {
  errors.currentPassword = '';
  errors.newPassword = '';
  errors.confirmPassword = '';
}

const validatePasswordForm = () => {
  let isValid = true;

  resetValidationErrors();

  if (!form.currentPassword) {
    errors.currentPassword = 'Введите текущий пароль';
    isValid = false;
  }

  if (!form.newPassword) {
    errors.newPassword = 'Введите новый пароль';
    isValid = false;
  } else if (form.newPassword.length < 8) {
    errors.newPassword = 'Пароль должен содержать не менее 8 символов';
    isValid = false;
  } else if (!/[A-Za-z]/.test(form.newPassword) || !/[0-9]/.test(form.newPassword)) {
    errors.newPassword = 'Пароль должен содержать буквы и цифры';
    isValid = false;
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = 'Подтвердите новый пароль';
    isValid = false;
  } else if (form.newPassword !== form.confirmPassword) {
    errors.confirmPassword = 'Пароли не совпадают';
    isValid = false;
  }

  return isValid;
};

const changePassword = async () => {
  if (loading.value) return;

  error.value = '';

  if (!validatePasswordForm()) {
    return;
  }

  loading.value = true;

  try {
    const success = await userStore.changePassword(form.currentPassword, form.newPassword);

    if (success) {
      emit('success');
      const notify = useNotificationStore();
      notify.notify('Пароль изменён', 'success');
    } else {
      error.value = 'Не удалось изменить пароль';
    }
  } catch (err) {
    console.error('Ошибка при изменении пароля:', err);
    error.value = 'Не удалось изменить пароль';
  } finally {
    loading.value = false;
  }
};
</script>

<style>
.profile-control {
  height: 100%;
}

.profile-control__form {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.change-password-form__group {
  display: flex;
  flex-direction: column;
  gap: var(--size-4);
}

.change-password-form__actions {
  display: flex;
  gap: 5px;
  margin-top: auto;
}
</style>
