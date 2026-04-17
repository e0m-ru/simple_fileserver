const dropZone = document.getElementById('dropZone');
const fileInput = document.getElementById('fileInput');
const uploadStatus = document.getElementById('uploadStatus');
const filesList = document.getElementById('filesList');

// Load files on page load
loadFiles();

// Click to select file
dropZone.addEventListener('click', () => fileInput.click());

// Drag and drop events
['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
    dropZone.addEventListener(eventName, preventDefaults, false);
});

function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
}

['dragenter', 'dragover'].forEach(eventName => {
    dropZone.addEventListener(eventName, highlight, false);
});

['dragleave', 'drop'].forEach(eventName => {
    dropZone.addEventListener(eventName, unhighlight, false);
});

function highlight(e) {
    dropZone.classList.add('active');
}

function unhighlight(e) {
    dropZone.classList.remove('active');
}

dropZone.addEventListener('drop', handleDrop, false);

function handleDrop(e) {
    const dt = e.dataTransfer;
    const files = dt.files;
    fileInput.files = files;
    uploadFile(files[0]);
}

fileInput.addEventListener('change', (e) => {
    if (e.target.files.length > 0) {
        uploadFile(e.target.files[0]);
    }
});

function uploadFile(file) {
    const formData = new FormData();
    formData.append('file', file);

    uploadStatus.innerHTML = '<div class="alert alert-info">Uploading...</div>';

    fetch('/api/upload', {
        method: 'POST',
        body: formData
    })
        .then(response => response.text())
        .then(data => {
            uploadStatus.innerHTML = '<div class="alert alert-success">✓ ' + data + '</div>';
            loadFiles();
            setTimeout(() => {
                uploadStatus.innerHTML = '';
            }, 3000);
        })
        .catch(error => {
            uploadStatus.innerHTML = '<div class="alert alert-danger">✗ Upload failed: ' + error + '</div>';
        });
}

function loadFiles() {
    fetch('/api/files')
        .then(response => response.json())
        .then(files => {
            if (!files || files.length === 0) {
                filesList.innerHTML = '<div class="list-group-item text-muted">No files uploaded yet</div>';
                return;
            }

            filesList.innerHTML = '';
            files.forEach(file => {
                const size = formatFileSize(file.size);
                const item = document.createElement('div');
                item.className = 'list-group-item d-flex justify-content-between align-items-center';
                item.innerHTML = `
                    <a href="/uploads/${encodeURIComponent(file.name)}" download="${file.name}" class="text-decoration-none flex-grow-1">
                        <span>📄 ${file.name}</span>
                    </a>
                    <div class="d-flex gap-2">
                        <span class="badge bg-secondary">${size}</span>
                        <button class="btn btn-sm btn-danger" onclick="deleteFile('${file.name}')">
                            🗑️
                        </button>
                    </div>
                `;
                filesList.appendChild(item);
            });
        })
        .catch(error => {
            filesList.innerHTML = '<div class="list-group-item text-danger">Failed to load files</div>';
            console.error(error);
        });
}

function deleteFile(filename) {
    let formData = new FormData();
    formData.append('filename', filename);
    fetch('/api/delete', {
        method: 'POST',
        body: formData
    })
        .then(response => {
            if (response.ok) {
                loadFiles();
                uploadStatus.innerHTML = '<div class="alert alert-success">✓ File deleted successfully</div>';
                setTimeout(() => {
                    uploadStatus.innerHTML = '';
                }, 2000);
            } else {
                uploadStatus.innerHTML = '<div class="alert alert-danger">✗ Failed to delete file ' + response.statusText + '</div>';
            }
        })
        .catch(error => {
            uploadStatus.innerHTML = '<div class="alert alert-danger">✗ Delete failed: ' + error + '</div>';
            console.error(error);
        });
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}