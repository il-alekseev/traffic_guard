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
import type { RequestAnalyticShortData } from '~/types/reports';
import { formatCompactNumber } from '~/helpers';

const props = defineProps<{
  data?: DashboardRequestObj | RequestAnalyticShortData,
  label: string
  graphHeight: string
  showX?: boolean,
  showY?: boolean,
}>()

const isDashboardRequestObj = (
  d: DashboardRequestObj | RequestAnalyticShortData | undefined
): d is DashboardRequestObj => {
  return !!(d && typeof d === 'object' && 'type' in d);
};

const safeData = computed(() => {
  const d = props.data;

  if (!d) {
    return { data: [], time: [] };
  }

  if (isDashboardRequestObj(d)) {
    return {
      data: d.data?.data ?? [],
      time: d.data?.time ?? []
    }
  }

  return {
    data: d.data ?? [],
    time: d.time ?? []
  }
});

const series = computed(() => [
  {
    name: props.label,
    data: safeData.value.data
  }
])


const chartOptions = computed<ApexOptions>(() => ({
  chart: {
    background: 'transparent',
    type: 'area',
    sparkline: { enabled: true },
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
    tickAmount: safeData.value.time.length,

    labels: {
      show: props.showX === true,
      rotate: 0,
      
      datetimeUTC: false,
      formatter: (value: string | number) => {
        const d = new Date(value);
        const day = String(d.getDate()).padStart(2, '0');
        const month = String(d.getMonth() + 1).padStart(2, '0');
        return `${day}.${month}`;
      }
    },

    axisBorder: {
      show: props.showX === true
    },

    tooltip: {
      enabled: false
    }
  },

  yaxis: {
    labels: {
      show: props.showY === true,
      formatter: (val: number) =>
        props.showY ? formatCompactNumber(val) : ''
    },

    axisTicks: {
      show: props.showY === true
    },

    axisBorder: {
      show: props.showY === true
    }
  },

  legend: {
    show: false
  },

  tooltip: {
    x: {
      show: false,
      format: 'dd MM'
    },
    y: {
      formatter: (val: number) => formatCompactNumber(val),
      title: { formatter: () => props.label }
    }
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
}));

</script>

<style scoped>
div {
  width: 100%;
  margin-top: -1rem;
}
</style>
