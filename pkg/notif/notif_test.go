package notif

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// =============================================================================
// Mock DBExecutor
// =============================================================================

type mockDBExecutor struct {
	// execErr dikembalikan sebagai gorm.DB.Error saat Exec dipanggil.
	execErr error
	// capturedSQL menyimpan query yang diterima.
	capturedSQL string
	// capturedArgs menyimpan argumen yang diterima.
	capturedArgs []interface{}
}

func (m *mockDBExecutor) Exec(sql string, values ...interface{}) *gorm.DB {
	m.capturedSQL = sql
	m.capturedArgs = values
	return &gorm.DB{Error: m.execErr}
}

// =============================================================================
// Helper
// =============================================================================

func newMock(err error) *mockDBExecutor {
	return &mockDBExecutor{execErr: err}
}

// =============================================================================
// Tests
// =============================================================================

func TestInsert_Success(t *testing.T) {
	mock := newMock(nil)

	err := Insert(mock,
		"KPI-2024-001",
		"Approval Penyusunan, ID : KPI-2024-001",
		"USER001",
		"APPROVER001",
		"approval_penyusunan",
	)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestInsert_CapturesCorrectSQL(t *testing.T) {
	mock := newMock(nil)

	Insert(mock,
		"KPI-2024-001",
		"Approval Penyusunan, ID : KPI-2024-001",
		"USER001",
		"APPROVER001",
		"approval_penyusunan",
	)

	if mock.capturedSQL != queryInsertNotif {
		t.Errorf("expected SQL:\n%s\ngot:\n%s", queryInsertNotif, mock.capturedSQL)
	}
}

func TestInsert_CapturesCorrectArgs(t *testing.T) {
	mock := newMock(nil)

	keyNotif := "KPI-2024-001"
	descNotif := "Approval Penyusunan, ID : KPI-2024-001"
	userPengirim := "USER001"
	userPenerima := "APPROVER001"
	jenisNotif := "approval_penyusunan"

	Insert(mock, keyNotif, descNotif, userPengirim, userPenerima, jenisNotif)

	if len(mock.capturedArgs) != 5 {
		t.Fatalf("expected 5 args, got: %d", len(mock.capturedArgs))
	}
	if mock.capturedArgs[0] != keyNotif {
		t.Errorf("args[0]: expected '%s', got '%v'", keyNotif, mock.capturedArgs[0])
	}
	if mock.capturedArgs[1] != descNotif {
		t.Errorf("args[1]: expected '%s', got '%v'", descNotif, mock.capturedArgs[1])
	}
	if mock.capturedArgs[2] != userPengirim {
		t.Errorf("args[2]: expected '%s', got '%v'", userPengirim, mock.capturedArgs[2])
	}
	if mock.capturedArgs[3] != userPenerima {
		t.Errorf("args[3]: expected '%s', got '%v'", userPenerima, mock.capturedArgs[3])
	}
	if mock.capturedArgs[4] != jenisNotif {
		t.Errorf("args[4]: expected '%s', got '%v'", jenisNotif, mock.capturedArgs[4])
	}
}

func TestInsert_DBError_ReturnsWrappedError(t *testing.T) {
	dbErr := errors.New("connection refused")
	mock := newMock(dbErr)

	err := Insert(mock,
		"KPI-2024-001",
		"Approval Penyusunan, ID : KPI-2024-001",
		"USER001",
		"APPROVER001",
		"approval_penyusunan",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("expected error to wrap '%v', got: '%v'", dbErr, err)
	}
}

func TestInsert_ErrorMessageContainsContext(t *testing.T) {
	mock := newMock(errors.New("duplicate entry"))

	err := Insert(mock, "KPI-001", "desc", "sender", "receiver", "jenis")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expected := "gagal insert notifikasi"
	if len(err.Error()) < len(expected) || err.Error()[:len(expected)] != expected {
		t.Errorf("expected error message to start with '%s', got: '%s'", expected, err.Error())
	}
}

func TestInsert_AllJenisNotif(t *testing.T) {
	jenisNotifList := []string{
		"approval_penyusunan",
		"penyusunan_ditolak",
		"approval_realisasi",
		"realisasi_ditolak",
		"approval_validasi",
		"validasi_ditolak",
	}

	for _, jenis := range jenisNotifList {
		t.Run(jenis, func(t *testing.T) {
			mock := newMock(nil)

			err := Insert(mock, "KPI-2024-001", "desc notif", "USER001", "APPROVER001", jenis)

			if err != nil {
				t.Errorf("jenis '%s': expected no error, got: %v", jenis, err)
			}
			if mock.capturedArgs[4] != jenis {
				t.Errorf("jenis '%s': expected args[4]='%s', got='%v'", jenis, jenis, mock.capturedArgs[4])
			}
		})
	}
}

func TestInsert_EmptyFields(t *testing.T) {
	tests := []struct {
		name         string
		keyNotif     string
		descNotif    string
		userPengirim string
		userPenerima string
		jenisNotif   string
	}{
		{"empty keyNotif", "", "desc", "sender", "receiver", "approval_penyusunan"},
		{"empty descNotif", "KPI-001", "", "sender", "receiver", "approval_penyusunan"},
		{"empty userPengirim", "KPI-001", "desc", "", "receiver", "approval_penyusunan"},
		{"empty userPenerima", "KPI-001", "desc", "sender", "", "approval_penyusunan"},
		{"empty jenisNotif", "KPI-001", "desc", "sender", "receiver", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := newMock(nil)

			// Insert tetap dipanggil — validasi field kosong bukan tanggung jawab package ini,
			// melainkan tanggung jawab caller (service/repo).
			err := Insert(mock, tc.keyNotif, tc.descNotif, tc.userPengirim, tc.userPenerima, tc.jenisNotif)

			if err != nil {
				t.Errorf("expected no error for empty field test '%s', got: %v", tc.name, err)
			}
		})
	}
}
