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

// Direct all console.log to stderr
function pipeStdin() {
    let buf = '';

    return new Promise(resolve => {

        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout
        });

        rl.on('line', (line) => {
            buf += line;
        });

        rl.on('close', () => {
            resolve(buf)
        })
    });
}

function IPCServer(fd) {

    srv = net.createConnection(fd)

    srv.on('connect', () => {
        srv.write("op|start");
    })
    srv.on('error', (err) => {
        syserr(`${err}`);
    })

    return {
        write: (data) => {
            return new Promise(resolve => {
                resolve(srv.write(data))
            })
        },
        read: () => {
            return new Promise(resolve => {
                srv.on('data', (data) => {
                    resolve(data.toString());
                })
            })
        },
        end: () => {
            srv.write("op|close")
        }


    }

}

module.exports.run = (handler, handlerName) => {

    let args = process.argv.slice(2);
    let fd = args[0];
    pipe = IPCServer(fd)
    pipe.read().then(raw => {
        try {
            const fn = handler[handlerName];
            let input = JSON.parse(raw)
            fn(input.params, input.ctx, (data, err) => {
                if (err !== undefined && err !== null) {
                    syserr(`${err}`)
                    pipe.end()
                } else {
                    pipe.write(data)
                    pipe.end()
                }

            })
        } catch (err) {
            syserr(`${err} input: ${raw}`);
            pipe.end()
        }
    });
}


