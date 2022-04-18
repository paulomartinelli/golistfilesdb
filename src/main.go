package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {

	sqlStatement := `INSERT INTO public.arquivos
	(file_md5, file_path, file_name, file_entry, ext, ext_u)
	VALUES($1, $2, $3, $4, $5, UPPER( $6 ));
	`

	args := os.Args[1:]
	fmt.Println(args)

	db, err := sql.Open("postgres", "host=localhost port=15432 user=root "+
		"password=root dbname=my_db sslmode=disable")

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Connected!")
	}
	defer db.Close()

	root := args[0]
	err = filepath.Walk(root, visit(db, sqlStatement, root))

	if err != nil {
		log.Fatal(err)
	}

}

func visit(db *sql.DB, sqlStatement string, root string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
		} else {
			if !info.IsDir() {
				fhash, err := calculateMD5(path)

				if err == nil {
					ext := filepath.Ext(path)
					name := info.Name()

					err = insertDB(db, sqlStatement, fhash, path, name, ext, root)

					if err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println(err)
				}

			}
		}
		return nil
	}
}

func calculateMD5(file string) (string, error) {
	f, err := os.Open(file)

	if err != nil {
		return "", err
	}

	defer f.Close()

	hash := md5.New()
	io.Copy(hash, f)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func insertDB(db *sql.DB, sqlStatement string, fhash string, file string, name string, ext string, root string) error {
	insert, errst := db.Prepare(sqlStatement)
	if errst != nil {
		return errst
	}
	_, errdb := insert.Exec(fhash, file, name, root, ext, ext)

	if errdb != nil {
		return errdb
	}

	return nil
}
