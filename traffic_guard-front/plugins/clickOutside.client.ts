import { defineNuxtPlugin } from "#app";
import clickOutside from "~/directives/clickOutside";

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.directive("click-outside", clickOutside);
});
