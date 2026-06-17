<template>
  <div class="page">
    <div class="home-grid">

      <!-- LEFT: File drop zone -->
      <div class="left-col">
        <div class="col-heading">Upload File</div>
        <p class="col-sub">Select or drag in your Obsidian <span class="mono">.md</span> note to begin.</p>

        <FileDropZone v-if="!report.filePath" @file-selected="onFile" class="drop-zone-area" />

        <div v-else class="file-loaded">
          <div class="file-icon-wrap">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
            </svg>
          </div>
          <div class="file-details">
            <div class="file-name">{{ report.fileName }}</div>
            <div class="file-path">{{ report.filePath }}</div>
          </div>
          <button class="clear-btn" @click="report.reset()" title="Remove file">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <!-- Paste JSON section -->
        <div class="paste-section">
          <button @click="showPasteJson = !showPasteJson" class="paste-toggle">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
                 class="paste-arrow" :class="{ 'paste-arrow--open': showPasteJson }">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
            Paste JSON from Claude.ai
          </button>

          <div v-if="showPasteJson" class="paste-body">
            <div class="paste-hint">
              Parse your note on <span class="mono accent">claude.ai</span> manually, paste the JSON output below — no API key needed.
              <button @click="copyPrompt" class="paste-link">Copy prompt ↗</button>
            </div>
            <textarea
              v-model="pastedJson"
              placeholder='{"reportType":"AR","dateStart":"...","entries":[...]}'
              class="input paste-textarea selectable"
            />
            <div v-if="pasteError" class="paste-error">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/>
              </svg>
              {{ pasteError }}
            </div>
            <button @click="usePastedJson" :disabled="!pastedJson.trim()" class="btn btn--primary">
              Generate .docx from JSON
            </button>
          </div>
        </div>
      </div>

      <!-- RIGHT: Configuration -->
      <div class="right-col">

        <!-- Date Range -->
        <div class="config-block">
          <div class="config-label">Date Range</div>
          <div class="date-row">
            <div class="date-field">
              <label class="field-label">Start Date</label>
              <div class="date-input-wrap">
                <svg class="date-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
                <input type="date" v-model="dateStart" @change="onDateChange" class="input date-input" />
              </div>
            </div>
            <div class="date-field">
              <label class="field-label">End Date</label>
              <div class="date-input-wrap">
                <svg class="date-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
                <input type="date" v-model="dateEnd" @change="onDateChange" class="input date-input" />
              </div>
            </div>
          </div>
        </div>

        <!-- Report Type -->
        <div class="config-block">
          <div class="config-label">Report Type</div>
          <p class="config-sub">Select if the file has an accomplishment for the selected date range</p>
          <div class="type-grid">
            <button
              v-for="rt in reportTypes"
              :key="rt.id"
              @click="setReportType(rt.id)"
              class="type-card"
              :class="{ 'type-card--active': reportType === rt.id }"
            >
              <i :class="['fi', rt.icon, 'type-icon']" />
              <span class="type-id">{{ rt.id }}</span>
              <span class="type-label">{{ rt.label }}</span>
              <div v-if="reportType === rt.id" class="type-dot" />
            </button>
          </div>
        </div>

        <!-- Employee Type (AR only) -->
        <div v-if="reportType === 'AR'" class="config-block">
          <div class="config-label">Employee Type</div>
          <div class="emp-grid">
            <button
              v-for="et in employeeTypes"
              :key="et.id"
              @click="setEmployeeType(et.id)"
              class="emp-btn"
              :class="{ 'emp-btn--active': employeeType === et.id }"
            >
              <i :class="['fi', et.icon]" />
              <span>{{ et.label }}</span>
              <div v-if="employeeType === et.id" class="emp-dot" />
            </button>
          </div>
        </div>

        <!-- Parse error -->
        <div v-if="parseError" class="error-box">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="error-icon">
            <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <p class="error-text selectable">{{ parseError }}</p>
        </div>

        <!-- Actions -->
        <div class="actions-row">
          <button
            @click="previewMd"
            :disabled="!report.filePath"
            class="btn btn--secondary"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
              <circle cx="12" cy="12" r="3"/>
            </svg>
            Preview MD
          </button>

          <button
            :disabled="!canParse"
            class="btn btn--primary"
            @click="parse"
          >
            <svg v-if="parsing" class="spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="3"/>
              <path d="M12 2v3M12 19v3M4.22 4.22l2.12 2.12M17.66 17.66l2.12 2.12M2 12h3M19 12h3M4.22 19.78l2.12-2.12M17.66 6.34l2.12-2.12"/>
            </svg>
            {{ parsing ? 'Parsing…' : parseLabel }}
          </button>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useReportStore } from '@/stores/report.js'
import { useSettingsStore } from '@/stores/settings.js'
import { cole } from '@/lib/cole.js'
import FileDropZone from '@/components/FileDropZone.vue'

