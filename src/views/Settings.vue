<template>
  <div class="page">
    <div class="settings-grid">

      <!-- LEFT: AI Provider -->
      <div class="col">
        <div class="section-head">
          <div class="section-title">AI Provider</div>
          <div class="section-sub">Choose how your notes are parsed into reports.</div>
        </div>

        <div class="card">
          <div class="provider-tabs">
            <button
              v-for="p in providers"
              :key="p.id"
              @click="setProvider(p.id)"
              class="provider-tab"
              :class="{ 'provider-tab--active': provider === p.id }"
            >
              {{ p.label }}
            </button>
          </div>

          <!-- Claude -->
          <template v-if="provider === 'claude'">
            <div class="field">
              <label class="field-label">API Key</label>
              <div class="input-row">
                <input v-model="apiKey" :type="showKey ? 'text' : 'password'"
                       placeholder="sk-ant-api03-…" class="input mono flex-1" />
                <button @click="showKey = !showKey" class="btn-ghost">{{ showKey ? 'Hide' : 'Show' }}</button>
              </div>
              <p class="field-hint">
                Encrypted with <span class="mono accent">safeStorage</span> — never plaintext.
              </p>
            </div>
            <div class="status-row">
              <div class="status-dot-row">
                <span class="status-dot" :class="settings.hasApiKey ? 'dot--green' : 'dot--red'" />
                <span class="status-text">{{ settings.hasApiKey ? 'Key saved' : 'No key saved' }}</span>
              </div>
              <button @click="saveApiKey" :disabled="!apiKey" class="btn-primary-sm">Save Key</button>
            </div>
          </template>

          <!-- Gemini -->
          <template v-if="provider === 'gemini'">
            <div class="info-box">
              Free tier: 15 req/min, 1500 req/day. Get a key at <span class="mono accent">aistudio.google.com</span>.
            </div>
            <div class="field">
              <label class="field-label">Gemini API Key</label>
              <div class="input-row">
                <input v-model="geminiKey" :type="showGeminiKey ? 'text' : 'password'"
                       placeholder="AIza…" class="input mono flex-1" />
                <button @click="showGeminiKey = !showGeminiKey" class="btn-ghost">{{ showGeminiKey ? 'Hide' : 'Show' }}</button>
              </div>
            </div>
            <div class="status-row">
              <div class="status-dot-row">
                <span class="status-dot" :class="settings.hasGeminiKey ? 'dot--green' : 'dot--red'" />
                <span class="status-text">{{ settings.hasGeminiKey ? 'Key saved' : 'No key saved' }}</span>
              </div>
              <button @click="saveGeminiKey" :disabled="!geminiKey" class="btn-primary-sm">Save Key</button>
            </div>
          </template>

          <!-- Ollama -->
          <template v-if="provider === 'ollama'">
            <div class="info-box">
              Ollama must be running locally. No API key required.
            </div>
            <div class="field">
              <label class="field-label">Ollama URL</label>
              <input v-model="ollamaUrl" placeholder="http://127.0.0.1:11434" class="input mono" />
            </div>
            <div class="field">
              <label class="field-label">Model Name</label>
              <select v-model="ollamaModel" class="input mono">
                <option v-if="ollamaModels.length === 0" value="" disabled>No models found — is Ollama running?</option>
                <option v-for="m in ollamaModels" :key="m" :value="m">{{ m }}</option>
              </select>
            </div>
            <div class="save-row">
              <button @click="saveOllama" class="btn-primary-sm">Save</button>
            </div>
          </template>
        </div>

        <!-- Default Paths -->
        <div class="section-head" style="margin-top: 8px;">
          <div class="section-title">Default Paths</div>
        </div>
        <div class="card">
          <div class="field">
            <label class="field-label">Obsidian Vault Root</label>
            <input v-model="vaultRoot" placeholder="Root folder of your vault" class="input mono" />
          </div>
          <div class="field">
            <label class="field-label">Default Input Folder</label>
            <input v-model="inputPath" placeholder="Where Cole looks for .md files first" class="input mono" />
          </div>
          <div class="field">
            <label class="field-label">Default Output Folder</label>
            <input v-model="outputPath" placeholder="Where generated .docx files are saved" class="input mono" />
          </div>
          <div class="save-row">
            <button @click="savePaths" class="btn-primary-sm">Save Paths</button>
          </div>
        </div>
      </div>

      <!-- RIGHT: Profiles -->
      <div class="col">
        <!-- Reporter Profile -->
        <div class="section-head">
          <div class="section-title">Reporter Profile</div>
          <div class="section-sub">Fills your name, position, and office in generated .docx files.</div>
        </div>
        <div class="card">
          <div class="field">
            <label class="field-label">Full Name (Header)</label>
            <input v-model="reporterName" placeholder="LAST, FIRST MIDDLE" class="input mono" />
          </div>
          <div class="field">
            <label class="field-label">Full Name (Signature)</label>
            <input v-model="reporterNameSig" placeholder="First M. Last" class="input mono" />
          </div>
          <div class="field">
            <label class="field-label">Position</label>
            <input v-model="position" placeholder="Position title" class="input mono" />
          </div>
          <div class="field">
            <label class="field-label">College / Office</label>
            <input v-model="office" placeholder="College or office name" class="input mono" />
          </div>
          <div class="save-row">
            <button @click="saveProfile" class="btn-primary-sm">Save Profile</button>
          </div>
        </div>

        <!-- Supervisor Profile -->
        <div class="section-head" style="margin-top: 8px;">
          <div class="section-title">Supervisor Profile</div>
          <div class="section-sub">Used for the approval signature block in reports.</div>
        </div>
        <div class="card">
          <div class="field">
            <label class="field-label">Supervisor Name</label>
            <input v-model="supervisorName" placeholder="Supervisor full name" class="input mono" />
          </div>
          <div class="field">
            <label class="field-label">Title Line 1</label>
            <input v-model="supervisorPos1" placeholder="e.g. Chief, CBOO" class="input mono" />
          </div>
          <div class="field">
            <label class="field-label">Title Line 2</label>
            <input v-model="supervisorPos2" placeholder="e.g. Administrative Officer V" class="input mono" />
          </div>
          <div class="save-row">
            <button @click="saveSupervisor" class="btn-primary-sm">Save Supervisor</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <transition name="fade">
      <div v-if="toast" class="toast">{{ toast }}</div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings.js'
