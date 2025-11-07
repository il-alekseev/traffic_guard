<template>
  <div class="dashboard-card">
    <div class="card__header">
      <h3 class="card__title">{{ title }}</h3>
      <NuxtLink v-if="props.link" :to="link" class="details-btn">
        Подробнее
      </NuxtLink>
      <div v-if="props.legend" class="chart-legend">
        <div v-for="(item, key) in props.legend" :key="key" class="legend__item">
          <div class="legend__dot" :style="{ backgroundColor: item.color }"></div>
          <div class="legend__name">{{ item.name }}</div>
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
import type { DashboardTrafficLegend } from '~/types/dashboard';

const props = defineProps<{
  title: string;
  link?: string;
  legend?: DashboardTrafficLegend
}>();
</script>

<style scoped lang="scss">
.dashboard-card {
  min-width: 512px;
  background: #fff;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
}

.card__header {
  margin-bottom: 0.75rem;
  width: 100%;
  display: flex;
  justify-content: space-between;
}

.card__title {
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.25rem;
  color: #3F3F46;
}

.card__content {
  flex-grow: 1;
}

a.details-btn {
  color: #2563EB;
  text-decoration: none;
  font-size: 1rem;
  line-height: 1.25rem;
  font-weight: 400;
  

  &:hover {
    cursor: pointer;
    text-decoration: underline;
  }
}

.chart-legend {
  display: flex;
  gap: 0.75rem;
}

.legend__item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.legend__dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend__name {
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.5rem;
  color: #A1A1AA;
  white-space: nowrap;
}
</style>