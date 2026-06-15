import Anthropic from '@anthropic-ai/sdk'
import { app, safeStorage } from 'electron'
import { readFile } from 'node:fs/promises'
import path from 'node:path'

const SECRETS_PATH  = path.join(app.getPath('userData'), 'secrets.bin')
const SETTINGS_PATH = path.join(app.getPath('userData'), 'settings.json')

async function getProviderSettings() {
  try {
    const raw = await readFile(SETTINGS_PATH, 'utf-8')
    const s = JSON.parse(raw)
    return {
      provider:    s.provider    ?? 'claude',
      ollamaUrl:   s.ollamaUrl   ?? 'http://127.0.0.1:11434',
      ollamaModel: s.ollamaModel ?? 'llama3',
    }
  } catch {
    return { provider: 'claude', ollamaUrl: 'http://127.0.0.1:11434', ollamaModel: 'llama3' }
  }
}

async function getApiKey() {
  try {
    const raw = await readFile(SECRETS_PATH, 'utf-8')
    const secrets = JSON.parse(raw)
    if (!secrets.apiKey) return null
    return safeStorage.decryptString(Buffer.from(secrets.apiKey, 'base64'))
  } catch {
    return null
  }
}

const MONTHS = ['january','february','march','april','may','june','july','august','september','october','november','december']

// Parse a loose date string like "June 1", "June 4", "June 8, 2026", "2026-06-01"
// Returns a Date or null. Falls back to the note's year when no year is present.
function parseLooseDate(str, fallbackYear) {
  if (!str) return null
  str = str.trim()

  // ISO: 2026-06-01
  let m = str.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (m) return new Date(+m[1], +m[2] - 1, +m[3])

  // "June 1, 2026" or "June 1"
  m = str.match(/^(January|February|March|April|May|June|July|August|September|October|November|December)\s+(\d{1,2})(?:,?\s*(\d{4}))?$/i)
  if (m) {
    const mo  = MONTHS.indexOf(m[1].toLowerCase())
    const day = +m[2]
    const yr  = m[3] ? +m[3] : fallbackYear
    if (mo !== -1 && day) return new Date(yr, mo, day)
  }

  return null
}

// Detect the year of the note from its title line (e.g. "# 📅 June 2026")
function detectNoteYear(content) {
  const m = content.match(/\b(20\d{2})\b/)
  return m ? +m[1] : new Date().getFullYear()
}

// Parse a markdown table with a "Date" column into structured rows.
// Returns [{ date: Date, task: string, notes: string }] or null if no table found.
function parseMarkdownTable(content, fallbackYear) {
  const lines = content.split('\n')
  const rows  = []
  let headers = null

  for (const line of lines) {
    if (!line.startsWith('|')) { headers = null; continue }
    const cells = line.split('|').map(c => c.trim()).filter((_, i, a) => i > 0 && i < a.length - 1)
    if (!cells.length) continue

    // Separator row
    if (cells.every(c => /^[-:]+$/.test(c))) continue

    // Header row
    const lower = cells.map(c => c.toLowerCase())
    if (lower.some(c => c === 'task') && lower.some(c => c === 'date')) {
      headers = lower
      continue
    }

    if (!headers) continue

    const get = (key) => {
      const i = headers.indexOf(key)
      return i !== -1 ? cells[i] ?? '' : ''
    }

    const dateStr = get('date')
    const task    = get('task')
    const notes   = get('notes')

    if (!task || task === '---') continue

    const date = parseLooseDate(dateStr, fallbackYear)
    rows.push({ date, task, notes, dateStr })
  }

  return rows.length ? rows : null
}

// Convert parsed table rows into a flat text block the model can easily read.
// Filters to the requested date range when dates are parseable.
function tableRowsToText(rows, start, end) {
  const inRange = rows.filter(r => {
    if (!r.date) return true  // keep undated rows
    return r.date >= start && r.date <= end
  })

  if (!inRange.length) return null

  return inRange.map(r => {
    const d = r.date
      ? `${r.date.getFullYear()}-${String(r.date.getMonth()+1).padStart(2,'0')}-${String(r.date.getDate()).padStart(2,'0')}`
      : r.dateStr
    const lines = [`Date: ${d}`, `Task: ${r.task}`]
    if (r.notes) lines.push(`Notes: ${r.notes}`)
    return lines.join('\n')
  }).join('\n\n')
}

