package notif

import (
	"fmt"

	"gorm.io/gorm"
)

const queryInsertNotif = `
	INSERT INTO log_notif (key_notif, desc_notif, user_pengirim, user_penerima, jenis_notif, entry_date)
	VALUES (?, ?, ?, ?, ?, NOW())`

// Insert menyimpan satu record notifikasi ke tabel log_notif.
// Dipanggil di dalam transaksi (tx) agar atomic bersama operasi utama.
//
//   - keyNotif     : id_pengajuan yang menjadi kunci notifikasi
//   - descNotif    : deskripsi notifikasi, mis. "Approval Penyusunan, ID : ..."
//   - userPengirim : PERNR pengirim (dari header userq)
//   - userPenerima : PERNR penerima (approval_posisi berikutnya atau entry_user)
//   - jenisNotif   : jenis notifikasi, mis. "approval_penyusunan", "penyusunan_ditolak"
func Insert(tx *gorm.DB, keyNotif, descNotif, userPengirim, userPenerima, jenisNotif string) error {
	if err := tx.Exec(queryInsertNotif,
		keyNotif, descNotif, userPengirim, userPenerima, jenisNotif,
	).Error; err != nil {
		return fmt.Errorf("gagal insert notifikasi: %w", err)
	}
	return nil
}
