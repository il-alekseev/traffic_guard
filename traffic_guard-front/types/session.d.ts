import type { Categories } from "~/types/categories"
import type { TableMeta } from "~/types/meta";
import type { NGFWType } from "~/types/NGFW";

export interface Session {
  id: string;
  datetime_utc: string;
  status: SessionStatus;
  url: string;
  proto: string,
  hostname: string,
  src_ip: string;
  src_port: string,
  src_country: string,
  username: string,
  dst_ip: string,
  dst_port: string,
  dst_country: string,
  category: Categories
}

export type SessionStatus = 'Разрешен' | 'Запрещен' | 'Ожидает' | 'Аномалия';

export type SessionTypes = 
  | 'Разрешен'
  | 'Запрещен'
  | 'VPN'

export type SessionOrderType = 
  | 'id'
  | 'datetime_utc'
  | 'type'
  | 'status'
  | 'url'
  | 'proto'
  | 'hostname'
  | 'src_ip'
  | 'src_port'
  | 'src_country'
  | 'username'
  | 'dst_ip'
  | 'dst_port'
  | 'dst_country'
  | 'category'

export interface SessionTable {
  data: Session[];
  meta: TableMeta;
}

export interface SessionsState {
  sessions: Session[] | null;
}

export interface SessionsActions {
  fetchSessions(page: number = 1, limit: number = 11, search?: string): Promise<SessionTable>
}
