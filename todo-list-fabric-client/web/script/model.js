class Model {
  
  constructor() {
    this.listeners = [];

    this.account = null;
    this.locations = null;
    this.peers = null;
    this.tasks = null;
  }

  addEventListener( label, callback ) {
    this.listeners.push( {
      label: label,
      callback: callback
    } );
  }

  emit( label, evt ) {
    for( let h = 0; h < this.listeners.length; h++ ) {
      if( this.listeners[h].label == label ) {
        this.listeners[h].callback( evt );
      }
    }
  }

  login( username, password ) {
    Blockchain.request( {
      method: Blockchain.QUERY,
      operation: 'account_read',
      values: [username, password]
    } ).then( result => {
      if( result == null ) {
        console.log( 'Not logged in.' );
        throw( new Error( 'Not logged in.' ) );
      } else {
        console.log( result );
        this.account = result;
        this.emit( Model.LOGIN, null );
        return Blockchain.request( {
          method: Blockchain.QUERY,
          operation: 'location_browse',
          values: [this.account.id]
        } );
      }
    } ).then( result => {      
      let any = {
        id: 'any',
        account: null,
        name: 'Any'
      };

      if( result == null ) {
        console.log( 'No locations.' );
        this.locations = [any];
      } else{
        console.log( result );        
        this.locations = result.slice( 0 );
        this.locationSort();
      }

      return Blockchain.request( {
        method: Blockchain.QUERY,
        operation: 'task_browse',
        values: [this.account.id]
      } );      
    } ).then( result => {
      if( result == null ) {
        console.log( 'No tasks.' );
        this.tasks = [];
      } else {
        console.log( result );          
        this.tasks = result.slice( 0 );
        this.taskSort();
      }

      return Blockchain.request( {
        method: Blockchain.QUERY,
        operation: 'account_browse',
        values: [this.account.id]
      } );  
    } ).then( result => {
      for( let r = 0; r < result.length; r++ ) {
        result[r].name = result[r].first + ' ' + result[r].last;
      }

      this.peers = result.slice( 0 );
      this.peerSort();
      this.emit( Model.READY, null );      
    } ).catch( error => {
      this.emit( Model.WRONG, null );        
    } );
  }

  locationAdd( location ) {
    Blockchain.request( {
      method: Blockchain.INVOKE,
      operation: 'location_add',
      values: [
        location.id, 
        location.account, 
        location.name
      ]
    } ).then( result => {
      console.log( result );
      this.emit( Model.LOCATION_ADD, null );       
    } );    
  }

  locationDelete( location ) {
    for( let p = 0; p < this.locations.length; p++ ) {
      if( this.locations[p].id == location.id ) {
        this.locations.splice( p, 1 );
        break;
      }
    }

    Blockchain.request( {
      method: Blockchain.INVOKE,
      operation: 'location_delete',
      values: [location.id], 
    } ).then( result => {
      console.log( result );     
    } );    
  }  

  locationSort() {
    if( this.locations[0].id == 'any' ) {
      this.locations.splice( 0, 1 );      
    }

    // Sort locations by name
    this.locations = this.locations.sort( ( a, b ) => {
      if( a.name.toUpperCase() < b.name.toUpperCase() ) {
        return -1;
      }

      if( a.name.toUpperCase() > b.name.toUpperCase() ) {
        return 1;
      }

      return 0;
    } );        

    this.locations.unshift( {
      id: 'any',
      account: null,
      name: 'Any'
    } );    
  }

  peerSort() {
    // Sort peer accounts by last name
    this.peers = this.peers.sort( ( a, b ) => {
      if( a.last.toUpperCase() < b.last.toUpperCase() ) {
        return -1;
      }

      if( a.last.toUpperCase() > b.last.toUpperCase() ) {
        return 1;
      }

      return 0;      
    } );
  }

  reset( accounts, locations, tasks ) {
    Blockchain.request( {
      method: Blockchain.INVOKE,
      operation: 'reset_data',
      values: [
        accounts,
        locations,
        tasks
      ]
    } ).then( result => {
      console.log( result );
    } );        
  }

  taskAdd( task ) {
    Blockchain.request( {
      method: Blockchain.INVOKE,
      operation: 'task_add',
      values: [
        task.id, 
        task.account, 
        task.name, 
        task.due.toString(),
        task.location, 
        task.duration.toString(), 
        task.energy.toString(), 
        task.created.toString()
      ]
    } ).then( result => {
      console.log( result );
    } );    
  }

  taskAssign( task ) {
    for( let t = 0; t < this.tasks.length; t++ ) {
      if( this.tasks[t].id == task.id ) {
        this.tasks.splice( t, 1 );
        break;
      }
    }

    this.taskEdit( task );    
  }

  taskDelete( task ) {
    for( let t = 0; t < this.tasks.length; t++ ) {
      if( this.tasks[t].id == task.id ) {
        this.tasks.splice( t, 1 );
        break;
      }
    }

    Blockchain.request( {
      method: Blockchain.INVOKE,
      operation: 'task_delete',
      values: [task.id]
    } ).then( result => {
      console.log( result );
    } );  
  }

  taskEdit( task ) {
    Blockchain.request( {
      method: Blockchain.INVOKE,
      operation: 'task_edit',
      values: [
        task.id, 
        task.account, 
        task.due.toString(), 
        task.location, 
        task.duration.toString(), 
        task.energy.toString(), 
        task.tags, 
        task.notes, 
        task.complete.toString(), 
        task.name]
    } ).then( result => {
      console.log( result );
    } );
  }

  taskSort() {
    // Sort tasks by date
    this.tasks = this.tasks.sort( ( a, b ) => {
      if( a.due < b.due ) {
        return -1;
      }

      if( a.due > b.due ) {
        return 1;
      }

      return 0;
    } );

    // Sort tasks by complete
    this.tasks = this.tasks.sort( ( a, b ) => {
      if( a.complete ) {
        return 1;
      }

      return -1;
    } );
  }

}

Model.DURATION = [
  {id: '0', name: 'Any'},
  {id: '1', name: '&lt; 30 mins.'},
  {id: '2', name: '30 - 60 mins.'},
  {id: '3', name: '1 - 2 hrs.'},
  {id: '4', name: '2 - 4 hrs.'},
  {id: '5', name: '&gt 4 hrs.'}  
];

Model.ENERGY = [
  {id: '0', name: 'Any'},
  {id: '1', name: 'Low'},
  {id: '2', name: 'Normal'},
  {id: '3', name: 'High'},      
];

Model.LOCATION_ADD = 'model_location_add';

Model.LOGIN = 'model_login';
Model.READY = 'model_ready';
Model.WRONG = 'model_wrong';
