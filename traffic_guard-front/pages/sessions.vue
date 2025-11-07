<template>
  <div class="sessions-page">
    <div class="sessions__header">
      <h1 class="sessions__title page-title">Подозрительные сессии</h1>
      <div class="sessions__toolbar">
        <BaseSearch
          v-model="searchQuery"
          class="sessions__search"
          placeholder="Поиск"
          @search="applyFilters"
        />
        <div class="dashboard__datepicker-container">
          <DatePicker v-model="dateRange" />
        </div>

        <FilterButton @click="showFilters"/>
      </div>
    </div>

    <ErrorBlock v-if="fetchError !== ''" :fetch-error="fetchError" />

    <BaseTable
      :items="sessions"
      :columns="columns"
      :loading="loading"
      :fetch-error="fetchError"
      :total-items="totalSessions"
      :total-pages="totalPages"
      :current-page="currentPage"
      :items-per-page="itemsPerPage"
      prefix="sessions"
      @page-change="handleChangePage"
      @action-click="handleActionClick"
    >
      <template #cell-status="{ value }">
        <span :class="getBadgeClassByStatus(value)">
          {{ value }}
        </span>
      </template>
      <template #cell-hostname="{ value }">
        <span
          class="sessions__table-cell__badge"
          :style="{
            backgroundColor: generateColor(value).background,
            color: generateColor(value).color
          }"
        >
          {{ value || 'Неизвестно' }}
        </span>
      </template>
      <template #cell-url="{ value }">
        <span class="sessions__table-cell-url">{{ value || '-' }}</span>
      </template>
      <template #cell-type="{ value }">
        <span class="sessions__table-cell-type">{{ value }}</span>
      </template>
      <template #cell-category="{ value }">
        <span class="sessions__table-cell-category">{{ value }}</span>
      </template>
      <template #cell-datetime_utc="{ value }">
        <span class="sessions__table-cell-datetime_utc">
          {{ new Date(value).toLocaleString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
          }) }}
        </span>
      </template>
      <template #action-button="{ item }">
        <ContextMenuDotsIcon class="table__action-icon sessions__action-icon" />
      </template>
    </BaseTable>

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
import type { Session, SessionOrderType, SessionStatus, SessionTable } from '~/types/session';
import { useSessionsStore } from '~/stores/session';
import { getBadgeClassByStatus, getNgfwBadgeClass, getCurrentDateWithOffset, isCategory, isSessionStatus, isSessionTypes, isValidDateString } from '~/helpers/index';
import DatePicker from '~/components/ui/DatePicker.vue';
import BaseSearch from '~/components/ui/BaseSearch.vue';
import FilterButton from '~/components/ui/FilterButton.vue';
import ErrorBlock from '~/components/ui/ErrorBlock.vue';
import BaseTable from '~/components/ui/BaseTable.vue';
import SideModal from '~/components/ui/SideModal.vue';
import FilterForm, { type SessionFilter } from '~/components/filters/SessionsFilterForm.vue';
import { useDeviceColors } from '~/composables/useDeviceColors';
import ContextMenuDotsIcon from '~/assets/img/context-menu-btn.svg';
import type { OrderDir } from '~/types/otherTypes';
import type { Categories } from '~/types/categories';


definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const ORDER_BY: SessionOrderType = 'id';
const ORDER_DIR: OrderDir = 'desc'

const { generateColor } = useDeviceColors();

const route = useRoute();
const router = useRouter();

const sessionsStore = useSessionsStore();

const dateRange = ref<{ from: Date | null; to: Date | null }>({
  from: getCurrentDateWithOffset(-1, 'd'),
  to: getCurrentDateWithOffset()
})

const loading = ref(true);
const fetchError = ref('');

const sessions = ref<Session[]>([]);
const columns = [
  { key: 'status', label: 'Статус' },
  { key: 'url', label: 'URL' },
  { key: 'dst_ip', label: 'IP' },
  { key: 'hostname', label: 'NGFW' },
  { key: 'src_ip', label: 'IP SRC' },
  { key: 'type', label: 'Тип сессии' },
  { key: 'category', label: 'Категория' },
  { key: 'datetime_utc', label: 'Дата и время' },
];
const currentPage = ref(1);
const itemsPerPage = ref(11);
const tableMetaData = ref<{ total: number; pages: number } | null>(null);
const totalSessions = computed(() =>
  tableMetaData.value ? tableMetaData.value.total : sessions.value.length
);
const totalPages = computed(() =>
  tableMetaData.value ? tableMetaData.value.pages : 1
);



const fetchSessions = async () => {
  loading.value = true;
  fetchError.value = '';

  try {
    const result: SessionTable = await sessionsStore.fetchSessions(
      dateRange.value.from?.toISOString(),
      dateRange.value.to?.toISOString(),
      currentPage.value,
      itemsPerPage.value,
      ORDER_BY,
      ORDER_DIR,
      searchQuery.value,
      isSessionStatus(statusFilter.value) ? statusFilter.value : undefined ,
      isCategory(categoryFilter.value) ? categoryFilter.value as Categories : undefined,
      isSessionTypes(typesFilter.value) ? typesFilter.value : undefined,
      deviceFilter.value
    );

    if (result) {
      sessions.value = result.data;
      tableMetaData.value = result.meta;
    } else {
      sessions.value = [];
    }
  } catch (error) {
    console.error('Ошибка при загрузке сессий:', error);
    fetchError.value = 'Произошла ошибка при загрузке сессий';
  } finally {
    loading.value = false;
  }
};

const handleChangePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  updateUrlParams();
};

const isShowFilters = ref(false);

const showFilters = () => {
  isShowFilters.value = true;
}

const closeFilters = () => {
  isShowFilters.value = false;
}
const searchQuery = ref('');
const statusFilter = ref<string | undefined>();
const categoryFilter = ref<string | undefined>();
const typesFilter = ref<string | undefined>();
const deviceFilter = ref<string | undefined>();

const filtersData = computed<SessionFilter | null>(() => {
  const status = statusFilter.value ?? '';
  const category = categoryFilter.value ?? '';
  const types = typesFilter.value ?? '';
  const device = deviceFilter.value ?? '';

  if (!status && !category && !types && !device) {
    return null;
  }

  return {
    status: { id: status, name: status },
    category: { id: category, name: category },
    types: { id: types, name: types },
    device: { id: device, name: device },
  };
});

const initFiltersFromUrl = () => {
  const query = route.query;

  currentPage.value = Number(query.page) || 1;
  itemsPerPage.value = Number(query.per_page) || 11;
  searchQuery.value = typeof query.search === 'string' ? query.search : '';
  statusFilter.value = query.status != null ? String(query.status) : undefined;
  categoryFilter.value = query.category != null ? String(query.category) : undefined;
  typesFilter.value = query.types != null ? String(query.types) : undefined;
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
  if (itemsPerPage.value !== 11) query.per_page = itemsPerPage.value;
  if (searchQuery.value.trim()) query.search = searchQuery.value.trim();
  if (dateRange.value.from) query.from = dateRange.value.from.toISOString();
  if (dateRange.value.to) query.to = dateRange.value.to.toISOString();
  if (statusFilter.value && statusFilter.value !== '') query.status = statusFilter.value;
  if (categoryFilter.value && categoryFilter.value !== '') query.category = categoryFilter.value;
  if (typesFilter.value && typesFilter.value !== '') query.types = typesFilter.value;
  if (deviceFilter.value && deviceFilter.value !== '') query.device = deviceFilter.value;

  router.replace({ query });
};

const handleSetFilters = (filtersData?: SessionFilter) => {
  currentPage.value = 1;

  if (filtersData) {
    filtersData.status.id !== '' ? statusFilter.value = filtersData.status.id : statusFilter.value = undefined;
    filtersData.category.id !== '' ? categoryFilter.value = filtersData.category.id : categoryFilter.value = undefined;
    filtersData.types.id !== '' ? typesFilter.value = filtersData.types.id : typesFilter.value = undefined;
    filtersData.device.id !== '' ? deviceFilter.value = filtersData.device.id : deviceFilter.value = undefined;
  }
  
  updateUrlParams()
}

const applyFilters = () => {
  currentPage.value = 1;
  updateUrlParams();
};

const handleActionClick = (session: Session) => {
  console.log('Action clicked for session:', session);
};

onMounted(async () => {
  initFiltersFromUrl();
  await fetchSessions();
});
watch(
  () => route.query,
  async () => {

    initFiltersFromUrl();
    await fetchSessions();
  },
  { deep: true }
);

watch(dateRange, () => {
  applyFilters();
})
</script>

<style scoped>
.sessions__header {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.sessions__title {
  min-width: 23rem;
}

.sessions__toolbar {
  width: 100%;
  display: flex;
  gap: 0.5rem;
}

.sessions__search {
  max-width: 66.625rem;
  min-width: 16rem;
  width: 100%;
}

:deep(.sessions__table-cell-status span) {
  padding: 0.25rem 0.5rem;
  font-weight: 500;
  font-size: 0.75rem;
  line-height: 1rem;
  background: #71717a69;
  color: #3F3F46;
  border-radius: 6px;
}

:deep(.sessions__table-cell-status span.session-status-badge_greeen) {
  color: #008236;
  background-color: #DCFCE7;
}

:deep(.sessions__table-cell-status span.session-status-badge_red) {
  color: #C10007;
  background-color: #FFE2E2;
}

:deep(.sessions__table-cell-status span.session-status-badge_yellow) {
  color: #894B00;
  background-color: #FEF9C2;
}

:deep(.sessions__table-cell__badge) {
  font-weight: 500;
  font-size: 0.75rem;
  line-height: 1rem;
  padding: 0.25rem 0.5rem;
  text-transform: uppercase;
  border-radius: 6px;
}

:deep(.sessions__table-cell-url),
:deep(.sessions__table-cell-category),
:deep(.sessions__table-cell-datetime_utc) {
  color: #3F3F46;
}

:deep(.sessions__table-cell-type) {
  color: #2563EB;
}

:deep(.table__action-icon) {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563EB;
}

:deep(.sessions__table-column-status) {
  width: 7%;
}
:deep(.sessions__table-column-url) {
  width: 24%;
}
:deep(.sessions__table-column-dst_ip) {
  width: 10%;
}
:deep(.sessions__table-column-hostname) {
  width: 13%;
}
:deep(.sessions__table-column-src_ip) {
  width: 10%;
}
:deep(.sessions__table-column-type) {
  width: 11%;
}
:deep(.sessions__table-column-category) {
  width: 11%;
}
:deep(.sessions__table-column-datetime_utc) {
  width: 12%;
}
:deep(.sessions__table-column-button) {
  width: 7%;
}
</style>