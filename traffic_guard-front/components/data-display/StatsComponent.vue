<template>
  <div class="stats-container">
    <dl class="stats-grid">
      <div v-for="item in stats" :key="item.name" class="stat-item">
        <dt class="stat-name">{{ item.name }}</dt>
        <dd class="stat-content">
          <div class="stat-value-wrapper">
            <div v-if="!isLoading" class="stat-value">
              {{ item.stat }}
              <span class="stat-previous">было {{ item.previousStat }}</span>
            </div>
            <div v-else class="skeleton skeleton-value"></div>
          </div>

          <div v-if="!isLoading" :class="[
            'stat-badge',
            {
              'stat-badge--increase': item.changeType === 'increase',
              'stat-badge--decrease': item.changeType === 'decrease',
              'stat-badge--neutral': item.changeType === 'neutral'
            }
          ]">
            <ArrowIcon
              v-if="item.changeType !== 'neutral'"
              :size="20"
              :class="[
                'stat-icon',
                {
                  'stat-icon--increase': item.changeType === 'increase',
                  'stat-icon--decrease': item.changeType === 'decrease'
                }
              ]"
            />
            <span class="sr-only">
              {{
                item.changeType === 'increase' ? 'Increased' :
                item.changeType === 'decrease' ? 'Decreased' :
                'No change'
              }} by
            </span>
            {{ item.change }}
          </div>
          <div v-else class="skeleton skeleton-badge"></div>
        </dd>
      </div>
    </dl>
  </div>
</template>

<script setup lang="ts">
import ArrowIcon from '~/assets/img/full-arrow-down.svg';
import type { StatItem } from '~/types/statistics';

interface Props {
  stats: StatItem[]
  isLoading?: boolean
}

withDefaults(defineProps<Props>(), {
  isLoading: false
})
</script>

<style lang="scss" scoped>
.stats-container {
  margin-bottom: 1.5rem;
  
  .stats-title {
    font-size: 1rem;
    font-weight: 600;
    line-height: 1.5rem;
    color: #111827;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(1, minmax(0, 1fr));
    overflow: hidden;
    border-radius: 0.5rem;
    background-color: #ffffff;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
    
    @media (min-width: 768px) {
      grid-template-columns: repeat(3, minmax(0, 1fr));
      
      .stat-item {
        border-right: 1px solid #e5e7eb;
        border-bottom: none;
        
        &:last-child {
          border-right: none;
        }
      }
    }
  }

  .stat-item {
    padding: 1.25rem 1rem;
    border-bottom: 1px solid #e5e7eb;
    
    &:last-child {
      border-bottom: none;
    }
    
    @media (min-width: 640px) {
      padding: 1.5rem;
    }
  }

  .stat-name {
    font-size: 1rem;
    font-weight: 400;
    color: #111827;
  }

  .stat-content {
    margin-top: 0.25rem;
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    
    @media (min-width: 768px) {
      display: block;
    }
    
    @media (min-width: 1024px) {
      display: flex;
    }
  }

  .stat-value-wrapper {
    display: flex;
    align-items: baseline;
  }

  .stat-value {
    display: flex;
    align-items: baseline;
    font-size: 1.5rem;
    font-weight: 600;
    color: #4f46e5;
  }

  .stat-previous {
    margin-left: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: #6b7280;
  }

  .stat-badge {
    display: inline-flex;
    align-items: baseline;
    border-radius: 9999px;
    padding: 0.125rem 0.625rem;
    font-size: 0.875rem;
    font-weight: 500;

    @media (min-width: 768px) {
      margin-top: 0.5rem;
    }

    @media (min-width: 1024px) {
      margin-top: 0;
    }

    &--increase {
      background-color: #d1fae5;
      color: #065f46;
    }

    &--decrease {
      background-color: #fee2e2;
      color: #991b1b;
    }

    &--neutral {
      background-color: #f3f4f6;
      color: #4b5563;
    }
  }

  .stat-icon {
    width: 1.25rem;
    height: 1.25rem;
    margin-left: -0.25rem;
    margin-right: 0.125rem;
    flex-shrink: 0;
    align-self: center;

    &--increase {
      color: #10b981;
      transform: rotate(180deg);
    }

    &--decrease {
      color: #ef4444;
      transform: rotate(0deg);
    }
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border-width: 0;
  }

  .skeleton {
    background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%);
    background-size: 200% 100%;
    animation: loading 1.5s ease-in-out infinite;
    border-radius: 0.25rem;
  }

  .skeleton-value {
    height: 2rem;
    width: 8rem;
  }

  .skeleton-badge {
    height: 1.75rem;
    width: 4rem;
    border-radius: 9999px;
    
    @media (min-width: 768px) {
      margin-top: 0.5rem;
    }
    
    @media (min-width: 1024px) {
      margin-top: 0;
    }
  }

  @keyframes loading {
    0% {
      background-position: 200% 0;
    }
    100% {
      background-position: -200% 0;
    }
  }
}
</style>