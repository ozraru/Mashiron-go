package mashiron

import (
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"os"
	"regexp"
)

type DB struct {
	id string
	file *bbolt.DB
}

func DB_CreateRootBacket(root string,dir *Dir) {
	db := DBLoader(dir)
	db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		return nil
	})
}

func DB_AddBucket(root string,dir *Dir,BucketName string,Content [][]string) {
	db := DBLoader(dir)
	db.Update(func(tx *bbolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		var bucket *bbolt.Bucket
		if BucketName != "" {
			bucket, err = root.CreateBucketIfNotExists([]byte(BucketName))
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			}
		} else {
			bucket = root
		}
		for _,c := range Content {
			err = bucket.Put([]byte(c[0]), []byte(c[1]))
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			}
		}
		return nil
	})
}

func DB_IfBucketExists(root string,dir *Dir,cmd string) bool {
	db := DBLoader(dir)
	res := false
	db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte(root))
		if root == nil {
			fmt.Fprint(os.Stderr, "bucket not found")
			return nil
		}
		if root.Bucket([]byte(cmd)) != nil {
			res = true
		} else if root.Get([]byte(cmd)) != nil{
			res = true
		}
		return nil
	})
	return res
}

func DB_GetBucket(root string,dir *Dir,name string,key []string) []string {
	db := DBLoader(dir)
	value := make([]string,len(key))
	db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte(root))
		if root == nil {
			fmt.Fprint(os.Stderr, "bucket not found")
		}
		var bucket *bbolt.Bucket
		if name != "" {
			bucket = root.Bucket([]byte(name))
		} else {
			bucket = root
		}
		if bucket == nil {
			return nil
		}
		for i,v := range key {
			value[i] = string(bucket.Get([]byte(v)))
		}
		return nil
	})
	return value
}

func DB_DeleteBucket(root string,dir *Dir, NestedBucket string, name string) {
	db := DBLoader(dir)
	db.Update(func(tx *bbolt.Tx) error {
		bucket,_ := tx.CreateBucketIfNotExists([]byte(root))
		if NestedBucket != "" {
			bucket,_ = bucket.CreateBucketIfNotExists([]byte(NestedBucket))
		}
		v := bucket.Get([]byte(name))
		var err error
		if v == nil {
			err = bucket.DeleteBucket([]byte(name))
		} else {
			err = bucket.Delete([]byte(name))
		}
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
		return nil
	})
}

func DB_GetBucketList(root string,dir *Dir) [][]string {
	db := DBLoader(dir)
	var res [][]string
	db.View(func(tx *bbolt.Tx) error {
		root := tx.Bucket([]byte(root))
		if root == nil {
			fmt.Fprint(os.Stderr, "bucket not found")
		}
		root.ForEach(func(k []byte, v []byte) error {
			res = append(res, []string{string(k),string(v)})
			return nil
		})
		return nil
	})
	return res
}

//Use DB_GetBucket if you know the name of keys.
func DB_GetFullKVList(root string,dir *Dir, name string) [][]string {
	db := DBLoader(dir)
	var res [][]string
	db.View(func (tx *bbolt.Tx) error{
		root := tx.Bucket([]byte(root))
		if root == nil {
			return errors.New("root bucket not found")
		}
		bucket := root.Bucket([]byte(name))
		if bucket == nil {
			return errors.New("root bucket not found")
		}
		bucket.ForEach(func(k []byte, v []byte) error {
			res = append(res,[]string{string(k),string(v)})
			return nil
		})
		return nil
	})
	return res
}

func DB_Regex(root string,BucketName string,req string, dir *Dir) []string {
	db := DBLoader(dir)
	res := make([]string, 0)
	db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(root))
		if bucket == nil {
			return errors.New("bucket not found")
		}
		if BucketName != "" {
			bucket := bucket.Bucket([]byte(BucketName))
			if bucket == nil {
				fmt.Fprint(os.Stderr, "bucket not found")
			}
		}
		bucket.ForEach(func(key []byte, value []byte) error {
			hit, _ := regexp.MatchString(string(key), req)
			if hit {
				res = append(res, string(value))
			}
			return nil
		})
		return nil
	})
	return res
}

var dbs []DB
func DBLoader(dir *Dir) *bbolt.DB {
	for _,v := range dbs {
		if v.id == dir.CmdDataDir && v.file != nil {
			return v.file
		}
	}
	db, err := bbolt.Open(dir.CmdDataDir+"user.db", 0600, nil)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	dbs = append(dbs,DB{
		id:   dir.CmdDataDir,
		file: db,
	})
	return db
}
