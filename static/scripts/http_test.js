// Function to handle the file click and process the download/upload
async function fileClick(filename, uploadType = 'normal') {
    const dlurl = `http://localhost:8080/static/downloads/${filename}`;
    const posturl = "http://localhost:8080/upload";
    
    try {
        let result = await downloadAndPostFile(dlurl, posturl, uploadType);
        console.log("Result from server:", result);
    } catch (error) {
        console.error("Error during file download and posting:", error);
    }
}

// Function to download and post the file (normal or secure)
async function downloadAndPostFile(fileUrl, postUrl, uploadType) {
    const filename = fileUrl.substring(fileUrl.lastIndexOf('/') + 1);
    const status_span = document.getElementById(`${filename}-status`);

    try {
        // Download the file as a Blob
        const blob = await downloadFileToBlob(fileUrl);

        if (uploadType === 'normal') {
            await normalUpload(blob, filename, postUrl, status_span);
        } else if (uploadType === 'secure') {
            await secureUpload(blob, filename, status_span);
        }

        console.log('File posted successfully');
    } catch (error) {
        if (status_span) {
            status_span.textContent = `POST not successful - ${error.message}`;
        }
        console.error('Error downloading or posting file:', error);
        throw error;
    }
}

// Function to upload file normally (to your server)
async function normalUpload(file, filename, postUrl, status_span) {
    const formData = new FormData();
    formData.append('file', file, filename);

    try {
        const postResponse = await fetch(postUrl, {
            method: 'POST',
            body: formData,
        });

        if (!postResponse.ok) {
            const errorText = await postResponse.text();
            if (status_span) {
                status_span.textContent = `POST failed - Status: ${postResponse.status}`;
            }
            throw new Error(`POST error: ${postResponse.status}, Message: ${errorText}`);
        }

        if (status_span) {
            status_span.textContent = "POST successful";
        }

        return await postResponse.json();
    } catch (error) {
        if (status_span) {
            status_span.textContent = `POST not successful - ${error.message}`;
        }
        console.error('Error during normal upload:', error);
        throw error;
    }
}

// Function to securely upload the file to S3 using pre-signed URL
async function secureUpload(file, filename, status_span) {
    try {
        const response = await fetch('/generate-s3-url', { method: 'GET' });

        if (!response.ok) {
            throw new Error('Failed to retrieve S3 URL.');
        }

        const data = await response.json();
        const preSignedUrl = data.url;
        const formData = new FormData();
        formData.append('file', file, filename);

        const uploadResponse = await fetch(preSignedUrl, {
            method: 'POST',
            body: formData
        });

        if (!uploadResponse.ok) {
            throw new Error('Failed to upload file to S3.');
        }

        if (status_span) {
            status_span.textContent = "File securely uploaded to S3!";
        }
    } catch (error) {
        if (status_span) {
            status_span.textContent = `Error during secure upload: ${error.message}`;
        }
        console.error('Error during secure upload:', error);
    }
}
