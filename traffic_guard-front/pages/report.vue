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
        :date-range="reportConfig.dateRange"
        :generated-date="reportConfig.generatedDate"
      />

      <ActivityPage
        ref="activityPageRef"
        :traffic="mockTraffic"
        :categoriesTop="mockCategories"
        :resourcesTop="mockResources"
      />

      <AnalyticsTableAllNGFW
        ref="analyticsPageRef"
        :data="mockAllNGFW"
      />

      <AnomaliesTableAllNGFW
        ref="anomaliesPageRef"
        :data="mockAnomaliesAllNGFW"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import html2canvas from 'html2canvas'
import WelcomePage from '~/components/report/WelcomePage.vue';
import ActivityPage from '~/components/report/ActivityPage.vue';
import AnalyticsTableAllNGFW from '~/components/report/AnalyticsTableAllNGFW.vue';
import AnomaliesTableAllNGFW from '~/components/report/AnomaliesTableAllNGFW.vue';
import type { AnomalyReportItem, NGFWReportItem, ReportCategoryItem, ReportResourceItem } from '~/types/reports';
import type { DashboardTraffic } from '~/types/dashboard';


const mockTraffic: DashboardTraffic = {
  count: 20,
  data: {
    Time: [
      "2025-11-04T11:57:23.371048622Z",
      "2025-11-04T13:09:23.371048622Z",
      "2025-11-04T14:21:23.371048622Z",
      "2025-11-04T15:33:23.371048622Z",
      "2025-11-04T16:45:23.371048622Z",
      "2025-11-04T17:57:23.371048622Z",
      "2025-11-04T19:09:23.371048622Z",
      "2025-11-04T20:21:23.371048622Z",
      "2025-11-04T21:33:23.371048622Z",
      "2025-11-04T22:45:23.371048622Z",
      "2025-11-04T23:57:23.371048622Z",
      "2025-11-05T01:09:23.371048622Z",
      "2025-11-05T02:21:23.371048622Z",
      "2025-11-05T03:33:23.371048622Z",
      "2025-11-05T04:45:23.371048622Z",
      "2025-11-05T05:57:23.371048622Z",
      "2025-11-05T07:09:23.371048622Z",
      "2025-11-05T08:21:23.371048622Z",
      "2025-11-05T09:33:23.371048622Z",
      "2025-11-05T10:45:23.371048622Z"
    ],
    Input: [
      100400996,
      100486904,
      96787968,
      96053229,
      100738850,
      100825040,
      92498802,
      100990342,
      101076739,
      92729233,
      101241709,
      101327882,
      97589696,
      96863024,
      101580493,
      101667004,
      97923086,
      97245841,
      101979728,
      89705634
    ],
    Output: [
      122752450,
      122926140,
      119920253,
      116180355,
      123438560,
      123612851,
      113463812,
      123940406,
      124122287,
      113933181,
      124457789,
      124632833,
      121579160,
      117772393,
      125142155,
      125324050,
      122263910,
      118481013,
      125883479,
      112302485
    ]
  }
}

const mockCategories: ReportCategoryItem[] = [
  { category: "Прокси и анонимайзеры", stat: {before_block: 3044, after_block: 3044, pending: 129} },
  { category: "Наркотики", stat: {before_block: 2588, after_block: 2288, pending: 151} },
  { category: "Порнография и секс", stat: {before_block: 2637, after_block: 1437, pending: 151}  },
  { category: "Прокси и анонимайзеры", stat: {before_block: 1511, after_block: 1711, pending: 151}  },
  { category: "Торренты и P2P-сети", stat: {before_block: 1511, after_block: 1711, pending: 15}  },
  { category: "Агрессия, расизм, терроризм", stat: {before_block: 1511, after_block: 1711, pending: 15}  },
  { category: "Азартные игры", stat: {before_block: 1511, after_block: 1711, pending: 15}  },
  { category: "Ботнеты", stat: {before_block: 1512, after_block: 1000, pending: 15}  },
  { category: "Депрессивный контент", stat: {before_block: 512, after_block: 1000, pending: 15}  },
  { category: "Фишинг", stat: {before_block: 512, after_block: 100, pending: 1}  },
];

