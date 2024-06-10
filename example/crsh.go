package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	stdpath "path"
	"strconv"
	"strings"
	"time"

	json "github.com/ZxwyProject/zson"
	"github.com/ZxwyWebSite/cr-go-sdk"
	"github.com/ZxwyWebSite/cr-go-sdk/service/explorer"
	"github.com/pterm/pterm"
)

var (
	path string
	site = &cr.SiteObj{
		Addr: `https://cloudreveplus-demo.onrender.com/`,
		Users: &cr.UserObj{
			Mail: `admin@cloudreve.org`,
			Pass: `CloudrevePlusDemo`,
		},
	}
)

const help = `欢迎使用crsh - cr-go-sdk 演示程序

可尝试输入以下命令：
  # 终端
  clear - 清屏
  print <arg1> <arg2> - fmt.Printf(arg1, arg2)
  set <var> <value> - 设置变量(beta)
  sh <args> - 终端执行
  # 站点
  init <site> - 载入站点
  captcha - 获取验证码
  login <user> <pass> <?code> - 用户登录
  # 用户
  whoami - 用户详情
  users - 用户列表
  su <user> - 切换用户
  # 文件
  ls <?path> - 列出文件
  ll <?path> - 列出详情
  cd <dir> - 进入目录
  # 其它
  exit - 退出程序
  ??? - 敬请期待
`

// 相对目录
func cdpath(p string) string {
	if stdpath.IsAbs(p) {
		return stdpath.Clean(p)
	} else {
		return stdpath.Join(path, p)
	}
}

