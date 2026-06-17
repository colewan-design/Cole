<template>
  <div class="preview-page">
    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-left">
        <button class="back-btn" @click="$router.push('/')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>
        <div class="file-meta">
          <span class="file-meta-name">{{ report.fileName || 'No file loaded' }}</span>
          <span v-if="report.dateStart" class="file-meta-range">
            {{ report.dateStart }} → {{ report.dateEnd }}
          </span>
          <span v-if="report.reportType" class="type-badge">{{ report.reportType }}</span>
        </div>
      </div>

      <div class="toolbar-right">
        <!-- Tab switcher -->
        <div class="tab-bar">
          <button
            v-for="tab in availableTabs"
            :key="tab"
            @click="activeTab = tab"
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === tab }"
          >{{ tab }}</button>
        </div>

        <button
          v-if="report.parsedData"
          @click="generateDocx"
          :disabled="generating"
          class="export-btn"
        >
          <svg v-if="generating" class="spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>
          </svg>
          {{ generating ? 'Generating…' : 'Export .docx' }}
        </button>
      </div>
    </div>

    <!-- Export error -->
    <div v-if="genError" class="gen-error">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="gen-error-icon">
        <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
      </svg>
      <p class="gen-error-text selectable">{{ genError }}</p>
    </div>

    <!-- Content -->
    <div class="content-area">

      <!-- No file -->
      <div v-if="!report.rawContent" class="empty-state">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" class="empty-icon">
          <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
        </svg>
        <p class="empty-text">No file loaded.</p>
        <button class="empty-link" @click="$router.push('/')">← Back to Generate</button>
      </div>

      <!-- Document tab -->
      <div v-else-if="activeTab === 'Document'" class="doc-panel">
        <div v-if="docxRendering" class="doc-loading">
          <svg class="spin doc-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
          </svg>
          <p>Rendering document…</p>
        </div>
        <div v-else-if="docxError" class="doc-error">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/>
          </svg>
          <p class="selectable">{{ docxError }}</p>
        </div>
        <div ref="docxContainer" class="docx-wrap" :class="{ 'docx-wrap--hidden': docxRendering || docxError }" />
      </div>

      <!-- Parsed Data tab -->
      <div v-else-if="activeTab === 'Parsed Data'" class="parsed-panel">
        <div v-if="!report.parsedData" class="empty-state">
          <p class="empty-text">No parsed data yet.</p>
          <button class="empty-link" @click="$router.push('/')">← Go back to parse</button>
        </div>

        <div v-else class="parsed-content">
          <!-- Summary strip -->
          <div class="summary-strip">
            <div class="summary-chip">
              <span class="chip-label">Type</span>
              <span class="chip-val">{{ report.parsedData.reportType }}</span>
            </div>
            <div class="summary-chip">
              <span class="chip-label">Period</span>
              <span class="chip-val">{{ report.parsedData.dateStart }} → {{ report.parsedData.dateEnd }}</span>
            </div>
            <div v-if="report.parsedData.totalHours" class="summary-chip">
              <span class="chip-label">Total Hours</span>
              <span class="chip-val accent">{{ report.parsedData.totalHours }}h</span>
            </div>
            <div v-if="report.parsedData.entries?.length" class="summary-chip">
              <span class="chip-label">Entries</span>
              <span class="chip-val">{{ report.parsedData.entries.length }}</span>
            </div>
          </div>

          <!-- Summary text -->
          <div v-if="report.parsedData.summary" class="summary-card">
            <p class="summary-text selectable">{{ report.parsedData.summary }}</p>
          </div>

          <!-- Task table -->
          <TaskTable
            v-if="report.parsedData.entries?.length"
            :rows="report.parsedData.entries"
            @remove="report.removeEntry($event)"
          />

          <!-- Ongoing tasks -->
          <div v-if="report.parsedData.ongoingTasks?.length">
            <h3 class="section-label">Ongoing Tasks</h3>
            <div v-for="task in report.parsedData.ongoingTasks" :key="task.name" class="ongoing-card">
              <div class="ongoing-head">
                <span class="ongoing-name">{{ task.name }}</span>
                <span class="ongoing-pct">{{ task.percentComplete }}%</span>
              </div>
              <div class="progress-track">
                <div class="progress-fill" :style="{ width: task.percentComplete + '%' }" />
              </div>
              <ul v-if="task.nextSteps?.length" class="steps-list">
                <li v-for="step in task.nextSteps" :key="step" class="step-item">
                  <span class="step-dot" />
                  <span class="selectable">{{ step }}</span>
                </li>
              </ul>
            </div>
          </div>

          <!-- Raw JSON -->
          <details class="raw-details">
            <summary class="raw-summary">Raw JSON</summary>
            <pre class="raw-json selectable">{{ JSON.stringify(report.parsedData, null, 2) }}</pre>
          </details>
        </div>
      </div>

      <!-- Rendered markdown -->
      <div
        v-else-if="activeTab === 'Rendered'"
        class="md-wrap md-preview selectable"
        v-html="rendered"
      />

      <!-- Raw markdown -->
      <pre v-else class="raw-md selectable">{{ report.rawContent }}</pre>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, toRaw, watch, nextTick } from 'vue'
