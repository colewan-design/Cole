<template>
  <div class="page">
    <div class="page-inner">

      <!-- Step 1 — File -->
      <div class="card">
        <div class="step-head">
          <StepBadge :n="1" :done="!!report.filePath" />
          <span class="step-label">Select Markdown file</span>
        </div>

        <FileDropZone v-if="!report.filePath" @file-selected="onFile" />

        <div v-else class="file-row">
          <div class="file-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
            </svg>
          </div>
          <div class="file-info">
            <div class="file-name">{{ report.fileName }}</div>
            <div class="file-path">{{ report.filePath }}</div>
          </div>
          <button class="clear-btn" @click="report.reset()" title="Clear file">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Step 2 — Date range -->
      <div class="card">
        <div class="step-head">
          <StepBadge :n="2" :done="!!(report.dateStart && report.dateEnd)" />
          <span class="step-label">Date range</span>
        </div>
        <DateRangePicker
          v-model:modelStart="dateStart"
          v-model:modelEnd="dateEnd"
          @update:modelStart="onDateChange"
          @update:modelEnd="onDateChange"
        />
      </div>

      <!-- Step 3 — Report type -->
      <div class="card">
        <div class="step-head">
          <StepBadge :n="3" :done="true" />
          <span class="step-label">Report type</span>
        </div>
        <ReportTypeSelector v-model="reportType" @update:modelValue="report.setReportType($event)" />
      </div>

      <!-- Step 4 — Employee type (AR only) -->
      <div v-if="reportType === 'AR'" class="card">
        <div class="step-head">
          <StepBadge :n="4" :done="true" />
          <span class="step-label">Employee type</span>
        </div>
        <div class="emp-grid">
          <button
            v-for="et in employeeTypes"
            :key="et.id"
            @click="setEmployeeType(et.id)"
            class="emp-btn"
            :class="{ 'emp-btn--active': employeeType === et.id }"
          >
            <i :class="['fi', et.icon, 'text-base leading-none']" />
            <span class="emp-label">{{ et.label }}</span>
            <div v-if="employeeType === et.id" class="emp-dot" />
          </button>
        </div>
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

      <!-- Parse error -->
      <div v-if="parseError" class="error-box">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="error-icon">
          <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <p class="error-text selectable">{{ parseError }}</p>
      </div>

      <!-- OR divider -->
      <div class="divider">
        <div class="divider-line" />
        <span class="divider-label">or</span>
        <div class="divider-line" />
      </div>

      <!-- Paste JSON -->
      <div class="card">
        <button @click="showPasteJson = !showPasteJson" class="paste-toggle">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
               class="paste-arrow" :class="{ 'paste-arrow--open': showPasteJson }">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
          Paste JSON from Claude.ai (no API tokens)
        </button>

        <div v-if="showPasteJson" class="paste-body">
          <div class="paste-hint">
            Parse your note manually on <span class="paste-code">claude.ai</span>, paste the JSON output below, then generate the .docx — no API key needed.
            <button @click="copyPrompt" class="paste-link">Copy prompt ↗</button>
          </div>
          <textarea
            v-model="pastedJson"
            placeholder='Paste JSON here — e.g. {"reportType":"AR","dateStart":"...","entries":[...]}'
            class="input paste-textarea selectable"
          />
          <div v-if="pasteError" class="paste-error">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/>
            </svg>
            {{ pasteError }}
          </div>
          <button
            @click="usePastedJson"
            :disabled="!pastedJson.trim()"
            class="btn btn--primary"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
            </svg>
            Generate .docx from pasted JSON
          </button>
        </div>
      </div>

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

const dateStart    = ref(report.dateStart)
const dateEnd      = ref(report.dateEnd)
const reportType   = ref(report.reportType)
const employeeType = ref(report.employeeType)

const employeeTypes = [
  { id: 'casual',    icon: 'fi-rr-user',         label: 'Casual' },
  { id: 'plantilla', icon: 'fi-rr-briefcase',     label: 'Plantilla' },
  { id: 'cos',       icon: 'fi-rr-file-contract', label: 'COS' },
  { id: 'external',  icon: 'fi-rr-building',      label: 'External' },
]

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
  return 'Parse with Claude'
})

