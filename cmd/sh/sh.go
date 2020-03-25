package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
	"gopkg.in/ini.v1"
)

type Request struct {
	version string
	ishook  bool
	API     string
	ROOM    string
	USER    string
	PRIV    []string
	CONTENT string
}
type Dir struct {
	roomdir    string
	cmddatadir string
	cmddir     string
}
type Conf struct {
	priv_run     []string
	priv_edit    []string
	priv_regex   []string
	priv_admin   []string
	global_cache bool
	prefix       string
}
type Cmd struct {
	name   string
	author string
	time   string
	file   string
	cache  string
}

func main() {
	req := parse()
	if req.version != "0" {
		fmt.Println("Ask bot admin to update me, This is V0 and request is " + req.version)
		return
	}
	dir := dirstr(&req)
	conf := parseconf(&dir)
	if check(&req, &conf.priv_run) {
		cmd(&req, &conf, &dir)
	}
}
func parse() Request {
	ishook, _ := strconv.ParseBool(os.Args[2])
	priv := strings.Split(os.Args[4], ",")
	req := Request{
		version: os.Args[1],
		ishook:  ishook,
		API:     os.Args[3],
		ROOM:    os.Args[5],
		USER:    os.Args[6],
		PRIV:    priv,
		CONTENT: os.Args[7],
	}
	return req
}
func dirstr(req *Request) Dir {
	roomdir := "data/" + req.API + "/" + req.ROOM + "/"
	cmddir := "cmd/sh/"
	return Dir{
		roomdir:    roomdir,
		cmddatadir: roomdir + cmddir,
		cmddir:     cmddir,
	}
}
func parseconf(dir *Dir) Conf {
	c, err := ini.Load(dir.roomdir + "user.ini")
	if err != nil {
		fmt.Println(err)
	}
	c2, err := ini.Load("mashiron.ini")
	if err != nil {
		fmt.Println(err)
	}
	return Conf{
		priv_edit:    c.Section("sh").Key("priv_conf").Strings(" "),
		priv_run:     c.Section("sh").Key("priv_run").Strings(" "),
		priv_regex:   c.Section("sh").Key("priv_regex").Strings(" "),
		priv_admin:   c.Section("sh").Key("priv_admin").Strings(" "),
		global_cache: c2.Section("sh").Key("cache").MustBool(),
		prefix:       c.Section("core").Key("prefix").String(),
	}
}

//check privs
func check(req *Request, privs *[]string) bool {
	if len(*privs) == 0 {
		return true
	}
	for _, priv := range *privs {
		for _, req_priv := range req.PRIV {
			if req_priv == priv {
				return true
			}
		}
	}
	return false
}

