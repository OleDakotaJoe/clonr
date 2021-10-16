#!/usr/bin/env node
"use strict";

var request = require('request'),
    path = require('path'),
    tar = require('tar'),
    zlib = require('zlib'),
    mkdirp = require('mkdirp'),
    fs = require('fs'),
    exec = require('child_process').exec;

var VERSION = "1.0.11"; // Mapping from Node's `process.arch` to Golang's `$GOARCH`

var ARCH_MAPPING = {
  "ia32": "386",
  "x64": "amd64",
  "arm": "arm"
}; // Mapping between Node's `process.platform` to Golang's 

var PLATFORM_MAPPING = {
  "darwin": "darwin",
  "linux": "linux",
  "win32": "windows",
  "freebsd": "freebsd"
};

function getNpmBinLocation(callback) {
  // `npm bin` will output the path where binary files should be installed
  exec("npm bin -g", function (err, stdout, stderr) {
    var dir = null;

    if (err || stderr || !stdout || stdout.length === 0) {
      // We couldn't infer path from `npm bin`. Let's try to get it from
      // Environment variables set by NPM when it runs.
      // npm_config_prefix points to NPM's installation directory where `bin` folder is available
      // Ex: /Users/foo/.nvm/versions/node/v4.3.0
      var env = process.env;

      if (env && env.npm_config_prefix) {
        dir = path.join(env.npm_config_prefix, "bin");
      }
    } else {
      dir = stdout.trim();
    }

    mkdirp.sync(dir);
    callback(null, dir);
  });
}

function verifyAndPlaceBinary(binName, binPath, callback) {
  console.log("binName: ".concat(binName, ", binPath: ").concat(binPath, ", ").concat(path.join(binPath, binName)));
  console.log(fs.existsSync("clonr.exe"));
  if (!fs.existsSync(binName)) return callback("Downloaded binary does not contain the binary specified in configuration - ".concat(binName));
  getNpmBinLocation(function (err, installationPath) {
    if (err) return callback("Error getting binary installation path from `npm bin`"); // Move the binary file

    fs.renameSync(binName, path.join(installationPath, binName));
    callback(null);
  });
}

function getDownloadData() {
  if (!(process.arch in ARCH_MAPPING)) {
    console.error("Installation is not supported for this architecture: " + process.arch);
    return;
  }

  if (!(process.platform in PLATFORM_MAPPING)) {
    console.error("Installation is not supported for this platform: " + process.platform);
    return;
  }

  var binName = "clonr";
  var binPath = "./bin";
  var url = "https://github.com/oledakotajoe/clonr/releases/download/v{{version}}/clonr_{{version}}_{{platform}}_{{arch}}.tar.gz";
  var version = VERSION;
  if (version[0] === 'v') version = version.substr(1); // strip the 'v' if necessary v0.0.1 => 0.0.1
  // Binary name on Windows has .exe suffix

  if (process.platform === "win32") {
    binName += ".exe";
  } // Interpolate variables in URL, if necessary


  url = url.replace(/{{arch}}/g, ARCH_MAPPING[process.arch]);
  url = url.replace(/{{platform}}/g, PLATFORM_MAPPING[process.platform]);
  url = url.replace(/{{version}}/g, version);
  url = url.replace(/{{bin_name}}/g, binName);
  return {
    binName: binName,
    binPath: binPath,
    url: url,
    version: version
  };
}

var INVALID_INPUT = "Invalid inputs";

function install(callback) {
  var opts = getDownloadData();
  if (!opts) return callback(INVALID_INPUT);
  mkdirp.sync(opts.binPath);
  var ungz = zlib.createGunzip();
  var untar = tar.extract({
    path: opts.binPath
  });
  ungz.on('error', callback);
  untar.on('error', callback); // First we will Un-GZip, then we will untar. So once untar is completed,
  // binary is downloaded into `binPath`. Verify the binary and call it good

  untar.on('end', verifyAndPlaceBinary.bind(null, opts.binName, opts.binPath, callback));
  console.log("Downloading from URL: " + opts.url);
  var req = request({
    uri: opts.url
  });
  req.on('error', callback.bind(null, "Error downloading from URL: " + opts.url));
  req.on('response', function (res) {
    if (res.statusCode !== 200) return callback("Error downloading binary. HTTP Status Code: " + res.statusCode);
    req.pipe(ungz).pipe(untar);
  });
}

function uninstall(callback) {
  var opts = getDownloadData();
  getNpmBinLocation(function (err, installationPath) {
    if (err) callback("Error finding binary installation directory");

    try {
      fs.unlinkSync(path.join(installationPath, opts.binName));
    } catch (ex) {// Ignore errors when deleting the file.
    }

    return callback(null);
  });
} // Parse command line arguments and call the right method


var actions = {
  "install": install,
  "uninstall": uninstall
};
var argv = process.argv;

if (argv && argv.length > 2) {
  var cmd = process.argv[2];

  if (!actions[cmd]) {
    console.log("Invalid command to go-clonr. `install` and `uninstall` are the only supported commands");
    process.exit(1);
  }

  actions[cmd](function (err) {
    if (err) {
      console.error(err);
      process.exit(1);
    } else {
      process.exit(0);
    }
  });
}