<template>
  <div class="profile-control">
    <BaseAlert v-if="error">{{ error }}</BaseAlert>
    <form @submit.prevent="changePassword" class="profile-control__form" >
      <div class="change-password-form__group">
        <BaseInput
          id="password"
          name="password"
          type="password"
          label="Введите старый пароль"
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
          label="Подтверждение пароля"
          v-model="form.confirmPassword"
          :error="errors.confirmPassword"
          placeholder="••••••••"
          autocomplete="current-password"
        />
      </div>

      <div class="password-rules">
        <div class="password-rules__header">
          <div class="password-rules__icon" :class="allPasswordRulesPassed ? 'password-rules__icon_blue' : ''">
            <CheckMarkCircleIcon />
          </div>
          <div class="password-rules__title">
            Пароль должен содержать:
          </div>
        </div>
        <div class="password-rules__list">
          <div
            class="password-rules__rule"
            v-for="(value, ruleKey) in passwordRules"
            :key="ruleKey"
            :class="{ 'password-rules__rule_valid': value === true, 'password-rules__rule_invalid': value === false }"
          >
            {{ getNameOfPasswordRule(ruleKey as PasswordRuleKey) }}
          </div>
        </div>
      </div>

      <div class="change-password-form__actions">
        <BaseButton
          type="button"
          variant="secondary"
          :disabled="loading"
          loadingText="Смена пароля..."
          @click="emit('close')"
        >
          Отменить
        </BaseButton>

        <BaseButton
          type="submit"
          variant="primary"
          :loading="loading"
          loadingText="Смена пароля..."
        >
          Сохранить
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
import { useNotificationStore } from '~/stores/notification';
import CheckMarkCircleIcon from '~/assets/img/checkmark-circle.svg'
import type { PasswordRuleKey, PasswordRules } from '~/types/user';

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

const validatePasswordDetail = (password: string) => {
  const pwd = password.trim();

  passwordRules.minLength = pwd.length >= 8;
  passwordRules.lettersUpperAndLowerCase =
    /[a-z]/.test(pwd) && /[A-Z]/.test(pwd);
  passwordRules.numbers = /\d/.test(pwd);
  passwordRules.specialSymbols =
    /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd);
};


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

const passwordRuleLabels: Record<PasswordRuleKey, string> = {
  minLength: 'Не менее 8 символов',
  lettersUpperAndLowerCase: 'Буквы верхнего и нижнего регистра',
  numbers: 'Цифры',
  specialSymbols: 'Специальные символы'
};

const getNameOfPasswordRule = (ruleKey: PasswordRuleKey): string => {
  return passwordRuleLabels[ruleKey];
};

const passwordRules = reactive<PasswordRules>({
  minLength: undefined,
  lettersUpperAndLowerCase: undefined,
  numbers: undefined,
  specialSymbols: undefined
})

const allPasswordRulesPassed = computed(() => {
  return Object.values(passwordRules).every((val) => val === true)
})

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

watch(
  () => form.newPassword,
  (newPassword) => {
    validatePasswordDetail(newPassword);
  }
);
</script>

<style>
.profile-control {
  height: 100%;
}

.profile-control__form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  height: 100%;
}

.change-password-form__group {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.change-password-form__actions {
  display: flex;
  gap: 0.75rem;
  margin-top: auto;
}

.password-rules {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.password-rules__header {
  display: flex;
  gap: 0.375rem;
  align-items: center;
}

.password-rules__icon {
  max-width: 1.5rem;
  max-height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #A1A1AA;
}

.password-rules__icon svg {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.password-rules__icon_blue {
  color: #2563EB;
}

.password-rules__title {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.5rem;
}

.password-rules__list {
  display: flex;
  flex-direction: column;
  padding-left: 1.875rem;
}

.password-rules__rule {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #A1A1AA;
}

.password-rules__rule_valid {
  text-decoration: line-through;
  font-weight: 500;
  color: #2563EB;
}

.password-rules__rule_invalid {
  text-decoration: none;
  color: #C10007;

}
</style>
