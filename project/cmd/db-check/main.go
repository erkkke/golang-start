package main

import (
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/models"
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

//	query := `SELECT cp.*,
//       crt.id "certificates.id",
//       crt.name "certificates.name",
//       crt.real_price "certificates.real_price",
//       crt.discount "certificates.discount",
//       crt.price_with_discount "certificates.price_with_discount",
//       crt.count "certificates.count",
//       crt.count_of_sales "certificates.count_of_sales"
//FROM coupons AS cp
//         INNER JOIN certificates AS crt
//                    ON cp.id = crt.coupon_id;`

	//query2 :=

	getCouponsQuery := `SELECT * FROM coupons`
	if err := conn.Select(&coupons, getCouponsQuery); err != nil {
		panic(err)
	}

	getCertificatesQuery := `SELECT id, name, real_price, discount, price_with_discount, count, count_of_sales FROM certificates WHERE coupon_id = $1`
	for _, cp := range coupons {
		if err := conn.Select(&cp.Certificates, getCertificatesQuery, cp.ID); err != nil {
			panic(err)
		}
	}
	res, err := json.Marshal(coupons)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

}
