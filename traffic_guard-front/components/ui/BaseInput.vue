<template>
  <div class="base-form-field">
    <label :for="id" class="base-label">{{ label }} <span class="requires-msg" v-if="props.required === true && props.showRequiredLabel">Обязательное поле</span></label>
    <div v-if="type === 'password'" class="base-password-container">
      <input
        :required="props.required"
        :id="id"
        :name="name"
        :type="showPassword ? 'text' : 'password'"
        :value="modelValue"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        class="base-input"
        :class="{ 'base-input--error': error, 'base-input--disabled': props.disabled}"
        :placeholder="placeholder"
        :autocomplete="autocomplete"
        :disabled="props.disabled"
      />
      <button
        type="button"
        class="base-password-toggle"
        @click="togglePassword"
        tabindex="-1"
      >
        <EyeIcon v-if="!showPassword" alt="Показать"/>
        <EyeOffIcon v-else alt="Скрыть"/>
      </button>
    </div>
    <input
      v-else
      :required="props.required"
      :id="id"
      :name="name"
      :type="type"
      :value="modelValue"
      @input="$emit(
        'update:modelValue',
        type === 'number'
          ? Number(($event.target as HTMLInputElement).value)
          : ($event.target as HTMLInputElement).value
      )"
      class="base-input"
      :class="{ 'base-input--error': error, 'base-input--disabled': props.disabled}"
      :placeholder="placeholder"
      :autocomplete="autocomplete"
      :disabled="props.disabled"
    />
    <p v-if="props.extraInfo" class="base-form-hint-text">
      {{ props.extraInfo }}
    </p>
    <p v-if="error && showTextError" class="base-form-error-text">
      {{ error }}
    </p>
  </div>
</template>

<script setup lang="ts">
import EyeIcon from '~/assets/img/eye.svg';
import EyeOffIcon from '~/assets/img/eye-off.svg';

const props = defineProps({
  id: {
    type: String,
    required: true
  },
  name: {
    type: String,
    required: true
  },
  label: {
    type: String,
    required: true
  },
  showRequiredLabel: {
    type: Boolean,
    default: false
  },
  modelValue: {
    type: [String, Number],
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  placeholder: {
    type: String,
    default: ''
  },
  error: {
    type: String,
    default: ''
  },
  showTextError: {
    type: Boolean,
    default: true
  },
  disabled: {
    type: Boolean,
    default: false
  },
  extraInfo: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean, 
    default: true,
  },
  autocomplete: {
    type: String,
    default: ''
  }
});

defineEmits(['update:modelValue']);

const showPassword = ref(false);

const togglePassword = () => {
  showPassword.value = !showPassword.value;
};
</script>

<style>
.base-form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.base-label {
  font-size: 1rem;
  font-weight: 500;
  line-height: 1.5rem;
  color: #3F3F46;
  display: flex;
  gap: 0.375rem;
}

.base-input {
  width: 100%;
  padding: 0.5rem;
  font-size: 1rem;
  line-height: 1rem;
  color: #3F3F46;
  background-color: #FFFFFF;
  outline: 1px solid #D4D4D8;
  border-radius: 0.3125rem;
  border: 0;
  transition: all 0.2s ease;
}

.base-input:focus {
  box-shadow: 0 0 0 0.125rem #cbd5e1, 0 0 0 0.25rem rgba(255, 255, 255, 0.8);
}

.base-input--error {
  outline-color: #EF4444;
}

.base-input--disabled {
  background-color: #EDF1F2;
}

.base-input--error:focus {
  box-shadow: 0 0 0 0.125rem rgba(239, 68, 68, 0.3);
}

.base-input::placeholder {
  color: #848F93;
}

.base-password-container {
  position: relative;
}

.base-password-toggle {
  position: absolute;
  right: 0.625rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 1rem;
  height: 1rem;
}

.base-form-error-text {
  font-size: 1rem;
  color: #EF4444;
}

.base-form-hint-text {
  font-size: 0.75rem;
  line-height: 0.75rem;
  color: #71717A;
}

.requires-msg {
  color: #4A5565;
  font-size: 0.75rem;
  line-height: 1rem;
  font-weight: 400;
  padding: 0.25rem 0.5rem;
  border-radius: 33554400px;
  background-color: #F3F4F6;
}
</style>