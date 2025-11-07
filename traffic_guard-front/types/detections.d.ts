import type { TableMeta } from "~/types/meta";
import type { NGFWType } from "~/types/NGFW";
import type { Categories } from "./categories";

export type StatusType = 'blocking' | 'verification' | 'blocked' | 'unknown'

export type ActionType = 'Не решено' | 'Разрешено' | 'Заблокировано' | '';

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
  request_count: number;
  hostname: string;
  category: Categories;
  description: string;
  categorized_at: string;
  action: ActionType
}

export interface DetectionTable {
  data: Detection[];
  meta: TableMeta;
}

export interface DetectionsState {
  detections: Detection[] | null;
}

export interface DetectionsActions {
  fetchDetections(page: number = 1, limit: number = 11): Promise<DetectionTable>
}
