<template>
  <div class="filter-chips">
    <button
      v-for="option in options"
      :key="option.value"
      @click="handleSelect(option.value)"
      :class="['chip', { 'chip--active': isSelected(option.value) }]"
    >
      <span
        :style="{ backgroundColor: option.color }"
        class="chip__indicator"
      />
      <span class="chip__label">{{ option.label }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import type { ActionFilterOption } from '~/types/detections';


interface Props {
  options: ActionFilterOption[]
  modelValue?: 'Разрешено' | 'Заблокировано' | 'Не решено' | undefined
}

interface Emits {
  (e: 'update:modelValue', value: Props['modelValue']): void
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const handleSelect = (value: ActionFilterOption['value']) => {
  if (isSelected(value)) emit('update:modelValue', undefined);
  else emit('update:modelValue', value);
}

const isSelected = (value: ActionFilterOption['value']) => {
  return value === props.modelValue;
}
</script>

<style scoped lang="scss">
.filter-chips {
  display: flex;
  background-color: #FFFFFF;
  box-shadow: 0px 1px 2px 0px #0000000D;
  border: 0px solid transparent;
  border-radius: 6px;
}

.chip {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.4375rem 0.75rem;
  background-color: #FFFFFF;
  border: 1px solid transparent;
  font-weight: 500;

  &:first-child {
    border-top-left-radius: 6px;
    border-bottom-left-radius: 6px;
  }

  &:last-child {
    border-top-right-radius: 6px;
    border-bottom-right-radius: 6px;
  }

  &:nth-child(2) {
    border-left: 1px solid #E4E4E7;
    border-right: 1px solid #E4E4E7;
  }

  &__label {
    font-size: 0.875rem;
    line-height: 1.25rem;
    color: #3F3F46;
  }

  &--active {
    background-color: #E4E4E7;
  }
}

.chip__indicator {
  width: 0.375rem;
  height: 0.375rem;
  border-radius: 50%;
  flex-shrink: 0;
}

.chip__label {
  white-space: nowrap;
}
</style>