function onFile(info) { report.setFile(info) }
function onDateChange() { report.setDateRange(dateStart.value, dateEnd.value) }
function previewMd() { router.push('/preview') }

async function parse() {
  if (!canParse.value) return
  parsing.value    = true
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

const StepBadge = {
  props: ['n', 'done'],
  setup(props) {
    return () => h(
      'div',
      { class: `step-badge ${props.done ? 'step-badge--done' : ''}` },
      props.done
        ? [h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'3' },
            [h('polyline', { points:'20 6 9 17 4 12' })])]
        : [String(props.n)]
    )
  },
}
</script>

<style scoped>
.page {
  min-height: 100%;
  background: var(--color-bg);
  padding: 28px 24px;
}
.page-inner {
  max-width: 560px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Card */
.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  padding: 18px 20px;
  box-shadow: var(--shadow-card);
}

/* Step header */
.step-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}
.step-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text);
}

/* Step badge */
:deep(.step-badge) {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  flex-shrink: 0;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  color: var(--color-text-muted);
  transition: background 0.15s, color 0.15s;
}
:deep(.step-badge--done) {
  background: #007AFF;
  border-color: #007AFF;
  color: #fff;
}
:deep(.step-badge svg) { width: 12px; height: 12px; }

/* File row */
.file-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 12px;
}
.file-icon {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: rgba(0,122,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: #007AFF;
}
.file-icon svg { width: 16px; height: 16px; }
.file-info { flex: 1; min-width: 0; }
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
.clear-btn:hover { color: #FF3B30; background: rgba(255,59,48,0.1); }
.clear-btn svg { width: 14px; height: 14px; }

/* Employee type grid */
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
  padding: 10px 8px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  background: var(--color-surface-2);
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: 11px;
  transition: all 0.12s;
}
.emp-btn:hover {
  border-color: #007AFF;
  color: var(--color-text);
}
.emp-btn--active {
  border-color: #007AFF;
  background: rgba(0,122,255,0.08);
  color: #007AFF;
}
.emp-label {
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.02em;
  line-height: 1;
}
.emp-dot {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #007AFF;
}

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
  border-color: #007AFF;
}
.btn--primary {
  flex: 1;
  background: #007AFF;
  color: #fff;
  box-shadow: 0 1px 6px rgba(0,122,255,0.2);
}
.btn--primary:hover:not(:disabled) { background: #0071e3; }

/* Spin animation */
.spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

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
.error-icon { width: 15px; height: 15px; color: #FF3B30; flex-shrink: 0; margin-top: 1px; }
.error-text { font-size: 12px; color: #FF3B30; margin: 0; line-height: 1.5; }

/* OR divider */
.divider {
  display: flex;
  align-items: center;
  gap: 12px;
}
.divider-line { flex: 1; height: 1px; background: var(--color-border); }
.divider-label { font-size: 11px; color: var(--color-text-muted); }

/* Paste JSON */
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
.paste-arrow {
  width: 14px;
  height: 14px;
  transition: transform 0.15s;
}
.paste-arrow--open { transform: rotate(90deg); }

.paste-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 12px;
}
.paste-hint {
  font-size: 11px;
  color: var(--color-text-muted);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 10px 12px;
  line-height: 1.5;
}
.paste-code { color: #007AFF; font-family: 'JetBrains Mono', Consolas, monospace; }
.paste-link { color: #007AFF; background: none; border: none; cursor: pointer; font-size: 11px; padding: 0; margin-left: 6px; }
.paste-link:hover { text-decoration: underline; }
.paste-textarea { height: 120px; resize: none; font-family: 'JetBrains Mono', Consolas, monospace; font-size: 11px; }
.paste-error {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: #FF3B30;
}
.paste-error svg { width: 13px; height: 13px; flex-shrink: 0; }
</style>
