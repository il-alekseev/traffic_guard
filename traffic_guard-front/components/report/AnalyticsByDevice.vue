<template>
  <BaseLayout
    :title="'Аналитика'"
    :coloredTitle="`по ${deviceName}`"
  >
    <template #data>
      <div class="analytics-data">
        <BaseCard
          v-if="props.traffic.count > 0"
          :title="'Объём трафика'"
          :legend="trafficLegend"
          class="traffic-card"
        >
          <template #DashboardStatistic>
            <TrafficSplineChart :traffic-data="props.traffic" :graph-height="'130px'" />
          </template>
        </BaseCard>
        <div class="request-types-graphs">
          <BaseCard
            v-if="props.requestsTraffic.allowed"
            :title="'Разрешенные запросы'"
            :extra-titles="extraTitles.allowed"
            class="requests-card"
          >
            <template #DashboardStatistic>
              <RequestGraph 
                :data="props.requestsTraffic.allowed" 
                :label="'Разрешенные запросы'" 
                :graph-height="'120px'"
                :show-x="true"
                :show-y="true"
              />
            </template>
          </BaseCard>
          <BaseCard
            v-if="props.requestsTraffic.blocked"
            :title="'Запрещенные запросы'"
            :extra-titles="extraTitles.blocked"
            class="requests-card"
          >
            <template #DashboardStatistic>
              <RequestGraph 
                :data="props.requestsTraffic.blocked" 
                :label="'Запрещенные запросы'" 
                :graph-height="'120px'"
                :show-x="true"
                :show-y="true"
              />
            </template>
          </BaseCard>
          <BaseCard
            v-if="props.requestsTraffic.pending"
            :title="'Запросы в ожидании'"
            :extra-titles="extraTitles.pending"
            class="requests-card"
          >
            <template #DashboardStatistic>
              <RequestGraph 
                :data="props.requestsTraffic.pending" 
                :label="'Запросы в ожидании'" 
                :graph-height="'120px'"
                :show-x="true"
                :show-y="true"
              />
            </template>
          </BaseCard>
        </div>
        <div class="anomalies-and-blocks">
          <div class="anomalies analytics__data-blocks">
            <div class="anomalies-title analytics__data-blocks-title">
              Аномалии
            </div>
            <div class="analytics__data-values-container">
              <span class="anomalies-value analytics__data-value">{{ props.anomaliesAndBlockStat.anomalies }}/{{ props.anomaliesAndBlockStat.all }}</span>
              <span class="anomalies-percent analytics__data-percent">{{ calculatePercent(props.anomaliesAndBlockStat.anomalies, props.anomaliesAndBlockStat.all) }}</span>
            </div>
          </div>
          <div class="blocks analytics__data-blocks">
            <div class="blocks-title analytics__data-blocks-title">
              Блокировки
            </div>
            <div class="analytics__data-values-container">
              <span class="blocks-value analytics__data-value">{{ props.anomaliesAndBlockStat.anomalies }}/{{ props.anomaliesAndBlockStat.all }}</span>
              <span class="blocks-percent analytics__data-percent">{{ calculatePercent(props.anomaliesAndBlockStat.anomalies, props.anomaliesAndBlockStat.all) }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </BaseLayout>
</template>

<script setup lang="ts">
import BaseLayout from '~/components/report/BaseLayout.vue';
import BaseCard from '~/components/report/BaseCard.vue';
import TrafficSplineChart from '~/components/dashboard-grid/TrafficSplineChart.vue';
import RequestGraph from '~/components/dashboard-grid/RequestGraph.vue';
import type { DashboardTraffic } from '~/types/dashboard';
import type { AnomaliesAndBlocksStat, GraphLegend, RequestsAnalyticsTraffic, RequestsGraphExtraTitle } from '~/types/reports';
import { calculatePercent } from '~/helpers';


const props = defineProps<{
  traffic: DashboardTraffic
  requestsTraffic: RequestsAnalyticsTraffic
  anomaliesAndBlockStat: AnomaliesAndBlocksStat
  deviceName: string
}>();


const trafficLegend: GraphLegend = [
  {
    name: 'Входящий',
    color: '#37C84F'
  },
  {
    name: 'Исходящий',
    color: '#2563EB'
  }
];

const extraTitles = computed<Record<string, RequestsGraphExtraTitle>>(() => {
  const calc = (arr: number[]) => {
    const last = arr[arr.length - 1] ?? 0;
    const prev = arr[arr.length - 2] ?? 0;

    let percent = prev === 0
      ? 0 
      : ((last - prev) / prev) * 100;

    percent = Math.round(percent);

    const sign: 'positive' | 'negative' | 'neutral' =
      percent > 0 ? 'positive'
      : percent < 0 ? 'negative'
      : 'neutral';

    return { last, percent, sign };
  };

  const allowed = calc(props.requestsTraffic.allowed.data);
  const blocked = calc(props.requestsTraffic.blocked.data);
  const pending = calc(props.requestsTraffic.pending.data);

  return {
    allowed: {
      mainValue: allowed.last.toString(),
      percentValue: allowed.percent.toString() + '%',
      dynamicSign: allowed.sign,
    },
    blocked: {
      mainValue: blocked.last.toString(),
      percentValue: blocked.percent.toString() + '%' ,
      dynamicSign: blocked.sign,
    },
    pending: {
      mainValue: pending.last.toString(),
      percentValue: pending.percent.toString() + '%',
      dynamicSign: pending.sign,
    },
  };
});

</script>

<style scoped lang="scss">
.analytics-data {
  display: flex;
  flex-direction: column;
}

.request-types-graphs {
  display: flex;
  gap: 12px;
}

.requests-card {
  min-width: 32.5%;
}

.traffic-card {
  width: 100%;
}

.anomalies-and-blocks {
  margin-top: 15px;
  width: 100%;
  display: flex;
  justify-content: space-between;

}

.analytics__data-blocks {
  display: flex;
  flex-direction: column;
  width: 48%;
  max-width: 50%;
}

.analytics__data-values-container {
  padding-top: 6px;
  text-align: center;
  justify-content: center;
  align-items: center;
  display: flex;
  gap: 6px
}

.analytics__data-blocks-title {
  width: 100%;
  font-weight: 400;
  font-size: 15px;
  line-height: 14px;
  color: #1D4ED8;
  padding-bottom: 6px;
  border-bottom: 1px solid #1D4ED8;
}

.analytics__data-value {
  font-weight: 700;
  font-size: 48px;
  line-height: 80px;
  color: #3F3F46;
}

.analytics__data-percent {
  font-weight: 700;
  font-size: 48px;
  line-height: 80px;
  color: #A1A1AA;
}

</style>
