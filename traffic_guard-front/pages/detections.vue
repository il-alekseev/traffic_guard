<template>
  <div class="detections-page">
    <div class="detections__header">
      <StatsComponent :stats="stats" />
      <div class="detections__header-block">
        <div class="detections__header-title-block">
          <h1 class="detections__title page-title">{{ detectionsTitle }}</h1>
          <div class="detections__subtitle-block page-subtitle-block">
            <p class="detections__subtitle page-subtitle">Обновлено {{ minutesAgo }} мин назад</p>
            <div class="detections__refresh-btn page-refresh-btn" @click="fetchDetections">
              <ReloadIcon />
            </div>
          </div>
        </div>
        <div class="detections__filters">
          <div class="detections__datepicker-container">
            <DatePicker v-model="dateRange" />
          </div>
          <DetectionsChipFilter 
            v-model="actionFilter"
            :options="actionFilterOptions"
          />
          <button class="detections__filter-button" @click="showFilters">
            <div class="detections__filter-button-icon">
              <FilterIcon />
            </div>
            <span class="detections__filter-button-text">Фильтр</span>
          </button>
          <DownloadButton />
        </div>
      </div>
    </div>

    <div class="detections-data">
      <ErrorBlock v-if="fetchDetectionsError !== ''" :fetch-error="fetchDetectionsError" />

      <div v-if="loadingDetections" class="detections__loadingDetections">
        Загрузка выявлений...
      </div>

      <div v-else-if="detections.length === 0" class="detections__empty">
        Выявления не найдены
      </div>

      <div v-else class="resources-grid">
        <ResourceCard
          v-for="resource in detections"
          :item="resource"
          @confirm="handleConfirm(resource)"
          @reject="handleReject(resource)"
        />
      </div>

      <BasePagination
        class="detectios__pagination"
        v-if="fetchDetectionsError == '' && detections.length > 0"
        :total="totalDetections"
        :totalPages="totalPages"
        :allowedItemsCount="[9, 18, 27]"
        :currentPage="currentPage"
        :selectedItemsCount="itemsPerPage"
        @update:current-page="handleChangePage"
        @update:selected-items-count="handleChangeItemsPerPage"
      />
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
import { definePageMeta } from '#imports';
import type { ActionFilterOption, Detection, DetectionStats, DetectionTable } from '~/types/detections';
import { useDetectionsStore } from '~/stores/detections';
import type { StatItem } from '~/types/statistics';
import StatsComponent from '~/components/data-display/StatsComponent.vue';
import ResourceCard from '~/components/resource-cards/ResourceCard.vue';
import DatePicker from '~/components/ui/DatePicker.vue';
import DownloadButton from '~/components/ui/DownloadButton.vue';
import ErrorBlock from '~/components/ui/ErrorBlock.vue';
import SideModal from '~/components/ui/SideModal.vue';
import BasePagination from '~/components/ui/BasePagination.vue';
import FilterForm, { type DetectionsFilter } from '~/components/filters/DetectionsFilterForm.vue';
import ReloadIcon from "~/assets/img/reload.svg"
import FilterIcon from "~/assets/img/filter-icon.svg"
import ArrowLeftIcon from "~/assets/img/arrow-left.svg"
import { getCurrentDateWithOffset, isCategory, isValidDateString } from '~/helpers';
import type { Categories } from '~/types/categories';
import DetectionsChipFilter from '~/components/filters/DetectionsChipFilter.vue';


definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const route = useRoute();
const router = useRouter();

const detectionsStore = useDetectionsStore();

const loadingStats = ref(true);
const fetchStatsError = ref('');
const detectionStats = ref<{old: DetectionStats | null, current: DetectionStats | null}>({old: null, current: null})

const stats = computed<StatItem[]>(() => {
  const current = detectionStats.value.current;
  const old = detectionStats.value.old;

  const getValue = (obj: DetectionStats | null, key: keyof DetectionStats): number => {
    return obj?.[key] ?? 0;
  };

  const formatStat = (value: number): string => value.toString();

  const calculateChange = (currentVal: number, oldVal: number): { change: string; changeType: 'increase' | 'decrease' | 'neutral' } => {
    if (oldVal === 0) {
      if (currentVal === 0) {
        return { change: '0%', changeType: 'neutral' };
      } else {
        return { change: '100%', changeType: 'increase' };
      }
    }

    const diff = currentVal - oldVal;
    const percent = Math.round((Math.abs(diff) / oldVal) * 100);

    if (diff > 0) {
      return { change: `${percent}%`, changeType: 'increase' };
    } else if (diff < 0) {
      return { change: `${percent}%`, changeType: 'decrease' };
    } else {
      return { change: '0%', changeType: 'neutral' };
    }
  };

  const items = [
    {
      key: 'detected' as const,
      name: 'Количество выявлений',
    },
    {
      key: 'accepted' as const,
      name: 'Разрешено',
    },
    {
      key: 'denied' as const,
      name: 'Заблокировано',
    }
  ];

  return items.map(({ key, name }) => {
    const currentVal = getValue(current, key);
    const oldVal = getValue(old, key);
    const { change, changeType } = calculateChange(currentVal, oldVal);

    return {
      name,
      stat: formatStat(currentVal),
      previousStat: formatStat(oldVal),
      change,
      changeType,
    };
  });
});

