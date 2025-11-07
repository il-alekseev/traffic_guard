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
import type { DashboardRequestObj } from '~/types/dashboard';

const props = defineProps<{
  data?: DashboardRequestObj,
  label: string
  graphHeight: string
}>()

const safeData = computed(() => ({
  data: props.data?.data?.Data ?? [],
  time: props.data?.data?.Time ?? []
}))

const series = computed(() => [
  {
    name: props.label,
    data: safeData.value.data
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
    background: 'transparent',
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
  grid: {
    show: false
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
      show: false,
      datetimeUTC: false,
      format: 'HH:mm'
    },
    tooltip: {
      enabled: false
    },
    axisTicks: {
      show: false,
    },
    axisBorder: {
      show: false,
    }
  },
  yaxis: {
    labels: {
      show: false,
      formatter: (val: number) => formatCompactNumber(val)
    }
  },
  legend: {
    show: false,
    position: 'top',
    horizontalAlign: 'right', 
  },
  tooltip: {
    x: {
      show: false,
      format: 'dd MMM',
      formatter: undefined,
    },
    y: {
      formatter: (val: number) => formatCompactNumber(val),
      title: {
          formatter: () => props.label
      },
    },
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
  colors: ['#37C84F']
}))
</script>

<style scoped>
div {
  width: 100%;
  margin-top: -30px;
}
</style>
