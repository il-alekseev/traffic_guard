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
    
    <div
      v-if="props.extraTitles"
      class="card__extra-title"
    >
      <div class="card__title-value">
        {{ props.extraTitles.mainValue }}
      </div>

      <div class="card__title-dynamincs-contatiner" :class="extraTitles?.dynamicSign === 'positive' ? 'dynamic-positive' : 'dynamic-negative'">
        <div
          class="card__arrow"
          v-if="props.extraTitles.dynamicSign !== 'neutral'"
          :style="arrowStyle"
          v-html="arrowSvg"
        ></div>
        <div class="card__dynamic-percent-value">
          {{ props.extraTitles.percentValue }}
        </div>
      </div>
    </div>

    <div class="card__content">
      <slot name="DashboardStatistic"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { GraphLegend, RequestsGraphExtraTitle } from '~/types/reports';

const props = defineProps<{
  title: string;
  extraTitles?: RequestsGraphExtraTitle
  legend?: GraphLegend
}>();


const arrowSvg = `<svg width="14" height="12" viewBox="0 0 14 12" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M7 12L1.04907e-06 -1.22392e-06L14 0L7 12Z" fill="currentColor"/></svg>`;

const arrowStyle = computed(() => ({
  transform: props.extraTitles?.dynamicSign === 'negative' ? 'rotate(0deg)' : 'rotate(180deg)',
}));

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

.card__extra-title {
  padding-top: 12px;
  padding-left: 12px;
  margin-bottom: 40px;
  display: flex;
  gap: 8px;
}

.card__title-value {
 font-weight: 700;
 color: #3F3F46;
 font-size: 32px;
 line-height: 32px;
}

.card__title-dynamincs-contatiner {
  align-self: flex-end;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}

.card__arrow {
  height: 12px;
  width: 14px;
  font-size: 14px;
  line-height: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card__dynamic-percent-value {
  font-weight: 400;
  font-size: 14px;
  line-height: 14px;
}

.dynamic-positive {
  color: #37C84F;
}

.dynamic-negative {
  color: #FB2904;
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