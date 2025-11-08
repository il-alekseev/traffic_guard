import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { Categories, CategoriesState } from "~/types/categories";
import type { ReportData, ReportsState } from "~/types/reports";


export const useReportsStore = defineStore("reports", {
  state: (): ReportsState => {
    return {
      report: null,
    };
  },

  getters: {
    // report: (state) => state.report,
  },

  actions: {
    async fetchReportAllDevices(from: string = 'now-24h', to: string = 'now'): Promise<ReportData> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить отчёт. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить отчёт. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to
        };

        const reportData = await $api.get<ReportData>('/analytics/reports', {
          params,
          ...getTokenHeaders(token)
        });


        if (reportData) {
          this.report = reportData;
          return reportData;
        } else {
          throw new Error("Не удалось получить отчёт");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении отчёта");
      }
    }
  },
});
