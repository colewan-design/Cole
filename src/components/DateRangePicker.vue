<template>
  <div class="grid grid-cols-2 gap-3">
    <div>
      <label class="block text-xs text-j-muted mb-1.5 font-medium">From</label>
      <input
        type="date"
        :value="modelStart"
        @input="$emit('update:modelStart', $event.target.value)"
        class="w-full bg-j-card border border-j-border rounded-lg px-3 py-2 text-sm text-j-text
               focus:outline-none focus:border-j-accent focus:ring-1 focus:ring-j-accent/30
               [color-scheme:dark] cursor-pointer"
      />
    </div>
    <div>
      <label class="block text-xs text-j-muted mb-1.5 font-medium">To</label>
      <input
        type="date"
        :value="modelEnd"
        @input="$emit('update:modelEnd', $event.target.value)"
        class="w-full bg-j-card border border-j-border rounded-lg px-3 py-2 text-sm text-j-text
               focus:outline-none focus:border-j-accent focus:ring-1 focus:ring-j-accent/30
               [color-scheme:dark] cursor-pointer"
      />
    </div>
  </div>
</template>

<script setup>
defineProps({ modelStart: String, modelEnd: String })
defineEmits(['update:modelStart', 'update:modelEnd'])
</script>