import { cole } from '@/lib/cole.js'

const settings = useSettingsStore()

const apiKey        = ref('')
const showKey       = ref(false)
const geminiKey     = ref('')
const showGeminiKey = ref(false)
const toast         = ref('')
const provider      = ref('claude')

const ollamaUrl    = ref('http://127.0.0.1:11434')
const ollamaModel  = ref('llama3')
const ollamaModels = ref([])

const reporterName    = ref('')
const reporterNameSig = ref('')
const position        = ref('')
const office          = ref('')
const supervisorName  = ref('')
const supervisorPos1  = ref('')
const supervisorPos2  = ref('')
const vaultRoot       = ref('')
const inputPath       = ref('')
const outputPath      = ref('')

const providers = [
  { id: 'claude', label: 'Claude (API)' },
  { id: 'gemini', label: 'Gemini (free)' },
  { id: 'ollama', label: 'Ollama (offline)' },
]

async function fetchOllamaModels() {
  ollamaModels.value = await cole.listOllamaModels(ollamaUrl.value) ?? []
}

onMounted(async () => {
  await settings.load()
  provider.value        = settings.provider
  ollamaUrl.value       = settings.ollamaUrl
  ollamaModel.value     = settings.ollamaModel
  if (provider.value === 'ollama') await fetchOllamaModels()
  reporterName.value    = settings.reporterName    || 'COLEWAN, CHRISTIAN FIARAWE'
  reporterNameSig.value = settings.reporterNameSig || 'CHRISTIAN F. COLEWAN'
  position.value        = settings.position        || 'Information System Analyst I'
  office.value          = settings.office          || 'Compensation, Benefits, and Other Obligations'
  supervisorName.value  = settings.supervisorName  || 'SUSAN P. BUASEN-OCASEN'
  supervisorPos1.value  = settings.supervisorPos1  || 'Chief, CBOO'
  supervisorPos2.value  = settings.supervisorPos2  || 'Administrative Officer V'
  vaultRoot.value       = settings.vaultRoot
  inputPath.value       = settings.inputPath
  outputPath.value      = settings.outputPath
})

function showToast(msg) {
  toast.value = msg
  setTimeout(() => { toast.value = '' }, 2200)
}

async function setProvider(value) {
  provider.value = value
  await settings.save('provider', value)
  if (value === 'ollama') await fetchOllamaModels()
  showToast(`Provider set to ${value}`)
}

