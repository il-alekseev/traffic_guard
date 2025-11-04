<template>
  <div class="reports-page">
    <div class="reports__header">
      <h1 class="reports__title page-title">Отчеты</h1>
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