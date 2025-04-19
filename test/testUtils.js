// testUtils.js - Common test utilities for all test suites
const fs = require('fs');
const path = require('path');

module.exports = {
  sleep: (ms) => new Promise(resolve => setTimeout(resolve, ms)),
  mockEnv: (vars) => {
    const oldEnv = { ...process.env };
    Object.assign(process.env, vars);
    return () => { process.env = oldEnv; };
  },
  // File management and cleanup utilities
  testFileUtils: {
    ensureDir: (dir) => {
      if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
    },
    removeFile: (file) => {
      if (fs.existsSync(file)) fs.unlinkSync(file);
    },
    cleanupDir: (dir) => {
      if (fs.existsSync(dir)) {
        fs.readdirSync(dir).forEach(f => {
          const filePath = path.join(dir, f);
          if (fs.lstatSync(filePath).isDirectory()) {
            module.exports.testFileUtils.cleanupDir(filePath);
            fs.rmdirSync(filePath);
          } else {
            fs.unlinkSync(filePath);
          }
        });
      }
    }
  },
  // Common assertion utilities
  assertUtils: {
    isDefined: (val, msg) => {
      if (val === undefined) throw new Error(msg || 'Value is undefined');
    },
    isNotNull: (val, msg) => {
      if (val === null) throw new Error(msg || 'Value is null');
    },
    deepEqual: (a, b, msg) => {
      if (JSON.stringify(a) !== JSON.stringify(b)) throw new Error(msg || 'Values are not deeply equal');
    }
  },
  // Test data generator
  testDataGen: {
    randomString: (len = 8) => Math.random().toString(36).substring(2, 2 + len),
    randomInt: (min = 0, max = 100) => Math.floor(Math.random() * (max - min + 1)) + min,
  },
  // Mocking interfaces
  mockFs: {
    withMockedFs: (mockImpl, fn) => {
      const origFs = { ...fs };
      Object.assign(fs, mockImpl);
      try { fn(); } finally { Object.assign(fs, origFs); }
    }
  },
  // I/O capture utilities
  captureStdout: function(fn) {
    const originalWrite = process.stdout.write;
    let output = '';
    process.stdout.write = (chunk, encoding, cb) => {
      output += chunk;
      if (cb) cb();
      return true;
    };
    try { fn(); } finally { process.stdout.write = originalWrite; }
    return output;
  },
  captureStderr: function(fn) {
    const originalWrite = process.stderr.write;
    let output = '';
    process.stderr.write = (chunk, encoding, cb) => {
      output += chunk;
      if (cb) cb();
      return true;
    };
    try { fn(); } finally { process.stderr.write = originalWrite; }
    return output;
  },
};
