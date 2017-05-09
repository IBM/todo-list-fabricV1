'use strict';
var express = require('express')
var bodyParser = require('body-parser');
var app = express()
var path = require('path');
var app = express();
var cors = require('cors');
var fs = require('fs');
var os = require('os');
var winston = require('winston');								//logginer module
var util = require("util");
// --- Set Our Things --- //
var logger = new (winston.Logger)({
	level: 'debug',
	transports: [
		new (winston.transports.Console)({ colorize: true }),
	]
});
var helper = require(__dirname + '/utils/helper.js')(process.env.creds_filename, logger);
var fcw = require('./utils/fc_wrangler/index.js')({ block_delay: helper.getBlockDelay() }, logger);
var enrollObj = null;
var toodles_lib = null;
var port = helper.getPort();
var host = 'localhost';
//enroll an admin with the CA for this peer/channel
function enroll_admin(attempt, cb) {
	fcw.enroll(helper.makeEnrollmentOptions(0), function (errCode, obj) {
		if (errCode != null) {
			logger.error('could not enroll...');

			// --- Try Again ---  //
			if (attempt >= 2) {
				if (cb) cb(errCode);
			} else {
				try {
					logger.warn('removing older kvs and trying to enroll again');
					rmdir(makeKVSpath());				//delete old kvs folder
					logger.warn('removed older kvs');
					enroll_admin(++attempt, cb);
				} catch (e) {
					logger.error('could not delete old kvs', e);
				}
			}
		} else {
			enrollObj = obj;
			if (cb) cb(null);
		}
	});
}


// remove any kvs from last run
function rmdir(dir_path) {
	if (fs.existsSync(dir_path)) {
		fs.readdirSync(dir_path).forEach(function (entry) {
			var entry_path = path.join(dir_path, entry);
			if (fs.lstatSync(entry_path).isDirectory()) {
				rmdir(entry_path);
			}
			else {
				fs.unlinkSync(entry_path);
			}
		});
		fs.rmdirSync(dir_path);
	}
}



// make the path to the kvs we use
function makeKVSpath() {
	var temp = helper.makeEnrollmentOptions(0);
	return path.join(os.homedir(), '.hfc-key-store/', temp.uuid);
}
app.options('*', cors());
app.use(cors());
app.use(bodyParser.json()); // support json encoded bodies
app.use(bodyParser.urlencoded({ extended: true })); // support encoded bodies

app.get('/enrollAdmin', function(req, res) {

	enroll_admin(1, function (e) {
		if (e == null) {
			res.json({message: "Admin Enrolled! "})
		}else{
			res.json({message: "Error in enrolling Admin! "+e})
		}
	});

})


//{"jsonrpc":"2.0","result":{"status":"OK","message":"bb7d1ade-c223-4775-ad4e-dc1556a3a2db"},"id":1}

app.post('/chaincode', function(req, res) {

  var data = typeof req.body !== "string"?req.body:JSON.parse(req.body);
  var values = typeof data.params !== "string"?data.params:JSON.parse(data.params);;


	var opts = helper.makeLibOptions();
	toodles_lib = require('./utils/toodles_cc_lib.js')(enrollObj, opts, fcw, logger);
	if(values.ctorMsg.args.length == 1 && values.ctorMsg.args[0]==null){
		values.ctorMsg.args =[""]
	}

	var options = {
		func: values.ctorMsg.function,
		args: values.ctorMsg.args
	};
	logger.debug("Query parameters : "+options);
	if(data.method.includes('invoke')){
	toodles_lib.invoke_function(options,function (err, resp) {
		logger.debug("Results from invoke_function invoke_cc: ");
		logger.debug(util.inspect(resp, {showHidden: false, depth: null}));
		logger.debug(util.inspect(err, {showHidden: false, depth: null}));
			if (err != null) {
				//console.log(res);
				res.json({"result":{"status":"failed","message":err}})
			}else{
				//res.send(resp);
				res.json({"result":{"status":"OK","message":resp}});
			}
	});
	}
	else{
		toodles_lib.query_function(options,function (err, resp) {
			logger.debug("Results from query_function invoke_cc: ");
			logger.debug(util.inspect(resp, {showHidden: false, depth: null}));
			logger.debug(util.inspect(err, {showHidden: false, depth: null}));
			if (err != null) {
				//console.log(res);
				res.json({"result":{"status":"failed","message":err}})
			}else{
					res.json({"result":{"status":"OK","message":resp.raw_peer_payloads}});
			}
		});
	}
	//res.json({message: "success"})
})

app.listen(port)
console.log('------------------------------------------ Server Up - ' + host + ':' + port + ' ------------------------------------------');
