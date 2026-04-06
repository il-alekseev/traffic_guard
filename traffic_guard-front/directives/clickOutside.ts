import type { ObjectDirective, DirectiveBinding } from "vue";

interface ClickOutsideValue {
  handler: (event: Event) => void;
  ignore?: string | string[];
}

interface ClickOutsideElement extends HTMLElement {
  __clickOutsideHandler__?: (event: MouseEvent) => void;
}

function matchesIgnore(target: Node, ignore?: string | string[]) {
  if (!ignore) return false;

  const selectors = Array.isArray(ignore) ? ignore : [ignore];

  return selectors.some((selector) => {
    const el = document.querySelector(selector);
    return el && (el === target || el.contains(target as Node));
  });
}

const clickOutside: ObjectDirective = {
  beforeMount(el: ClickOutsideElement, binding: DirectiveBinding<ClickOutsideValue>) {
    el.__clickOutsideHandler__ = (event: MouseEvent) => {
      const value = binding.value;

      if (!value || typeof value.handler !== "function") {
        console.warn("v-click-outside requires { handler: Function }");
        return;
      }

      const target = event.target as Node;

      if (matchesIgnore(target, value.ignore)) {
        return;
      }

      if (el === target || el.contains(target)) {
        return;
      }

      value.handler(event);
    };

    document.addEventListener("click", el.__clickOutsideHandler__);
  },

  unmounted(el: ClickOutsideElement) {
    if (el.__clickOutsideHandler__) {
      document.removeEventListener("click", el.__clickOutsideHandler__);
      delete el.__clickOutsideHandler__;
    }
  }
};

export default clickOutside;