import { marked } from 'marked'
import { renderAsync } from 'docx-preview'
import { useReportStore } from '@/stores/report.js'
import { useSettingsStore } from '@/stores/settings.js'
import { cole } from '@/lib/cole.js'
import TaskTable from '@/components/TaskTable.vue'

const report    = useReportStore()
const settings  = useSettingsStore()
const generating = ref(false)
const genError   = ref('')

const activeTab = ref(report.parsedData ? 'Document' : 'Rendered')

const availableTabs = computed(() => {
  const tabs = ['Rendered', 'Raw']
  if (report.parsedData) tabs.unshift('Parsed Data', 'Document')
  return tabs
})

const rendered = computed(() =>
  report.rawContent ? marked.parse(report.rawContent) : ''
)

const docxContainer = ref(null)
const docxRendering = ref(false)
const docxError     = ref('')
const docxRendered  = ref(false)

watch(activeTab, async (tab) => {
  if (tab === 'Document' && !docxRendered.value && report.parsedData) {
    await nextTick()
    await renderDocx()
  }
}, { immediate: true })

watch(() => report.parsedData, () => {
  docxRendered.value = false
  if (activeTab.value === 'Document') renderDocx()
})

async function renderDocx() {
  if (!report.parsedData) return
  docxRendering.value = true
  docxError.value     = ''
  try {
    await settings.load()
    const uint8 = await cole.renderReport({
      parsedData:   toRaw(report.parsedData),
      userInfo:     buildUserInfo(),
      employeeType: report.employeeType,
    })
    await nextTick()
    if (docxContainer.value) {
      await renderAsync(uint8, docxContainer.value, null, {
        className:    'docx',
        inWrapper:    true,
        ignoreWidth:  false,
        ignoreHeight: false,
        useBase64URL: true,
      })
      docxRendered.value = true
    }
  } catch (err) {
    docxError.value = err.message || 'Failed to render preview.'
  } finally {
    docxRendering.value = false
  }
}

function buildUserInfo() {
  return {
    reporterName:    settings.reporterName,
    reporterNameSig: settings.reporterNameSig,
    position:        settings.position,
    office:          settings.office,
    supervisorName:  settings.supervisorName,
    supervisorPos1:  settings.supervisorPos1,
    supervisorPos2:  settings.supervisorPos2,
  }
}

async function generateDocx() {
  if (!report.parsedData) return
  generating.value = true
  genError.value   = ''
  try {
    await settings.load()
    const d = toRaw(report.parsedData)
    const defaultName = `Accomplishment Report ${d.dateStart} to ${d.dateEnd}.docx`
    const savePath = await cole.showSaveDialog({
      defaultPath: defaultName,
      filters: [{ name: 'Word Document', extensions: ['docx'] }],
    })
    if (!savePath) return
    await cole.generateReport({
      parsedData:   d,
      userInfo:     buildUserInfo(),
      employeeType: report.employeeType,
      outputPath:   savePath,
    })
    await cole.openPath(savePath)
  } catch (err) {
    genError.value = err.message || 'Failed to generate report.'
  } finally {
    generating.value = false
  }
}
</script>

<style scoped>
.preview-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--color-bg);
}

/* Toolbar */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 52px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
  gap: 12px;
}
.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.back-btn {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.12s, color 0.12s;
}
.back-btn:hover { background: var(--color-surface-2); color: var(--color-text); }
.back-btn svg { width: 14px; height: 14px; }

.file-meta { display: flex; align-items: center; gap: 8px; min-width: 0; }
.file-meta-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 200px;
}
.file-meta-range {
  font-size: 11px;
  color: var(--color-text-muted);
  white-space: nowrap;
}
.type-badge {
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.04em;
  padding: 2px 7px;
  border-radius: 20px;
  background: rgba(0,122,255,0.1);
  color: var(--color-primary);
  border: 1px solid rgba(0,122,255,0.2);
  white-space: nowrap;
}

