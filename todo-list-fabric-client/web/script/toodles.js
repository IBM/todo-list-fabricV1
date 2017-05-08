class Toodles {

  constructor() {
    model.addEventListener( Model.WRONG, evt => this.doLoginError( evt ) );    
    model.addEventListener( Model.LOGIN, evt => this.doLoginSuccess( evt ) );    
    model.addEventListener( Model.READY, evt => this.doModelReady( evt ) );
    model.addEventListener( Model.LOCATION_ADD, evt => this.doMessageCancel( evt ) );

    this.authentication = document.querySelector( 'toodles-authentication' );
    this.authentication.addEventListener( Authentication.SIGN_IN, evt => this.doSignIn( evt ) );
    this.authentication.addEventListener( Authentication.RESET, evt => this.doReset( evt ) );    

    this.toolbar = document.querySelector( 'toodles-toolbar' );
    this.toolbar.addEventListener( Toolbar.EXIT, evt => this.doExit( evt ) );

    this.list = document.querySelector( 'toodles-list' );

    this.filter = document.querySelector( 'toodles-filter' );
    this.filter.addEventListener( Filter.CHANGE, evt => this.doFilter( evt ) );
    this.filter.addEventListener( Location.DELETE, evt => this.doLocationDelete( evt ) );    
    
    this.action = document.querySelector( 'toodles-action' );
    this.action.addEventListener( Action.CREATE_TASK, evt => this.doTaskCreate( evt ) );
    this.action.addEventListener( Action.CREATE_LOCATION, evt => this.doLocationCreate( evt ) );    

    this.doMessageAccept = this.doMessageAccept.bind( this );
    this.doMessageCancel = this.doMessageCancel.bind( this );
  }

  doFilter( evt ) {
    this.list.filter( evt.detail );
  }

  doLocationCreate( evt ) {
    let message = document.querySelector( 'toodles-message' );

    if( message ) {
      message.remove();
    }

    message = document.createElement( 'toodles-message' );
    message.addEventListener( Message.CANCEL, this.doMessageCancel );
    message.addEventListener( Message.OK, this.doMessageAccept );    
    document.body.appendChild( message );
    message.show();
  }

  doLocationDelete( evt ) {
    model.locationDelete( evt.detail );
    this.filter.places = model.locations;
  }

  doLoginError( evt ) {
    this.authentication.loading = false;
    this.authentication.shake();
  }

  doLoginSuccess( evt ) {
    this.authentication.id = model.account.id;
    this.authentication.hide();
    this.toolbar.show();
  }

  doMessageAccept( evt ) {
    let location = {
      id: uuid.v4(),
      account: model.account.id,
      name: evt.detail.value
    };

    model.locationAdd( location );
    model.locations.push( location );
    model.locationSort();

    this.filter.places = model.locations;
  }

  doMessageCancel( evt ) {
    let message = document.querySelector( 'toodles-message' );
    message.removeEventListener( Message.CANCEL, this.doMessageCancel );
    message.removeEventListener( Message.OK, this.doMessageAccept );
    message.remove();
  }

  doModelReady( evt ) {
    this.filter.places = model.locations;
    this.filter.show();

    this.list.data = model.tasks;
    this.list.show();

    this.action.show();
  }

  doExit( evt ) {
    this.action.hide();

    // Clear list
    this.list.hide();
    this.list.data = [];

    // Reset filtering
    this.filter.hide();
    this.filter.location = {
      id: 'any'
    };
    this.filter.duration = 0;
    this.filter.energy = 0;

    this.toolbar.hide();

    this.authentication.show();
  }

  doReset( evt ) {
    let xhr = new XMLHttpRequest();
    xhr.addEventListener( 'load', evt => { 
      let data = JSON.parse( xhr.responseText );
      model.reset( 
        JSON.stringify( data.accounts ),
        JSON.stringify( data.locations ),
        JSON.stringify( data.tasks )
      );
    } );
    xhr.open( 'GET', '/data/reset.json', true );
    xhr.send( null );
  }

  doSignIn( evt ) {
    model.login( evt.detail.username, evt.detail.password );
  }

  doTaskCreate( evt ) {
    let task = {
      id: uuid.v4(),
      account: model.account.id,
      complete: false,
      due: 0,
      duration: 0,
      energy: 0,
      location: 'any',
      name: 'New To Do',
      notes: '',
      tags: '',
      created: new Date().getTime()
    };

    model.taskAdd( task );
    this.list.add( task );
  }

  // https://www.kirupa.com/html5/get_element_position_using_javascript.htm
  static getPosition( element ) {
    let xPos = 0;
    let yPos = 0;

    while( element ) {
      if( element.tagName == 'BODY' ) {
        let xScroll = element.scrollLeft || document.documentElement.scrollLeft;
        let yScroll = element.scrollTop || document.documentElement.scrollTop;

        xPos += ( element.offsetLeft - xScroll + element.clientLeft);
        yPos += ( element.offsetTop - yScroll + element.clientTop );
      } else {
        xPos += ( element.offsetLeft - element.scrollLeft + element.clientLeft );
        yPos += ( element.offsetTop - element.scrollTop + element.clientTop );
      }

      element = element.offsetParent;
    }

    return {
      x: xPos,
      y: yPos
    };
  }

}

let model = new Model();
let app = new Toodles();