// Pre-process content before sending to the model.
// 1. If the note uses a markdown table with a "Task" + "Date" column, convert to plain text.
// 2. Otherwise fall back to section-header date filtering.
function filterContentByDateRange(content, dateStart, dateEnd) {
  const start       = new Date(dateStart)
  const end         = new Date(dateEnd)
  end.setHours(23, 59, 59, 999)
  const fallbackYear = detectNoteYear(content)

  // Try table extraction first
  const rows = parseMarkdownTable(content, fallbackYear)
  if (rows) {
    const text = tableRowsToText(rows, start, end)
    if (text) return text
  }

  // Fall back: section-header based filtering
  const DATE_PATTERNS = [
    /\b(\d{4}-\d{2}-\d{2})\b/,
    /\b(January|February|March|April|May|June|July|August|September|October|November|December)\s+(\d{1,2}),?\s+(\d{4})\b/i,
    /\b(\d{1,2})\/(\d{1,2})\/(\d{4})\b/,
  ]

  function parseHeaderDate(line) {
    if (!/^#{1,4}\s/.test(line)) return null
    for (const pat of DATE_PATTERNS) {
      const m = line.match(pat)
      if (!m) continue
      const d = new Date(m[0])
      if (!isNaN(d)) return d
    }
    return null
  }

  const lines = content.split('\n')
  const out   = []
  let inRange   = false
  let hasHeader = false

  for (const line of lines) {
    const d = parseHeaderDate(line)
    if (d !== null) {
      hasHeader = true
      inRange   = d >= start && d <= end
    }
    if (!hasHeader || inRange) out.push(line)
  }

  const filtered = out.join('\n').trim()
  return filtered.length > 100 ? filtered : content
}

const SYSTEM = `You are a government report assistant. Extract work entries from the provided Obsidian daily work log.
Return ONLY valid JSON — no markdown fences, no explanation, no extra text. Just the raw JSON object.`

function buildPrompt(reportType, dateStart, dateEnd, content) {
  const filtered = filterContentByDateRange(content, dateStart, dateEnd)
  const base = `Date range: ${dateStart} to ${dateEnd}\n\nWork log:\n${filtered}\n\n`

  if (reportType === 'AR') return base + `
You are writing an official Accomplishment Report for a government IT employee. The audience is HR and supervisors — non-technical readers.

Rules for writing task descriptions:
- Rewrite each task as a clear, plain-language sentence that a non-technical reader can understand.
- Describe WHAT was done and WHY it matters, not HOW it was coded.
- Avoid: component names, file paths, CSS properties, hex colors, variable names, code syntax, or technical jargon.
- Good example: "Redesigned the employee payroll edit form with reorganized deduction headers for better usability."
- Bad example: "Replaced edlg-net-box, grouped GSIS/HDMF headers via scoped CSS flex layout."
- If multiple tasks happened on the same day, list each as a separate string in the "tasks" array.
- If the log says WFH, begin the sentence with "WFH: ".
- Leaves, holidays, and weekends should have an empty tasks array [].

Extract all work entries within the date range and return this exact JSON:
{
  "reportType": "AR",
  "dateStart": "${dateStart}",
  "dateEnd": "${dateEnd}",
  "entries": [
    { "date": "YYYY-MM-DD", "tasks": ["Plain sentence describing what was accomplished."], "hoursWorked": 8 }
  ],
  "summary": "A 1-2 sentence plain-language summary of the period's accomplishments."
}`

  if (reportType === 'DTR') return base + `
Extract daily time records. Write task descriptions in plain, non-technical language suitable for HR. Return this exact JSON:
{
  "reportType": "DTR",
  "dateStart": "${dateStart}",
  "dateEnd": "${dateEnd}",
  "entries": [
    { "date": "YYYY-MM-DD", "timeIn": "08:00", "timeOut": "17:00", "hoursWorked": 8, "tasks": ["Plain sentence describing what was done."] }
  ],
  "totalHours": 0
}`

  return base + `
Extract progress on ongoing tasks. Write all descriptions in plain, non-technical language suitable for a supervisor. Return this exact JSON:
{
  "reportType": "PRG",
  "dateStart": "${dateStart}",
  "dateEnd": "${dateEnd}",
  "entries": [
    { "date": "YYYY-MM-DD", "tasks": ["Plain sentence describing what was done."], "hoursWorked": 8 }
  ],
  "ongoingTasks": [
    { "name": "Plain task name", "percentComplete": 75, "status": "In Progress", "nextSteps": ["Plain next step description."] }
  ],
  "summary": "A 1-2 sentence plain-language summary of overall progress."
}`
}

async function parseWithOllama({ ollamaUrl, ollamaModel, prompt }) {
  const res = await fetch(`${ollamaUrl}/api/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      model: ollamaModel,
      messages: [
        { role: 'system', content: SYSTEM },
        { role: 'user',   content: prompt },
      ],
      stream: false,
      options: { num_ctx: 16384, num_predict: 8192 },
    }),
  })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(`Ollama error ${res.status}: ${err}`)
  }
  const data = await res.json()
  return data.message?.content ?? ''
}

export function registerClaudeHandlers(ipcMain) {
  ipcMain.handle('claude:parse', async (_, { content, dateStart, dateEnd, reportType }) => {
    const { provider, ollamaUrl, ollamaModel } = await getProviderSettings()
    const prompt = buildPrompt(reportType, dateStart, dateEnd, content)

    let raw
    if (provider === 'ollama') {
      raw = await parseWithOllama({ ollamaUrl, ollamaModel, prompt })
    } else {
      const apiKey = await getApiKey()
      if (!apiKey) throw new Error('No API key found. Add your Claude API key in Settings first.')
      const client  = new Anthropic({ apiKey })
      const message = await client.messages.create({
        model: 'claude-sonnet-4-6',
        max_tokens: 8192,
        system: [{ type: 'text', text: SYSTEM, cache_control: { type: 'ephemeral' } }],
        messages: [{ role: 'user', content: prompt }],
      })
      if (message.stop_reason === 'max_tokens') {
        throw new Error('Response was cut off — your notes may be too large. Try a narrower date range.')
      }
      raw = message.content[0].text.trim()
    }

    const json = raw.replace(/^```(?:json)?\n?/, '').replace(/\n?```$/, '').trim()
    try {
      return JSON.parse(json)
    } catch {
      throw new Error(`LLM returned invalid JSON. Try a narrower date range.\n\nRaw response:\n${raw.slice(0, 300)}…`)
    }
  })
}
