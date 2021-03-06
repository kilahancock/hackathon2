package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func newDB(user, pass string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(mysql:3306)/testDB?parseTime=true", user, pass)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DataStore) GetPerson(uname string, pwd string) (Person, error){
	res := new(Person)

	rows, err := d.db.Query(`SELECT * FROM Person WHERE username=?`,uname)
	if err != nil {
		return Person{}, err
	}

	if ok := rows.Next(); !ok{
		return Person{}, errors.New(fmt.Sprintf("Didn't find this user: %s",uname))
	}

	err = rows.Scan(&res.Id, &res.Username, &res.Password, &res.Email, &res.Zipcode)
	if err != nil {
		return Person{}, err
	}

	if err = rows.Err(); err != nil {
		return Person{}, err
	}

	return *res, err
}

func (d *DataStore) SavePerson(p Person) (int64, error) {

	// TODO delete
	d.Logger.Info().Msg(fmt.Sprintf("We're adding this person to the db: %+v", p))

	res, err := d.db.Exec(`INSERT INTO Person (username,password,email,zipcode)
VALUES (?, ?, ?, ?)`,p.Username, p.Password, p.Email, p.Zipcode)

	if err != nil{
		return 0, err
	}
	pid, err := res.LastInsertId()

	return pid, err
}

func (d *DataStore) SaveResource(r Resource) (int64, error) {

	// TODO delete
	d.Logger.Info().Msg(fmt.Sprintf("We're adding this resource to the db: %+v", r))

	res, err := d.db.Exec(`INSERT INTO Resources (pid, rname, rtype, request, dsc, zipcode) VALUES (?, ?, ?, ?, ?, ?)`,
		r.Pid, r.Rname, r.Rtype, r.Request, r.Dsc, r.Zipcode)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

func (d *DataStore) GetResourceByZip(zipCode string) ([]Resource, error){
	res := make([]Resource, 0)

	rows, err := d.db.Query(`SELECT * FROM Resources`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rs := new(Resource)
		err := rows.Scan(&rs.Id, &rs.Rname, &rs.Pid, &rs.Request, &rs.Rtype, &rs.Dsc, &rs.Zipcode)
		if err != nil {
			return nil, err
		}
		res = append (res, *rs)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, err
}

func (d *DataStore) SaveCharity(c Charity) (int64, error) {

	// TODO delete
	d.Logger.Info().Msg(fmt.Sprintf("We're adding this charity to the db: %+v", c))

	res, err := d.db.Exec(`INSERT INTO Charities (id, pid, cname, cURL, ccity, cstate) VALUES (?, ?, ?, ?, ?, ?)`,
		c.Id, c.Pid, c.Cname, c.CURL, c.Ccity, c.Cstate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

func (d *DataStore) GetCharitiesByUser(pid int) ([]Charity, error){
	res := make([]Charity, 0)

	rows, err := d.db.Query(`SELECT * FROM Charities WHERE pid=?`,pid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ch := new(Charity)
		err := rows.Scan(&ch.Id, &ch.Pid, &ch.Cname, &ch.CURL, &ch.Ccity, &ch.Cstate)
		if err != nil {
			return nil, err
		}
		res = append (res, *ch)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, err
}