// 指令列表
var cmds = map[string]func(args ...string) error{
	// 显示帮助
	`help`: func(args ...string) error {
		println(help)
		return nil
	},
	// 清屏
	`clear`: func(args ...string) error {
		print("\033c")
		return nil
	},
	// 格式化输出
	`echo`: func(args ...string) error {
		if len(args) < 2 {
			return errors.New(`参数不足`)
		}
		_, err := fmt.Printf(args[0], args[1:])
		println("\n")
		return err
	},
	// 退出程序
	`exit`: func(args ...string) error {
		os.Exit(0)
		return nil
	},
	// 执行命令
	`sh`: func(args ...string) error {
		if len(args) == 0 {
			println(`usage:`, `sh <command...>`)
			return nil
		}
		cmd := exec.Command(`bash`, `-c`, strings.Join(args, ` `))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = cmd.Stdout
		return cmd.Run()
	},
	// 输出错误
	`err`: func(args ...string) error {
		if len(args) == 0 {
			return errors.New(`EOF`)
		}
		return errors.New(strings.Join(args, `: `))
	},
	// 特殊调用
	`cr`: func(args ...string) error {
		if len(args) == 0 {
			return nil
		}
		return nil
	},
	// 登录账号 (全参)
	`loginf`: func(args ...string) error {
		if len(args) < 3 {
			return errors.New(`参数不足`)
		}
		site.Addr = args[0]
		site.Users = &cr.UserObj{
			Mail: args[1],
			Pass: args[2],
		}
		println(`wait...`)
		err := site.SdkInit()
		if err == nil {
			err = site.SdkLogin()
		}
		if err != nil {
			site.Users = nil
			return err
		} else {
			path = `/`
			println(`userlogin:`, site.Users.Cookie.String())
		}
		return nil
	},
	// 登录账号 (无参)
	`login`: func(args ...string) error {
		println(`wait...`)
		err := site.SdkLogin()
		if err != nil {
			site.Users = nil
			return err
		} else {
			path = `/`
			println(`userlogin:`, site.Users.Cookie.String())
		}
		return nil
	},
	// 设置账号
	`iuser`: func(args ...string) error {
		length := len(args)
		if length == 0 {
			println(`usage:`, `iuser <mail> <pass>`)
			return nil
		}
		if length < 2 {
			return errors.New(`参数不足`)
		}
		site.Users = &cr.UserObj{
			Mail: args[0],
			Pass: args[1],
		}
		return nil
	},
	// 账号信息
	`user`: func(args ...string) error {
		if site.Users == nil {
			return errors.New(`暂无账号`)
		}
		mi, _ := json.MarshalIndent(site.Users, ``, `    `)
		println(string(mi))
		return nil
	},
	// 设置站点
	`isite`: func(args ...string) error {
		if len(args) == 0 {
			println(`usage:`, `isite <site>`)
			return nil
		}
		site.Addr = args[0]
		if !strings.HasSuffix(site.Addr, `/`) {
			site.Addr += `/`
		}
		return site.SdkInit()
	},
	// 站点信息
	`site`: func(args ...string) error {
		if site.Config == nil {
			return errors.New(`暂无站点`)
		}
		mi, _ := json.MarshalIndent(site.Config, ``, `    `)
		println(string(mi))
		return nil
	},
	// 文件列表
	`ls`: func(args ...string) error {
		if site.Users == nil {
			return errors.New(`未登录`)
		}
		var dir string = path
		if len(args) >= 1 {
			dir = cdpath(args[0])
		}
		list, err := site.Directory(dir)
		if err != nil {
			return err
		}
		fmt.Println(`Total`, len(list.Objects))
		table := pterm.TableData{
			{`Num`, `Type`, `ID`, `Size`, `CreateDate`, `Name`},
		}
		for i, o := range list.Objects {
			table = append(table, []string{
				strconv.Itoa(i), o.Type, o.ID,
				cr.SizeToString(o.Size), // strconv.FormatInt(int64(o.Size), 10),
				o.CreateDate.Format(time.DateTime), o.Name,
			})
		}
		return pterm.DefaultTable.WithHasHeader().WithHeaderRowSeparator(`-`).WithData(table).Render()
	},
	// 文件下载
	`dl`: func(args ...string) error {
		if len(args) == 0 {
			println(`usage:`, `dl <id>`)
			return nil
		}
		link, err := site.FileDownload(args[0])
		if err != nil {
			return err
		}
		println(*link)
		return nil
	},
	// 获取直链
	`sc`: func(args ...string) error {
		length := len(args)
		if length == 0 {
			println(`usage:`, `sc <ids...>`)
			return nil
		}
		src, err := site.FileSource(cr.GenerateSrc(false, args...))
		if err != nil {
			return err
		}
		table := pterm.TableData{
			{`Num`, `Name`, `URL`},
		}
		for i, o := range *src {
			if o.Error != `` || o.URL == `` {
				o.URL = o.Error
			}
			table = append(table, []string{
				strconv.Itoa(i), o.Name, o.URL,
			})
		}
		return pterm.DefaultTable.WithHasHeader().WithHeaderRowSeparator(`-`).WithData(table).Render()
	},
	// 文件搜索
	`find`: func(args ...string) error {
		length := len(args)
		if length == 0 {
			println(`usage:`, `find <type> <keywords> <path?>`)
			return nil
		}
		if length < 2 {
			return errors.New(`参数不足`)
		}
		var path string = `/`
		if length == 3 {
			path = args[2]
		}
		res, err := site.FileSearch(&explorer.ItemSearchService{
			Type:     args[0],
			Keywords: args[1],
			Path:     path,
		})
		if err != nil {
			return err
		}
		mi, _ := json.MarshalIndent(res, ``, `    `)
		println(string(mi))
		return nil
	},
	// 文件详情
	`file`: func(args ...string) error {
		return errors.New(`...?`)
	},
	// 切换目录
	`cd`: func(args ...string) error {
		if len(args) == 0 {
			println(`usage:`, `cd <dir>`)
			return nil
		}
		if stdpath.IsAbs(args[0]) {
			path = stdpath.Clean(args[0])
		} else {
			path = stdpath.Join(path, args[0])
		}
		return nil
	},
	// 上传文件
	`scp`: func(args ...string) error {
		length := len(args)
		if length == 0 {
			println(`usage:`, `scp <file> <dir> <name?>`)
			return nil
		}
		if length < 2 {
			return errors.New(`参数不足`)
		}
		var name string
		if length == 3 {
			name = args[2]
		}
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		// err = site.SdkUpload(args[1], file, name)
		err = func() error {
			info, err := file.Stat()
			if err != nil {
				return err
			}
			size := uint64(info.Size())
			if name == `` {
				name = info.Name()
			}
			mbuf := make([]byte, 512)
			mn, err := file.Read(mbuf)
			if err != nil {
				return err
			}
			mime := http.DetectContentType(mbuf[:mn])
			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				return err
			}
			task := &cr.UploadTask{
				Site:    site,
				File:    file,
				Size:    size,
				Name:    name,
				Mime:    mime,
				ModTime: info.ModTime().UnixMilli(),
			}
			err = task.In(args[1])
			if err != nil {
				return err
			}
			pterm.DefaultBasicText.Printfln(
				"创建上传会话: %v, 分片: %v, 大小: %v",
				task.Sess.SessionID,
				len(task.Chunks),
				cr.SizeToString(task.Size),
			)
			multi := pterm.DefaultMultiPrinter //.WithUpdateDelay(time.Millisecond * 100)
			//.WithBarCharacter(`▬`).WithLastCharacter(`▬`).WithBarFiller(pterm.Gray(`▬`))
			chunks := make([]*pterm.ProgressbarPrinter, len(task.Chunks))
			for c, s := range task.Chunks {
				pb, _ := pterm.DefaultProgressbar.WithWriter(multi.NewWriter()).WithTotal(int(s)).Start(strconv.Itoa(c))
				chunks[c] = pb
			}
			multi.Start()
			task.Callback = func(c, n int, e error) {
				if e != nil && e != io.EOF {
					chunks[c].Stop()
					pterm.Error.Printfln(`Chunk%v: %v`, c, e)
				} else {
					chunks[c].Add(n)
				}
			}
			err = task.Go()
			/*for _, v := range chunks {
				v.Stop()
			}*/
			multi.Stop()
			if err != nil {
				println(`Error!`)
			} else {
				println(`Success!`)
			}
			return err
		}()

		file.Close()
		return err
	},
	// 容量配额
	`df`: func(args ...string) error {
		us, err := site.UserStorage()
		if err != nil {
			return err
		}
		_, err = fmt.Printf(
			"Used: %v / Free: %v / Total: %v\n",
			cr.SizeToString(us.Used),
			cr.SizeToString(us.Free),
			cr.SizeToString(us.Total),
		)
		return err
	},
	// 删除文件
	`rm`: func(args ...string) error {
		length := len(args)
		if length == 0 {
			println(`usage:`, `rm <type> <id...>`)
			return nil
		}
		if length < 2 {
			return errors.New(`参数不足`)
		}
		err := site.ObjectDel(cr.GenerateSrc(args[0] == `dir`, args[1:]...))
		if err != nil {
			return err
		}
		println(`Success!`)
		return nil
	},
	// 退出登录
	`logout`: func(args ...string) error {
		return site.UserSessionDel()
	},
}