const router   = useRouter()
const report   = useReportStore()
const settings = useSettingsStore()

onMounted(() => settings.load())

const dateStart    = ref(report.dateStart)
const dateEnd      = ref(report.dateEnd)
const reportType   = ref(report.reportType)
const employeeType = ref(report.employeeType)

const reportTypes = [
  { id: 'AR',  icon: 'fi-rr-list-check',   label: 'Accomplishment' },
  { id: 'DTR', icon: 'fi-rr-clock',         label: 'Daily Time Record' },
  { id: 'PRG', icon: 'fi-rr-chart-line-up', label: 'Progress Report' },
]

const employeeTypes = [
  { id: 'casual',    icon: 'fi-rr-user',         label: 'Casual' },
  { id: 'plantilla', icon: 'fi-rr-briefcase',     label: 'Plantilla' },
  { id: 'cos',       icon: 'fi-rr-file-contract', label: 'COS' },
  { id: 'external',  icon: 'fi-rr-building',      label: 'External' },
]

function setReportType(id) {
  reportType.value = id
  report.setReportType(id)
}

function setEmployeeType(id) {
  employeeType.value = id
  report.setEmployeeType(id)
}

const parsing    = ref(false)
const parseError = ref('')
const showPasteJson = ref(false)
const pastedJson    = ref('')
const pasteError    = ref('')

const canParse = computed(() =>
  !!report.filePath && !!dateStart.value && !!dateEnd.value && !parsing.value
)

const parseLabel = computed(() => {
  if (settings.provider === 'ollama') return `Parse with ${settings.ollamaModel}`
  if (settings.provider === 'gemini') return 'Parse with Gemini'
  return 'Generate Report'
})

function onFile(info) { report.setFile(info) }
function onDateChange() { report.setDateRange(dateStart.value, dateEnd.value) }
function previewMd() { router.push('/preview') }

