import { dialog, shell } from 'electron'
import { readFile, writeFile } from 'node:fs/promises'
import path from 'node:path'

export function registerFileHandlers(ipcMain) {
  ipcMain.handle('file:open-dialog', async (_, options = {}) => {
    const result = await dialog.showOpenDialog({
      properties: ['openFile'],
      filters: [{ name: 'Markdown Files', extensions: ['md'] }],
      ...options,
    })
    if (result.canceled) return null
    return result.filePaths[0]
  })

  ipcMain.handle('file:read', async (_, filePath) => {
    const content = await readFile(filePath, 'utf-8')
    return {
      content,
      filePath,
      fileName: path.basename(filePath),
      dirName: path.dirname(filePath),
    }
  })

  ipcMain.handle('file:save', async (_, filePath, buffer) => {
    await writeFile(filePath, Buffer.from(buffer))
    return filePath
  })

  ipcMain.handle('file:save-dialog', async (_, options = {}) => {
    const result = await dialog.showSaveDialog({
      filters: [{ name: 'Word Document', extensions: ['docx'] }],
      ...options,
    })
    if (result.canceled) return null
    return result.filePath
  })

  ipcMain.handle('file:open-path', async (_, filePath) => {
    await shell.openPath(filePath)
  })
}
