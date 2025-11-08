<template>
  <div>
    <ClientOnly>
      <VueApexCharts
        type="area"
        :height="props.graphHeight"
        :options="chartOptions"
        :series="series"
      />
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import type { DashboardTraffic } from '~/types/dashboard';

const props = defineProps<{
  trafficData?: DashboardTraffic
  graphHeight: string
}>()

const safeData = computed(() => ({
  input: props.trafficData?.data?.input ?? [],
  output: props.trafficData?.data?.output ?? [],
  time: props.trafficData?.data?.time ?? []
}))

const series = computed(() => [
  {
    name: 'Входящий',
    data: safeData.value.input
  },
  {
    name: 'Исходящий',
    data: safeData.value.output
  }
])

const formatCompactNumber = (val: number, locale: string = 'ru'): string => {
  return new Intl.NumberFormat(locale, {
    notation: 'compact',
    maximumFractionDigits: 1
  }).format(val);
};

const chartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'area',
    height: props.graphHeight,
    toolbar: { show: false },
    zoom: { enabled: false },
    locales: [{
      name: 'ru',
      options: {
        months: ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'],
        shortMonths: ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'],
        days: ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота'],
        shortDays: ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'],
        toolbar: {
          exportToSVG: 'Скачать SVG',
          exportToPNG: 'Скачать PNG',
          exportToCSV: 'Скачать CSV',
          menu: 'Меню',
          selection: 'Выбор',
          selectionZoom: 'Выбор и масштабирование',
          zoomIn: 'Приблизить',
          zoomOut: 'Отдалить',
          pan: 'Перемещение',
          reset: 'Сбросить масштаб'
        }
      }
    }],
    defaultLocale: 'ru',
  },
  dataLabels: {
    enabled: false
  },
  stroke: {
    curve: 'smooth',
    width: 2
  },
  xaxis: {
    type: 'datetime',
    categories: safeData.value.time,
    labels: {
      datetimeUTC: false,
      format: 'HH:mm'
    }
  },
  yaxis: {
    labels: {
      formatter: (val: number) => formatCompactNumber(val)
    }
  },
  legend: {
    show: false,
    position: 'top',
    horizontalAlign: 'right', 
  },
  tooltip: {
    x: { format: 'dd MMM HH:mm' }
  },
  fill: {
    type: 'gradient',
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.4,
      opacityTo: 0.1,
      stops: [0, 90, 100]
    }
  },
  colors: ['#2563EB', '#37C84F']
}))
</script>

<style scoped>
div {
  width: 100%;
}
</style>
