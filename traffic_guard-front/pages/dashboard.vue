<template>
  <div class="dashboard-page">
    <div class="dashboard__header">
      <div class="dashboard__header-block">
        <div class="dashboard__header-title-block">
          <h1 class="dashboard__title page-title">{{ dashboardTitle }}</h1>
          <div class="dashboard__subtitle-block page-subtitle-block">
            <p class="dashboard__subtitle page-subtitle">Обновлено {{ minutesAgo }} мин назад</p>
            <div class="dashboard__refresh-btn page-refresh-btn" @click="fetchDashboardData">
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
        <DashboardCard class="dashboard__trafic" title="Трафик" :legend="{input: {name: 'Входящий', color: '#37C84F'}, output: {name: 'Исходящий', color: '#2563EB'}}">
          <template v-if="loading.trafic" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.trafic !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.trafic }}</p>
          </template>
          <template v-else-if="dashboardData.trafic" #DashboardStatistic>
            <DashboardTrafficChart :trafficData="dashboardData.trafic!" :graphHeight="dashboardGraphHeight.DashboardTrafficChart"/>
          </template>
          <template v-else-if="!dashboardData.trafic" #EmptyData>
            <p class="empty-data">Данные отсутствуют</p>
          </template>
        </DashboardCard>
        <div class="dashboard__requests">
          <DashboardRequestCard
            class="dashboard__request"
            title="Разрешенные запросы"
            :data="{
              currentValue: dashboardData?.requests?.allowed?.data?.data?.at(-1) ?? 0,
              oldValue: dashboardData?.requests?.allowed?.data?.data?.at(-2) ?? 0
            }"
          >
            <template v-if="loading.requests" #LoadingData>
              <p class="loading-data">Загрузка...</p>
            </template>
            <template v-else-if="fetchError.requests !== ''" #ErrorData>
              <p class="error-data">{{ fetchError.requests }}</p>
            </template>
            <template v-else-if="dashboardData.requests && dashboardData.requests.allowed" #DashboardStatistic>
              <RequestGraph :data="dashboardData.requests.allowed!" :label="'Разрешенные запросы'" :graphHeight="dashboardGraphHeight.RequestGraph"/>
            </template>
            <template v-else-if="!dashboardData.requests.allowed || dashboardData.requests.allowed.data.data.length == 0" #EmptyData>
              <p class="empty-data">Данные отсутствуют</p>
            </template>
          </DashboardRequestCard>
          <DashboardRequestCard
            class="dashboard__request"
            title="Запросы до блокировки"
            :data="{
              currentValue: dashboardData?.requests?.before_block?.data?.data?.at(-1) ?? 0,
              oldValue: dashboardData?.requests?.before_block?.data?.data?.at(-2) ?? 0
            }"  
          >
            <template v-if="loading.requests" #LoadingData>
              <p class="loading-data">Загрузка...</p>
            </template>
            <template v-else-if="fetchError.requests !== ''" #ErrorData>
              <p class="error-data">{{ fetchError.requests }}</p>
            </template>
            <template v-else-if="dashboardData.requests && dashboardData.requests.before_block" #DashboardStatistic>
              <RequestGraph :data="dashboardData.requests.before_block!" :label="'Запросы до блокировки'" :graphHeight="dashboardGraphHeight.RequestGraph"/>
            </template>
            <template v-else-if="!dashboardData.requests.before_block || dashboardData.requests.before_block.data.data.length == 0" #EmptyData>
              <p class="empty-data">Данные отсутствуют</p>
            </template>
          </DashboardRequestCard>
          <DashboardRequestCard
            class="dashboard__request"
            title="Заблокированные запросы"
            :data="{
              currentValue: dashboardData?.requests?.blocked?.data?.data?.at(-1) ?? 0,
              oldValue: dashboardData?.requests?.blocked?.data?.data?.at(-2) ?? 0
            }"  
          >
            <template v-if="loading.requests" #LoadingData>
              <p class="loading-data">Загрузка...</p>
            </template>
            <template v-else-if="fetchError.requests !== ''" #ErrorData>
              <p class="error-data">{{ fetchError.requests }}</p>
            </template>
            <template v-else-if="dashboardData.requests && dashboardData.requests.blocked" #DashboardStatistic>
              <RequestGraph :data="dashboardData.requests.blocked!" :label="'Заблокированные запросы'" :graphHeight="dashboardGraphHeight.RequestGraph"/>
            </template>
            <template v-else-if="!dashboardData.requests.blocked || dashboardData.requests.blocked.data.data.length == 0" #EmptyData>
              <p class="empty-data">Данные отсутствуют</p>
            </template>
          </DashboardRequestCard>
          <DashboardRequestCard
            class="dashboard__request"
            title="Запросы в ожидании"
            :data="{
              currentValue: dashboardData?.requests?.pending?.data?.data?.at(-1) ?? 0,
              oldValue: dashboardData?.requests?.pending?.data?.data?.at(-2) ?? 0
            }"  
          >
            <template v-if="loading.requests" #LoadingData>
              <p class="loading-data">Загрузка...</p>
            </template>
            <template v-else-if="fetchError.requests !== ''" #ErrorData>
              <p class="error-data">{{ fetchError.requests }}</p>
            </template>
            <template v-else-if="dashboardData.requests && dashboardData.requests.pending" #DashboardStatistic>
              <RequestGraph :data="dashboardData.requests.pending!" :label="'Запросы в ожидании'" :graphHeight="dashboardGraphHeight.RequestGraph"/>
            </template>
            <template v-else-if="!dashboardData.requests.pending || dashboardData.requests.pending.data.data.length == 0" #EmptyData>
              <p class="empty-data">Данные отсутствуют</p>
            </template>
          </DashboardRequestCard>
        </div>
        <DashboardCard class="dashboard__top-categories" title="Рейтинг запрещенных категорий" link="/reports">
          <template v-if="loading.topCategories" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.topCategories !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.topCategories }}</p>
          </template>
          <template v-else-if="dashboardData.topCategories && dashboardData.topCategories.length > 0" #DashboardStatistic>
            <DashboardTopCategories :categories="dashboardData.topCategories!" />
          </template>
          <template v-else-if="!dashboardData.topCategories || dashboardData.topCategories.length == 0" #EmptyData>
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
        <DashboardCard class="dashboard__anomalies" title="Аномалии" link="/sessions">
          <template v-if="loading.anomalies" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.anomalies !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.anomalies }}</p>
          </template>
          <template v-else-if="dashboardData.anomalies && dashboardData.anomalies.host_anomalies.length > 0" #DashboardStatistic>
            <DashboardAnomalies :anomaliesData="dashboardData.anomalies!" :graphHeight="dashboardGraphHeight.DashboardAnomalies"/>
          </template>
          <template v-else-if="!dashboardData.topDetections || dashboardData.topDetections.length == 0" #EmptyData>
            <p class="empty-data">Данные отсутствуют</p>
          </template>
        </DashboardCard>
        <DashboardCard class="dashboard__logs" title="События" link="/logs">
          <template v-if="loading.logs" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.logs !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.logs }}</p>
          </template>
          <template v-else-if="dashboardData.logs && dashboardData.logs.length > 0" #DashboardStatistic>
            <ShortLogs :data="dashboardData.logs" />
          </template>
          <template v-else-if="!dashboardData.logs || dashboardData.logs.length == 0" #EmptyData>
            <p class="empty-data">Данные отсутствуют</p>
          </template>
        </DashboardCard>
        <DashboardCard class="dashboard__proh-activity" title="График запрещенной активности">
          <template v-if="loading.proh_activity" #LoadingData>
            <p class="loading-data">Загрузка...</p>
          </template>
          <template v-else-if="fetchError.proh_activity !== ''" #ErrorData>
            <p class="error-data">{{ fetchError.proh_activity }}</p>
          </template>
          <template v-else-if="dashboardData.proh_activity" #DashboardStatistic>
            <HeatmapGraph
              :values="dashboardData.proh_activity.data"
              :height="dashboardGraphHeight.HeatmapGraph" 
              :days="['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']" 
              :months="['Янв', 'Фев', 'Мар', 'Апр',  'Май', 'Июн',  'Июл', 'Авг',  'Сен', 'Окт',  'Дек']" 
            />
          </template>
          <template v-else-if="!dashboardData.logs || dashboardData.logs.length == 0" #EmptyData>
            <p class="empty-data">Данные отсутствуют</p>
          </template>
        </DashboardCard>
      </div>
    </div>

    <template v-if="devices && devices?.length > 0">
      <div class="dashboard-nodes__header">
        <div class="dashboard-nodes__header-block">
          <div class="dashboard-nodes__header-title-block">
            <h2 class="dashboard-nodes__title page-title">Состояние узлов</h2>
            <div class="dashboard-nodes__subtitle-block page-subtitle-block">
              <p class="dashboard__subtitle page-subtitle">Обновлено {{ devicesMinutesAgo }} мин назад</p>
              <div class="dashboard__refresh-btn page-refresh-btn" @click="fetchDevicesData">
                <ReloadIcon />
              </div>
            </div>
          </div>
          <div class="dashboard-nodes__extra">
            <DownloadButton />
          </div>
        </div>
      </div>

      <div class="dashboard-nodes-data" v-if="devicesRequestsData.length > 0">
        <div class="dashboard-grid">
          <DashboardCard
            v-for="deviceRequestData in devicesRequestsData"
            class="dashboard-node__trafic" 
            :title="deviceRequestData.name"
            :titleClass="'ngfw-bagde ngfw-bagde_small'"
            :titleStyle="{
              backgroundColor: generateColor(deviceRequestData.name).background,
              color: generateColor(deviceRequestData.name).color
            }"
            :legend="{input: {name: 'Заблокированные ресурсы', color: '#FB2904'}, output: {name: ' Ожидающие решения', color: '#EFB100'}}"
            legend-pos="bottom-right"
          >
            <template #titleStatus>
              <div class="device__status" :class="deviceRequestData.data.time.length > 0 ? 'active' : 'inactive'">
                <div class="device__status-dot"></div>
                <div class="device__status-text">{{ deviceRequestData.data.time.length > 0 ? 'Активен' : 'Не активен' }}</div>
              </div>
            </template>
            <template v-if="devicesLoading[deviceRequestData.name + '_blocked'] || devicesLoading[deviceRequestData.name + '_pending']" #LoadingData>
              <p class="loading-data">Загрузка...</p>
            </template>
            <template v-else-if="devicesFetchErrors[deviceRequestData.name + '_blocked'] !== '' ||  devicesFetchErrors[deviceRequestData.name + '_pending'] !== ''" #ErrorData>
              <p class="error-data">Ошибка при получении данных</p>
            </template>
            <template v-else-if="deviceRequestData.data.blocked.length > 0 || deviceRequestData.data.pending.length > 0" #DashboardStatistic>
              <RequestsSplineChart :data="deviceRequestData" :graphHeight="dashboardGraphHeight.RequestsSplineChart"/>
            </template>
            <template v-else-if="!dashboardData.trafic" #EmptyData>
              <p class="empty-data">Данные отсутствуют</p>
            </template>
          </DashboardCard>
        </div>
      </div>
    </template>

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
import { useDetectionsStore } from '~/stores/detections';
import DatePicker from '~/components/ui/DatePicker.vue';
import DownloadButton from '~/components/ui/DownloadButton.vue';
import SideModal from '~/components/ui/SideModal.vue';
import FilterForm, { type DashboardFilter } from '~/components/filters/DashboardFilterForm.vue';
import DashboardCard from "~/components/dashboard-grid/BaseCard.vue"
import DashboardRequestCard from "~/components/dashboard-grid/BaseRequestCard.vue"
import RequestGraph from '~/components/dashboard-grid/RequestGraph.vue';
import DashboardTopCategories from "~/components/dashboard-grid/TopCategories.vue"
import DashboardTopDetections from "~/components/dashboard-grid/TopDetections.vue"
import DashboardTrafficChart from "~/components/dashboard-grid/TrafficSplineChart.vue"
import DashboardAnomalies from '~/components/dashboard-grid/DashboardAnomalies.vue';
import ReloadIcon from "~/assets/img/reload.svg"
import FilterIcon from "~/assets/img/filter-icon.svg"
import { getCurrentDateWithOffset, isValidDateString, remToPx } from '~/helpers';
import type { DashboardDeviceRequestsSplineChart, DashboardRequestsType } from '~/types/dashboard';
import RequestsSplineChart from '~/components/dashboard-grid/devices/RequestsSplineChart.vue';
import { useDeviceColors } from '~/composables/useDeviceColors';
import ShortLogs from '~/components/dashboard-grid/ShortLogs.vue';
import HeatmapGraph from '~/components/dashboard-grid/HeatmapGraph.vue';


