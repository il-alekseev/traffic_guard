<template>
  <div class="report-card">
    <div class="card__header">
      <h3 class="card__title">{{ title }}</h3>
      <div v-if="props.legend" class="chart-legend">
        <div v-for="item in props.legend" class="legend__item">
          <div class="legend__dot" :style="{ backgroundColor: item.color }"></div>
          <div class="legend__name">{{ item.name }}</div>
        </div>
      </div>
    </div>
    
    <div class="card__content">
      <slot name="DashboardStatistic"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GraphLegend } from '~/types/reports';

const props = defineProps<{
  title: string;
  legend?: GraphLegend
}>();
</script>

<style scoped lang="scss">
.report-card {
  min-width: 48%;
  background: #FFFFFF;
  display: flex;
  flex-direction: column;
  height: fit-content;
}

.card__header {
  margin-bottom: 6px;
  padding-bottom: 6px;
  border-bottom: 1px solid #2563EB;
  width: 100%;
  display: flex;
  gap: 8px;
}

.card__title {
  font-size: 15px;
  line-height: 14px;
  font-weight: 400;
  color: #1D4ED8;
}

.chart-legend {
  display: flex;
  gap: 6px;
}

.legend__item {
  display: flex;
  align-items: center;
  gap: 2px;
}

.legend__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend__name {
  font-weight: 400;
  font-size: 8px;
  line-height: 100%;
  color: #A1A1AA;
  white-space: nowrap;
}

.card__content {
  flex-grow: 1;
}
</style>