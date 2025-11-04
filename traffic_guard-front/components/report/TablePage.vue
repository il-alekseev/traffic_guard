<template>
  <div class="report-page table-page">
    <div class="page-header">
      <div class="header-title">
        <span class="title-text">Список</span>
        <span class="title-badge">запрещенных ресурсов</span>
      </div>
      <div class="header-logo"></div>
    </div>

    <table class="data-table">
      <thead>
        <tr>
          <th>URL</th>
          <th>Категория</th>
          <th>NGFW</th>
          <th>Трафик</th>
          <th>Обращения</th>
          <th>Аномалии</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in data" :key="index">
          <td>{{ item.url }}</td>
          <td>{{ item.category }}</td>
          <td>
            <div class="badges">
              <span
                v-for="ngfw in item.ngfw"
                :key="ngfw"
                class="badge"
                :class="`badge-${ngfw.toLowerCase()}`"
              >
                {{ ngfw }}
              </span>
            </div>
          </td>
          <td>{{ item.traffic }}</td>
          <td>{{ item.requests }}</td>
          <td>
            <span class="anomaly" :class="item.anomalyClass">
              {{ item.anomaly }}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
interface TableRow {
  url: string
  category: string
  ngfw: string[]
  traffic: string
  requests: string
  anomaly: string
  anomalyClass: string
}

interface Props {
  data: TableRow[]
}

defineProps<Props>()
</script>

<style scoped>
.report-page {
  width: 297mm;
  height: 210mm;
  background: white;
  padding: 30px;
  box-sizing: border-box;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-text {
  font-size: 28px;
  color: #1a1a1a;
  font-weight: 400;
}

.title-badge {
  background: #2563eb;
  color: white;
  padding: 8px 24px;
  border-radius: 24px;
  font-size: 24px;
  font-weight: 500;
}

.header-logo {
  width: 300px;
  height: 80px;
  background: #003366;
  border-radius: 4px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.data-table thead {
  background: #f8f9fa;
}

.data-table th {
  padding: 12px 8px;
  text-align: left;
  font-weight: 600;
  color: #4b5563;
  border-bottom: 2px solid #e5e7eb;
}

.data-table td {
  padding: 12px 8px;
  border-bottom: 1px solid #e5e7eb;
  color: #1f2937;
}

.badges {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.badge {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.badge-ngfw-1 {
  background: #d1fae5;
  color: #065f46;
}

.badge-ngfw-2 {
  background: #dbeafe;
  color: #1e40af;
}

.badge-ngfw-3 {
  background: #fef3c7;
  color: #92400e;
}

.anomaly {
  font-weight: 500;
}

.anomaly.present {
  color: #2563eb;
}

.anomaly.absent {
  color: #6b7280;
}
</style>
