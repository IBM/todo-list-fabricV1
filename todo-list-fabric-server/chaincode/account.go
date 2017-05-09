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


// ============================================================================================================================
// Browse accounts - Used to cross-assign tasks
// ============================================================================================================================

func ( t *SimpleChaincode ) account_browse( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_accounts" )
  if err != nil {
    fmt.Printf("Query Response failed:%s\n", bytes)
    return shim.Error( "Failed to get accounts " + err.Error())
  }

  var accounts []Account
  var peers []Account

  // From JSON to data structure
  err = json.Unmarshal( bytes, &accounts )

  if args[0] == "all" {
    bytes, err = json.Marshal( accounts )
    fmt.Printf("Query Response:%s\n ", bytes)
    fmt.Printf("Query Response accounts:%s\n", accounts)
    return shim.Success(bytes)
//   return bytes, nil
  }

  // Scrub passwords
  for a := 0; a < len( accounts ); a++ {
    if accounts[a].Id != args[0] {
      accounts[a].Password = ""
      peers = append( peers, accounts[a] )
    }
  }

  bytes, err = json.Marshal( peers )
  fmt.Printf("Query Response:%s\n", bytes)
  fmt.Printf("Query Response accounts:%s\n", accounts)
  return shim.Success(bytes)
//  return bytes, nil
}


// ============================================================================================================================
// Read Account - Used to read the account information
// Inputs - Name, Password
// ============================================================================================================================

func ( t *SimpleChaincode ) account_read( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_accounts" )
  if err != nil {
    return shim.Error( "Unable to get accounts." +err.Error())
  }

  var accounts []Account

  // From JSON to data structure
  err = json.Unmarshal( bytes, &accounts )
  found := false

  // Look for match
  for _, account := range accounts {
    // Match
    if account.Name == args[0] && account.Password == args[1] {
      // Sanitize
      account.Password = ""

      // JSON encode
      bytes, err = json.Marshal( account )
      found = true
      break
    }
  }

  // Nope
  if found != true {
    bytes, err = json.Marshal( nil )
  }
  //fmt.Printf("Query Response:%s\n", bytes)
  return shim.Success(bytes)
}


// ============================================================================================================================
// Edit Account - Used to edit the account information
// Inputs - ID, First, Last, Name, Password
// ============================================================================================================================

func ( t *SimpleChaincode ) account_edit( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil {
     return shim.Error( "Unable to get accounts." )
  }

  var accounts []Account

  // From JSON to data structure
  err = json.Unmarshal( bytes, &accounts )

  // Look for match
  for a := 0; a < len( accounts ); a++ {
    // Match
    if accounts[a].Id == args[0] {
      accounts[a].First = args[1]
      accounts[a].Last = args[2]
      accounts[a].Name = args[3]
      accounts[a].Password = args[4]
      break
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( accounts )
  err = stub.PutState( "toodles_accounts", bytes )
  //fmt.Printf("Query Response:%s\n", bytes)
  return shim.Success(nil)
}


// ============================================================================================================================
// Add Account - Used to add new account
// Inputs - ID, First, Last, Name, Password
// ============================================================================================================================

func ( t *SimpleChaincode ) account_add( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil {
    return shim.Error( "Unable to get accounts." )
  }

  var account Account

  // Build JSON values
  id := "\"id\": \"" + args[0] + "\", "
  first := "\"first\": \"" + args[1] + "\", "
  last := "\"last\": \"" + args[2] + "\", "
  name := "\"name\": \"" + args[3] + "\", "
  password := "\"password\": \"" + args[4] + "\""

  // Make into a complete JSON string
  // Decode into a single account value
  content := "{" + id + first + last + name + password + "}"
  err = json.Unmarshal( []byte( content ), &account )
  fmt.Printf("Query Response  content :\n", content)
  var accounts []Account

  // Decode JSON into account array
  // Add latest account
  err = json.Unmarshal( bytes, &accounts )
  accounts = append( accounts, account )

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( accounts )
  err = stub.PutState( "toodles_accounts", bytes )
  return shim.Success(nil)
}


// ============================================================================================================================
// Delete Account - Used to remove the account
// Inputs - ID
// ============================================================================================================================

func ( t *SimpleChaincode ) account_delete( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_accounts" )

  if err != nil {
    return shim.Error( "Unable to get accounts." )
  }

  var accounts []Account

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &accounts )

  for a := 0; a < len( accounts ); a++ {
    // Match
    if accounts[a].Id == args[0] {
      accounts = append( accounts[:a], accounts[a + 1:]... )
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( accounts )
  err = stub.PutState( "toodles_accounts", bytes )
  //fmt.Printf("Query Response:%s\n", bytes)
  return shim.Success(nil)
}


// ============================================================================================================================
// Reset data -  Load the initial the data
// Inputs - List of accounts,locations and tasks
// ============================================================================================================================

func ( t *SimpleChaincode ) reset_data( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  stub.PutState( "toodles_accounts", []byte( args[0] ) )
  stub.PutState( "toodles_locations", []byte( args[1] ) )
  stub.PutState( "toodles_tasks", []byte( args[2] ) )

  return shim.Success(nil)
}
