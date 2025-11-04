<template>
  <div class="detections-page">
    <div class="detections__header">
      <StatsComponent :stats="stats" />
      <div class="detections__header-block">
        <div class="detections__header-title-block">
          <h1 class="detections__title page-title">У вас {{ totalDetections }} нерешенных выявлений</h1>
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
          @confirm="handleConfirm(resource.id)"
          @reject="handleReject(resource.id)"
        />
      </div>
      <div v-if="fetchDetectionsError == '' && detections.length > 0" class="detections__footer">
        <div class="detections__info">
          <p>Показано от {{ startIndex }} до {{ endIndex }} из {{ totalDetections }} результатов</p>
        </div>
        <div class="detections__pagination">
          <div
            class="detections__pagination-item detections__pagination-back"
            :class="currentPage === 1 ? 'detections__pagination-item_disabled' : ''"
            @click="currentPage !== 1 ? handleChangePage(currentPage - 1) : null"
          >
            <ArrowLeftIcon class="arrow-icon" />
          </div>

          <div
            v-for="(page, index) in pagesToShow"
            :key="index"
            class="detections__pagination-item"
            :class="[
              page === currentPage ? 'detections__pagination-item_active' : '',
              page === '...' ? 'detections__pagination-item_dots' : ''
            ]"
            @click="typeof page === 'number' ? handleChangePage(page) : handleDotsClick(index === 1 ? 'left' : 'right')"
          >
            {{ page }}
          </div>

          <div
            class="detections__pagination-item detections__pagination-next"
            :class="currentPage === totalPages ? 'detections__pagination-item_disabled' : ''"
            @click="currentPage !== totalPages ? handleChangePage(currentPage + 1) : null"
          >
            <ArrowLeftIcon class="arrow-icon" />
          </div>
        </div>
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
import { definePageMeta } from '#imports';
import type { Detection, DetectionStats, DetectionTable } from '~/types/detectionsControl';
import { useDetectionsControlStore } from '~/stores/detectionsControl';
import type { StatItem } from '~/types/statistics';
import StatsComponent from '~/components/data-display/StatsComponent.vue';
import ResourceCard from '~/components/resource-cards/ResourceCard.vue';
import DatePicker from '~/components/ui/DatePicker.vue';
import DownloadButton from '~/components/ui/DownloadButton.vue';
import ErrorBlock from '~/components/ui/ErrorBlock.vue';
import SideModal from '~/components/ui/SideModal.vue';
import FilterForm, { type DetectionsFilter } from '~/components/filters/DetectionsFilterForm.vue';
import ReloadIcon from "~/assets/img/reload.svg"
import FilterIcon from "~/assets/img/filter-icon.svg"
import ArrowLeftIcon from "~/assets/img/arrow-left.svg"
import { getCurrentDateWithOffset, isCategory } from '~/helpers';


definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const route = useRoute();
const router = useRouter();

const detectionControlStore = useDetectionsControlStore();

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
const itemsPerPage = ref(6);
const tableMetaData = ref<{ total: number; pages: number } | null>(null);
const totalDetections = computed(() =>
  tableMetaData.value ? tableMetaData.value.total : detections.value.length
);
const totalPages = computed(() =>
  tableMetaData.value ? tableMetaData.value.pages : 1
);
const pagesToShow = computed(() => {
  const pages: (number | string)[] = [];
  const total = totalPages.value;
  const current = currentPage.value;

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i);
  } else {
    if (current <= 3) {
      pages.push(1, 2, 3, '...', total - 2, total - 1, total);
    } else if (current >= total - 2) {
      pages.push(1, 2, '...', total - 2, total - 1, total);
    } else {
      pages.push(1, '...', current - 1, current, current + 1, '...', total);
    }
  }

  return pages;
});
const startIndex = computed(() => {
  return (currentPage.value - 1) * itemsPerPage.value + 1;
});
const endIndex = computed(() => {
  const end = currentPage.value * itemsPerPage.value;
  return end > totalDetections.value ? totalDetections.value : end;
});

const lastUpdated = ref<Date | null>(null);
const minutesAgo = ref(0);
let intervalId: ReturnType<typeof setInterval> | null = null;
  

const handleChangePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  updateUrlParams();
};

const handleDotsClick = (dotsPosition: 'left' | 'right') => {
  const total = totalPages.value;
  const current = currentPage.value;

  if (dotsPosition === 'left') {
    handleChangePage(Math.max(1, current - 3));
  } else {
    handleChangePage(Math.min(total, current + 3));
  }
}

const fetchDetections = async () => {
  loadingDetections.value = true;
  fetchDetectionsError.value = '';

  try {
    const result: DetectionTable = await detectionControlStore.fetchDetections(
      dateRange.value.from?.toISOString(),
      dateRange.value.to?.toISOString(),
      currentPage.value,
      itemsPerPage.value,
      statusFilter.value,
      isCategory(categoryFilter.value) ? categoryFilter.value : undefined,
      deviceFilter.value
    );

    if (result) {
      detections.value = result.data;
      tableMetaData.value = result.meta;
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
      detectionControlStore.fetchDetectionStats(from.toISOString(), to.toISOString()),
      detectionControlStore.fetchDetectionStats(oldFrom.toISOString(), oldTo.toISOString())
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
  itemsPerPage.value = Number(query.per_page) || 6;
  statusFilter.value = query.status != null ? String(query.status) : undefined;
  categoryFilter.value = query.category != null ? String(query.category) : undefined;
  deviceFilter.value = query.device != null ? String(query.device) : undefined;
  dateRange.value.from = typeof query.from === 'string' ? new Date(query.from) : getCurrentDateWithOffset(-1, 'd');
  dateRange.value.to = typeof query.to === 'string' ? new Date(query.to) : getCurrentDateWithOffset();
};

const updateUrlParams = () => {
  const query: Record<string, string | number> = {};

  if (currentPage.value > 1) query.page = currentPage.value;
  if (itemsPerPage.value !== 6) query.per_page = itemsPerPage.value;
  if (statusFilter.value && statusFilter.value !== '') query.status = statusFilter.value;
  if (categoryFilter.value && categoryFilter.value !== '') query.category = categoryFilter.value;
  if (deviceFilter.value && deviceFilter.value !== '') query.device = deviceFilter.value;
  if (dateRange.value.from) query.from = dateRange.value.from.toISOString();
  if (dateRange.value.to) query.to = dateRange.value.to.toISOString();

  router.replace({ query });
};

const handleSetFilters = (filtersData?: DetectionsFilter) => {
  currentPage.value = 1;

  if (filtersData) {
    filtersData.status.id !== '' ? statusFilter.value = filtersData.status.id : statusFilter.value = undefined;
    filtersData.category.id !== '' ? categoryFilter.value = filtersData.category.id : categoryFilter.value = undefined;
    filtersData.device.id !== '' ? deviceFilter.value = filtersData.device.id : deviceFilter.value = undefined;
  }
  
  updateUrlParams()
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

const handleConfirm = (id: number) => {
  console.log('Подтверждено:', id)
}

const handleReject = (id: number) => {
  console.log('Отклонено:', id)
}

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
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(480px, 1fr));
  gap: 1.5rem;
  
  @media (max-width: 1200px) {
    grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  }
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
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
  width: 20px;
  height: 20px;
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


.detections__footer {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: auto;
  padding-bottom: 1px;
  padding-top: 1rem;
}

.detections__info {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
}

.detections__pagination {
  display: flex;
  align-items: center;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.detections__pagination-item {
  cursor: pointer;
  min-width: 2.5rem;
  height: 2.25rem;
  background: #FFFFFF;
  outline: 1px solid #E4E4E7;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
}

.detections__pagination-back {
  border-top-left-radius: 6px;
  border-bottom-left-radius: 6px;
}

.detections__pagination-item_active {
  background: #2563EB;
  outline: 1px solid #2563EB;
  color: #FFFFFF;
}

.arrow-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #A1A1AA;
}

.detections__pagination-next {
  border-top-right-radius: 6px;
  border-bottom-right-radius: 6px;
}

.detections__pagination-next .arrow-icon {
  transform: scaleX(-1);
}

.detections__pagination-item_disabled {
  cursor: not-allowed;
}

</style>