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
const loop_ms = 500;

let chat_input = document.getElementById('user_input');
let chat_name = document.getElementById('user_name');

let chat_count = 1;

function url_encode(URL) { return URL.replace(/[#%&+;]/g, i => ENCODE[i]); }
function url_decode(URL) { return URL.replace(/<[0-4]>/g, i => DECODE[i]); }

function send_chat() {

    const user_msg = chat_input.value;
    const user_name = chat_name.value;

    if (user_msg != "" & user_name != "") {

        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/chat?req=add&name='+url_encode(user_name)+'&data='+url_encode(user_msg), true);
        xhr.send();
    }
}

function sync_chat() {
    
    const xhr = new XMLHttpRequest();
    
    xhr.onload = function() {
        
        const svr_res = xhr.responseText;
        if (svr_res != "x") {

            document.getElementById('server_msg').innerHTML += url_decode(svr_res);
            msg_beacon.focus();
            chat_input.focus();

            chat_count++;
        }
    };
    
    xhr.onerror = function() {
        console.log('RES-ERROR!');
    };

    xhr.open('POST', '/chat?req=sync&no='+chat_count, true);
    xhr.send();
}
setInterval(sync_chat, loop_ms);


