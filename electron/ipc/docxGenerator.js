import {
  Document, Packer, Paragraph, Table, TableRow, TableCell, TextRun,
  AlignmentType, WidthType, BorderStyle, VerticalAlign, ImageRun, Header,
} from 'docx'
import { writeFile, readFile } from 'node:fs/promises'
import path from 'node:path'

// ─── Constants (matching the original template measurements) ─────────────────

const LOGO_PATH  = path.join(process.env.APP_ROOT || '', 'resources', 'bsu-logo.png')
const FONT_TITLE = 'Copperplate Gothic Bold'
const FONT_DAY   = 'Cambria'
const FONT_TASK  = 'Arial Narrow'
const FONT_HDR   = 'Calibri'    // header info paragraphs

const SZ_TITLE = 20   // 10pt
const SZ_HDR   = 20   // 10pt
const SZ_DAY   = 20   // 10pt
const SZ_TASK  = 20   // 10pt

// Column widths in DXA (twips) from the original document
const COL_DAY   = 729
const COL_TASK  = 9859
const TBL_WIDTH = COL_DAY + COL_TASK   // 10588

// Page margins in twips (from original sectPr)
const MARGIN = { top: 1440, right: 720, bottom: 720, left: 720 }
// A4 in twips
const PAGE = { width: 11906, height: 16838 }

const MONTHS = ['January','February','March','April','May','June',
                'July','August','September','October','November','December']
const DAYS_OF_WEEK = ['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday']

// ─── Helpers ─────────────────────────────────────────────────────────────────

function singleBorder(color = '000000', size = 4) {
  const s = { style: BorderStyle.SINGLE, size, color }
  return { top: s, bottom: s, left: s, right: s }
}

function noBorder() {
  const s = { style: BorderStyle.NONE, size: 0, color: 'FFFFFF' }
  return { top: s, bottom: s, left: s, right: s }
}

function t(text, opts = {}) {
  return new TextRun({ text: String(text ?? ''), font: FONT_TASK, size: SZ_TASK, ...opts })
}

function formatPeriod(dateStart, dateEnd) {
  const s = new Date(dateStart), e = new Date(dateEnd)
  const sm = MONTHS[s.getMonth()], em = MONTHS[e.getMonth()]
  const sd = s.getDate(), ed = e.getDate(), yr = e.getFullYear()
  return sm === em
    ? `${sm} ${sd}-${ed}, ${yr}`
    : `${sm} ${sd} - ${em} ${ed}, ${yr}`
}

function allDaysInRange(dateStart, dateEnd) {
  const days = []
  const cur  = new Date(dateStart)
  const end  = new Date(dateEnd)
  while (cur <= end) {
    days.push(new Date(cur))
    cur.setDate(cur.getDate() + 1)
  }
  return days
}

function toDateKey(d) {
  return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
}

function taskRuns(taskStr) {
  return [t(taskStr)]
}

function dottedBorder(color = '000000', size = 4) {
  const s = { style: BorderStyle.DOTTED, size, color }
  return { top: s, bottom: s, left: s, right: s }
}

// ─── Letterhead row: [Logo] [BSU text] [AR-HRMO form dashed box] ─────────────

async function letterheadRow() {
  let logoImage = null
  try {
    const logoData = await readFile(LOGO_PATH)
    logoImage = new ImageRun({ data: logoData, transformation: { width: 70, height: 70 }, type: 'png' })
  } catch (e) {
    console.error('Logo load failed:', e.message, LOGO_PATH)
  }

  const logoW   = Math.round(TBL_WIDTH * 0.15)
  const textW   = Math.round(TBL_WIDTH * 0.62)
  const spacerW = Math.round(TBL_WIDTH * 0.04)
  const formW   = TBL_WIDTH - logoW - textW - spacerW

  return new Table({
    width: { size: TBL_WIDTH, type: WidthType.DXA },
    borders: noBorder(),
    rows: [
      new TableRow({
        children: [
          // Logo — left
          new TableCell({
            width: { size: logoW, type: WidthType.DXA },
            borders: noBorder(),
            verticalAlign: VerticalAlign.CENTER,
            children: [new Paragraph({
              alignment: AlignmentType.CENTER,
              children: logoImage ? [logoImage] : [],
            })],
          }),
          // BSU text — center
          new TableCell({
            width: { size: textW, type: WidthType.DXA },
            borders: noBorder(),
            verticalAlign: VerticalAlign.CENTER,
            children: [
              new Paragraph({
                alignment: AlignmentType.CENTER,
                spacing: { before: 0, after: 0 },
                children: [new TextRun({ text: 'Republic of the Philippines', font: FONT_HDR, size: 18 })],
              }),
              new Paragraph({
                alignment: AlignmentType.CENTER,
                spacing: { before: 0, after: 0 },
                children: [new TextRun({ text: 'Benguet State University', font: 'Old English Text MT', size: 36, color: '1B5E20' })],
              }),
              new Paragraph({
                alignment: AlignmentType.CENTER,
                spacing: { before: 0, after: 0 },
                children: [new TextRun({ text: 'La Trinidad, Benguet, 2601', font: FONT_HDR, size: 18 })],
              }),
            ],
          }),
          // Spacer — pushes form box to the right edge
          new TableCell({
            width: { size: spacerW, type: WidthType.DXA },
            borders: noBorder(),
            children: [new Paragraph({ children: [] })],
          }),
          // AR-HRMO form — nested single-cell table; border on inner cell hugs text height
          new TableCell({
            width: { size: formW, type: WidthType.DXA },
            borders: noBorder(),
            verticalAlign: VerticalAlign.CENTER,
            children: [
              new Table({
                width: { size: formW - 40, type: WidthType.DXA },
                borders: noBorder(),
                rows: [
                  new TableRow({
                    children: [
                      new TableCell({
                        borders: {
                          top:    { style: BorderStyle.DOTTED, size: 4, color: '000000' },
                          bottom: { style: BorderStyle.DOTTED, size: 4, color: '000000' },
                          left:   { style: BorderStyle.DOTTED, size: 4, color: '000000' },
                          right:  { style: BorderStyle.DOTTED, size: 4, color: '000000' },
                        },
                        margins: { top: 60, bottom: 60, left: 80, right: 80 },
                        children: [new Paragraph({
                          alignment: AlignmentType.CENTER,
                          spacing: { before: 0, after: 0 },
                          children: [new TextRun({ text: 'AR-HRMO form 1.0', font: 'Brush Script MT', size: 18, italics: true })],
                        })],
                      }),
                    ],
                  }),
                ],
              }),
            ],
          }),
        ],
      }),
    ],
  })
}

