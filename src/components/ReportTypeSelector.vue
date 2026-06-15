<template>
  <div class="grid grid-cols-3 gap-2">
    <button
      v-for="type in types"
      :key="type.id"
      @click="$emit('update:modelValue', type.id)"
      class="relative flex flex-col items-center gap-1.5 p-3 rounded-xl border transition-all duration-150 text-center"
      :class="[
        modelValue === type.id
          ? 'border-j-accent bg-j-accent/10 shadow-glow-sm text-white'
          : 'border-j-border bg-j-card text-j-muted hover:border-j-accent/40 hover:text-j-text',
      ]"
    >
      <span class="text-lg">{{ type.emoji }}</span>
      <span class="text-xs font-semibold tracking-wide">{{ type.id }}</span>
      <span class="text-[10px] leading-tight opacity-75">{{ type.label }}</span>
      <div v-if="modelValue === type.id"
           class="absolute top-1.5 right-1.5 w-1.5 h-1.5 rounded-full bg-j-accent" />
    </button>
  </div>
</template>

<script setup>
defineProps({ modelValue: String })
defineEmits(['update:modelValue'])

const types = [
  { id: 'AR',  emoji: '📋', label: 'Accomplishment Report' },
  { id: 'DTR', emoji: '🕐', label: 'Daily Time Record' },
  { id: 'PRG', emoji: '📈', label: 'Progress Report' },
]
</script>
