const readline = require('readline');
const net = require('net');
const fs = require('fs');

syslog = (function() {
    let orig = console.log
    return function() {
        let tmp = process.stdout
        try {
            arguments[0] = `__1|${arguments[0]}`
            orig.apply(console, arguments)
        } finally {
            process.stdout = tmp
        }
    };
})();

syserr = (function() {
    let orig = console.error
    return function() {
        let tmp = process.stderr
        try {
            arguments[0] = `__2|${arguments[0]}_`
            orig.apply(console, arguments)
        } finally {
            process.stderr = tmp
        }
    }
})()

module.exports.run = (handler, handlerName) => {
    let args = process.argv.slice(2);
    let raw = args[0];
    try {
        syslog(args[0])
        const fn = handler[handlerName];
        let input = JSON.parse(raw)
        fn(input.params, input.ctx).then((data, err) => {
            if (err !== undefined && err !== null) {
                syserr(`${err}`)
            } else {
                syslog(data)
            }
        })
    } catch (err) {
        syserr(`${err} input: ${raw}`);

    }
}


