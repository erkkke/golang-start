package main

import (
	"fmt"
	"github.com/erkkke/golang-start/hw6/internal/models"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	//ctx := context.Background()
	urlExample := "postgres://postgres:postgres@localhost:5432/golang_project"
	conn, err := sqlx.Connect("pgx", urlExample)
	if err != nil {
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}

	//_, err = conn.ExecContext(ctx,
	//	"INSERT INTO coupons(category, company, name, description, address, phone_number) " +
	//	"VALUES ($1, $2, $3, $4, $5, $6)",
	//	1, "coupon.Company", "coupon.Name", "coupon.Description", "coupon.Address", "c123123",
	//)
	//if err != nil {
	//	panic(err)
	//}
	//
	//var id []int
	//err = conn.Select(&id, "SELECT currval(pg_get_serial_sequence('coupons', 'id'))")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(id)
	coupons := make([]*models.Coupon, 0)

	err = conn.Select(&coupons,
		"SELECT cp.id, category, company, cp.name, description, crt.id, crt.name, crt.real_price, crt.discount, crt.price_with_discount, crt.count, crt.count_of_sales, address, phone_number "+
			"FROM coupons AS cp "+
			"INNER JOIN certificates AS crt "+
			"ON cp.id = crt.coupon_id")
	if err != nil {
		panic(err)
	}

	fmt.Println(coupons)

}
