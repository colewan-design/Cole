<template>
  <div
    class="relative border-2 border-dashed rounded-xl p-8 text-center transition-all duration-200 cursor-pointer"
    :class="[
      isDragging
        ? 'border-j-accent bg-j-accent/10 shadow-glow'
        : 'border-j-border bg-j-card hover:border-j-accent/50 hover:bg-j-card/80',
    ]"
    @click="pick"
    @dragover.prevent="isDragging = true"
    @dragleave.prevent="isDragging = false"
    @drop.prevent="onDrop"
  >
    <div class="flex flex-col items-center gap-3">
      <div class="w-12 h-12 rounded-full flex items-center justify-center"
           :class="isDragging ? 'bg-j-accent/20' : 'bg-j-surface'">
        <svg viewBox="0 0 24 24" class="w-6 h-6" :class="isDragging ? 'text-j-glow' : 'text-j-muted'"
             fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
        </svg>
      </div>
      <div>
        <p class="text-sm font-medium" :class="isDragging ? 'text-j-glow' : 'text-j-text'">
          {{ isDragging ? 'Drop your note here' : 'Pick your Obsidian note' }}
        </p>
        <p class="text-xs text-j-muted mt-1">Click or drag a <span class="font-mono text-j-glow">.md</span> file</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['file-selected'])
const isDragging = ref(false)

async function pick() {
  const filePath = await window.cole.openFileDialog()
  if (!filePath) return
  const info = await window.cole.readFile(filePath)
  emit('file-selected', info)
}

async function onDrop(e) {
  isDragging.value = false
  const file = e.dataTransfer.files[0]
  if (!file || !file.name.endsWith('.md')) return
  const info = await window.cole.readFile(file.path)
  emit('file-selected', info)
}
</script>
