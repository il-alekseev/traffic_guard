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

export interface AnomalyReportItem {
  hostname: string;
  requests: number;
  traffic: {
    input: number;
    output: number;
  };

  url?: string;
  status?: string;  
  live_count?: string

  detections: {
    all: 0,
    allowed: 0,
    blocked: 0,
    unresolved: 0
  },

  anomaly_block_stat: {
    all: 0,
    anomalies: 0,
    blocks: 0
  },
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
