<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal" @click.stop>
      <div v-if="$slots.image" class="modal__image">
        <slot name="image" />
      </div>

      <header v-if="$slots.header" class="modal__header">
        <slot name="header" />
      </header>

      <main v-if="$slots.main" class="modal__content">
        <slot name="main"/>
      </main>

      <div v-if="$slots.error" class="modal__error">
        <slot name="error"/>
      </div>

      <footer v-if="$slots.footer" class="modal__footer">
        <slot name="footer" />
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits(['close']);

const closeModal = () => {
  emit('close');
};
</script>

<style>
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(90, 90, 90, 0.3);
  backdrop-filter: blur(1.5px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999999;
}

.modal {
  background-color: white;
  width: 100%;
  max-width: 24rem;
  padding: 1.5rem;
  box-shadow: 0px 8px 10px -6px #0000001A;
  box-shadow: 0px 20px 25px -5px #0000001A;

  border-radius: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.modal__header {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  align-items: center;
  text-align: center;
}

.modal__header-title-block {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.modal__image {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 4.75rem;
  height: 4.75rem;
  color: #002033;
  align-self: center;
}

.modal__footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.modal__title {
  font-size: 1rem;
  line-height: 1.5rem;
  font-weight: 600;
  color: #101828;
  margin: 0;
}

.modal__text {
  font-weight: 400;
  font-size: 0.875rem;
  line-height: 1.25rem;
  color: #6A7282;
}
</style>
