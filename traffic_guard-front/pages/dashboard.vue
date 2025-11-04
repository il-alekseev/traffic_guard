<template>
  <div class="dashboard-page">
    <div class="dashboard__header">
      <div class="dashboard__header-block">
        <div class="dashboard__header-title-block">
          <h1 class="dashboard__title page-title">Нет новых выявлений</h1>
          <div class="dashboard__subtitle-block page-subtitle-block">
            <p class="dashboard__subtitle page-subtitle">Обновлено {{ minutesAgo }} мин назад</p>
            <div class="dashboard__refresh-btn page-refresh-btn" @click="fetchData">
              <ReloadIcon />
            </div>
          </div>
        </div>
        <div class="dashboard__filters">
          <div class="dashboard__datepicker-container">
            <DatePicker v-model="dateRange" />
          </div>
          <button class="dashboard__filter-button" @click="showFilters">
            <div class="dashboard__filter-button-icon">
              <FilterIcon />
            </div>
            <span class="dashboard__filter-button-text">Фильтр</span>
          </button>
          <DownloadButton />
        </div>
      </div>
    </div>

    <div class="dashboard-data">
      <div class="dashboard-grid">
        <!-- <div class="dashboard__trafic">
          DASHBOARD TRAFIC DATA
        </div>
        <div class="dashboard__requests">
          <div class="dashboard__request dashboard__request-allowed">
            DASHBOARD REQUEST ALLOW
          </div>
          <div class="dashboard__request dashboard__request-before-blocked">
            DASHBOARD REQUEST BEFORE BLOCKED
          </div>
          <div class="dashboard__request dashboard__request-blocked">
            DASHBOARD REQUEST BLOCKED
          </div>
          <div class="dashboard__request dashboard__request-waiting">
            DASHBOARD REQUEST WAITING
          </div>
        </div> -->
        <DashboardCard class="dashboard__trafic" title="Трафик" :legend="{input: {name: 'Входящий', color: '#37C84F'}, output: {name: 'Исходящий', color: '#2563EB'}}">
          <template v-if="loading.trafic" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.trafic !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.trafic }}</p>
          </template>
          <template v-else-if="dashboardData.trafic" #DashboardStatistic>
            <DashboardTrafficChart :trafficData="dashboardData.trafic!"/>
          </template>
          <template v-else-if="!dashboardData.trafic" #EmptyData>
            <p class="empty-data">Данные отсутствуют</p>
          </template>
        </DashboardCard>
        <DashboardCard class="dashboard__top-categories" title="Рейтинг запрещенных категорий" link="/report">
          <template v-if="loading.topCategories" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.topCategories !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.topCategories }}</p>
          </template>
          <template v-else-if="dashboardData.topCategories && dashboardData.topCategories.length > 0" #DashboardStatistic>
            <DashboardTopCategories :categories="dashboardData.topCategories!" />
          </template>
          <template v-else-if="!dashboardData.topCategories" #EmptyData>
            <p class="empty-data">Данные отсутствуют</p>
          </template>
        </DashboardCard>
        <DashboardCard class="dashboard__top-detections" title="Топ нерешенных выявлений" link="/detections">
          <template v-if="loading.topDetections" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.topDetections !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.topDetections }}</p>
          </template>
          <template v-else-if="dashboardData.topDetections && dashboardData.topDetections.length > 0" #DashboardStatistic>
            <DashboardTopDetections :detections="dashboardData.topDetections!" />
          </template>
          <template v-else-if="!dashboardData.topDetections || dashboardData.topDetections.length == 0" #EmptyData>
            <p class="empty-data">Данные отсутствуют</p>
          </template>
        </DashboardCard>
        <!-- <div class="dashboard__top-detections">
          DASHBOARD TOP DETECTIONS
        </div>
        <div class="dashboard__anomalies">
          DASHBOARD TOP ANOMALIES
        </div>
        <div class="dashboard__proh-activity">
          DASHBOARD TOP PROH ACTIVITY
        </div> -->
      </div>
      <div class="devices-grid">
        <!-- <div class="devices-stat">
          STAT BY DEVICE
        </div> -->
      </div>
    </div>

    <SideModal
      v-model="isShowFilters"
      title="Фильтры"
      @close="closeFilters"
    >
      <FilterForm
        :filtersData="filtersData"
        @close="closeFilters"
        @setFilters="handleSetFilters"
      />

    </SideModal>
  </div>
</template>

<script setup lang="ts">
import {definePageMeta} from '#imports';
import { useRouter } from 'vue-router';
import { useDashboardStore } from '~/stores/dashboard';
import DatePicker from '~/components/ui/DatePicker.vue';
import DownloadButton from '~/components/ui/DownloadButton.vue';
import SideModal from '~/components/ui/SideModal.vue';
import FilterForm, { type DashboardFilter } from '~/components/filters/DashboardFilterForm.vue';
import DashboardCard from "~/components/dashboard-grid/BaseCard.vue"
import DashboardTopCategories from "~/components/dashboard-grid/TopCategories.vue"
import DashboardTopDetections from "~/components/dashboard-grid/TopDetections.vue"
import DashboardTrafficChart from "~/components/dashboard-grid/TrafficSplineChart.vue"
import ReloadIcon from "~/assets/img/reload.svg"
import FilterIcon from "~/assets/img/filter-icon.svg"
import { getCurrentDateWithOffset } from '~/helpers';


const route = useRoute();
const router = useRouter();

definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const dashboardStore = useDashboardStore();

const loading = ref({
  topCategories: true,
  topDetections: true,
  trafic: true,
  requests: false,
  anomalies: false,
  events: false,
  proh_activity: false,
  devicesState: false
});
const isLoading = computed(() => 
  Object.values(loading.value).some(status => status === true)
);