const loadingDetections = ref(true);
const fetchDetectionsError = ref('');

const detections = ref<Detection[]>([]);
const currentPage = ref(1);
const itemsPerPage = ref(9);
const tableMetaData = ref<{ total: number; pages: number } | null>(null);
const totalDetections = computed(() =>
  tableMetaData.value ? tableMetaData.value.total : detections.value.length
);
const detectionsTitle = computed<string>(() => {
  const detectionCount = totalDetections.value;

  if (!detectionCount) return 'Нет нерешенных выявлений';
  
  const pluralForm = (() => {
    const n = detectionCount % 100;
    if (n === 1) return 'нерешенное выявление';
    if (n >= 2 && n <= 4) return 'нерешенных выявления';
    return 'нерешенных выявлений';
  })();

  return `У вас ${detectionCount} ${pluralForm}`;
});
const totalPages = computed(() =>
  tableMetaData.value ? tableMetaData.value.pages : 1
);

const lastUpdated = ref<Date | null>(null);
const minutesAgo = ref(0);
let intervalId: ReturnType<typeof setInterval> | null = null;

const handleChangeItemsPerPage = (value: number) => {
  itemsPerPage.value = value;
  currentPage.value = 1;
  updateUrlParams();
}

const handleChangePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  updateUrlParams();
};

const fetchDetections = async () => {
  loadingDetections.value = true;
  fetchDetectionsError.value = '';

  try {
    const result: DetectionTable = await detectionsStore.fetchDetections(
      dateRange.value.from?.toISOString(),
      dateRange.value.to?.toISOString(),
      currentPage.value,
      itemsPerPage.value,
      statusFilter.value,
      isCategory(categoryFilter.value) ? categoryFilter.value as Categories : undefined,
      deviceFilter.value,
      actionFilter.value
    );

    if (result) {
      detections.value = result.data;
      tableMetaData.value = result.meta;
      loadingCardActs.value = detections.value.map(() => false)
    } else {
      detections.value = [];
    }
  } catch (error) {
    console.error('Ошибка при загрузке выявлений:', error);
    fetchDetectionsError.value = 'Произошла ошибка при загрузке выявлений';
  } finally {
    loadingDetections.value = false;
  }

  lastUpdated.value = new Date();
  updateMinutesAgo();
};

const fetchDetectionStats = async () => {
  const { from, to } = dateRange.value;
  if (!from || !to) return;

  loadingStats.value = true;
  fetchStatsError.value = '';

  const duration = to.getTime() - from.getTime();

  const oldTo = new Date(from.getTime() - 1);
  const oldFrom = new Date(oldTo.getTime() - duration + 1);

  try {
    const [currentResult, oldResult] = await Promise.all([
      detectionsStore.fetchDetectionStats(from.toISOString(), to.toISOString()),
      detectionsStore.fetchDetectionStats(oldFrom.toISOString(), oldTo.toISOString())
    ]);

    detectionStats.value.current = currentResult || null;
    detectionStats.value.old = oldResult || null;
  } catch (error) {
    console.error('Ошибка при загрузке статистики:', error);
    fetchStatsError.value = 'Произошла ошибка при загрузке статистики';
    detectionStats.value.current = null;
    detectionStats.value.old = null;
  } finally {
    loadingStats.value = false;
  }
};

const updateMinutesAgo = () => {
  if (!lastUpdated.value) return;
  const diffMs = Date.now() - lastUpdated.value.getTime();
  minutesAgo.value = Math.floor(diffMs / 60000);
}

const isShowFilters = ref(false);

const showFilters = () => {
  isShowFilters.value = true;
}

const closeFilters = () => {
  isShowFilters.value = false;
}

const statusFilter = ref<string | undefined>();
const categoryFilter = ref<string | undefined>();
const deviceFilter = ref<string | undefined>();
const actionFilter = ref<'Разрешено' | 'Заблокировано' | 'Не решено' | undefined>();
const actionFilterOptions: ActionFilterOption[] = [
  { value: 'Не решено', label: 'Ожидают', color: '#EFB100' },
  { value: 'Заблокировано', label: 'Заблокированы', color: '#FB2C36' },
  { value: 'Разрешено', label: 'Разрешены', color: '#05DF72' }
]

const filtersData = computed<DetectionsFilter | null>(() => {
  const status = statusFilter.value ?? '';
  const category = categoryFilter.value ?? '';
  const device = deviceFilter.value ?? '';

  if (!status && !category && !device) {
    return null;
  }

  return {
    status: { id: status, name: status },
    category: { id: category, name: category },
    device: { id: device, name: device },
  };
});

