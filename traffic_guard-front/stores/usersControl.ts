import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import type { User, UsersTable } from "~/types/user";
import { getTokenHeaders } from "~/helpers";
import type { UsersControlState } from "~/types/usersControl";
import type { defaultResponse } from "~/types/api";
import { useUserStore } from "./user";

export const useUsersControlStore = defineStore("userControl", {
  state: (): UsersControlState => {
    return {
      users: null,
    };
  },

  getters: {
    // users: (state) => state.users,
  },

  actions: {
    async fetchUsers(page: number = 1, limit: number = 10, search?: string, role?: string, ): Promise<UsersTable> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список пользователей. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список пользователей. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          page,
          limit,
          ...(search ? { search } : {}),
          ...(role ? { role } : {})
        };

        const users = await $api.get<UsersTable>('/users', {
          params,
          ...getTokenHeaders(token)
        });


        if (users) {
          this.users = users.data;
          return users;
        } else {
          throw new Error("Не удалось получить список пользователей");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении списка пользователей");
      }
    },

    async createUser(user:  Partial<User>): Promise<User>  {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось создать пользователя. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось создать пользователя. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();
        const dataToSend = { ...user };
        const createdUser = await $api.post<User>(
          `/users`,
          dataToSend,
          getTokenHeaders(token)
        );
        if (createdUser) {
          return createdUser;
        } else {
          throw new Error("Не удалось создать пользователя");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при создании пользователя");
      }
    },

    addUser(user: User) {
      this.users?.push(user);
    },

    async updateUserInfo(userId: string, userData: Partial<User>): Promise<boolean> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось обновить данные пользователя. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось обновить данные пользователя. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();
        const dataToSend = { ...userData };

        const response = await $api.put<defaultResponse>(
          `/users/${userId}`,
          dataToSend,
          getTokenHeaders(token)
        );

        if (response) {
          return true;
        } else {
          throw new Error("Не удалось обновить данные пользователя");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при обновлении данных пользователя");
      }
    },

    updateUser(user: User) {
      const userToChange = this.users?.find((u) => u.user_id === user.user_id);
      if (userToChange) {
        Object.assign(userToChange, user);
      }
    },

    async resetPassword(userId: string, password: string): Promise<boolean> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось обновить данные пользователя. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось обновить данные пользователя. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const data = {
          new_password: password,
        };

        const updatedUser = await $api.put<defaultResponse>(
          `/users/${userId}/pass/otp`,
          data,
          getTokenHeaders(token)
        );

        if (updatedUser) {
          return true;
        } else {
          throw new Error("Не удалось обновить данные пользователя");
        }
      } catch (error: any) {
        throw new Error(error.message || "Не удалось обновить данные пользователя");
      }
    },

    async deleteUser(userId: string): Promise<boolean> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось удалить пользователя. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось удалить пользователя. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const deletedUser = await $api.delete<defaultResponse>(
          `/users/${userId}`,
          getTokenHeaders(token)
        );

        if (deletedUser) {
          return true;
        } else {
          throw new Error("Не удалось удалить пользователя");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при удалении пользователя");
      }
    },

    removeUser(userId: string) {
      this.users = (this.users ?? []).filter((u) => u.user_id !== userId);
    },

  },
});
