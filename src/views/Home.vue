<template>
  <div class="p-6 max-w-2xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-xl font-semibold text-white">Generate Report</h1>
      <p class="text-sm text-j-muted mt-1">Pick an Obsidian note, set a date range, choose a report type.</p>
    </div>

    <!-- Step 1 — File -->
    <section class="mb-5">
      <div class="flex items-center gap-2 mb-3">
        <StepBadge :n="1" :done="!!report.filePath" />
        <span class="text-sm font-medium text-j-text">Select Markdown file</span>
      </div>

      <FileDropZone v-if="!report.filePath" @file-selected="onFile" />

      <div v-else class="flex items-center gap-3 bg-j-card border border-j-border rounded-xl px-4 py-3">
        <div class="w-8 h-8 rounded-lg bg-j-accent/10 flex items-center justify-center">
          <svg viewBox="0 0 24 24" class="w-4 h-4 text-j-glow" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
          </svg>
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium text-white truncate">{{ report.fileName }}</div>
          <div class="text-xs text-j-muted truncate font-mono">{{ report.filePath }}</div>
        </div>
        <button @click="report.reset()" class="text-j-muted hover:text-j-error transition-colors p-1 rounded">
          <svg viewBox="0 0 24 24" class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
    </section>

    <!-- Step 2 — Date range -->
    <section class="mb-5">
      <div class="flex items-center gap-2 mb-3">
        <StepBadge :n="2" :done="!!(report.dateStart && report.dateEnd)" />
        <span class="text-sm font-medium text-j-text">Date range</span>
      </div>
      <DateRangePicker
        v-model:modelStart="dateStart"
        v-model:modelEnd="dateEnd"
        @update:modelStart="onDateChange"
        @update:modelEnd="onDateChange"
      />
    </section>

    <!-- Step 3 — Report type -->
    <section class="mb-6">
      <div class="flex items-center gap-2 mb-3">
        <StepBadge :n="3" :done="true" />
        <span class="text-sm font-medium text-j-text">Report type</span>
      </div>
      <ReportTypeSelector v-model="reportType" @update:modelValue="report.setReportType($event)" />
    </section>

    <!-- Actions -->
    <div class="flex gap-3">
      <button
        @click="previewMd"
        :disabled="!report.filePath"
        class="flex items-center gap-2 px-4 py-2.5 rounded-lg border border-j-border bg-j-card text-sm font-medium
               hover:border-j-accent/50 hover:text-white transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
      >
        <svg viewBox="0 0 24 24" class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
          <circle cx="12" cy="12" r="3"/>
        </svg>
        Preview MD
      </button>

      <button
        :disabled="!canParse"
        class="relative flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg text-sm font-semibold transition-all
               bg-j-accent hover:bg-blue-500 text-white shadow-glow disabled:opacity-40 disabled:cursor-not-allowed disabled:shadow-none"
        @click="parse"
      >
        <svg v-if="parsing" class="w-4 h-4 animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/>
          <path d="M12 2v3M12 19v3M4.22 4.22l2.12 2.12M17.66 17.66l2.12 2.12M2 12h3M19 12h3M4.22 19.78l2.12-2.12M17.66 6.34l2.12-2.12"/>
        </svg>
        {{ parsing ? 'Parsing…' : settings.provider === 'ollama' ? `Parse with ${settings.ollamaModel}` : 'Parse with Claude' }}
      </button>
    </div>

    <!-- Error -->
    <div v-if="parseError" class="mt-4 flex items-start gap-2.5 p-3 bg-red-950/40 border border-j-error/40 rounded-lg">
      <svg viewBox="0 0 24 24" class="w-4 h-4 text-j-error flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      <p class="text-xs text-red-300 leading-relaxed">{{ parseError }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, h, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useReportStore } from '@/stores/report.js'
import { useSettingsStore } from '@/stores/settings.js'
import FileDropZone from '@/components/FileDropZone.vue'
import DateRangePicker from '@/components/DateRangePicker.vue'
import ReportTypeSelector from '@/components/ReportTypeSelector.vue'

const router   = useRouter()
const report   = useReportStore()
const settings = useSettingsStore()

onMounted(() => settings.load())

const dateStart  = ref(report.dateStart)
const dateEnd    = ref(report.dateEnd)
const reportType = ref(report.reportType)

const parsing   = ref(false)
const parseError = ref('')

const canParse = computed(() =>
  !!report.filePath && !!dateStart.value && !!dateEnd.value && !parsing.value
)

function onFile(info) {
  report.setFile(info)
}

function onDateChange() {
  report.setDateRange(dateStart.value, dateEnd.value)
}

function previewMd() {
  router.push('/preview')
}

async function parse() {
  if (!canParse.value) return
  parsing.value   = true
  parseError.value = ''
  try {
    const result = await window.cole.parseWithClaude({
      content:    report.rawContent,
      dateStart:  report.dateStart,
      dateEnd:    report.dateEnd,
      reportType: report.reportType,
    })
    report.setParsedData(result)
    router.push('/preview')
  } catch (err) {
    parseError.value = err.message || 'Something went wrong.'
  } finally {
    parsing.value = false
  }
}

// Sync store → local refs when navigating back
watch(() => report.dateStart, v => { dateStart.value = v })
watch(() => report.dateEnd,   v => { dateEnd.value   = v })
watch(() => report.reportType, v => { reportType.value = v })

const StepBadge = {
  props: ['n', 'done'],
  setup(props) {
    return () => h(
      'div',
      { class: `w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0 transition-colors ${props.done ? 'bg-j-accent text-white' : 'bg-j-surface border border-j-border text-j-muted'}` },
      props.done
        ? [h('svg', { viewBox: '0 0 24 24', class: 'w-3 h-3', fill: 'none', stroke: 'currentColor', 'stroke-width': '3' },
            [h('polyline', { points: '20 6 9 17 4 12' })])]
        : [String(props.n)]
    )
  },
}
</script>
