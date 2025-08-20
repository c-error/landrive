
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

function url_encode(URL) { return URL.replace(/[#%&+;]/g, i => ENCODE[i]); }
function url_decode(URL) { return URL.replace(/<[0-4]>/g, i => DECODE[i]); }

let p = document.getElementById('pin');

if (SVR_PAR.has("p")) p.value = url_decode(SVR_PAR.get("p"));

window.onload = function() {
    p.focus();
};

document.getElementById('pin').addEventListener('keypress', function(event) {
    
    if (event.key === 'Enter') location.href = '/login?pin=' + url_encode(p.value);
});
