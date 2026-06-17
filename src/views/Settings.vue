<template>
  <div class="p-6 max-w-xl mx-auto">
    <div class="mb-6">
      <h1 class="text-xl font-semibold text-white">Settings</h1>
      <p class="text-sm text-j-muted mt-1">API key, profile info, and default paths — stored locally.</p>
    </div>

    <!-- LLM Provider -->
    <section class="mb-5">
      <h2 class="section-label">LLM Provider</h2>
      <div class="card space-y-4">
        <div class="flex gap-2">
          <button @click="setProvider('claude')"
                  :class="['btn-provider', provider === 'claude' ? 'btn-provider--active' : '']">
            Claude (API)
          </button>
          <button @click="setProvider('gemini')"
                  :class="['btn-provider', provider === 'gemini' ? 'btn-provider--active' : '']">
            Gemini (free)
          </button>
          <button @click="setProvider('ollama')"
                  :class="['btn-provider', provider === 'ollama' ? 'btn-provider--active' : '']">
            Ollama (offline)
          </button>
        </div>

        <template v-if="provider === 'claude'">
          <div>
            <label class="field-label">API Key</label>
            <div class="flex gap-2 mt-1.5">
              <input v-model="apiKey" :type="showKey ? 'text' : 'password'"
                     placeholder="sk-ant-api03-…" class="input flex-1 font-mono" />
              <button @click="showKey = !showKey" class="btn-ghost px-3">{{ showKey ? 'Hide' : 'Show' }}</button>
            </div>
            <p class="text-xs text-j-muted mt-1.5">
              Encrypted with <span class="font-mono text-j-glow">safeStorage</span> — never plaintext.
            </p>
          </div>
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-1.5">
              <div class="w-2 h-2 rounded-full" :class="settings.hasApiKey ? 'bg-j-success' : 'bg-j-muted'" />
              <span class="text-xs text-j-muted">{{ settings.hasApiKey ? 'Key saved' : 'No key saved' }}</span>
            </div>
            <button @click="saveApiKey" :disabled="!apiKey" class="btn-primary text-xs px-3 py-1.5">Save Key</button>
          </div>
        </template>

        <template v-if="provider === 'gemini'">
          <div class="text-xs text-j-muted bg-j-surface border border-j-border rounded-lg px-3 py-2">
            Free tier: 15 requests/min, 1500 requests/day. Get a key at <span class="font-mono text-j-glow">aistudio.google.com</span>.
          </div>
          <div>
            <label class="field-label">Gemini API Key</label>
            <div class="flex gap-2 mt-1.5">
              <input v-model="geminiKey" :type="showGeminiKey ? 'text' : 'password'"
                     placeholder="AIza…" class="input flex-1 font-mono" />
              <button @click="showGeminiKey = !showGeminiKey" class="btn-ghost px-3">{{ showGeminiKey ? 'Hide' : 'Show' }}</button>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-1.5">
              <div class="w-2 h-2 rounded-full" :class="settings.hasGeminiKey ? 'bg-j-success' : 'bg-j-muted'" />
              <span class="text-xs text-j-muted">{{ settings.hasGeminiKey ? 'Key saved' : 'No key saved' }}</span>
            </div>
            <button @click="saveGeminiKey" :disabled="!geminiKey" class="btn-primary text-xs px-3 py-1.5">Save Key</button>
          </div>
        </template>

        <template v-if="provider === 'ollama'">
          <div class="text-xs text-j-muted bg-j-surface border border-j-border rounded-lg px-3 py-2">
            Ollama must be running locally. No API key required.
          </div>
          <div>
            <label class="field-label">Ollama URL</label>
            <input v-model="ollamaUrl" placeholder="http://127.0.0.1:11434"
                   class="input w-full mt-1.5 text-xs font-mono" />
          </div>
          <div>
            <label class="field-label">Model Name</label>
            <select v-model="ollamaModel" class="input w-full mt-1.5 text-xs font-mono">
              <option v-if="ollamaModels.length === 0" value="" disabled>No models found — is Ollama running?</option>
              <option v-for="m in ollamaModels" :key="m" :value="m">{{ m }}</option>
            </select>
          </div>
          <div class="flex justify-end">
            <button @click="saveOllama" class="btn-primary text-xs px-4 py-1.5">Save</button>
          </div>
        </template>
      </div>
    </section>

    <!-- Report Profile -->
    <section class="mb-5">
      <h2 class="section-label">Report Profile</h2>
      <p class="text-xs text-j-muted mb-3">Used to fill your name, position, and supervisor info in generated .docx files.</p>
      <div class="card space-y-3">
        <div>
          <label class="field-label">Full Name (Header)</label>
          <input v-model="reporterName" placeholder="LAST, FIRST MIDDLE"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div>
          <label class="field-label">Full Name (Signature)</label>
          <input v-model="reporterNameSig" placeholder="First M. Last"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div>
          <label class="field-label">Position</label>
          <input v-model="position" placeholder="Position title"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div>
          <label class="field-label">College / Office Assignment</label>
          <input v-model="office" placeholder="College or office name"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <hr class="border-j-border" />
        <div>
          <label class="field-label">Supervisor Name</label>
          <input v-model="supervisorName" placeholder="Supervisor full name"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div>
          <label class="field-label">Supervisor Title Line 1</label>
          <input v-model="supervisorPos1" placeholder="Title line 1"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div>
          <label class="field-label">Supervisor Title Line 2</label>
          <input v-model="supervisorPos2" placeholder="Title line 2"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div class="flex justify-end pt-1">
          <button @click="saveProfile" class="btn-primary text-xs px-4 py-1.5">Save Profile</button>
        </div>
      </div>
    </section>

    <!-- Default Paths -->
    <section class="mb-5">
      <h2 class="section-label">Default Paths</h2>
      <div class="card space-y-3">
        <div>
          <label class="field-label">Obsidian Vault Root</label>
          <input v-model="vaultRoot" placeholder="Root folder of your vault"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div>
          <label class="field-label">Default Input Folder</label>
          <input v-model="inputPath" placeholder="Where Cole looks for .md files first"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div>
          <label class="field-label">Default Output Folder</label>
          <input v-model="outputPath" placeholder="Where generated .docx files are saved"
                 class="input w-full mt-1.5 text-xs font-mono" />
        </div>
        <div class="flex justify-end pt-1">
          <button @click="savePaths" class="btn-primary text-xs px-4 py-1.5">Save Paths</button>
        </div>
      </div>
    </section>

    <!-- Toast -->
    <transition name="fade">
      <div v-if="toast" class="fixed bottom-6 right-6 bg-j-success text-black text-sm font-semibold px-4 py-2.5 rounded-xl shadow-lg">
        {{ toast }}
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings.js'