const fetchError = ref({
  topCategories: '',
  topDetections: '',
  trafic: '',
  requests: '',
  anomalies: '',
  events: '',
  proh_activity: '',
  devicesState: ''
});
const hasErrors = computed(() => 
  Object.values(fetchError.value).some(error => error !== '')
);

const dashboardData = computed(() => ({
  topCategories: dashboardStore.topCategories,
  topDetections: dashboardStore.topDetections,
  trafic: dashboardStore.trafic,
  requests: dashboardStore.requests,
  anomalies: dashboardStore.anomalies,
  events: dashboardStore.events,
  proh_activity: dashboardStore.events,
  devicesState: dashboardStore.devicesState
}));

const lastUpdated = ref<Date | null>(null);
const minutesAgo = ref(0);
let intervalId: ReturnType<typeof setInterval> | null = null;

const updateMinutesAgo = () => {
  if (!lastUpdated.value) return;
  const diffMs = Date.now() - lastUpdated.value.getTime();
  minutesAgo.value = Math.floor(diffMs / 60000);
}

const fetchData = async () => {
  Object.keys(fetchError.value).forEach(key => {
    fetchError.value[key as keyof typeof fetchError.value] = '';
  });

  const from = dateRange.value.from?.toISOString();
  const to = dateRange.value.to?.toISOString();
  const TOPS_COUNT = 5;
  const TRAFFIC_COUNT = 20;
  const hostname = deviceFilter.value;

  try {
    await Promise.all([
      dashboardStore.fetchTopCategories(from, to, TOPS_COUNT, hostname)
        .then(() => {
          loading.value.topCategories = false;
        })
        .catch((error: Error) => {
          loading.value.topCategories = false;
          fetchError.value.topCategories = error.message;
        }),

      dashboardStore.fetchDetections(from, to, TOPS_COUNT, hostname)
        .then(() => {
          loading.value.topDetections = false;
        })
        .catch((error: Error) => {
          loading.value.topDetections = false;
          fetchError.value.topDetections = error.message;
        }),

      dashboardStore.fetchTraffic(from, to, TRAFFIC_COUNT, hostname)
        .then(() => {
          loading.value.trafic = false;
        })
        .catch((error: Error) => {
          loading.value.trafic = false;
          fetchError.value.trafic = error.message;
        })
    ]);

  } catch (error) {
    console.error('Ошибка при загрузке данных dashboard:', error);
  }

  lastUpdated.value = new Date();
  updateMinutesAgo();
}

const dateRange = ref<{ from: Date | null; to: Date | null }>({
  from: getCurrentDateWithOffset(-1, 'd'),
  to: getCurrentDateWithOffset()
})

const isShowFilters = ref(false);
const showFilters = () => {
  isShowFilters.value = true;
}
const closeFilters = () => {
  isShowFilters.value = false;
}

const deviceFilter = ref<string | undefined>();
const filtersData = computed<DashboardFilter | null>(() => {
  const device = deviceFilter.value ?? '';

  if (!device) {
    return null;
  }

  return {
    device: { id: device, name: device },
  };
});

const initFiltersFromUrl = () => {
  const query = route.query;

  deviceFilter.value = query.device != null ? String(query.device) : undefined;
};

const updateUrlParams = () => {
  const query: Record<string, string | number> = {};

  if (deviceFilter.value && deviceFilter.value !== '') query.device = deviceFilter.value;

  router.replace({ query });
};

const handleSetFilters = (filtersData: DashboardFilter) => {
  filtersData.device.id !== '' ? deviceFilter.value = filtersData.device.id : deviceFilter.value = undefined;

  updateUrlParams()
}

onMounted(async () => {
  initFiltersFromUrl();
  await fetchData();
  intervalId = setInterval(updateMinutesAgo, 60 * 1000);
});
onUnmounted(() => {
  if (intervalId) clearInterval(intervalId);
})
watch(
  () => route.query,
  async (newQuery, oldQuery) => {
    const withoutCreate = (q: typeof newQuery) => {
      const { create, ...rest } = q;
      return JSON.stringify(rest);
    };

    if (withoutCreate(newQuery) === withoutCreate(oldQuery)) {
      return;
    }

    initFiltersFromUrl();
    await fetchData();
  },
  { deep: true }
);


</script>

<style lang="scss" scoped>

.dashboard__header-block {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dashboard__header-title-block {
  display: flex;
  gap: 0.75rem;
}

.dashboard__subtitle-block {
  align-self: flex-end;
}

.dashboard__loading,
.dashboard__empty {
  text-align: center;
  padding: 2rem;
  color: #3F3F46;
}

.dashboard-data {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  overflow-x: auto;
  min-height: calc(100vh - 21rem);
}

.dashboard-grid {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}


.dashboard__filter-button {
  background-color: #FFFFFF;
  box-shadow: 0px 1px 2px 0px #0000000D;
  padding: 0.5rem 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
}

.dashboard__filter-button-icon {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 1.25rem;
  height: 1.25rem;
  color: #A1A1AA;
}

.dashboard__filter-button-icon svg {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.dashboard__filter-button-text {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
}

.dashboard__filter-selected-area {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.dashboard__filters {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.dashboard__download-btn {
  padding: 0.5rem;
  background-color: #FFFFFF;
  box-shadow: 0px 1px 2px 0px #0000000D;
  border-radius: 6px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dashboard__download-btn svg {
  width: 20px;
  height: 20px;
  object-fit: contain;
}

.loading-data, .empty-data {
  font-size: 1rem;
  line-height: 1rem;
  color: #3F3F46;
  min-height: 148px;
}
.error-data {
  font-size: 1rem;
  line-height: 1rem;
  color: #FB2904;
}

</style>
