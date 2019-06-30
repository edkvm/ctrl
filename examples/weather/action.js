const http = require('http');
const apiKey = 'ba9f19affad8428980ee3a66462295a9';

exports.main = (event, callback) => {

	let city = event.city || 'New York';
	let url = `http://api.openweathermap.org/data/2.5/weather?q=${city}&units=imperial&appid=${apiKey}`

	http.get(url, (res) => {
		let body = '';

		res.on('data', (chunk) => {
			body += chunk;
		})

		res.on('end', () => {
			if(res.statusCode === 200){
				let weather = JSON.parse(body)
				let answer = `It is ${weather.main.temp} degrees in ${weather.name}!`;
				callback(null, answer);
			} else {
				console.error(err)
				callback(err, '')
			}

		})

	});

};

