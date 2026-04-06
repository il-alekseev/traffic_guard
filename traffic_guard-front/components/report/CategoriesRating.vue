<template>
  <BaseLayout
    :title="'Рейтинг'"
    :coloredTitle="'категорий'"
  >
    <template #data>
      <div class="column-chart">
        <BaseCard
          v-if="props.data.length > 0"
          :title="'Объём трафика'"
          class="traffic-card"
        >
          <template #DashboardStatistic>
            <DistributedColumnsChart 
              v-if="requestsDataByCategories.length > 0" 
              :data="requestsDataByCategories" 
              :columnNames="categoriesNames" 
              :height="'130px'"
            />
          </template>
        </BaseCard>
      </div>

      <div class="table">
        <div class="table__wrapper">
          <table class="table__content">
            <thead class="table__head">
              <tr class="table__row table__row--header">
                <th class="table__cell table__cell--header table__cell--category">Категория</th>
                <th class="table__cell table__cell--header table__cell--traffic">Трафик</th>
                <th class="table__cell table__cell--header table__cell--requests">Обращения</th>
                <th class="table__cell table__cell--header table__cell--before-block">До блокировки</th>
                <th class="table__cell table__cell--header table__cell--waiting">В ожидании</th>
                <th class="table__cell table__cell--header table__cell--after-block">После блокировки</th>
              </tr>
            </thead>
            <tbody class="table__body">
              <tr v-for="item in props.data" class="table__row">
                <td class="table__cell table__cell--category">
                  <span class="category-circle" :class="`category-badge--reverse--${getModificatorByCategory(item.category as Categories)}`"></span>
                  <span>{{ item.category || '–' }}</span>
                </td>
                <td class="table__cell table__cell--traffic">
                  <div class="traffic">
                    <span class="traffic__item traffic__item--up">
                      ↑ {{ formatTraffic(item.traffic.output) }}
                    </span>
                    <span class="traffic__item traffic__item--down">
                      ↓ {{ formatTraffic(item.traffic.input) }}
                    </span>
                  </div>
                </td>
                <td class="table__cell table__cell--requests">
                  <span class="requests">{{ formatTraffic(item.stat.all) }}</span>
                </td>
                <td class="table__cell table__cell--before-block">
                  <span class="before-block">{{ formatTraffic(item.stat.before_block) }}</span>
                </td>
                <td class="table__cell table__cell--waiting">
                  <span class="waiting">{{ formatTraffic(item.stat.pending) }}</span>
                </td>
                <td class="table__cell table__cell--after-block">
                  <span class="after-block">{{ formatTraffic(item.stat.after_block) }}</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </BaseLayout>
</template>

<script setup lang="ts">
import type { CategoriesDeviceData } from '~/types/reports';
import BaseLayout from '~/components/report/BaseLayout.vue';
import { formatTraffic, getMainColorByCategory } from '~/helpers';
import { getModificatorByCategory } from '~/helpers';
import type { Categories } from '~/types/categories';
import BaseCard from '~/components/report/BaseCard.vue';
import DistributedColumnsChart from './DistributedColumnsChart.vue';


const props = defineProps<{
  data: CategoriesDeviceData[];
}>();

const categoriesNames = computed<string[]>(() => props.data.map((c) => c.category))

const requestsDataByCategories = computed<{name: Categories, value: number, color: string}[]>(() => {
  return props.data.map((cat) => {
    return {
      name: cat.category,
      value: cat.stat.all,
      color: getMainColorByCategory(cat.category)
    }
  })
})

</script>

<style scoped lang="scss">
.table {
  width: 100%;
  overflow-x: auto;

  &__wrapper {
    min-width: 100%;
  }

  &__content {
    width: 100%;
    border-collapse: collapse;
    font-size: 10px;
    line-height: 12px;
    font-weight: 400;
    color: #3F3F46;
  }

  &__head {
    font-weight: 600;
  }

  &__row {
    border-bottom: 0.66px solid #E4E4E7;

    &--group-header {
      border-bottom: 0.66px solid #2563EB;
    }
  }

  &__cell {
    padding: 4px;
    text-align: center;
    vertical-align: middle;

    &--header {
      font-weight: 600;
      user-select: none;
    }

    &--group-header {
      padding: 4px;
      text-align: left;
    }

    &--category {
      display: flex;
      gap: 6px;
      align-items: center;
    }

    &--status {
      text-align: center;
    }

    &--traffic {
      text-align: center;
    }

    &--requests {
      text-align: center;
    }

    &--before-block {
      text-align: center;
    }

    &--waiting {
      text-align: center;
    }

    &--after-block {
      text-align: center;
    }
  }
}

.url-link {
  text-decoration: none;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 2px 6px;
  border-radius: 6px;
  font-size: 8px;
  line-height: 8px;
  font-weight: 500;
  overflow-wrap: anywhere;

  &--green {
    background-color: #DCFCE7;
    color: #008235;
  }

  &--red {
    background-color: #FFE2E2;
    color: #C10007;
  }

  &--yellow {
    background-color: #FEF9C2;
    color: #894B00;
  }
}

.category-circle {
  width: 4px;
  height: 4px;
  border-radius: 50%;
}

.group-header {
  display: inline-block;
  color: #2563EB;
}

.traffic {
  display: flex;
  justify-content: center;
  gap: 8px;
  align-items: center;

  &__item {
    display: inline-flex;
    align-items: center;
    gap: 2px;
    white-space: nowrap;

    &--up {
      color: #2563EB;
    }

    &--down {
      color: #3F3F46;
    }
  }
}

.requests {
  color: #71717A;
}

.after-block, .before-block, .waiting {
  color: #2563EB;
}
</style>
