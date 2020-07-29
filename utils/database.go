package utils

import (
	"database/sql"
	"errors"
	"fmt"

	"reflect"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/swallyu/iot-users/log"
)

// RowHandler is used to callback when query
type RowHandler func(rows []map[string]string) (interface{}, error)

// Config database config
type Config struct {
	Host         string
	Port         int
	UserName     string
	Password     string
	Dbms         string
	DatabaseName string
}

// DbUtils is used to execute sql
type DbUtils struct {
	conf *Config
	db   *sql.DB
}

// DbHolder is value data
type DbHolder struct {
	data []map[string]string
	Err  error
}

// NewDbUtils create database utils
func NewDbUtils(conf Config) DbUtils {
	return DbUtils{conf: &conf}
}

func (d *DbUtils) open() (*sql.DB, error) {

	if d.conf.Dbms == "mysql" {
		log.Info("mysql info")
		connURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.conf.UserName, d.conf.Password, d.conf.Host, d.conf.Port, d.conf.DatabaseName)
		db, err := sql.Open("mysql", connURL)
		if err != nil {
			log.WithFields(log.Fields{"host": d.conf.Host, "port": d.conf.Port, "name": d.conf.DatabaseName}).Error("connect database fail")
		}
		return db, err
	}
	return nil, fmt.Errorf("unkonwn database %s", d.conf.Dbms)
}

// Query run query
func (d *DbUtils) Query(statement string, args ...interface{}) *DbHolder {
	db, err := d.open()
	var holder = &DbHolder{}
	if err != nil {
		log.Error("open fail ")
		holder.Err = err
		return holder
	}
	defer db.Close()
	rows, err := db.Query(statement, args...)
	if err != nil {
		log.Errorf("query fail %s %v", statement, err)
		return nil
	}
	column, err := rows.Columns()
	values := make([]sql.RawBytes, len(column))
	scans := make([]interface{}, len(column))
	for i := range values {
		scans[i] = &values[i]
	}
	results := make([]map[string]string, 0)

	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			log.Error(err)
			holder.Err = err
			return holder
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {
			//每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}

	holder.data = results
	return holder
}

// Unique return single value
func (d *DbHolder) Unique(in interface{}) error {
	if len(d.data) > 0 {
		return d.mapping(d.data[0], reflect.ValueOf(in))
	}
	return nil
}

// List list value
func (d *DbHolder) List(in interface{}) error {
	if d.Err != nil {
		return d.Err
	}

	length := len(d.data)

	if length > 0 {
		v := reflect.ValueOf(in).Elem()
		newv := reflect.MakeSlice(v.Type(), 0, length)
		v.Set(newv)
		v.SetLen(length)

		index := 0
		for i := 0; i < length; i++ {
			k := v.Type().Elem()
			newObj := reflect.New(k)
			err := d.mapping(d.data[i], newObj)
			if err != nil {
				return err
			}
			v.Index(index).Set(newObj.Elem())
			index++
		}
		v.SetLen(index)
	}
	return nil
}

// Rows return rows handle
func (d *DbHolder) Rows(handle RowHandler) (interface{}, error) {
	return handle(d.data)
}

func (d *DbHolder) mapping(m map[string]string, v reflect.Value) error {
	t := v.Type()
	val := v.Elem()
	typ := t.Elem()

	if !val.IsValid() {
		return errors.New("数据类型不正确")
	}

	for i := 0; i < val.NumField(); i++ {

		value := val.Field(i)
		kind := value.Kind()
		tag := typ.Field(i).Tag.Get("col")

		if len(tag) > 0 {
			meta, ok := m[tag]
			if !ok {
				continue
			}

			if !value.CanSet() {
				return errors.New("结构体字段没有读写权限")
			}

			if len(meta) == 0 {
				continue
			}

			if kind == reflect.String {
				value.SetString(meta)
			} else if kind == reflect.Float32 {
				f, err := strconv.ParseFloat(meta, 32)
				if err != nil {
					return err
				}
				value.SetFloat(f)
			} else if kind == reflect.Float64 {
				f, err := strconv.ParseFloat(meta, 64)
				if err != nil {
					return err
				}
				value.SetFloat(f)
			} else if kind == reflect.Int64 {
				integer64, err := strconv.ParseInt(meta, 10, 64)
				if err != nil {
					return err
				}
				value.SetInt(integer64)
			} else if kind == reflect.Int {
				integer, err := strconv.Atoi(meta)
				if err != nil {
					return err
				}
				value.SetInt(int64(integer))
			} else if kind == reflect.Bool {
				b, err := strconv.ParseBool(meta)
				if err != nil {
					return err
				}
				value.SetBool(b)
			} else {
				return errors.New("数据库映射存在不识别的数据类型")
			}
		}
	}
	return nil
}
