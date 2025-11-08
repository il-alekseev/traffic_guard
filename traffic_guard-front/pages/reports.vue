<template>
  <div class="reports-page">
    <div class="reports__generator">
      <div class="page-title">
        Генератор отчётов
      </div>
      <div class="reports__config-form">
        <ReportForm @submit="getReport" :loading="formLoading" :success="formSucces" />
      </div>
    </div>
    <div class="reports__data">
      <div class="page-title">
        Сформированные отчёты
      </div>
      <div v-if="reports.length > 0" class="reports__list">
        <ReportItem v-for="report in reports" :report="report" />
      </div>
      <div v-else class="reports__empty-block">
        <div v-if="loadingReports" class="reports__loading">
          <p>Загрузка...</p>
        </div>
        <div v-else-if="reportsFetchError !== ''" class="reports__error">
          <ErrorBlock v-if="reportsFetchError !== ''" :fetch-error="reportsFetchError" />
        </div>
        <div v-else class="reports__empty-data">
          <div class="reports__empty-img">
            <EmptyDataIcon />
          </div>
          <div class="reports__empty-title">
            <p>У вас пока нет отчетов</p>
          </div>
          <div class="reports__empty-subtitle">
            Для начала создайте свой первый отчет, используя генератор отчётов
          </div>
        </div>
      </div>
    </div>

    <div class="hidden-report" v-if="formSucces && reportData !== null">
      <WelcomePage 
        ref="welcomePageRef"
        :date-range="reportConfig.dateRange"
        :generated-date="reportConfig.generatedDate"
      />

      <ActivityPage
        v-if="reportData.main_activity_page.traffic.count > 0 || reportData.main_activity_page.top_categories.length > 0 || reportData.main_activity_page.top_resources.length > 0"
        ref="activityPageRef"
        :traffic="reportData.main_activity_page.traffic"
        :categoriesTop="reportData.main_activity_page.top_categories"
        :resourcesTop="reportData.main_activity_page.top_resources"
      />

      <AnalyticsTableAllNGFW
        v-if="reportData.device_analytics_page.analytics"
        ref="analyticsPageRef"
        :data="reportData.device_analytics_page.analytics"
      />

      <AnomaliesTableAllNGFW
        v-if="reportData.anomalies_list_page.anomalies"
        ref="anomaliesPageRef"
        :data="reportData.anomalies_list_page.anomalies"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {definePageMeta} from '#imports';
import { useReportsStore } from '~/stores/reports';
import ErrorBlock from '~/components/ui/ErrorBlock.vue';
import EmptyDataIcon from "~/assets/img/empty-data.svg"
import ReportForm from "~/components/reports/ReportGeneratorForm.vue"
import type { ReportData, ReportFormData } from '~/types/reports';
import ReportItem from '~/components/reports/ReportItem.vue';
import html2canvas from 'html2canvas'
import WelcomePage from '~/components/report/WelcomePage.vue';
import ActivityPage from '~/components/report/ActivityPage.vue';
import AnalyticsTableAllNGFW from '~/components/report/AnalyticsTableAllNGFW.vue';
import AnomaliesTableAllNGFW from '~/components/report/AnomaliesTableAllNGFW.vue';


definePageMeta({
  layout: 'dashboard',
  middleware: ['auth']
});

const reportStore = useReportsStore();

const loadingReports = ref(true);
const reportsFetchError = ref('');
const reports = ref<any[]>([]);

const fetchReports = async () => {
  loadingReports.value = false;
  // pass
}

const formLoading = ref(false);
const formError = ref('');
const formSucces = ref(false);

const reportData = ref<ReportData | null>(null);

const formatDate = (date: string | Date, withYear: boolean = true): string => {
  const d = new Date(date);

  const options: Intl.DateTimeFormatOptions = {
    day: 'numeric',
    month: 'long',
  };

  if (withYear) {
    options.year = 'numeric';
  }

  const formatted = d.toLocaleDateString('ru-RU', options);

  return formatted.replace(/(\p{L})/u, c => c.toUpperCase());
};

const reportConfig = ref({
  dateRange: '',
  generatedDate: formatDate(new Date(), false)
})