const route = useRoute();
const router = useRouter();

definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const { generateColor } = useDeviceColors();

const TOPS_COUNT = 5;
const TRAFFIC_COUNT = 20;
const REQUESTS_COUNT = 10;
const REQUESTS_BY_DEVICE_COUNT = 5;
const LOGS_PAGE = 1;
const LOGS_LIMIT = 20;
const YEAR = new Date().getFullYear().toString();

const dashboardStore = useDashboardStore();
const logsStore = useLogsStore();
const detectionsStore = useDetectionsStore();
const deviceStore = useDevicesStore();

const loading = ref({
  detectionsStat: true,
  topCategories: true,
  topDetections: true,
  trafic: true,
  requests: true,
  anomalies: false,
  logs: false,
  proh_activity: false,
  
});
const isLoading = computed(() => 
  Object.values(loading.value).some(status => status === true)
);

const fetchError = ref({
  detectionsStat: '',
  topCategories: '',
  topDetections: '',
  trafic: '',
  requests: '',
  anomalies: '',
  logs: '',
  proh_activity: '',
});
const hasErrors = computed(() => 
  Object.values(fetchError.value).some(error => error !== '')
);

const dashboardData = computed(() => ({
  detectionsStat: detectionsStore.detectionsStat,
  topCategories: dashboardStore.topCategories,
  topDetections: dashboardStore.topDetections,
  trafic: dashboardStore.trafic,
  requests: dashboardStore.requests,
  anomalies: dashboardStore.anomalies,
  logs: dashboardStore.logs,
  proh_activity: dashboardStore.proh_activity,
}));

