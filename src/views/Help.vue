<template>
  <div class="p-6 max-w-xl mx-auto">
    <div class="mb-6">
      <h1 class="text-xl font-semibold text-white">Help & FAQs</h1>
      <p class="text-sm text-j-muted mt-1">How to use Cole and answers to common questions.</p>
    </div>

    <!-- Quick start -->
    <section class="mb-5">
      <h2 class="section-label mb-3">Quick Start</h2>
      <div class="card space-y-0">
        <div v-for="(step, i) in steps" :key="i"
             class="flex gap-3 py-3" :class="i < steps.length - 1 ? 'border-b border-j-border' : ''">
          <div class="w-6 h-6 rounded-full bg-j-accent flex items-center justify-center text-xs font-bold text-white flex-shrink-0 mt-0.5">
            {{ i + 1 }}
          </div>
          <div>
            <div class="text-sm font-medium text-white">{{ step.title }}</div>
            <div class="text-xs text-j-muted mt-0.5 leading-relaxed">{{ step.desc }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- FAQs -->
    <section class="mb-5">
      <h2 class="section-label mb-3">Frequently Asked Questions</h2>
      <div class="space-y-2">
        <div v-for="(faq, i) in faqs" :key="i" class="card">
          <button @click="open === i ? open = null : open = i"
                  class="flex items-start justify-between gap-3 w-full text-left">
            <span class="text-sm font-medium text-white leading-snug">{{ faq.q }}</span>
            <i :class="['fi fi-rr-angle-down text-j-muted transition-transform flex-shrink-0 mt-0.5', open === i ? 'rotate-180' : '']" />
          </button>
          <div v-if="open === i" class="mt-2.5 text-xs text-j-muted leading-relaxed border-t border-j-border pt-2.5">
            {{ faq.a }}
          </div>
        </div>
      </div>
    </section>

    <!-- Providers -->
    <section class="mb-5">
      <h2 class="section-label mb-3">AI Provider Comparison</h2>
      <div class="card overflow-hidden p-0">
        <table class="w-full text-xs">
          <thead>
            <tr class="border-b border-j-border">
              <th class="text-left px-4 py-2.5 text-j-muted font-semibold">Provider</th>
              <th class="text-left px-4 py-2.5 text-j-muted font-semibold">Cost</th>
              <th class="text-left px-4 py-2.5 text-j-muted font-semibold">Quality</th>
              <th class="text-left px-4 py-2.5 text-j-muted font-semibold">Internet</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in providers" :key="p.name"
                class="border-b border-j-border last:border-0">
              <td class="px-4 py-2.5 text-white font-medium">{{ p.name }}</td>
              <td class="px-4 py-2.5 text-j-muted">{{ p.cost }}</td>
              <td class="px-4 py-2.5">
                <span :class="['px-1.5 py-0.5 rounded text-[10px] font-semibold', p.qualityClass]">{{ p.quality }}</span>
              </td>
              <td class="px-4 py-2.5 text-j-muted">{{ p.internet }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <!-- Troubleshooting -->
    <section>
      <h2 class="section-label mb-3">Troubleshooting</h2>
      <div class="space-y-2">
        <div v-for="t in troubles" :key="t.issue" class="card flex gap-3">
          <i class="fi fi-rr-triangle-warning text-yellow-400 flex-shrink-0 mt-0.5" />
          <div>
            <div class="text-sm font-medium text-white">{{ t.issue }}</div>
            <div class="text-xs text-j-muted mt-0.5 leading-relaxed">{{ t.fix }}</div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const open = ref(null)

const steps = [
  { title: 'Select your Obsidian note',    desc: 'Click the file zone or drag your .md file onto it.' },
  { title: 'Set the date range',           desc: 'Choose the start and end dates for the report period.' },
  { title: 'Choose a report type',         desc: 'AR = Accomplishment Report, DTR = Daily Time Record, PRG = Progress Report.' },
  { title: 'Parse the note',               desc: 'Click "Parse with Claude/Gemini/Ollama" or use the free Claude.ai paste workflow below.' },
  { title: 'Review & export',              desc: 'Check the parsed entries in the Preview tab, then click "Export .docx" to save your report.' },
]

const faqs = [
  {
    q: 'Do I need an API key?',
    a: 'No. You can use the free "Paste JSON from Claude.ai" workflow: copy the prompt, paste your note on claude.ai, then paste the JSON back into Cole. No API key needed.',
  },
  {
    q: 'What is the Claude.ai paste workflow?',
    a: 'Open the "Paste JSON from Claude.ai" section on the Generate page. Click "Copy prompt ↗" — this copies the full prompt with your note content included. Paste it into claude.ai, copy the JSON it returns, paste it back into Cole, and click "Generate .docx from pasted JSON".',
  },
  {
    q: 'How do I reduce API token costs?',
    a: 'Switch to the Ollama provider in Settings (free, runs locally). Or use Gemini free tier. Or use the Claude.ai paste workflow entirely. Prompt caching is already enabled for Claude API calls to reduce repeat costs.',
  },
  {
    q: 'Why does the .docx have blank fields?',
    a: 'Go to Settings → Report Profile and click "Save Profile". Your name, position, office, and supervisor details must be saved before they appear in generated documents.',
  },
  {
    q: 'Why are some tasks missing from the report?',
    a: 'The AI filters entries by your selected date range. If a task has no date or uses an unrecognized format, it may be excluded. Check that dates in your .md file match common formats (June 1, 2026-06-01, etc.).',
  },
  {
    q: 'Ollama returns an error about context size',
    a: 'The note is too large for the default model context. Cole automatically sets num_ctx: 16384 for Ollama. If it still fails, try a smaller date range to reduce the input size.',
  },
  {
    q: 'The parsed JSON looks wrong',
    a: 'Try narrowing the date range. Large notes with many entries can confuse smaller models. For best results use Claude API or the Claude.ai paste workflow.',
  },
  {
    q: 'Can I edit the parsed data before exporting?',
    a: 'Not yet — editing is a planned feature. For now, you can manually edit the JSON in the "Paste JSON" textarea before clicking Generate, or edit the .docx after export.',
  },
]

const providers = [
  { name: 'Claude API',   cost: 'Pay per token',  quality: 'Best',   qualityClass: 'bg-green-900/50 text-green-400',  internet: 'Required' },
  { name: 'Claude.ai',    cost: 'Free (Pro plan)', quality: 'Best',   qualityClass: 'bg-green-900/50 text-green-400',  internet: 'Required' },
  { name: 'Gemini',       cost: 'Free tier',       quality: 'Good',   qualityClass: 'bg-blue-900/50 text-blue-400',    internet: 'Required' },
  { name: 'Ollama',       cost: 'Free',            quality: 'Varies', qualityClass: 'bg-yellow-900/50 text-yellow-400', internet: 'No' },
]

const troubles = [
  { issue: 'No models shown in Ollama dropdown',     fix: 'Make sure Ollama is running (ollama serve) before opening Settings. The app fetches the model list on load.' },
  { issue: 'Gemini key shows "No key saved"',        fix: 'Restart the app after saving the key — the main process needs a full restart to pick up new settings.' },
  { issue: '"Invalid JSON" when pasting from Claude.ai', fix: 'Cole auto-extracts the JSON block from surrounding text. If it still fails, manually copy only the { ... } block.' },
  { issue: 'Logo or BSU logo missing in .docx',      fix: 'The BSU logo is loaded from resources/bsu-logo.png inside the app package. Make sure the file exists at that path.' },
]
</script>

<style scoped>
.section-label { @apply text-xs font-semibold uppercase tracking-wider text-j-muted mb-3; }
.card          { @apply bg-j-card border border-j-border rounded-xl p-4; }
</style>
