import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import type { SessionTable, SessionsState, SessionOrderType, SessionTypes, SessionStatus } from '~/types/session';
import { useUserStore } from "./user";
import { getTokenHeaders } from "~/helpers";
import type { Categories } from "~/types/categories";
import type { OrderDir } from "~/types/otherTypes";


export const useSessionsStore = defineStore("sessions", {
  state: (): SessionsState => {
    return {
      sessions: null,
    };
  },

  getters: {
    // sessions: (state) => state.sessions,
  },

  actions: {
    async fetchSessions(from: string = 'now-10m', to: string = 'now', page: number = 1, limit: number = 11, order_by: SessionOrderType = 'id', order_dir: OrderDir = 'asc', search?: string, status?: SessionStatus, category?: Categories, type?: SessionTypes, hostname?: string,): Promise<SessionTable> {
      const userStore = useUserStore();
      try {
        await userStore.ensureValidToken();
      } catch {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список cессий. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        userStore.clearToken();
        throw new Error(
          "Не удалось получить список cессий. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const params: Record<string, string | number> = {
          from,
          to,
          page,
          limit,
          order_by,
          order_dir,
          ...(search ? { search } : {}),
          ...(hostname ? { hostname } : {}),
          ...(category ? { category } : {}),
          ...(type ? { type } : {}),
          ...(status ? { status } : {}),
        };

        const sessions = await $api.get<SessionTable>('/sessions', {
          params,
          ...getTokenHeaders(token)
        });


        if (sessions) {
          this.sessions = sessions.data;
          return sessions;
        } else {
          throw new Error("Не удалось получить список cессий");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при получении списка cессий");
      }
    }
  },
});