const rootFontSize = useRootFontSize();

const dashboardGraphHeight = computed(() => ({
  DashboardAnomalies: rootFontSize.value * 16.25 + 'px',
  HeatmapGraph: rootFontSize.value * 11.5 + 'px',
  RequestGraph: rootFontSize.value * 5.625 + 'px',
  DashboardTrafficChart: rootFontSize.value * 8.25 + 'px',
  RequestsSplineChart: rootFontSize.value * 8.25 + 'px',
}));

const dashboardTitle = computed<string>(() => {
  const detectionCount = dashboardData.value.detectionsStat?.detected;

  if (!detectionCount) return 'Нет новых выявлений';
  
  const pluralForm = (() => {
    const n = detectionCount % 100;
    if (n === 1) return 'новое выявление';
    if (n >= 2 && n <= 4) return 'новых выявления';
    return 'новых выявлений';
  })();

  return `У вас ${detectionCount} ${pluralForm}`;
});

const lastUpdated = ref<Date | null>(null);
const minutesAgo = ref(0);
let intervalId: ReturnType<typeof setInterval> | null = null;

const updateMinutesAgo = () => {
  if (!lastUpdated.value) return;
  const diffMs = Date.now() - lastUpdated.value.getTime();
  minutesAgo.value = Math.floor(diffMs / 60000);
}

