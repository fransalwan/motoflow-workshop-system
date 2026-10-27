package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
)

type PrintController struct {
	db *gorm.DB
}

func NewPrintController(db *gorm.DB) *PrintController {
	return &PrintController{db: db}
}

type ThermalReceiptResponse struct {
	WidthMM        int              `json:"width_mm"`
	RawText        string           `json:"raw_text"`
	WorkOrder      models.WorkOrder `json:"work_order"`
	StoreName      string           `json:"store_name"`
	DealerName     string           `json:"dealer_name"`
	StoreAddress   string           `json:"store_address"`
	StorePhone     string           `json:"store_phone"`
	WarrantyNotice string           `json:"warranty_notice"`
}

// ThermalReceipt generates formatted thermal text (58mm or 80mm) ready for ESC/POS printing
func (ctrl *PrintController) ThermalReceipt(c *gin.Context) {
	id := c.Param("id")
	widthStr := c.DefaultQuery("width", "58")
	width := 58
	if widthStr == "80" {
		width = 80
	}

	var wo models.WorkOrder
	if err := ctrl.db.
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceAdvisor").
		Preload("Items").
		First(&wo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Work Order tidak ditemukan", err.Error()))
		return
	}

	lineLength := 32 // 58mm standard character width
	if width == 80 {
		lineLength = 48
	}

	var sb strings.Builder
	separator := strings.Repeat("-", lineLength) + "\n"
	doubleSeparator := strings.Repeat("=", lineLength) + "\n"

	center := func(s string) string {
		if len(s) >= lineLength {
			return s[:lineLength] + "\n"
		}
		spaces := (lineLength - len(s)) / 2
		return strings.Repeat(" ", spaces) + s + "\n"
	}

	row := func(left, right string) string {
		maxLeft := lineLength - len(right) - 1
		if maxLeft < 0 {
			maxLeft = 0
		}
		if len(left) > maxLeft {
			left = left[:maxLeft]
		}
		spacing := lineLength - len(left) - len(right)
		if spacing < 1 {
			spacing = 1
		}
		return left + strings.Repeat(" ", spacing) + right + "\n"
	}

	// Header AHASS NSS
	sb.WriteString(center("AHASS 01234 - NUSANTARA SAKTI"))
	sb.WriteString(center("PT NUSANTARA SAKTI GROUP"))
	sb.WriteString(center("Jl. Cempaka Putih Raya No. 128, Jakarta"))
	sb.WriteString(center("Telp: (021) 420-5678 / 0812-3456-7890"))
	sb.WriteString(doubleSeparator)

	// Transaction info
	sb.WriteString(row("No. Faktur:", wo.WONumber))
	sb.WriteString(row("Tanggal:", wo.CreatedAt.Format("02/01/2006 15:04")))
	if wo.Customer != nil {
		sb.WriteString(row("No. Polisi:", wo.Customer.PlateNumber))
		sb.WriteString(row("Motor:", wo.Customer.MotorModel))
		sb.WriteString(row("Pemilik:", wo.Customer.Name))
		if wo.EntryKM > 0 {
			sb.WriteString(row("Kilometer:", fmt.Sprintf("%d KM", wo.EntryKM)))
		}
	}
	if wo.Mechanic != nil {
		sb.WriteString(row("Mekanik:", fmt.Sprintf("%s (Pit %d)", wo.Mechanic.Name, wo.Mechanic.PitNumber)))
	}
	if wo.ServiceAdvisor != nil {
		sb.WriteString(row("Advisor:", wo.ServiceAdvisor.FullName))
	}
	sb.WriteString(separator)

	// Items
	sb.WriteString(row("RINCIAN PART & JASA", "TOTAL"))
	sb.WriteString(separator)

	for _, item := range wo.Items {
		itemTitle := item.ItemNameSnapshot
		if item.ItemType == models.ItemTypeService {
			itemTitle = "[JASA] " + itemTitle
		}
		sb.WriteString(itemTitle + "\n")
		qtyPrice := fmt.Sprintf(" %dx Rp %s", item.Quantity, item.PriceSnapshot.StringFixed(0))
		sub := "Rp " + item.Subtotal.StringFixed(0)
		sb.WriteString(row(qtyPrice, sub))
	}

	sb.WriteString(separator)
	sb.WriteString(row("Subtotal Suku Cadang:", "Rp "+wo.PartsSubtotal.StringFixed(0)))
	sb.WriteString(row("Subtotal Jasa Servis:", "Rp "+wo.ServicesSubtotal.StringFixed(0)))

	if wo.TotalDiscount.GreaterThan(decimal.Zero) {
		sb.WriteString(row("Diskon:", "-Rp "+wo.TotalDiscount.StringFixed(0)))
	}

	if wo.TaxAmount.GreaterThan(decimal.Zero) {
		taxLabel := fmt.Sprintf("PPN (%s%%):", wo.TaxRate.StringFixed(0))
		sb.WriteString(row(taxLabel, "Rp "+wo.TaxAmount.StringFixed(0)))
	}

	sb.WriteString(doubleSeparator)
	sb.WriteString(row("GRAND TOTAL:", "Rp "+wo.GrandTotal.StringFixed(0)))

	if wo.PaymentStatus == models.PaymentStatusPaid {
		payMethod := strings.ToUpper(string(wo.PaymentMethod))
		sb.WriteString(row("Metode Pembayaran:", payMethod))
		sb.WriteString(row("Tunai/Bayar:", "Rp "+wo.PaidAmount.StringFixed(0)))
		sb.WriteString(row("Kembalian:", "Rp "+wo.ChangeAmount.StringFixed(0)))
	} else {
		sb.WriteString(row("Status:", "BELUM DIBAYAR (DRAFT)"))
	}

	sb.WriteString(doubleSeparator)
	sb.WriteString(center("GARANSI SERVIS AHASS 500 KM"))
	sb.WriteString(center("Atau 1 Minggu Sejak Servis"))
	sb.WriteString(center("Gunakan Selalu Honda Genuine Parts"))
	sb.WriteString(center("PT NUSANTARA SAKTI - SATU HATI"))
	sb.WriteString("\n\n")

	resp := ThermalReceiptResponse{
		WidthMM:        width,
		RawText:        sb.String(),
		WorkOrder:      wo,
		StoreName:      "AHASS 01234 - NSS",
		DealerName:     "PT NUSANTARA SAKTI GROUP",
		StoreAddress:   "Jl. Cempaka Putih Raya No. 128, Jakarta Pusat",
		StorePhone:     "(021) 420-5678",
		WarrantyNotice: "Garansi servis 500 KM / 7 hari kerja. Hubungi kami jika kendala berlanjut.",
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Struk kasir thermal berhasil di-generate", resp))
}

// ServiceInvoicePDF generates a formal A4 service invoice with dual signatures and digital seal
func (ctrl *PrintController) ServiceInvoicePDF(c *gin.Context) {
	id := c.Param("id")
	var wo models.WorkOrder

	if err := ctrl.db.
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceAdvisor").
		Preload("Items").
		First(&wo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Work Order tidak ditemukan", err.Error()))
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Red AHASS Banner Header
	pdf.SetFillColor(204, 0, 0) // AHASS Racing Red
	pdf.Rect(15, 15, 180, 8, "F")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(15, 16.5)
	pdf.CellFormat(180, 5, "HONDA AHASS - PT NUSANTARA SAKTI GROUP (OFFICIAL AUTHORIZED WORKSHOP)", "", 0, "C", false, 0, "")

	// Main Company Title
	pdf.SetTextColor(30, 41, 59)
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(15, 26)
	pdf.Cell(120, 8, "FAKTUR BIAYA SERVIS & SUKU CADANG")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(15, 34)
	pdf.MultiCell(100, 4.5, "AHASS 01234 NSS Cempaka Putih\nJl. Cempaka Putih Raya No. 128, Jakarta Pusat\nTelp: (021) 420-5678 | Email: ahass.nss@nusantara-sakti.com", "", "L", false)

	// Work Order Badge Box (Top Right)
	pdf.SetFillColor(241, 245, 249)
	pdf.SetDrawColor(203, 213, 225)
	pdf.RoundedRect(125, 26, 70, 22, 2, "1234", "FD")
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(127, 28)
	pdf.Cell(66, 5, fmt.Sprintf("NO: %s", wo.WONumber))
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetXY(127, 34)
	pdf.Cell(66, 4, fmt.Sprintf("Tanggal: %s", wo.CreatedAt.Format("02 January 2006, 15:04 WIB")))
	pdf.SetXY(127, 39)
	statusText := strings.ToUpper(string(wo.Status))
	if wo.PaymentStatus == models.PaymentStatusPaid {
		statusText += " (LUNAS)"
	}
	pdf.Cell(66, 4, fmt.Sprintf("Status: %s", statusText))

	pdf.Ln(18)

	// Vehicle & Customer Info Box
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(15, 52, 180, 28, 2, "1234", "FD")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetXY(18, 54)
	pdf.Cell(85, 5, "DATA KONSUMEN & KENDARAAN HONDA")
	pdf.SetXY(105, 54)
	pdf.Cell(85, 5, "PETUGAS BENGKEL AHASS")

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(51, 65, 85)

	custName := "-"
	custPlate := "-"
	custModel := "-"
	custKM := "-"
	if wo.Customer != nil {
		custName = fmt.Sprintf("%s (%s)", wo.Customer.Name, wo.Customer.Phone)
		custPlate = wo.Customer.PlateNumber
		custModel = fmt.Sprintf("%s (%d)", wo.Customer.MotorModel, wo.Customer.Year)
		custKM = fmt.Sprintf("%d KM", wo.EntryKM)
	}
	saName := "-"
	if wo.ServiceAdvisor != nil {
		saName = wo.ServiceAdvisor.FullName
	}
	mekName := "-"
	if wo.Mechanic != nil {
		mekName = fmt.Sprintf("%s (Pit %d - %s)", wo.Mechanic.Name, wo.Mechanic.PitNumber, wo.Mechanic.Specialization)
	}

	pdf.SetXY(18, 60)
	pdf.Cell(85, 4.5, fmt.Sprintf("Pemilik      : %s", custName))
	pdf.SetXY(18, 65)
	pdf.Cell(85, 4.5, fmt.Sprintf("No. Polisi  : %s  |  Model: %s", custPlate, custModel))
	pdf.SetXY(18, 70)
	pdf.Cell(85, 4.5, fmt.Sprintf("Odometer : %s", custKM))

	pdf.SetXY(105, 60)
	pdf.Cell(85, 4.5, fmt.Sprintf("Service Advisor : %s", saName))
	pdf.SetXY(105, 65)
	pdf.Cell(85, 4.5, fmt.Sprintf("Mekanik AHASS   : %s", mekName))
	pdf.SetXY(105, 70)
	pdf.Cell(85, 4.5, fmt.Sprintf("Metode Bayar    : %s", strings.ToUpper(string(wo.PaymentMethod))))

	// Items Table Header
	pdf.SetXY(15, 85)
	pdf.SetFillColor(220, 38, 38) // Red Header
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8.5)
	pdf.CellFormat(10, 7, "NO", "1", 0, "C", true, 0, "")
	pdf.CellFormat(28, 7, "KODE PART/JASA", "1", 0, "C", true, 0, "")
	pdf.CellFormat(62, 7, "DESKRIPSI PEKERJAAN / PART", "1", 0, "L", true, 0, "")
	pdf.CellFormat(15, 7, "QTY", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, "HARGA (RP)", "1", 0, "R", true, 0, "")
	pdf.CellFormat(15, 7, "DISC", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, "SUBTOTAL (RP)", "1", 1, "R", true, 0, "")

	// Items Table Body
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Arial", "", 8)
	var fill bool
	for idx, item := range wo.Items {
		if fill {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		discStr := "-"
		if item.DiscountPercent.GreaterThan(decimal.Zero) {
			discStr = fmt.Sprintf("%s%%", item.DiscountPercent.StringFixed(0))
		}

		itemTypeBadge := item.ItemNameSnapshot
		if item.ItemType == models.ItemTypeService {
			itemTypeBadge = "[JASA] " + itemTypeBadge
		}

		pdf.SetX(15)
		pdf.CellFormat(10, 6, fmt.Sprintf("%d", idx+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(28, 6, item.PartNumberSnapshot, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(62, 6, itemTypeBadge, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(15, 6, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(25, 6, item.PriceSnapshot.StringFixed(0), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(15, 6, discStr, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(25, 6, item.Subtotal.StringFixed(0), "1", 1, "R", fill, 0, "")
		fill = !fill
	}

	// Financial Totals Summary Block
	yPos := pdf.GetY() + 4
	pdf.SetY(yPos)

	// Left: Warranty & Complaints
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(15, yPos)
	complaintText := "Keluhan Konsumen: " + wo.Complaints
	if len(complaintText) > 60 {
		complaintText = complaintText[:60] + "..."
	}
	pdf.Cell(100, 4, complaintText)

	pdf.SetXY(15, yPos+5)
	pdf.Cell(100, 4, "* Syarat garansi servis: Bawa faktur ini, masa berlaku 500 KM / 7 hari.")
	pdf.SetXY(15, yPos+10)
	pdf.Cell(100, 4, "* Suku cadang Honda Genuine Parts (HGP) dijamin 100% Asli AHM.")

	// Right: Summary Table
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(15, 23, 42)

	drawSummaryRow := func(label, value string, isBold bool) {
		if isBold {
			pdf.SetFont("Arial", "B", 9)
		} else {
			pdf.SetFont("Arial", "", 8.5)
		}
		pdf.SetXY(120, yPos)
		pdf.Cell(40, 5, label)
		pdf.CellFormat(35, 5, value, "", 1, "R", false, 0, "")
		yPos += 5
	}

	drawSummaryRow("Subtotal Suku Cadang:", "Rp "+wo.PartsSubtotal.StringFixed(0), false)
	drawSummaryRow("Subtotal Jasa Servis:", "Rp "+wo.ServicesSubtotal.StringFixed(0), false)
	if wo.TotalDiscount.GreaterThan(decimal.Zero) {
		drawSummaryRow("Total Diskon:", "-Rp "+wo.TotalDiscount.StringFixed(0), false)
	}
	if wo.TaxAmount.GreaterThan(decimal.Zero) {
		drawSummaryRow(fmt.Sprintf("PPN (%s%%):", wo.TaxRate.StringFixed(0)), "Rp "+wo.TaxAmount.StringFixed(0), false)
	}

	// Grand Total Highlight
	pdf.SetFillColor(241, 245, 249)
	pdf.Rect(120, yPos, 75, 7, "F")
	drawSummaryRow("GRAND TOTAL:", "Rp "+wo.GrandTotal.StringFixed(0), true)

	if wo.PaymentStatus == models.PaymentStatusPaid {
		drawSummaryRow("Jumlah Bayar ("+strings.ToUpper(string(wo.PaymentMethod))+"):", "Rp "+wo.PaidAmount.StringFixed(0), false)
		drawSummaryRow("Kembalian:", "Rp "+wo.ChangeAmount.StringFixed(0), false)
	}

	// Signatures and Digital Seal
	sigY := yPos + 10
	if sigY > 240 {
		pdf.AddPage()
		sigY = 25
	}

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(30, 41, 59)

	// Customer Signature Box
	pdf.SetXY(25, sigY)
	pdf.Cell(50, 4, "Tanda Tangan Pelanggan,")
	pdf.SetXY(25, sigY+22)
	customerDisplayName := "............................................."
	if wo.Customer != nil && wo.Customer.Name != "" {
		customerDisplayName = wo.Customer.Name
	}
	pdf.Cell(50, 4, fmt.Sprintf("( %s )", customerDisplayName))

	// Digital Seal
	pdf.SetDrawColor(220, 38, 38)
	pdf.SetTextColor(220, 38, 38)
	pdf.SetFont("Arial", "B", 7)
	pdf.Circle(105, sigY+11, 10, "D")
	pdf.SetXY(97, sigY+9)
	pdf.Cell(16, 4, "AHASS NSS")

	// Service Advisor / Cashier Signature Box
	pdf.SetTextColor(30, 41, 59)
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetXY(135, sigY)
	pdf.Cell(50, 4, "Service Advisor / Kasir AHASS,")
	pdf.SetXY(135, sigY+22)
	advisorDisplayName := "Petugas Kasir"
	if wo.ServiceAdvisor != nil {
		advisorDisplayName = wo.ServiceAdvisor.FullName
	}
	pdf.Cell(50, 4, fmt.Sprintf("( %s )", advisorDisplayName))

	// Output PDF buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menghasilkan dokumen PDF", err.Error()))
		return
	}

	fileName := fmt.Sprintf("Faktur-AHASS-%s.pdf", wo.WONumber)
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
