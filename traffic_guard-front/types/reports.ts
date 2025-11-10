import type { Categories } from "./categories"
import type { DashboardTraffic } from "./dashboard"

export interface ReportFormData {
  dateFrom: string
  dateTo: string
  deviceSelection: 'all' | 'specific'
  selectedDevice?: string
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
  }
}

export interface ReportsState {
  report: ReportData | null
}
