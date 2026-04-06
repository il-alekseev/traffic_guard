<template>
  <BaseLayout
    :title="'Список аномалий'"
    :coloredTitle="`по ${deviceName}`"
  >
    <template #data>
      <div class="table">
        <div class="table__wrapper">
          <table class="table__content">
            <thead class="table__head">
              <tr class="table__row table__row--header">
                <th class="table__cell table__cell--header table__cell--url">URL</th>
                <th class="table__cell table__cell--header table__cell--status">Статус/срок жизни</th>
                <th class="table__cell table__cell--header table__cell--traffic">Трафик</th>
                <th class="table__cell table__cell--header table__cell--requests">Обращения</th>
                <th class="table__cell table__cell--header table__cell--before-block">До блокировки</th>
                <th class="table__cell table__cell--header table__cell--waiting">В ожидании</th>
                <th class="table__cell table__cell--header table__cell--after-block">После блокировки</th>
              </tr>
            </thead>
            <tbody class="table__body">
              <tr v-for="item in props.data" class="table__row">
                <td class="table__cell table__cell--url">
                  <span>{{ item.hostname || '–' }}</span>
                </td>
                <td class="table__cell table__cell--status">
                  <span v-if="item.anomaly_block_stat.status && item.anomaly_block_stat.live_count" :class="['status-badge', `status-badge--${getStatusColor(item.anomaly_block_stat.status)}`]">
                    {{ item.anomaly_block_stat.live_count || '–' }}
                  </span>
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
                  <span class="requests">{{ formatTraffic(item.requests) }}</span>
                </td>
                <td class="table__cell table__cell--before-block">
                  <span class="before-block">{{ formatTraffic(item.anomaly_block_stat.before_block) }}</span>
                </td>
                <td class="table__cell table__cell--waiting">
                  <span class="waiting">{{ formatTraffic(item.anomaly_block_stat.pending) }}</span>
                </td>
                <td class="table__cell table__cell--after-block">
                  <span class="after-block">{{ formatTraffic(item.anomaly_block_stat.after_block) }}</span>
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
import type { DeviceAnomaliesData } from '~/types/reports';
import BaseLayout from '~/components/report/BaseLayout.vue';
import { formatTraffic } from '~/helpers';


const props = defineProps<{
  data: DeviceAnomaliesData[];
  deviceName: string
}>();


const getStatusColor = (status: string): string => {
  return status === 'Заблокировано' ? 'red' : status === 'Разрешено' ? 'green' : 'yellow';
};
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

    &--url {
      width: 140px;
      text-align: left;
    }

    &--status {
      width: 106px;
      text-align: center;
    }

    &--traffic {
      width: 106px;
      text-align: center;
    }

    &--requests {
      width: 106px;
      text-align: center;
    }

    &--before-block {
      width: 106px;
      text-align: center;
    }

    &--waiting {
      width: 106px;
      text-align: center;
    }

    &--after-block {
      width: 106px;
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
