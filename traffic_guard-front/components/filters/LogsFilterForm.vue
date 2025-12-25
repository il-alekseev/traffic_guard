<template>
  <div class="filters">
    <BaseAlert v-if="error">{{ error }}</BaseAlert>
    <form @submit.prevent="submitForm" class="filters__form">
      <div class="filters-form__group">

        <div class="base-form-field">
          <label for="status" class="base-label">Роль</label>
          <Dropdown
            :search="true"
            :localSearch="true"
            id="role"
            :items="allowedRoles"
            searchPlaceholder="Поиск..."
            v-model="form.role"
          />
        </div>

        <div class="base-form-field">
          <label for="category" class="base-label">Контекст</label>
          <Dropdown
            :search="true"
            :localSearch="true"
            id="category"
            :items="allowedContexts"
            searchPlaceholder="Поиск..."
            v-model="form.contextId"
          />
        </div>
      </div>

      <div class="filters-form__actions actions">
        <button type="button" class="actions__extra" @click="resetFilters">
          Очистить фильтры
        </button>

        <div class="actions__main">
          <BaseButton
            class="cancel-btn"
            type="button"
            variant="secondary"
            :disabled="loading"
            loadingText="Cохранение..."
            @click="emit('close')"
          >
            Отменить
          </BaseButton>

          <BaseButton
            class="submit-btn"
            type="submit"
            variant="primary"
            :loading="loading"
            :loadingText="'Загрузка...'"
          >
            Применить
          </BaseButton>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useDeviceColors } from '~/composables/useDeviceColors';
import { useDevicesStore } from '~/stores/devices';
import { useCategoriesStore } from '~/stores/categories';
import type { DropdownItem } from '~/types/dropdown';
import BaseButton from '~/components/ui/BaseButton.vue';
import BaseAlert from '~/components/ui/BaseAlert.vue';
import Dropdown from '~/components/ui/Dropdown.vue';
import type { LogFilter } from '~/types/logs';


const emit = defineEmits<{
  setFilters: [filterData: LogFilter];
  close: [];
}>();

const props = defineProps<{
  filtersData: LogFilter | null
}>();


const { generateColor } = useDeviceColors();
const devicesStore = useDevicesStore();

const loading = ref(false);
const error = ref('');

const form = reactive<LogFilter>({
  role: {
    id: '',
    name: ''
  },
  contextId: {
    id: '',
    name: ''
  },
});

const allowedRoles = [
  {
    id: 'SA',
    name: 'Системный администратор'
  },
  {
    id: 'CA',
    name: 'Администратор узла NGFW'
  }
];

const allowedContexts = [
  {
    id: '1',
    name: '1'
  },
  {
    id: '2',
    name: '2'
  },
  {
    id: '3',
    name: '3'
  },
];

const resetFilters = () => {
  form.role = {id: '', name: ''}
  form.contextId = {id: '', name: ''}
}

const submitForm = () => {
  emit('setFilters', form)
  emit('close');
}

const initForm = async () => {
  if (props.filtersData === null) return;

  form.role = props.filtersData.role;
  form.contextId = props.filtersData.contextId;
}

const onGlobalKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter') {
    e.preventDefault();
    submitForm();
  }
};

onMounted(() => {
  initForm();
  window.addEventListener('keydown', onGlobalKeydown);
})

onBeforeUnmount(() => {
  window.addEventListener('keydown', onGlobalKeydown);
});
</script>

<style lang="scss" scoped>
.filters {
  height: 100%;
}

.filters .base-alert {
  margin: 1.5rem 1.5rem;
}

.info-block {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-block__title {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
}

.info-block__text {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #71717A;
}

.filters__form {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.filters-form__group {
  padding: 0 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.actions {
  margin-top: auto;
  display: flex;
  justify-content: space-between;
  border-top: 1px solid #E4E4E7;
}

.actions__extra {
  border: none;
  background-color: #FFFFFF;
  color: #A1A1AA;
  font-weight: 500;
  font-size: 0.875rem;
  line-height: 1.25rem;
  margin-left: 1.5rem;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  align-self: center;
}

.actions__main {
  display: flex;
  gap: 0.5rem;
  margin-left: auto;
  margin-top: auto;
  padding: 1rem;
  justify-content: flex-end;
}

.extras {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.extras__title {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
}

.extras__btns {
  display: flex;
  gap: 0.5rem;
}

.cancel-btn {
  max-width: 5.8125rem;
}

.submit-btn {
  max-width: 13rem;
}

.user-form__error-text {
  font-size: var(--font-size-xs);
  color: #EF4444;
}
</style>
