"use strict";

// loading modules
const express = require('express');
const cors = require('cors');
const morgan = require('morgan');

// require func can also load json files
const zips = require('./zips.json');
// console.log('loaded %d zips', zips.length);

// let zipCityIndex;
// for (var i = 0; i < zips.length; i++) {
//     var key = zips[i].city;
//     var value = zips[i].zip;
//     zipCityIndex[key.toLowerCase()] = value;
//     console.log(zipCityIndex[i]) 
// }

const zipCityIndex = zips.reduce((index, record) => {
    let cityLower = record.city.toLowerCase();
    let zipsForCity = index[cityLower];
    if (!zipsForCity) {
        index[cityLower] = zipsForCity = [];
    }
    zipsForCity.push(record);
    return index;
}, {});
// console.log('there are %d zips in Seattle', zipCityIndex.seattle.length);

const app = express();

const port = process.env.PORT || 80; // if not set, default to 80
const host = process.env.HOST || '';

// morgan('dev') returns a middleware function
// app.use says: use that middleware function for EVERY request
app.use(morgan('dev'));
app.use(cors());

app.get('/zips/city/:cityName', (req, res) => {
    let zipsForCity = zipCityIndex[req.params.cityName.toLowerCase()];
    if (!zipsForCity) {
        res.status(404).send('invalid city name');
    } else {
        res.json(zipsForCity);
    }
});

// request object first, respond object second (opposite of golang)
app.get('/hello/:name', (req, res) => {
    res.send(`Hello ${req.params.name}!`); // shows on screen at the provided path
});

// add handler

app.listen(port, host, () => {
    console.log(`server is listening at http://${host}:${port}...`);
});
