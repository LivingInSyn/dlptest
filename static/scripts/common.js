function uploadFile(uploadType) {
    const fileInput = document.getElementById("fileInput");
    const file = fileInput.files[0];
    const progressBar = document.getElementById("uploadProgress");
    const message = document.getElementById("message");

    if (!file) {
        message.textContent = "Please select a file to upload.";
        return;
    }

    progressBar.style.display = "block";
    progressBar.value = 0;
    
    if (uploadType === "normal") {
        normalUpload(file);
    } else if (uploadType === "secure") {
        secureUpload(file);
    } else {
        message.textContent = "Invalid upload type.";
    }
}

function normalUpload(file) {
    const message = document.getElementById("message");
    const progressBar = document.getElementById("uploadProgress");

    const formData = new FormData();
    formData.append("file", file);

    fetch("/upload", {
        method: "POST",
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        message.textContent = "File uploaded successfully!";
        progressBar.value = 100;
    })
    .catch(error => {
        message.textContent = "Error uploading file: " + error;
    });
}

function secureUpload(file) {
    const message = document.getElementById("message");
    const progressBar = document.getElementById("uploadProgress");

    fetch("/generate-s3-token")
    .then(response => response.json())
    .then(data => {
        const formData = new FormData();
        formData.append("file", file);

        return fetch(data.uploadUrl, {
            method: "PUT",
            body: file,
            headers: { "Content-Type": file.type }
        });
    })
    .then(response => {
        if (response.ok) {
            message.textContent = "Secure file uploaded to S3 successfully!";
            progressBar.value = 100;
        } else {
            throw new Error("S3 Upload Failed");
        }
    })
    .catch(error => {
        message.textContent = "Error in secure upload: " + error;
    });
}
