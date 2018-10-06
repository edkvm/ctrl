const readline = require('readline');
const net = require('net');
const fs = require('fs');

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

function pipeIPC(fd) {
    return new Promise(resolve => {
        client = net.createConnection(fd)
            .on('connect', () => {
                client.write("__connected");
            })
            .on('error', (err) => {
                console.error(err);
            })
            .on('data', (data) => {
                resolve(data.toString());
            })
    })

}
function main() {

    let args = process.argv.slice(2);
    const executor = require(args[1]);

    // pipeStdin().then(input => {
    //     executor.handler(input, () => {
    //         console.log('Done');
    //     })
    // });
    let fd = args[0];
    console.log(fd)
    pipeIPC(fd).then(input => {
        executor.handler(input, () => {
            console.log('Done');
        })
    });
}

main();

