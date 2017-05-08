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
  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


// ============================================================================================================================
// Browse location
// Inputs - Account ID
// ============================================================================================================================

func ( t *SimpleChaincode ) location_browse( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil {
    return shim.Error( "Unable to get locations." )
  }

  // Back door for viewing all data
  if( args[0] == "all" ) {
    return shim.Success(bytes)
  }

  var locations []Location
  var items []Location

  // From JSON to data structure
  err = json.Unmarshal( bytes, &locations )

  // Look for match
  for _, location := range locations {
    // Match
    if location.AccountId == args[0] {
      items = append( items, location )
    }
  }

  // JSON encode
  bytes, err = json.Marshal( items )

  return shim.Success(bytes)
}

// ============================================================================================================================
// Read location
// Inputs -Id
// ============================================================================================================================

func ( t *SimpleChaincode ) location_read( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil {
    return shim.Error( "Unable to get locations." )
  }

  var locations []Location

  // From JSON to data structure
  err = json.Unmarshal( bytes, &locations )
  found := false

  // Look for match
  for _, location := range locations {
    // Match
    if location.Id == args[0] {
      // JSON encode
      bytes, err = json.Marshal( location )
      found = true
      break
    }
  }

  // Nope
  if found != true {
    bytes, err = json.Marshal( nil )
  }

  return shim.Success(bytes)
}


// ============================================================================================================================
// Edit location
// Inputs - ID, Account ID, Name
// ============================================================================================================================

func ( t *SimpleChaincode ) location_edit( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil {
    return shim.Error( "Unable to get locations." )
  }

  var locations []Location

  // From JSON to data structure
  err = json.Unmarshal( bytes, &locations )

  // Look for match
  for g := 0; g < len( locations ); g++ {
    // Match
    if locations[g].Id == args[0] {
      locations[g].AccountId = args[1]
      locations[g].Name = args[2]
      break
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( locations )
  err = stub.PutState( "toodles_locations", bytes )

  return shim.Success(nil)
}

// ============================================================================================================================
// Add location
// Inputs - ID, Account ID, Name
// ============================================================================================================================

func ( t *SimpleChaincode ) location_add( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil {
    return shim.Error( "Unable to get locations." )
  }

  var location Location

  // Build JSON values
  id := "\"id\": \"" + args[0] + "\", "
  account := "\"account\": \"" + args[1] + "\", "
  name := "\"name\": \"" + args[2] + "\""

  // Make into a complete JSON string
  // Decode into structure instance
  content := "{" + id + account + name + "}"
  err = json.Unmarshal( []byte(content), &location )

  var locations []Location

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &locations )
  locations = append( locations, location )

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( locations )
  err = stub.PutState( "toodles_locations", bytes )

  return shim.Success(nil)
}


// ============================================================================================================================
// Delete location
// Inputs - ID
// ============================================================================================================================

func ( t *SimpleChaincode ) location_delete( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_locations" )

  if err != nil {
    return shim.Error( "Unable to get locations." )
  }

  var locations []Location

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &locations )

  for g := 0; g < len( locations ); g++ {
    // Match
    if locations[g].Id == args[0] {
      locations = append( locations[:g], locations[g + 1:]... )
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( locations )
  err = stub.PutState( "toodles_locations", bytes )

  return shim.Success(nil)
}
