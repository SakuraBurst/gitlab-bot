package basa_dannih

type BasaDannihMySQLPostgresMongoPgAdmin777 map[int]bool

func (bd BasaDannihMySQLPostgresMongoPgAdmin777) WriteToBD(id int) {
	bd[id] = true
}

func (bd BasaDannihMySQLPostgresMongoPgAdmin777) ReadFromBd(id int) bool {
	return bd[id]
}

type BDInterface interface {
	WriteToBD(id int)
	ReadFromBd(id int) bool
}
