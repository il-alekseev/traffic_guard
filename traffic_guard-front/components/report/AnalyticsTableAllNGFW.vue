<template>
  <BaseLayout
    :title="'Аналитика'"
    :coloredTitle="'по всем NGFW'"
  >
    <template #data>
      <div class="table">
        <div class="table__wrapper">
          <table class="table__content">
            <thead class="table__head">
              <tr class="table__row table__row--header">
                <th class="table__cell table__cell--header table__cell--name">Название</th>
                <th class="table__cell table__cell--header table__cell--traffic">Трафик</th>
                <th class="table__cell table__cell--header table__cell--requests">Обращения</th>
                <th class="table__cell table__cell--header table__cell--anomalies">Аномалии</th>
                <th class="table__cell table__cell--header table__cell--blocks">Блокировки</th>
                <th class="table__cell table__cell--header table__cell--detections">Выявления</th>
              </tr>
            </thead>
            <tbody class="table__body">
              <tr v-for="item in props.data" :key="item.hostname" class="table__row">
                <td class="table__cell table__cell--name">
                  <span class="table__text">{{ item.hostname || '–' }}</span>
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
                  <span class="table__text requests">{{ formatTraffic(item.requests) }}</span>
                </td>
                <td class="table__cell table__cell--anomalies">
                  <span class="metric">
                    {{ item.anomaly_block_stat.anomalies || '–' }}/{{ item.anomaly_block_stat.all || '–' }}
                    <span class="metric__percent">{{ calculatePercent(item.anomaly_block_stat.anomalies, item.anomaly_block_stat.all) }}</span>
                  </span>
                </td>
                <td class="table__cell table__cell--blocks">
                  <span class="metric">
                    {{ item.anomaly_block_stat.blocks || '–' }}/{{ item.anomaly_block_stat.all || '–' }}
                    <span class="metric__percent">{{ calculatePercent(item.anomaly_block_stat.blocks, item.anomaly_block_stat.all) }}</span>
                  </span>
                </td>
                <td class="table__cell table__cell--detections">
                  <div class="detections">
                    <span class="detections__count">{{ item.detections.all || '–' }}</span>
                    <div class="detections__badges">
                      <span class="badge badge--allowed" :title="`Разрешено: ${item.detections.allowed || '–'}`">
                        ● {{ item.detections.allowed || '0' }}
                      </span>
                      <span class="badge badge--blocked" :title="`Заблокировано: ${item.detections.blocked || '–'}`">
                        ● {{ item.detections.blocked || '0' }}
                      </span>
                      <span class="badge badge--unresolved" :title="`Неопределено: ${item.detections.unresolved || '–'}`">
                        ● {{ item.detections.unresolved || '0' }}
                      </span>
                    </div>
                  </div>
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
import type { NGFWReportItem } from '~/types/reports';
import BaseLayout from '~/components/report/BaseLayout.vue';
import { calculatePercent, formatTraffic } from '~/helpers';

const props = defineProps<{
  data: NGFWReportItem[]
}>();
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
  }

  &__cell {
    padding: 4px;
    text-align: center;
    vertical-align: middle;

    &--header {
      font-weight: 600;
      user-select: none;
    }

    &--name {
      min-width: 100px;
      text-align: left;
      width: 12%;
    }

    &--traffic {
      min-width: 140px;
      width: 15%;
    }

    &--requests {
      min-width: 100px;
      width: 12%;
    }

    &--anomalies {
      min-width: 110px;
      width: 13%;
    }

    &--blocks {
      min-width: 110px;
      width: 13%;
    }

    &--detections {
      min-width: 180px;
      width: 20%;
    }
  }

  &__text {
    color: #333;
    font-weight: 400;
  }
}

.requests {
  color: #71717A;
}

.traffic {
  display: flex;
  justify-content: center;
  gap: 8px;

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

.metric {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-weight: 400;

  &__percent {
    color: #A1A1AA;
  }
}

.detections {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;

  &__badges {
    width: 75%;
    display: flex;
    gap: 8px;
    justify-content: space-between;
    flex-wrap: wrap;
  }
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  white-space: nowrap;

  &--allowed {
    color: #37C84F;
  }

  &--blocked {
    color: #FB2C36;
  }

  &--unresolved {
    color: #EFB100;
  }
}
</style>