func cmd(req *Request, conf *Conf, dir *Dir) {
	os.MkdirAll(dir.cmddatadir, 0777)
	db, err := bolt.Open(dir.cmddatadir+"user.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	db_default_create([]string{"hook", "cmd", "cache"}, db)
	if req.ishook {
		for _, i := range db_gen_regex(&req.CONTENT, "hook", db) {
			//run
			i := db_cmd_info(&i, db)
			vm(&req.CONTENT, dir, &i)
		}
	} else if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.") {
		//someone calls me
		if check(req, &conf.priv_edit) {
			if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.add ") {
				//add script
				req_split := strings.SplitN(req.CONTENT, " ", 4)
				req_splitline := strings.TrimLeft(strings.SplitN(req.CONTENT, "\n", 2)[0], conf.prefix+"sh.add ")
				req_splitline = strings.TrimSuffix(req_splitline, "```sh")
				req_splitline = strings.TrimSuffix(req_splitline, "```")
				req_splitline = strings.TrimSuffix(req_splitline, " ")
				if strings.HasSuffix(req_split[0], "\n") || req_splitline == "" {
					fmt.Println("> Please include script name before file.")
					return
				}
				if db_cmd_exists(&req_splitline, db) {
					fmt.Println("> Script already exists.")
				} else {
					index := 2
					cmd := Cmd{
						author: req.USER,
						cache:  "true",
						time:   time.Now().String(),
						name:   req_splitline,
					}
					if req_split[index] == "--no-cache" {
						cmd.cache = "false"
					} else {
						cmd.cache = "true"
					}
					//trim
					cmdtmp := strings.SplitN(req.CONTENT, "\n", 3)
					c := 0
					for {
						if strings.HasSuffix(cmdtmp[c], "```sh") || strings.HasSuffix(cmdtmp[c], "```") {
							cmd.file = strings.Join(cmdtmp[c+1:], "\n")
							break
						} else if c == 3 {
							fmt.Println("> Cannot find script...?")
							return
						} else {
							c++
						}
					}
					cmd.file = strings.TrimRight(cmd.file, "```")
					out, _ := exec.Command(dir.cmddir+"shchk.sh", cmd.file).Output()
					fmt.Print(string(out))
					db_cmd_add(&cmd, db)
					fmt.Println("> Added script. Type `" + conf.prefix + "sh.info " + cmd.name + "` for details.")
				}
			}
			if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.rm ") {
				//delete script
				req_split := strings.SplitN(req.CONTENT, " ", 2)
				if len(req_split) != 2 {
					fmt.Println("> Request split error.")
				} else if db_cmd_exists(&req_split[1], db) {
					info := db_cmd_info(&req_split[1], db)
					if info.author == req.USER || check(req, &conf.priv_admin) {
						db_cmd_rm(&req_split[1], db)
						fmt.Println("> Deleted `" + info.name + "` .")
					} else {
						fmt.Println("> You are not allowed to delete this command.")
					}
				} else {
					fmt.Println("> No such command in database.")
				}
			}
			if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.info ") {
				req_split := strings.SplitN(req.CONTENT, " ", 2)
				if len(req_split) != 2 {
					fmt.Println("> Request split error.")
				} else if db_cmd_exists(&req_split[1], db) {
					info := db_cmd_info(&req_split[1], db)
					fmt.Printf(">>> Name: `" + info.name + "`\nBy: `" + info.author + "`\nAt: `" + info.time + "`\nCache: `" + info.cache + "`\n File:```sh\n" + info.file + "```\n")
				} else {
					fmt.Println("> No such script in database.")
				}
			}
			if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.ls") {
				list := db_cmd_list(db)
				if len(list) == 0 {
					fmt.Println("> There are no script in database.")
					return
				}
				fmt.Println("> There are " + strconv.Itoa(len(list)) + " command(s) in database.")
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
			if check(req, &conf.priv_regex) {
				if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.hook.add ") {
					//add hook regex
					req_split := strings.SplitN(req.CONTENT, " ", 3)
					if len(req_split) != 3 {
						fmt.Println("> Request split error.")
					} else if db_cmd_exists(&req_split[2], db) {
						_, err := regexp.Compile(req_split[1])
						if err != nil {
							fmt.Println("> Regex error.")
							fmt.Println(err.Error())
						} else {
							db_gen_add(&req_split[1], &req_split[2], "hook", db)
							fmt.Println(">>> Added to DB.")
							fmt.Println("RegEx: `" + req_split[1] + "`")
							fmt.Println("Cmd: `" + req_split[2] + "`")
						}
					} else {
						fmt.Println("> No such script in database.")
					}
				}
				if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.hook.rm ") {
					//delete hook regex
					req_split := strings.SplitN(req.CONTENT, " ", 2)
					if len(req_split) != 2 {
						fmt.Println("> Request split error.")
					} else if db_gen_exists(&req_split[1], "hook", db) {
						db_gen_rm(&req_split[1], "hook", db)
						fmt.Println("> Deleted `" + req_split[1] + "`.")
					} else {
						fmt.Println("> No such regex in database.")
					}
				}
				if strings.HasPrefix(req.CONTENT, conf.prefix+"sh.hook.ls") {
					//hook listing
					list := db_gen_list("hook", db)
					if len(list) == 0 {
						fmt.Println("> There are no regex in database.")
						return
					}
					fmt.Println("> There are " + strconv.Itoa(len(list)) + " regex(s) in database.")
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
						fmt.Println("> " + i)
					}
				}
			}
			if check(req, &conf.priv_run) {
				req_split := strings.SplitN(req.CONTENT, " ", 2)
				c := strings.Split(req_split[0], ".")
				req_cmd := c[len(c)-1]
				if db_cmd_exists(&req_cmd, db) {
					info := db_cmd_info(&req_cmd, db)
					v := ""
					if len(req_split) > 1 {
						v = req_split[1]
					}
					vm(&v, dir, &info)
				}
			}
		}
	}
}

