<template>
  <button
    class="sub-button"
    :class="`sub-button--${props.color}`"
    :disabled="props.disabled"
  >
    <span class="sub-button__icon">
      <slot name="icon" />
    </span>
    <span class="sub-button__text">
      <slot />
    </span>
  </button>
</template>

<script setup lang="ts">
const props = defineProps({
  disabled: {
    type: Boolean,
    default: false
  },
  color: {
    type: String,
    default: 'white',
    validator: (value: string) => ['white', 'blue'].includes(value),
  }
});
</script>

<style>
.sub-button {
  width: 100%;
  max-height: 2.5rem;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--size-2);
  padding-block: 0.625rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--color-border-primary);
  font-size: var(--font-size-sm);
  line-height: 1.25rem;
  font-weight: 500;
  transition: all var(--transition);
  cursor: pointer;
  color: var(--color-typo-primary);
  background-color: transparent;
}

.sub-button--white {
  color: var(--color-typo-primary);
  background-color: transparent;
}

.sub-button--blue {
  color: var(--color-white);
  background-color: var(--color-blue-main);
  border: 0px;
}

.sub-button:hover:not(:disabled) {
  filter: brightness(85%);
}

.sub-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.sub-button:focus {
  outline: none;
  box-shadow: 0 0 0 0.125rem var(--color-slate-300), 0 0 0 0.25rem rgba(255, 255, 255, 0.8);
}

.sub-button__icon {
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--color-icon-primary);
  width: 1rem;
  height: 1rem;
}

.sub-button--white .sub-button__icon {
  color: var(--color-typo-primary);
}

.sub-button--blue .sub-button__icon {
  color: var(--color-white);
}
</style>