async function parse() {
  if (!canParse.value) return
  parsing.value    = true
  parseError.value = ''
  try {
    const result = await cole.parseWithClaude({
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

function usePastedJson() {
  pasteError.value = ''
  try {
    const stripped = pastedJson.value.replace(/```(?:json)?/g, '').trim()
    const start = stripped.indexOf('{')
    const end   = stripped.lastIndexOf('}')
    if (start === -1 || end === -1) throw new Error('No JSON object found in pasted text.')
    const data = JSON.parse(stripped.slice(start, end + 1))
    if (!data.entries || !data.reportType) throw new Error('Missing required fields: reportType or entries.')
    report.setParsedData(data)
    router.push('/preview')
  } catch (err) {
    pasteError.value = err.message || 'Invalid JSON.'
  }
}

function copyPrompt() {
  const prompt = `You are a government report assistant. Extract work entries from the provided Obsidian daily work log.
Return ONLY valid JSON — no markdown fences, no explanation, no extra text. Just the raw JSON object.

Date range: ${report.dateStart || 'YYYY-MM-DD'} to ${report.dateEnd || 'YYYY-MM-DD'}

Extract all work entries within the date range and return this exact JSON:
{
  "reportType": "AR",
  "dateStart": "${report.dateStart || ''}",
  "dateEnd": "${report.dateEnd || ''}",
  "entries": [
    { "date": "YYYY-MM-DD", "tasks": ["Plain sentence describing what was accomplished."], "hoursWorked": 8 }
  ],
  "summary": "A 1-2 sentence plain-language summary of the period's accomplishments."
}

Work log:
${report.rawContent || '(no file loaded — select a file first)'}`
  navigator.clipboard.writeText(prompt)
}

watch(() => report.dateStart,    v => { dateStart.value    = v })
watch(() => report.dateEnd,      v => { dateEnd.value      = v })
watch(() => report.reportType,   v => { reportType.value   = v })
watch(() => report.employeeType, v => { employeeType.value = v })
</script>

<style scoped>
.page {
  height: 100%;
  overflow: auto;
  background: var(--color-bg);
  display: flex;
  flex-direction: column;
}

/* Two-column grid */
.home-grid {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  padding: 24px;
  align-items: stretch;
}

/* ── Left column ── */
.left-col {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Make the drop zone grow to fill remaining left-column height */
.drop-zone-area {
  flex: 1;
  min-height: 180px;
}

.col-heading {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text);
  letter-spacing: -0.01em;
}
.col-sub {
  font-size: 12px;
  color: var(--color-text-muted);
  margin: -6px 0 0;
  line-height: 1.5;
}
.mono { font-family: 'JetBrains Mono', Consolas, monospace; }
.accent { color: var(--color-primary); }

/* File loaded row */
.file-loaded {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  box-shadow: var(--shadow-card);
}
.file-icon-wrap {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: rgba(0,122,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--color-primary);
}
.file-icon-wrap svg { width: 18px; height: 18px; }
.file-details { flex: 1; min-width: 0; }
.file-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.file-path {
  font-size: 10px;
  color: var(--color-text-muted);
  font-family: 'JetBrains Mono', Consolas, monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}
.clear-btn {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: color 0.12s, background 0.12s;
}
.clear-btn:hover { color: var(--color-danger); background: rgba(255,59,48,0.1); }
.clear-btn svg { width: 14px; height: 14px; }

/* Paste section */
.paste-section {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  padding: 14px 16px;
  box-shadow: var(--shadow-card);
}
.paste-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 12px;
  color: var(--color-text-muted);
  padding: 0;
  transition: color 0.12s;
}
.paste-toggle:hover { color: var(--color-text); }
.paste-arrow { width: 14px; height: 14px; transition: transform 0.15s; }
.paste-arrow--open { transform: rotate(90deg); }
.paste-body { display: flex; flex-direction: column; gap: 10px; margin-top: 12px; }
.paste-hint {
  font-size: 11px;
  color: var(--color-text-muted);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 10px 12px;
  line-height: 1.5;
}
.paste-link { color: var(--color-primary); background: none; border: none; cursor: pointer; font-size: 11px; padding: 0; margin-left: 6px; }
.paste-link:hover { text-decoration: underline; }
.paste-textarea { height: 100px; resize: none; font-family: 'JetBrains Mono', Consolas, monospace; font-size: 11px; }
.paste-error { display: flex; align-items: center; gap: 6px; font-size: 11px; color: var(--color-danger); }
.paste-error svg { width: 13px; height: 13px; flex-shrink: 0; }

/* ── Right column ── */
.right-col {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.config-block {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  padding: 16px 18px;
  box-shadow: var(--shadow-card);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.config-label {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
}
.config-sub {
  font-size: 11px;
  color: var(--color-text-muted);
  margin: -6px 0 0;
  line-height: 1.5;
}

/* Date row */
.date-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}
.date-field { display: flex; flex-direction: column; gap: 6px; }
.field-label { font-size: 11px; font-weight: 500; color: var(--color-text-muted); }
.date-input-wrap { position: relative; }
.date-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  width: 14px;
  height: 14px;
  color: var(--color-text-muted);
  pointer-events: none;
}
.date-input { padding-left: 32px; color-scheme: dark; }

/* Report type grid */
.type-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}
.type-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 8px 10px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  background: var(--color-surface-2);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all 0.12s;
  text-align: center;
}
.type-card:hover {
  border-color: var(--color-primary);
  color: var(--color-text);
}
.type-card--active {
  border-color: var(--color-primary);
  background: rgba(0,122,255,0.08);
  color: var(--color-primary);
}
.type-icon { font-size: 18px; line-height: 1; }
.type-id { font-size: 11px; font-weight: 700; letter-spacing: 0.04em; }
.type-label { font-size: 9px; line-height: 1.3; opacity: 0.75; }
.type-dot {
  position: absolute;
  top: 7px;
  right: 7px;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-primary);
}

/* Employee grid */
.emp-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}
.emp-btn {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 6px;
  border-radius: 10px;
  border: 1px solid var(--color-border);
  background: var(--color-surface-2);
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: 10px;
  font-weight: 600;
  transition: all 0.12s;
}
.emp-btn:hover { border-color: var(--color-primary); color: var(--color-text); }
.emp-btn--active {
  border-color: var(--color-primary);
  background: rgba(0,122,255,0.08);
  color: var(--color-primary);
}
.emp-dot {
  position: absolute;
  top: 5px;
  right: 5px;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--color-primary);
}

/* Error box */
.error-box {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  background: rgba(255,59,48,0.07);
  border: 1px solid rgba(255,59,48,0.2);
  border-radius: 12px;
}
.error-icon { width: 15px; height: 15px; color: var(--color-danger); flex-shrink: 0; margin-top: 1px; }
.error-text { font-size: 12px; color: var(--color-danger); margin: 0; line-height: 1.5; }

/* Actions */
.actions-row {
  display: flex;
  gap: 10px;
}
.btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 0 16px;
  height: 40px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: background 0.12s, opacity 0.12s;
}
.btn:disabled { opacity: 0.4; cursor: not-allowed; }
.btn svg { width: 15px; height: 15px; flex-shrink: 0; }
.btn--secondary {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  color: var(--color-text-muted);
}
.btn--secondary:hover:not(:disabled) {
  background: var(--color-surface-2);
  color: var(--color-text);
  border-color: var(--color-primary);
}
.btn--primary {
  flex: 1;
  background: var(--color-primary);
  color: #fff;
  box-shadow: 0 1px 6px rgba(0,122,255,0.25);
}
.btn--primary:hover:not(:disabled) { background: #0071e3; }

.spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* Shared input */
.input {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 10px;
  padding: 8px 12px;
  font-size: 13px;
  color: var(--color-text);
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  width: 100%;
}
.input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(0,122,255,0.12);
}
.input::placeholder { color: var(--color-text-muted); }
</style>
