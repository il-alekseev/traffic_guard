<template>
  <div class="report">
    <div class="report__info">
      <div class="report__title">
        <p>Отчет от {{ new Date(props.report.datetime).toLocaleString('ru-RU', {
              day: '2-digit',
              month: 'long',
              year: 'numeric',
            }) }}</p>
      </div>
      <div class="report__meta">
        <span v-if="props.report.device === 'all'" class="report__meta-device">
          Все устройства
        </span>
        <span v-else
          class="report__meta-device-name" 
          :style="{
            backgroundColor: deviceColors.background,
            color: deviceColors.color
          }"
        >
          {{ props.report.device }}
        </span>
        <span>&nbsp;|&nbsp;</span>
        <span class="report__meta-date">
          Создано {{
            new Date(props.report.datetime).toLocaleString('ru-RU', {
              day: '2-digit',
              month: '2-digit',
              year: 'numeric',
              hour: '2-digit',
              minute: '2-digit'
            })
          }}
        </span>
      </div>
    </div>
    <button class="report__download" @click="saveReport" :disabled="isGenerating">
      <div class="report__download-icon">
        <DownloadIcon />
      </div>
      <span>Скачать отчёт</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import DownloadIcon from '~/assets/img/download.svg';
import { useDeviceColors } from '~/composables/useDeviceColors';


const props = defineProps<{
  report: {
    id: string,
    device: string,
    datetime: string
  }
}>();

const { generateColor } = useDeviceColors();
const deviceColors = generateColor(props.report.device);

const isLoading = ref(false);
const fetchError = ref('');
const isGenerating = ref(false);

const saveReport = async () => {
  console.log('saveReport')
}

</script>

<style scoped lang="scss">
.hidden-report {
  position: absolute;
  left: -9999px;
  top: 0;
}

.report {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.75rem;
  padding-bottom:  0.75rem;
  border-bottom: 1px solid #E4E4E7;

  &:last-child {
    border-bottom: none;
  }
}

.report__info {
  display: flex;
  flex-direction: column;
}

.report__info p {
  font-weight: 400;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
}

.report__meta {
  display: flex;
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.5rem;
  color: #A1A1AA;
}

.report__meta-device-name {
  padding: 0.25rem 0.5rem;
  border-radius: 6px;
  font-weight: 500;
  font-size: 0.75rem;
  line-height: 1rem;
  text-transform: uppercase;
}

.report__download {
  background-color: #FFFFFF;
  border: 1px solid #D4D4D8;
  border-radius: 6px;
  box-shadow: 0px 1px 2px 0px #0000000D;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.75rem;
  height: 2.25rem;
  cursor: pointer;
}

.report__download-icon svg {
  color: #3F3F46;
}

.report__download span {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 600;
  color: #3F3F46;
}
</style>