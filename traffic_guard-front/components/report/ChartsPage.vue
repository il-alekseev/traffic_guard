<template>
  <div class="report-page charts-page">
    <h2 class="page-title">Статистика и аналитика</h2>
    
    <div class="charts-grid">
      <div class="chart-container">
        <h3 class="chart-title">Распределение по категориям</h3>
        <canvas ref="categoryChart"></canvas>
      </div>
      
      <div class="chart-container">
        <h3 class="chart-title">Трафик по дням</h3>
        <canvas ref="trafficChart"></canvas>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

interface Props {
  categoryData: { label: string; value: number }[]
  trafficData: { date: string; value: number }[]
}

const props = defineProps<Props>()

const categoryChart = ref<HTMLCanvasElement | null>(null)
const trafficChart = ref<HTMLCanvasElement | null>(null)

const initCharts = () => {
  if (categoryChart.value) {
    new Chart(categoryChart.value, {
      type: 'doughnut',
      data: {
        labels: props.categoryData.map(d => d.label),
        datasets: [{
          data: props.categoryData.map(d => d.value),
          backgroundColor: [
            '#3b82f6',
            '#10b981',
            '#f59e0b',
            '#ef4444',
            '#8b5cf6',
          ],
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
          legend: {
            position: 'bottom',
          }
        }
      }
    })
  }

  if (trafficChart.value) {
    new Chart(trafficChart.value, {
      type: 'line',
      data: {
        labels: props.trafficData.map(d => d.date),
        datasets: [{
          label: 'Трафик (MB)',
          data: props.trafficData.map(d => d.value),
          borderColor: '#3b82f6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          tension: 0.4,
          fill: true,
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
          legend: {
            display: false,
          }
        },
        scales: {
          y: {
            beginAtZero: true
          }
        }
      }
    })
  }
}

onMounted(() => {
  initCharts()
})

defineExpose({
  initCharts
})
</script>

<style scoped>
.report-page {
  width: 297mm;
  height: 210mm;
  background: white;
  padding: 40px;
  box-sizing: border-box;
}

.page-title {
  font-size: 32px;
  font-weight: bold;
  color: #1a1a1a;
  margin-bottom: 40px;
  text-align: center;
}

.charts-grid {
  display: flex;
  justify-content: space-between;
}

.chart-container {
  background: #f9fafb;
  padding: 24px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.chart-title {
  font-size: 18px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 20px;
}

canvas {
  max-height: 300px;
}
</style>
