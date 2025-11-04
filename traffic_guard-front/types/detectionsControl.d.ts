import type { TableMeta } from "~/types/meta";
import type { NGFWType } from "~/types/NGFW";
import type { Categories } from "./categories";

export type StatusType = 'blocking' | 'verification' | 'blocked' | 'unknown'

export interface DetectionStats {
  accepted: number,
  denied: number,
  detected: number,
  unresolved: number
}

export interface Detection {
  id: number;
  ip: string;
  location: string;
  domain: string;
  access_count: number;
  host_name: string;
  category: Categories;
  description: string;
  decision: StatusType;
  last_access_datetime: string;
  action: '' | 'Разрешено' | 'Заблокировано'
}

export interface DetectionTable {
  data: Detection[];
  meta: TableMeta;
}

export interface DetectionControlState {
  detections: Detection[] | null;
}

export interface DetectionControlActions {
  fetchDetections(page: number = 1, limit: number = 11): Promise<DetectionTable>
}
