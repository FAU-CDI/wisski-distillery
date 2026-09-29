/**
 * CSV export for HTML tables.
 *
 * Wire a button with data-csv-export pointing at a table selector, e.g.:
 *   <button type="button" data-csv-export=".instances-table">Export CSV</button>
 *
 * Per-cell overrides on th/td:
 *   data-export="value"  – use this string instead of the cell text
 *   data-skip-export     – omit this cell from the export
 */

const csvEscape = (value: string): string => {
  if (/[",\r\n]/.test(value)) {
    return '"' + value.replace(/"/g, '""') + '"'
  }
  return value
}

const cellValue = (cell: HTMLTableCellElement): string | null => {
  if (cell.hasAttribute('data-skip-export')) {
    return null
  }
  if (cell.hasAttribute('data-export')) {
    return cell.getAttribute('data-export') ?? ''
  }
  const text = cell.innerText !== '' ? cell.innerText : (cell.textContent ?? '')
  return text.replace(/\s+/g, ' ').trim()
}

const rowValues = (row: HTMLTableRowElement): string[] => {
  const values: string[] = []
  row.querySelectorAll<HTMLTableCellElement>('th, td').forEach(cell => {
    const value = cellValue(cell)
    if (value !== null) {
      values.push(value)
    }
  })
  return values
}

const tableToCSV = (table: HTMLTableElement): string => {
  const lines: string[] = []

  const headRow = table.tHead?.rows[0]
  if (headRow !== undefined) {
    lines.push(rowValues(headRow).map(csvEscape).join(','))
  }

  const bodies = table.tBodies.length > 0 ? Array.from(table.tBodies) : [table]
  bodies.forEach(body => {
    Array.from(body.rows).forEach(row => {
      // skip header-like rows accidentally in tbody
      if (row.parentElement?.tagName === 'THEAD') {
        return
      }
      lines.push(rowValues(row).map(csvEscape).join(','))
    })
  })

  return lines.join('\r\n') + '\r\n'
}

const downloadCSV = (filename: string, content: string): void => {
  const blob = new Blob([content], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}

document.querySelectorAll<HTMLElement>('[data-csv-export]').forEach(button => {
  const selector = button.dataset.csvExport
  if (typeof selector !== 'string' || selector.length === 0) {
    console.warn('Element', button, 'has invalid or missing data-csv-export attribute')
    return
  }

  button.addEventListener('click', (ev) => {
    ev.preventDefault()
    const table = document.querySelector<HTMLTableElement>(selector)
    if (table === null) {
      console.warn('No table found for selector', selector)
      return
    }

    const filename = button.dataset.csvFilename ?? 'export.csv'
    downloadCSV(filename, tableToCSV(table))
  })
})
