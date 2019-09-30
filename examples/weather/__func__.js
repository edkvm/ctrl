const http = require('http');

/**
 * Weather Action returns the weather form the specified city
 * @param {string} city any city
 * @returns {string} The temperature in that city
 */
module.exports.action = (params, ctx, callback) => {

    // ENV
    let apiKey = ctx.config.apiKey;
    let unitSys = ctx.config.unitSys;

    // Parse Event
    let city = params.city;


    let url = buildUrl(city, unitSys, apiKey);
    getApi(url)
        .then(function(data) {
            callback(`It's ${data.main.temp} degrees in ${data.name}!`, null);
        })
        .catch(function(err) {
            callback(null, err);
        });
};

function buildUrl(city, unitSys, apiKey) {
    return `http://api.openweathermap.org/data/2.5/weather?q=${city}&units=${unitSys}&appid=${apiKey}`;
}

function getApi(url) {
    return new Promise(function(resolve, reject) {
        http.get(url, (res) => {
            const { statusCode } = res;
            const contentType = res.headers['content-type'];

            let error;
            if (statusCode !== 200) {
                error = new Error('request failed.\n' +
                    `status code: ${statusCode}`);
            } else if (!/^application\/json/.test(contentType)) {
                error = new Error('invalid content-type.\n' +
                    `expected application/json but received ${contentType}`);
            }

            if (error) {
                reject(error);
            }

            res.setEncoding('utf8');

            let rawData = '';
            res.on('data', (chunk) => { rawData += chunk; });

            res.on('end', () => {
                try {
                    let parsedData = JSON.parse(rawData);
                    resolve(parsedData);
                } catch (err) {
                    reject(err);
                }

            })

        }).on('error', (err) => {
            reject(err);
        });
    });
}

