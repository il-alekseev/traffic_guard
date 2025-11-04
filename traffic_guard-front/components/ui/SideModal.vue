<template>
  <transition name="slide-fade">
    <div v-if="modelValue" class="side-modal">
      <div class="side-modal__overlay">
        <div class="side-modal__overlay-blur"></div>
      </div>
      
      <div class="side-modal__content">
        <div class="side-modal__header">
          <div class="side-modal__title-container">
            <h2 class="side-modal__title">{{ title }}</h2>
            <span class="side-modal__steps" v-if="maxStep !== '' && maxStep !== undefined && currentStep !== '' && currentStep !== undefined">{{ currentStep }}/{{ maxStep }}</span>
          </div>
          
          <CrossIcon class="side-modal__close" @click="emit('close')" />
        </div>

        <div class="side-modal__body">
          <slot />
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import CrossIcon from '@/assets/img/cross.svg'

defineProps<{
  modelValue: boolean
  title: string
  currentStep?: string;
  maxStep?: string;
}>()

const emit = defineEmits(['close']);
</script>

<style scoped>
.side-modal {
  position: fixed;
  inset: 0;
  z-index: 150;
  display: flex;
  justify-content: flex-end;
}
.side-modal__overlay {
  position: absolute;
  inset: 0;
  background: transparent;
}
.side-modal__overlay-blur {
  z-index: 9999;
  position: absolute;
  top: 0;
  width: calc(100vw - 18rem);
  margin-left: 18rem;
  height: calc(100vh - 4rem);
  margin-top: 4rem;
  background: #D7D7D766;
  backdrop-filter: blur(10px);
  border-top-left-radius: 1.75rem;
}
.side-modal__content {
  z-index: 99999;
  position: relative;
  width: 28rem;
  height: calc(100vh - 4rem);
  margin-top: auto;
  background: #FFFFFF;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0px 20px 25px -5px #0000001A;
}
.side-modal__header {
  padding: 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.side-modal__title-container {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}
.side-modal__title {
  font-weight: 500;
  font-size: 1rem;
  line-height: 1.5rem;
  color: #3F3F46;
}
.side-modal__steps {
  font-weight: 600;
  font-size: var(--font-size-3xl);
  line-height: 2.125rem;
  color: var(--color-typo-placeholder);
}
.side-modal__body {
  overflow-y: auto;
  flex: 1;
  height: 100%;
}
.side-modal__close {
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  width: 1.5rem;
  height: 1.5rem;
  color: #A1A1AA;
}
.slide-fade-enter-active {
  transition: all 0.3s ease;
}
.slide-fade-leave-active {
  transition: all 0.3s ease;
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(100%);
  opacity: 0;
}
</style>
