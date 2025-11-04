<template>
  <div class="base-alert" :class="`base-alert--${type}`">
    <div class="base-alert-icon">!</div>
    <div class="base-alert-message"><slot /></div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  type: {
    type: String as () => 'error' | 'warning' | 'success' | 'info',
    default: 'error'
  }
});
</script>

<style>
.base-alert {
  display: flex;
  align-items: flex-start;
  padding: 0.75rem;
  border-radius: var(--radius-md);
  border-left: 0.25rem solid;
}

.base-alert--error {
  background-color: #FEE2E2;
  border-left-color: #EF4444;
}

.base-alert-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.25rem;
  height: 1.25rem;
  border-radius: 50%;
  background-color: #EF4444;
  color: white;
  font-weight: bold;
  font-size: var(--font-size-xs);
  margin-right: 0.75rem;
  flex-shrink: 0;
}

.base-alert-message {
  font-size: var(--font-size-sm);
  color: #B91C1C;
}
</style>