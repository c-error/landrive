const SVR_PAR = new URLSearchParams(window.location.search);

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

let check_filter_all = document.getElementById("filter_all");
let check_filter_fo = document.getElementById("filter_fo");
let check_filter_fi = document.getElementById("filter_fi");

let search_name = document.getElementById("filter_name");
let search_size = document.getElementById("filter_size");

let up_data = document.getElementById('uploding');

let upload_count = 0;
let noti_count = 0;

function url_encode(URL) { return URL.replace(/[#%&+;]/g, i => ENCODE[i]); }
function url_decode(URL) { return URL.replace(/<[0-4]>/g, i => DECODE[i]); }

if (SVR_PAR.has("f")) {

    const filter = SVR_PAR.get("f").split(",");

    if (filter[0] === 'fo') check_filter_fo.checked = true;
    else if (filter[0] === 'fi') check_filter_fi.checked = true;
    else check_filter_all.checked = true;

    if (filter.length >= 2) search_name.value = url_decode(filter[1]);
    if (filter.length >= 3) search_size.value = url_decode(filter[2]);
}

document.body.addEventListener('click', e => {

    const filter_all = e.target.closest('#filter_all');
    const filter_fo = e.target.closest('#filter_fo');
    const filter_fi = e.target.closest('#filter_fi');

    const cell = e.target.closest('cl');

    if (cell) {
        const data_type = cell.querySelector('nm a')?.textContent;
        const data_name = cell.querySelector('nm b')?.textContent;

        if (data_type && data_name) {

            if (data_type === 'FO:') {

                let url = '/path?fo=' + SVR_URL;
                if (!url.endsWith('/')) url += '/';
                location.href = url + url_encode(data_name);

            } else if (data_type === 'FI:') {

                let url = '/path?fi=' + SVR_URL;
                if (!url.endsWith('/')) url += '/';
                // url += data_name;
                window.open(url + url_encode(data_name), "_blank");
            }
        }
    }

    if (filter_all) {

        check_filter_fo.checked = false;
        check_filter_fi.checked = false;

    } else if (filter_fo) {

        check_filter_all.checked = false;
        check_filter_fi.checked = false;

    } else if (filter_fi) {

        check_filter_all.checked = false;
        check_filter_fo.checked = false;
    }
});

function svr_search() {

    let url = '/path?fo=' + SVR_URL;
    if (!url.endsWith("/")) url += "/";

    if (check_filter_fo.checked) url += "&f=fo,";
    else if (check_filter_fi.checked) url += "&f=fi,";
    else url += "&f=all,";

    url += url_encode(search_name.value) + ",";
    url += url_encode(search_size.value);
    location.href = url;
}

function search_clr() {

    let url = "/path?fo=" + SVR_URL;
    if (!url.endsWith("/")) url += "/";
    location.href = url;
}



function add_noti(msg) {

    const i = noti_count++;

    document.getElementById('popup').innerHTML += `<a id="noti_cell_id_${i}">💡 ${msg}</a>`;
    setTimeout(() => {
        document.getElementById(`noti_cell_id_${i}`).remove();
    }, 3000);
}





function uploadFiles() {
    const fileInput = document.getElementById('fileInput');
    const files = fileInput.files;
    
    up_data.innerHTML = '';
    
    if (files.length === 0) {
        // alert('Please select files to upload');
        return;
    }
    
    document.getElementById('uploding_shell').style.display = '';

    for (let i = 0; i < files.length; i++) {
        const file = files[i];
        const formData = new FormData();
        formData.append('files', file);
        // formData.append('dir', "__file__");
        
        // Create progress element
        const progressDiv = document.createElement('cell');
        progressDiv.id = "cell_id_"+i;
        progressDiv.innerHTML = '<up-info><a>'+file.name+'</a><b id="upload_per_'+i+'">0%</b></up-info><div class="progress_shell"><div id="upload_width_'+i+'" class="progress"></div></div>';

        up_data.appendChild(progressDiv);
        
        const xhr = new XMLHttpRequest();
        
        // Progress event
        xhr.upload.addEventListener('progress', function(e) {
            if (e.lengthComputable) {
                const percent = ((e.loaded / e.total) * 100).toFixed(3);
                document.getElementById('upload_width_'+i).style.width = percent + '%';
                document.getElementById('upload_per_'+i).innerText = percent + '%';

            }
        });
        
        xhr.onreadystatechange = function() {

            if (xhr.readyState === 4 && xhr.status === 200) {

                upload_count++;
                
                add_noti('DONE:'+file.name);
                document.getElementById('cell_id_'+i).style.display = 'none';
                // console.log('Server response:', xhr.responseText);
            }

            if (upload_count === files.length) {
                
                upload_count = 0;
                document.getElementById('uploding_shell').style.display = 'none';

                add_noti('Upload Complete ...');

                setTimeout(() => {
                    location.reload();
                }, 4000);
            }
            // console.log('Server response: ', upload_count, files.length);
            // }
        };

        const url = '/set?fo=' + SVR_URL;

        xhr.open('POST', url, true);
        xhr.send(formData);
    }
}