func db_cmd_add(cmd *Cmd, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("cmd"))
		if err != nil {
			fmt.Println(err)
		}
		bucket, err := root.CreateBucketIfNotExists([]byte(cmd.name))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte("author"), []byte(cmd.author))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte("time"), []byte(cmd.time))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte("file"), []byte(cmd.file))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte("cache"), []byte(cmd.cache))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
}
func db_cmd_info(name *string, db *bolt.DB) Cmd {
	var cmd Cmd
	cmd.name = *name
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("cmd"))
		if root == nil {
			fmt.Println("bucket not found")
		}
		bucket := root.Bucket([]byte(*name))
		if bucket == nil {
			fmt.Println("cmd not found")
		}
		cmd.author = string(bucket.Get([]byte("author")))
		cmd.time = string(bucket.Get([]byte("time")))
		cmd.file = string(bucket.Get([]byte("file")))
		cmd.cache = string(bucket.Get([]byte("cache")))
		return nil
	})
	return cmd
}

func db_cmd_rm(name *string, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("cmd"))
		if err != nil {
			fmt.Println(err)
		}
		err = root.DeleteBucket([]byte(*name))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
}

func db_cmd_exists(cmd *string, db *bolt.DB) bool {
	res := true
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("cmd"))
		if root == nil {
			fmt.Println("bucket not found")
		}
		if root.Bucket([]byte(*cmd)) == nil {
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

func db_gen_add(key *string, value *string, bucketname string, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Put([]byte(*key), []byte(*value))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
}

func db_gen_rm(key *string, bucketname string, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketname))
		if err != nil {
			fmt.Println(err)
		}
		err = bucket.Delete([]byte(*key))
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
}

func db_gen_list(bucketname string, db *bolt.DB) []string {
	res := make([]string, 0)
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(bucketname))
		if root == nil {
			fmt.Println("bucket not found")
		}
		root.ForEach(func(key []byte, value []byte) error {
			res = append(res, string(key)+" => "+string(value))
			return nil
		})
		return nil
	})
	return res
}

func db_gen_regex(req *string, bucket string, db *bolt.DB) []string {
	res := make([]string, 0)
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(bucket))
		if root == nil {
			fmt.Println("bucket not found")
		}
		root.ForEach(func(key []byte, value []byte) error {
			hit, _ := regexp.MatchString(string(key), *req)
			if hit {
				res = append(res, string(value))
			}
			return nil
		})
		return nil
	})
	return res
}

func db_gen_match(req *string, bucket string, db *bolt.DB) string {
	var res string
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(bucket))
		if root == nil {
			fmt.Println("bucket not found")
		}
		root.ForEach(func(key []byte, value []byte) error {
			if *req == string(key) {
				res = string(value)
				return nil
			}
			return nil
		})
		return nil
	})
	return res
}

func db_gen_exists(key *string, bucket string, db *bolt.DB) bool {
	res := true
	db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(bucket))
		if root == nil {
			fmt.Println("bucket not found")
		}
		if root.Get([]byte(*key)) == nil {
			res = false
		}
		return nil
	})
	return res
}

func db_default_create(defaultarr []string, db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		for _, d := range defaultarr {
			_, err := tx.CreateBucketIfNotExists([]byte(d))
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil
	})
}

func vm(req *string, dir *Dir, cmd *Cmd) {
	//Systemd-nspawn needs root priv.
	c := exec.Command("sudo", append([]string{dir.cmddir + "run.sh", cmd.file, ""}, strings.Split(*req, " ")...)...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stdout
	c.Run()
}
