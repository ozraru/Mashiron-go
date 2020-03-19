package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
	"gopkg.in/ini.v1"
)

type Request struct {
	version string
	API     string
	ROOM    string
	USER    string
	PRIV    []string
	CONTENT string
}
type Dir struct {
	roomdir    string
	cmddatadir string
}
type Conf struct {
	priv_read  []string
	priv_edit  []string
	priv_admin []string
	prefix     string
}
type Man struct {
	name   string
	author string
	time   string
	file   string
}

func main() {
	req := parse()
	if req.version != "0" {
		fmt.Println("Ask bot admin to update me, This is V0 and request is " + req.version)
		return
	}
	dir := dirstr(req)
	conf := parseconf(dir)
	if check(req, conf.priv_read) {
		cmd(req, conf, dir)
	}
}
func parse() Request {
	priv := strings.Split(os.Args[4], ",")
	req := Request{
		version: os.Args[1],
		API:     os.Args[3],
		ROOM:    os.Args[5],
		USER:    os.Args[6],
		PRIV:    priv,
		CONTENT: os.Args[7],
	}
	return req
}
func dirstr(req Request) Dir {
	roomdir := "data/" + req.API + "/" + req.ROOM + "/"
	cmddir := "cmd/man/"
	return Dir{
		roomdir:    roomdir,
		cmddatadir: roomdir + cmddir,
	}
}
func parseconf(dir Dir) Conf {
	c, err := ini.Load(dir.roomdir + "user.ini")
	if err != nil {
		fmt.Println(err)
	}
	return Conf{
		priv_edit:  c.Section("man").Key("priv_edit").Strings(" "),
		priv_read:  c.Section("man").Key("priv_read").Strings(" "),
		priv_admin: c.Section("man").Key("priv_admin").Strings(" "),
		prefix:     c.Section("core").Key("prefix").String(),
	}
}
func check(req Request, privs []string) bool {
	if len(privs) == 0 {
		return true
	}
	for _, priv := range privs {
		for _, req_priv := range req.PRIV {
			if req_priv == priv {
				return true
			}
		}
	}
	return false
}
func cmd(req Request, conf Conf, dir Dir) {
	os.MkdirAll(dir.cmddatadir, 0777)
	db, err := bolt.Open(dir.cmddatadir+"user.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	db_default_create("man", db)
	if strings.HasPrefix(req.CONTENT, conf.prefix+"man.") {
		//someone calls me
		if check(req, conf.priv_edit) {
			if strings.HasPrefix(req.CONTENT, conf.prefix+"man.add") {
				//add command
				req_split := strings.SplitN(req.CONTENT, " ", 4)
				if strings.Contains(req_split[1], "\n") {
					fmt.Println("> Please include man name before file.")
					return
				}
				if db_cmd_exists(req_split[1], db) {
					fmt.Println("> Command already exists.")
				} else {
					man := Man{
						author: req.USER,
						time:   time.Now().String(),
						name:   req_split[1],
					}
					//trim
					cmdtmp := strings.SplitN(req.CONTENT, "\n", 5)
					c := 0
					for {
						if strings.HasSuffix(cmdtmp[c], "```") {
							man.file = strings.Join(cmdtmp[c+1:], "")
							break
						} else if c == 5 {
							fmt.Println("> Cannot find man...?")
							return
						} else {
							c++
						}
					}
					man.file = strings.TrimRight(man.file, "```")
					db_cmd_add(man, db)
					fmt.Println("> Added man. Type `" + conf.prefix + "man." + man.name + "` for details.")
				}
			}
			if strings.HasPrefix(req.CONTENT, conf.prefix+"man.rm ") {
				//delete command
				req_split := strings.SplitN(req.CONTENT, " ", 2)
				if len(req_split) != 2 {
					fmt.Println("> Request split error.")
				} else if db_cmd_exists(req_split[1], db) {
					info := db_cmd_info(req_split[1], db)
					if info.author == req.USER || check(req, conf.priv_admin) {
						db_cmd_rm(req_split[1], db)
						fmt.Println("> Deleted `" + info.name + "` .")
					} else {
						fmt.Println("> You are not allowed to delete this command.")
					}
				} else {
					fmt.Println("> No such command in database.")
				}
			}
			if strings.HasPrefix(req.CONTENT, conf.prefix+"man.ls") {
				list := db_cmd_list(db)
				if len(list) == 0 {
					fmt.Println("> There are no mans in database.")
					return
				}
				fmt.Println("> There are " + strconv.Itoa(len(list)) + " mans in database.")
				req_split := strings.SplitN(req.CONTENT, " ", 2)
				if len(list) > 10 {
					page := 1
					if len(req_split) == 2 {
						page, err = strconv.Atoi(req_split[1])
						if err == nil {
							if page > 0 {
								if page*10-1 > len(list) {
									list = list[(page-1)*10 : len(list)]
								} else {
									list = list[(page-1)*10 : page*10-1]
								}
							}
						}
					}
					fmt.Println("> Pager enabled,Current page is: " + string(page))
				}
				for _, i := range list {
					fmt.Println(i)
				}
			}
			if check(req, conf.priv_read) {
				req_split := strings.SplitN(req.CONTENT, " ", 2)
				c := strings.Split(req_split[0], ".")
				req_cmd := c[len(c)-1]
				if db_cmd_exists(req_cmd, db) {
					i := db_cmd_info(req_cmd, db)
					fmt.Print(">>> Man file `" + i.name + "` \n```\nAuthor: `" + i.author + "\nTime: " + i.time + "\n```\nMan:\n" + i.file)
				}
			}
		}
	}
}
func db_default_create(defaultstr string, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultstr))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
}
func db_cmd_add(man Man, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("man"))
		if err != nil {
			fmt.Println(err)
		}
		bucket, err := root.CreateBucketIfNotExists([]byte(man.name))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte("author"), []byte(man.author))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte("time"), []byte(man.time))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte("file"), []byte(man.file))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
}
func db_cmd_info(name string, db *bolt.DB) Man {
	var man Man
	man.name = name
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("man"))
		if root == nil {
			fmt.Println("bucket not found")
		}
		bucket := root.Bucket([]byte(name))
		if bucket == nil {
			fmt.Println("man not found")
		}
		man.author = string(bucket.Get([]byte("author")))
		man.time = string(bucket.Get([]byte("time")))
		man.file = string(bucket.Get([]byte("file")))
		return nil
	})
	return man
}

func db_cmd_rm(name string, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("man"))
		err := root.DeleteBucket([]byte(name))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
}

func db_cmd_exists(cmd string, db *bolt.DB) bool {
	res := true
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("man"))
		if root == nil {
			fmt.Println("bucket not found")
		}
		if root.Bucket([]byte(cmd)) == nil {
			res = false
		}
		return nil
	})
	return res
}

func db_cmd_list(db *bolt.DB) []string {
	res := make([]string, 0)
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("cmd"))
		if root == nil {
			fmt.Println("bucket not found")
		}
		root.ForEach(func(name []byte, v []byte) error {
			res = append(res, string(name))
			return nil
		})
		return nil
	})
	return res
}
