<template>
  <div class="report-page">
    <div class="report__header">
      <div class="report__title">
        {{ props.title }}
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
  width: 297mm;
  height: 210mm;
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
}

.report__title_colored {
  color: #FFFFFF;
  background-color: #2563EB;
  padding: 10px 20px;
  border-radius: 999px;
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