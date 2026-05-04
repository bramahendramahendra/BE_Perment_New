package excel

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/xuri/excelize/v2"
)

// =============================================================================
// Helper
// =============================================================================

// newExcelFileHeader membuat multipart.FileHeader dari konten xlsx in-memory.
func newExcelFileHeader(t *testing.T, filename string, buildFn func(*excelize.File)) *multipart.FileHeader {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()
	buildFn(f)

	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatalf("gagal menulis excel ke buffer: %v", err)
	}

	return newFileHeader(t, filename, buf.Bytes())
}

// newFileHeader membuat multipart.FileHeader dari raw bytes.
func newFileHeader(t *testing.T, filename string, data []byte) *multipart.FileHeader {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	h.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	fw, err := mw.CreatePart(h)
	if err != nil {
		t.Fatalf("gagal membuat part: %v", err)
	}
	fw.Write(data)
	mw.Close()

	mr := multipart.NewReader(&b, mw.Boundary())
	form, err := mr.ReadForm(10 << 20)
	if err != nil {
		t.Fatalf("gagal membaca form: %v", err)
	}
	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("tidak ada file dalam form")
	}
	return files[0]
}

// =============================================================================
// TestGetCell
// =============================================================================

func TestGetCell(t *testing.T) {
	tests := []struct {
		name  string
		row   []string
		index int
		want  string
	}{
		{
			name:  "index valid, value ada",
			row:   []string{"A", "B", "C"},
			index: 1,
			want:  "B",
		},
		{
			name:  "index 0",
			row:   []string{"first"},
			index: 0,
			want:  "first",
		},
		{
			name:  "index melebihi panjang row",
			row:   []string{"A", "B"},
			index: 5,
			want:  "",
		},
		{
			name:  "row kosong",
			row:   []string{},
			index: 0,
			want:  "",
		},
		{
			name:  "value ada spasi — harus di-trim",
			row:   []string{"  hello  "},
			index: 0,
			want:  "hello",
		},
		{
			name:  "value hanya spasi — return kosong",
			row:   []string{"   "},
			index: 0,
			want:  "",
		},
		{
			name:  "index tepat di batas akhir",
			row:   []string{"A", "B", "C"},
			index: 2,
			want:  "C",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCell(tt.row, tt.index)
			if got != tt.want {
				t.Errorf("GetCell() = %q, want %q", got, tt.want)
			}
		})
	}
}

// =============================================================================
// TestParseFloat
// =============================================================================

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:    "string kosong — return 0 tanpa error",
			input:   "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "angka bulat",
			input:   "25",
			want:    25.0,
			wantErr: false,
		},
		{
			name:    "angka desimal",
			input:   "12.50",
			want:    12.50,
			wantErr: false,
		},
		{
			name:    "angka dengan persen — strip persen",
			input:   "75.5%",
			want:    75.5,
			wantErr: false,
		},
		{
			name:    "angka dengan spasi",
			input:   "  30.00  ",
			want:    30.0,
			wantErr: false,
		},
		{
			name:    "angka negatif",
			input:   "-10.5",
			want:    -10.5,
			wantErr: false,
		},
		{
			name:    "presisi 2 desimal — rounding",
			input:   "3.145",
			want:    3.15,
			wantErr: false,
		},
		{
			name:    "bukan angka — return error",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "teks campuran — return error",
			input:   "12abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "nol",
			input:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "angka besar",
			input:   "1000000.99",
			want:    1000000.99,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFloat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseFloat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// =============================================================================
// TestNullableString
// =============================================================================

func TestNullableString(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		isActive bool
		wantNil  bool
		wantVal  string
	}{
		{
			name:     "isActive true — return pointer berisi nilai",
			val:      "hello",
			isActive: true,
			wantNil:  false,
			wantVal:  "hello",
		},
		{
			name:     "isActive false — return nil",
			val:      "hello",
			isActive: false,
			wantNil:  true,
		},
		{
			name:     "isActive true, val kosong — return pointer string kosong",
			val:      "",
			isActive: true,
			wantNil:  false,
			wantVal:  "",
		},
		{
			name:     "isActive false, val kosong — return nil",
			val:      "",
			isActive: false,
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullableString(tt.val, tt.isActive)
			if tt.wantNil {
				if got != nil {
					t.Errorf("NullableString() = %v, want nil", *got)
				}
				return
			}
			if got == nil {
				t.Errorf("NullableString() = nil, want %q", tt.wantVal)
				return
			}
			if *got != tt.wantVal {
				t.Errorf("NullableString() = %q, want %q", *got, tt.wantVal)
			}
		})
	}
}

