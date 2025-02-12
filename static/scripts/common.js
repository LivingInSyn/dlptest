/*common.js*/

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

async function normalUpload(file) {
    const message = document.getElementById("message");
    const progressBar = document.getElementById("uploadProgress");

    const formData = new FormData();
    formData.append("file", file);

    try {
        const response = await fetch("/upload", {
            method: "POST",
            body: formData
        });

        if (!response.ok) throw new Error(`Upload failed: ${response.statusText}`);

        message.textContent = "File uploaded successfully!";
        progressBar.value = 100;
    } catch (error) {
        message.textContent = "Error uploading file: " + error.message;
    }
}

async function secureUpload(file) {
    const message = document.getElementById("message");
    const progressBar = document.getElementById("uploadProgress");

    try {
        const response = await fetch("/generate-s3-token");
        if (!response.ok) throw new Error("Failed to get S3 URL");

        const data = await response.json();
        const uploadUrl = data.uploadUrl;

        const uploadResponse = await fetch(uploadUrl, {
            method: "PUT",
            body: file,
            headers: { "Content-Type": file.type }
        });

        if (!uploadResponse.ok) throw new Error("S3 Upload Failed");

        message.textContent = "Secure file uploaded to S3 successfully!";
        progressBar.value = 100;
    } catch (error) {
        message.textContent = "Error in secure upload: " + error.message;
    }
}
