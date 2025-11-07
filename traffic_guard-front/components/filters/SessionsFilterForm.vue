<template>
  <div class="filters">
    <BaseAlert v-if="error">{{ error }}</BaseAlert>
    <form @submit.prevent="submitForm" class="filters__form">
      <div class="filters-form__group">

        <div class="base-form-field">
          <label for="status" class="base-label">Статус</label>
          <Dropdown
            :search="true"
            :localSearch="true"
            id="status"
            :items="allowedStatus"
            searchPlaceholder="Поиск..."
            v-model="form.status"
            :selectedTextClass="getBadgeClassByStatus(form.status.id as SessionStatus)"
          />
        </div>

        <div class="base-form-field">
          <label for="category" class="base-label">Категория</label>
          <Dropdown
            :search="true"
            :localSearch="true"
            id="category"
            :items="allowedCategories"
            searchPlaceholder="Поиск..."
            v-model="form.category"
            :selectedTextClass="`category-badge category-badge--${getModificatorByCategory(form.category.id as Categories)}`"
          />
        </div>

        <div class="base-form-field">
          <label for="types" class="base-label">Тип сессии</label>
          <Dropdown
            :search="true"
            :localSearch="true"
            id="types"
            :items="allowedTypes"
            searchPlaceholder="Поиск..."
            v-model="form.types"
          />
        </div>

        <div class="base-form-field">
          <label for="device" class="base-label">Устройство</label>
          <Dropdown
            :search="true"
            :localSearch="true"
            id="device"
            :items="allowedDevices"
            searchPlaceholder="Поиск..."
            v-model="form.device"
            :selectedTextClass="'ngfw-bagde'"
            :selectedTextStyle="{
              backgroundColor: generateColor(form.device.id).background,
              color: generateColor(form.device.id).color
            }"
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
import { getModificatorByCategory, getBadgeClassByStatus } from '~/helpers';
import { useDeviceColors } from '~/composables/useDeviceColors';
import { useDevicesStore } from '~/stores/devices';
import { useCategoriesStore } from '~/stores/categories';
import type { DropdownItem } from '~/types/dropdown';
import BaseButton from '~/components/ui/BaseButton.vue';
import BaseAlert from '~/components/ui/BaseAlert.vue';
import Dropdown from '~/components/ui/Dropdown.vue';
import type { Categories } from '~/types/categories';
import type { SessionStatus } from '~/types/session';


const emit = defineEmits<{
  setFilters: [filterData: SessionFilter];
  close: [];
}>();

const props = defineProps<{
  filtersData: SessionFilter | null
}>();


const { generateColor } = useDeviceColors();
const devicesStore = useDevicesStore();
const categoriesStore = useCategoriesStore()

const loading = ref(false);
const error = ref('');

export interface SessionFilter {
  status: DropdownItem,
  category: DropdownItem,
  types: DropdownItem,
  device: DropdownItem
}

const form = reactive<SessionFilter>({
  status: {
    id: '',
    name: ''
  },
  category: {
    id: '',
    name: ''
  },
  types: {
    id: '',
    name: ''
  },
  device: {
    id: '',
    name: ''
  },
});


const allowedStatus = ref<DropdownItem[]>([
  {
    id: '',
    name: 'Все',
  },
  {
    id: 'Разрешен',
    name: 'Разрешен',
  },
  {
    id: 'Запрещен',
    name: 'Запрещен',
  },
  {
    id: 'Ожидает',
    name: 'Ожидает',
  },
  {
    id: 'Аномалия',
    name: 'Аномалия',
  },
]);

const allowedCategories = computed<DropdownItem[]>(() => {
  const categories = categoriesStore.categories;
  if (!categories) return [];
  const result = categories.map((category) => {
    return {
      id: category,
      name: category
    }
  })
  return result;
})

const allowedTypes = ref<DropdownItem[]>([
  {
    id: '',
    name: 'Все',
  },
  {
    id: 'Разрешен',
    name: 'Разрешен',
  },
  {
    id: 'Запрещен',
    name: 'Запрещен',
  },
  {
    id: 'VPN',
    name: 'VPN',
  },
]);

const allowedDevices = computed<DropdownItem[]>(() => {
  const devices = devicesStore.devices;
  if (!devices) return [];
  const result = devices.map((device) => {
    return {
      id: device,
      name: device
    }
  })
  return result;
})

const resetFilters = () => {
  form.category = {id: '', name: ''}
  form.device = {id: '', name: ''}
  form.status = {id: '', name: ''}
}

const submitForm = () => {
  emit('setFilters', form)
  emit('close');
}

const initForm = async () => {
  if (props.filtersData === null) return;

  form.status = props.filtersData.status;
  form.category = props.filtersData.category;
  form.types = props.filtersData.types;
  form.device = props.filtersData.device;
}

onMounted(() => {
  initForm()
})
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
