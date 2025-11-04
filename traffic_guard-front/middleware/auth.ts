import { defineNuxtRouteMiddleware, navigateTo } from "#app";
import { useCookie } from "nuxt/app";

export default defineNuxtRouteMiddleware(async (to, from) => {
  const token = useCookie('auth_token');
  const refreshToken = useCookie('refresh_token')

  if (!token.value && !refreshToken.value) {
    return navigateTo('/auth');
  }
});