func loop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// 构建当前路径
		var b strings.Builder
		if site.Users != nil {
			b.WriteString(site.Users.Mail)
		} else {
			b.WriteString(`cr-go-sdk`)
			b.WriteByte('@')
			b.WriteString(cr.Cr_Version)
		}
		b.WriteByte(':')
		b.WriteString(path)
		if site.Config != nil {
			if site.Config.User.Group.ID == 1 {
				b.WriteByte('#')
			} else {
				b.WriteByte('$')
			}
		} else {
			b.WriteByte('>')
		}
		b.WriteByte(' ')
		print(b.String())
		// 读取输入指令
		line, _, err := reader.ReadLine()
		if err != nil {
			println(`err:`, err)
			continue
		}
		if len(line) == 0 {
			continue
		}
		sep := strings.Split(string(line), ` `)
		fmt.Printf("sep: %s\n\n", sep)
		if len(sep) == 0 {
			println(`err:`, `参数不全`)
			continue
		}
		// 解析指令结果
		cmd, ok := cmds[sep[0]]
		if !ok {
			println(`err:`, `未定义的命令`)
			continue
		}
		if err := cmd(sep[1:]...); err != nil {
			println(`uncaught:`, err.Error())
		}
	}
}

func main() {
	if site.Addr != `` {
		if err := site.SdkInit(); err != nil {
			println(err.Error())
			os.Exit(0)
		}
		if site.Users != nil {
			if site.Users.Cookie == nil {
				if err := site.SdkLogin(); err != nil {
					println(err.Error())
				}
			}
			if site.Config.User.Anonymous {
				println(`当前正在以游客身份登录`)
			}
		}
	}
	// 列出指令
	cmds[`cmds`] = func(args ...string) error {
		for k := range cmds {
			println(k)
		}
		return nil
	}
	loop()
}
