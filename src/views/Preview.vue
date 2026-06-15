<template>
  <div class="flex flex-col h-full">
    <!-- Toolbar -->
    <div class="flex items-center justify-between px-6 py-3 border-b border-j-border bg-j-surface flex-shrink-0">
      <div class="flex items-center gap-3">
        <button @click="$router.push('/')" class="text-j-muted hover:text-white transition-colors p-1 rounded">
          <svg viewBox="0 0 24 24" class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>
        <div>
          <span class="text-sm font-medium text-white">{{ report.fileName || 'No file loaded' }}</span>
          <span v-if="report.dateStart" class="text-xs text-j-muted ml-2">
            {{ report.dateStart }} → {{ report.dateEnd }}
          </span>
          <span v-if="report.reportType" class="ml-2 text-[10px] font-bold bg-j-accent/20 text-j-glow border border-j-accent/30 px-1.5 py-0.5 rounded">
            {{ report.reportType }}
          </span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <!-- Generate .docx button — only when parsed data exists -->
        <button
          v-if="report.parsedData"
          @click="generateDocx"
          :disabled="generating"
          class="flex items-center gap-1.5 px-3 py-1.5 bg-j-accent hover:bg-blue-500 text-white text-xs font-semibold rounded-lg transition-colors shadow-glow-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg v-if="generating" class="w-3.5 h-3.5 animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          {{ generating ? 'Generating…' : 'Export .docx' }}
        </button>

        <div class="flex rounded-lg border border-j-border overflow-hidden">
          <button
            v-for="tab in availableTabs"
            :key="tab"
            @click="activeTab = tab"
            class="px-3 py-1.5 text-xs font-medium transition-colors"
            :class="activeTab === tab ? 'bg-j-accent text-white' : 'bg-j-card text-j-muted hover:text-white'"
          >
            {{ tab }}
          </button>
        </div>
      </div>
    </div>

    <!-- Export error — always visible regardless of active tab -->
    <div v-if="genError" class="flex items-start gap-2.5 px-6 py-3 bg-red-950/60 border-b border-j-error/40 flex-shrink-0">
      <svg viewBox="0 0 24 24" class="w-4 h-4 text-j-error flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      <p class="text-xs text-red-300 leading-relaxed whitespace-pre-wrap">{{ genError }}</p>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-auto p-6">
      <!-- No file -->
      <div v-if="!report.rawContent" class="flex flex-col items-center justify-center h-full text-center">
        <svg viewBox="0 0 24 24" class="w-12 h-12 text-j-border mb-4" fill="none" stroke="currentColor" stroke-width="1">
          <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
        </svg>
        <p class="text-j-muted text-sm">No file loaded.</p>
        <button @click="$router.push('/')" class="mt-3 text-xs text-j-accent hover:text-j-glow transition-colors">← Back to Generate</button>
      </div>

      <!-- Parsed Data tab -->
      <div v-else-if="activeTab === 'Parsed Data'">
        <div v-if="!report.parsedData" class="flex flex-col items-center justify-center h-48 text-center">
          <p class="text-j-muted text-sm">No parsed data yet.</p>
          <button @click="$router.push('/')" class="mt-2 text-xs text-j-accent hover:text-j-glow transition-colors">← Go back to parse</button>
        </div>
        <div v-else class="max-w-3xl mx-auto space-y-5">
          <!-- Meta / summary -->
          <div class="bg-j-card border border-j-border rounded-xl p-4">
            <div class="flex items-center justify-between mb-2">
              <span class="text-xs font-semibold uppercase tracking-wider text-j-muted">
                {{ report.parsedData.reportType }} · {{ report.parsedData.dateStart }} → {{ report.parsedData.dateEnd }}
              </span>
              <span v-if="report.parsedData.totalHours" class="text-xs font-mono text-j-glow">
                Total: {{ report.parsedData.totalHours }}h
              </span>
            </div>
            <p v-if="report.parsedData.summary" class="text-sm text-j-text leading-relaxed">
              {{ report.parsedData.summary }}
            </p>
          </div>

          <!-- Task table -->
          <TaskTable v-if="report.parsedData.entries?.length" :rows="report.parsedData.entries" />

          <!-- Ongoing tasks (PRG) -->
          <div v-if="report.parsedData.ongoingTasks?.length" class="space-y-2">
            <h3 class="text-xs font-semibold uppercase tracking-wider text-j-muted mb-2">Ongoing Tasks</h3>
            <div v-for="task in report.parsedData.ongoingTasks" :key="task.name"
                 class="bg-j-card border border-j-border rounded-xl p-4">
              <div class="flex items-center justify-between mb-2">
                <span class="text-sm font-medium text-white">{{ task.name }}</span>
                <span class="text-xs font-mono text-j-glow">{{ task.percentComplete }}%</span>
              </div>
              <div class="w-full bg-j-surface rounded-full h-1.5 mb-3">
                <div class="bg-j-accent h-1.5 rounded-full" :style="{ width: task.percentComplete + '%' }" />
              </div>
              <ul v-if="task.nextSteps?.length" class="space-y-0.5">
                <li v-for="step in task.nextSteps" :key="step" class="flex items-start gap-2 text-xs text-j-text">
                  <span class="mt-1.5 w-1 h-1 rounded-full bg-j-accent flex-shrink-0" />
                  {{ step }}
                </li>
              </ul>
            </div>
          </div>

          <!-- Raw JSON toggle -->
          <details>
            <summary class="text-xs text-j-muted cursor-pointer hover:text-white transition-colors select-none">Raw JSON</summary>
            <pre class="mt-2 selectable text-xs font-mono text-j-muted bg-j-card border border-j-border rounded-xl p-4 overflow-auto">{{ JSON.stringify(report.parsedData, null, 2) }}</pre>
          </details>

        </div>
      </div>

      <!-- Rendered markdown -->
      <div v-else-if="activeTab === 'Rendered'"
           class="md-preview selectable max-w-3xl mx-auto"
           v-html="rendered" />

      <!-- Raw markdown -->
      <pre v-else class="selectable text-xs font-mono text-j-muted whitespace-pre-wrap leading-relaxed max-w-3xl mx-auto">{{ report.rawContent }}</pre>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, toRaw } from 'vue'
