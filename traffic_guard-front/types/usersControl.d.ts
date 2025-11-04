import { User } from "./user";

export interface UsersControlState {
  users: User[] | null;
}

export interface UserControlActions {
  fetchUsers(page: number = 1, perPage: number = 10, like?: string, role?: string): Promise<UsersTable>;
  createUser(user:  Partial<User>): Promise<User>;
  addUser(user:  User): void;
  updateUserInfo(userId: string, userData: Partial<User>): Promise<boolean>;
  updateUser(user: User): void;
  resetPassword(userId: string, password: string): Promise<boolean>;
  deleteUser(userId: string): Promise<boolean>;
  removeUser(userId: string): void;
}
