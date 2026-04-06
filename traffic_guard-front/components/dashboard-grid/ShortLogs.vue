<template>
  <div class="logs">
    <div class="logs__log log" v-for="data in props.data">
      <div class="log__icon">
        <component :is="getIconByLogItem(data.entity)" />
      </div>
      <div class="log__text">
        <span class="log__user-name">{{ data.user_name }}:</span>
        {{ data.description }}
        <span class="log__time"> {{ getTimeByLogItem(data.timestamp) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Log } from '~/types/logs';
import BaseUserCircleIcon from "@/assets/img/base-user-circle.svg"
import DataDocumentIcon from "@/assets/img/empty-data.svg"

const props = defineProps<{
  data: Log[]
}>();

const getIconByLogItem = (entity: 'user' | 'context') => {
  if (entity === 'user') return BaseUserCircleIcon
  if (entity === 'context') return DataDocumentIcon

  return BaseUserCircleIcon
}

const getTimeByLogItem = (timestamp: string): string => {
  const past = new Date(timestamp).getTime()
  const now = Date.now()

  if (isNaN(past)) return ''

  const diffMs = now - past
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))

  if (diffHours <= 0) return 'только что'

  return `${diffHours}ч назад`
}

</script>

<style lang="scss">
.logs {
  display: flex;
  flex-direction: column;
  gap: 0;
  max-width: 29rem;
  max-height: 14rem;
  overflow-y: auto;

  @media screen and (min-width: 1921px) {
    max-width: 100%;
  }

  @media screen and (max-width: 1919px) {
    max-width: unset;
    max-height: 11.5rem;
  }

  &::-webkit-scrollbar {
    width: 0.375rem;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }

  &::-webkit-scrollbar-thumb {
    background-color: #d1d5db;
    border-radius: 3px;

    &:hover {
      background-color: #9ca3af;
    }
  }
}

.log {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 0;
  position: relative;

  &:first-child{
    padding-top: 0;
  }
  
  &:not(:last-child) {
    &::before {
      content: '';
      position: absolute;
      left: 1rem;
      top: 3rem;
      bottom: -1rem;
      width: 2px;
      background-color: #E4E4E7;
    }
  }
  
  &__icon {
    flex-shrink: 0;
    width: 2rem;
    height: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #F4F4F5;
    border-radius: 50%;
    position: relative;
    z-index: 1;
    
    svg {
      width: 1rem;
      height: 1rem;
    }
  }

  &__user-name {
    font-size: 0.875rem;
    font-weight: 500;
    line-height: 1.25;
    color: #27272A;
  }
  
  &__text {
    flex: 1;
    font-size: 0.875rem;
    font-weight: 400;
    line-height: 1.25;
    color: #3F3F46;
  }
  
  &__time {
    flex-shrink: 0;
    font-size: 0.875rem;;
    line-height: 1.25rem;
    color: #A1A1AA;
    white-space: nowrap;
  }
}
</style>