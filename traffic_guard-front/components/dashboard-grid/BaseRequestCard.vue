<template>
  <div class="dashboard-card">
    <div class="card__header">
      <h3 class="card__title">{{ title }}</h3>
      <div class="card__data-values">
        <div class="card__current-value">
          {{ formatNumber(props.data?.currentValue) }}
        </div>

        <div
          v-if="props.data"
          class="card__old-value-container"
          :class="{
            'card__old-value-container_red': percentChange < 0,
            'card__old-value-container_green': percentChange > 0,
            'card__old-value-container_gray': percentChange === 0
          }"
        >
          <div
            class="card__arrow"
            v-if="percentChange !== 0"
            :style="arrowStyle"
            v-html="arrowSvg"
          ></div>
          <div class="card__old-value">
            {{ formatPercent(percentChange) }}
          </div>
        </div>
      </div>
    </div>

    <div class="card__content">
      <slot name="DashboardStatistic"></slot>
      <slot name="LoadingData"></slot>
      <slot name="ErrorData"></slot>
      <slot name="EmptyData"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  title: string;
  data?: {
    currentValue: number;
    oldValue: number;
  };
}>();

function formatNumber(value?: number): string {
  if (value == null) return '-';
  if (value >= 1_000_000) return (value / 1_000_000).toLocaleString('ru-RU', { maximumFractionDigits: 1 }) + ' млн';
  if (value >= 1_000) return (value / 1_000).toLocaleString('ru-RU', { maximumFractionDigits: 1 }) + ' тыс';
  return value.toLocaleString('ru-RU');
}

const percentChange = computed(() => {
  if (!props.data || props.data.oldValue === 0) return 0;
  return ((props.data.currentValue - props.data.oldValue) / props.data.oldValue) * 100;
});

function formatPercent(value: number): string {
  const sign = value > 0 ? '+' : '';
  return `${sign}${value.toFixed(1)}%`;
}

const arrowSvg = `<svg width="14" height="12" viewBox="0 0 14 12" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M7 12L1.04907e-06 -1.22392e-06L14 0L7 12Z" fill="currentColor"/></svg>`;

const arrowStyle = computed(() => ({
  transform: percentChange.value < 0 ? 'rotate(0deg)' : 'rotate(180deg)',
}));
</script>

<style scoped lang="scss">
.dashboard-card {
  width: 32rem;
  height: 6.75rem;
  background: #fff;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  gap: 1.5rem;
}

.card__header {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.card__title {
  width: 13.375rem;
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.25rem;
  color: #3F3F46;
}

.card__data-values {
  display: flex;
  gap: 0.5rem;
}

.card__current-value {
  font-weight: 700;
  font-size: 2rem;
  line-height: 2rem;
  color: #3F3F46;
}

.card__old-value-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.25rem;
  color: #3F3F46;
  align-self: flex-end;
  margin-bottom: 0.125rem;
}

.card__old-value-container_red {
  color: #FB2904;
}

.card__old-value-container_green {
  color: #37C84F;
}

.card__old-value-container_gray {
  color: #A1A1AA;
}

.card__arrow {
  height: 0.75rem;
  width: 0.875rem;
  font-size: 0.875rem;
  line-height: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card__old-value {
  font-weight: 400;
  font-size: 0.75rem;
  line-height: 0.75rem;
}

.card__content {
  flex-grow: 1;
  max-width: 14.125rem;
  height: 5.625rem;
}
</style>
