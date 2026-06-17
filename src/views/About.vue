<template>
  <div class="p-6 max-w-xl mx-auto">

    <!-- Logo + name -->
    <div class="flex flex-col items-center text-center mb-8 mt-4">
      <img src="@/assets/logo.png" alt="Cole" class="w-20 h-20 rounded-2xl object-contain mb-4 shadow-glow" />
      <h1 class="text-2xl font-bold text-white">Cole</h1>
      <p class="text-sm text-j-muted mt-1">AI Report Assistant</p>
      <span class="mt-2 text-xs font-mono bg-j-card border border-j-border text-j-glow px-2.5 py-1 rounded-full">v0.1.0</span>
    </div>

    <!-- Description -->
    <div class="card mb-4">
      <p class="text-sm text-j-text leading-relaxed">
        Cole is a desktop app that turns your Obsidian daily work logs into formatted government accomplishment reports (.docx) using AI. It supports Claude, Gemini, Ollama, and a manual paste workflow — so you stay in control of your API costs.
      </p>
    </div>

    <!-- Stack -->
    <div class="card mb-4">
      <h2 class="section-label mb-3">Built with</h2>
      <div class="grid grid-cols-2 gap-2">
        <div v-for="tech in stack" :key="tech.name"
             class="flex items-center gap-2.5 bg-j-surface border border-j-border rounded-lg px-3 py-2">
          <i :class="['fi', tech.icon, 'text-base text-j-glow']" />
          <div>
            <div class="text-xs font-semibold text-white">{{ tech.name }}</div>
            <div class="text-[10px] text-j-muted">{{ tech.role }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Developer -->
    <div class="card mb-4">
      <h2 class="section-label mb-3">Developer</h2>
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-full bg-j-accent/20 border border-j-accent/30 flex items-center justify-center">
          <i class="fi fi-rr-user text-j-glow text-sm" />
        </div>
        <div>
          <div class="text-sm font-semibold text-white">Christian Fiarawe Colewan</div>
          <div class="text-xs text-j-muted">Information System Analyst I · BSU CBOO</div>
        </div>
      </div>
    </div>

    <!-- Links -->
    <div class="card">
      <h2 class="section-label mb-3">Links</h2>
      <div class="space-y-2">
        <a @click.prevent="openLink('https://github.com/colewan-design/Cole')"
           href="#"
           class="flex items-center gap-2 text-sm text-j-muted hover:text-white transition-colors">
          <i class="fi fi-rr-code-branch text-j-glow" />
          github.com/colewan-design/Cole
        </a>
        <a @click.prevent="openLink('https://claude.ai')"
           href="#"
           class="flex items-center gap-2 text-sm text-j-muted hover:text-white transition-colors">
          <i class="fi fi-rr-brain text-j-glow" />
          claude.ai — manual parse workflow
        </a>
      </div>
    </div>

    <p class="text-center text-xs text-j-muted mt-6">© 2026 Christian Colewan · BSU CBOO</p>
  </div>
</template>

<script setup>
const stack = [
  { name: 'Electron',    icon: 'fi-rr-computer',      role: 'Desktop shell' },
  { name: 'Vue 3',       icon: 'fi-rr-layers',        role: 'UI framework' },
  { name: 'Tailwind',    icon: 'fi-rr-palette',       role: 'Styling' },
  { name: 'Claude API',  icon: 'fi-rr-robot',         role: 'AI parsing' },
  { name: 'docx',        icon: 'fi-rr-document',      role: 'Report generation' },
  { name: 'Pinia',       icon: 'fi-rr-database',      role: 'State management' },
]

function openLink(url) {
  window.cole.openExternal(url)
}
</script>

<style scoped>
.section-label { @apply text-xs font-semibold uppercase tracking-wider text-j-muted; }
.card          { @apply bg-j-card border border-j-border rounded-xl p-4; }
</style>