import { marked } from 'marked'
import { useReportStore } from '@/stores/report.js'
import { useSettingsStore } from '@/stores/settings.js'
import TaskTable from '@/components/TaskTable.vue'

const report    = useReportStore()
const settings  = useSettingsStore()
const generating = ref(false)
const genError   = ref('')

const activeTab = ref(report.parsedData ? 'Parsed Data' : 'Rendered')

const availableTabs = computed(() => {
  const tabs = ['Rendered', 'Raw']
  if (report.parsedData) tabs.unshift('Parsed Data')
  return tabs
})

const rendered = computed(() =>
  report.rawContent ? marked.parse(report.rawContent) : ''
)

async function generateDocx() {
  if (!report.parsedData) return
  generating.value = true
  genError.value   = ''
  try {
    await settings.load()
    const userInfo = {
      reporterName:    settings.reporterName,
      reporterNameSig: settings.reporterNameSig,
      position:        settings.position,
      office:          settings.office,
      supervisorName:  settings.supervisorName,
      supervisorPos1:  settings.supervisorPos1,
      supervisorPos2:  settings.supervisorPos2,
    }

    const d = toRaw(report.parsedData)
    const defaultName = `Accomplishment Report ${d.dateStart} to ${d.dateEnd}.docx`
    const savePath = await window.cole.showSaveDialog({
      defaultPath: defaultName,
      filters: [{ name: 'Word Document', extensions: ['docx'] }],
    })
    if (!savePath) return

    await window.cole.generateReport({ parsedData: d, userInfo, outputPath: savePath })
    await window.cole.openPath(savePath)
  } catch (err) {
    genError.value = err.message || 'Failed to generate report.'
  } finally {
    generating.value = false
  }
}
</script>
