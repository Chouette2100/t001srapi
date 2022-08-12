/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

Ver. 0.0.0

*/
package main

import (
	"log"
	"os"
	"sort"
	"time"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
)

type Config struct {
	SR_acct    string //	SHOWROOMのアカウント名
	SR_pswd    string //	SHOWROOMのパスワード
	Category   string
	Genre   []string
}


/*
	配信中ルームの一覧を取得する

	$ cd ~/go/src/t001srapi
	$ vi t001srapi.go
	$ go mod init					<== 注意：パッケージ部分のソースをダウンロードした場合はimport部分はそのままにしておいて
	$ go mod tidy					<== 	  go.modに“replace github.com/Chouette2100/srapi ../srapi”みたいなのを追加します。
	$ go build t001srapi.go
	$ cat config.yml 
	sr_acct: ${SRACCT}				<== ログインアカウントを環境変数 SRACCT で与えます。ここに直接アカウントを書くこともできます。
	sr_pswd: ${SRPSWD}				<== ログインパスワードを環境変数 SRPSWD で与えます。ここに直接パスワードを書くこともできます。
	category: Free					<== "Free"|"Official"|"All"
	genre:
	- アイドル						<== "人気"|"フリー"|"アイドル"|"タレント・モデル"|...
	- タレント・モデル
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRACCT xxxxxxxxx
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRPSWD xxxxxxxxx
	$ ./t000srapi config.yml

	v0.0.0
*/
func main() {

	//	ログファイルを設定する。
	logfile := exsrapi.CreateLogfile("", "")
	defer logfile.Close()

	if len(os.Args) != 2 {
		//      引数が足りない(設定ファイル名がない)
		log.Printf("usage:  %s NameOfConfigFile\n", os.Args[0])
		return
	}

	//	設定ファイルを読み込む
	var config Config
	err := exsrapi.LoadConfig(os.Args[1], &config)
	if err != nil {
		log.Printf("exsrapi.LoadConfig: %s\n", err.Error())
		return
	}	

	//	cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient(config.SR_acct)
	if err != nil {
		log.Printf("CreateNewClient: %s\n", err.Error())
		return
	}
	//	すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	//	配信しているルームの一覧を取得する
	roomlives, err := srapi.ApiLiveOnlives(client)
	if err != nil {
		log.Printf("ApiLiveOnlives(): %s\n", err.Error())
		return
	}
	log.Printf("*****************************************************************\n")
	log.Printf("配信中ルーム数\n")
	log.Printf("\n")
	log.Printf("　ジャンル数= %d\n", len(roomlives.Onlives))
	log.Printf("\n")
	log.Printf("　ルーム数　ジャンル　ジャンル名\n")
	for _, roomlive := range roomlives.Onlives {
		log.Printf("%10d%10d  %s\n", len(roomlive.Lives) , roomlive.Genre_id, roomlive.Genre_name)
	}
	//	指定したカテゴリー（Free|Official|All）のルームの一覧を取得する
	roomlive, err := roomlives.ExtrByCtg(config.Category)
	if err != nil {
		log.Printf("ExtrRoomLiveByCtg: %s\n", err.Error())
		return
	}
	log.Printf("\n")
	log.Printf("カテゴリ[%s]のルーム数 = %d\n", config.Category, len(*roomlive))
	log.Printf("\n")

	log.Printf("*****************************************************************\n")
	log.Printf("ジャンル別配信中ルーム数\n")
	//	指定したジャンル（アイドル|タレント・モデル|...）のルームの一覧を取得する
	gnrmap := map[string]bool {}
	for _, gnre := range config.Genre {
		gnrmap[gnre] = true
	}
	log.Printf("　指定ジャンル = %+v\n", config.Genre)
	//	log.Printf("Gnrmap= %+v\n", gnrmap)
	roomlive, err =  roomlives.ExtrByGnr(gnrmap)
	if err != nil {
		log.Printf("ExtrRoomLiveByGnr: %s\n", err.Error())
		return
	}
	log.Printf("\n")
	log.Printf("　指定したジャンルのルーム数 = %d\n", len(*roomlive))

	log.Printf("\n")
	log.Printf("　*** 指定したジャンルのルーム一覧\n")
	for _, room := range *roomlive {
		log.Printf("  started at %s %s\n", time.Unix(room.Started_at, 0).Format("02 15:04:05"), room.Main_name)
	}

	//	ルームを配信開始時刻の降順にソートする
	sort.Sort(*roomlive)

	log.Printf("\n")
	log.Printf("　*** 指定したジャンルのルーム一覧（ソート後）\n")
	for _, room := range *roomlive {
		log.Printf("  started at %s %s\n", time.Unix(room.Started_at, 0).Format("02 15:04:05"), room.Main_name)
	}
}