const setReportConfig = (from: string, to: string) => {
  if (!from || !to) return;

  reportConfig.value = {
    dateRange: `${formatDate(from)} — ${formatDate(to)}`,
    generatedDate: formatDate(new Date(), false)
  }
}

const getReport = async (formData: ReportFormData) => {
  if (formData.deviceSelection === 'specific') return;

  formLoading.value = true;
  formSucces.value = false;
  formError.value = '';

  try {
    const reportDataResponse = await reportStore.fetchReportAllDevices(new Date(formData.dateFrom)?.toISOString(), new Date(formData.dateTo)?.toISOString());
    reportData.value = reportDataResponse;
    setReportConfig(reportData.value.from, reportData.value.to);
    formSucces.value = true;

    await downloadReport();
  } catch (error: any) {
    formError.value = error;
    formSucces.value = false;
  } finally {
    formLoading.value = false;
  }
}


const welcomePageRef = ref<InstanceType<typeof WelcomePage> | null>(null);

const activityPageRef = ref<InstanceType<typeof ActivityPage> | null>(null);

const analyticsPageRef = ref<InstanceType<typeof AnalyticsTableAllNGFW> | null>(null);

const anomaliesPageRef = ref<InstanceType<typeof AnomaliesTableAllNGFW> | null>(null);

const downloadReport = async () => {
  const { jsPDF } = await import('jspdf');
  console.log('reportData.value', reportData.value);
  console.log('welcomePageRef.value', welcomePageRef.value);
  console.log('activityPageRef.value', activityPageRef.value);
  console.log('analyticsPageRef.value', analyticsPageRef.value);
  console.log('anomaliesPageRef.value', anomaliesPageRef.value);
  
  if (reportData.value == null) {
    console.log("LOX");
    return
  }

  try {
    const pdf = new jsPDF({
      orientation: 'landscape',
      unit: 'mm',
      format: 'a4'
    })

    const pages: { ref: HTMLElement; name: string }[] = [
      welcomePageRef.value && { ref: welcomePageRef.value.$el, name: 'Welcome' },
      activityPageRef.value && { ref: activityPageRef.value.$el, name: 'Activity' },
      analyticsPageRef.value && { ref: analyticsPageRef.value.$el, name: 'Analytics' },
      anomaliesPageRef.value && { ref: anomaliesPageRef.value.$el, name: 'Anomalies' },
    ].filter((p): p is { ref: HTMLElement; name: string } => Boolean(p));


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
    throw new Error('Произошла ошибка при создании отчета')
  }
}

onMounted( async () => {
  await fetchReports();
})

</script>

<style lang="scss">
.hidden-report {
  position: absolute;
  left: -9999px;
  top: 0;
}

.page-title {
  font-weight: 500;
  font-size: 1.875rem;
  line-height: 2.25rem;
  color: #3F3F46;
}

.reports-page {
  display: flex;
  gap: 1rem;
}

.reports__generator {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  max-width: 34%;
  width: 100%;
}

.reports__data {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  max-width: 66%;
  width: 100%;
}

.reports__empty-block {
  height: 29.5rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #FFFFFF;
  border-radius: 24px;
}

.reports__empty-data {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.reports__empty-img {
  align-self: center;
  width: 88px;
  height: 88px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.reports__empty-img svg {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.reports__empty-title {
  font-weight: 400;
  font-size: 2rem;
  line-height: 2.25rem;
  text-align: center;
  color: #3F3F46;
}

.reports__empty-subtitle {
  font-weight: 400;
  font-size: 1.25rem;
  line-height: 1.5rem;
  text-align: center;
  color: #A1A1AA;
}

.reports__loading {
  font-weight: 400;
  font-size: 2rem;
  line-height: 2.25rem;
  text-align: center;
  color: #3F3F46;
}

.reports__config-form {
  padding: 1.5rem;
  background-color: #FFFFFF;
  border-radius: 24px;
}

.reports__list {
  padding: 0.75rem 1.5rem;
  background-color: #FFFFFF;
  border-radius: 24px;
}


</style>