const fetchDashboardData = async () => {
  Object.keys(fetchError.value).forEach(key => {
    fetchError.value[key as keyof typeof fetchError.value] = '';
  });

  const from = dateRange.value.from?.toISOString();
  const to = dateRange.value.to?.toISOString();
  const hostname = deviceFilter.value;

  try {
    await Promise.all([
      detectionsStore.fetchDetectionStats(from, to)
        .then(() => {
          loading.value.detectionsStat = false;
        })
        .catch((error: Error) => {
          loading.value.detectionsStat = false;
          fetchError.value.detectionsStat = error.message;
        }),

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
        }),

      dashboardStore.fetchRequests(from, to, REQUESTS_COUNT, hostname)
        .then(() => {
          loading.value.requests = false;
        })
        .catch((error: Error) => {
          loading.value.requests = false;
          fetchError.value.requests = error.message;
        }),
      
      dashboardStore.fetchAnomalies(from, to, hostname)
        .then(() => {
          loading.value.anomalies = false;
        })
        .catch((error: Error) => {
          loading.value.anomalies = false;
          fetchError.value.anomalies = error.message;
        }),

      logsStore.fetchLogs(LOGS_PAGE, LOGS_LIMIT)
        .then((res) => {
          dashboardStore.logs = res.data;
          loading.value.logs = false;
        })
        .catch((error: Error) => {
          loading.value.logs = false;
          fetchError.value.logs = error.message;
        }),

      dashboardStore.fetchProhActivity(YEAR, hostname)
        .then(() => {
          loading.value.proh_activity = false;
        })
        .catch((error: Error) => {
          loading.value.proh_activity = false;
          fetchError.value.proh_activity = error.message;
        }),
    ]);

  } catch (error) {
    console.error('Ошибка при загрузке данных dashboard:', error);
  }

  lastUpdated.value = new Date();
  updateMinutesAgo();
}

