/*
 * Copyright (c) 2016 General Electric Company. All rights reserved.
 *
 * The copyright to the computer software herein is the property of
 * General Electric Company. The software may be used and/or copied only
 * with the written permission of General Electric Company or in accordance
 * with the terms and conditions stipulated in the agreement/contract
 * under which the software has been supplied.
 *
 * author: chia.chang@ge.com
 */

package main

import (
	//"os"
	"flag"

	"net"
	"net/url"
	util "github.build.ge.com/212359746/wzutil"
	plugin "github.build.ge.com/212359746/wzplugin"
	"gopkg.in/yaml.v2"
	"encoding/base64"

)

var ()

const (
	//YML_TLS_FLAG = "tls"
	REV = "v1"

)

func GetTLSSetting()(map[string]interface{}, error){
	
	plg:=flag.String("plg","","Enable support for EC TLS Plugin.")
	flag.Parse()
	util.DbgLog(*plg)
	
	f, err := base64.StdEncoding.DecodeString(*plg)
	if err!=nil{
		return nil, err
	}
	t:=make(map[string]interface{})
	
	err=yaml.Unmarshal(f, &t)
	if err!=nil{
		return nil,err
	}
	
	return t,nil

}

func init(){
	util.Init("tls",true)
}

func main(){

	defer func(){
		if r:=recover();r!=nil{
			util.PanicRecovery(r)
		} else {
			util.InfoLog("plugin undeployed.")
		}
	}()

	t,err:=GetTLSSetting()
	if err!=nil{
		panic(err)
	}
	util.DbgLog(t)
	
	_, err= net.ResolveTCPAddr("tcp",t["hostname"].(string)+":"+t["tlsport"].(string))
	if err != nil {
		panic(err)
	}

	p:=plugin.NewProxy(REV)

	u, err := url.Parse(t["proxy"].(string))
	if err != nil {
		panic(err)
	}

	if err:=p.Init(t["schema"].(string)+"://"+t["hostname"].(string)+":"+t["tlsport"].(string),u);err!=nil{
		panic(err)
	}

	p.Start(t["port"].(string))
	
	
}
