<template>
  <div class="page-container">
    <h1 class="page-title">Генерация отчетов</h1>
    
    <button 
      @click="generatePDF" 
      :disabled="isGenerating"
      class="generate-button"
    >
      {{ isGenerating ? 'Генерация...' : 'Скачать отчёт' }}
    </button>

    <div class="hidden-report">
      <WelcomePage 
        ref="welcomePageRef"
        :logo-url="reportConfig.logoUrl"
        :date-range="reportConfig.dateRange"
        :generated-date="reportConfig.generatedDate"
      />
      
      <TablePage 
        ref="tablePageRef"
        :data="tableData"
      />
      
      <ChartsPage 
        ref="chartsPageRef"
        :category-data="categoryData"
        :traffic-data="trafficData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
// import jsPDF from 'jspdf'
import html2canvas from 'html2canvas'
import WelcomePage from '../components/report/WelcomePage.vue'
import TablePage from '../components/report/TablePage.vue'
import ChartsPage from '../components/report/ChartsPage.vue'

interface TableRow {
  url: string
  category: string
  ngfw: string[]
  traffic: string
  requests: string
  anomaly: string
  anomalyClass: string
}

const welcomePageRef = ref<InstanceType<typeof WelcomePage> | null>(null)
const tablePageRef = ref<InstanceType<typeof TablePage> | null>(null)
const chartsPageRef = ref<InstanceType<typeof ChartsPage> | null>(null)
const isGenerating = ref(false)

const reportConfig = ref({
  logoUrl: '',
  dateRange: '01.01.2025 - 31.01.2025',
  generatedDate: new Date().toLocaleDateString('ru-RU')
})

const tableData = ref<TableRow[]>([
  {
    url: 'Lindsay_Walton-12.org',
    category: 'Насилие',
    ngfw: ['NGFW-1', 'NGFW-2', 'NGFW-3'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Отсутствует',
    anomalyClass: 'absent'
  },
  {
    url: 'Courtney_Henry.ru',
    category: 'Ботнеты',
    ngfw: ['NGFW-2'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Отсутствует',
    anomalyClass: 'absent'
  },
  {
    url: 'Tom_Cook.su',
    category: 'Наркотики',
    ngfw: ['NGFW-2', 'NGFW-3'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Отсутствует',
    anomalyClass: 'absent'
  },
  {
    url: 'Whitney_Francis.design',
    category: 'Фишинг',
    ngfw: ['NGFW-1', 'NGFW-2'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Присутствует',
    anomalyClass: 'present'
  },
  {
    url: 'Leonard_Krasner.com',
    category: 'Торрент',
    ngfw: ['NGFW-1'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Присутствует',
    anomalyClass: 'present'
  },
  {
    url: 'Emily_Selman',
    category: 'Алкоголь',
    ngfw: ['NGFW-1', 'NGFW-2'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Присутствует',
    anomalyClass: 'present'
  },
  {
    url: 'Kristin_Watson',
    category: 'Суицид',
    ngfw: ['NGFW-1', 'NGFW-2', 'NGFW-3'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Присутствует',
    anomalyClass: 'present'
  },
  {
    url: 'Emma_Dorsey',
    category: 'Азартные игры',
    ngfw: ['NGFW-1'],
    traffic: '10MB',
    requests: '12.5k',
    anomaly: 'Отсутствует',
    anomalyClass: 'absent'
  }
])

const categoryData = ref([
  { label: 'Насилие', value: 25 },
  { label: 'Наркотики', value: 20 },
  { label: 'Фишинг', value: 30 },
  { label: 'Торрент', value: 15 },
  { label: 'Другое', value: 10 }
])

const trafficData = ref([
  { date: '01.01', value: 120 },
  { date: '02.01', value: 150 },
  { date: '03.01', value: 180 },
  { date: '04.01', value: 140 },
  { date: '05.01', value: 200 },
  { date: '06.01', value: 170 },
  { date: '07.01', value: 190 }
])

async function fetchReportData() {
  try {
    // const response = await $fetch('/api/reports/data')
    // tableData.value = response.tableData
    // categoryData.value = response.categoryData
    // trafficData.value = response.trafficData
    // reportConfig.value.dateRange = response.dateRange
  } catch (error) {
    console.error('[v0] Ошибка загрузки данных:', error)
  }
}

async function generatePDF() {
  const { jsPDF } = await import('jspdf');
  
  if (!welcomePageRef.value || !tablePageRef.value || !chartsPageRef.value) {
    return
  }

  isGenerating.value = true

  try {
    const pdf = new jsPDF({
      orientation: 'landscape',
      unit: 'mm',
      format: 'a4'
    })

    const pages = [
      { ref: welcomePageRef.value.$el, name: 'Welcome' },
      { ref: tablePageRef.value.$el, name: 'Table' },
      { ref: chartsPageRef.value.$el, name: 'Charts' }
    ]

    for (let i = 0; i < pages.length; i++) {
      const page = pages[i]

      const canvas = await html2canvas(page.ref, {
        scale: 2,
        useCORS: true,
        logging: false,
        backgroundColor: '#ffffff',
        width: page.ref.offsetWidth,
        height: page.ref.offsetHeight
      })

      const imgData = canvas.toDataURL('image/png')
      const imgWidth = 297
      const imgHeight = 210

      if (i > 0) {
        pdf.addPage()
      }

      pdf.addImage(imgData, 'PNG', 0, 0, imgWidth, imgHeight)
      
    }

    const fileName = `report_${new Date().toISOString().split('T')[0]}.pdf`
    pdf.save(fileName)
    
  } catch (error) {
    alert('Произошла ошибка при генерации отчета')
  } finally {
    isGenerating.value = false
  }
}

onMounted(() => {
  // fetchReportData()
})
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.page-title {
  font-size: 32px;
  font-weight: 600;
  margin-bottom: 24px;
  color: #111827;
}

.generate-button {
  background: #2563eb;
  color: white;
  padding: 12px 24px;
  border-radius: 8px;
  border: none;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.generate-button:hover:not(:disabled) {
  background: #1d4ed8;
}

.generate-button:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.hidden-report {
  position: absolute;
  left: -9999px;
  top: 0;
}
</style>
