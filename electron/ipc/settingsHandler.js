import { app, safeStorage } from 'electron'
import { readFile, writeFile, mkdir } from 'node:fs/promises'
import path from 'node:path'

const SETTINGS_PATH = path.join(app.getPath('userData'), 'settings.json')
const SECRETS_PATH  = path.join(app.getPath('userData'), 'secrets.bin')

async function loadSettings() {
  try {
    const raw = await readFile(SETTINGS_PATH, 'utf-8')
    return JSON.parse(raw)
  } catch {
    return {}
  }
}

async function saveSettings(settings) {
  await mkdir(path.dirname(SETTINGS_PATH), { recursive: true })
  await writeFile(SETTINGS_PATH, JSON.stringify(settings, null, 2), 'utf-8')
}

async function loadSecret(key) {
  try {
    const raw = await readFile(SECRETS_PATH, 'utf-8')
    const secrets = JSON.parse(raw)
    if (!secrets[key]) return null
    const buf = Buffer.from(secrets[key], 'base64')
    return safeStorage.decryptString(buf)
  } catch {
    return null
  }
}

async function saveSecret(key, value) {
  let secrets = {}
  try {
    const raw = await readFile(SECRETS_PATH, 'utf-8')
    secrets = JSON.parse(raw)
  } catch {}
  const encrypted = safeStorage.encryptString(value)
  secrets[key] = Buffer.from(encrypted).toString('base64')
  await mkdir(path.dirname(SECRETS_PATH), { recursive: true })
  await writeFile(SECRETS_PATH, JSON.stringify(secrets, null, 2), 'utf-8')
}

export function registerSettingsHandlers(ipcMain) {
  ipcMain.handle('ollama:list-models', async (_, url) => {
    try {
      const res = await fetch(`${url}/api/tags`)
      const data = await res.json()
      return (data.models || []).map(m => m.name)
    } catch {
      return []
    }
  })

  ipcMain.handle('settings:get', async (_, key) => {
    if (key === 'apiKey') return loadSecret('apiKey')
    const settings = await loadSettings()
    return settings[key] ?? null
  })

  ipcMain.handle('settings:set', async (_, key, value) => {
    if (key === 'apiKey') {
      await saveSecret('apiKey', value)
      return
    }
    const settings = await loadSettings()
    settings[key] = value
    await saveSettings(settings)
  })

  ipcMain.handle('settings:get-all', async () => {
    const settings = await loadSettings()
    const hasApiKey = !!(await loadSecret('apiKey'))
    return { ...settings, hasApiKey }
  })
}
