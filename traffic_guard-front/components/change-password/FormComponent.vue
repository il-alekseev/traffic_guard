<template>
  <form class="change-password-form" @submit.prevent="handleSubmit">
    <div class="change-password-form__group">
      <BaseInput
        id="password"
        name="password"
        type="password"
        label="Введите новый пароль"
        v-model="form.password"
        :error="errors.password"
        :showTextError="showTextErrors.password"
        placeholder=""
        autocomplete="current-password"
      />


      <BaseInput
        id="passwordRepeat"
        name="passwordRepeat"
        type="password"
        label="Подтверждение пароля"
        v-model="form.passwordRepeat"
        :error="errors.passwordRepeat"
        :showTextError="showTextErrors.passwordRepeat"
        placeholder=""
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
        type="submit"
        variant="login"
        :loading="props.loading"
        loadingText="Смена пароля..."
      >
        Сохранить
      </BaseButton>
    </div>
  </form>
</template>

<script setup lang="ts">
import BaseInput from '~/components/ui/BaseInput.vue';
import BaseButton from '~/components/ui/BaseButton.vue';
import CheckMarkCircleIcon from '~/assets/img/checkmark-circle.svg'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['change-password']);

interface PasswordRules {
  minLength: undefined | boolean,
  lettersUpperAndLowerCase: undefined | boolean,
  numbers: undefined | boolean,
  specialSymbols: undefined | boolean,
}

type PasswordRuleKey = keyof PasswordRules;

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

const form = reactive({
  password: '',
  passwordRepeat: ''
});

const errors = reactive({
  password: '',
  passwordRepeat: ''
});

const showTextErrors = reactive({
  password: true,
  passwordRepeat: true
})

const validatePasswordDetail = (password: string) => {
  const pwd = password.trim();
  showTextErrors.password = false;

  passwordRules.minLength = pwd.length >= 8;
  passwordRules.lettersUpperAndLowerCase = /[a-z]/.test(pwd) && /[A-Z]/.test(pwd);
  passwordRules.numbers = /\d/.test(pwd);
  passwordRules.specialSymbols = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd);

  if (!allPasswordRulesPassed.value) {
    errors.password = 'Пожалуйста, введите корректный пароль';
  } else {
    errors.password = '';
  }
}

const validateForm = () => {
  let isValid = true;

  (Object.keys(errors) as (keyof typeof errors)[]).forEach(key => {
    errors[key] = '';
  });

  const pwd = form.password.trim()
  if (!pwd) {
    showTextErrors.password = true;
    errors.password = 'Пожалуйста, введите новый пароль';
    isValid = false;
  } else {
    showTextErrors.password = false;

    passwordRules.minLength = pwd.length >= 8;
    passwordRules.lettersUpperAndLowerCase = /[a-z]/.test(pwd) && /[A-Z]/.test(pwd);
    passwordRules.numbers = /\d/.test(pwd);
    passwordRules.specialSymbols = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(pwd);

    if (!allPasswordRulesPassed.value) {
      errors.password = 'Пожалуйста, введите корректный пароль';
      isValid = false;
    }
  }

  if (!form.passwordRepeat.trim()) {
    showTextErrors.passwordRepeat = true;
    errors.passwordRepeat = 'Пожалуйста, подтвердите пароль';
    isValid = false;
  } else if (form.password !== form.passwordRepeat) {
    showTextErrors.passwordRepeat = true;
    errors.passwordRepeat = 'Пароли не совпадают';
    isValid = false;
  }

  return isValid;
};

watch(
  () => form.password,
  (newPassword, _oldPassword) => {
    validatePasswordDetail(newPassword);
  }
);

const handleSubmit = async () => {
  if (!validateForm()) {
    return;
  }

  emit('change-password', form.password);
};
</script>

<style>
.change-password-form {
  padding: 3rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  background-color: #FFFFFF;
  box-shadow: 0px 1px 3px 0px #0000001A;
  border-radius: 1.5rem;
}

.change-password-form__group {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
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

.change-password-form__actions {
  width: 100%;
}
</style>