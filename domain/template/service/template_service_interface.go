package service

import (
	dto "permen_api/domain/template/dto"
	repo "permen_api/domain/template/repo"
)

type (
	TemplateServiceInterface interface {
		// GenerateFormatPenyusunanKpi digunakan oleh endpoint POST /template/format-penyusunan-kpi.
		// Menghasilkan file Excel template penyusunan KPI sesuai triwulan (tanpa isi data baris).
		// Sheet 1 — nama sheet mengikuti nilai triwulan dari request (TW1, TW2, TW3, TW4):
		//   Jika triwulan TW1/TW3 → kolom A–O (format base).
		//   Jika triwulan TW2/TW4 → kolom A–U (format extended).
		// Sheet 2 — nama sheet "KPI":
		//   Kolom A (KPI) dan B (Polarisasi) dari join mst_kpi dan mst_polarisasi.
		GenerateFormatPenyusunanKpi(req *dto.FormatPenyusunanKpiRequest) (fileBytes []byte, filename string, err error)

		// GenerateRevisionPenyusunanKpi digunakan oleh endpoint POST /template/tolakan-penyusunan-kpi.
		// Menghasilkan file Excel yang sudah terisi data baris sub KPI berdasarkan id_pengajuan.
		// Format kolom mengikuti triwulan dari DB:
		//   TW1/TW3 → kolom A–O (format base).
		//   TW2/TW4 → kolom A–U (format extended, kolom P–U terisi data result/method/challenge).
		// Sheet 2 — nama sheet "KPI":
		//   Kolom A (KPI) dan B (Polarisasi) dari join mst_kpi dan mst_polarisasi.
		GenerateRevisionPenyusunanKpi(req *dto.RevisionPenyusunanKpiRequest) (fileBytes []byte, filename string, err error)

		// GenerateRevisionRealisasiKpi digunakan oleh endpoint POST /template/revision-realisasi-kpi.
		// Menghasilkan file Excel realisasi KPI yang sudah terisi data realisasi berdasarkan id_pengajuan,
		// sehingga user dapat langsung merevisi dan mengupload ulang via /realisasi-kpi/revision.
		// Kolom A–I dari penyusunan; kolom J–M pre-filled data realisasi sebelumnya.
		// Format kolom extended mengikuti triwulan dari request:
		//   TW1/TW3 → A–M.
		//   TW2/TW4 → A–Y (N,O,R,S,V,W dari penyusunan; P,Q,T,U,X,Y pre-filled realisasi).
		GenerateRevisionRealisasiKpi(req *dto.RevisionRealisasiKpiRequest) (fileBytes []byte, filename string, err error)

		// GenerateFormatRealisasiKpi digunakan oleh endpoint POST /template/format-realisasi-kpi.
		// Menghasilkan file Excel template realisasi KPI berdasarkan id_pengajuan dan triwulan dari request.
		// Kolom A–I terisi data dari DB; kolom J–M dikosongkan untuk diisi user.
		// Format kolom extended mengikuti triwulan dari request:
		//   TW1/TW3 → kolom A–S (kolom N–S terisi data result/process/context dari DB).
		//   TW2/TW4 → kolom A–Y (kolom N, O, R, S, V, W dari DB; kolom P, Q, T, U, X, Y kosong untuk user).
		// Row 1 adalah header kolom; data dimulai dari row 2.
		GenerateFormatRealisasiKpi(req *dto.FormatRealisasiKpiRequest) (fileBytes []byte, filename string, err error)
	}

	templateService struct {
		repo repo.TemplateRepoInterface
	}
)

func NewTemplateService(repo repo.TemplateRepoInterface) *templateService {
	return &templateService{repo: repo}
}
