<template>
  <div class="shell">
    <!-- Sidebar -->
    <aside class="sidebar">
      <!-- Electron drag / title-bar region -->
      <div class="drag-zone app-drag">
        <span class="drag-label app-no-drag">COLE</span>
      </div>

      <!-- Logo -->
      <div class="sidebar-logo">
        <img src="@/assets/logo.png" alt="Cole" class="logo-img" />
        <div>
          <div class="logo-name">Report Assistant</div>
          <div class="logo-ver">v0.1 · Phase 1</div>
        </div>
      </div>

      <!-- Main nav -->
      <nav class="sidebar-nav">
        <router-link
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="nav-item"
          active-class="nav-item--active"
        >
          <component :is="item.icon" class="nav-icon" />
          <span class="nav-label">{{ item.label }}</span>
        </router-link>
      </nav>

      <div class="sidebar-spacer" />

      <!-- Bottom nav -->
      <nav class="sidebar-nav sidebar-nav--bottom">
        <router-link
          v-for="item in navBottom"
          :key="item.to"
          :to="item.to"
          class="nav-item"
          active-class="nav-item--active"
        >
          <component :is="item.icon" class="nav-icon" />
          <span class="nav-label">{{ item.label }}</span>
        </router-link>
      </nav>

      <!-- Active file card -->
      <div v-if="report.fileName" class="file-card">
        <div class="file-card-label">Active file</div>
        <div class="file-card-name">{{ report.fileName }}</div>
        <div v-if="report.dateStart" class="file-card-dates">
          {{ report.dateStart }} → {{ report.dateEnd }}
        </div>
      </div>
    </aside>

    <!-- Right panel: navbar + content -->
    <div class="right-panel">
      <!-- Top navbar -->
      <header class="navbar app-drag">
        <div class="navbar-left app-no-drag">
          <h1 class="page-title">{{ pageTitle }}</h1>
        </div>
        <div class="navbar-right app-no-drag">
          <!-- Provider chip -->
          <div class="provider-chip">
            <span class="provider-dot" :class="providerDotClass" />
            {{ providerLabel }}
          </div>
          <!-- User -->
          <div class="user-info">
            <span v-if="settings.reporterNameSig" class="user-name">
              {{ settings.reporterNameSig }}
            </span>
            <div class="user-avatar">{{ avatarInitials }}</div>
          </div>
        </div>
      </header>

      <!-- Router content -->
      <main class="content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { h, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useReportStore } from '@/stores/report.js'
import { useSettingsStore } from '@/stores/settings.js'

const report   = useReportStore()
const settings = useSettingsStore()
const route    = useRoute()

onMounted(() => settings.load())

const pageTitle = computed(() => route.meta?.title ?? 'Cole')

const providerLabel = computed(() => {
  if (settings.provider === 'ollama') return settings.ollamaModel || 'Ollama'
  if (settings.provider === 'gemini') return 'Gemini'
  return 'Claude'
})

const providerDotClass = computed(() => {
  if (settings.provider === 'ollama') return 'dot--orange'
  if (settings.provider === 'gemini') return 'dot--blue'
  return settings.hasApiKey ? 'dot--green' : 'dot--red'
})

const avatarInitials = computed(() => {
  const name = settings.reporterNameSig || settings.reporterName || ''
  if (!name) return 'U'
  return name.split(/\s+/).slice(0, 2).map(w => w[0] ?? '').join('').toUpperCase()
})

