package entities

type Employee struct {
	UUID      string `validate:"required"`
	Name      string `validate:"required" label:"Nama"`
	Email     string `validate:"required,email,isunique=employee-email"`
	Phone     string `validate:"required,numeric,gte=10" label:"No.Handphone"`
	Address   string `validate:"required" label:"Alamat"`
	NIK       string `validate:"required,isunique=employee-nik"`
	Gender    string `validate:"required,oneof=F M" label:"Jenis Kelamin"`
	BirthDate string `validate:"required" label:"Tanggal Lahir"`
	Photo     string `validate:"-"`
	Password  string `validate:"-"`
	IsAdmin   bool   `validate:"-"`
}

type EditEmployee struct {
	UUID      string
	Name      string `validate:"required" label:"Nama"`
	Email     string `validate:"required,email"`
	Phone     string `validate:"required,numeric,gte=10" label:"No.Handphone"`
	Address   string `validate:"required" label:"Alamat"`
	NIK       string `validate:"required"`
	Gender    string `validate:"required,oneof=F M" label:"Jenis Kelamin"`
	BirthDate string `validate:"required" label:"Tanggal Lahir"`
	IsAdmin   bool   `validate:"-"`
}

// type LoginEmployee struct {
// 	Name      string
// 	Email     string
// 	Phone     string
// 	NIK       string
// 	Gender    string
// 	BirthDate string
// 	Photo     sql.NullString
// 	Password  string
// 	IsAdmin   bool
// }