const devices = computed(() => deviceStore.devices);
const devicesLoading = ref<Record<string, boolean>>({});
const devicesFetchErrors = ref<Record<string, string>>({});
const devicesRequestsData = ref<DashboardDeviceRequestsSplineChart[]>([]);

const initDeviceLoading = () => {
  if (!devices.value) return;

  devicesLoading.value = {};
  devices.value.forEach(device => {
    devicesLoading.value[device + '_blocked'] = false;
    devicesLoading.value[device + '_pending'] = false;
  });
};

const initDevicesFetchErrors = () => {
  if (!devices.value) return;

  devicesFetchErrors.value = {};
  devices.value.forEach(device => {
    devicesFetchErrors.value[device + '_blocked'] = '';
    devicesFetchErrors.value[device + '_pending'] = '';
  });
};

const initDevicesGraphData = () => {
  if (!devices.value) return;

  devicesRequestsData.value = [];
  devices.value.forEach(device => {
    devicesRequestsData.value.push({
      name: device,
      data: {
        blocked: [],
        pending: [],
        time: []
      }
    })
  });
}

const prepareDevicesData = () => {
  initDeviceLoading();
  initDevicesFetchErrors();
  initDevicesGraphData();
}

const lastDeviceUpdated = ref<Date | null>(null);
const devicesMinutesAgo = ref(0);
let devicesIntervalId: ReturnType<typeof setInterval> | null = null;

const updateDevicesMinutesAgo = () => {
  if (!lastDeviceUpdated.value) return;
  const diffMs = Date.now() - lastDeviceUpdated.value.getTime();
  devicesMinutesAgo.value = Math.floor(diffMs / 60000);
}

