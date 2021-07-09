/* globals gauge*/
'use strict';
const assert = require('assert');
const axios = require('axios');
const moment = require('moment');

step('Generate random username', async () => {
    const username = moment().format('YYMMDDhhmmssSSS');
    gauge.dataStore.scenarioStore.put('random_username', username);
});

step('Register random username with <password>', async (password) => {
    const username = gauge.dataStore.scenarioStore.get('random_username');
    await axios.post(
        `${process.env.API_BASE}/v1/auth/login`,
        {
            'username': username,
            'password': password,
        })
        .then(res => {
            assert.strictEqual(res.status, 201, 'status should be 201');
            assert.strictEqual(res.data.hasOwnProperty('token'), true, 'register user should return token');
            gauge.dataStore.scenarioStore.put('token', res.data.token);
        });
});

step('Login random username with <password>', async (password) => {
    const username = gauge.dataStore.scenarioStore.get('random_username');
    await axios.post(
        `${process.env.API_BASE}/v1/auth/login`,
        {
            'username': username,
            'password': password,
        })
        .then(res => {
            assert.strictEqual(res.status, 200, 'status should be 200');
            assert.strictEqual(res.data.hasOwnProperty('token'), true, 'register user should return token');
            gauge.dataStore.scenarioStore.put('token', res.data.token);
        });
});