// ─── Title box — centered on its own row ──────────────────────────────────────

function titleBoxSection() {
  const boxW = Math.round(TBL_WIDTH * 0.56)
  const padW = Math.floor((TBL_WIDTH - boxW) / 2)

  return new Table({
    width: { size: TBL_WIDTH, type: WidthType.DXA },
    borders: noBorder(),
    rows: [
      new TableRow({
        children: [
          new TableCell({
            width: { size: padW, type: WidthType.DXA },
            borders: noBorder(),
            children: [new Paragraph({ children: [] })],
          }),
          new TableCell({
            width: { size: boxW, type: WidthType.DXA },
            borders: singleBorder('000000', 6),
            verticalAlign: VerticalAlign.CENTER,
            children: [new Paragraph({
              alignment: AlignmentType.CENTER,
              spacing: { before: 60, after: 60 },
              children: [new TextRun({ text: 'Accomplishment Report for CASUAL Personnel', font: FONT_TITLE, size: SZ_TITLE })],
            })],
          }),
          new TableCell({
            width: { size: TBL_WIDTH - padW - boxW, type: WidthType.DXA },
            borders: noBorder(),
            children: [new Paragraph({ children: [] })],
          }),
        ],
      }),
    ],
  })
}

// ─── Header info paragraphs ───────────────────────────────────────────────────

function hRun(text, bold = false) {
  return new TextRun({ text, font: FONT_HDR, size: SZ_HDR, bold })
}

function headerParagraphs(period, userInfo) {
  const half = Math.floor(TBL_WIDTH / 2)

  const namePositionTable = new Table({
    width: { size: TBL_WIDTH, type: WidthType.DXA },
    borders: noBorder(),
    rows: [
      new TableRow({
        children: [
          new TableCell({
            width: { size: half, type: WidthType.DXA },
            borders: noBorder(),
            children: [new Paragraph({
              children: [hRun('NAME:  '), hRun(userInfo.reporterName || '______________________________', true)],
              spacing: { after: 0 },
            })],
          }),
          new TableCell({
            width: { size: TBL_WIDTH - half, type: WidthType.DXA },
            borders: noBorder(),
            children: [new Paragraph({
              children: [hRun('Position:  '), hRun(userInfo.position || '______________________________')],
              spacing: { after: 0 },
            })],
          }),
        ],
      }),
    ],
  })

  return [
    new Paragraph({
      alignment: AlignmentType.CENTER,
      spacing: { before: 60, after: 60 },
      children: [hRun('For the Period Covered: ', true), hRun(period, true)],
    }),
    namePositionTable,
    new Paragraph({
      spacing: { before: 0, after: 40 },
      children: [new TextRun({ text: '(Family Name, First Name, Middle Name)', font: FONT_HDR, size: 18, italics: true })],
      indent: { left: 720 },
    }),
    new Paragraph({
      spacing: { after: 60 },
      children: [
        new TextRun({ text: 'College/Office Assignment:  ', font: FONT_HDR, size: SZ_HDR, italics: true }),
        new TextRun({ text: userInfo.office || '______________________________', font: FONT_HDR, size: SZ_HDR, bold: true, italics: true }),
      ],
    }),
  ]
}

// ─── Main 2-column task table ─────────────────────────────────────────────────