// ── Inline SVG icon components ────────────────────────────────────────────
const IconHome = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('path', { d:'M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2z' }),
  h('polyline', { points:'9 22 9 12 15 12 15 22' }),
]) }
const IconEye = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('path', { d:'M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z' }),
  h('circle', { cx:'12', cy:'12', r:'3' }),
]) }
const IconSettings = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('circle', { cx:'12', cy:'12', r:'3' }),
  h('path', { d:'M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z' }),
]) }
const IconHelp = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('circle', { cx:'12', cy:'12', r:'10' }),
  h('path', { d:'M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3' }),
  h('line', { x1:'12', y1:'17', x2:'12.01', y2:'17' }),
]) }
const IconInfo = { render: () => h('svg', { viewBox:'0 0 24 24', fill:'none', stroke:'currentColor', 'stroke-width':'2' }, [
  h('circle', { cx:'12', cy:'12', r:'10' }),
  h('line', { x1:'12', y1:'8', x2:'12', y2:'12' }),
  h('line', { x1:'12', y1:'16', x2:'12.01', y2:'16' }),
]) }

const nav = [
  { to: '/',         label: 'Generate',   icon: IconHome },
  { to: '/preview',  label: 'Preview',    icon: IconEye },
  { to: '/settings', label: 'Settings',   icon: IconSettings },
]
const navBottom = [
  { to: '/help',  label: 'Help & FAQs', icon: IconHelp },
  { to: '/about', label: 'About',       icon: IconInfo },
]
</script>

<style scoped>
/* ── Shell ─────────────────────────────────────────────────────────────── */
.shell {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: var(--color-bg);
}

/* ── Sidebar ────────────────────────────────────────────────────────────── */
.sidebar {
  width: 200px;
  flex-shrink: 0;
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-sm);
}

.drag-zone {
  height: 36px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  border-bottom: 1px solid var(--color-border);
}
.drag-label {
  font-size: 10px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  color: var(--color-text-muted);
  letter-spacing: 0.2em;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--color-border);
}
.logo-img {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  object-fit: contain;
}
.logo-name {
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text);
  line-height: 1.2;
}
.logo-ver {
  font-size: 10px;
  color: var(--color-text-muted);
  margin-top: 1px;
}

.sidebar-nav {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.sidebar-nav--bottom {
  border-top: 1px solid var(--color-border);
}
.sidebar-spacer { flex: 1; }

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px 8px 13px;
  border-radius: 10px;
  font-size: 13px;
  color: var(--color-text-muted);
  text-decoration: none;
  transition: background 0.12s, color 0.12s;
  border-left: 3px solid transparent;
}
.nav-item:hover {
  background: var(--color-surface-2);
  color: var(--color-text);
}
.nav-item--active {
  background: var(--color-surface-2);
  color: #007AFF;
  border-left-color: #007AFF;
  font-weight: 500;
}
.nav-icon { width: 16px; height: 16px; flex-shrink: 0; }
.nav-label { flex: 1; }

.file-card {
  margin: 8px;
  padding: 10px 12px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 12px;
}
.file-card-label {
  font-size: 10px;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 3px;
}
.file-card-name {
  font-size: 11px;
  color: #007AFF;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.file-card-dates {
  font-size: 10px;
  color: var(--color-text-muted);
  margin-top: 2px;
}

/* ── Right panel ────────────────────────────────────────────────────────── */
.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ── Navbar ─────────────────────────────────────────────────────────────── */
.navbar {
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: var(--color-surface);
  border-bottom: 1px solid var(--color-border);
  flex-shrink: 0;
}
.navbar-left { display: flex; align-items: center; }
.page-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
  letter-spacing: -0.01em;
}
.navbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.provider-chip {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  background: var(--color-surface-2);
  border: 1px solid var(--color-border);
  border-radius: 20px;
  font-size: 11px;
  color: var(--color-text-muted);
  font-weight: 500;
}
.provider-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}
.dot--green  { background: #34C759; }
.dot--red    { background: #FF3B30; }
.dot--blue   { background: #007AFF; }
.dot--orange { background: #FF9500; }

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}
.user-name {
  font-size: 12px;
  color: var(--color-text-muted);
  max-width: 130px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.user-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: #007AFF;
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* ── Content ────────────────────────────────────────────────────────────── */
.content {
  flex: 1;
  overflow: auto;
}
</style>
