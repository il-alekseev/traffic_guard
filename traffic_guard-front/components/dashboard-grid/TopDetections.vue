<template>
  <div class="top-detections">
    <div 
      v-for="(detection, index) in localDetections" 
      :key="index"
      class="detection__row"
    >
      <div class="detection__info">
        <div class="detection__position">
          <span class="detection__number">{{ index + 1 }}</span>
          <span class="trend-indicator" :class="getTrendClass(detection)">
            {{ getTrendSymbol(detection) }}
          </span>
        </div>
        <span class="detection__name">{{ detection.domain }}</span>
      </div>
      
      <div class="detections__bar">
        <div 
          class="bar-segment after" 
          :style="{ width: calculatePercentage(detection, 'after') + '%' }"
        >
        </div>
        <div 
          class="bar-segment before" 
          :style="{ width: calculatePercentage(detection, 'before') + '%' }"
        >
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { DashboardDetection, TopDetections } from '~/types/dashboard';

const props = defineProps<{
  detections: TopDetections;
}>();

const localDetections = ref<TopDetections>([]);

onMounted(() => {
  localDetections.value = props.detections.map(detection => ({
    ...detection,
    requests_before: 0,
    requests_after: 0,
    requests_all: 0
  }));

  setTimeout(() => {
    localDetections.value = [...props.detections];
  }, 50);  
})

const getTrendSymbol = (detection: any) => {
  if (detection.requests_after > detection.requests_before) return '▲';
  if (detection.requests_after < detection.requests_before) return '▼';
  return '–';
};

const getTrendClass = (detection: any) => {
  if (detection.requests_after > detection.requests_before) return 'trend-up';
  if (detection.requests_after < detection.requests_before) return 'trend-down';
  return 'trend-neutral';
};

const calculatePercentage = (detection: DashboardDetection, type: 'before' | 'after'): number => {
  const { requests_all } = detection;
  
  if (requests_all <= 0) return 0;
  
  const value = type === 'before' ? detection.requests_before : detection.requests_after;
  return (value / requests_all) * 100;
};


</script>

<style scoped lang="scss">
.top-detections {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.detection__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detection__position {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.trend-indicator {
  font-size: 0.875rem;
  line-height: 0.75rem;
  font-weight: 400;
  min-width: 0.875rem;
  
  &.trend-up {
    color: #37C84F;
  }
  
  &.trend-down {
    color: #FB2904;
  }
  
  &.trend-neutral {
    color: #9CA3AF;
  }
}

.detection__info {
  display: flex;
  align-items: center;
  gap: 0.875rem;
}

.detection__number {
  font-size: 0.875rem;
  line-height: 0.875rem;
  font-weight: 400;
  color: #3F3F46;
}

.detection__name {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 400;
  color: #2563EB;
}

.detections__bar {
  display: flex;
  width: 66%;
  height: 20px;
  overflow: hidden;
  gap: 0.375rem;
}

.bar-segment {
  border-radius: 6px;
  height: 100%;
  transition: width 0.35s ease;
  
  &.before {
    background-color: #FB2904;
  }
  
  &.after {
    background-color: #37C84F;
  }
}

</style>