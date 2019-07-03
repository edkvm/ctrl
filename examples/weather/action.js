const http = require('http');



/**
 * Weather Action returns the weather form the specified city
 * @param {string} city any city
 *
 * @return {string} The temperature in that city
 */
exports.main = (event, callback) => {

	// Parse Event
	let city = event.params.city;
	let apiKey = event.__action_config.apiKey;
	let unitSys = event.__action_config.unitSys;

	let url = buildUrl(city, unitSys, apiKey);
	getWeather(url, callback);
};

function buildUrl(city, unitSys, apiKey) {
	return `http://api.openweathermap.org/data/2.5/weather?q=${city}&units=${unitSys}&appid=${apiKey}`;
}

function getWeather(url, callback) {
	http.get(url, (res) => {
		let body = '';

		res.on('data', (chunk) => {
			body += chunk;
		})

		res.on('end', () => {
			if(res.statusCode === 200){
				let weather = JSON.parse(body);
				let answer = `It is ${weather.main.temp} degrees in ${weather.name}!`;
				callback(null, answer);
			} else {
				console.error(err);
				callback(err, '');
			}

		})

	});
}

