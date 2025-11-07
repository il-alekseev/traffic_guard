<template>
  <form @submit.prevent="handleSubmit" class="report-form">
    <section class="form__section">
      <h2 class="section__title">Выберите временной диапазон</h2>
      <div class="quick-range-grid">
        <button
          v-for="range in quickRanges"
          :key="range.value"
          type="button"
          :class="['quick-range-btn', { 
            'active': selectedQuickRange === range.value,
          }]"
          @click="selectQuickRange(range.value)"
        >
          {{ range.label }}
        </button>
      </div>
      <div v-if="dateRangeError" class="error-message">
        {{ dateRangeError }}
      </div>
    </section>

    <section class="form__section">
      <h2 class="section__title">Настраиваемый диапазон</h2>
      <div class="custom-range-grid">
        <div class="date-picker-wrapper">
          <SingleDatePicker
            v-model="customDateFrom"
            :disabled="selectedQuickRange !== null"
            @update:modelValue="onCustomDateChange"
          />
        </div>
        <div class="date-picker-wrapper">
          <SingleDatePicker
            v-model="customDateTo"
            :disabled="selectedQuickRange !== null"
            @update:modelValue="onCustomDateChange"
          />
        </div>
      </div>
    </section>

    <section class="form__section">
      <h2 class="section__title">Устройства</h2>
      <div class="devices__options">
        <label class="radio__option">
          <input
            v-model="deviceSelection"
            type="radio"
            value="all"
            class="radio__input"
          />
          <div class="radio__content">
            <span class="radio__title">Все устройства</span>
            <span class="radio__description">Сгенерировать отчет для всех устройств</span>
          </div>
        </label>

        <label class="radio__option radio__option_disabled">
          <input
            v-model="deviceSelection"
            :disabled="true"
            type="radio"
            value="specific"
            class="radio__input"
          />
          <div class="radio__content">
            <span class="radio__title">Конкретное устройство</span>
            <span class="radio__description">Выберите конкретное устройство</span>
          </div>
        </label>

        <div v-if="deviceSelection === 'specific'" class="device-select-wrapper">
          <div class="base-form-field">
            <label for="device" class="base-label">Выберите устройство из списка</label>
            <Dropdown
              :search="true"
              :localSearch="true"
              id="device"
              :items="allowedDevices"
              searchPlaceholder="Поиск..."
              v-model="selectedDevice"
              :selectedTextClass="'ngfw-bagde'"
              :selectedTextStyle="{
                backgroundColor: generateColor(selectedDevice.id).background,
                color: generateColor(selectedDevice.id).color
              }"
              :class="{ 'dropdown-error': !!deviceError }"
            />
            <div v-if="deviceError" class="error-message">
              {{ deviceError }}
            </div>
          </div>
        </div>
      </div>
    </section>

    <BaseButton
      type="submit"
      variant="login"
      :loading="loading"
      :disabled="loading"
      loadingText="Отчет формируется..."
      :class="{ 'button-success': isSuccess, 'button-loading': loading }"
      :style="isSuccess ? { backgroundColor: '#05DF72' } : {}"
    >
      {{ isSuccess ? 'Отчет сформирован' : 'Сформировать отчёт' }}
    </BaseButton>
  </form>
</template>

<script setup lang="ts">
import { useDeviceColors } from '~/composables/useDeviceColors';
import { useDevicesStore } from '~/stores/devices';
import type { DropdownItem } from '~/types/dropdown';
import BaseButton from '~/components/ui/BaseButton.vue';
import Dropdown from '~/components/ui/Dropdown.vue';
import SingleDatePicker from '~/components/ui/SingleDatePicker.vue';
import { computed, ref, watch } from 'vue';
import { getCurrentDateWithOffset } from '~/helpers';
import type { ReportFormData } from '~/types/reports';

interface Props {
  loading?: boolean;
  success?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  success: false
})

const emit = defineEmits<{
  submit: [data: ReportFormData]
}>()

const { generateColor } = useDeviceColors();

type QuickRangeValue = 'week' | 'month' | 'quarter' | 'year'

interface QuickRange {
  value: QuickRangeValue
  label: string
}

interface ValidationErrors {
  dateRange?: string
  device?: string
}

const quickRanges: QuickRange[] = [
  { value: 'week', label: 'За последнюю неделю' },
  { value: 'month', label: 'За последний месяц' },
  { value: 'quarter', label: 'За последний квартал' },
  { value: 'year', label: 'За последний год' }
];

const devicesStore = useDevicesStore();

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
});

const selectedQuickRange = ref<QuickRangeValue | null>(null);
const customDateFrom = ref<Date | null>(getCurrentDateWithOffset(-1, 'd'))
const customDateTo = ref<Date | null>(getCurrentDateWithOffset())
const deviceSelection = ref<'all' | 'specific'>('all');
const selectedDevice = ref<DropdownItem>({
  id: '',
  name: ''
});
const errors = ref<ValidationErrors>({});
const isSuccess = ref(false);

