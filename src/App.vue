<template>
  <div class="flex h-screen bg-j-bg text-j-text overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-52 flex-shrink-0 bg-j-surface border-r border-j-border flex flex-col">
      <!-- Title bar region (leaves space for traffic lights / overlay buttons) -->
      <div class="h-9 app-drag flex items-center px-4 border-b border-j-border">
        <span class="text-xs font-mono text-j-muted tracking-widest app-no-drag">COLE</span>
      </div>

      <!-- Logo block -->
      <div class="p-4 border-b border-j-border flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-j-accent flex items-center justify-center shadow-glow">
          <svg viewBox="0 0 24 24" class="w-4 h-4 text-white" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 2v3M12 19v3M4.22 4.22l2.12 2.12M17.66 17.66l2.12 2.12M2 12h3M19 12h3M4.22 19.78l2.12-2.12M17.66 6.34l2.12-2.12"/>
          </svg>
        </div>
        <div>
          <div class="text-sm font-semibold text-white leading-tight">Report Assistant</div>
          <div class="text-xs text-j-muted">v0.1 · Phase 1</div>
        </div>
      </div>

      <!-- Nav -->
      <nav class="flex-1 p-2 space-y-0.5">
        <router-link
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-j-muted hover:text-white hover:bg-j-card transition-colors"
          active-class="!text-white bg-j-card border border-j-border shadow-glow-sm"
        >
          <component :is="item.icon" class="w-4 h-4 flex-shrink-0" />
          {{ item.label }}
        </router-link>
      </nav>

      <!-- Report store status -->
      <div v-if="report.fileName" class="m-2 p-2 bg-j-card border border-j-border rounded-lg">
        <div class="text-xs text-j-muted mb-0.5">Active file</div>
        <div class="text-xs text-j-glow font-mono truncate">{{ report.fileName }}</div>
        <div v-if="report.dateStart" class="text-xs text-j-muted mt-0.5">
          {{ report.dateStart }} → {{ report.dateEnd }}
        </div>
      </div>

      <div class="p-3 text-xs text-j-muted border-t border-j-border">
        Obsidian → Claude → .docx
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-auto">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { h, onMounted } from 'vue'
import { useReportStore } from '@/stores/report.js'
import { useSettingsStore } from '@/stores/settings.js'

const report   = useReportStore()
const settings = useSettingsStore()

onMounted(() => settings.load())

// Inline SVG icon components
const IconHome = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('path', { d:'M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z' }),
  h('polyline', { points:'9 22 9 12 15 12 15 22' }),
]) }

const IconEye = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('path', { d:'M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z' }),
  h('circle', { cx:'12', cy:'12', r:'3' }),
]) }

const IconSettings = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('circle', { cx:'12', cy:'12', r:'3' }),
  h('path', { d:'M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z' }),
]) }

const nav = [
  { to: '/',         label: 'Generate',    icon: IconHome },
  { to: '/preview',  label: 'MD Preview',  icon: IconEye },
  { to: '/settings', label: 'Settings',    icon: IconSettings },
]
</script>
