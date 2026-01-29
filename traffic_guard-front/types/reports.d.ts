import type { Categories } from "./categories"
import type { DashboardTraffic } from "./dashboard"

export interface ReportFormData {
  dateFrom: string
  dateTo: string
  deviceSelection: 'all' | 'specific'
  selectedDevice?: string
}

export interface ReportConfig {
  dateRange: string,
  generatedDate: string,
  deviceName: string | undefined
}

export interface GraphLegendItem {
  name: string,
  color: string
}

export type GraphLegend = GraphLegendItem[]

export interface NGFWReportItem {
  hostname: string;
  traffic: {
    input: number;
    output: number;
  }
  requests: number;
  anomaly_block_stat: {
    all: number;
    blocks: number;
    anomalies: number;
  }
  detections: {
    all: number
    allowed: number
    blocked: number
    unresolved: number;
  };
}

export interface AnomalyStatItem {
  live_count: number,
  stat: {
    after_block: number,
    all: number,
    before_block: number,
    pending: number
  },
  status: "Разрешено" | "Заблокировано" | "Не решено",
  traffic: {
    input: number,
    output: number
  },
  url: string
}

export interface AnomalyReportItem {
  hostname: string;
  anomaly_stat: AnomalyStatItem,
}

export interface ReportCategoryItem {
  category: Categories,
  stat: {
    before_block: number,
    after_block: number,
    pending: number
  }
}

export interface ReportResourceItem {
  resource: string,
  categories: Categories[],
  stat: {
    before_block: number,
    after_block: number,
    pending: number
  }
}

export interface RequestsGraphExtraTitle {
  mainValue: string,
  percentValue: string,
  dynamicSign: 'positive' | 'negative' | 'neutral'
}

export interface RequestsAnalyticsTraffic {
  allowed: RequestAnalyticShortData,
  blocked: RequestAnalyticShortData,
  pending: RequestAnalyticShortData
}

export interface RequestAnalyticShortData {
  count: number,
  data: number[],
  time: string[]
}

export interface AnomaliesAndBlocksStat {
  all: number,
  anomalies: number,
  blocks: number,
  status?: string,
  live_count?: string,
  before_block?: number,
  pending?: number,
  after_block?: number
}

export interface DeviceAnalyticsPageData {
  traffic: DashboardTraffic,
  requests_analytics: RequestsAnalyticsTraffic,
  anomaly_block_stat: AnomaliesAndBlocksStat
}

export interface DeviceAnomaliesData {
  hostname: string,
  requests: number,
  anomaly_block_stat: AnomaliesAndBlocksStat,
  detections: {
    all: number,
    allowed: number,
    blocked: number,
    unresolved: number
  },
  traffic: {
    input: number,
    output: number
  }
}

export interface CategoriesDeviceData {
  category: Categories,
  stat: {
    after_block: number,
    all: number,
    before_block: number,
    pending: number
  },
  traffic: {
    input: number,
    output: number
  }
}

export interface ReportData {
  from: string,
  to: string,
  main_activity_page: {
    top_categories: ReportCategoryItem[],
    top_resources: ReportResourceItem[],
    traffic: DashboardTraffic,
  }
  device_analytics_page: {
    analytics: NGFWReportItem[];
  }
  anomalies_list_page: {
    anomalies: AnomalyReportItem[]
  },
  top_anomalies_page: {
    device_anomaly: DeviceAnomaliesData[]
  }
  top_categories_page: {
    categories: CategoriesDeviceData[]
  }
}

export interface ReportDataByDevice {
  from: string,
  to: string,
  hostname: string,
  device_analytics_page: DeviceAnalyticsPageData,
  anomalies_list_page: {
    device_anomaly: DeviceAnomaliesData[]
  },
  categories_page: {
    categories: CategoriesDeviceData[]
  }
}

export interface ReportsState {
  report: ReportData | ReportDataByDevice | null
}
