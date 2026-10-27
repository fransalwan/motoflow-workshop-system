package seeders

import (
	"log"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/models"
	"motoflow-backend/services"
)

func SeedAll(db *gorm.DB) error {
	// Check if already seeded
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("[Seeder] Database already populated with seed data. Skipping.")
		return nil
	}

	log.Println("[Seeder] Populating database with realistic Honda AHASS data...")

	// 1. Seed Users
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@nusantara-sakti.com",
			FullName: "Bambang Haryanto (Kepala Bengkel AHASS)",
			Role:     models.RoleAdmin,
			Phone:    "081211112222",
			IsActive: true,
		},
		{
			Username: "advisor",
			Email:    "advisor@nusantara-sakti.com",
			FullName: "Budi Santoso, S.T. (Service Advisor)",
			Role:     models.RoleServiceAdvisor,
			Phone:    "081233334444",
			IsActive: true,
		},
		{
			Username: "kasir",
			Email:    "kasir@nusantara-sakti.com",
			FullName: "Siti Rahmawati (Kasir POS AHASS)",
			Role:     models.RoleCashier,
			Phone:    "081255556666",
			IsActive: true,
		},
		{
			Username: "mekanik",
			Email:    "mekanik@nusantara-sakti.com",
			FullName: "Agus Setiawan (Mekanik Senior)",
			Role:     models.RoleMechanic,
			Phone:    "081277778888",
			IsActive: true,
		},
	}

	for i := range users {
		_ = users[i].SetPassword("password123")
		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}
	}

	// 2. Seed Categories
	categories := []models.Category{
		{Name: "Oli & Pelumas AHM Oil", Code: "OIL", Description: "Pelumas resmi Astra Honda Motor MPX dan SPX"},
		{Name: "Sistem Pengereman (Brake)", Code: "BRK", Description: "Kampas rem cakram dan tromol Honda Genuine Parts"},
		{Name: "Transmisi & CVT (Drive)", Code: "CVT", Description: "V-Belt, roller, kampas ganda matic Honda"},
		{Name: "Kelistrikan & Pengapian", Code: "ELC", Description: "Busi NGK/Denso, aki MF, dan sekring"},
		{Name: "Filter & Perawatan", Code: "FLT", Description: "Saringan udara dan pembersih injeksi"},
		{Name: "Jasa Servis Resmi AHASS", Code: "SRV", Description: "Paket jasa servis mekanik tersertifikasi Honda"},
	}
	for i := range categories {
		if err := db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	// 3. Seed Spareparts & Services
	catOil := categories[0].ID
	catBrk := categories[1].ID
	catCvt := categories[2].ID
	catElc := categories[3].ID
	catFlt := categories[4].ID
	catSrv := categories[5].ID

	parts := []models.Sparepart{
		// Suku Cadang (Honda Genuine Parts)
		{
			PartNumber:   "08232-2MB-K0NZ1",
			Name:         "AHM Oil MPX2 0.8L (Matic)",
			CategoryID:   catOil,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(46000),
			SellingPrice: decimal.NewFromInt(58000),
			Stock:        45,
			MinStock:     10,
			Unit:         "Botol",
			Location:     "Rak Pelumas A-01",
			IsActive:     true,
		},
		{
			PartNumber:   "08232-M99-K1NZ1",
			Name:         "AHM Oil MPX1 0.8L (Bebek / Sport)",
			CategoryID:   catOil,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(44000),
			SellingPrice: decimal.NewFromInt(56000),
			Stock:        30,
			MinStock:     10,
			Unit:         "Botol",
			Location:     "Rak Pelumas A-02",
			IsActive:     true,
		},
		{
			PartNumber:   "08234-2MB-K0NZ1",
			Name:         "AHM Oil SPX2 Fully Synthetic 0.8L",
			CategoryID:   catOil,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(60000),
			SellingPrice: decimal.NewFromInt(75000),
			Stock:        25,
			MinStock:     8,
			Unit:         "Botol",
			Location:     "Rak Pelumas A-03",
			IsActive:     true,
		},
		{
			PartNumber:   "06455-KRE-K01",
			Name:         "Pad Set Front (Kampas Rem Depan Cakram)",
			CategoryID:   catBrk,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(34000),
			SellingPrice: decimal.NewFromInt(45000),
			Stock:        18,
			MinStock:     5,
			Unit:         "Set",
			Location:     "Rak Rem B-01",
			IsActive:     true,
		},
		{
			PartNumber:   "43130-KZL-930",
			Name:         "Shoe Set Brake (Kampas Rem Belakang Tromol)",
			CategoryID:   catBrk,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(26000),
			SellingPrice: decimal.NewFromInt(35000),
			Stock:        12,
			MinStock:     5,
			Unit:         "Set",
			Location:     "Rak Rem B-02",
			IsActive:     true,
		},
		{
			PartNumber:   "23100-K0S-V01",
			Name:         "Drive Belt Kit V-Belt + Roller Vario 160",
			CategoryID:   catCvt,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(145000),
			SellingPrice: decimal.NewFromInt(185000),
			Stock:        8,
			MinStock:     5,
			Unit:         "Set",
			Location:     "Rak CVT C-01",
			IsActive:     true,
		},
		{
			PartNumber:   "23100-K44-V01",
			Name:         "Drive Belt Kit Beat FI / Scoopy eSP",
			CategoryID:   catCvt,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(110000),
			SellingPrice: decimal.NewFromInt(140000),
			Stock:        3, // LOW STOCK TRIGGER (threshold 5)
			MinStock:     5,
			Unit:         "Set",
			Location:     "Rak CVT C-02",
			IsActive:     true,
		},
		{
			PartNumber:   "31919-K25-601",
			Name:         "Busi NGK MR9C-9N (Beat / Scoopy / Vario)",
			CategoryID:   catElc,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(18000),
			SellingPrice: decimal.NewFromInt(25000),
			Stock:        35,
			MinStock:     10,
			Unit:         "Pcs",
			Location:     "Rak Listrik D-01",
			IsActive:     true,
		},
		{
			PartNumber:   "17210-K0J-N00",
			Name:         "Filter Udara Element Comp Air Beat LED / Genio",
			CategoryID:   catFlt,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(38000),
			SellingPrice: decimal.NewFromInt(52000),
			Stock:        15,
			MinStock:     5,
			Unit:         "Pcs",
			Location:     "Rak Filter E-01",
			IsActive:     true,
		},
		{
			PartNumber:   "06406-K45-N01",
			Name:         "Drive Chain Kit Rantai Roda CB150R / CBR150R",
			CategoryID:   catCvt,
			Type:         models.ItemTypePart,
			CostPrice:    decimal.NewFromInt(260000),
			SellingPrice: decimal.NewFromInt(320000),
			Stock:        2, // LOW STOCK TRIGGER
			MinStock:     5,
			Unit:         "Set",
			Location:     "Rak Rantai C-05",
			IsActive:     true,
		},

		// Jasa Servis Mekanik AHASS
		{
			PartNumber:   "SRV-RINGAN",
			Name:         "Jasa Servis Ringan & Safety Check AHASS",
			CategoryID:   catSrv,
			Type:         models.ItemTypeService,
			CostPrice:    decimal.Zero,
			SellingPrice: decimal.NewFromInt(65000),
			Stock:        9999,
			MinStock:     0,
			Unit:         "Paket",
			Location:     "Area Pit Servis",
			IsActive:     true,
		},
		{
			PartNumber:   "SRV-LENGKAP-MATIC",
			Name:         "Jasa Servis Lengkap Matic & Injeksi PGM-FI",
			CategoryID:   catSrv,
			Type:         models.ItemTypeService,
			CostPrice:    decimal.Zero,
			SellingPrice: decimal.NewFromInt(120000),
			Stock:        9999,
			MinStock:     0,
			Unit:         "Paket",
			Location:     "Area Pit Servis",
			IsActive:     true,
		},
		{
			PartNumber:   "SRV-LENGKAP-SPORT",
			Name:         "Jasa Servis Lengkap Sport 150cc",
			CategoryID:   catSrv,
			Type:         models.ItemTypeService,
			CostPrice:    decimal.Zero,
			SellingPrice: decimal.NewFromInt(150000),
			Stock:        9999,
			MinStock:     0,
			Unit:         "Paket",
			Location:     "Area Pit Servis",
			IsActive:     true,
		},
		{
			PartNumber:   "SRV-CVT",
			Name:         "Jasa Servis & Pembersihan CVT Matic",
			CategoryID:   catSrv,
			Type:         models.ItemTypeService,
			CostPrice:    decimal.Zero,
			SellingPrice: decimal.NewFromInt(45000),
			Stock:        9999,
			MinStock:     0,
			Unit:         "Paket",
			Location:     "Area Pit Servis",
			IsActive:     true,
		},
		{
			PartNumber:   "SRV-OIL-CHANGE",
			Name:         "Jasa Ganti Oli Mesin & Transmisi Gardan",
			CategoryID:   catSrv,
			Type:         models.ItemTypeService,
			CostPrice:    decimal.Zero,
			SellingPrice: decimal.NewFromInt(15000),
			Stock:        9999,
			MinStock:     0,
			Unit:         "Paket",
			Location:     "Area Pit Servis",
			IsActive:     true,
		},
	}
	for i := range parts {
		if err := db.Create(&parts[i]).Error; err != nil {
			return err
		}
	}

	// 4. Seed Mechanics
	mechanics := []models.Mechanic{
		{Code: "MEK-01", Name: "Agus Setiawan", Phone: "081277778888", Specialization: "Injeksi PGM-FI & Kelistrikan", PitNumber: 1, IsAvailable: true},
		{Code: "MEK-02", Name: "Hendra Wijaya", Phone: "081388889999", Specialization: "Transmisi & CVT Matic Honda", PitNumber: 2, IsAvailable: true},
		{Code: "MEK-03", Name: "Dedi Pratama", Phone: "081599990000", Specialization: "Overhaul Mesin Bebek & Sport", PitNumber: 3, IsAvailable: true},
		{Code: "MEK-04", Name: "Rizky Ananda", Phone: "081700001111", Specialization: "Fast Moving & Servis Berkala", PitNumber: 4, IsAvailable: true},
	}
	for i := range mechanics {
		if err := db.Create(&mechanics[i]).Error; err != nil {
			return err
		}
	}

	// 5. Seed Customers & Motorcycles
	customers := []models.Customer{
		{
			Name:         "Bambang Pamungkas",
			Phone:        "081812345671",
			Email:        "bambang.p@gmail.com",
			Address:      "Jl. Rawamangun Muka No. 12, Jakarta Timur",
			PlateNumber:  "B 3456 KRF",
			MotorModel:   "Honda Vario 160 ABS",
			Year:         2023,
			FrameNumber:  "MH1KF5118PK120934",
			EngineNumber: "KF51E1120934",
			CurrentKM:    14250,
		},
		{
			Name:         "Ibu Rina Wardani",
			Phone:        "081812345672",
			Email:        "rina.wardani@yahoo.com",
			Address:      "Jl. Percetakan Negara No. 45, Jakarta Pusat",
			PlateNumber:  "B 4567 TYU",
			MotorModel:   "Honda BeAT Deluxe CBS-ISS",
			Year:         2024,
			FrameNumber:  "MH1JM8118RK540192",
			EngineNumber: "JM81E1540192",
			CurrentKM:    6800,
		},
		{
			Name:         "Dimas Prasetyo",
			Phone:        "081812345673",
			Email:        "dimas.prasetyo@outlook.com",
			Address:      "Jl. Utan Kayu Utara No. 88, Jakarta Timur",
			PlateNumber:  "B 5678 WXY",
			MotorModel:   "Honda PCX 160 CBS",
			Year:         2023,
			FrameNumber:  "MH1KF8116PK009218",
			EngineNumber: "KF81E1009218",
			CurrentKM:    21500,
		},
		{
			Name:         "Maya Anggraini",
			Phone:        "081812345674",
			Email:        "maya.anggraini@gmail.com",
			Address:      "Jl. Salemba Tengah No. 23, Jakarta Pusat",
			PlateNumber:  "B 6789 ZZA",
			MotorModel:   "Honda Scoopy Prestige Smart Key",
			Year:         2022,
			FrameNumber:  "MH1JM3119NK892014",
			EngineNumber: "JM31E1892014",
			CurrentKM:    18300,
		},
		{
			Name:         "Eko Kurniawan",
			Phone:        "081812345675",
			Email:        "eko.kurniawan@company.id",
			Address:      "Jl. Cempaka Baru Tengah No. 15, Jakarta Pusat",
			PlateNumber:  "B 7890 NMS",
			MotorModel:   "Honda CB150R Streetfire",
			Year:         2021,
			FrameNumber:  "MH1KC7119MK301982",
			EngineNumber: "KC71E1301982",
			CurrentKM:    32400,
		},
	}
	for i := range customers {
		if err := db.Create(&customers[i]).Error; err != nil {
			return err
		}
	}

	// 6. Seed Realistic Work Orders (Paid, Completed, In-Progress, Draft)
	finService := services.NewFinancialService()
	now := time.Now()
	saID := users[1].ID // Budi Santoso (Advisor)

	// WO 1: Paid - Vario 160 (MPX2, Kampas Rem Cakram, Servis Lengkap Matic)
	wo1Items := []models.WorkOrderItem{
		{
			SparepartID:        &parts[0].ID, // MPX2
			ItemType:           models.ItemTypePart,
			PartNumberSnapshot: parts[0].PartNumber,
			ItemNameSnapshot:   parts[0].Name,
			Quantity:           1,
			PriceSnapshot:      parts[0].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[0].SellingPrice,
		},
		{
			SparepartID:        &parts[3].ID, // Kampas Rem Depan
			ItemType:           models.ItemTypePart,
			PartNumberSnapshot: parts[3].PartNumber,
			ItemNameSnapshot:   parts[3].Name,
			Quantity:           1,
			PriceSnapshot:      parts[3].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[3].SellingPrice,
		},
		{
			SparepartID:        &parts[11].ID, // Servis Lengkap Matic
			ItemType:           models.ItemTypeService,
			PartNumberSnapshot: parts[11].PartNumber,
			ItemNameSnapshot:   parts[11].Name,
			Quantity:           1,
			PriceSnapshot:      parts[11].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[11].SellingPrice,
		},
	}
	wo1Totals := finService.CalculateWorkOrderTotals(wo1Items, decimal.Zero, decimal.NewFromInt(11), decimal.NewFromInt(250000))

	paidTime1 := now.Add(-2 * time.Hour)
	wo1 := models.WorkOrder{
		WONumber:         "WO-AHASS-" + now.Format("20060102") + "-0001",
		CustomerID:       customers[0].ID, // Bambang Vario 160
		MechanicID:       &mechanics[0].ID,
		ServiceAdvisorID: &saID,
		EntryKM:          14250,
		Complaints:       "Tarikan awal agak berat, rem depan berdecit saat macet.",
		ServiceDiagnosis: "Servis lengkap injeksi, ganti oli MPX2, ganti pad set rem cakram depan.",
		Status:           models.WOStatusPaid,
		PaymentMethod:    models.PaymentCash,
		PaymentStatus:    models.PaymentStatusPaid,
		PartsSubtotal:    wo1Totals.PartsSubtotal,
		ServicesSubtotal: wo1Totals.ServicesSubtotal,
		TotalDiscount:    wo1Totals.TotalDiscount,
		TaxRate:          wo1Totals.TaxRate,
		TaxAmount:        wo1Totals.TaxAmount,
		GrandTotal:       wo1Totals.GrandTotal,
		PaidAmount:       decimal.NewFromInt(250000),
		ChangeAmount:     wo1Totals.ChangeAmount,
		CompletedAt:      &paidTime1,
		PaidAt:           &paidTime1,
		Items:            wo1Items,
	}
	db.Create(&wo1)

	// WO 2: Completed (Awaiting Cashier Payment) - Scoopy (Oli SPX2 + Servis CVT)
	wo2Items := []models.WorkOrderItem{
		{
			SparepartID:        &parts[2].ID, // SPX2
			ItemType:           models.ItemTypePart,
			PartNumberSnapshot: parts[2].PartNumber,
			ItemNameSnapshot:   parts[2].Name,
			Quantity:           1,
			PriceSnapshot:      parts[2].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[2].SellingPrice,
		},
		{
			SparepartID:        &parts[13].ID, // Servis CVT
			ItemType:           models.ItemTypeService,
			PartNumberSnapshot: parts[13].PartNumber,
			ItemNameSnapshot:   parts[13].Name,
			Quantity:           1,
			PriceSnapshot:      parts[13].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[13].SellingPrice,
		},
	}
	wo2Totals := finService.CalculateWorkOrderTotals(wo2Items, decimal.Zero, decimal.NewFromInt(11), decimal.Zero)
	completedTime2 := now.Add(-30 * time.Minute)
	wo2 := models.WorkOrder{
		WONumber:         "WO-AHASS-" + now.Format("20060102") + "-0002",
		CustomerID:       customers[3].ID, // Maya Scoopy
		MechanicID:       &mechanics[1].ID, // Hendra CVT
		ServiceAdvisorID: &saID,
		EntryKM:          18300,
		Complaints:       "Getar pada CVT saat rpm rendah / angkatan pertama.",
		ServiceDiagnosis: "Pembersihan rumah mangkok ganda CVT dan ganti oli synthetic SPX2.",
		Status:           models.WOStatusCompleted,
		PaymentMethod:    models.PaymentQRIS,
		PaymentStatus:    models.PaymentStatusUnpaid,
		PartsSubtotal:    wo2Totals.PartsSubtotal,
		ServicesSubtotal: wo2Totals.ServicesSubtotal,
		TotalDiscount:    wo2Totals.TotalDiscount,
		TaxRate:          wo2Totals.TaxRate,
		TaxAmount:        wo2Totals.TaxAmount,
		GrandTotal:       wo2Totals.GrandTotal,
		PaidAmount:       decimal.Zero,
		ChangeAmount:     decimal.Zero,
		CompletedAt:      &completedTime2,
		Items:            wo2Items,
	}
	db.Create(&wo2)

	// WO 3: In Progress - Honda PCX 160 (Drive Belt Kit + Busi + Servis Lengkap)
	wo3Items := []models.WorkOrderItem{
		{
			SparepartID:        &parts[5].ID, // V-Belt Vario/PCX 160
			ItemType:           models.ItemTypePart,
			PartNumberSnapshot: parts[5].PartNumber,
			ItemNameSnapshot:   parts[5].Name,
			Quantity:           1,
			PriceSnapshot:      parts[5].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[5].SellingPrice,
		},
		{
			SparepartID:        &parts[7].ID, // Busi NGK
			ItemType:           models.ItemTypePart,
			PartNumberSnapshot: parts[7].PartNumber,
			ItemNameSnapshot:   parts[7].Name,
			Quantity:           1,
			PriceSnapshot:      parts[7].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[7].SellingPrice,
		},
		{
			SparepartID:        &parts[11].ID, // Servis Lengkap
			ItemType:           models.ItemTypeService,
			PartNumberSnapshot: parts[11].PartNumber,
			ItemNameSnapshot:   parts[11].Name,
			Quantity:           1,
			PriceSnapshot:      parts[11].SellingPrice,
			DiscountPercent:    decimal.Zero,
			Subtotal:           parts[11].SellingPrice,
		},
	}
	wo3Totals := finService.CalculateWorkOrderTotals(wo3Items, decimal.Zero, decimal.NewFromInt(11), decimal.Zero)
	wo3 := models.WorkOrder{
		WONumber:         "WO-AHASS-" + now.Format("20060102") + "-0003",
		CustomerID:       customers[2].ID, // Dimas PCX
		MechanicID:       &mechanics[1].ID,
		ServiceAdvisorID: &saID,
		EntryKM:          21500,
		Complaints:       "Jadwal penggantian V-Belt berkala 20.000 KM dan tarikan tersendat sesekali.",
		ServiceDiagnosis: "Penggantian V-Belt set roller, penggantian busi, servis sistem injeksi PGM-FI.",
		Status:           models.WOStatusInProgress,
		PaymentMethod:    models.PaymentCash,
		PaymentStatus:    models.PaymentStatusUnpaid,
		PartsSubtotal:    wo3Totals.PartsSubtotal,
		ServicesSubtotal: wo3Totals.ServicesSubtotal,
		TotalDiscount:    wo3Totals.TotalDiscount,
		TaxRate:          wo3Totals.TaxRate,
		TaxAmount:        wo3Totals.TaxAmount,
		GrandTotal:       wo3Totals.GrandTotal,
		PaidAmount:       decimal.Zero,
		ChangeAmount:     decimal.Zero,
		Items:            wo3Items,
	}
	db.Create(&wo3)

	log.Println("[Seeder] Successfully seeded AHASS users, categories, parts, mechanics, and work orders.")
	return nil
}

