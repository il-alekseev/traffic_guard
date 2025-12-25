import type { TableMeta } from "~/types/meta";
import type { DropdownItem } from '~/types/dropdown';


export interface Log {
  id: 0,
  context: string,
  description: string,
  entity: 'user' | 'context',
  entity_id: string,
  event_type: 'CREATE' | 'UPDATE' | 'DELETE',
  new_value: string,
  old_value: string,
  timestamp: string,
  user_name: string,
  user_role: string
}

export interface LogTable {
  data: Log[];
  meta: TableMeta;
}

export interface LogsState {
  logs: Log[] | null;
}

export interface LogFilter {
  role: DropdownItem;
  contextId: DropdownItem;
}

export interface LogsActions {
  fetchLogs(page: number = 1, limit: number = 10, search?: string, role?: string, contextId?: string): Promise<LogTable>
}