const dateRangeError = computed(() => errors.value.dateRange);
const deviceError = computed(() => errors.value.device);

const customDatesSelected = computed(() => {
  return !!(customDateFrom.value || customDateTo.value)
});

const calculateDateRange = (range: QuickRangeValue): { from: string; to: string } => {
  const today = new Date()
  const to = today.toISOString().split('T')[0]
  const from = new Date(today)

  switch (range) {
    case 'week':
      from.setDate(today.getDate() - 7)
      break
    case 'month':
      from.setMonth(today.getMonth() - 1)
      break
    case 'quarter':
      from.setMonth(today.getMonth() - 3)
      break
    case 'year':
      from.setFullYear(today.getFullYear() - 1)
      break
  }

  return {
    from: from.toISOString().split('T')[0],
    to
  }
}

const selectQuickRange = (range: QuickRangeValue) => {
  if (selectedQuickRange.value === range) {
    selectedQuickRange.value = null
    customDateFrom.value = getCurrentDateWithOffset(1, 'd')
    customDateTo.value = getCurrentDateWithOffset()
    return
  }

  selectedQuickRange.value = range
  const { from, to } = calculateDateRange(range)

  customDateFrom.value = new Date(from)
  customDateTo.value = new Date(to)
  errors.value.dateRange = undefined
}


const onCustomDateChange = () => {
  if (customDatesSelected.value) {
    selectedQuickRange.value = null
    errors.value.dateRange = undefined
  }
}

const resetCustomDates = () => {
  customDateFrom.value = getCurrentDateWithOffset(1, 'd')
  customDateTo.value = getCurrentDateWithOffset()
}

const validateForm = (): boolean => {
  errors.value = {}
  let isValid = true

  if (!customDateFrom.value || !customDateTo.value) {
    errors.value.dateRange = 'Необходимо выбрать временной диапазон'
    isValid = false
  }

  if (deviceSelection.value === 'specific' && !selectedDevice.value.id) {
    errors.value.device = 'Необходимо выбрать устройство'
    isValid = false
  }

  return isValid
}

const handleSubmit = () => {
  if (!validateForm()) return

  const formData: ReportFormData = {
    dateFrom: customDateFrom.value?.toISOString().split('T')[0] ?? '',
    dateTo: customDateTo.value?.toISOString().split('T')[0] ?? '',
    deviceSelection: deviceSelection.value,
    ...(deviceSelection.value === 'specific' && { selectedDevice: selectedDevice.value.id })
  }

  emit('submit', formData)
}


watch(deviceSelection, () => {
  if (deviceSelection.value === 'all') {
    errors.value.device = undefined
  }
})

watch(selectedDevice, () => {
  if (selectedDevice.value.id) {
    errors.value.device = undefined
  }
})

watch(() => props.success, (newValue) => {
  if (newValue) {
    isSuccess.value = true
    setTimeout(() => {
      isSuccess.value = false
    }, 3000)
  }
})
</script>

<style scoped lang="scss">
.report-form {
  margin: 0 auto;
}

.form__section {
  margin-bottom: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.section__title {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 500;
  color: #3F3F46;
}

.error-message {
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: #EF4444;
  font-weight: 500;
}

.quick-range-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.quick-range-btn {
  padding: 0.5rem 2rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  color: #3F3F46;
  background: #ffffff;
  border: 1px solid #D4D4D8;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover:not(:disabled) {
    border-color: #2563EB;
  }

  &.active {
    border-color: #2563EB;
    background: #EFF6FF;
    color: #2563EB;
  }

  &.disabled,
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.custom-range-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.date-picker-wrapper {
  position: relative;
}

.reset-btn {
  font-size: 0.875rem;
  color: #2563EB;
  background: transparent;
  border: none;
  cursor: pointer;
  transition: opacity 0.2s ease;

  &:hover {
    opacity: 0.7;
  }
}

.devices__options {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.radio__option {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.5rem 0.75rem;
  background: #ffffff;
  border: 1px solid #D4D4D8;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover:not(:disabled) {
    border-color: #2563EB;
  }

  &_disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.radio__input {
  margin-top: 4px;
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: #2563EB;
}

.radio__content {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.radio__title {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 500;
  color: #3F3F46;
}

.radio__description {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 400;
  color: #A1A1AA;
}

.device-select-wrapper {
  margin-top: 0.5rem;
}

.base-form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.base-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #3F3F46;
}

.dropdown-error {
  border-color: #EF4444 !important;
}

.button-loading {
  opacity: 0.5;
  cursor: not-allowed;
}

.button-success {
  transition: background-color 0.3s ease;
}
</style>
