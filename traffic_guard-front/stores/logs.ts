import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { Categories } from "~/types/categories";
import type { OrderDir } from "~/types/otherTypes";
import type { LogsState, LogTable } from "~/types/logs";


export const useLogsStore = defineStore("logs", {
  state: (): LogsState => {
    return {
      logs: null,
    };
  },

  getters: {
    // logs: (state) => state.logs,
  },

  actions: {
    async fetchLogs(page: number = 1, limit: number = 10, search?: string, role?: string, contextId?: string): Promise<LogTable> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список логов. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список логов. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          page,
          limit,
          ...(search ? { search } : {}),
          ...(role ? { role } : {}),
          ...(contextId ? { contextId } : {}),
        };

        const logs = await $api.get<LogTable>('/blog/logs', {
          params,
          ...getTokenHeaders(token)
        });


        if (logs) {
          this.logs = logs.data;
          return logs;
        } else {
          throw new Error("Не удалось получить список логов");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении списка логов");
      }
    }
  },
});
