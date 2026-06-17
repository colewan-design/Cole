import { defineStore } from 'pinia'
import { ref, watchEffect } from 'vue'

const STORAGE_KEY = 'cole-report-state'

export const useReportStore = defineStore('report', () => {
  const filePath     = ref('')
  const fileName     = ref('')
  const rawContent   = ref('')
  const dateStart    = ref('')
  const dateEnd      = ref('')
  const reportType   = ref('AR')
  const employeeType = ref('casual')
  const parsedData   = ref(null)
  const outputPath   = ref('')

  // Hydrate from localStorage on init
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const s = JSON.parse(saved)
      if (s.filePath)     filePath.value     = s.filePath
      if (s.fileName)     fileName.value     = s.fileName
      if (s.rawContent)   rawContent.value   = s.rawContent
      if (s.dateStart)    dateStart.value    = s.dateStart
      if (s.dateEnd)      dateEnd.value      = s.dateEnd
      if (s.reportType)   reportType.value   = s.reportType
      if (s.employeeType) employeeType.value = s.employeeType
      if (s.parsedData)   parsedData.value   = s.parsedData
    }
  } catch {}

  // Auto-save whenever state changes
  watchEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        filePath:     filePath.value,
        fileName:     fileName.value,
        rawContent:   rawContent.value,
        dateStart:    dateStart.value,
        dateEnd:      dateEnd.value,
        reportType:   reportType.value,
        employeeType: employeeType.value,
        parsedData:   parsedData.value,
      }))
    } catch {}
  })

  function setFile(info) {
    filePath.value   = info.filePath
    fileName.value   = info.fileName
    rawContent.value = info.content
    parsedData.value = null
  }

  function setDateRange(start, end) {
    dateStart.value  = start
    dateEnd.value    = end
    parsedData.value = null
  }

  function setReportType(type) {
    reportType.value = type
    parsedData.value = null
  }

  function setEmployeeType(type) {
    employeeType.value = type
  }

  function setParsedData(data) {
    parsedData.value = data
  }

  function removeEntry(index) {
    if (!parsedData.value?.entries) return
    parsedData.value = {
      ...parsedData.value,
      entries: parsedData.value.entries.filter((_, i) => i !== index),
    }
  }

  function reset() {
    filePath.value     = ''
    fileName.value     = ''
    rawContent.value   = ''
    dateStart.value    = ''
    dateEnd.value      = ''
    reportType.value   = 'AR'
    employeeType.value = 'casual'
    parsedData.value   = null
    outputPath.value   = ''
  }

  return {
    filePath, fileName, rawContent,
    dateStart, dateEnd, reportType, employeeType,
    parsedData, outputPath,
    setFile, setDateRange, setReportType, setEmployeeType, setParsedData, removeEntry, reset,
  }
})
