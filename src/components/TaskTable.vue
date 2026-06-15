<template>
  <div class="overflow-auto rounded-xl border border-j-border">
    <table class="w-full text-sm">
      <thead>
        <tr class="bg-j-surface border-b border-j-border">
          <th class="text-left px-4 py-2.5 text-xs font-semibold text-j-muted uppercase tracking-wider">Date</th>
          <th class="text-left px-4 py-2.5 text-xs font-semibold text-j-muted uppercase tracking-wider">Tasks</th>
          <th class="text-right px-4 py-2.5 text-xs font-semibold text-j-muted uppercase tracking-wider">Hours</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(row, i) in rows"
          :key="i"
          class="border-b border-j-border/50 hover:bg-j-card/50 transition-colors"
        >
          <td class="px-4 py-3 font-mono text-xs text-j-muted whitespace-nowrap">{{ row.date }}</td>
          <td class="px-4 py-3">
            <ul class="space-y-0.5">
              <li v-for="(task, j) in row.tasks" :key="j" class="flex items-start gap-2">
                <span class="mt-1.5 w-1 h-1 rounded-full bg-j-accent flex-shrink-0" />
                <span class="text-j-text text-xs">{{ task }}</span>
              </li>
            </ul>
          </td>
          <td class="px-4 py-3 text-right font-mono text-xs text-j-glow whitespace-nowrap">
            {{ row.hoursWorked }}h
          </td>
        </tr>
      </tbody>
      <tfoot>
        <tr class="bg-j-surface">
          <td class="px-4 py-2.5 text-xs font-semibold text-j-muted" colspan="2">Total</td>
          <td class="px-4 py-2.5 text-right font-mono text-sm font-bold text-j-glow">
            {{ totalHours }}h
          </td>
        </tr>
      </tfoot>
    </table>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ rows: { type: Array, default: () => [] } })

const totalHours = computed(() =>
  props.rows.reduce((sum, r) => sum + (Number(r.hoursWorked) || 0), 0)
)
</script>
