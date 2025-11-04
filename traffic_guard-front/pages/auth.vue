<template>
  <div class="login">
    <div class="login__container">
      <div class="login__header">
        <div class="login__logo">
          <FazendaLogo />
        </div>
        <h2 class="login__title">Вход в Fazenda</h2>
      </div>

      <BaseAlert v-if="error">{{ error }}</BaseAlert>

      <FormComponent
        @login="handleLogin"
        :loading="loading"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useUserStore } from "~/stores/user";
import BaseAlert from '~/components/ui/BaseAlert.vue';
import FormComponent from '~/components/login/FormComponent.vue';
import FazendaLogo from "~/assets/img/logo.svg"

const router = useRouter();
const userStore = useUserStore();

const error = ref('');
const loading = ref(false);

type LoginData = {
  login: string;
  password: string;
};

const handleLogin = async ({ login, password }: LoginData) => {
  error.value = '';
  loading.value = true;

  try {
    const success = await userStore.login(login, password);

    if (success) {
      const user = await userStore.fetchUserInfo();
      if (!user.is_need_to_change_password) {
        router.push('/dashboard');
      } else {
        localStorage.setItem('old_pass', password);
        router.push('/change-password');
      }
    } else {
      error.value = 'Неверный логин или пароль';
    }
  } catch (err) {
    console.error('Ошибка при входе:', err);
    error.value = 'Произошла ошибка при входе в систему';
  } finally {
    loading.value = false;
  }
};
</script>

<style>
.login {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #FAFAFA;
}

@media screen and (max-width: 520px) {
  .login {
    padding: 1.5rem;
  }
}

.login__container {
  max-width: 30rem;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

@media screen and (max-width: 520px) {
  .login__container {
    padding: 2rem;
    padding-top: 1.5rem;
  }
  
}

.login__header {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.login__logo {
  max-width: 3rem;
  max-height: 3rem;
  margin-left: auto;
  margin-right: auto;
}

.login__logo svg {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.login__title {
  font-size: 1.5rem;
  line-height: 2rem;
  font-weight: 700;
  color: #3F3F46;
}

@media screen and (max-width: 520px) {
  .login__title {
    font-size: var(--font-size-2xl);
  }
  
}

.login__subtitle {
  font-size: var(--font-size-xs);
  line-height: var(--line-height-xs-tight);
  font-weight: 400;
}
</style>