async function saveApiKey() {
  await settings.save('apiKey', apiKey.value)
  apiKey.value = ''
  showToast('API key saved')
}

async function saveGeminiKey() {
  await settings.save('geminiKey', geminiKey.value)
  geminiKey.value = ''
  await settings.load()
  showToast('Gemini key saved')
}

async function saveOllama() {
  await settings.save('ollamaUrl', ollamaUrl.value)
  await settings.save('ollamaModel', ollamaModel.value)
  showToast('Ollama settings saved')
}

async function saveProfile() {
  await settings.save('reporterName',    reporterName.value)
  await settings.save('reporterNameSig', reporterNameSig.value)
  await settings.save('position',        position.value)
  await settings.save('office',          office.value)
  showToast('Profile saved')
}

async function saveSupervisor() {
  await settings.save('supervisorName',  supervisorName.value)
  await settings.save('supervisorPos1',  supervisorPos1.value)
  await settings.save('supervisorPos2',  supervisorPos2.value)
  showToast('Supervisor saved')
}

async function savePaths() {
  await settings.save('vaultRoot',   vaultRoot.value)
  await settings.save('inputPath',   inputPath.value)
  await settings.save('outputPath',  outputPath.value)
  showToast('Paths saved')
}
</script>

<style scoped>
.page {
  height: 100%;
  overflow: auto;
  background: var(--color-bg);
  display: flex;
  flex-direction: column;
}

.settings-grid {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  padding: 24px;
  align-items: start;
}

.col {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* Section heading */
.section-head { display: flex; flex-direction: column; gap: 2px; }
.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
  letter-spacing: -0.01em;
}
.section-sub {
  font-size: 11px;
  color: var(--color-text-muted);
  line-height: 1.4;
}

/* Card */
.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 14px;
  padding: 16px;
  box-shadow: var(--shadow-card);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* Provider tabs */
.provider-tabs {
  display: flex;
  gap: 6px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 10px;
  padding: 4px;
}
.provider-tab {
  flex: 1;
  padding: 6px 8px;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: var(--color-text-muted);
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.12s;
  white-space: nowrap;
}
.provider-tab:hover { color: var(--color-text); }
.provider-tab--active {
  background: var(--color-surface);
  color: var(--color-primary);
  box-shadow: var(--shadow-sm);
}

/* Field */
.field { display: flex; flex-direction: column; gap: 5px; }
.field-label { font-size: 11px; font-weight: 500; color: var(--color-text-muted); }
.field-hint { font-size: 10px; color: var(--color-text-muted); margin: 0; }
.mono { font-family: 'JetBrains Mono', Consolas, monospace; }
.accent { color: var(--color-primary); }

/* Input row */
.input-row { display: flex; gap: 8px; }
.flex-1 { flex: 1; }

/* Input */
.input {
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 9px;
  padding: 8px 10px;
  font-size: 12px;
  color: var(--color-text);
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  width: 100%;
}
.input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(0,122,255,0.12);
}
.input::placeholder { color: var(--color-text-muted); opacity: 0.5; }

/* Info box */
.info-box {
  font-size: 11px;
  color: var(--color-text-muted);
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 8px 10px;
  line-height: 1.5;
}

/* Status row */
.status-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.status-dot-row { display: flex; align-items: center; gap: 6px; }
.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot--green { background: var(--color-success); }
.dot--red   { background: var(--color-danger); }
.status-text { font-size: 11px; color: var(--color-text-muted); }

/* Save row */
.save-row { display: flex; justify-content: flex-end; }

/* Buttons */
.btn-primary-sm {
  padding: 6px 14px;
  height: 32px;
  border-radius: 9px;
  border: none;
  background: var(--color-primary);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.12s, opacity 0.12s;
  display: flex;
  align-items: center;
}
.btn-primary-sm:hover:not(:disabled) { background: #0071e3; }
.btn-primary-sm:disabled { opacity: 0.4; cursor: not-allowed; }

.btn-ghost {
  padding: 0 12px;
  height: 36px;
  border-radius: 9px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
  font-size: 12px;
  cursor: pointer;
  white-space: nowrap;
  transition: color 0.12s, border-color 0.12s;
}
.btn-ghost:hover { color: var(--color-text); border-color: var(--color-primary); }

/* Toast */
.toast {
  position: fixed;
  bottom: 20px;
  right: 24px;
  background: var(--color-success);
  color: #000;
  font-size: 13px;
  font-weight: 600;
  padding: 10px 18px;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.2);
  z-index: 999;
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