const fetchDeviceData = async (deviceName: string, type: 'blocked' | 'pending') => {
  if (devicesLoading.value[deviceName + '_' + type]) return;

  devicesLoading.value[deviceName + '_' + type] = true;
  devicesFetchErrors.value[deviceName  + '_' + type] = '';

  const from = dateRange.value.from?.toISOString();
  const to = dateRange.value.to?.toISOString();

  try {
    const result = await dashboardStore.fetchRequest(from, to, type, REQUESTS_BY_DEVICE_COUNT, deviceName, false);

    if (result.data) {
      const deviceObj = devicesRequestsData.value.find((dev) => dev.name === deviceName);
      if (!deviceObj) return;

      if (result.type === 'requests_blocked') {
        deviceObj.data.blocked = result.data.data;
        deviceObj.data.time = result.data.time;
      } else if (result.type === 'requests_pending') {
        deviceObj.data.pending = result.data.data;
        deviceObj.data.time = result.data.time;
      }
    }
  } catch (error: any) {
    console.error(`Ошибка при загрузке данных  по девайсу: ${deviceName}`, error);
    devicesFetchErrors.value[deviceName  + '_' + type] = error?.message || 'Неизвестная ошибка';
  } finally {
    devicesLoading.value[deviceName + '_' + type] = false;
  }
}

const fetchDevicesData = async () => {
  try {
    if (!devices.value || devices.value.length === 0) {
      return;
    }
    
    prepareDevicesData();
    
    const promises: Promise<void>[] = [];
    
    devices.value.forEach(device => {
      promises.push(fetchDeviceData(device, 'blocked'));
      promises.push(fetchDeviceData(device, 'pending'));
    });
    
    await Promise.all(promises);
  } catch (error) {
    console.error('Ошибка при загрузке данных устройств:', error);
  }

  lastDeviceUpdated.value = new Date();
  updateDevicesMinutesAgo();
}

const fetchAllData = async () => {
  try {
     await fetchDashboardData();
     await fetchDevicesData();
  } catch (error) {
    console.error('Ошибка при загрузке данных:', error);
  }
}

const isShowFilters = ref(false);
const showFilters = () => {
  isShowFilters.value = true;
}
const closeFilters = () => {
  isShowFilters.value = false;
}

const dateRange = ref<{ from: Date | null; to: Date | null }>({
  from: getCurrentDateWithOffset(-1, 'd'),
  to: getCurrentDateWithOffset()
})
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

  const fromStr = typeof query.from === 'string' ? query.from : null;
  const toStr = typeof query.to === 'string' ? query.to : null;
  
  dateRange.value.from = isValidDateString(fromStr)
    ? new Date(fromStr!)
    : getCurrentDateWithOffset(-1, 'd');

  dateRange.value.to = isValidDateString(toStr)
    ? new Date(toStr!)
    : getCurrentDateWithOffset();
};

const updateUrlParams = () => {
  const query: Record<string, string | number> = {};

  if (deviceFilter.value && deviceFilter.value !== '') query.device = deviceFilter.value;
  if (dateRange.value.from) query.from = dateRange.value.from.toISOString();
  if (dateRange.value.to) query.to = dateRange.value.to.toISOString();

  router.replace({ query });
};

const handleSetFilters = (filtersData?: DashboardFilter) => {
  
  if (filtersData) {
    filtersData.device.id !== '' ? deviceFilter.value = filtersData.device.id : deviceFilter.value = undefined;
  }

  updateUrlParams()
}

onMounted(async () => {
  initFiltersFromUrl();
  await fetchAllData();
  intervalId = setInterval(updateMinutesAgo, 60 * 1000);
  devicesIntervalId = setInterval(updateDevicesMinutesAgo, 60 * 1000);
});

onUnmounted(() => {
  if (intervalId) clearInterval(intervalId);
  if (devicesIntervalId) clearInterval(devicesIntervalId);
})

watch(dateRange, () => {
  handleSetFilters();
})

