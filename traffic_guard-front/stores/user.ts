import { defineStore } from "pinia";
import { useNuxtApp } from "#app";
import { useCookie } from "#app";
import type { User, UserState } from "~/types/user";
import { getTokenHeaders, getRefreshTokenHeaders, getLogoutTokenHeaders, getJwtExpMs } from "~/helpers";
import type { defaultResponse } from "~/types/api";

export const useUserStore = defineStore("user", {
  state: (): UserState => {

    return {
      user: null,
    };
  },

  getters: {
    firstName: (state) => state.user?.first_name || "",
    lastName: (state) => state.user?.last_name || "",
    fullName: (state) => state.user?.first_name + " " + state.user?.last_name || "",
    isSuperAdmin: (state) => state.user?.is_super_admin || false,
    role: (state) => state.user?.role || "",
  },

  actions: {
    setRefreshToken(token: string): void {
      const config = useRuntimeConfig();
      const nodeEnvLevel = config.public.nodeEnv;

      useCookie('refresh_token', {
        maxAge: 60 * 60 * 24 * 14,
        secure: nodeEnvLevel === "production",
        sameSite: "strict",
      }).value = token;
    },

    clearRefreshToken(): void {
      const config = useRuntimeConfig();
      const nodeEnvLevel = config.public.nodeEnv;

      useCookie('refresh_token', {
        maxAge: 60 * 60 * 24 * 14,
        secure: nodeEnvLevel === "production",
        sameSite: "strict",
      }).value = null;
    },

    setToken(token: string, refreshToken?: string, expiresIn?: number): void {
      const config = useRuntimeConfig();
      const nodeEnvLevel = config.public.nodeEnv;

      let maxAgeSec: number;

      if (typeof expiresIn === "number" && isFinite(expiresIn) && expiresIn > 0) {
        maxAgeSec = Math.max(1, Math.floor(expiresIn));
      } else {
        const expMs = getJwtExpMs(token);
        if (expMs && expMs > Date.now()) {
          maxAgeSec = Math.max(1, Math.floor((expMs - Date.now()) / 1000));
        } else {
          maxAgeSec = 60 * 60;
        }
      }

      useCookie("auth_token", {
        maxAge: maxAgeSec,
        secure: nodeEnvLevel === "production",
        sameSite: "strict",
      }).value = token;

      if (refreshToken) {
        this.setRefreshToken(refreshToken);
      }
    },

    clearToken(): void {
      const config = useRuntimeConfig();
      const nodeEnvLevel = config.public.nodeEnv;

      useCookie("auth_token", {
        maxAge: 60 * 60, // ignore
        secure: nodeEnvLevel === "production",
        sameSite: "strict",
      }).value = null;

      this.clearRefreshToken();
    },

    isTokenExpiringSoon(thresholdMs: number = 60 * 1000): boolean {
      const token = useCookie<string | null>("auth_token").value || null;
      if (!token) return true;

      const expMs = getJwtExpMs(token);

      // fail-safe
      if (!expMs || !Number.isFinite(expMs)) return true;

      return expMs <= Date.now() + thresholdMs;
    },

    async ensureValidToken(): Promise<void> {
      const hasRefresh = !!useCookie("refresh_token").value;

      if (this.isTokenExpiringSoon()) {
        if (hasRefresh) {
          try {
            await this.refreshAuthToken();
          } catch {
            this.clearToken();
            throw new Error("Срок действия токена истек");
          }
        } else {
          this.clearToken();
          throw new Error("Срок действия токена истек");
        }
      }
    },

    async refreshAuthToken(): Promise<boolean> {
      const refreshToken = useCookie('refresh_token').value;

      if (!refreshToken) {
        this.clearToken();
        throw new Error("Невозомжно обновить токен. Срок действия токена истек");
      }

      try {
        const { $api } = useNuxtApp();
        const response = await $api.post<{
          access_token: string;
          expires_in: number;
          refresh_token: string | undefined;
        }>(
          "/auth/refresh",
          {},
          getRefreshTokenHeaders(refreshToken),
        );

        if (response && response.access_token) {
          this.setToken(
            response.access_token,
            response.refresh_token,
            response.expires_in,
          );
          return true;
        } else {
          throw new Error("Неверный формат ответа при обновлении токена");
        }
      } catch (error: any) {
        this.clearToken();
        throw new Error("Произошла ошибка при обновлении токена");
      }
    },

    async login(login: string, password: string): Promise<boolean> {
      try {
        const { $api } = useNuxtApp();
        const response = await $api.post<{
          access_token: string;
          expires_in: number;
          refresh_token: string | undefined;
        }>("/auth/sign-in", { login: login, password: password });

        if (response && response.access_token) {
          this.setToken(
            response.access_token,
            response.refresh_token,
            response.expires_in,
          );
          return true;
        } else {
          throw new Error("Неверный формат ответа от сервера");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при входе в систему");
      }
    },

    async logout() {
      const token = useCookie('auth_token').value;
      const refreshToken = useCookie('refresh_token').value;

      if (!token || !refreshToken) {
        this.user = null;
        this.clearToken();
        return;
      }

      try {
        const { $api } = useNuxtApp();
        await $api.post<{ message: string }>(
          "/auth/sign-out",
          undefined,
          getLogoutTokenHeaders(token, refreshToken),
        );
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при выходе из системы");
      } finally {
        this.user = null;
        this.clearToken();
      }
    },

    async fetchUserInfo(): Promise<User> {
      try {
        await this.ensureValidToken();
      } catch {
        this.clearToken();
        throw new Error(
          "Не удалось получить данные пользователя. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token) {
        this.clearToken();
        throw new Error(
          "Не удалось получить данные пользователя. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();
        const userData = await $api.get<User>(
          "/users/profile",
          getTokenHeaders(token),
        );

        if (userData) {
          this.user = userData;
          return userData;
        } else {
          throw new Error("Не удалось получить данные пользователя");
        }
      } catch (error: any) {
        this.clearToken();
        throw new Error(
          error.message || "Ошибка при получении данных пользователя",
        );
      }
    },

    async updateUserInfo(userData: Partial<User>): Promise<boolean> {
      try {
        await this.ensureValidToken();
      } catch {
        this.clearToken();
        throw new Error(
          "Не удалось получить данные пользователя. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token || !this.user) {
        this.clearToken();
        throw new Error(
          "Не удалось обновить данные пользователя. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();

        const dataToSend = { ...userData };

        const updatedUser = await $api.put<defaultResponse>(
          `/users/${this.user.user_id}`,
          dataToSend,
          getTokenHeaders(token),
        );

        if (updatedUser) {
          this.user = {
            ...this.user,
            ...dataToSend,
          };
          return true;
        } else {
          throw new Error("Не удалось обновить данные пользователя");
        }
      } catch (error: any) {
        throw new Error(
          error.message || "Ошибка при обновлении данных пользователя",
        );
      }
    },

    async changePassword(oldPassword: string, newPassword: string): Promise<boolean> {
      try {
        await this.ensureValidToken();
      } catch {
        this.clearToken();
        throw new Error(
          "Не удалось получить данные пользователя. Пользователь неавторизован",
        );
      }

      const token = useCookie('auth_token').value;

      if (!token || !this.user) {
        this.clearToken();
        throw new Error(
          "Не удалось изменить пароль. Пользователь неавторизован",
        );
      }

      try {
        const { $api } = useNuxtApp();
        const dataToSend = { old_password: oldPassword, new_password: newPassword };
        const response = await $api.put<defaultResponse>(
          "/users/profile/pass",
          dataToSend,
          getTokenHeaders(token),
        );

        if (response) {
          return true;
        } else {
          throw new Error("Неверный формат ответа от сервера");
        }
      } catch (error: any) {
        throw new Error(error.message || "Ошибка при изменении пароля");
      }
    }
  },
});
