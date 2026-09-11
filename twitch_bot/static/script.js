"use strict";

const client_id = 'wltjlsycrctd59xi1f7xg3xa1qrve1'
const redirect_uri = 'http://localhost:8080/oauth'
const scope = 'channel:bot'


function addBot() {
    window.location.href = encodeURI(`https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=${client_id}&redirect_uri=${redirect_uri}&scope=${scope}`)
}
