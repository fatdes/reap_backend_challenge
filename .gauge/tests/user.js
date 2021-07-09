/* globals gauge*/
'use strict';
const assert = require('assert');
const axios = require('axios');
const FormData = require('form-data');
const fs = require('fs');
const util = require('util');

step('Create new post <filename> <file> and <description>', async (filename, file, description) => {
    let form = new FormData();
    form.append('image_name', filename)
    form.append('description', description);

    let fileStream = fs.createReadStream(file, { encoding: 'binary' });
    form.append('image', fileStream);

    const formHeaders = form.getHeaders(
        {
            'Authorization': `Bearer ${gauge.dataStore.scenarioStore.get('token')}`
        }
    );

    await axios.post(
        `${process.env.API_BASE}/v1/user/post`,
        form,
        {
            headers: formHeaders,
        })
        .then(res => {
            assert.strictEqual(res.status, 200, 'status should be 200');
            assert.strictEqual(res.data.hasOwnProperty('url'), true, 'create post succeed should return "url"');
        });
});

step('List post posted by myself (random user) order by created at <sortOrder> <table>', async (sortOrder, table) => {
    let username = gauge.dataStore.scenarioStore.get('random_username');

    await axios.get(
        `${process.env.API_BASE}/v1/user/post?username=${username}&created_at_sort_order=${sortOrder}`,
        {
            headers: {
                'Authorization': `Bearer ${gauge.dataStore.scenarioStore.get('token')}`
            }
        })
        .then(res => {
            assert.strictEqual(res.status, 200, 'status should be 200');
            assert.strictEqual(res.data.hasOwnProperty('post_list'), true, 'list post should have "post_list"');
            assert.strictEqual(Array.isArray(res.data.post_list), true, 'list post should be array');
            assert.strictEqual(res.data.post_list.length, table.rows.length, `list post should have ${table.rows.length} number of post, but is ${res.data.post_list.length}`);
            var i = 0;
            table.entries(function (entry) {
                assert.strictEqual(entry['description'], res.data.post_list[i].description);
                i++;
            });
        });
});
