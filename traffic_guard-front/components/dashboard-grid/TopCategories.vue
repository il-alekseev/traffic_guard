<template>
  <div class="top-categories">
    <div 
      v-for="(category, index) in categories" 
      :key="index"
      class="category__row"
    >
      <div class="category__info">
        <span class="category__number">{{ index + 1 }}</span>
        <span class="category__name">{{ category.name }}</span>
      </div>
      <span class="category__count">{{ category.access_count }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TopCategories } from '~/types/dashboard';

defineProps<{
  categories: TopCategories;
}>();
</script>

<style scoped lang="scss">
.top-categories {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.category__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.category__info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.category__number {
  font-size: 0.875rem;
  line-height: 0.875rem;
  font-weight: 400;
  color: #3F3F46;
}

.category__name {
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 400;
  color: #3F3F46;
}

.category__count {
  font-size: 0.875rem;
  line-height: 1rem;
  font-weight: 500;
  color: #C10007;

  &::before {
    content: '';
    width: 0.375rem;
    height: 0.375rem;
    border-radius: 100%;
    background-color: #FB2C36;
    display: inline-block;
    margin-right: 0.3125rem;
    margin-bottom: 1px;
  }
}
</style>