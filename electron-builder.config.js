export default {
  appId: 'com.cole.report-assistant',
  productName: 'Cole',
  directories: {
    output: 'release',
  },
  files: [
    'dist/**/*',
    'dist-electron/**/*',
    'public/**/*',
    'resources/**/*',
  ],
  win: {
    target: [{ target: 'nsis', arch: ['x64'] }],
    icon: 'public/icon.ico',
  },
  nsis: {
    oneClick: false,
    allowToChangeInstallationDirectory: true,
    createDesktopShortcut: true,
    createStartMenuShortcut: true,
  },
}
