<template>
  <button
    :type="type"
    class="base-button"
    :class="[`base-button--${variant}`, `base-button--${extras}`, { 'base-button--loading': loading }]"
    :disabled="disabled || loading"
  >
    <span v-if="loading">{{ loadingText }}</span>
    <template v-else>
      <span v-if="$slots.icon" class="base-button__icon">
        <slot name="icon" />
      </span>
      <span class="base-button__text">
        <slot />
      </span>
    </template>
  </button>
</template>

<script setup lang="ts">
defineProps({
  type: {
    type: String as () => 'button' | 'submit' | 'reset',
    default: 'button'
  },
  variant: {
    type: String as () => 'primary' | 'secondary' | 'extra' | 'danger' | 'login',
    default: 'primary'
  },
  extras: {
    type: String as () => 'default' | 'small',
    default: 'default'
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  loadingText: {
    type: String,
    default: 'Загрузка...'
  }
});
</script>

<style>
.base-button {
  width: 100%;
  max-height: 2.25rem;
  display: flex;
  justify-content: center;
  align-items: center;
  border-width: 0;
  transition: all 0.2s ease;
  cursor: pointer;
  border-radius: 0.375rem;
  box-shadow: 0px 1px 2px 0px #0000000D;
  padding-inline: 0.75rem;
}

.base-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.base-button--primary {
  padding-block: 0.5rem;
  font-size: 1rem;
  line-height: 1.25rem;
  font-weight: 600;
  background-color: #2563EB;
  color: #FFFFFF;
}

.base-button--primary:hover:not(:disabled) {
  filter: brightness(85%);
}

.base-button--secondary {
  outline: 1px solid #E4E4E7;
  padding-block: 0.5rem;
  font-size: 1rem;
  line-height: 1.25rem;
  font-weight: 500;
  background-color: #FFFFFF;
  color: #101828;
}

.base-button--login {
  line-height: 1.25rem;
  padding-block: 0.375rem;
  border-radius: 999px;
  color: #FFFFFF;
  background-color: #2563EB;

}

.base-button--extra {
  padding-block: 0.5rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 600;
  background-color: #EFF6FF;
  color: #2563EB;
}

.base-button--danger {
  padding-block: 0.5rem;
  font-size: 0.875rem;
  line-height: 1.25rem;
  font-weight: 600;
  background-color: #FFE2E2;
  color: #C10007;
}


.base-button:focus {
  outline: none;
  box-shadow: 0 0 0 0.125rem #cbd5e1, 0 0 0 0.25rem rgba(255, 255, 255, 0.8);
}

.base-button--small {
  font-size: 0.875rem;
}

.base-button--secondary.base-button--small {
  outline: 1px solid #D4D4D8;
}

.base-button__icon {
  display: flex;
  justify-content: center;
  align-items: center;
  color: var(--color-icon-primary);
  width: 1.25rem;
  height: 1.25rem;
  margin-right: 0.375rem;
}

.base-button--primary .base-button__icon {
  color: #FFFFFF;
}
</style>