/* Export button */
.export-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 14px;
  height: 32px;
  border-radius: 10px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.12s;
}
.export-btn:hover:not(:disabled) { background: #0071e3; }
.export-btn:disabled { opacity: 0.45; cursor: not-allowed; }
.export-btn svg { width: 13px; height: 13px; }

/* Tab bar */
.tab-bar {
  display: flex;
  border: 1px solid var(--color-border);
  border-radius: 10px;
  overflow: hidden;
  background: var(--color-surface-2);
}
.tab-btn {
  padding: 0 12px;
  height: 30px;
  border: none;
  border-right: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.12s, color 0.12s;
}
.tab-btn:last-child { border-right: none; }
.tab-btn:hover { background: var(--color-surface-3); color: var(--color-text); }
.tab-btn--active { background: var(--color-primary); color: #fff; }

/* Gen error */
.gen-error {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 20px;
  background: rgba(255,59,48,0.07);
  border-bottom: 1px solid rgba(255,59,48,0.15);
  flex-shrink: 0;
}
.gen-error-icon { width: 15px; height: 15px; color: var(--color-danger); flex-shrink: 0; margin-top: 1px; }
.gen-error-text { font-size: 12px; color: var(--color-danger); margin: 0; line-height: 1.5; }

/* Content area */
.content-area { flex: 1; overflow: auto; }

/* Empty state */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 60%;
  text-align: center;
  padding: 40px;
}
.empty-icon { width: 48px; height: 48px; color: var(--color-border); margin-bottom: 16px; }
.empty-text { font-size: 13px; color: var(--color-text-muted); margin: 0 0 10px; }
.empty-link {
  font-size: 12px;
  color: var(--color-primary);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
}
.empty-link:hover { text-decoration: underline; }

/* Document panel */
.doc-panel { height: 100%; }
.doc-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  gap: 12px;
  color: var(--color-text-muted);
  font-size: 12px;
}
.doc-spin { width: 22px; height: 22px; color: var(--color-primary); }
.doc-error {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin: 20px;
  padding: 14px;
  background: rgba(255,59,48,0.07);
  border: 1px solid rgba(255,59,48,0.2);
  border-radius: 12px;
  color: var(--color-danger);
  font-size: 12px;
}
.doc-error svg { width: 15px; height: 15px; flex-shrink: 0; margin-top: 1px; }
.docx-wrap {
  background: var(--color-surface-2);
  min-height: 100%;
  padding: 24px;
}
.docx-wrap--hidden { display: none; }
:deep(.docx-wrapper) { background: transparent !important; padding: 0 !important; }
:deep(.docx) { box-shadow: 0 4px 24px rgba(0,0,0,0.15); margin: 0 auto; }

/* Parsed panel */
.parsed-panel { padding: 20px 24px; }
.parsed-content { max-width: 900px; margin: 0 auto; display: flex; flex-direction: column; gap: 14px; }

/* Summary strip */
.summary-strip {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.summary-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  box-shadow: var(--shadow-sm);
}
.chip-label {
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--color-text-muted);
}
.chip-val {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text);
}
.chip-val.accent { color: var(--color-primary); font-family: 'JetBrains Mono', Consolas, monospace; }

/* Summary card */
.summary-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-left: 3px solid var(--color-primary);
  border-radius: 12px;
  padding: 12px 16px;
  box-shadow: var(--shadow-card);
}
.summary-text {
  font-size: 13px;
  color: var(--color-text);
  line-height: 1.6;
  margin: 0;
}

/* Section label */
.section-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  margin: 4px 0 8px;
}

/* Ongoing tasks */
.ongoing-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 14px 16px;
  box-shadow: var(--shadow-card);
  margin-bottom: 8px;
}
.ongoing-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.ongoing-name { font-size: 13px; font-weight: 500; color: var(--color-text); }
.ongoing-pct { font-size: 12px; font-weight: 700; color: var(--color-primary); font-family: 'JetBrains Mono', Consolas, monospace; }
.progress-track { height: 4px; background: var(--color-surface-2); border-radius: 2px; margin-bottom: 10px; }
.progress-fill { height: 100%; background: var(--color-primary); border-radius: 2px; }
.steps-list { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 4px; }
.step-item { display: flex; align-items: baseline; gap: 6px; font-size: 12px; color: var(--color-text); }
.step-dot { width: 5px; height: 5px; border-radius: 50%; background: var(--color-primary); flex-shrink: 0; }

/* Raw JSON */
.raw-summary {
  font-size: 11px;
  color: var(--color-text-muted);
  cursor: pointer;
  user-select: none;
  transition: color 0.12s;
}
.raw-summary:hover { color: var(--color-text); }
.raw-json {
  margin-top: 10px;
  font-size: 11px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  color: var(--color-text-muted);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 10px;
  padding: 14px;
  overflow: auto;
  line-height: 1.5;
}

/* Rendered markdown */
.md-wrap { padding: 28px; max-width: 680px; margin: 0 auto; }

/* Raw markdown */
.raw-md {
  padding: 28px;
  font-size: 11px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  color: var(--color-text-muted);
  white-space: pre-wrap;
  line-height: 1.7;
  max-width: 680px;
  margin: 0 auto;
}

/* Spin */
.spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
