<template>
  <BaseLayout
    :title="'Отчёт'"
    :coloredTitle="'по основной активности'"
  >
    <template #data>
      <div class="activity-data">
        <BaseCard
          :title="'Рейтинг категорий'"
          :legend="ratingLegend"
          class="category-card"
        >
          <template #DashboardStatistic>
            <BarChartForRating :items="categoiersChartData" />
          </template>
        </BaseCard>
        <BaseCard
          :title="'Рейтинг ресурсов'"
          :legend="ratingLegend"
          class="resources-card"
        >
          <template #DashboardStatistic>
            <BarChartForRating :items="resourcesChartData" showSubcategories />
          </template>
        </BaseCard>
        <BaseCard
          :title="'Объём трафика'"
          :legend="trafficLegend"
          class="traffic-card"
        >
          <template #DashboardStatistic>
            <TrafficSplineChart :traffic-data="props.traffic" :graph-height="'130px'" />
          </template>
        </BaseCard>
      </div>
    </template>
  </BaseLayout>
</template>

<script setup lang="ts">
import BaseLayout from '~/components/report/BaseLayout.vue';
import BaseCard from '~/components/report/BaseCard.vue';
import TrafficSplineChart from '~/components/dashboard-grid/TrafficSplineChart.vue';
import BarChartForRating, { type ChartItem } from '~/components/report/BarChartForRating.vue';
import type { DashboardTraffic } from '~/types/dashboard';
import type { GraphLegend, ReportCategoryItem, ReportResourceItem } from '~/types/reports';


const props = defineProps<{
  traffic: DashboardTraffic;
  categoriesTop: ReportCategoryItem[];
  resourcesTop: ReportResourceItem[];
}>();

const categoiersChartData = computed<ChartItem[]>(() => {
  return props.categoriesTop.map((c) => {
    return {
      name: c.category,
      before_block: c.stat.before_block,
      after_block: c.stat.after_block,
      waiting: c.stat.pending
    }
  })
})

const resourcesChartData = computed<ChartItem[]>(() => {
  return props.resourcesTop.map((r) => {
    return {
      name: r.resource,
      before_block: r.stat.before_block,
      after_block: r.stat.after_block,
      waiting: r.stat.pending,
      categoies: r.categories
    }
  })
})

const ratingLegend: GraphLegend = [
  {
    name: 'До блокировки',
    color: '#37C84F'
  },
  {
    name: 'После блокировки',
    color: '#FB2904'
  },
  {
    name: 'В ожидании',
    color: '#EFB100'
  },
];

const trafficLegend: GraphLegend = [
  {
    name: 'Входящий',
    color: '#37C84F'
  },
  {
    name: 'Исходящий',
    color: '#2563EB'
  }
]
</script>

<style scoped lang="scss">
.activity-data {
  display: flex;
  row-gap: 15px;
  justify-content: space-between;
  flex-wrap: wrap;
}

.traffic-card {
  width: 100%;
}

</style>
