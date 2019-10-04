const http = require('http');

function handler(req, res) {
    let buf = null;

    // listen for incoming data
    req.on('data', data => {
        console.log(data);
        if (buf === null) {
            buf = data;
        } else {
            buf = buf + data;
        }
    });

    req.on('end', () => {
        console.log(data);
    })
}

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
            arguments[0] = `__2|${arguments[0]}`
            orig.apply(console, arguments)
        } finally {
            process.stderr = tmp
        }
    }
})()

module.exports.run = (handler, handlerName) => {

    let args = process.argv.slice(2);
    let fd = args[0];
    pipe = IPCServer(fd)
    pipe.read().then(raw => {
        try {
            const fn = handler[handlerName];
            let input = JSON.parse(raw)
            fn(input.params, input.ctx)
            .then(data => {
                pipe.write(data)
                pipe.end()
            })
            .catch(error => {
                pipe.write(error)
                pipe.end()

            });
        } catch (err) {
            syserr(`${err} input: ${raw}`);
            pipe.end()
        }
    });
}



