import type { TableMeta } from "~/types/meta";
import type { DropdownItem } from "~/types/dropdown";


export interface User {
  user_id: string;
  login: string;
  email: string;
  first_name: string;
  last_name: string;
  patronymic: string;
  role: string;
  area: string;
  is_super_admin: boolean;
  created_at: string;
  is_need_to_change_password: boolean;
}

export interface UsersTable {
  data: User[];
  meta: TableMeta;
}

export interface UserForm {
  login: string;
  first_name: string;
  last_name: string;
  patronymic: string;
  role: DropdownItem;
  email: string;
}

export interface UserState {
  user: User | null;
}

export interface UserActions {
  isTokenExpiringSoon(thresholdMs: number): boolean
  setToken(token: string, refreshToken?: string, expiresIn?: number): void;
  clearToken(): void;
  setRefreshToken(token: string): void;
  clearRefreshToken(): void;
  setExpiresIn(expiresIn: number): void
  clearExpiresIn(): void
  ensureValidToken(): Promise<void>
  login(login: string, password: string): Promise<boolean>;
  logout(): void;
  fetchUserInfo(): Promise<User>;
  updateUserInfo(userData: Partial<User>): Promise<boolean>;
  changePassword(newPassword: string): Promise<boolean>;
}
