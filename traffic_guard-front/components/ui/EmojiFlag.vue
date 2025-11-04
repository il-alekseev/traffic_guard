<template>
  <span>{{ emojiFlag(code) }}</span>
</template>

<script setup>
const props = defineProps({
  code: { type: String, required: true }
})

const emojiFlag = (code) => {
  return code
    .toUpperCase()
    .replace(/./g, char =>
      String.fromCodePoint(0x1F1E6 + char.charCodeAt(0) - 65)
    )
}
</script>
