import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSettingsStore = defineStore('settings', () => {
  const hasApiKey        = ref(false)
  const inputPath        = ref('')
  const outputPath       = ref('')
  const vaultRoot        = ref('')
  // LLM provider
  const provider         = ref('claude')   // 'claude' | 'ollama'
  const ollamaUrl        = ref('http://127.0.0.1:11434')
  const ollamaModel      = ref('llama3')
  // Profile
  const reporterName     = ref('')   // e.g. COLEWAN, CHRISTIAN FIARAWE
  const reporterNameSig  = ref('')   // e.g. CHRISTIAN F. COLEWAN (signature display)
  const position         = ref('')   // e.g. Information System Analyst I
  const office           = ref('')   // e.g. Compensation, Benefits, and Other Obligations
  const supervisorName   = ref('')   // e.g. SUSAN P. BUASEN-OCASEN
  const supervisorPos1   = ref('')   // e.g. Chief, CBOO
  const supervisorPos2   = ref('')   // e.g. Administrative Officer V

  async function load() {
    const all = await window.cole.getAllSettings()
    hasApiKey.value       = all.hasApiKey       ?? false
    inputPath.value       = all.inputPath       ?? ''
    outputPath.value      = all.outputPath      ?? ''
    vaultRoot.value       = all.vaultRoot       ?? ''
    provider.value        = all.provider        ?? 'claude'
    ollamaUrl.value       = all.ollamaUrl       ?? 'http://127.0.0.1:11434'
    ollamaModel.value     = all.ollamaModel     ?? 'llama3'
    reporterName.value    = all.reporterName    ?? ''
    reporterNameSig.value = all.reporterNameSig ?? ''
    position.value        = all.position        ?? ''
    office.value          = all.office          ?? ''
    supervisorName.value  = all.supervisorName  ?? ''
    supervisorPos1.value  = all.supervisorPos1  ?? ''
    supervisorPos2.value  = all.supervisorPos2  ?? ''
  }

  async function save(key, value) {
    await window.cole.setSetting(key, value)
    const map = {
      inputPath, outputPath, vaultRoot,
      reporterName, reporterNameSig, position, office,
      supervisorName, supervisorPos1, supervisorPos2,
      provider, ollamaUrl, ollamaModel,
    }
    if (key === 'apiKey') { hasApiKey.value = !!value; return }
    if (map[key]) map[key].value = value
  }

  return {
    hasApiKey, inputPath, outputPath, vaultRoot,
    reporterName, reporterNameSig, position, office,
    supervisorName, supervisorPos1, supervisorPos2,
    provider, ollamaUrl, ollamaModel,
    load, save,
  }
})
