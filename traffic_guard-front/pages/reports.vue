<template>
  <div class="reports-page">
    <div class="reports__generator">
      <div class="page-title">
        Генератор отчётов
      </div>
      <div class="reports__config-form">
        <ReportForm @submit="getReport" :loading="formLoading" :success="!formLoading && formSucces" :error="formError" />
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

    <div class="hidden-report" v-if="formSucces && (reportData !== null || reportDataByDevice !== null)">
      <WelcomePage
        ref="welcomePageRef"
        :date-range="reportConfig.dateRange"
        :generated-date="reportConfig.generatedDate"
        :deviceName="reportConfig.deviceName"
      />

      <ActivityPage
        v-if="reportData && 
          (
            reportData.main_activity_page.traffic.count > 0 || 
            reportData.main_activity_page.top_categories.length > 0 || 
            reportData.main_activity_page.top_resources.length > 0
          )"
        ref="activityPageRef"
        :traffic="reportData.main_activity_page.traffic"
        :categoriesTop="reportData.main_activity_page.top_categories"
        :resourcesTop="reportData.main_activity_page.top_resources"
      />

      <AnalyticsByDevice
        v-if="reportDataByDevice &&
          (
            reportDataByDevice.device_analytics_page.traffic.count > 0 || 
            reportDataByDevice.device_analytics_page.anomaly_block_stat || 
            reportDataByDevice.device_analytics_page.requests_analytics
          )"
        ref="analiticsByDevicePageRef"
        :traffic="reportDataByDevice.device_analytics_page.traffic"
        :anomaliesAndBlockStat="reportDataByDevice.device_analytics_page.anomaly_block_stat"
        :requestsTraffic="reportDataByDevice.device_analytics_page.requests_analytics"
        :deviceName="reportDataByDevice.hostname"
      />

      <AnalyticsTableAllNGFW
        v-if="reportData && reportData.device_analytics_page.analytics"
        ref="analyticsPageRef"
        :data="reportData.device_analytics_page.analytics"
      />

      <AnomaliesTableAllNGFW
        v-if="reportData && reportData.anomalies_list_page.anomalies"
        ref="anomaliesPageRef"
        :data="reportData.anomalies_list_page.anomalies"
      />

      <AnomaliesTableByDevice
        v-if="reportDataByDevice && reportDataByDevice.anomalies_list_page.device_anomaly"
        ref="anomaliesByDevicePageRef"
        :data="reportDataByDevice.anomalies_list_page.device_anomaly"
        :deviceName="reportDataByDevice.hostname"
      />

      <AnomaliesRating
        v-if="reportData && reportData.top_anomalies_page.device_anomaly"
        ref="anomaliesRatingPageRef"
        :data="reportData.top_anomalies_page.device_anomaly"
      />

      <CategoriesRating
        v-if="reportData && reportData.top_categories_page.categories"
        ref="categoriesRatingAllNGFWPageRef"
        :data="reportData.top_categories_page.categories"
      />

      <CategoriesRating
        v-if="reportDataByDevice && reportDataByDevice.categories_page.categories"
        ref="categoriesRatingPageRef"
        :data="reportDataByDevice.categories_page.categories"
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
import type { ReportConfig, ReportData, ReportDataByDevice, ReportFormData } from '~/types/reports';
import ReportItem from '~/components/reports/ReportItem.vue';
import html2canvas from 'html2canvas'
import WelcomePage from '~/components/report/WelcomePage.vue';
import ActivityPage from '~/components/report/ActivityPage.vue';
import AnalyticsTableAllNGFW from '~/components/report/AnalyticsTableAllNGFW.vue';
import AnomaliesTableAllNGFW from '~/components/report/AnomaliesTableAllNGFW.vue';
import AnalyticsByDevice from '~/components/report/AnalyticsByDevice.vue'
import AnomaliesTableByDevice from '~/components/report/AnomaliesTableByDevice.vue';
import AnomaliesRating from '~/components/report/AnomaliesRating.vue';
import CategoriesRating from '~/components/report/CategoriesRating.vue';
import { normalizeEndDate, normalizeStartDate } from '~/helpers';


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
const reportDataByDevice = ref<ReportDataByDevice | null>(null);
const reportDataDeviceName = ref<string | undefined>();

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

const reportConfig = ref<ReportConfig>({
  dateRange: '',
  generatedDate: formatDate(new Date(), false),
  deviceName: undefined,
})

const setReportConfig = (from: string, to: string, deviceName?: string) => {
  if (!from || !to) return;

  reportConfig.value = {
    dateRange: `${formatDate(from)} — ${formatDate(to)}`,
    generatedDate: formatDate(new Date(), false),
    deviceName: deviceName
  }
}

