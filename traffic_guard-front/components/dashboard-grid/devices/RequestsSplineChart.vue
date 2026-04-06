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
import { formatCompactNumber } from '~/helpers';
import type { DashboardDeviceRequestsSplineChart } from '~/types/dashboard';


const props = defineProps<{
  data?: DashboardDeviceRequestsSplineChart
  graphHeight: string
}>()

const safeData = computed(() => ({
  blocked: props.data?.data?.blocked ?? [],
  pending: props.data?.data?.pending ?? [],
  time: props.data?.data?.time ?? []
}))

const series = computed(() => [
  {
    name: 'Заблокированные ресурсы',
    data: safeData.value.blocked
  },
  {
    name: ' Ожидающие решения',
    data: safeData.value.pending
  }
])

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
    tickAmount: safeData.value.time.length,

    labels: {
      show: true,
      rotate: 0,
      
      datetimeUTC: false,
      formatter: (value: string | number) => {
        const d = new Date(value);
        const day = String(d.getDate()).padStart(2, '0');
        const month = String(d.getMonth() + 1).padStart(2, '0');
        return `${day}.${month}`;
      }
    },

    axisTicks: {
      show: true
    },

    axisBorder: {
      show: true
    },

    tooltip: {
      enabled: false
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
  colors: ['#FB2904', '#EFB100']
}))

</script>

<style lang="scss">

</style>