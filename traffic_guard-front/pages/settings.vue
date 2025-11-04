<template>
  <div class="settings-page">
    <div class="settings__header">
      <h1 class="settings__title page-title">Настройки</h1>
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