<template>
  <div class="report-page">
    <div class="report__header">
      <div class="report__title">
        <p>{{ props.title }}</p>
        <span class="report__title_colored">{{ props.coloredTitle }}</span>
      </div>
      <div class="report__title-logo">
        <LogoWithNameIcon />
      </div>
    </div>
    <div class="report__main">
      <slot name="data" />
    </div>
  </div>

</template>

<script setup lang="ts">
import LogoWithNameIcon from "~/assets/img/logo-with-name.svg"


const props = defineProps<{
  title: string,
  coloredTitle: string
}>();

</script>

<style lang="scss" scoped>
.report-page {
  width: 842px;
  height: 595px;
  background: #FFFFFF;
  padding: 20px 30px;
  box-sizing: border-box;
  position: relative;
}

.report__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.report__title {
  font-size: 24px;
  line-height: 20px;
  font-weight: 400;
  color: #2563EB;
  display: flex;
  align-items: center;
  gap: 6px;
}

.report__title_colored {
  color: #FFFFFF;
  background-color: #2563EB;
  padding: 10px 20px;
  border-radius: 999px;
  line-height: 20px;
  display: inline-block;
}

.report__title-logo {
  max-width: 149px;
  max-height: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.report__title-logo svg {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
</style>