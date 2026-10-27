/**
 * MotoFlow AHASS - Professional Thermal & A4 Printing Engine
 * Designed for PT Nusantara Sakti Group (NSS) Workshop & POS Cashier Terminals
 */

export interface ThermalPrintOptions {
  widthMm?: number
  title?: string
  barcodeText?: string
}

/**
 * Prints thermal receipt text (58mm or 80mm ESC/POS layout)
 * via an invisible iframe without browser header/footer clutter.
 */
export function printThermalReceipt(
  receiptText: string,
  options: ThermalPrintOptions = {}
): void {
  const widthMm = options.widthMm || 58
  const title = options.title || 'Struk Kasir AHASS NSS'
  const barcode = options.barcodeText || ''

  // Create isolated printing iframe
  const iframe = document.createElement('iframe')
  iframe.style.position = 'fixed'
  iframe.style.right = '0'
  iframe.style.bottom = '0'
  iframe.style.width = '0'
  iframe.style.height = '0'
  iframe.style.border = '0'
  iframe.setAttribute('aria-hidden', 'true')
  document.body.appendChild(iframe)

  const doc = iframe.contentWindow?.document
  if (!doc) {
    document.body.removeChild(iframe)
    return
  }

  const barcodeHtml = barcode
    ? `
      <div class="barcode-area">
        <div class="barcode-bars"></div>
        <div class="barcode-text">${barcode}</div>
      </div>
    `
    : ''

  const html = `
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>${title}</title>
        <style>
          @page {
            size: ${widthMm === 58 ? '58mm auto' : '80mm auto'};
            margin: 0;
          }
          * {
            box-sizing: border-box;
          }
          html, body {
            margin: 0;
            padding: 4mm 2mm 8mm 2mm;
            background: #ffffff;
            color: #000000;
            font-family: 'Courier New', Courier, monospace;
            font-size: ${widthMm === 58 ? '9.5pt' : '10pt'};
            line-height: 1.25;
            -webkit-print-color-adjust: exact;
            print-color-adjust: exact;
          }
          .receipt-container {
            display: table;
            margin: 0 auto;
            text-align: left;
          }
          pre {
            margin: 0;
            padding: 0;
            white-space: pre;
            word-break: normal;
            font-family: inherit;
          }
          .barcode-area {
            margin-top: 10px;
            padding-top: 6px;
            border-top: 1px dashed #444;
            text-align: center;
          }
          .barcode-bars {
            height: 28px;
            display: inline-block;
            width: ${widthMm === 58 ? '160px' : '220px'};
            background-image: linear-gradient(90deg, #000 50%, transparent 50%);
            background-size: 4px 100%;
            background-repeat: repeat-x;
          }
          .barcode-text {
            font-size: 9pt;
            letter-spacing: 2px;
            margin-top: 3px;
            font-weight: bold;
            color: #111;
          }
        </style>
      </head>
      <body>
        <div class="receipt-container">
          <pre>${escapeHtml(receiptText)}</pre>
          ${barcodeHtml}
        </div>
      </body>
    </html>
  `

  doc.open()
  doc.write(html)
  doc.close()

  iframe.contentWindow?.focus()
  setTimeout(() => {
    try {
      iframe.contentWindow?.print()
    } catch (e) {
      console.error('Thermal print error:', e)
    } finally {
      // Clean up iframe after print dialog closes
      setTimeout(() => {
        if (document.body.contains(iframe)) {
          document.body.removeChild(iframe)
        }
      }, 1000)
    }
  }, 250)
}

/**
 * Prints formal A4 document from an HTML container element.
 */
export function printA4Document(element: HTMLElement, title: string = 'Faktur Servis AHASS'): void {
  const iframe = document.createElement('iframe')
  iframe.style.position = 'fixed'
  iframe.style.right = '0'
  iframe.style.bottom = '0'
  iframe.style.width = '0'
  iframe.style.height = '0'
  iframe.style.border = '0'
  iframe.setAttribute('aria-hidden', 'true')
  document.body.appendChild(iframe)

  const doc = iframe.contentWindow?.document
  if (!doc) {
    document.body.removeChild(iframe)
    return
  }

  // Extract all stylesheets and head styles for accurate styling
  const headStyles = Array.from(document.querySelectorAll('style, link[rel="stylesheet"]'))
    .map((node) => node.outerHTML)
    .join('\n')

  const contentClone = element.cloneNode(true) as HTMLElement
  // Ensure background is pure white and width is standard
  contentClone.style.maxWidth = '100%'
  contentClone.style.boxShadow = 'none'
  contentClone.style.border = 'none'

  const html = `
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>${title}</title>
        ${headStyles}
        <style>
          @page {
            size: A4 portrait;
            margin: 8mm 10mm 10mm 10mm;
          }
          @media print {
            body {
              background: #ffffff !important;
              padding: 0 !important;
              margin: 0 !important;
              -webkit-print-color-adjust: exact !important;
              print-color-adjust: exact !important;
            }
          }
          body {
            background: #ffffff;
            margin: 0;
            padding: 0;
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
          }
        </style>
      </head>
      <body>
        <div style="width: 100%; max-width: 210mm; margin: 0 auto;">
          ${contentClone.outerHTML}
        </div>
      </body>
    </html>
  `

  doc.open()
  doc.write(html)
  doc.close()

  iframe.contentWindow?.focus()
  setTimeout(() => {
    try {
      iframe.contentWindow?.print()
    } catch (e) {
      console.error('A4 print error:', e)
    } finally {
      setTimeout(() => {
        if (document.body.contains(iframe)) {
          document.body.removeChild(iframe)
        }
      }, 1000)
    }
  }, 350)
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}
