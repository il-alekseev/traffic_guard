<template>
  <div class="error">
    {{ props.fetchError }}
  </div>
</template>

<script setup lang="ts">

const props = defineProps<{
  fetchError: string
}>();


</script>

<style>
.error {
  margin-top: 1rem;
  margin-bottom: 1rem;
  padding: 1.25rem;
  background-color: #FFE2E2;
  color: #C10007;
  border-radius: 0.375rem;
  font-size: 1rem;
  line-height: 1.5rem;
}
</style>