const mockResources: ReportResourceItem[] =  [
  { resource: "youtube.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 129} },
  { resource: "telegram.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 15} },
  { resource: "discord.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 15}  },
  { resource: "vk.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 15} },
  { resource: "vtb.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 15} },
  { resource: "dinamo.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 15} },
  { resource: "ok.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 15} },
  { resource: "twitch.com", categories: ["Наркотики", "Порнография и секс"], stat: { before_block: 2588, after_block: 2288, pending: 15 }},
  { resource: "minecraft.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 2588, after_block: 2288, pending: 15}  },
  { resource: "tiktok.com", categories: ["Наркотики", "Порнография и секс"], stat: {before_block: 25, after_block: 22, pending: 15} },
];

const mockAllNGFW: NGFWReportItem[] = [
  {
    name: 'NGFW-1',
    input: 1500,
    output: 1200,
    requests: 4500,
    all: 7200,
    anomalies: 23,
    blocks: 156,
    detections: {
      all: 345,
      allowed: 189,
      blocked: 145,
      unresolved: 11
    }
  },
  {
    name: 'NGFW-2',
    input: 2300,
    output: 1800,
    requests: 6200,
    all: 10300,
    anomalies: 45,
    blocks: 234,
    detections: {
      all: 512,
      allowed: 278,
      blocked: 221,
      unresolved: 13
    }
  },
  {
    name: 'NGFW-3',
    input: 980,
    output: 750,
    requests: 2800,
    all: 4530,
    anomalies: 12,
    blocks: 89,
    detections: {
      all: 187,
      allowed: 98,
      blocked: 86,
      unresolved: 3
    }
  },
  {
    name: 'NGFW-4',
    input: 3200,
    output: 2600,
    requests: 8900,
    all: 14700,
    anomalies: 67,
    blocks: 345,
    detections: {
      all: 678,
      allowed: 333,
      blocked: 332,
      unresolved: 13
    }
  },
  {
    name: 'NGFW-5',
    input: 1200,
    output: 950,
    requests: 3500,
    all: 5650,
    anomalies: 18,
    blocks: 123,
    detections: {
      all: 234,
      allowed: 111,
      blocked: 119,
      unresolved: 4
    }
  },
  {
    name: 'NGFW-6',
    input: 4500,
    output: 3800,
    requests: 12500,
    all: 20800,
    anomalies: 89,
    blocks: 478,
    detections: {
      all: 892,
      allowed: 414,
      blocked: 465,
      unresolved: 13
    }
  },
  {
    name: 'NGFW-7',
    input: 800,
    output: 620,
    requests: 2100,
    all: 3520,
    anomalies: 8,
    blocks: 67,
    detections: {
      all: 145,
      allowed: 78,
      blocked: 65,
      unresolved: 2
    }
  },
  {
    name: 'NGFW-8',
    input: 2800,
    output: 2200,
    requests: 7500,
    all: 12500,
    anomalies: 34,
    blocks: 198,
    detections: {
      all: 423,
      allowed: 225,
      blocked: 192,
      unresolved: 6
    }
  },
  {
    name: 'NGFW-9',
    input: 1900,
    output: 1500,
    requests: 5200,
    all: 8600,
    anomalies: 29,
    blocks: 167,
    detections: {
      all: 378,
      allowed: 211,
      blocked: 162,
      unresolved: 5
    }
  },
  {
    name: 'NGFW-10',
    input: 3700,
    output: 3100,
    requests: 10500,
    all: 17300,
    anomalies: 78,
    blocks: 412,
    detections: {
      all: 756,
      allowed: 344,
      blocked: 398,
      unresolved: 14
    }
  },
  {
    name: 'NGFW-11',
    input: 1400,
    output: 1100,
    requests: 4100,
    all: 6600,
    anomalies: 21,
    blocks: 142,
    detections: {
      all: 312,
      allowed: 170,
      blocked: 137,
      unresolved: 5
    }
  },
  {
    name: 'NGFW-12',
    input: 2600,
    output: 2100,
    requests: 7100,
    all: 11800,
    anomalies: 42,
    blocks: 256,
    detections: {
      all: 487,
      allowed: 265,
      blocked: 215,
      unresolved: 7
    }
  },
  {
    name: 'NGFW-13',
    input: 1100,
    output: 850,
    requests: 3300,
    all: 5250,
    anomalies: 15,
    blocks: 108,
    detections: {
      all: 198,
      allowed: 102,
      blocked: 92,
      unresolved: 4
    }
  },
  {
    name: 'NGFW-14',
    input: 3400,
    output: 2800,
    requests: 9400,
    all: 15600,
    anomalies: 72,
    blocks: 389,
    detections: {
      all: 723,
      allowed: 355,
      blocked: 358,
      unresolved: 10
    }
  },
  {
    name: 'NGFW-15',
    input: 1700,
    output: 1300,
    requests: 4800,
    all: 7800,
    anomalies: 26,
    blocks: 178,
    detections: {
      all: 345,
      allowed: 187,
      blocked: 152,
      unresolved: 6
    }
  },
  {
    name: 'NGFW-16',
    input: 4100,
    output: 3400,
    requests: 11200,
    all: 18700,
    anomalies: 83,
    blocks: 445,
    detections: {
      all: 812,
      allowed: 378,
      blocked: 422,
      unresolved: 12
    }
  },
  {
    name: 'NGFW-17',
    input: 900,
    output: 700,
    requests: 2500,
    all: 4100,
    anomalies: 10,
    blocks: 74,
    detections: {
      all: 167,
      allowed: 89,
      blocked: 75,
      unresolved: 3
    }
  },
  {
    name: 'NGFW-18',
    input: 3000,
    output: 2400,
    requests: 8200,
    all: 13600,
    anomalies: 58,
    blocks: 312,
    detections: {
      all: 567,
      allowed: 289,
      blocked: 268,
      unresolved: 10
    }
  },
  {
    name: 'NGFW-19',
    input: 2100,
    output: 1700,
    requests: 5800,
    all: 9600,
    anomalies: 37,
    blocks: 201,
    detections: {
      all: 412,
      allowed: 223,
      blocked: 182,
      unresolved: 7
    }
  },
  {
    name: 'NGFW-20',
    input: 3900,
    output: 3300,
    requests: 10800,
    all: 18000,
    anomalies: 91,
    blocks: 467,
    detections: {
      all: 789,
      allowed: 367,
      blocked: 409,
      unresolved: 13
    }
  },
  {
    name: 'NGFW-21',
    input: 1300,
    output: 1000,
    requests: 3800,
    all: 6100,
    anomalies: 19,
    blocks: 134,
    detections: {
      all: 278,
      allowed: 145,
      blocked: 128,
      unresolved: 5
    }
  },
  {
    name: 'NGFW-22',
    input: 2400,
    output: 1900,
    requests: 6700,
    all: 11000,
    anomalies: 48,
    blocks: 278,
    detections: {
      all: 523,
      allowed: 278,
      blocked: 237,
      unresolved: 8
    }
  },
  {
    name: 'NGFW-23',
    input: 1050,
    output: 800,
    requests: 3000,
    all: 4850,
    anomalies: 14,
    blocks: 95,
    detections: {
      all: 176,
      allowed: 94,
      blocked: 79,
      unresolved: 3
    }
  },
  {
    name: 'NGFW-24',
    input: 3600,
    output: 3000,
    requests: 9800,
    all: 16400,
    anomalies: 76,
    blocks: 423,
    detections: {
      all: 745,
      allowed: 367,
      blocked: 368,
      unresolved: 10
    }
  }
];

const mockAnomaliesAllNGFW: AnomalyReportItem[] = [
  {
    name: "NGFW-1",
    url: "https://example.com/api/1",
    status: "normal",
    live_count: "10 дней",
    traffic: { input: 120, output: 95 },
    requests: 450,
    before_block: 20,
    waiting: 5,
    after_block: 10,
  },
  {
    name: "NGFW-2",
    url: "https://example.com/api/2",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 300, output: 280 },
    requests: 920,
    before_block: 50,
    waiting: 12,
    after_block: 25,
  },
  {
    name: "NGFW-3",
    url: "https://example.com/api/3",
    status: "normal",
    live_count: '10 дней',
    traffic: { input: 560, output: 510 },
    requests: 1350,
    before_block: 70,
    waiting: 18,
    after_block: 35,
  },
  {
    name: "NGFW-4",
    url: "https://example.com/api/4",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 220, output: 190 },
    requests: 600,
    before_block: 30,
    waiting: 8,
    after_block: 14,
  },
  {
    name: "NGFW-1",
    url: "https://example.com/api/5",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 180, output: 140 },
    requests: 520,
    before_block: 25,
    waiting: 6,
    after_block: 11,
  },
  {
    name: "NGFW-2",
    url: "https://example.com/api/6",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 420, output: 390 },
    requests: 980,
    before_block: 40,
    waiting: 10,
    after_block: 20,
  },
  {
    name: "NGFW-3",
    url: "https://example.com/api/7",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 610, output: 580 },
    requests: 1430,
    before_block: 75,
    waiting: 22,
    after_block: 38,
  },
  {
    name: "NGFW-4",
    url: "https://example.com/api/8",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 250, output: 230 },
    requests: 640,
    before_block: 28,
    waiting: 9,
    after_block: 15,
  },
  {
    name: "NGFW-1",
    url: "https://example.com/api/9",
    status: "normal",
    live_count: '10 дней',
    traffic: { input: 310, output: 270 },
    requests: 770,
    before_block: 33,
    waiting: 7,
    after_block: 16,
  },
  {
    name: "NGFW-2",
    url: "https://example.com/api/10",
    status: "normal",
    live_count: '10 дней',
    traffic: { input: 390, output: 340 },
    requests: 850,
    before_block: 36,
    waiting: 11,
    after_block: 22,
  },
  {
    name: "NGFW-3",
    url: "https://example.com/api/11",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 570, output: 540 },
    requests: 1260,
    before_block: 65,
    waiting: 17,
    after_block: 30,
  },
  {
    name: "NGFW-4",
    url: "https://example.com/api/12",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 210, output: 190 },
    requests: 580,
    before_block: 26,
    waiting: 7,
    after_block: 13,
  },
  {
    name: "NGFW-1",
    url: "https://example.com/api/13",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 280, output: 240 },
    requests: 690,
    before_block: 32,
    waiting: 9,
    after_block: 17,
  },
  {
    name: "NGFW-2",
    url: "https://example.com/api/14",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 440, output: 400 },
    requests: 910,
    before_block: 45,
    waiting: 10,
    after_block: 19,
  },
  {
    name: "NGFW-3",
    url: "https://example.com/api/15",
    status: "normal",
    live_count: '10 дней',
    traffic: { input: 600, output: 560 },
    requests: 1380,
    before_block: 72,
    waiting: 20,
    after_block: 33,
  },
  {
    name: "NGFW-4",
    url: "https://example.com/api/16",
    status: "normal",
    live_count: '10 дней',
    traffic: { input: 230, output: 200 },
    requests: 610,
    before_block: 29,
    waiting: 8,
    after_block: 14,
  },
  {
    name: "NGFW-1",
    url: "https://example.com/api/17",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 260, output: 220 },
    requests: 710,
    before_block: 31,
    waiting: 9,
    after_block: 15,
  },
  {
    name: "NGFW-2",
    url: "https://example.com/api/18",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 410, output: 370 },
    requests: 940,
    before_block: 42,
    waiting: 11,
    after_block: 21,
  },
  {
    name: "NGFW-3",
    url: "https://example.com/api/19",
    status: "critical",
    live_count: '10 дней',
    traffic: { input: 550, output: 520 },
    requests: 1320,
    before_block: 68,
    waiting: 18,
    after_block: 31,
  },
  {
    name: "NGFW-4",
    url: "https://example.com/api/20",
    status: "normal",
    live_count: '10 дней',
    traffic: { input: 240, output: 210 },
    requests: 630,
    before_block: 27,
    waiting: 8,
    after_block: 12,
  },
];