watch(
  () => route.query,
  async () => {
    initFiltersFromUrl();
    fetchAllData();
  },
  { deep: true }
);

</script>

<style lang="scss" scoped>

.dashboard__header-block, .dashboard-nodes__header-block {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dashboard__header-title-block, .dashboard-nodes__header-title-block {
  display: flex;
  gap: 0.75rem;
}

.dashboard__subtitle-block, .dashboard-nodes__subtitle-block {
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
  margin-bottom: 1.5rem;
}

.dashboard-nodes-data {
  margin-top: 1.5rem;
  display: flex;
}

.dashboard-grid {
  width: 100%;
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  overflow: hidden;
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
  width: 1.25rem;
  height: 1.25rem;
  object-fit: contain;
}

.loading-data, .empty-data {
  font-size: 1rem;
  line-height: 1rem;
  color: #3F3F46;
  min-height: 9.25rem;
}
.error-data {
  font-size: 1rem;
  line-height: 1rem;
  color: #FB2904;
}

.dashboard__trafic {
  min-height: 14.25rem;
}

.dashboard__proh-activity {
  max-width: 65rem;
  width: 100%;
}

.dashboard__requests {
  display: flex;
  row-gap: 0.75rem;
  column-gap: 1rem;
  flex-wrap: wrap;
  max-width: 65rem;
}

.device__status {
  padding: 0.25rem 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  border-radius: 6px;
}

.device__status.active {
  background-color: #DCFCE7;
  color: #008236;
}
.device__status.inactive {
  background-color: #FFE2E2;
  color: #C10007;
}
.device__status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.device__status.active .device__status-dot {
  background-color: #00C950;
}

.device__status.inactive .device__status-dot {
  background-color: #FB2C36;
}

.device__status-text {
  font-weight: 500;
  font-size: 0.75rem;
  line-height: 1rem;
}

.dashboard__logs {
  padding-bottom: 0;
}

@media screen and (min-width: 1921px) {
  .dashboard-grid {
    max-width: 130rem;
  }

  .dashboard__trafic {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__proh-activity {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__requests {
    max-width: 49%;
    width: 100%;
  }

  .dashboard__request {
    max-width: 49%;
    width: 100%;
  }

  .dashboard__top-categories {
    width: 24.65%;
  }

  .dashboard__top-detections {
    width: 24.65%;
  }

  .dashboard__anomalies {
    width: 49.5%;
  }

  .dashboard__logs {
    width: 49.5%;
  }
}

@media screen and (max-width: 1919px) {
  .dashboard__trafic {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__proh-activity {
    max-width: 100%;
    width: 100%;
  }


  .dashboard__requests {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__request {
    width: 49.3%;
  }

  .dashboard__top-categories {
    width: 49.3%;
  }

  .dashboard__top-detections {
    width: 49.3%;
  }

  .dashboard__anomalies {
    width: 49.3%;
  }

  .dashboard__logs {
    width: 49.3%;
  }

  .dashboard-node__trafic {
    width: 49.3%;
  }
}

@media screen and (max-width: 1500px) {
  .dashboard__trafic {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__requests {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__request {
    width: 49%;
  }

  .dashboard__top-categories {
    width: 49%;
  }

  .dashboard__top-detections {
    width: 49%;
  }

  .dashboard__anomalies {
    width: 49%;
  }

  .dashboard__logs {
    width: 49%;
  }

  .dashboard-node__trafic {
    width: 49%;
  }
}

@media screen and (max-width: 780px) {
  .dashboard__trafic {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__requests {
    max-width: 100%;
    width: 100%;
  }

  .dashboard__request {
    width: 100%;
  }

  .dashboard__top-categories {
    width: 100%;
  }

  .dashboard__top-detections {
    width: 100%;
  }

  .dashboard__anomalies {
    width: 100%;
  }

  .dashboard__logs {
    width: 100%;
  }

  .dashboard-node__trafic {
    width: 100%;
  }
}

</style>
