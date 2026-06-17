<template>
  <div class="table-wrap">
    <!-- Panel header row -->
    <div class="panel-header">
      <span class="panel-meta">Total Entries: <strong>{{ rows.length }}</strong></span>
      <span class="panel-meta">Total Hours: <strong>{{ totalHours }}h</strong></span>
    </div>

    <!-- Table -->
    <div class="table-scroll">
      <table class="tbl">
        <thead>
          <tr>
            <th class="th th--date">Date</th>
            <th class="th th--tasks">Tasks</th>
            <th class="th th--hours">Hours</th>
            <th class="th th--status">Status</th>
            <th class="th th--actions">Actions</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="(row, i) in pageRows" :key="globalIndex(i)">
            <tr
              class="trow"
              :class="{ 'trow--expanded': expandedRow === globalIndex(i) }"
            >
              <td class="td td--date">{{ row.date }}</td>
              <td class="td td--tasks">
                <!-- Collapsed: first task only -->
                <div v-if="expandedRow !== globalIndex(i)" class="task-preview">
                  <span class="task-dot" />
                  <span class="task-text">{{ row.tasks?.[0] || '—' }}</span>
                  <span v-if="(row.tasks?.length ?? 0) > 1" class="task-more">
                    +{{ row.tasks.length - 1 }} more
                  </span>
                </div>
                <!-- Expanded: all tasks -->
                <ul v-else class="task-list">
                  <li v-for="(t, j) in row.tasks" :key="j" class="task-item">
                    <span class="task-dot" />
                    <span class="task-text selectable">{{ t }}</span>
                  </li>
                </ul>
              </td>
              <td class="td td--hours">{{ row.hoursWorked }}h</td>
              <td class="td td--status">
                <span class="badge badge--success">Logged</span>
              </td>
              <td class="td td--actions">
                <!-- View / collapse toggle -->
                <button
                  class="act-btn act-btn--outline"
                  :title="expandedRow === globalIndex(i) ? 'Collapse' : 'Expand tasks'"
                  @click="toggleExpand(globalIndex(i))"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path v-if="expandedRow === globalIndex(i)" d="M18 15l-6-6-6 6"/>
                    <path v-else d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle v-if="expandedRow !== globalIndex(i)" cx="12" cy="12" r="3"/>
                  </svg>
                </button>
                <!-- Remove -->
                <button
                  class="act-btn act-btn--danger"
                  title="Remove entry"
                  @click="$emit('remove', globalIndex(i))"
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6"/>
                    <path d="M10 11v6M14 11v6"/>
                  </svg>
                </button>
              </td>
            </tr>
          </template>

          <!-- Empty state -->
          <tr v-if="rows.length === 0">
            <td colspan="5" class="empty-cell">
              <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="empty-icon">
                  <path d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
                </svg>
                <span>No entries</span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Footer: total row + pagination -->
    <div class="table-footer">
      <span class="footer-total">Total: <strong>{{ totalHours }}h</strong></span>
      <div class="pagination">
        <button
          class="pag-btn"
          :disabled="currentPage === 1"
          @click="currentPage--"
        >Previous</button>

        <button
          v-for="p in totalPages"
          :key="p"
          class="pag-btn pag-btn--num"
          :class="{ 'pag-btn--active': p === currentPage }"
          @click="currentPage = p"
        >{{ p }}</button>

        <button
          class="pag-btn"
          :disabled="currentPage === totalPages"
          @click="currentPage++"
        >Next</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({ rows: { type: Array, default: () => [] } })
defineEmits(['remove'])

const PAGE_SIZE = 8

const currentPage  = ref(1)
const expandedRow  = ref(null)

const totalPages = computed(() => Math.max(1, Math.ceil(props.rows.length / PAGE_SIZE)))
const pageRows   = computed(() => {
  const start = (currentPage.value - 1) * PAGE_SIZE
  return props.rows.slice(start, start + PAGE_SIZE)
})
const totalHours = computed(() =>
  props.rows.reduce((sum, r) => sum + (Number(r.hoursWorked) || 0), 0)
)

function globalIndex(pageLocalIndex) {
  return (currentPage.value - 1) * PAGE_SIZE + pageLocalIndex
}

function toggleExpand(idx) {
  expandedRow.value = expandedRow.value === idx ? null : idx
}
</script>

