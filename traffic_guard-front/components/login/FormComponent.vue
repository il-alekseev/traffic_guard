<template>
  <form class="login-form" @submit.prevent="handleSubmit">
    <div class="login-form-group">
      <BaseInput
        id="login"
        name="login"
        label="Логин"
        v-model="form.login"
        :error="errors.login"
      />

      <BaseInput
        id="password"
        name="password"
        type="password"
        label="Пароль"
        v-model="form.password"
        :error="errors.password"
        placeholder=""
        autocomplete="current-password"
      />
    </div>

    <div class="login-form-actions">
      <BaseButton
        type="submit"
        variant="login"
        :loading="props.loading"
        loadingText="Вход..."
      >
        Войти
      </BaseButton>
    </div>
  </form>
</template>

<script setup lang="ts">
import BaseInput from '~/components/ui/BaseInput.vue';
import BaseButton from '~/components/ui/BaseButton.vue';

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['login']);

const form = reactive({
  login: '',
  password: ''
});

const errors = reactive({
  login: '',
  password: ''
});

const validateForm = () => {
  let isValid = true;

  (Object.keys(errors) as (keyof typeof errors)[]).forEach(key => {
    errors[key] = '';
  });

  if (!form.login.trim()) {
    errors.login = 'Пожалуйста, введите логин';
    isValid = false;
  } else if (form.login.length < 3) {
    errors.login = 'Логин должен содержать не менее 3 символов';
    isValid = false;
  } else if (form.login.length > 255) {
    errors.login = 'Логин должен содержать менее 255 символов';
  } else if (!/^[a-zA-Z0-9_-]+$/.test(form.login)) {
    errors.login = 'Логин может содержать только латинские буквы, цифры и символ подчеркивания';
    isValid = false;
  }

  if (!form.password.trim()) {
    errors.password = 'Пожалуйста, введите пароль';
    isValid = false;
  }

  return isValid;
};
const handleSubmit = async () => {
  if (!validateForm()) {
    return;
  }

  emit('login', {
    login: form.login,
    password: form.password
  });
};
</script>

<style>
.login-form {
  padding: 3rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  background-color: #FFFFFF;
  box-shadow: 0px 1px 3px 0px #0000001A;
  border-radius: 1.5rem;
}

.login-form-group {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.login-form-actions {
  width: 100%;
}
</style>