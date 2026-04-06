<template>
  <div class="logs-page">
    <div class="logs__header">
      <h1 class="logs__title page-title">Логи</h1>
      <div class="logs__toolbar">
        <BaseSearch
          v-model="searchQuery"
          class="logs__search"
          placeholder="Поиск по пользователю, сущности, описанию"
          @search="applyFilters"
        />
        <FilterButton @click="showFilters"/>
      </div>
    </div>

    <ErrorBlock v-if="fetchError !== ''" :fetch-error="fetchError" />

    <BaseTable
      :items="logs"
      :columns="columns"
      :loading="loading"
      :fetch-error="fetchError"
      :total-items="totalLogs"
      :total-pages="totalPages"
      :current-page="currentPage"
      :items-per-page="itemsPerPage"
      :allowedItemsCount="[11, 33, 66]"
      prefix="logs"
      @page-change="handleChangePage"
      @set-items-per-page="handleChangeItemsPerPage"
      @action-click="handleActionClick"
    >
      <template #cell-status="{ value }">
        <span class="logs__table-cell-status">
          {{ value ? value : '–' }}
        </span>
      </template>
      <template #cell-event_type="{ value }">
        <span class="logs__table-cell-event_type" :class="getEventTypeClass(value)">
          {{ value ? value : '–' }}
        </span>
      </template>
      <template #cell-timestamp="{ value }">
        <span class="logs__table-cell-timestamp">
          {{ new Date(value).toLocaleString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
          }) }}
        </span>
      </template>
      <template #cell-user_name="{ value }">
        <span class="logs__table-cell-user_name">
          {{ value ? value : '–' }}
        </span>
      </template>
      <template #cell-description="{ value }">
        <span class="logs__table-cell-description">
          {{ value ? value : '–' }}
        </span>
      </template>
      <template #action-button="{ item }">
        <ContextMenuDotsIcon class="table__action-icon logs__action-icon" />
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
import SideModal from '~/components/ui/SideModal.vue';
import FilterButton from '~/components/ui/FilterButton.vue';
import { useLogsStore } from '~/stores/logs';
import BaseSearch from '~/components/ui/BaseSearch.vue';
import ErrorBlock from '~/components/ui/ErrorBlock.vue';
import BaseTable from '~/components/ui/BaseTable.vue';
import FilterForm from '~/components/filters/LogsFilterForm.vue';
import ContextMenuDotsIcon from '~/assets/img/context-menu-btn.svg';
import type { Log, LogFilter, LogTable } from '~/types/logs';
import { getEventTypeClass, getRoleDisplayName } from '~/helpers';


definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const logsStore = useLogsStore();

const route = useRoute();
const router = useRouter();

const loading = ref(true);
const fetchError = ref('');

const logs = ref<Log[]>([]);
const columns = [
  { key: 'id', label: 'ID' },
  { key: 'event_type', label: 'Тип операции' },
  { key: 'timestamp', label: 'Дата и время' },
  { key: 'entity', label: 'Источник' },
  { key: 'user_name', label: 'Пользователь' },
  { key: 'description', label: 'Описание' },
];
const currentPage = ref(1);
const itemsPerPage = ref(11);
const tableMetaData = ref<{ total: number; pages: number } | null>(null);
const totalLogs = computed(() =>
  tableMetaData.value ? tableMetaData.value.total : logs.value.length
);
const totalPages = computed(() =>
  tableMetaData.value ? tableMetaData.value.pages : 1
);

const fetchLogs = async () => {
  loading.value = true;
  fetchError.value = '';

  try {
    const result: LogTable = await logsStore.fetchLogs(
      currentPage.value,
      itemsPerPage.value,
      searchQuery.value,
      roleFilter.value,
      contextIdFilter.value
    );

    if (result) {
      logs.value = result.data;
      tableMetaData.value = result.meta;
    } else {
      logs.value = [];
    }
  } catch (error) {
    console.error('Ошибка при загрузке логов:', error);
    fetchError.value = 'Произошла ошибка при загрузке логов';
  } finally {
    loading.value = false;
  }
};

const handleChangePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  updateUrlParams();
};

const handleChangeItemsPerPage = (value: number) => {
  itemsPerPage.value = value;
  currentPage.value = 1;
  updateUrlParams();
}

const isShowFilters = ref(false);

const showFilters = () => {
  isShowFilters.value = true;
}

const closeFilters = () => {
  isShowFilters.value = false;
}
const searchQuery = ref('');
const roleFilter = ref<string | undefined>();
const contextIdFilter = ref<string | undefined>();

const filtersData = computed<LogFilter | null>(() => {
  const role = roleFilter.value ?? '';
  const context = contextIdFilter.value ?? '';

  if (!role && !context) {
    return null;
  }

  return {
    role: { id: role, name: getRoleDisplayName(role) },
    contextId: { id: context, name: context },
  };
});

const initFiltersFromUrl = () => {
  const query = route.query;

  currentPage.value = Number(query.page) || 1;
  itemsPerPage.value = Number(query.per_page) || 11;
  searchQuery.value = typeof query.search === 'string' ? query.search : '';
  roleFilter.value = query.role != null ? String(query.role) : undefined;
  contextIdFilter.value = query.contextId != null ? String(query.contextId) : undefined;
};

const updateUrlParams = () => {
  const query: Record<string, string | number> = {};

  if (currentPage.value > 1) query.page = currentPage.value;
  if (itemsPerPage.value !== 11) query.per_page = itemsPerPage.value;
  if (searchQuery.value.trim()) query.search = searchQuery.value.trim();
  if (roleFilter.value && roleFilter.value !== '') query.role = roleFilter.value;
  if (contextIdFilter.value && contextIdFilter.value !== '') query.contextId = contextIdFilter.value;

  router.replace({ query });
};

const handleSetFilters = (filtersData?: LogFilter) => {
  currentPage.value = 1;

  if (filtersData) {
    filtersData.role.id !== '' ? roleFilter.value = filtersData.role.id : roleFilter.value = undefined;
    filtersData.contextId.id !== '' ? contextIdFilter.value = filtersData.contextId.id : contextIdFilter.value = undefined;
  }
  
  updateUrlParams()
}

const applyFilters = () => {
  currentPage.value = 1;
  updateUrlParams();
};

const handleActionClick = (log: Log) => {
  // pass
};


onMounted(async () => {
  initFiltersFromUrl();
  await fetchLogs();
});
watch(
  () => route.query,
  async () => {

    initFiltersFromUrl();
    await fetchLogs();
  },
  { deep: true }
);
</script>

<style lang="scss" scoped>
.logs__header {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.logs__toolbar {
  width: 100%;
  display: flex;
  gap: 0.5rem;
}

:deep(.logs__table) {
  min-width: 77.5rem;
}


:deep(.logs__table-cell-user) {
  color: #2563EB;
}

:deep(.table__action-icon) {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563EB;
}

:deep(.logs__table-column-id) {
  width: 5%;
}
:deep(.logs__table-column-status) {
  width: 7%;
}
:deep(.logs__table-column-timestamp) {
  width: 14%;
}
:deep(.logs__table-column-entity) {
  width: 9%;
}
:deep(.logs__table-column-user) {
  width: 14%;
}
:deep(.logs__table-column-description) {
  width: 44%;
}
:deep(.logs__table-column-button) {
  width: 7%;
}

.logs__table-cell-description {
  color: #3F3F46;
}

.logs__table-cell-timestamp {
  color: #3F3F46;
}

.logs__table-cell-user_name {
  color: #2563EB;
}

.logs-event-type-badge_greeen {
  max-width: fit-content;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 500;
  color: #008236;
  background-color: #DCFCE7;
}

.logs-event-type-badge_red {
  max-width: fit-content;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 500;
  color: #C10007;
  background-color: #FFE2E2;
}

.logs-event-type-badge_yellow {
  max-width: fit-content;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 500;
  color: #894B00;
  background-color: #FEF9C2;
}
</style>