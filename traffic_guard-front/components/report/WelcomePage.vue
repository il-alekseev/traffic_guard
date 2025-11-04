<template>
  <div class="report-page welcome-page">
    <div class="welcome-content">
      <div class="logo-container">
        <div class="company-logo">
          <img v-if="logoUrl" :src="logoUrl" alt="Company Logo" />
          <div v-else class="logo-placeholder">LOGO</div>
        </div>
      </div>
      
      <h1 class="report-title">Отчет о запрещенных ресурсах</h1>
      
      <div class="report-info">
        <div class="info-block">
          <span class="info-label">Период:</span>
          <span class="info-value">{{ dateRange }}</span>
        </div>
        
        <div class="info-block">
          <span class="info-label">Дата формирования:</span>
          <span class="info-value">{{ generatedDate }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  logoUrl?: string
  dateRange: string
  generatedDate: string
}

defineProps<Props>()
</script>

<style scoped>
.report-page {
  width: 297mm;
  height: 210mm;
  background: white;
  padding: 40px;
  box-sizing: border-box;
  position: relative;
}

.welcome-page {
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-content {
  text-align: center;
  max-width: 600px;
}

.logo-container {
  margin-bottom: 60px;
}

.company-logo img {
  max-width: 200px;
  height: auto;
}

.logo-placeholder {
  width: 200px;
  height: 80px;
  background: #003366;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: bold;
  margin: 0 auto;
  border-radius: 8px;
}

.report-title {
  font-size: 36px;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 60px;
  line-height: 1.3;
}

.report-info {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-block {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 20px;
  background: #f5f5f5;
  border-radius: 8px;
}

.info-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.info-value {
  font-size: 20px;
  color: #1a1a1a;
  font-weight: 600;
}
</style>
