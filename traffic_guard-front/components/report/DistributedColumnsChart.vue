<template>
  <div>
    <ClientOnly>
      <VueApexCharts
        width="100%"
        type="bar"
        :height="props.height"
        :options="chartOptions"
        :series="series"
      />
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import { computed } from 'vue';
import { formatCompactNumber } from '~/helpers';

const props = defineProps<{
  data: { name: string; value: number; color: string }[];
  height: string
}>();

const safeData = computed(() => ({
  categories: props.data.map(i => i.name),
  colors: props.data.map(i => i.color),
  data: props.data.map(i => i.value)
}));

const series = computed(() => [
  {
    name: "",
    data: safeData.value.data
  }
]);

const chartOptions = computed<ApexOptions>(() => ({
  chart: {
    type: 'bar',
    toolbar: { show: false },
    sparkline: { enabled: false }
  },

  plotOptions: {
    bar: {
      distributed: true,
      horizontal: false,
      borderRadius: 6,
      columnWidth: '50%'
    }
  },

  colors: safeData.value.colors,

  xaxis: {
    categories: safeData.value.categories,
    labels: { show: false },
    axisTicks: { show: false },
    axisBorder: { show: false },
    floating: true
  },

  yaxis: {
    labels: {
      show: true,
      formatter: (value: number) => formatCompactNumber(Math.round(value))
    },
    title: {
      text: '',
      style: { fontSize: '12px' }
    }
  },

  grid: {
    show: true,
    strokeDashArray: 0,
    yaxis: {
      lines: { show: true }
    },
    xaxis: {
      lines: { show: false }
    },
    padding: {
      bottom: 8
    }
  },

  dataLabels: { enabled: false },

  tooltip: {
    y: { title: { formatter: () => "" } },
    x: {
      formatter: (_, opts) => props.data[opts.dataPointIndex].name
    }
  },

  legend: { show: false }
}));
</script>

<style lang="scss">

</style>
