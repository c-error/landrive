const ENCODE = {
    '#': '<0>',
    '%': '<1>',
    '&': '<2>',
    '+': '<3>',
    ';': '<4>'
};

const DECODE = {
    '<0>': '#',
    '<1>': '%',
    '<2>': '&',
    '<3>': '+',
    '<4>': ';'
};

const msg_beacon = document.getElementById('end_beacon');

let chat_input = document.getElementById('user_input');
let chat_name = document.getElementById('user_name');

let chat_count = 1;
let loop_ms = 500;

function url_encode(URL) { return URL.replace(/[#%&+;]/g, i => ENCODE[i]); }
function url_decode(URL) { return URL.replace(/<[0-4]>/g, i => DECODE[i]); }
// let SRV_MSG = document.getElementById('server_msg');

function runSyncChat() {
    sync_chat();
}
let syncInterval = setInterval(runSyncChat, loop_ms);

function sync_chat() {
    
    const xhr = new XMLHttpRequest();
    
    xhr.onload = function() {
        
        const svr_res = xhr.responseText;
        if (svr_res != "x") {

            loop_ms = 10;
            clearInterval(syncInterval);
            syncInterval = setInterval(runSyncChat, loop_ms);

            document.getElementById('server_msg').innerHTML += svr_res;
            msg_beacon.focus();
            chat_input.focus();

            chat_count++;



        } else {

            loop_ms = 500;
            clearInterval(syncInterval);
            syncInterval = setInterval(runSyncChat, loop_ms);
        }
        // console.log(svr_res);
    };
    
    xhr.onerror = function() {
        console.log('RES-ERROR!');
    };

    xhr.open('POST', '/chat?req=sync&no='+chat_count, true);
    xhr.send();
}


function send_chat() {

    const user_msg = chat_input.value;
    const user_name = chat_name.value;

    if (user_msg != "" & user_name != "") {

        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/chat?req=add&name='+url_encode(user_name)+'&data='+url_encode(user_msg), true);
        xhr.send();
    }
}
