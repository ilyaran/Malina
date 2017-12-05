/**
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package entity


import (
	"time"
	"database/sql"
)

type Account struct {
	Id       int64       	`json:"id"`
	Role     int64       	`json:"role"`
	Email    string      	`json:"email"`
	Phone    string      	`json:"phone"`
	Nick     string        	`json:"nick"`
	FirstName     string   	`json:"first_name"`
	LastName     string    	`json:"last_name"`
	Sex     bool           	`json:"sex"`
	Skype    string        	`json:"skype"`
	Steam    string        	`json:"steam"`
	Trade    string        	`json:"trade"`
	Game     int64         	`json:"game"`
	State    int64         	`json:"state"`
	Img string             	`json:"img"`

	Ban bool        		`json:"ban"`
	BanReason string		`json:"ban_reason"`

	Created time.Time		`json:"created"`
	Updated time.Time		`json:"updated"`

	LastLogged time.Time	`json:"last_logged"`

	Team     []int64       	`json:"team"`


	Balance float64        	`json:"-"`
	LastIp string			`json:"-"`
	Provider    string      `json:"-"`
	Token    string      	`json:"-"`
	Password string      	`json:"-"`
	Newpass string      	`json:"-"`
}

func (s *Account) Appending(rows *sql.Rows,size int64)*[]Scanable{
	list := make([]Scanable, 0, size)
	for rows.Next() {
		item := &Account{}
		if item.Scanning(nil, rows) == 'o' {
			list = append(list, item)
		}else {
			return nil
		}
	}
	return &list
}
func (s *Account) Item(row *sql.Row) Scanable {
	item := Account{}
	if item.Scanning(row, nil) == 'o' {
		return &item
	}
	return nil
}


func (s *Account) GetId() int64 {return s.Id}
func (s *Account) SetId(id int64) {s.Id = id}

func (s *Account) Scanning(row *sql.Row,rows *sql.Rows) byte {
	var err error
	if row!=nil {
		err = row.Scan(
			&s.Id,
			&s.Role,
			&s.Email,
			&s.Phone,
			&s.Nick,
			&s.FirstName,
			&s.LastName,
			&s.Skype,
			&s.State,
			&s.Img,
			&s.Balance,
			&s.Ban,
			&s.BanReason,
			&s.Created,
			&s.Updated,
			&s.LastLogged,
			&s.LastIp,
			&s.Token,
			&s.Password,
			&s.Provider,
		)
	}
	if rows!=nil {
		err = rows.Scan(

			&s.Id,
			&s.Role,
			&s.Email,
			&s.Phone,
			&s.Nick,
			&s.FirstName,
			&s.LastName,
			&s.Skype,
			&s.State,
			&s.Img,
			&s.Balance,
			&s.Ban,
			&s.BanReason,
			&s.Created,
			&s.Updated,
			&s.LastLogged,
			&s.LastIp,
			&s.Token,
			&s.Password,
			&s.Provider,
		)
	}
	if err == sql.ErrNoRows {
		//panic(err)
		return 'n'
	}
	if err != nil {
		panic(err)
		return 'e'
	}

	return 'o'
}