const getReport = async (formData: ReportFormData) => {
  reportData.value = null;
  reportDataByDevice.value = null;
  formLoading.value = true;
  formSucces.value = false;
  formError.value = '';
  reportDataDeviceName.value = undefined;

  if (formData.deviceSelection === 'specific' && formData.selectedDevice) {
    try {
      const reportDataResponse = await reportStore.fetchReportByDevice(normalizeStartDate(new Date(formData.dateFrom)).toISOString(), normalizeEndDate(new Date(formData.dateTo)).toISOString(), formData.selectedDevice);
      reportDataByDevice.value = reportDataResponse;
      reportDataDeviceName.value = formData.selectedDevice;
      setReportConfig(reportDataByDevice.value.from, reportDataByDevice.value.to, reportDataByDevice.value.hostname);
      formSucces.value = true;
      
      await downloadReport(); 
      formLoading.value = false;
      reportDataByDevice.value = null;
    } catch (error: any) {
      formLoading.value = false;
      formError.value = error;
    }
  } else {
    try {
      const reportDataResponse = await reportStore.fetchReportAllDevices(normalizeStartDate(new Date(formData.dateFrom)).toISOString(), normalizeEndDate(new Date(formData.dateTo)).toISOString());
      reportData.value = reportDataResponse;
      setReportConfig(reportData.value.from, reportData.value.to);
      formSucces.value = true;
      
      await downloadReport();
      formLoading.value = false;
      reportData.value = null;
    } catch (error: any) {
      formLoading.value = false;
      formError.value = error;
    }
  }
}

const welcomePageRef = ref<InstanceType<typeof WelcomePage> | null>(null);
const activityPageRef = ref<InstanceType<typeof ActivityPage> | null>(null);
const analyticsPageRef = ref<InstanceType<typeof AnalyticsTableAllNGFW> | null>(null);
const anomaliesPageRef = ref<InstanceType<typeof AnomaliesTableAllNGFW> | null>(null);
const anomaliesRatingPageRef = ref<InstanceType<typeof AnomaliesRating> | null>(null);
const categoriesRatingAllNGFWPageRef = ref<InstanceType<typeof CategoriesRating> | null>(null);

const analiticsByDevicePageRef = ref<InstanceType<typeof AnalyticsByDevice> | null>(null);
const anomaliesByDevicePageRef = ref<InstanceType<typeof AnomaliesTableByDevice> | null>(null);
const categoriesRatingPageRef = ref<InstanceType<typeof CategoriesRating> | null>(null);

const downloadReport = async () => {
  const { jsPDF } = await import('jspdf');
  
  if (reportData.value == null && reportDataByDevice.value === null) {
    return
  }

  try {
    const pdf = new jsPDF({
      orientation: 'landscape',
      unit: 'px',
      format: 'a4'
    })

    const pages: { ref: HTMLElement; name: string }[] = [];
    if (welcomePageRef.value) pages.push({ ref: welcomePageRef.value.$el, name: 'Welcome' });

    if (reportDataDeviceName.value) {
      if (analiticsByDevicePageRef.value) pages.push({ ref: analiticsByDevicePageRef.value.$el, name: 'Analytics By Device' })
      if (anomaliesByDevicePageRef.value) pages.push({ ref: anomaliesByDevicePageRef.value.$el, name: 'Anomalies By Device' })
      if (categoriesRatingPageRef.value) pages.push({ ref: categoriesRatingPageRef.value.$el, name: 'Categories Rating' })
    } else {
      if (activityPageRef.value) pages.push({ ref: activityPageRef.value.$el, name: 'Activity' })
      if (analyticsPageRef.value)pages.push({ ref: analyticsPageRef.value.$el, name: 'Analytics' })
      if (anomaliesPageRef.value) pages.push({ ref: anomaliesPageRef.value.$el, name: 'Anomalies' })
      if (anomaliesRatingPageRef.value) pages.push({ ref: anomaliesRatingPageRef.value.$el, name: 'Anomalies Rating' })
      if (categoriesRatingAllNGFWPageRef.value) pages.push({ ref: categoriesRatingAllNGFWPageRef.value.$el, name: 'Categories Rating' })
    }


    for (let i = 0; i < pages.length; i++) {
      const page = pages[i]

      const canvas = await html2canvas(page.ref, {
        scale: 2,
        useCORS: true,
        logging: false,
      })

      const imgData = canvas.toDataURL('image/png')
      const pdfWidth = pdf.internal.pageSize.getWidth();
      const pdfHeight = pdf.internal.pageSize.getHeight();
      const imgWidth = canvas.width;
      const imgHeight = canvas.height;

      const ratio = Math.min(pdfWidth / imgWidth, pdfHeight / imgHeight);

      const finalWidth = imgWidth * ratio;
      const finalHeight = imgHeight * ratio;

      if (i > 0) {
        pdf.addPage()
      }

      pdf.addImage(imgData, 'PNG', 0, 0, finalWidth, finalHeight)
      
    }


    let fileName = `report_${new Date().toISOString().split('T')[0]}.pdf`
    if (reportDataDeviceName.value) {
      fileName = `report_${reportDataDeviceName.value.toLowerCase()}_${new Date().toISOString().split('T')[0]}.pdf`
    }
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
  width: 5.5rem;
  height: 5.5rem;
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