const settings = useSettingsStore()

const apiKey       = ref('')
const showKey      = ref(false)
const geminiKey    = ref('')
const showGeminiKey = ref(false)
const toast    = ref('')
const provider = ref('claude')

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

async function fetchOllamaModels() {
  ollamaModels.value = await window.cole.listOllamaModels(ollamaUrl.value) ?? []
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
  showToast(`Provider set to ${value === 'ollama' ? 'Ollama' : 'Claude'}`)
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
  await settings.save('supervisorName',  supervisorName.value)
  await settings.save('supervisorPos1',  supervisorPos1.value)
  await settings.save('supervisorPos2',  supervisorPos2.value)
  showToast('Profile saved')
}

async function savePaths() {
  await settings.save('vaultRoot',   vaultRoot.value)
  await settings.save('inputPath',   inputPath.value)
  await settings.save('outputPath',  outputPath.value)
  showToast('Paths saved')
}
</script>

<style scoped>
.section-label        { @apply text-xs font-semibold uppercase tracking-wider text-j-muted mb-3; }
.card                 { @apply bg-j-card border border-j-border rounded-xl p-4; }
.field-label          { @apply block text-xs font-medium text-j-muted; }
.input                { @apply bg-j-surface border border-j-border rounded-lg px-3 py-2 text-sm text-j-text focus:outline-none focus:border-j-accent focus:ring-1 focus:ring-j-accent/30 placeholder:text-j-muted/40; }
.btn-primary          { @apply bg-j-accent hover:bg-blue-500 text-white font-semibold rounded-lg transition-colors disabled:opacity-40 disabled:cursor-not-allowed; }
.btn-ghost            { @apply rounded-lg border border-j-border text-j-muted hover:text-white hover:border-j-accent/50 transition-colors; }
.btn-provider         { @apply flex-1 py-1.5 text-xs font-semibold rounded-lg border border-j-border text-j-muted hover:text-white transition-colors; }
.btn-provider--active { @apply border-j-accent text-j-accent bg-j-accent/10; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to       { opacity: 0; }
</style>
