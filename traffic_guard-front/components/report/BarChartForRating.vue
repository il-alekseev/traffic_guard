<template>
  <div class="report-chart">
    <div 
      v-for="(item, index) in items" 
      :key="index"
      class="report-chart__row"
    >
      <div class="report-chart__info">
        <span class="report-chart__name">{{ item.name || '–' }}</span>
        <span v-if="showSubcategories && item.categories" class="report-chart__subcategories">
          {{ item.categories.join('/') }}
        </span>
      </div>
      
      <div class="report-chart__bars">
        <div 
          class="report-chart__bar-group"
          :style="{ width: calculatePercentage(item, 'before') + '%' }"
        >
          <div 
            class="report-chart__bar report-chart__bar--before"
          >
            <span class="report-chart__dot"></span>
            <span class="report-chart__value">{{ item.before_block || '–' }}</span>
          </div>
        </div>

        <div 
          class="report-chart__bar-group"
          :style="{ width: calculatePercentage(item, 'after') + '%' }"
        >
          <div 
            class="report-chart__bar report-chart__bar--after"
          >
            <span class="report-chart__dot"></span>
            <span class="report-chart__value">{{ item.after_block || '–' }}</span>
          </div>
        </div>
        
        <div 
          class="report-chart__bar-group"
          :style="{ width: calculatePercentage(item, 'waiting') + '%' }"
        >
          <div 
            class="report-chart__bar report-chart__bar--waiting"
          >
            <span class="report-chart__dot"></span>
            <span class="report-chart__value">{{ item.waiting || '–' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface ChartItem {
  name: string
  before_block: number
  after_block: number
  waiting: number
  categories?: string[]
}

interface Props {
  items: ChartItem[]
  showSubcategories?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showSubcategories: false
})

const calculatePercentage = (item: ChartItem, type: 'before' | 'after' | 'waiting'): number => {
  const requestAll = props.items[0].before_block + props.items[0].after_block + props.items[0].waiting;
  
  if (requestAll <= 0) return 0;
  
  const value = type === 'before' ? item.before_block : type === 'after' ? item.after_block : item.waiting;
  let result = (value / requestAll) * 100;
  if (result < 10) {
    result = 10
  }
  return result;
};
</script>

<style scoped lang="scss">
.report-chart {
  display: flex;
  flex-direction: column;
  gap: 8px;

  &__row {
    display: flex;
    justify-content: flex-start;
    align-items: center;
    gap: 12px;
  }

  &__info {
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    width: 124px;
  }

  &__name {
    font-size: 9px;
    line-height: 14px;
    font-weight: 400;
    color: #3F3F46;
  }

  &__subcategories {
    font-size: 6px;
    line-height: 6px;
    font-weight: 400;
    color: #3F3F46;
  }

  &__bars {
    width: 100%;
    display: flex;
    gap: 4px;
    height: 16px;
  }

  &__bar-group {
    display: flex;
    min-height: 16px;
  }

  &__bar {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    border-radius: 4px;
    height: 100%;
    width: 100%;
    position: relative;
    padding: 2px 6px;
    gap: 4px;

    &--after {
      background-color: #fb2904;
    }

    &--before {
      background-color: #37c84f;
    }

    &--waiting {
      background-color: #facc15;
    }
  }

  &__dot {
    display: inline-flex;
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background-color: #ffffff;
    flex-shrink: 0;
  }

  &__value {
    font-size: 10px;
    line-height: 10px;
    font-weight: 500;
    color: #ffffff;
    white-space: nowrap;
  }
}
</style>
