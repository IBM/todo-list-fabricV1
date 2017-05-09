/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package main

import (
  "encoding/json"
  "fmt"
  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)

// Task duration
const DURATION_ANY  = 0
const DURATION_30   = 1
const DURATION_60   = 2
const DURATION_2    = 3
const DURATION_4    = 4
const DURATION_LONG = 5

// Task energy level
const ENERGY_ANY    = 0
const ENERGY_LOW    = 1
const ENERGY_NORMAL = 2
const ENERGY_HIGH   = 3

// SimpleChaincode example simple Chaincode implementation
type  SimpleChaincode struct {
}

// ============================================================================================================================
// Asset Definitions - The ledger will store account, location and task
// ============================================================================================================================

// Account
type Account struct {
  Id       string `json:"id"`
  First    string `json:"first"`
  Last     string `json:"last"`
  Name     string `json:"name"`
  Password string `json:"password"`
}

// Location
type Location struct {
  Id        string `json:"id"`
  AccountId string `json:"account"`
  Name      string `json:"name"`
}

// To do task
type Task struct {
  Id         string `json:"id"`
  AccountId  string `json:"account"`
  Due        int    `json:"due"`
  LocationId string `json:"location"`
  Duration   int    `json:"duration"`
  Energy     int    `json:"energy"`
  Tags       string `json:"tags"`
  Notes      string `json:"notes"`
  Complete   bool   `json:"complete"`
  Name       string `json:"name"`
  CreatedAt  int    `json:"created"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================

// Main
func main() {
  err := shim.Start( new( SimpleChaincode ) )
  if err != nil {
    fmt.Printf( "Error starting chaincode: %s", err )
  }
}


// ============================================================================================================================
// Init - initialize the chaincode
// ============================================================================================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
  // Accounts
  fmt.Println("Toodles Is Starting Up")
  var accounts []Account
  bytes, err := json.Marshal( accounts )

  if err != nil {
    return shim.Error("Error initializing accounts.")
  }

  err = stub.PutState( "toodles_accounts", bytes )

  // Locations
  // TODO: Empty array versus nil
  var locations []string

  bytes, err = json.Marshal( locations )

  if err != nil {
    return shim.Error("Error initializing locations.")
  }

  err = stub.PutState( "toodles_locations", bytes )

  // Tasks
  var tasks []string

  bytes, err = json.Marshal( tasks )

  if err != nil {
    return shim.Error("Error initializing tasks.")
  }

  err = stub.PutState( "toodles_tasks", bytes )

  return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
  fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)
  if function == "init" {
    return t.Init(stub)
  } else if function == "account_add" {
    return t.account_add( stub, args )
  } else if function == "account_edit" {
    return t.account_edit( stub, args )
  } else if function == "account_delete" {
    return t.account_delete( stub, args )
  } else if function == "task_add" {
    return t.task_add( stub, args )
  } else if function == "task_edit" {
    return t.task_edit( stub, args )
  } else if function == "task_delete" {
    return t.task_delete( stub, args )
  } else if function == "location_add" {
    return t.location_add( stub, args )
  } else if function == "location_edit" {
    return t.location_edit( stub, args )
  } else if function == "location_delete" {
    return t.location_delete( stub, args )
  } else if function == "reset_data" {
    return t.reset_data( stub, args )
  } else if function == "account_browse" {
    return t.account_browse( stub, args )
  } else if function == "account_read" {
    return t.account_read( stub, args )
  } else if function == "task_browse" {
    return t.task_browse( stub, args )
  } else if function == "task_read" {
    return t.task_read( stub, args )
  } else if function == "location_browse" {
    return t.location_browse( stub, args )
  } else if function == "location_read" {
    return t.location_read( stub, args )
  }

  return shim.Error( "Function with the name " + function + " does not exist.")
}
