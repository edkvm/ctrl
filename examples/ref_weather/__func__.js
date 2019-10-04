const http = require('http');

/**
 * Weather Action returns the weather form the specified city
 * @param {string} city any city
 * @returns {string} The temperature in that city
 */
module.exports.action = async (params, ctx) => {

    // ENV
    let apiKey = process.env.API_KEY;
    let unitSys = params.unit || process.env.DEFAULT_UNIT_SYS;

    // Parse Event
    let city = params.city;
    let language = params.city;


    let url = buildUrl(city, unitSys, apiKey);

    const data =
        await getApi(url)
             .catch(err => {
                 return "e"
             });
    return `It's ${data.main.temp} degrees in ${data.name}!`
};

function buildUrl(city, unitSys, apiKey) {
    return `http://api.openweathermap.org/data/2.5/weather?q=${city}&units=${unitSys}&appid=${apiKey}`;
}

function getApi(url) {
    return new Promise((resolve, reject) => {
        http.get(url, (res) => {
            const { statusCode } = res;
            const contentType = res.headers['content-type'];

            let error;
            if (statusCode == 404) {
                console.error(`request failed. status code: ${statusCode}`)
                error = new Error(`location not found`);
            } else if (statusCode !== 200) {
                console.error(`request failed. status code: ${statusCode}`)

            } else if (!/^application\/json/.test(contentType)) {
                error = new Error('invalid content-type.\n' +
                    `expected application/json but received ${contentType}`);
            }

            if (error) {
                reject(error);
            }

            res.setEncoding('utf8');

            let rawData = '';
            res.on('data', (dataChunk) => {
                rawData += dataChunk;
            });

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

