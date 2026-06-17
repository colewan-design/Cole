import { contextBridge, ipcRenderer } from 'electron'

contextBridge.exposeInMainWorld('cole', {
  // File I/O
  openFileDialog: (options) => ipcRenderer.invoke('file:open-dialog', options),
  readFile:       (filePath) => ipcRenderer.invoke('file:read', filePath),
  saveFile:       (filePath, buffer) => ipcRenderer.invoke('file:save', filePath, buffer),
  showSaveDialog: (options) => ipcRenderer.invoke('file:save-dialog', options),
  openPath:       (filePath) => ipcRenderer.invoke('file:open-path', filePath),

  // Settings (encrypted via safeStorage)
  getSetting: (key)        => ipcRenderer.invoke('settings:get', key),
  setSetting: (key, value) => ipcRenderer.invoke('settings:set', key, value),
  getAllSettings: ()       => ipcRenderer.invoke('settings:get-all'),

  // Shell
  openExternal: (url) => ipcRenderer.invoke('shell:open-external', url),

  // Ollama
  listOllamaModels: (url) => ipcRenderer.invoke('ollama:list-models', url),

  // Claude API (Phase 2)
  parseWithClaude: (payload) => ipcRenderer.invoke('claude:parse', payload),

  // Docx generation (Phase 3)
  generateReport: (payload) => ipcRenderer.invoke('docx:generate', payload),
  renderReport:   (payload) => ipcRenderer.invoke('docx:render',   payload),
})
