class Blockchain {

  static request( route ) {
    return new Promise( ( resolve, reject ) => {
      let xhr = new XMLHttpRequest();
      xhr.addEventListener( 'error', evt => {
        reject( xhr.statusText );
      } );
      xhr.addEventListener( 'load', evt => {
        console.log(xhr);
        let data = JSON.parse( xhr.responseText );
        let result = null;
        console.log(data);
        try {
          result = JSON.parse( data.result.message );
        } catch( e ) {
          result = data.result.message;
        }

        resolve( result );
      } );
      xhr.open( 'POST', Blockchain.URL, true );
      xhr.setRequestHeader( 'Content-Type', 'application/json' );
      xhr.send( JSON.stringify( {
        jsonrpc: '2.0',
        method: route.method,
        params: {
          chaincodeID: {
            name: Blockchain.CHAINCODE
          },
          ctorMsg: {
            function: route.operation,
            args: route.values
          },
          secureContext: Blockchain.USER,
          type: 1
        },
        id: 1
      } ) );
    } );
  }

}

Blockchain.CHAINCODE = 'end2end';
Blockchain.URL = 'http://localhost:3000/chaincode';
Blockchain.USER = 'user_type1_0';

Blockchain.QUERY = 'query';
Blockchain.INVOKE = 'invoke';