const reportConfig = {
  dateRange: '1 Августа 2025 — 1 Сентября 2025',
  generatedDate: '6 ноября',
}

// @ts-ignore
const welcomePageRef = ref<InstanceType<typeof WelcomePage> | null>(null);
// @ts-ignore
const activityPageRef = ref<InstanceType<typeof activityPageRef> | null>(null);
// @ts-ignore
const analyticsPageRef = ref<InstanceType<typeof analyticsPageRef> | null>(null);
// @ts-ignore
const anomaliesPageRef = ref<InstanceType<typeof anomaliesPageRef> | null>(null);

const isGenerating = ref(false);


async function fetchReportData() {
  try {
    // const response = await $fetch('/api/reports/data')
    // tableData.value = response.tableData
    // categoryData.value = response.categoryData
    // trafficData.value = response.trafficData
    // reportConfig.value.dateRange = response.dateRange
  } catch (error) {
    console.error('Ошибка загрузки данных:', error)
  }
}

async function generatePDF() {
  const { jsPDF } = await import('jspdf');
  
  if (!welcomePageRef.value) {
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
      { ref: activityPageRef.value.$el, name: 'Activity' },
      { ref: analyticsPageRef.value.$el, name: 'Analytics' },
      { ref: anomaliesPageRef.value.$el, name: 'Anomalies' },
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
