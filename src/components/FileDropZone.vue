<template>
  <div
    class="drop-zone"
    :class="{ 'drop-zone--over': isDragging }"
    @click="pick"
    @dragover.prevent="isDragging = true"
    @dragleave.prevent="isDragging = false"
    @drop.prevent="onDrop"
  >
    <div class="drop-inner">
      <div class="drop-icon-wrap" :class="{ 'drop-icon-wrap--over': isDragging }">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
        </svg>
      </div>
      <div class="drop-text-block">
        <p class="drop-title">{{ isDragging ? 'Drop your note here' : 'Drag & drop your note file here' }}</p>
        <p class="drop-sub">Click or drag a <span class="mono accent">.md</span> file to upload</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { cole } from '@/lib/cole.js'

const emit = defineEmits(['file-selected'])
const isDragging = ref(false)

async function pick() {
  const filePath = await cole.openFileDialog()
  if (!filePath) return
  const info = await cole.readFile(filePath)
  emit('file-selected', info)
}

async function onDrop(e) {
  isDragging.value = false
  const file = e.dataTransfer.files[0]
  if (!file || !file.name.endsWith('.md')) return
  if (!file.path) return
  const info = await cole.readFile(file.path)
  emit('file-selected', info)
}
</script>

<style scoped>
.drop-zone {
  border: 2px dashed var(--color-border);
  border-radius: 16px;
  padding: 32px 24px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;
  background: var(--color-surface);
  box-shadow: var(--shadow-card);
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.drop-zone:hover {
  border-color: var(--color-primary);
  background: rgba(0,122,255,0.03);
}
.drop-zone--over {
  border-color: var(--color-primary);
  background: rgba(0,122,255,0.06);
  box-shadow: 0 0 0 4px rgba(0,122,255,0.1);
}

.drop-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
}

.drop-icon-wrap {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: var(--color-surface-2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-muted);
  transition: background 0.15s, color 0.15s;
}
.drop-icon-wrap svg { width: 26px; height: 26px; }
.drop-icon-wrap--over {
  background: rgba(0,122,255,0.12);
  color: var(--color-primary);
}

.drop-text-block { display: flex; flex-direction: column; gap: 4px; }
.drop-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text);
  margin: 0;
}
.drop-sub {
  font-size: 11px;
  color: var(--color-text-muted);
  margin: 0;
}
.mono { font-family: 'JetBrains Mono', Consolas, monospace; }
.accent { color: var(--color-primary); }
</style>
