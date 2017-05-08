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
  "strconv"
  pb "github.com/hyperledger/fabric/protos/peer"
)


// ============================================================================================================================
// Browse task
// Inputs -  Account ID
// ============================================================================================================================

func ( t *SimpleChaincode ) task_browse( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil {
    return shim.Error( "Unable to get tasks." )
  }

  // Back door for viewing all data
  if( args[0] == "all" ) {
    return shim.Success(bytes)
  }

  var tasks []Task
  var items []Task

  // From JSON to data structure
  err = json.Unmarshal( bytes, &tasks )

  // Look for match
  for _, task := range tasks {
    // Match
    if task.AccountId == args[0] {
      items = append( items, task )
    }
  }

  // JSON encode
  bytes, err = json.Marshal( items )

  return shim.Success(bytes)
}


// ============================================================================================================================
// Task read
// Inputs - Id
// ============================================================================================================================

func ( t *SimpleChaincode ) task_read( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil {
    return shim.Error( "Unable to get tasks." )
  }

  var tasks []Task

  // From JSON to data structure
  err = json.Unmarshal( bytes, &tasks )
  found := false

  // Look for match
  for _, task := range tasks {
    // Match
    if task.Id == args[0] {
      // JSON encode
      bytes, err = json.Marshal( task )
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
// Edit task
// Inputs - ID, Account ID, Due, Location ID, Duration, Energy, Tags, Notes, Complete, Name
// ============================================================================================================================

func ( t *SimpleChaincode ) task_edit( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil {
    return shim.Error( "Unable to get tasks." )
    //return nil, errors.New( "Unable to get tasks." )
  }

  var tasks []Task

  // From JSON to data structure
  err = json.Unmarshal( bytes, &tasks )

  // Look for match
  for t := 0; t < len( tasks ); t++ {
    // Match
    if tasks[t].Id == args[0] {
      // String arguments to integer values
      // TODO: Deal with errors
      due, _ := strconv.Atoi( args[2] )
      duration, _ := strconv.Atoi( args[4] )
      energy, _ := strconv.Atoi( args[5] )

      // No ternary operator
      complete := false
      if args[8] == "true" {
        complete = true
      }

      tasks[t].AccountId = args[1]
      tasks[t].Due = due
      tasks[t].LocationId = args[3]
      tasks[t].Duration = duration
      tasks[t].Energy =  energy
      tasks[t].Tags = args[6]
      tasks[t].Notes = args[7]
      tasks[t].Complete = complete
      tasks[t].Name = args[9]
      break
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( tasks )
  err = stub.PutState( "toodles_tasks", bytes )

  return shim.Success(nil)
}

// ============================================================================================================================
// Add task
// Inputs - ID, Account ID, Name, Due, Location, Duration, Energy, Created
// ============================================================================================================================

func ( t *SimpleChaincode ) task_add( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil {
    return shim.Error( "Unable to get tasks." )
  }

  var task Task

  // Build JSON values
  id := "\"id\": \"" + args[0] + "\", "
  account := "\"account\": \"" + args[1] + "\", "
  name := "\"name\": \"" + args[2] + "\", "
  due := "\"due\": " + args[3] + ", "
  location := "\"location\": \"" + args[4] + "\", "
  duration := "\"duration\": " + args[5] + ", "
  energy := "\"energy\": " + args[6] + ", "
  tags := "\"tags\": \"\", "
  notes := "\"notes\": \"\", "
  complete := "\"complete\": false, "
  created := "\"created\": " + args[7]

  // Make into a complete JSON string
  // Decode into structure instance
  content := "{" + id + account + name + due + location + duration + energy + tags + notes + complete + created + "}"
  err = json.Unmarshal( []byte(content), &task )

  var tasks []Task

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &tasks )
  tasks = append( tasks, task )

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( tasks )
  err = stub.PutState( "toodles_tasks", bytes )

  return shim.Success(nil)
}

// Arguments: ID
func ( t *SimpleChaincode ) task_delete( stub shim.ChaincodeStubInterface, args []string ) pb.Response {
  bytes, err := stub.GetState( "toodles_tasks" )

  if err != nil {
    return shim.Error( "Unable to get tasks." )
  //  return nil, errors.New( "Unable to get tasks." )
  }

  var tasks []Task

  // Decode JSON collection into array
  // Add latest instance value
  err = json.Unmarshal( bytes, &tasks )

  for t := 0; t < len( tasks ); t++ {
    // Match
    if tasks[t].Id == args[0] {
      tasks = append( tasks[:t], tasks[t + 1:]... )
    }
  }

  // Encode as JSON
  // Put back on the block
  bytes, err = json.Marshal( tasks )
  err = stub.PutState( "toodles_tasks", bytes )

  return shim.Success(nil)
}
