<template>
  <div class="anomalies">
    <div class="anomalies__container">
      <div class="anomalies__chart-wrapper">
        <ClientOnly>
          <VueApexCharts
            :options="chartOptions"
            :series="chartSeriesWithFiller"
            type="donut"
            :height="props.graphHeight"
          />
        </ClientOnly>
        
        <div class="anomalies__stats">
          <div class="anomalies__percentage">{{ anomalyPercentage }}%</div>
          <div class="anomalies__count">{{ totalAnomalyCount }} ресурсов</div>
        </div>
      </div>

      <div class="anomalies__domains-wrapper">
        <div class="anomalies__domains-list">
          <template v-if="allDomains.length > 0">
            <div
              v-for="(domain, index) in allDomains"
              :key="domain"
              class="anomalies__domain-item"
            >
              <div 
                class="anomalies__domain-indicator"
                :style="{ backgroundColor: chartColors[index % chartColors.length] }"
              />
              <span class="anomalies__domain-name" :style="{ color: chartColors[index % chartColors.length] }">{{ domain }}</span>
            </div>
          </template>
          <template v-else>
            <p class="anomalies-empty">Данные отсутствуют</p>
          </template>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'

export interface HostAnomaly {
  anomaly_count: number
  hostname: string
  domains: string[]
}

export interface DashboardAnomalies {
  block_count: number
  host_count: number
  host_anomalies: HostAnomaly[]
}

const props = defineProps<{
  anomaliesData?: DashboardAnomalies
  graphHeight: string
}>()

const chartColors = [
  '#FFA726',
  '#EF5350',
  '#AB47BC',
  '#29B6F6',
  '#26A69A',
  '#FFCA28',
  '#8D6E63',
  '#BDBDBD',
]

const totalAnomalyCount = computed(() => {
  if (!props.anomaliesData) return 0
  return props.anomaliesData.host_anomalies.reduce(
    (sum, host) => sum + (host?.anomaly_count || 0),
    0
  )
})

const anomalyPercentage = computed(() => {
  if (!props.anomaliesData || !props.anomaliesData.block_count || !props.anomaliesData.host_count) {
    return 0
  }
  const totalRules = props.anomaliesData.block_count * props.anomaliesData.host_count
  const percentage = Math.round((totalAnomalyCount.value / totalRules) * 100)
  return Math.min(percentage, 100)
})

const allDomains = computed(() => {
  if (!props.anomaliesData) return []
  const domainsSet = new Set<string>()
  props.anomaliesData.host_anomalies.forEach((host) => {
    host.domains.forEach((domain) => {
      domainsSet.add(domain)
    })
  })
  return Array.from(domainsSet)
})

const chartSeries = computed(() => {
  if (!props.anomaliesData) return []
  return props.anomaliesData.host_anomalies.map((host) => host.anomaly_count || 0)
})

const chartSeriesWithFiller = computed(() => {
  const series = chartSeries.value
  const sum = series.reduce((a, b) => a + b, 0)
  const fillerValue = Math.max(0, 100 - anomalyPercentage.value)
  return fillerValue > 0 ? [...series, fillerValue] : series
})

const chartLabels = computed(() => {
  if (!props.anomaliesData) return []
  const labels = props.anomaliesData.host_anomalies.map((host) => host.hostname)
  const fillerValue = Math.max(0, 100 - anomalyPercentage.value)
  return fillerValue > 0 ? [...labels, 'Empty'] : labels
})

const chartColorsWithFiller = computed(() => {
  const colors = chartColors.slice(0, chartSeries.value.length)
  const fillerValue = Math.max(0, 100 - anomalyPercentage.value)
  return fillerValue > 0 ? [...colors, '#E5E7EB'] : colors
})

const chartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'donut',
    toolbar: {
      show: false,
    },
    sparkline: {
      enabled: false,
    },
  },
  colors: chartColorsWithFiller.value,
  labels: chartLabels.value,
  legend: {
    show: false,
  },
  plotOptions: {
    pie: {
      donut: {
        size: '75%',
        background: 'transparent',
        labels: {
          show: false,
        },
      },
      startAngle: -90,
      endAngle: 90,
    },
  },
  dataLabels: {
    enabled: false,
  },
  tooltip: {
    enabled: true,
  },
  responsive: [
    {
      breakpoint: 480,
      options: {
        chart: {
          width: 200,
        },
      },
    },
  ],
}))
</script>

<style scoped lang="scss">
.anomalies {
  width: 100%;
  height: 100%;
  max-width: 29rem;
  max-height: 9.25rem;
  display: flex;
  align-items: center;
  justify-content: center;

  @media screen and (max-width: 1919px) {
    max-width: unset;
  }

}

.anomalies__container {
  display: flex;
  gap: 0.75rem;
  width: 100%;
  height: 100%;

  @media screen and (max-width: 1919px) {
    justify-content: space-between;
  }
}

.anomalies__chart-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  max-width: 15.75rem;
  margin-top: 8.75rem;
}

.anomalies__stats {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -100%);
  text-align: center;
  z-index: 10;
}

.anomalies__percentage {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
  margin-bottom: 0.25rem;
}

.anomalies__count {
  font-size: 0.875rem;
  color: #9ca3af;
  font-weight: 500;
}

.anomalies__domains-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #F4F4F5;
  border-radius: 12px;
  min-width: 0;
  overflow: hidden;
  max-height: 9.25rem;
  max-width: 12.5rem;

  @media screen and (max-width: 1919px) {
    max-width: unset;
  }
}

.anomalies__domains-list {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
  overflow-y: auto;
  padding: 0.75rem 0.25rem 0.75rem 0.75rem;

  &::-webkit-scrollbar {
    width: 0.375rem;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background-color: #d1d5db;
    border-radius: 3px;

    &:hover {
      background-color: #9ca3af;
    }
  }
}

.anomalies__domain-item {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.anomalies__domain-indicator {
  flex-shrink: 0;
  width: 0.75rem;
  height: 0.75rem;
  border-radius: 50%;
}

.anomalies__domain-name {
  font-size: 0.875rem;
  color: #1f2937;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media screen and (min-width: 1921px) {
  .anomalies {
    max-width: 100%;
  }

  .anomalies__container {
    width: 100%;
    justify-content: space-between;
  }

  .anomalies__domains-wrapper {
    max-width: 35rem;
  }
}

@media (max-width: 768px) {
  .anomalies__container {
    flex-direction: column;
    gap: 1rem;
  }

  .anomalies__chart-wrapper {
    min-width: 100%;
  }

  .anomalies__percentage {
    font-size: 2rem;
  }
}
</style>
