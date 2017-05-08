//-------------------------------------------------------------------
// Toodles Chaincode Library
//-------------------------------------------------------------------
var util = require("util");
module.exports = function (enrollObj, g_options, fcw, logger) {
	var toodles_chaincode = {};

	// Chaincode -------------------------------------------------------------------------------
  //Invoke a function
	toodles_chaincode.invoke_function = function (options, cb) {
		console.log('');
		logger.info('Invoking function '+options.func +"  with arguments "+options.args);

		var opts = {
			channel_id: g_options.channel_id,
			chaincode_id: g_options.chaincode_id,
			chaincode_version: g_options.chaincode_version,
			event_url: g_options.event_url,
			peer_tls_opts: g_options.peer_tls_opts,
			cc_function: options.func,
			cc_args: options.args
		};
		fcw.invoke_chaincode(enrollObj, opts, cb);
	};

	toodles_chaincode.query_function = function (options, cb) {
		console.log('');
		logger.info('Fetching EVERYTHING...');

		var opts = {
			channel_id: g_options.channel_id,
			chaincode_version: g_options.chaincode_version,
			chaincode_id: g_options.chaincode_id,
			cc_function: options.func,
			cc_args: options.args
		};
		fcw.query_chaincode(enrollObj, opts, cb);
	};

	return toodles_chaincode;
};
