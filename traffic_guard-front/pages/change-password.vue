<template>
  <div class="change-password">
    <div class="change-password__container">
      <div class="change-password__header">
        <div class="change-password__logo">
          <FazendaLogo />
        </div>
        <div class="change-password__text"> 
          <h2 class="change-password__title">Смена временного пароля</h2>
          <p class="change-password__subtitle">
            Ваш текущий пароль является временным.<br>
            Задайте новый пароль для продолжения работы.
          </p>
        </div>
      </div>

      <BaseAlert v-if="error">{{ error }}</BaseAlert>

      <FormComponent
        @change-password="changePassword"
        :loading="loading"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useUserStore } from "~/stores/user";
import BaseAlert from '~/components/ui/BaseAlert.vue';
import FormComponent from '~/components/change-password/FormComponent.vue';
import FazendaLogo from "~/assets/img/logo.svg"


definePageMeta({
  middleware: ['auth']
});


const router = useRouter();
const userStore = useUserStore();

const error = ref('');
const loading = ref(false);

const changePassword = async (password: string) => {
  error.value = '';
  loading.value = true;

  if (!userStore.user) {
    error.value = 'Произошла ошибка при смене пароля';
    return;
  }

  try {
    const oldPassword = localStorage.getItem('old_pass');

    if (!oldPassword) {
      router.push('/auth');
      return;
    }

    const success = await userStore.changePassword(oldPassword, password);

    if (success) {
      localStorage.clear();
      userStore.user.is_need_to_change_password = false;
      router.push('/dashboard');
    } else {
      error.value = 'Не удалось сменить пароль';
    }
  } catch (err) {
    console.error('Ошибка при смене пароля:', err);
    error.value = 'Произошла ошибка при смене пароля';
  } finally {
    loading.value = false;
  }
};

</script>

<style>
.change-password {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #FAFAFA;
}

@media screen and (max-width: 520px) {
  .change-password {
    padding: 1.5rem;
  }
}

.change-password__container {
  max-width: 30rem;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

@media screen and (max-width: 520px) {
  .change-password__container {
    padding: 2rem;
    padding-top: 1.5rem;
  }
  
}

.change-password__header {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  max-width: 28rem;
}

.change-password__logo {
  max-width: 3rem;
  max-height: 3rem;
  margin-left: auto;
  margin-right: auto;
}

.change-password__logo svg {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.change-password__text {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  text-align: center;
  color: #3F3F46;
}

.change-password__title {
  font-size: 1.5rem;
  line-height: 2rem;
  font-weight: 700;
}

@media screen and (max-width: 520px) {
  .change-password__title {
    font-size: 1.5rem;
  }
}

.change-password__subtitle {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 400;
  color: #6A7282;
}
</style>