// =============================================================================
// TestReadSheet
// =============================================================================

func TestReadSheet(t *testing.T) {
	t.Run("sheet ada — return rows", func(t *testing.T) {
		fh := newExcelFileHeader(t, "test.xlsx", func(f *excelize.File) {
			f.SetCellValue("Sheet1", "A1", "KPI")
			f.SetCellValue("Sheet1", "B1", "SubKPI")
			f.SetCellValue("Sheet1", "A2", "KPI-001")
			f.SetCellValue("Sheet1", "B2", "Sub-001")
		})

		rows, err := ReadSheet(fh, "Sheet1")
		if err != nil {
			t.Fatalf("ReadSheet() error = %v", err)
		}
		if len(rows) < 2 {
			t.Errorf("ReadSheet() jumlah baris = %d, want >= 2", len(rows))
		}
		if rows[0][0] != "KPI" {
			t.Errorf("baris 1 kolom A = %q, want %q", rows[0][0], "KPI")
		}
		if rows[1][0] != "KPI-001" {
			t.Errorf("baris 2 kolom A = %q, want %q", rows[1][0], "KPI-001")
		}
	})

	t.Run("sheet tidak ada — return error", func(t *testing.T) {
		fh := newExcelFileHeader(t, "test.xlsx", func(f *excelize.File) {
			f.SetCellValue("Sheet1", "A1", "data")
		})

		_, err := ReadSheet(fh, "TIDAKADA")
		if err == nil {
			t.Error("ReadSheet() expected error untuk sheet tidak ada, got nil")
		}
	})

	t.Run("file bukan excel valid — return error", func(t *testing.T) {
		fh := newFileHeader(t, "bukan_excel.xlsx", []byte("ini bukan xlsx"))

		_, err := ReadSheet(fh, "Sheet1")
		if err == nil {
			t.Error("ReadSheet() expected error untuk file tidak valid, got nil")
		}
	})

	t.Run("sheet kosong — return rows kosong", func(t *testing.T) {
		fh := newExcelFileHeader(t, "test.xlsx", func(f *excelize.File) {
			f.NewSheet("KOSONG")
		})

		rows, err := ReadSheet(fh, "KOSONG")
		if err != nil {
			t.Fatalf("ReadSheet() error = %v", err)
		}
		if len(rows) != 0 {
			t.Errorf("ReadSheet() jumlah baris = %d, want 0", len(rows))
		}
	})

	t.Run("sheet dengan nama TW1 — berhasil dibaca", func(t *testing.T) {
		fh := newExcelFileHeader(t, "kpi.xlsx", func(f *excelize.File) {
			f.NewSheet("TW1")
			f.SetCellValue("TW1", "A1", "Header")
			f.SetCellValue("TW1", "A2", "Data1")
		})

		rows, err := ReadSheet(fh, "TW1")
		if err != nil {
			t.Fatalf("ReadSheet() error = %v", err)
		}
		if len(rows) < 2 {
			t.Errorf("ReadSheet() jumlah baris = %d, want >= 2", len(rows))
		}
	})
}

// =============================================================================
// Benchmark
// =============================================================================

func BenchmarkParseFloat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseFloat("12.50")
	}
}

func BenchmarkGetCell(b *testing.B) {
	row := []string{"A", "B", "C", "D", "E"}
	for i := 0; i < b.N; i++ {
		GetCell(row, 3)
	}
}
