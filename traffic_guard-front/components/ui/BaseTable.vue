<template>
  <div :class="`${prefix}__table-container table__table-container`">
    <div :class="`${prefix}__table-wrapper table__table-wrapper`">
      <table :class="`${prefix}__table table__table`">
        <thead>
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              :class="`${prefix}__table-column-${column.key}`"
            >
              {{ column.label }}
            </th>
            <th
              v-if="showActions"
              :class="`${prefix}__table-column-button`"
            ></th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td :colspan="columns.length + (showActions ? 1 : 0)" :class="`${prefix}__loading table__loading`">
              {{ loadingMessage }}
            </td>
          </tr>
          <tr v-else-if="items.length === 0">
            <td :colspan="columns.length + (showActions ? 1 : 0)" :class="`${prefix}__empty table__empty`">
              {{ emptyMessage }}
            </td>
          </tr>
          <tr v-else v-for="item in items" :key="item.id || item[itemKey]">
            <td
              v-for="column in columns"
              :key="column.key"
              :class="`${prefix}__table-cell-${column.key}`"
            >
              <slot :name="`cell-${column.key}`" :item="item" :value="item[column.key]">
                {{ item[column.key] }}
              </slot>
            </td>
            <td v-if="showActions" :class="`${prefix}__table-cell-button table__table-cell-button td-button`">
              <button
                :class="`${prefix}__action-button table__action-button`"
                @click="$emit('action-click', item)"
              >
                <slot name="action-button" :item="item">
                </slot>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="!loading && items.length > 0" :class="`${prefix}__footer table__footer`">
      <div :class="`${prefix}__info`">
        Показано от {{ startIndex }} до {{ endIndex }} из {{ totalItems }} результатов
      </div>

      <div :class="`${prefix}__pagination table__pagination`">
        <div
          :class="[
            `${prefix}__pagination-item table__pagination-item`,
            `${prefix}__pagination-prev table__pagination-back`,
            { [`${prefix}__pagination-item_disabled table__pagination-item_disabled`]: currentPage === 1 }
          ]"
          @click="changePage(currentPage - 1)"
        >
          <component :is="arrowLeftIcon" class="arrow-icon" />
        </div>

        <div
          v-for="(page, index) in pagesToShow"
          :key="index"
          :class="[
            `${prefix}__pagination-item table__pagination-item`,
            { [`${prefix}__pagination-item_active table__pagination-item_active`]: page === currentPage }
          ]"
          @click="typeof page === 'number' ? changePage(page) : handleDotsClick(index)"
        >
          {{ page }}
        </div>

        <div
          :class="[
            `${prefix}__pagination-item table__pagination-item`,
            `${prefix}__pagination-next table__pagination-next`,
            { [`${prefix}__pagination-item_disabled table__pagination-item_disabled`]: currentPage === totalPages }
          ]"
          @click="changePage(currentPage + 1)"
        >
          <component :is="arrowLeftIcon" class="arrow-icon" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ArrowLeftIcon from '~/assets/img/arrow-left.svg';

interface Column {
  key: string;
  label: string;
}

interface Props {
  items: any[];
  columns: Column[];
  loading?: boolean;
  fetchError?: string;
  totalItems: number;
  totalPages: number;
  currentPage: number;
  itemsPerPage: number;
  prefix?: string;
  showActions?: boolean;
  arrowLeftIcon?: any;
  emptyMessage?: string;
  loadingMessage?: string;
  itemKey?: string;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  fetchError: '',
  prefix: 'table',
  showActions: true,
  arrowLeftIcon: ArrowLeftIcon,
  emptyMessage: 'Нет данных для отображения',
  loadingMessage: 'Загрузка...',
  itemKey: 'id'
});

const emit = defineEmits<{
  'page-change': [page: number];
  'action-click': [item: any];
}>();

const pagesToShow = computed(() => {
  const pages: (number | string)[] = [];
  const total = props.totalPages;
  const current = props.currentPage;

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
  return (props.currentPage - 1) * props.itemsPerPage + 1;
});

const endIndex = computed(() => {
  const end = props.currentPage * props.itemsPerPage;
  return end > props.totalItems ? props.totalItems : end;
});

const changePage = (page: number) => {
  if (page < 1 || page > props.totalPages) return;
  emit('page-change', page);
};

const handleDotsClick = (index: number) => {
  const total = props.totalPages
  const current = props.currentPage
  
  let isLeftDots = false
  
  if (current <= 3) {
    isLeftDots = false
  } else if (current >= total - 2) {
    isLeftDots = true
  } else {
    isLeftDots = index === 1
  }
  
  if (isLeftDots) {
    changePage(Math.max(1, current - 3))
  } else {
    changePage(Math.min(total, current + 3))
  }
}
</script>

<style scoped>
.table__table-container {
  overflow-x: auto;
  background-color: #FFFFFF;
  border-radius: 1.75rem;
  min-height: calc(100vh - 15rem);
  padding: 1.5rem;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  justify-content: space-between;
}

.table__table-wrapper {
  overflow-x: auto;
}

.table__table {
  table-layout: fixed;
  width: 100%;
  border-collapse: collapse;
}

.table__table td,
.table__table th {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.table__table td.td-button {
  overflow: visible;
}

.table__table th {
  text-align: left;
  padding: 0.90625rem 0.75rem;
  font-weight: 600;
  color: #3F3F46;
  border-bottom: 1px solid #E4E4E7;
  line-height: 1.25rem;
  font-size: 1rem;
  max-width: 18.6875rem;
}

.table__table td {
  padding: 1rem 0.75rem;
  padding-right: 0px;
  border-bottom: 1px solid #E4E4E7;
  color: #71717A;
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.25rem;
  max-width: 16.6875rem;
}

.table__loading,
.table__empty {
  text-align: center;
  padding: 2rem;
  color: #3F3F46;
}

td.table__table-cell-button {
  position: relative;
  padding: 0;
}

.table__action-button {
  background: none;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s ease;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-left: auto;
  padding: 0.5rem;
}

.table__action-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #2563EB;
}

.table__footer {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
}

.table__info {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.25rem;
  color: #3F3F46;
}

.table__pagination {
  display: flex;
  align-items: center;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.table__pagination-item {
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
  transition: background-color 0.2s ease, outline-color 0.2s ease;
}

.table__pagination-back {
  border-top-left-radius: 6px;
  border-bottom-left-radius: 6px;
}

.table__pagination-item:hover:not(.table__pagination-item_disabled):not(.table__pagination-item_active) {
  background-color: #F4F4F5;
}

.table__pagination-item_active {
  background: #2563EB;
  outline: 1px solid #2563EB;
  color: #FFFFFF;
}

.table__pagination-item_disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.arrow-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: #A1A1AA;
}

.table__pagination-next {
  border-top-right-radius: 6px;
  border-bottom-right-radius: 6px;
}

.table__pagination-next .arrow-icon {
  transform: scaleX(-1);
}

@media (max-width: 768px) {
  .table__table-container {
    padding: 1rem;
  }

  .table__footer {
    flex-direction: column;
    align-items: flex-start;
  }

  .table__pagination {
    width: 100%;
    justify-content: center;
  }

  .table__pagination-item {
    min-width: 2rem;
    height: 2rem;
    font-size: 0.875rem;
  }
}
</style>