<style scoped>
.table-wrap {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: var(--shadow-card);
}

/* Panel header */
.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-surface-2);
}
.panel-meta {
  font-size: 11px;
  color: var(--color-text-muted);
}
.panel-meta strong {
  color: var(--color-text);
  font-weight: 600;
}

/* Table */
.table-scroll { overflow-x: auto; }
.tbl { width: 100%; border-collapse: collapse; }

/* Header */
.th {
  padding: 10px 16px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  background: var(--color-surface-2);
  border-bottom: 1px solid var(--color-border);
  text-align: left;
  white-space: nowrap;
  user-select: none;
}
.th--hours, .th--status, .th--actions { text-align: center; }

/* Column widths */
.th--date, .td--date { width: 110px; }
.th--hours, .td--hours { width: 72px; }
.th--status, .td--status { width: 96px; }
.th--actions, .td--actions { width: 100px; }

/* Rows */
.trow {
  border-bottom: 1px solid var(--color-border);
  transition: background 0.1s;
}
.trow:last-child { border-bottom: none; }
.trow:hover { background: var(--color-surface-2); }
.trow--expanded { background: rgba(0,122,255,0.03); }

.td {
  padding: 12px 16px;
  vertical-align: middle;
  font-size: 13px;
  color: var(--color-text);
}
.td--date {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: var(--color-text-muted);
  white-space: nowrap;
}
.td--hours {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  font-weight: 600;
  color: #007AFF;
  text-align: center;
}
.td--status { text-align: center; }
.td--actions { text-align: center; }

/* Task preview (collapsed) */
.task-preview {
  display: flex;
  align-items: baseline;
  gap: 6px;
}
.task-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #007AFF;
  flex-shrink: 0;
  margin-top: 2px;
}
.task-text {
  font-size: 12px;
  color: var(--color-text);
  line-height: 1.4;
}
.task-more {
  font-size: 11px;
  color: var(--color-text-muted);
  white-space: nowrap;
  flex-shrink: 0;
}

/* Task list (expanded) */
.task-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.task-item {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

/* Status badge */
.badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 8px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.01em;
  white-space: nowrap;
}
.badge--success {
  background: rgba(52,199,89,0.12);
  color: #34C759;
  border: 1px solid rgba(52,199,89,0.25);
}
.badge--danger {
  background: rgba(255,59,48,0.1);
  color: #FF3B30;
  border: 1px solid rgba(255,59,48,0.2);
}

/* Action buttons */
.act-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  transition: background 0.12s, color 0.12s;
}
.act-btn svg { width: 14px; height: 14px; }

.act-btn--outline {
  background: transparent;
  border: 1px solid var(--color-border);
  color: var(--color-text-muted);
  margin-right: 6px;
}
.act-btn--outline:hover {
  background: var(--color-surface-2);
  color: #007AFF;
  border-color: #007AFF;
}

.act-btn--danger {
  background: rgba(255,59,48,0.1);
  color: #FF3B30;
}
.act-btn--danger:hover {
  background: rgba(255,59,48,0.18);
}

/* Empty state */
.empty-cell { padding: 40px; }
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: var(--color-text-muted);
  font-size: 13px;
}
.empty-icon { width: 36px; height: 36px; }

/* Footer */
.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  border-top: 1px solid var(--color-border);
  background: var(--color-surface-2);
}
.footer-total {
  font-size: 12px;
  color: var(--color-text-muted);
}
.footer-total strong {
  color: #007AFF;
  font-weight: 700;
}

/* Pagination */
.pagination {
  display: flex;
  align-items: center;
  gap: 4px;
}
.pag-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 28px;
  padding: 0 10px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  background: var(--color-surface);
  color: var(--color-text-muted);
  font-size: 12px;
  cursor: pointer;
  transition: background 0.12s, color 0.12s, border-color 0.12s;
}
.pag-btn:hover:not(:disabled) {
  background: var(--color-surface-2);
  color: var(--color-text);
}
.pag-btn:disabled {
  opacity: 0.35;
  cursor: default;
}
.pag-btn--num { min-width: 28px; padding: 0; }
.pag-btn--active {
  background: #007AFF;
  border-color: #007AFF;
  color: #fff;
  font-weight: 600;
}
.pag-btn--active:hover { background: #0071e3; }
</style>
