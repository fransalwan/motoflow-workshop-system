package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	if err := os.MkdirAll("../docs", 0755); err != nil {
		log.Fatalf("Failed to create docs dir: %v", err)
	}

	generateServiceInvoicePDF("../docs/Faktur-Resmi-AHASS-Demo.pdf")
	generateSPKInvoicePDF("../docs/SPK-Resmi-Showroom-Demo.pdf")

	fmt.Println("Successfully generated official demo PDFs in docs/")
}

func generateServiceInvoicePDF(outputPath string) {
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
	pdf.SetFont("Arial", "B", 15)
	pdf.SetXY(15, 26)
	pdf.Cell(120, 8, "FAKTUR RESMI BIAYA SERVIS & SUKU CADANG")

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(15, 34)
	pdf.MultiCell(100, 4.2, "AHASS 01234 NSS Cempaka Putih\nJl. Cempaka Putih Raya No. 128, Jakarta Pusat\nTelp: (021) 420-5678 | Email: ahass.nss@nusantara-sakti.com", "", "L", false)

	// Work Order Badge Box (Top Right)
	pdf.SetFillColor(241, 245, 249)
	pdf.SetDrawColor(203, 213, 225)
	pdf.RoundedRect(125, 26, 70, 22, 2, "1234", "FD")
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(127, 28)
	pdf.Cell(66, 5, "NO: WO-AHASS-20261027-0001")
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetXY(127, 34)
	pdf.Cell(66, 4, "Tanggal: 27 Oktober 2026, 14:30 WIB")
	pdf.SetXY(127, 39)
	pdf.SetTextColor(22, 101, 52) // Green for Lunas
	pdf.SetFont("Arial", "B", 8.5)
	pdf.Cell(66, 4, "STATUS: SELESAI & LUNAS (PAID)")

	pdf.Ln(18)

	// Vehicle & Customer Info Box
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(15, 52, 180, 28, 2, "1234", "FD")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetXY(18, 54)
	pdf.Cell(85, 5, "DATA KONSUMEN & KENDARAAN HONDA")
	pdf.SetXY(105, 54)
	pdf.Cell(85, 5, "PETUGAS BENGKEL RESMI AHASS")

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(51, 65, 85)

	pdf.SetXY(18, 60)
	pdf.Cell(85, 4.5, "Pemilik      : Budi Santoso (0812-3456-7890)")
	pdf.SetXY(18, 65)
	pdf.Cell(85, 4.5, "No. Polisi  : B 3456 KRF  |  Model: Honda Vario 160 ABS (2024)")
	pdf.SetXY(18, 70)
	pdf.Cell(85, 4.5, "Odometer : 12.500 KM (Servis Rutin Berkala)")

	pdf.SetXY(105, 60)
	pdf.Cell(85, 4.5, "Service Advisor : Ahmad Fauzi (SA-012)")
	pdf.SetXY(105, 65)
	pdf.Cell(85, 4.5, "Mekanik AHASS   : Agus Setiawan (Pit 1 - Teknisi Injeksi)")
	pdf.SetXY(105, 70)
	pdf.Cell(85, 4.5, "Metode Bayar    : TUNAI (CASH)")

	// Items Table Header
	pdf.SetXY(15, 84)
	pdf.SetFillColor(204, 0, 0) // AHASS Racing Red
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8.5)
	pdf.CellFormat(10, 7, "NO", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 7, "NOMOR PART/JASA", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 7, "DESKRIPSI PEKERJAAN / SUKU CADANG HGP", "1", 0, "L", true, 0, "")
	pdf.CellFormat(15, 7, "QTY", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, "HARGA (RP)", "1", 0, "R", true, 0, "")
	pdf.CellFormat(15, 7, "DISC", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, "SUBTOTAL (RP)", "1", 1, "R", true, 0, "")

	// Sample Items Data
	type Item struct {
		No       int
		PartNo   string
		Name     string
		Qty      int
		Price    string
		Disc     string
		Subtotal string
	}
	items := []Item{
		{1, "08232-2MB-K0LN1", "AHM Oil MPX2 0.8L (Matic SJ 10W-30)", 1, "56.000", "-", "56.000"},
		{2, "31919-K25-601", "Busi NGK CPR9EA-9 (Original Honda)", 1, "25.000", "-", "25.000"},
		{3, "06455-KRE-K01", "Kampas Rem Cakram Depan HGP", 1, "45.000", "-", "45.000"},
		{4, "SRV-PKG-001", "[JASA] Paket Servis Lengkap Matik AHASS", 1, "90.000", "-", "90.000"},
	}

	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Arial", "", 8)
	var fill bool
	for _, it := range items {
		if fill {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.SetX(15)
		pdf.CellFormat(10, 6.5, fmt.Sprintf("%d", it.No), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(30, 6.5, it.PartNo, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(60, 6.5, it.Name, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(15, 6.5, fmt.Sprintf("%d", it.Qty), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(25, 6.5, it.Price, "1", 0, "R", fill, 0, "")
		pdf.CellFormat(15, 6.5, it.Disc, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(25, 6.5, it.Subtotal, "1", 1, "R", fill, 0, "")
		fill = !fill
	}

	// Financial Totals Summary Block
	yPos := pdf.GetY() + 4
	pdf.SetY(yPos)

	// Left: Warranty & Complaints
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(15, yPos)
	pdf.Cell(100, 4, "Keluhan: Tarikan awal tersendat & rem berdecit.")

	pdf.SetXY(15, yPos+5)
	pdf.Cell(100, 4, "* Syarat garansi servis: Bawa faktur ini, masa berlaku 500 KM / 7 hari.")
	pdf.SetXY(15, yPos+10)
	pdf.Cell(100, 4, "* Seluruh suku cadang dijamin 100% Asli Honda Genuine Parts (HGP).")
	pdf.SetXY(15, yPos+15)
	pdf.Cell(100, 4, "* Terima kasih atas kepercayaan Anda di PT NUSANTARA SAKTI GROUP.")

	// Right: Summary Table
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(15, 23, 42)

	drawRow := func(label, value string, isBold bool) {
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

	drawRow("Subtotal Suku Cadang:", "Rp 126.000", false)
	drawRow("Subtotal Jasa Servis:", "Rp 90.000", false)
	drawRow("PPN (11%):", "Rp 23.760", false)

	// Grand Total Highlight
	pdf.SetFillColor(241, 245, 249)
	pdf.Rect(120, yPos, 75, 7, "F")
	drawRow("GRAND TOTAL:", "Rp 239.760", true)
	drawRow("Jumlah Diterima (TUNAI):", "Rp 250.000", false)
	drawRow("Kembalian:", "Rp 10.240", false)

	// Signatures and Digital Seal
	sigY := yPos + 10
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(30, 41, 59)

	// Customer Signature Box
	pdf.SetXY(25, sigY)
	pdf.Cell(50, 4, "Tanda Tangan Pelanggan,")
	pdf.SetXY(25, sigY+22)
	pdf.Cell(50, 4, "( Budi Santoso )")

	// Digital Seal
	pdf.SetDrawColor(204, 0, 0)
	pdf.SetTextColor(204, 0, 0)
	pdf.SetFont("Arial", "B", 7)
	pdf.Circle(105, sigY+11, 11, "D")
	pdf.SetXY(95, sigY+9)
	pdf.Cell(20, 4, "AHASS 01234 NSS")
	pdf.SetXY(95, sigY+13)
	pdf.Cell(20, 4, "VERIFIED")

	// Service Advisor / Cashier Signature Box
	pdf.SetTextColor(30, 41, 59)
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetXY(135, sigY)
	pdf.Cell(50, 4, "Service Advisor / Kasir AHASS,")
	pdf.SetXY(135, sigY+22)
	pdf.Cell(50, 4, "( Ahmad Fauzi )")

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		log.Fatalf("Failed to save invoice PDF: %v", err)
	}
}

func generateSPKInvoicePDF(outputPath string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Red Header Banner
	pdf.SetFillColor(204, 0, 0) // AHASS Red
	pdf.Rect(15, 15, 180, 8, "F")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(15, 16.5)
	pdf.CellFormat(180, 5, "DEALER RESMI SEPEDA MOTOR HONDA - PT NUSANTARA SAKTI GROUP", "", 0, "C", false, 0, "")

	// Main Company Title
	pdf.SetTextColor(30, 41, 59)
	pdf.SetFont("Arial", "B", 15)
	pdf.SetXY(15, 26)
	pdf.Cell(120, 8, "SURAT PEMESANAN KENDARAAN (SPK)")

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(15, 34)
	pdf.MultiCell(100, 4.2, "PT Nusantara Sakti Group - Showroom Cempaka Putih\nJl. Cempaka Putih Raya No. 128, Jakarta Pusat\nTelp: (021) 420-5678 | Layanan Multifinance: NSC Finance", "", "L", false)

	// SPK Badge Box (Top Right)
	pdf.SetFillColor(241, 245, 249)
	pdf.SetDrawColor(203, 213, 225)
	pdf.RoundedRect(125, 26, 70, 22, 2, "1234", "FD")
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(127, 28)
	pdf.Cell(66, 5, "NO: SPK-NSS-20261027-0001")
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetXY(127, 34)
	pdf.Cell(66, 4, "Tanggal: 27 Oktober 2026, 15:45 WIB")
	pdf.SetXY(127, 39)
	pdf.SetTextColor(204, 0, 0)
	pdf.SetFont("Arial", "B", 8.5)
	pdf.Cell(66, 4, "STATUS: BOOKED SPK (UNIT ALLOCATED)")

	pdf.Ln(18)

	// Customer Info Box
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(15, 52, 180, 24, 2, "1234", "FD")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetXY(18, 54)
	pdf.Cell(85, 5, "DATA CALON PEMBELI / DEBITUR")
	pdf.SetXY(105, 54)
	pdf.Cell(85, 5, "PETUGAS SALES COUNTER NSS")

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(51, 65, 85)

	pdf.SetXY(18, 60)
	pdf.Cell(85, 4.5, "Nama Lengkap : Budi Santoso")
	pdf.SetXY(18, 65)
	pdf.Cell(85, 4.5, "No. Telepon / WA: 0812-3456-7890")
	pdf.SetXY(18, 70)
	pdf.Cell(85, 4.5, "Alamat STNK  : Jl. Percetakan Negara No. 45, Jakarta Pusat")

	pdf.SetXY(105, 60)
	pdf.Cell(85, 4.5, "Sales Counter : Siti Rahma (NSS-SC-09)")
	pdf.SetXY(105, 65)
	pdf.Cell(85, 4.5, "Dealer Cabang : NSS Cempaka Putih - Jakarta")
	pdf.SetXY(105, 70)
	pdf.Cell(85, 4.5, "Alokasi Gudang: Ready Showroom Unit #01")

	// Vehicle Specifications Box
	yUnit := 80.0
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(15, yUnit, 180, 32, 2, "1234", "FD")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetXY(18, yUnit+2)
	pdf.Cell(174, 5, "SPESIFIKASI KENDARAAN & NOMOR RANGKA ASTRA (VIN)")

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(51, 65, 85)

	pdf.SetXY(18, yUnit+8)
	pdf.Cell(85, 4.5, "Model Motor   : Honda Vario 160 ABS (Matte Black 2024)")
	pdf.SetXY(18, yUnit+13)
	pdf.Cell(85, 4.5, "Kapasitas CC   : 160cc 4-Katup eSP+ Liquid Cooled")
	pdf.SetXY(18, yUnit+18)
	pdf.Cell(85, 4.5, "Warna Pilihan  : Matte Black (Hitam Doff Premium)")
	pdf.SetXY(18, yUnit+23)
	pdf.Cell(85, 4.5, "Harga OTR Jkt  : Rp 29.500.000 (On-The-Road Resmi)")

	pdf.SetXY(105, yUnit+8)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Arial", "B", 8.5)
	pdf.Cell(85, 4.5, "Nomor Rangka  : MH1KF8115PK019842")
	pdf.SetXY(105, yUnit+13)
	pdf.Cell(85, 4.5, "Nomor Mesin   : KF81E1019842")
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(51, 65, 85)
	pdf.SetXY(105, yUnit+18)
	pdf.Cell(85, 4.5, "Status Unit   : Terverifikasi Pabrikan Astra Honda Motor")
	pdf.SetXY(105, yUnit+23)
	pdf.Cell(85, 4.5, "Garansi Rangka : Garansi Resmi eSAF 5 Tahun Astra")

	// Leasing Scheme Box
	yLease := yUnit + 36.0
	pdf.SetFillColor(254, 242, 242) // Light red background
	pdf.SetDrawColor(254, 202, 202)
	pdf.RoundedRect(15, yLease, 180, 48, 2, "1234", "FD")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(153, 27, 27)
	pdf.SetXY(18, yLease+2)
	pdf.Cell(174, 5, "SKEMA PEMBIAYAAN KREDIT NSC FINANCE (PT NUSA SURYA CIPTADANA)")

	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(51, 65, 85)

	pdf.SetXY(18, yLease+8)
	pdf.Cell(85, 4.5, "Skema Pembelian : Kredit Multifinance NSC")
	pdf.SetXY(18, yLease+13)
	pdf.Cell(85, 4.5, "Uang Muka (DP)   : Rp 5.900.000 (20.00% Sesuai Regulasi OJK)")
	pdf.SetXY(18, yLease+18)
	pdf.Cell(85, 4.5, "Pokok Pinjaman   : Rp 23.600.000 (Principal Loan)")
	pdf.SetXY(18, yLease+23)
	pdf.Cell(85, 4.5, "Jangka Waktu     : 35 Bulan (Tenor Pilihan Konsumen)")
	pdf.SetXY(18, yLease+28)
	pdf.Cell(85, 4.5, "Bunga Pembiayaan : 18.00% Flat p.a. (1.50% per bulan)")
	pdf.SetXY(18, yLease+33)
	pdf.Cell(85, 4.5, "Biaya Fidusia    : Rp 450.000 (Administrasi & Kontrak)")
	pdf.SetXY(18, yLease+38)
	pdf.Cell(85, 4.5, "Asuransi TLO     : Rp 350.000 (Kehilangan Total 35 Bln)")

	pdf.SetXY(105, yLease+8)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(185, 28, 28)
	pdf.Cell(85, 5, "CICILAN : Rp 1.028.000 / BULAN")
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(105, yLease+13)
	pdf.Cell(85, 4.5, "Jatuh Tempo Tiap Tanggal 15")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetXY(105, yLease+23)
	pdf.Cell(85, 4.5, "TOTAL BAYAR PERTAMA (TDP):")
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(204, 0, 0)
	pdf.SetXY(105, yLease+28)
	pdf.Cell(85, 6, "Rp 6.700.000")

	pdf.SetFont("Arial", "I", 7.5)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(105, yLease+35)
	pdf.MultiCell(75, 3.5, "* TDP mencakup: DP Murni (Rp 5.9jt) + Biaya Administrasi & Fidusia (Rp 450rb) + Asuransi Kehilangan TLO (Rp 350rb).", "", "L", false)

	// Terms and Conditions
	yTerms := yLease + 52.0
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(15, 23, 42)
	pdf.SetXY(15, yTerms)
	pdf.Cell(180, 4, "KETENTUAN PENGURUSAN STNK & BPKB ASTRA HONDA:")

	pdf.SetFont("Arial", "", 7.5)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(15, yTerms+4)
	pdf.Cell(180, 3.5, "1. Estimasi penerbitan STNK & Plat Nomor Resmi Samsat adalah 14 hari kerja sejak pelunasan TDP.")
	pdf.SetXY(15, yTerms+8)
	pdf.Cell(180, 3.5, "2. Buku Pemilik Kendaraan Bermotor (BPKB) diterbitkan dalam 60-90 hari kerja dan disimpan oleh PT Nusa Surya Ciptadana.")
	pdf.SetXY(15, yTerms+12)
	pdf.Cell(180, 3.5, "3. Surat Pemesanan Kendaraan (SPK) ini sah sebagai bukti pemesanan unit resmi dealer PT Nusantara Sakti Group.")

	// Signatures and Digital Seal
	sigY := yTerms + 20.0
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(30, 41, 59)

	// Customer Signature Box
	pdf.SetXY(25, sigY)
	pdf.Cell(50, 4, "Pemesan / Calon Debitur,")
	pdf.SetXY(25, sigY+22)
	pdf.Cell(50, 4, "( Budi Santoso )")

	// Digital Seal
	pdf.SetDrawColor(204, 0, 0)
	pdf.SetTextColor(204, 0, 0)
	pdf.SetFont("Arial", "B", 7)
	pdf.Circle(105, sigY+11, 11, "D")
	pdf.SetXY(95, sigY+9)
	pdf.Cell(20, 4, "NSS SHOWROOM")
	pdf.SetXY(95, sigY+13)
	pdf.Cell(20, 4, "VERIFIED")

	// Sales Signature Box
	pdf.SetTextColor(30, 41, 59)
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetXY(135, sigY)
	pdf.Cell(50, 4, "Sales Counter NSS Dealer,")
	pdf.SetXY(135, sigY+22)
	pdf.Cell(50, 4, "( Siti Rahma )")

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		log.Fatalf("Failed to save SPK PDF: %v", err)
	}
}
