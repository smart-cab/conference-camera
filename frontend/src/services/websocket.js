const os = require('os');

const networkInterfaces = os.networkInterfaces();
// const ip = networkInterfaces['eth0'][0]['address']

console.log(networkInterfaces);