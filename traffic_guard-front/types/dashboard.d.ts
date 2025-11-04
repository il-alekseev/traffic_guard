import type { Categories } from "./categories";

export interface DashboardCategory {
  category: Categories;
  count: number;
}

export type TopCategories = DashboardCategory[];

export interface DashboardDetection {
  domain: string;
  requests_after: number;
  requests_all: number;
  requests_before: number
}

export type TopDetections = DashboardDetection[];

export interface DashboardTraffic {
  count: number;
  data: {
    Input: number[];
    Output: number[];
    Time: string[];
  }
}

export interface DashboardTrafficLegend {
  input: {
    name: string,
    color: string
  }, 
  output: {
    name: string,
    color: string
  }
}

export interface DashboardState {
  topCategories: TopCategories | null;
  topDetections: TopDetections | null;
  trafic: DashboardTraffic | null
  requests: null,
  anomalies: null,
  events: null,
  proh_activity: null,
  devicesState: null
}

export interface DashboardActions {
  fetchTopCategories(from: string = 'now-10m', to: string = 'now', count: number = 5, hostname?: string): Promise<TopCategories>
  fetchDetections(from: string = 'now-10m', to: string = 'now', count: number = 5, hostname?: string): Promise<TopDetections>
  fetchTraffic(from: string = 'now-10m', to: string = 'now', count: number = 5, hostname?: string): Promise<DashboardTraffic>
}

