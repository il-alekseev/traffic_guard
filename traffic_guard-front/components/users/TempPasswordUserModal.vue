<template>
  <BaseModal @close="closeModal">
    <template #header>
      <div class="modal__image modal-password-image">
        <SuccesCheckImage />
      </div>
      <div class="modal__header-title-block">
        <h3 class="modal__title">{{ props.type === 'create' ? 'Пользователь создан!' : 'Пароль успешно сброшен!' }}</h3>
        <p class="modal__text">
          {{ props.type === 'create' ? 'Направьте временный пароль пользователю' : 'Направьте новый пароль пользователю' }}
          
        </p>
      </div>
    </template>

    <template #main>
      <div class="password">
        <div class="password__value-container" @click="copyToClipboard">
          <p class="password__value">{{ password }}</p>
          <div class="password__copy-icon">
            <svg v-if="!copied" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            <CheckIcon v-else class="password__copy-icon_green" />
          </div>
        </div>
        <!-- <p v-if="copied" class="password__note">
          Пароль скопирован в буфер обмена
        </p> -->
      </div>
    </template>

    <template #footer>
      <BaseButton type="button" variant="secondary" :disabled="false" @click="closeModal">
        Закрыть
      </BaseButton>
      <BaseButton type="button" variant="primary" :disabled="false" @click="copyToClipboard">
        Скопировать
      </BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import type { User } from '~/types/user';
import BaseModal from '~/components/ui/BaseModal.vue';
import BaseButton from '~/components/ui/BaseButton.vue';
import CheckIcon from '~/assets/img/check.svg';
import SuccesCheckImage from "~/assets/img/succes-check-image.svg"

const props = defineProps<{
  user: User | null;
  password: string;
  type: 'create' | 'reset'
}>();

const emit = defineEmits(['close']);

const copied = ref(false);

const closeModal = () => {
  emit('close');
};

const copyToClipboard = async () => {
  if (typeof navigator === 'undefined' || !navigator.clipboard) {
    console.error('Clipboard API is not available');
    return;
  }
  try {
    await navigator.clipboard.writeText(props.password);
    copied.value = true;
  } catch (err) {
    console.error('Failed to copy text: ', err);
  }
};


onMounted(() => {
  copyToClipboard();
});
</script>

<style scoped>
.modal-password-image {
  width: 3rem;
  height: 3rem;
}

.modal__text-full-name {
  font-weight: 600;
}

.password {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.password__value-container {
  display: flex;
  align-items: center;
  text-align: center;
  justify-content: space-between;
  background-color: #FFFFFF;
  padding: 0.625rem;
  border-radius: 6px;
  cursor: pointer;
  outline: 1px solid #D4D4D8;
}

.password__value {
  text-align: center;
  width: 100%;
  font-size: 1rem;
  line-height: 1.25rem;
  font-weight: 400;
  color: #3F3F46;
}

.password__copy-icon {
  color: #A1A1AA;
  width: 1rem;
  height: 1rem;
}

.password__copy-icon svg {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.password__copy-icon_green {
  color: #3E9B4F;
}
/* 
.password__note {
  font-weight: 400;
  font-size: 0.75rem;
  line-height: 1rem;
  color: #A1A1AA;
} */
</style>
