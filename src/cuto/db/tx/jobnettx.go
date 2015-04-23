// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package tx

import (
	"fmt"
	"time"

	"cuto/db"
	"cuto/log"
	"cuto/util"
)

// ジョブIDをキーに持つ
type JobMap map[string]*db.JobResult

// ジョブ実行結果を保持する。
type ResultMap struct {
	JobnetResult *db.JobNetworkResult // ジョブネットワーク情報の構造体。
	Jobresults   JobMap               // ジョブネットワーク内のジョブ状態を保存するMap。
	conn         db.IConnection       // DBコネクション
}

// ジョブネットワークの開始状態を記録する。
//
// param : jobnetName ジョブネットワーク名。
//
// param : dbname データベース名。
//
// return : ジョブ実行結果を保持する構造体ポインタ。
//
// return : error
func StartJobNetwork(jobnetName string, dbname string) (*ResultMap, error) {
	jn := db.NewJobNetworkResult(jobnetName, util.DateFormat(time.Now()), db.RUNNING)

	conn, err := db.Open(dbname)
	if err != nil {
		return nil, err
	}
	resMap := &ResultMap{jn, make(JobMap), conn}

	if err := resMap.insertJobNetwork(); err != nil {
		return nil, err
	}
	return resMap, nil
}

// ネットワーク終了時に結果情報を設定する。同時にDBコネクションも切断する。
//
// param : status ジョブネットワークのステータス。
//
// param : detail ジョブネットワークに記録する詳細メッセージ。
//
// return : error
func (r *ResultMap) EndJobNetwork(status int, detail string) error {
	if r.conn == nil {
		return fmt.Errorf("Can't access DB file.")
	}
	defer r.conn.Close()

	if r.JobnetResult == nil {
		return fmt.Errorf("Invalid Jobnetwork info.")
	}
	r.JobnetResult.EndDate = util.DateFormat(time.Now())
	r.JobnetResult.Status = status
	r.JobnetResult.Detail = detail

	if err := r.updateJobNetwork(); err != nil {
		return err
	}
	return nil
}

// DBコネクションを返す。
func (r *ResultMap) GetConnection() db.IConnection {
	return r.conn
}

// ジョブネットワークレコードをInsertする。
func (r *ResultMap) insertJobNetwork() error {
	var isCommit bool
	tx, err := r.conn.GetDbMap().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if !isCommit {
			tx.Rollback()
		}
	}()

	now := util.DateFormat(time.Now())
	r.JobnetResult.CreateDate = now
	r.JobnetResult.UpdateDate = now

	err = tx.Insert(r.JobnetResult)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("networkId[%v]", r.JobnetResult.ID))
	isCommit = true
	return nil
}

// ジョブネットワークレコードをUpdateする。
func (r *ResultMap) updateJobNetwork() error {
	var isCommit bool
	tx, err := r.conn.GetDbMap().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if !isCommit {
			tx.Rollback()
		}
	}()
	r.JobnetResult.UpdateDate = util.DateFormat(time.Now())

	if _, err = tx.Update(r.JobnetResult); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	isCommit = true
	return nil
}