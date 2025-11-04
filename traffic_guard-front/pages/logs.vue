<template>
  <div class="logs-page">
    <div class="logs__header">
      <h1 class="logs__title page-title">Логи</h1>
    </div>
  </div>
</template>

<script setup lang="ts">
import {definePageMeta} from '#imports';
definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});
</script>

<style>

</style>