const initFiltersFromUrl = () => {
  const query = route.query;

  currentPage.value = Number(query.page) || 1;
  itemsPerPage.value = Number(query.per_page) || 9;
  statusFilter.value = query.status != null ? String(query.status) : undefined;
  categoryFilter.value = query.category != null ? String(query.category) : undefined;
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

  if (currentPage.value > 1) query.page = currentPage.value;
  if (itemsPerPage.value !== 9) query.per_page = itemsPerPage.value;
  if (statusFilter.value && statusFilter.value !== '') query.status = statusFilter.value;
  if (categoryFilter.value && categoryFilter.value !== '') query.category = categoryFilter.value;
  if (deviceFilter.value && deviceFilter.value !== '') query.device = deviceFilter.value;
  if (dateRange.value.from) query.from = dateRange.value.from.toISOString();
  if (dateRange.value.to) query.to = dateRange.value.to.toISOString();
  if (actionFilter.value) query.action = actionFilter.value;

  router.replace({ query });
};

const handleSetFilters = (filtersData?: DetectionsFilter) => {
  currentPage.value = 1;

  if (filtersData) {
    filtersData.status.id !== '' ? statusFilter.value = filtersData.status.id : statusFilter.value = undefined;
    filtersData.category.id !== '' ? categoryFilter.value = filtersData.category.id : categoryFilter.value = undefined;
    filtersData.device.id !== '' ? deviceFilter.value = filtersData.device.id : deviceFilter.value = undefined;
  }
  
  updateUrlParams();
}

onMounted(async () => {
  initFiltersFromUrl();
  await fetchDetectionStats();
  await fetchDetections();
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
    await fetchDetectionStats();
    await fetchDetections();
  },
  { deep: true }
);

const dateRange = ref<{ from: Date | null; to: Date | null }>({
  from: getCurrentDateWithOffset(-1, 'd'),
  to: getCurrentDateWithOffset()
})

const loadingCardActs = ref<boolean[]>([])

const handleConfirm = async (item: Detection) => {
  const itemIndex = detections.value.findIndex((d) => d.domain === item.domain);
  if (itemIndex === -1) return;

  if (loadingCardActs.value[itemIndex]) return;

  loadingCardActs.value[itemIndex] = true;

  try {
    const res = await detectionsStore.actForDetection('allow', item.domain);
    if (res ) {
      item.action = 'Разрешено';
    }
  } catch (error: any) {
    console.error('error: ', error)
  } finally {
    loadingCardActs.value[itemIndex] = false;
  }
}

const handleReject = async (item: Detection) => {
  const itemIndex = detections.value.findIndex((d) => d.domain === item.domain);
  if (itemIndex === -1) return;

  if (loadingCardActs.value[itemIndex]) return;

  loadingCardActs.value[itemIndex] = true;

  try {
    const res = await detectionsStore.actForDetection('deny', item.domain);
    if (res ) {
      item.action = 'Заблокировано';
    }
  } catch (error: any) {
    console.error('error: ', error)
  } finally {
    loadingCardActs.value[itemIndex] = false;
  }
}

watch(actionFilter, () => {
  handleSetFilters()
})

watch(dateRange, () => {
  handleSetFilters();
})
</script>

<style lang="scss" scoped>

.detections__header-block {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detections__loadingDetections,
.detections__empty {
  text-align: center;
  padding: 2rem;
  color: #3F3F46;
}

.detections-data {
  display: flex;
  flex-direction: column;
  overflow-x: auto;
  min-height: calc(100vh - 21rem);
}

.resources-grid {
  margin-top: 1.5rem;
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.detections__filter-button {
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

.detections__filter-button-icon {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 1.25rem;
  height: 1.25rem;
  color: #A1A1AA;
}

.detections__filter-button-icon svg {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.detections__filter-button-text {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
}

.detections__filter-selected-area {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.detections__filters {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.detections__download-btn {
  padding: 0.5rem;
  background-color: #FFFFFF;
  box-shadow: 0px 1px 2px 0px #0000000D;
  border-radius: 6px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.detections__download-btn svg {
  width: 1.25rem;
  height: 1.25rem;
  object-fit: contain;
}

.resource-card {
  &__badge {
    padding: 0.25rem 0.5rem;
    border-radius: 6px;
    font-size: 0.75rem;
    line-height: 1rem;
    font-weight: 500;
    text-transform: uppercase;

    &_small {
      padding: 0.125rem 0.365rem;
    }

    &--ngfw-1 {
      background: #DCFCE7;
      color: #008236;
    }

    &--ngfw-2 {
      background: #DBEAFE;
      color: #1447E6;
    }
  }
}

.detectios__pagination {
  margin-top: auto;
}

@media screen and (max-width: 1919px) {
  .resource-card {
    max-width: 49%;
    width: 49%;
  }

  .detections-data {
    min-height: calc(100vh - 20rem);
  }
}

@media screen and (max-width: 1440px) {
  .detections-data {
    min-height: calc(100vh - 13rem);
  }
}



</style>