function mainTable(parsedData) {
  const entryMap = {}
  for (const e of (parsedData.entries ?? [])) {
    entryMap[e.date] = e
  }

  const days = allDaysInRange(parsedData.dateStart, parsedData.dateEnd)

  const rows = days.map(d => {
    const dayNum  = d.getDate()
    const dow     = d.getDay()      // 0=Sun, 6=Sat
    const dateKey = toDateKey(d)
    const entry   = entryMap[dateKey]

    // Day number cell
    const dayCell = new TableCell({
      width: { size: COL_DAY, type: WidthType.DXA },
      verticalAlign: VerticalAlign.CENTER,
      children: [new Paragraph({
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: String(dayNum), font: FONT_DAY, size: SZ_DAY })],
      })],
    })

    // Task cell content
    let taskParas = []

    if (dow === 0 || dow === 6) {
      // Weekend
      taskParas = [new Paragraph({
        children: [t(DAYS_OF_WEEK[dow], { italics: true, color: '888888' })],
      })]
    } else if (entry && entry.tasks?.length) {
      // Work day with tasks — each task as its own paragraph
      taskParas = entry.tasks.map(taskStr =>
        new Paragraph({ children: taskRuns(taskStr), spacing: { after: 20 } })
      )
    } else {
      // Work day, no entry (holiday or data missing)
      taskParas = [new Paragraph({ children: [t('')] })]
    }

    const taskCell = new TableCell({
      width: { size: COL_TASK, type: WidthType.DXA },
      verticalAlign: VerticalAlign.CENTER,
      children: taskParas,
    })

    return new TableRow({
      children: [dayCell, taskCell],
    })
  })

  return new Table({
    width: { size: TBL_WIDTH, type: WidthType.DXA },
    borders: singleBorder('000000', 4),
    rows,
  })
}

// ─── Signature block ──────────────────────────────────────────────────────────

function signatureTable(userInfo) {
  const half = Math.floor(TBL_WIDTH / 2)

  function sigCell(children) {
    return new TableCell({
      width: { size: half, type: WidthType.DXA },
      borders: noBorder(),
      children,
    })
  }

  // Signature labels use Cambria (matches original); names/titles use Calibri (theme default)
  function labelRun(text) {
    return new TextRun({ text, font: FONT_DAY, size: SZ_HDR })
  }
  function nameRun(text, bold = false) {
    return new TextRun({ text, font: FONT_HDR, size: SZ_HDR, bold })
  }

  return new Table({
    width: { size: TBL_WIDTH, type: WidthType.DXA },
    borders: noBorder(),
    rows: [
      // Row 1: labels
      new TableRow({
        children: [
          sigCell([new Paragraph({ alignment: AlignmentType.CENTER, children: [labelRun('Reported by:')] })]),
          sigCell([new Paragraph({ alignment: AlignmentType.CENTER, children: [labelRun('Certified Correct:')] })]),
        ],
      }),
      // Row 2: blank signature space
      new TableRow({
        children: [
          sigCell([new Paragraph({ children: [], spacing: { after: 600 } })]),
          sigCell([new Paragraph({ children: [], spacing: { after: 600 } })]),
        ],
      }),
      // Row 3: names
      new TableRow({
        children: [
          sigCell([
            new Paragraph({
              alignment: AlignmentType.CENTER,
              children: [nameRun(userInfo.reporterNameSig || userInfo.reporterName || '______________________________', true)],
            }),
            new Paragraph({ alignment: AlignmentType.CENTER, children: [nameRun('Signature Over Printed Name')] }),
          ]),
          sigCell([
            new Paragraph({
              alignment: AlignmentType.CENTER,
              children: [nameRun(userInfo.supervisorName || '______________________________', true)],
            }),
            new Paragraph({ alignment: AlignmentType.CENTER, children: [nameRun(userInfo.supervisorPos1 || '')] }),
            new Paragraph({ alignment: AlignmentType.CENTER, children: [nameRun(userInfo.supervisorPos2 || '')] }),
          ]),
        ],
      }),
    ],
  })
}

// ─── IPC handler ─────────────────────────────────────────────────────────────

export function registerDocxHandlers(ipcMain) {
  ipcMain.handle('docx:generate', async (_, { parsedData, userInfo = {}, outputPath }) => {
    if (!outputPath) throw new Error('No output path provided.')
    if (!parsedData?.dateStart || !parsedData?.dateEnd) throw new Error('Parsed data is missing date range.')
    const period = formatPeriod(parsedData.dateStart, parsedData.dateEnd)

    const letterhead = await letterheadRow()

    const doc = new Document({
      styles: {
        default: {
          document: {
            run: { font: 'Calibri', size: SZ_HDR },
          },
        },
      },
      sections: [{
        properties: {
          page: {
            size:   { width: PAGE.width, height: PAGE.height },
            margin: { ...MARGIN, header: 360 },
          },
        },
        headers: {
          default: new Header({ children: [letterhead] }),
        },
        children: [
          titleBoxSection(),
          new Paragraph({ children: [] }),
          ...headerParagraphs(period, userInfo),
          new Paragraph({ children: [] }),
          mainTable(parsedData),
          new Paragraph({ children: [] }),
          signatureTable(userInfo),
        ],
      }],
    })

    const buffer = await Packer.toBuffer(doc)
    await writeFile(outputPath, buffer)
    return outputPath
  })
}
