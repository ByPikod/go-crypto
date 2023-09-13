package Database

type DBConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (d *DBConnection) Connect() {

}
