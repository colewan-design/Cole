function getLegacyCole() {
  if (typeof window === 'undefined') return null
  if (window.go?.main?.App) return null
  const legacy = window.cole
  return legacy && typeof legacy === 'object' ? legacy : null
}

function getGoBinding(method) {
  if (typeof window === 'undefined') return null
  return window.go?.main?.App?.[method] ?? null
}

async function call(method, ...args) {
  const binding = getGoBinding(method)
  if (binding) return binding(...args)

  const legacy = getLegacyCole()
  if (legacy && typeof legacy[method] === 'function') {
    return legacy[method](...args)
  }

  throw new Error(`Cole backend method "${method}" is unavailable.`)
}

export const cole = {
  openFileDialog: (options) => call('OpenFileDialog', options),
  readFile: (filePath) => call('ReadFile', filePath),
  saveFile: (filePath, buffer) => call('SaveFile', filePath, buffer),
  showSaveDialog: (options) => call('ShowSaveDialog', options),
  openPath: (filePath) => call('OpenPath', filePath),
  getSetting: (key) => call('GetSetting', key),
  setSetting: (key, value) => call('SetSetting', key, value),
  getAllSettings: () => call('GetAllSettings'),
  openExternal: (url) => call('OpenExternal', url),
  listOllamaModels: (url) => call('ListOllamaModels', url),
  parseWithClaude: (payload) => call('ParseWithClaude', payload),
  generateReport: (payload) => call('GenerateReport', payload),
  renderReport: async (payload) => {
    const result = await call('RenderReport', payload)
    if (result instanceof Uint8Array) return result
    if (Array.isArray(result)) return new Uint8Array(result)
    if (typeof result === 'string') {
      const binary = atob(result)
      const bytes = new Uint8Array(binary.length)
      for (let i = 0; i < binary.length; i += 1) {
        bytes[i] = binary.charCodeAt(i)
      }
      return bytes
    }
    throw new Error('Unexpected render payload from backend.')
  },
}
