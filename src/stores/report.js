import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useReportStore = defineStore('report', () => {
  const filePath    = ref('')
  const fileName    = ref('')
  const rawContent  = ref('')
  const dateStart   = ref('')
  const dateEnd     = ref('')
  const reportType  = ref('AR')       // 'AR' | 'DTR' | 'PRG'
  const parsedData  = ref(null)       // JSON from Claude (Phase 2)
  const outputPath  = ref('')

  function setFile(info) {
    filePath.value   = info.filePath
    fileName.value   = info.fileName
    rawContent.value = info.content
    parsedData.value = null
  }

  function setDateRange(start, end) {
    dateStart.value = start
    dateEnd.value   = end
    parsedData.value = null
  }

  function setReportType(type) {
    reportType.value = type
    parsedData.value = null
  }

  function setParsedData(data) {
    parsedData.value = data
  }

  function reset() {
    filePath.value   = ''
    fileName.value   = ''
    rawContent.value = ''
    dateStart.value  = ''
    dateEnd.value    = ''
    reportType.value = 'AR'
    parsedData.value = null
    outputPath.value = ''
  }

  return {
    filePath, fileName, rawContent,
    dateStart, dateEnd, reportType,
    parsedData, outputPath,
    setFile, setDateRange, setReportType, setParsedData, reset,
  }
})
