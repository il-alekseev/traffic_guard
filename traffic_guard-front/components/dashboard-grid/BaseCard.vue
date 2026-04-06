<template>
  <div class="dashboard-card">
    <div v-if="isTopLegend" class="card__header">
      <div class="card__header-content">
        <h3 class="card__title" :class="props.titleClass" :style="props.titleStyle">{{ title }}</h3>
        <NuxtLink v-if="props.link" :to="link" class="details-btn">
          Подробнее
        </NuxtLink>
      </div>
      <div v-if="props.legend" :class="['chart-legend', `chart-legend--${legendPos}`]">
        <div v-for="(item, key) in props.legend" :key="key" class="legend__item">
          <div class="legend__dot" :style="{ borderColor: item.color }"></div>
          <div class="legend__name">{{ item.name }}</div>
        </div>
      </div>
    </div>

    <div v-else-if="!isTopLegend" class="card__header">
      <h3 class="card__title" :class="props.titleClass" :style="props.titleStyle">{{ title }}</h3>
      <NuxtLink v-if="props.link" :to="link" class="details-btn">
        Подробнее
      </NuxtLink>
      <slot name="titleStatus"></slot>
    </div>
    
    <div class="card__content">
      <slot name="DashboardStatistic"></slot>
      <slot name="LoadingData"></slot>
      <slot name="ErrorData"></slot>
      <slot name="EmptyData"></slot>
    </div>

    <div v-if="props.legend && !isTopLegend" :class="['chart-legend', 'chart-legend--bottom', `chart-legend--${legendPos}`]">
      <div v-for="(item, key) in props.legend" :key="key" class="legend__item">
        <div class="legend__dot" :style="{ borderColor: item.color }"></div>
        <div class="legend__name">{{ item.name }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DashboardTrafficLegend } from '~/types/dashboard';
import type { StyleValue } from 'vue';
  

const props = defineProps<{
  title: string;
  titleClass?: string,
  titleStyle?: StyleValue
  link?: string;
  legend?: DashboardTrafficLegend;
  legendPos?: 'top-left' | 'top-center' | 'top-right' | 'bottom-left' | 'bottom-center' | 'bottom-right'
}>();

const legendPos = computed(() => props.legendPos || 'top-right');

const isTopLegend = computed(() => {
  return legendPos.value.startsWith('top');
});
</script>

<style scoped lang="scss">
.dashboard-card {
  min-width: 32rem;
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
  align-items: flex-start;
  gap: 1rem;
}

.card__header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-grow: 1;
  justify-content: space-between;
}

.card__title {
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.25rem;
  color: #3F3F46;
  white-space: nowrap;
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
  white-space: nowrap;

  &:hover {
    cursor: pointer;
    text-decoration: underline;
  }
}

.chart-legend {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;

  &--top-left {
    justify-content: flex-start;
  }

  &--top-center {
    justify-content: center;
  }

  &--top-right {
    justify-content: flex-end;
  }

  &--bottom {
    margin-top: 0.75rem;
    width: 100%;

    &.chart-legend--bottom-left {
      justify-content: flex-start;
    }

    &.chart-legend--bottom-center {
      justify-content: center;
    }

    &.chart-legend--bottom-right {
      justify-content: flex-end;
    }
  }
}

.legend__item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.legend__dot {
  width: 0.5625rem;
  height: 0.5625rem;
  border-radius: 50%;
  flex-shrink: 0;
  border: 1.7px solid;
  background-color: transparent;
}

.legend__name {
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.5rem;
  color: #A1A1AA;
  white-space: nowrap;
}
</style>