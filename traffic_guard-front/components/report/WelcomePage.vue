<template>
  <div class="report-page welcome-page">
    <div class="welcome-content">
      <div class="logo">
        <div class="logo__icon">
          <AlternativeLogoIcon />
        </div>
        <div class="logo__text">
          Fazenda
        </div>
      </div>
      
      <h1 class="title">Отчёт за период</h1>
      
      <div class="info">
        <div class="info__block">
          <!-- <span class="info__label">Период:</span> -->
          <span class="info__value">{{ dateRange }}</span>
        </div>
        
        <!-- <div class="info__block">
          <span class="info__label">Дата формирования:</span>
          <span class="info__value">{{ generatedDate }}</span>
        </div> -->
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AlternativeLogoIcon from "~/assets/img/alternative-logo.svg"


interface Props {
  dateRange: string
  generatedDate?: string
}

defineProps<Props>()
</script>

<style scoped>
.report-page {
  width: 297mm;
  height: 210mm;
  background: #FFFFFF;
  padding: 40px;
  box-sizing: border-box;
  position: relative;
}

.welcome-page {
  background: #2563EB;
  color: #FFFFFF;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-content {
  text-align: center;
  max-width: 600px;
}

.logo {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.logo__icon {
  max-width: 67px;
  max-height: 67px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo__icon svg {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.logo__text {
  font-weight: 400;
  color: #FFFFFF;
  font-size: 48px;
  line-height: 100%;
}

.title {
  font-size: 72px;
  line-height: 100%;
  font-weight: 600;
  color: #FFFFFF;
}

.info {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info__block {
  display: flex;
  flex-direction: column;
}

.info__label {
  font-size: 14px;
  color: #FFFFFF;
  font-weight: 500;
}

.info__value {
  font-weight: 400;
  font-size: 24px;
  line-height: 100%;
  font-size: 20px;
  color: #FFFFFF;;
}
</style>
